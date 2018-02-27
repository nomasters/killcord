// Copyright (c) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package path

import (
	"fmt"
	"testing"

	"github.com/aristanetworks/goarista/key"
	"github.com/aristanetworks/goarista/value"
)

func TestNew(t *testing.T) {
	tcases := []struct {
		in  []interface{}
		out Path
	}{
		{
			in:  nil,
			out: Path{},
		}, {
			in:  []interface{}{},
			out: Path{},
		}, {
			in:  []interface{}{"foo", key.New("bar"), true},
			out: Path{key.New("foo"), key.New("bar"), key.New(true)},
		}, {
			in:  []interface{}{int8(5), int16(5), int32(5), int64(5)},
			out: Path{key.New(int8(5)), key.New(int16(5)), key.New(int32(5)), key.New(int64(5))},
		}, {
			in: []interface{}{uint8(5), uint16(5), uint32(5), uint64(5)},
			out: Path{key.New(uint8(5)), key.New(uint16(5)), key.New(uint32(5)),
				key.New(uint64(5))},
		}, {
			in:  []interface{}{float32(5), float64(5)},
			out: Path{key.New(float32(5)), key.New(float64(5))},
		}, {
			in:  []interface{}{customKey{i: &a}, map[string]interface{}{}},
			out: Path{key.New(customKey{i: &a}), key.New(map[string]interface{}{})},
		},
	}
	for i, tcase := range tcases {
		if p := New(tcase.in...); !Equal(p, tcase.out) {
			t.Fatalf("Test %d failed: %#v != %#v", i, p, tcase.out)
		}
	}
}

func TestClone(t *testing.T) {
	if !Equal(Clone(Path{}), Path{}) {
		t.Error("Clone(Path{}) != Path{}")
	}
	a := Path{key.New("foo"), key.New("bar")}
	b, c := Clone(a), Clone(a)
	b[1] = key.New("baz")
	if Equal(a, b) || !Equal(a, c) {
		t.Error("Clone is not making a copied path")
	}
}

func TestAppend(t *testing.T) {
	tcases := []struct {
		a      Path
		b      []interface{}
		result Path
	}{
		{
			a:      Path{},
			b:      []interface{}{},
			result: Path{},
		}, {
			a:      Path{key.New("foo")},
			b:      []interface{}{},
			result: Path{key.New("foo")},
		}, {
			a:      Path{},
			b:      []interface{}{"foo", key.New("bar")},
			result: Path{key.New("foo"), key.New("bar")},
		}, {
			a:      Path{key.New("foo")},
			b:      []interface{}{int64(0), key.New("bar")},
			result: Path{key.New("foo"), key.New(int64(0)), key.New("bar")},
		},
	}
	for i, tcase := range tcases {
		if p := Append(tcase.a, tcase.b...); !Equal(p, tcase.result) {
			t.Fatalf("Test %d failed: %#v != %#v", i, p, tcase.result)
		}
	}
}

func TestBase(t *testing.T) {
	if Base(Path{}) != nil {
		t.Fatal("Base of empty Path should be nil")
	}
	tcases := []struct {
		in  Path
		out key.Key
	}{
		{
			in:  Path{key.New("foo")},
			out: key.New("foo"),
		}, {
			in:  Path{key.New("foo"), key.New("bar")},
			out: key.New("bar"),
		},
	}
	for _, tcase := range tcases {
		if !Base(tcase.in).Equal(tcase.out) {
			t.Fatalf("Base of %#v != %#v", tcase.in, tcase.out)
		}
	}
}

type customKey struct {
	i *int
}

func (c customKey) String() string {
	return fmt.Sprintf("customKey=%d", *c.i)
}

func (c customKey) MarshalJSON() ([]byte, error) {
	return nil, nil
}

func (c customKey) ToBuiltin() interface{} {
	return nil
}

func (c customKey) Equal(other interface{}) bool {
	o, ok := other.(customKey)
	return ok && *c.i == *o.i
}

var (
	_ value.Value    = customKey{}
	_ key.Comparable = customKey{}
	a                = 1
	b                = 1
)

func TestEqual(t *testing.T) {
	tcases := []struct {
		a      Path
		b      Path
		result bool
	}{
		{
			a:      nil,
			b:      nil,
			result: true,
		}, {
			a:      nil,
			b:      Path{},
			result: true,
		}, {
			a:      Path{},
			b:      nil,
			result: true,
		}, {
			a:      Path{},
			b:      Path{},
			result: true,
		}, {
			a:      Path{},
			b:      Path{key.New("")},
			result: false,
		}, {
			a:      Path{Wildcard},
			b:      Path{key.New("foo")},
			result: false,
		}, {
			a:      Path{Wildcard},
			b:      Path{Wildcard},
			result: true,
		}, {
			a:      Path{key.New("foo")},
			b:      Path{key.New("foo")},
			result: true,
		}, {
			a:      Path{key.New(true)},
			b:      Path{key.New(false)},
			result: false,
		}, {
			a:      Path{key.New(int32(5))},
			b:      Path{key.New(int64(5))},
			result: false,
		}, {
			a:      Path{key.New("foo")},
			b:      Path{key.New("foo"), key.New("bar")},
			result: false,
		}, {
			a:      Path{key.New("foo"), key.New("bar")},
			b:      Path{key.New("foo")},
			result: false,
		}, {
			a:      Path{key.New(uint8(0)), key.New(int8(0))},
			b:      Path{key.New(int8(0)), key.New(uint8(0))},
			result: false,
		},
		// Ensure that we check deep equality.
		{
			a:      Path{key.New(map[string]interface{}{})},
			b:      Path{key.New(map[string]interface{}{})},
			result: true,
		}, {
			a:      Path{key.New(customKey{i: &a})},
			b:      Path{key.New(customKey{i: &b})},
			result: true,
		},
	}
	for i, tcase := range tcases {
		if result := Equal(tcase.a, tcase.b); result != tcase.result {
			t.Fatalf("Test %d failed: a: %#v; b: %#v, result: %t",
				i, tcase.a, tcase.b, tcase.result)
		}
	}
}

func TestMatch(t *testing.T) {
	tcases := []struct {
		a      Path
		b      Path
		result bool
	}{
		{
			a:      nil,
			b:      nil,
			result: true,
		}, {
			a:      nil,
			b:      Path{},
			result: true,
		}, {
			a:      Path{},
			b:      nil,
			result: true,
		}, {
			a:      Path{},
			b:      Path{},
			result: true,
		}, {
			a:      Path{},
			b:      Path{key.New("foo")},
			result: false,
		}, {
			a:      Path{Wildcard},
			b:      Path{key.New("foo")},
			result: true,
		}, {
			a:      Path{key.New("foo")},
			b:      Path{Wildcard},
			result: false,
		}, {
			a:      Path{Wildcard},
			b:      Path{key.New("foo"), key.New("bar")},
			result: false,
		}, {
			a:      Path{Wildcard, Wildcard},
			b:      Path{key.New(int64(0))},
			result: false,
		}, {
			a:      Path{Wildcard, Wildcard},
			b:      Path{key.New(int64(0)), key.New(int32(0))},
			result: true,
		}, {
			a:      Path{Wildcard, key.New(false)},
			b:      Path{key.New(true), Wildcard},
			result: false,
		},
	}
	for i, tcase := range tcases {
		if result := Match(tcase.a, tcase.b); result != tcase.result {
			t.Fatalf("Test %d failed: a: %#v; b: %#v, result: %t",
				i, tcase.a, tcase.b, tcase.result)
		}
	}
}

func TestHasPrefix(t *testing.T) {
	tcases := []struct {
		a      Path
		b      Path
		result bool
	}{
		{
			a:      nil,
			b:      nil,
			result: true,
		}, {
			a:      nil,
			b:      Path{},
			result: true,
		}, {
			a:      Path{},
			b:      nil,
			result: true,
		}, {
			a:      Path{},
			b:      Path{},
			result: true,
		}, {
			a:      Path{},
			b:      Path{key.New("foo")},
			result: false,
		}, {
			a:      Path{key.New("foo")},
			b:      Path{},
			result: true,
		}, {
			a:      Path{key.New(true)},
			b:      Path{key.New(false)},
			result: false,
		}, {
			a:      Path{key.New("foo"), key.New("bar")},
			b:      Path{key.New("bar"), key.New("foo")},
			result: false,
		}, {
			a:      Path{key.New(int8(0)), key.New(uint8(0))},
			b:      Path{key.New(uint8(0)), key.New(uint8(0))},
			result: false,
		}, {
			a:      Path{key.New(true), key.New(true)},
			b:      Path{key.New(true), key.New(true), key.New(true)},
			result: false,
		}, {
			a:      Path{key.New(true), key.New(true), key.New(true)},
			b:      Path{key.New(true), key.New(true)},
			result: true,
		}, {
			a:      Path{Wildcard, key.New(int32(0)), Wildcard},
			b:      Path{key.New(int64(0)), Wildcard},
			result: false,
		},
	}
	for i, tcase := range tcases {
		if result := HasPrefix(tcase.a, tcase.b); result != tcase.result {
			t.Fatalf("Test %d failed: a: %#v; b: %#v, result: %t",
				i, tcase.a, tcase.b, tcase.result)
		}
	}
}

func TestMatchPrefix(t *testing.T) {
	tcases := []struct {
		a      Path
		b      Path
		result bool
	}{
		{
			a:      nil,
			b:      nil,
			result: true,
		}, {
			a:      nil,
			b:      Path{},
			result: true,
		}, {
			a:      Path{},
			b:      nil,
			result: true,
		}, {
			a:      Path{},
			b:      Path{},
			result: true,
		}, {
			a:      Path{},
			b:      Path{key.New("foo")},
			result: false,
		}, {
			a:      Path{key.New("foo")},
			b:      Path{},
			result: true,
		}, {
			a:      Path{key.New("foo")},
			b:      Path{Wildcard},
			result: false,
		}, {
			a:      Path{Wildcard},
			b:      Path{key.New("foo")},
			result: true,
		}, {
			a:      Path{Wildcard},
			b:      Path{key.New("foo"), key.New("bar")},
			result: false,
		}, {
			a:      Path{Wildcard, key.New(true)},
			b:      Path{key.New(false), Wildcard},
			result: false,
		}, {
			a:      Path{Wildcard, key.New(int32(0)), key.New(int16(0))},
			b:      Path{key.New(int64(0)), key.New(int32(0))},
			result: true,
		},
	}
	for i, tcase := range tcases {
		if result := MatchPrefix(tcase.a, tcase.b); result != tcase.result {
			t.Fatalf("Test %d failed: a: %#v; b: %#v, result: %t",
				i, tcase.a, tcase.b, tcase.result)
		}
	}
}

func TestFromString(t *testing.T) {
	tcases := []struct {
		in  string
		out Path
	}{
		{
			in:  "",
			out: Path{},
		}, {
			in:  "/",
			out: Path{key.New("")},
		}, {
			in:  "//",
			out: Path{key.New(""), key.New("")},
		}, {
			in:  "foo",
			out: Path{key.New("foo")},
		}, {
			in:  "/foo",
			out: Path{key.New("foo")},
		}, {
			in:  "foo/bar",
			out: Path{key.New("foo"), key.New("bar")},
		}, {
			in:  "/foo/bar",
			out: Path{key.New("foo"), key.New("bar")},
		}, {
			in:  "foo/bar/baz",
			out: Path{key.New("foo"), key.New("bar"), key.New("baz")},
		}, {
			in:  "/foo/bar/baz",
			out: Path{key.New("foo"), key.New("bar"), key.New("baz")},
		}, {
			in:  "0/123/456/789",
			out: Path{key.New("0"), key.New("123"), key.New("456"), key.New("789")},
		}, {
			in:  "/0/123/456/789",
			out: Path{key.New("0"), key.New("123"), key.New("456"), key.New("789")},
		}, {
			in:  "`~!@#$%^&*()_+{}\\/|[];':\"<>?,./",
			out: Path{key.New("`~!@#$%^&*()_+{}\\"), key.New("|[];':\"<>?,."), key.New("")},
		}, {
			in:  "/`~!@#$%^&*()_+{}\\/|[];':\"<>?,./",
			out: Path{key.New("`~!@#$%^&*()_+{}\\"), key.New("|[];':\"<>?,."), key.New("")},
		},
	}
	for i, tcase := range tcases {
		if p := FromString(tcase.in); !Equal(p, tcase.out) {
			t.Fatalf("Test %d failed: %#v != %#v", i, p, tcase.out)
		}
	}
}

func TestString(t *testing.T) {
	tcases := []struct {
		in  Path
		out string
	}{
		{
			in:  Path{},
			out: "/",
		}, {
			in:  Path{key.New("")},
			out: "/",
		}, {
			in:  Path{key.New("foo")},
			out: "/foo",
		}, {
			in:  Path{key.New("foo"), key.New("bar")},
			out: "/foo/bar",
		}, {
			in:  Path{key.New("/foo"), key.New("bar")},
			out: "//foo/bar",
		}, {
			in:  Path{key.New("foo"), key.New("bar/")},
			out: "/foo/bar/",
		}, {
			in:  Path{key.New(""), key.New("foo"), key.New("bar")},
			out: "//foo/bar",
		}, {
			in:  Path{key.New("foo"), key.New("bar"), key.New("")},
			out: "/foo/bar/",
		}, {
			in:  Path{key.New("/"), key.New("foo"), key.New("bar")},
			out: "///foo/bar",
		}, {
			in:  Path{key.New("foo"), key.New("bar"), key.New("/")},
			out: "/foo/bar//",
		},
	}
	for i, tcase := range tcases {
		if s := tcase.in.String(); s != tcase.out {
			t.Fatalf("Test %d failed: %s != %s", i, s, tcase.out)
		}
	}
}
