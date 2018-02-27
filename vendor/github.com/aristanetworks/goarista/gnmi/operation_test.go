// Copyright (c) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gnmi

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aristanetworks/goarista/test"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"

	pb "github.com/openconfig/gnmi/proto/gnmi"
)

func TestNewSetRequest(t *testing.T) {
	pathFoo := &pb.Path{
		Element: []string{"foo"},
		Elem:    []*pb.PathElem{{Name: "foo"}},
	}
	pathCli := &pb.Path{
		Origin: "cli",
	}

	testCases := map[string]struct {
		setOps []*Operation
		exp    pb.SetRequest
	}{
		"delete": {
			setOps: []*Operation{{Type: "delete", Path: []string{"foo"}}},
			exp:    pb.SetRequest{Delete: []*pb.Path{pathFoo}},
		},
		"update": {
			setOps: []*Operation{{Type: "update", Path: []string{"foo"}, Val: "true"}},
			exp: pb.SetRequest{
				Update: []*pb.Update{{
					Path: pathFoo,
					Val: &pb.TypedValue{
						Value: &pb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("true")}},
				}},
			},
		},
		"replace": {
			setOps: []*Operation{{Type: "replace", Path: []string{"foo"}, Val: "true"}},
			exp: pb.SetRequest{
				Replace: []*pb.Update{{
					Path: pathFoo,
					Val: &pb.TypedValue{
						Value: &pb.TypedValue_JsonIetfVal{JsonIetfVal: []byte("true")}},
				}},
			},
		},
		"cli-replace": {
			setOps: []*Operation{{Type: "replace", Path: []string{"cli"},
				Val: "hostname foo\nip routing"}},
			exp: pb.SetRequest{
				Replace: []*pb.Update{{
					Path: pathCli,
					Val: &pb.TypedValue{
						Value: &pb.TypedValue_AsciiVal{AsciiVal: "hostname foo\nip routing"}},
				}},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := newSetRequest(tc.setOps)
			if err != nil {
				t.Fatal(err)
			}
			if diff := test.Diff(tc.exp, *got); diff != "" {
				t.Errorf("unexpected diff: %s", diff)
			}
		})
	}
}

func TestStrUpdateVal(t *testing.T) {
	anyBytes, err := proto.Marshal(&pb.ModelData{Name: "foobar"})
	if err != nil {
		t.Fatal(err)
	}
	anyMessage := &any.Any{TypeUrl: "gnmi/ModelData", Value: anyBytes}
	anyString := proto.CompactTextString(anyMessage)

	for name, tc := range map[string]struct {
		update *pb.Update
		exp    string
	}{
		"JSON Value": {
			update: &pb.Update{
				Value: &pb.Value{
					Value: []byte(`{"foo":"bar"}`),
					Type:  pb.Encoding_JSON}},
			exp: `{
  "foo": "bar"
}`,
		},
		"JSON_IETF Value": {
			update: &pb.Update{
				Value: &pb.Value{
					Value: []byte(`{"foo":"bar"}`),
					Type:  pb.Encoding_JSON_IETF}},
			exp: `{
  "foo": "bar"
}`,
		},
		"BYTES Value": {
			update: &pb.Update{
				Value: &pb.Value{
					Value: []byte{0xde, 0xad},
					Type:  pb.Encoding_BYTES}},
			exp: "3q0=",
		},
		"PROTO Value": {
			update: &pb.Update{
				Value: &pb.Value{
					Value: []byte{0xde, 0xad},
					Type:  pb.Encoding_PROTO}},
			exp: "3q0=",
		},
		"ASCII Value": {
			update: &pb.Update{
				Value: &pb.Value{
					Value: []byte("foobar"),
					Type:  pb.Encoding_ASCII}},
			exp: "foobar",
		},
		"INVALID Value": {
			update: &pb.Update{
				Value: &pb.Value{
					Value: []byte("foobar"),
					Type:  pb.Encoding(42)}},
			exp: "foobar",
		},
		"StringVal": {
			update: &pb.Update{Val: &pb.TypedValue{
				Value: &pb.TypedValue_StringVal{StringVal: "foobar"}}},
			exp: "foobar",
		},
		"IntVal": {
			update: &pb.Update{Val: &pb.TypedValue{
				Value: &pb.TypedValue_IntVal{IntVal: -42}}},
			exp: "-42",
		},
		"UintVal": {
			update: &pb.Update{Val: &pb.TypedValue{
				Value: &pb.TypedValue_UintVal{UintVal: 42}}},
			exp: "42",
		},
		"BoolVal": {
			update: &pb.Update{Val: &pb.TypedValue{
				Value: &pb.TypedValue_BoolVal{BoolVal: true}}},
			exp: "true",
		},
		"BytesVal": {
			update: &pb.Update{Val: &pb.TypedValue{
				Value: &pb.TypedValue_BytesVal{BytesVal: []byte{0xde, 0xad}}}},
			exp: "3q0=",
		},
		"FloatVal": {
			update: &pb.Update{Val: &pb.TypedValue{
				Value: &pb.TypedValue_FloatVal{FloatVal: 3.14}}},
			exp: "3.14",
		},
		"DecimalVal": {
			update: &pb.Update{Val: &pb.TypedValue{
				Value: &pb.TypedValue_DecimalVal{
					DecimalVal: &pb.Decimal64{Digits: 314, Precision: 2},
				}}},
			exp: "3.14",
		},
		"LeafListVal": {
			update: &pb.Update{Val: &pb.TypedValue{
				Value: &pb.TypedValue_LeaflistVal{
					LeaflistVal: &pb.ScalarArray{Element: []*pb.TypedValue{
						&pb.TypedValue{Value: &pb.TypedValue_BoolVal{BoolVal: true}},
						&pb.TypedValue{Value: &pb.TypedValue_AsciiVal{AsciiVal: "foobar"}},
					}},
				}}},
			exp: "[true, foobar]",
		},
		"AnyVal": {
			update: &pb.Update{Val: &pb.TypedValue{
				Value: &pb.TypedValue_AnyVal{AnyVal: anyMessage}}},
			exp: anyString,
		},
		"JsonVal": {
			update: &pb.Update{Val: &pb.TypedValue{
				Value: &pb.TypedValue_JsonVal{JsonVal: []byte(`{"foo":"bar"}`)}}},
			exp: `{
  "foo": "bar"
}`,
		},
		"JsonIetfVal": {
			update: &pb.Update{Val: &pb.TypedValue{
				Value: &pb.TypedValue_JsonIetfVal{JsonIetfVal: []byte(`{"foo":"bar"}`)}}},
			exp: `{
  "foo": "bar"
}`,
		},
		"AsciiVal": {
			update: &pb.Update{Val: &pb.TypedValue{
				Value: &pb.TypedValue_AsciiVal{AsciiVal: "foobar"}}},
			exp: "foobar",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := StrUpdateVal(tc.update)
			if got != tc.exp {
				t.Errorf("Expected: %q Got: %q", tc.exp, got)
			}
		})
	}
}

func TestExtractJSON(t *testing.T) {
	jsonFile, err := ioutil.TempFile("", "extractJSON")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(jsonFile.Name())
	if _, err := jsonFile.Write([]byte(`"jsonFile"`)); err != nil {
		jsonFile.Close()
		t.Fatal(err)
	}
	if err := jsonFile.Close(); err != nil {
		t.Fatal(err)
	}

	for val, exp := range map[string][]byte{
		jsonFile.Name(): []byte(`"jsonFile"`),
		"foobar":        []byte(`"foobar"`),
		`"foobar"`:      []byte(`"foobar"`),
		"Val: true":     []byte(`"Val: true"`),
		"host42":        []byte(`"host42"`),
		"42":            []byte("42"),
		"-123.43":       []byte("-123.43"),
		"0xFFFF":        []byte("0xFFFF"),
		// Int larger than can fit in 32 bits should be quoted
		"0x8000000000":  []byte(`"0x8000000000"`),
		"-0x8000000000": []byte(`"-0x8000000000"`),
		"true":          []byte("true"),
		"false":         []byte("false"),
		"null":          []byte("null"),
		"{true: 42}":    []byte("{true: 42}"),
		"[]":            []byte("[]"),
	} {
		t.Run(val, func(t *testing.T) {
			got := extractJSON(val)
			if !bytes.Equal(exp, got) {
				t.Errorf("Unexpected diff. Expected: %q Got: %q", exp, got)
			}
		})
	}
}
