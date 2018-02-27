// Copyright (c) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gnmi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	pb "github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc/codes"
)

// Get sents a GetRequest to the given client.
func Get(ctx context.Context, client pb.GNMIClient, paths [][]string) error {
	req, err := NewGetRequest(paths)
	if err != nil {
		return err
	}
	resp, err := client.Get(ctx, req)
	if err != nil {
		return err
	}
	for _, notif := range resp.Notification {
		for _, update := range notif.Update {
			fmt.Printf("%s:\n", StrPath(update.Path))
			fmt.Println(StrUpdateVal(update))
		}
	}
	return nil
}

// Capabilities retuns the capabilities of the client.
func Capabilities(ctx context.Context, client pb.GNMIClient) error {
	resp, err := client.Capabilities(ctx, &pb.CapabilityRequest{})
	if err != nil {
		return err
	}
	fmt.Printf("Version: %s\n", resp.GNMIVersion)
	for _, mod := range resp.SupportedModels {
		fmt.Printf("SupportedModel: %s\n", mod)
	}
	for _, enc := range resp.SupportedEncodings {
		fmt.Printf("SupportedEncoding: %s\n", enc)
	}
	return nil
}

// val may be a path to a file or it may be json. First see if it is a
// file, if so return its contents, otherwise return val
func extractJSON(val string) []byte {
	if jsonBytes, err := ioutil.ReadFile(val); err == nil {
		return jsonBytes
	}
	// Best effort check if the value might a string literal, in which
	// case wrap it in quotes. This is to allow a user to do:
	//   gnmi update ../hostname host1234
	//   gnmi update ../description 'This is a description'
	// instead of forcing them to quote the string:
	//   gnmi update ../hostname '"host1234"'
	//   gnmi update ../description '"This is a description"'
	maybeUnquotedStringLiteral := func(s string) bool {
		if s == "true" || s == "false" || s == "null" || // JSON reserved words
			strings.ContainsAny(s, `"'{}[]`) { // Already quoted or is a JSON object or array
			return false
		} else if _, err := strconv.ParseInt(s, 0, 32); err == nil {
			// Integer. Using byte size of 32 because larger integer
			// types are supposed to be sent as strings in JSON.
			return false
		} else if _, err := strconv.ParseFloat(s, 64); err == nil {
			// Float
			return false
		}

		return true
	}
	if maybeUnquotedStringLiteral(val) {
		out := make([]byte, len(val)+2)
		out[0] = '"'
		copy(out[1:], val)
		out[len(out)-1] = '"'
		return out
	}
	return []byte(val)
}

// StrUpdateVal will return a string representing the value within the supplied update
func StrUpdateVal(u *pb.Update) string {
	if u.Value != nil {
		// Backwards compatibility with pre-v0.4 gnmi
		switch u.Value.Type {
		case pb.Encoding_JSON, pb.Encoding_JSON_IETF:
			return strJSON(u.Value.Value)
		case pb.Encoding_BYTES, pb.Encoding_PROTO:
			return base64.StdEncoding.EncodeToString(u.Value.Value)
		case pb.Encoding_ASCII:
			return string(u.Value.Value)
		default:
			return string(u.Value.Value)
		}
	}
	return StrVal(u.Val)
}

// StrVal will return a string representing the supplied value
func StrVal(val *pb.TypedValue) string {
	switch v := val.GetValue().(type) {
	case *pb.TypedValue_StringVal:
		return v.StringVal
	case *pb.TypedValue_JsonIetfVal:
		return strJSON(v.JsonIetfVal)
	case *pb.TypedValue_JsonVal:
		return strJSON(v.JsonVal)
	case *pb.TypedValue_IntVal:
		return strconv.FormatInt(v.IntVal, 10)
	case *pb.TypedValue_UintVal:
		return strconv.FormatUint(v.UintVal, 10)
	case *pb.TypedValue_BoolVal:
		return strconv.FormatBool(v.BoolVal)
	case *pb.TypedValue_BytesVal:
		return base64.StdEncoding.EncodeToString(v.BytesVal)
	case *pb.TypedValue_DecimalVal:
		return strDecimal64(v.DecimalVal)
	case *pb.TypedValue_FloatVal:
		return strconv.FormatFloat(float64(v.FloatVal), 'g', -1, 32)
	case *pb.TypedValue_LeaflistVal:
		return strLeaflist(v.LeaflistVal)
	case *pb.TypedValue_AsciiVal:
		return v.AsciiVal
	case *pb.TypedValue_AnyVal:
		return v.AnyVal.String()
	default:
		panic(v)
	}
}

func strJSON(inJSON []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, inJSON, "", "  ")
	if err != nil {
		return fmt.Sprintf("(error unmarshalling json: %s)\n", err) + string(inJSON)
	}
	return out.String()
}

func strDecimal64(d *pb.Decimal64) string {
	var i, frac int64
	if d.Precision > 0 {
		div := int64(10)
		it := d.Precision - 1
		for it > 0 {
			div *= 10
			it--
		}
		i = d.Digits / div
		frac = d.Digits % div
	} else {
		i = d.Digits
	}
	if frac < 0 {
		frac = -frac
	}
	return fmt.Sprintf("%d.%d", i, frac)
}

// strLeafList builds a human-readable form of a leaf-list. e.g. [1, 2, 3] or [a, b, c]
func strLeaflist(v *pb.ScalarArray) string {
	var buf bytes.Buffer
	buf.WriteByte('[')

	for i, elm := range v.Element {
		buf.WriteString(StrVal(elm))
		if i < len(v.Element)-1 {
			buf.WriteString(", ")
		}
	}

	buf.WriteByte(']')
	return buf.String()
}

func update(p *pb.Path, val string) *pb.Update {
	var v *pb.TypedValue
	switch p.Origin {
	case "":
		v = &pb.TypedValue{
			Value: &pb.TypedValue_JsonIetfVal{JsonIetfVal: extractJSON(val)}}
	case "cli":
		v = &pb.TypedValue{
			Value: &pb.TypedValue_AsciiVal{AsciiVal: val}}
	default:
		panic(fmt.Errorf("unexpected origin: %q", p.Origin))
	}

	return &pb.Update{Path: p, Val: v}
}

// Operation describes an gNMI operation.
type Operation struct {
	Type string
	Path []string
	Val  string
}

func newSetRequest(setOps []*Operation) (*pb.SetRequest, error) {
	req := &pb.SetRequest{}
	for _, op := range setOps {
		p, err := ParseGNMIElements(op.Path)
		if err != nil {
			return nil, err
		}

		switch op.Type {
		case "delete":
			req.Delete = append(req.Delete, p)
		case "update":
			req.Update = append(req.Update, update(p, op.Val))
		case "replace":
			req.Replace = append(req.Replace, update(p, op.Val))
		}
	}
	return req, nil
}

// Set sends a SetRequest to the given client.
func Set(ctx context.Context, client pb.GNMIClient, setOps []*Operation) error {
	req, err := newSetRequest(setOps)
	if err != nil {
		return err
	}
	resp, err := client.Set(ctx, req)
	if err != nil {
		return err
	}
	if resp.Message != nil && codes.Code(resp.Message.Code) != codes.OK {
		return errors.New(resp.Message.Message)
	}
	// TODO: Iterate over SetResponse.Response for more detailed error message?

	return nil
}

// Subscribe sends a SubscribeRequest to the given client.
func Subscribe(ctx context.Context, client pb.GNMIClient, paths [][]string,
	respChan chan<- *pb.SubscribeResponse, errChan chan<- error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	stream, err := client.Subscribe(ctx)
	if err != nil {
		errChan <- err
		return
	}
	req, err := NewSubscribeRequest(paths)
	if err != nil {
		errChan <- err
		return
	}
	if err := stream.Send(req); err != nil {
		errChan <- err
		return
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return
			}
			errChan <- err
			return
		}
		respChan <- resp
	}
}

// LogSubscribeResponse logs update responses to stderr.
func LogSubscribeResponse(response *pb.SubscribeResponse) error {
	switch resp := response.Response.(type) {
	case *pb.SubscribeResponse_Error:
		return errors.New(resp.Error.Message)
	case *pb.SubscribeResponse_SyncResponse:
		if !resp.SyncResponse {
			return errors.New("initial sync failed")
		}
	case *pb.SubscribeResponse_Update:
		prefix := StrPath(resp.Update.Prefix)
		for _, update := range resp.Update.Update {
			fmt.Printf("%s = %s\n", path.Join(prefix, StrPath(update.Path)),
				StrUpdateVal(update))
		}
	}
	return nil
}
