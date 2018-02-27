// Copyright (c) 2015 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package key

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/aristanetworks/goarista/value"
)

// Key represents the Key in the updates and deletes of the Notification
// objects.  The only reason this exists is that Go won't let us define
// our own hash function for non-hashable types, and unfortunately we
// need to be able to index maps by map[string]interface{} objects.
type Key interface {
	Key() interface{}
	String() string
	Equal(other interface{}) bool
}

type keyImpl struct {
	key interface{}
}

type strKey string

type int8Key int8
type int16Key int16
type int32Key int32
type int64Key int64

type uint8Key int8
type uint16Key int16
type uint32Key int32
type uint64Key int64

type float32Key float32
type float64Key float64

type boolKey bool

// New wraps the given value in a Key.
// This function panics if the value passed in isn't allowed in a Key or
// doesn't implement value.Value.
func New(intf interface{}) Key {
	switch t := intf.(type) {
	case map[string]interface{}:
		return composite{sentinel, t}
	case string:
		return strKey(t)
	case int8:
		return int8Key(t)
	case int16:
		return int16Key(t)
	case int32:
		return int32Key(t)
	case int64:
		return int64Key(t)
	case uint8:
		return uint8Key(t)
	case uint16:
		return uint16Key(t)
	case uint32:
		return uint32Key(t)
	case uint64:
		return uint64Key(t)
	case float32:
		return float32Key(t)
	case float64:
		return float64Key(t)
	case bool:
		return boolKey(t)
	case value.Value:
		return keyImpl{key: intf}
	default:
		panic(fmt.Sprintf("Invalid type for key: %T", intf))
	}
}

func (k keyImpl) Key() interface{} {
	return k.key
}

func (k keyImpl) String() string {
	return stringify(k.key)
}

func (k keyImpl) GoString() string {
	return fmt.Sprintf("key.New(%#v)", k.Key())
}

func (k keyImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.Key())
}

func (k keyImpl) Equal(other interface{}) bool {
	o, ok := other.(keyImpl)
	return ok && keyEqual(k.key, o.key)
}

// Comparable types have an equality-testing method.
type Comparable interface {
	// Equal returns true if this object is equal to the other one.
	Equal(other interface{}) bool
}

func mapStringEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for k, av := range a {
		if bv, ok := b[k]; !ok || !keyEqual(av, bv) {
			return false
		}
	}
	return true
}

func keyEqual(a, b interface{}) bool {
	switch a := a.(type) {
	case map[string]interface{}:
		b, ok := b.(map[string]interface{})
		return ok && mapStringEqual(a, b)
	case map[Key]interface{}:
		b, ok := b.(map[Key]interface{})
		if !ok || len(a) != len(b) {
			return false
		}
		for k, av := range a {
			if bv, ok := b[k]; !ok || !keyEqual(av, bv) {
				return false
			}
		}
		return true
	case Comparable:
		return a.Equal(b)
	}

	return a == b
}

func (k strKey) Key() interface{} {
	return string(k)
}

func (k strKey) String() string {
	return string(k)
}

func (k strKey) GoString() string {
	return fmt.Sprintf("key.New(%q)", string(k))
}

func (k strKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(escape(string(k)))
}

func (k strKey) Equal(other interface{}) bool {
	o, ok := other.(strKey)
	return ok && k == o
}

// Key interface implementation for int8
func (k int8Key) Key() interface{} {
	return int8(k)
}

func (k int8Key) String() string {
	return strconv.FormatInt(int64(k), 10)
}

func (k int8Key) GoString() string {
	return fmt.Sprintf("key.New(%d)", int8(k))
}

func (k int8Key) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(k), 10)), nil
}

func (k int8Key) Equal(other interface{}) bool {
	o, ok := other.(int8Key)
	return ok && k == o
}

// Key interface implementation for int16
func (k int16Key) Key() interface{} {
	return int16(k)
}

func (k int16Key) String() string {
	return strconv.FormatInt(int64(k), 10)
}

func (k int16Key) GoString() string {
	return fmt.Sprintf("key.New(%d)", int16(k))
}

func (k int16Key) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(k), 10)), nil
}

func (k int16Key) Equal(other interface{}) bool {
	o, ok := other.(int16Key)
	return ok && k == o
}

// Key interface implementation for int32
func (k int32Key) Key() interface{} {
	return int32(k)
}

func (k int32Key) String() string {
	return strconv.FormatInt(int64(k), 10)
}

func (k int32Key) GoString() string {
	return fmt.Sprintf("key.New(%d)", int32(k))
}

func (k int32Key) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(k), 10)), nil
}

func (k int32Key) Equal(other interface{}) bool {
	o, ok := other.(int32Key)
	return ok && k == o
}

// Key interface implementation for int64
func (k int64Key) Key() interface{} {
	return int64(k)
}

func (k int64Key) String() string {
	return strconv.FormatInt(int64(k), 10)
}

func (k int64Key) GoString() string {
	return fmt.Sprintf("key.New(%d)", int64(k))
}

func (k int64Key) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(k), 10)), nil
}

func (k int64Key) Equal(other interface{}) bool {
	o, ok := other.(int64Key)
	return ok && k == o
}

// Key interface implementation for uint8
func (k uint8Key) Key() interface{} {
	return uint8(k)
}

func (k uint8Key) String() string {
	return strconv.FormatUint(uint64(k), 10)
}

func (k uint8Key) GoString() string {
	return fmt.Sprintf("key.New(%d)", uint8(k))
}

func (k uint8Key) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(k), 10)), nil
}

func (k uint8Key) Equal(other interface{}) bool {
	o, ok := other.(uint8Key)
	return ok && k == o
}

// Key interface implementation for uint16
func (k uint16Key) Key() interface{} {
	return uint16(k)
}

func (k uint16Key) String() string {
	return strconv.FormatUint(uint64(k), 10)
}

func (k uint16Key) GoString() string {
	return fmt.Sprintf("key.New(%d)", uint16(k))
}

func (k uint16Key) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(k), 10)), nil
}

func (k uint16Key) Equal(other interface{}) bool {
	o, ok := other.(uint16Key)
	return ok && k == o
}

// Key interface implementation for uint32
func (k uint32Key) Key() interface{} {
	return uint32(k)
}

func (k uint32Key) String() string {
	return strconv.FormatUint(uint64(k), 10)
}

func (k uint32Key) GoString() string {
	return fmt.Sprintf("key.New(%d)", uint32(k))
}

func (k uint32Key) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(k), 10)), nil
}

func (k uint32Key) Equal(other interface{}) bool {
	o, ok := other.(uint32Key)
	return ok && k == o
}

// Key interface implementation for uint64
func (k uint64Key) Key() interface{} {
	return uint64(k)
}

func (k uint64Key) String() string {
	return strconv.FormatUint(uint64(k), 10)
}

func (k uint64Key) GoString() string {
	return fmt.Sprintf("key.New(%d)", uint64(k))
}

func (k uint64Key) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(k), 10)), nil
}

func (k uint64Key) Equal(other interface{}) bool {
	o, ok := other.(uint64Key)
	return ok && k == o
}

// Key interface implementation for float32
func (k float32Key) Key() interface{} {
	return float32(k)
}

func (k float32Key) String() string {
	return "f" + strconv.FormatInt(int64(math.Float32bits(float32(k))), 10)
}

func (k float32Key) GoString() string {
	return fmt.Sprintf("key.New(%v)", float32(k))
}

func (k float32Key) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(k), 'g', -1, 32)), nil
}

func (k float32Key) Equal(other interface{}) bool {
	o, ok := other.(float32Key)
	return ok && k == o
}

// Key interface implementation for float64
func (k float64Key) Key() interface{} {
	return float64(k)
}

func (k float64Key) String() string {
	return "f" + strconv.FormatInt(int64(math.Float64bits(float64(k))), 10)
}

func (k float64Key) GoString() string {
	return fmt.Sprintf("key.New(%v)", float64(k))
}

func (k float64Key) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(k), 'g', -1, 64)), nil
}

func (k float64Key) Equal(other interface{}) bool {
	o, ok := other.(float64Key)
	return ok && k == o
}

// Key interface implementation for bool
func (k boolKey) Key() interface{} {
	return bool(k)
}

func (k boolKey) String() string {
	return strconv.FormatBool(bool(k))
}

func (k boolKey) GoString() string {
	return fmt.Sprintf("key.New(%v)", bool(k))
}

func (k boolKey) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatBool(bool(k))), nil
}

func (k boolKey) Equal(other interface{}) bool {
	o, ok := other.(boolKey)
	return ok && k == o
}
