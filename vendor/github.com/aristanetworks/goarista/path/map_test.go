// Copyright (c) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package path

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/aristanetworks/goarista/key"
	"github.com/aristanetworks/goarista/test"
	"github.com/aristanetworks/goarista/value"
)

func accumulator(counter map[int]int) VisitorFunc {
	return func(val interface{}) error {
		counter[val.(int)]++
		return nil
	}
}

type pseudoWildcard struct{}

func (w pseudoWildcard) Key() interface{} {
	return struct{}{}
}

func (w pseudoWildcard) String() string {
	return "*"
}

func (w pseudoWildcard) Equal(other interface{}) bool {
	o, ok := other.(pseudoWildcard)
	return ok && w == o
}

func TestWildcardTypeIsNotAKey(t *testing.T) {
	var intf interface{} = WildcardType{}
	_, ok := intf.(key.Key)
	if ok {
		t.Error("WildcardType should not implement key.Key")
	}
}

func TestWildcardTypeEqual(t *testing.T) {
	k1 := key.New(WildcardType{})
	k2 := key.New(WildcardType{})
	if !k1.Equal(k2) {
		t.Error("They should be equal")
	}
	if !Wildcard.Equal(k1) {
		t.Error("They should be equal")
	}
}

func TestWildcardTypeAsValue(t *testing.T) {
	var k value.Value = WildcardType{}
	w := WildcardType{}
	if k.ToBuiltin() != w {
		t.Error("Wildcard builtin is not correct")
	}
}

func TestWildcardMarshalJSON(t *testing.T) {
	b, err := json.Marshal(Wildcard)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"_wildcard":{}}`
	if string(b) != expected {
		t.Errorf("Invalid Wildcard json representation.\nExpected: %s\nReceived: %s",
			expected, string(b))
	}
}

func TestWildcardUniqueness(t *testing.T) {
	if Wildcard.Equal(pseudoWildcard{}) {
		t.Fatal("Wildcard is not unique")
	}
	if Wildcard.Equal(struct{}{}) {
		t.Fatal("Wildcard is not unique")
	}
	if Wildcard.Equal(key.New("*")) {
		t.Fatal("Wildcard is not unique")
	}
}

func TestMapVisit(t *testing.T) {
	m := Map{}
	m.Set(Path{key.New("foo"), key.New("bar"), key.New("baz")}, 1)
	m.Set(Path{Wildcard, key.New("bar"), key.New("baz")}, 2)
	m.Set(Path{Wildcard, Wildcard, key.New("baz")}, 3)
	m.Set(Path{Wildcard, Wildcard, Wildcard}, 4)
	m.Set(Path{key.New("foo"), Wildcard, Wildcard}, 5)
	m.Set(Path{key.New("foo"), key.New("bar"), Wildcard}, 6)
	m.Set(Path{key.New("foo"), Wildcard, key.New("baz")}, 7)
	m.Set(Path{Wildcard, key.New("bar"), Wildcard}, 8)

	m.Set(Path{}, 10)

	m.Set(Path{Wildcard}, 20)
	m.Set(Path{key.New("foo")}, 21)

	m.Set(Path{key.New("zap"), key.New("zip")}, 30)
	m.Set(Path{key.New("zap"), key.New("zip")}, 31)

	m.Set(Path{key.New("zip"), Wildcard}, 40)
	m.Set(Path{key.New("zip"), Wildcard}, 41)

	testCases := []struct {
		path     Path
		expected map[int]int
	}{{
		path:     Path{key.New("foo"), key.New("bar"), key.New("baz")},
		expected: map[int]int{1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 6: 1, 7: 1, 8: 1},
	}, {
		path:     Path{key.New("qux"), key.New("bar"), key.New("baz")},
		expected: map[int]int{2: 1, 3: 1, 4: 1, 8: 1},
	}, {
		path:     Path{key.New("foo"), key.New("qux"), key.New("baz")},
		expected: map[int]int{3: 1, 4: 1, 5: 1, 7: 1},
	}, {
		path:     Path{key.New("foo"), key.New("bar"), key.New("qux")},
		expected: map[int]int{4: 1, 5: 1, 6: 1, 8: 1},
	}, {
		path:     Path{},
		expected: map[int]int{10: 1},
	}, {
		path:     Path{key.New("foo")},
		expected: map[int]int{20: 1, 21: 1},
	}, {
		path:     Path{key.New("foo"), key.New("bar")},
		expected: map[int]int{},
	}, {
		path:     Path{key.New("zap"), key.New("zip")},
		expected: map[int]int{31: 1},
	}, {
		path:     Path{key.New("zip"), key.New("zap")},
		expected: map[int]int{41: 1},
	}}

	for _, tc := range testCases {
		result := make(map[int]int, len(tc.expected))
		m.Visit(tc.path, accumulator(result))
		if diff := test.Diff(tc.expected, result); diff != "" {
			t.Errorf("Test case %v: %s", tc.path, diff)
		}
	}
}

func TestMapVisitError(t *testing.T) {
	m := Map{}
	m.Set(Path{key.New("foo"), key.New("bar")}, 1)
	m.Set(Path{Wildcard, key.New("bar")}, 2)

	errTest := errors.New("Test")

	err := m.Visit(Path{key.New("foo"), key.New("bar")},
		func(v interface{}) error { return errTest })
	if err != errTest {
		t.Errorf("Unexpected error. Expected: %v, Got: %v", errTest, err)
	}
	err = m.VisitPrefixes(Path{key.New("foo"), key.New("bar"), key.New("baz")},
		func(v interface{}) error { return errTest })
	if err != errTest {
		t.Errorf("Unexpected error. Expected: %v, Got: %v", errTest, err)
	}
}

func TestMapGet(t *testing.T) {
	m := Map{}
	m.Set(Path{}, 0)
	m.Set(Path{key.New("foo"), key.New("bar")}, 1)
	m.Set(Path{key.New("foo"), Wildcard}, 2)
	m.Set(Path{Wildcard, key.New("bar")}, 3)
	m.Set(Path{key.New("zap"), key.New("zip")}, 4)
	m.Set(Path{key.New("baz"), key.New("qux")}, nil)

	testCases := []struct {
		path Path
		v    interface{}
		ok   bool
	}{{
		path: Path{},
		v:    0,
		ok:   true,
	}, {
		path: Path{key.New("foo"), key.New("bar")},
		v:    1,
		ok:   true,
	}, {
		path: Path{key.New("foo"), Wildcard},
		v:    2,
		ok:   true,
	}, {
		path: Path{Wildcard, key.New("bar")},
		v:    3,
		ok:   true,
	}, {
		path: Path{key.New("baz"), key.New("qux")},
		v:    nil,
		ok:   true,
	}, {
		path: Path{key.New("bar"), key.New("foo")},
		v:    nil,
	}, {
		path: Path{key.New("zap"), Wildcard},
		v:    nil,
	}}

	for _, tc := range testCases {
		v, ok := m.Get(tc.path)
		if v != tc.v || ok != tc.ok {
			t.Errorf("Test case %v: Expected (v: %v, ok: %t), Got (v: %v, ok: %t)",
				tc.path, tc.v, tc.ok, v, ok)
		}
	}
}

func countNodes(m *Map) int {
	if m == nil {
		return 0
	}
	count := 1
	count += countNodes(m.wildcard)
	for _, child := range m.children {
		count += countNodes(child)
	}
	return count
}

func TestMapDelete(t *testing.T) {
	m := Map{}
	m.Set(Path{}, 0)
	m.Set(Path{Wildcard}, 1)
	m.Set(Path{key.New("foo"), key.New("bar")}, 2)
	m.Set(Path{key.New("foo"), Wildcard}, 3)

	n := countNodes(&m)
	if n != 5 {
		t.Errorf("Initial count wrong. Expected: 5, Got: %d", n)
	}

	testCases := []struct {
		del      Path        // Path to delete
		expected bool        // expected return value of Delete
		visit    Path        // Path to Visit
		before   map[int]int // Expected to find items before deletion
		after    map[int]int // Expected to find items after deletion
		count    int         // Count of nodes
	}{{
		del:      Path{key.New("zap")}, // A no-op Delete
		expected: false,
		visit:    Path{key.New("foo"), key.New("bar")},
		before:   map[int]int{2: 1, 3: 1},
		after:    map[int]int{2: 1, 3: 1},
		count:    5,
	}, {
		del:      Path{key.New("foo"), key.New("bar")},
		expected: true,
		visit:    Path{key.New("foo"), key.New("bar")},
		before:   map[int]int{2: 1, 3: 1},
		after:    map[int]int{3: 1},
		count:    4,
	}, {
		del:      Path{Wildcard},
		expected: true,
		visit:    Path{key.New("foo")},
		before:   map[int]int{1: 1},
		after:    map[int]int{},
		count:    3,
	}, {
		del:      Path{Wildcard},
		expected: false,
		visit:    Path{key.New("foo")},
		before:   map[int]int{},
		after:    map[int]int{},
		count:    3,
	}, {
		del:      Path{key.New("foo"), Wildcard},
		expected: true,
		visit:    Path{key.New("foo"), key.New("bar")},
		before:   map[int]int{3: 1},
		after:    map[int]int{},
		count:    1, // Should have deleted "foo" and "bar" nodes
	}, {
		del:      Path{},
		expected: true,
		visit:    Path{},
		before:   map[int]int{0: 1},
		after:    map[int]int{},
		count:    1, // Root node can't be deleted
	}}

	for i, tc := range testCases {
		beforeResult := make(map[int]int, len(tc.before))
		m.Visit(tc.visit, accumulator(beforeResult))
		if diff := test.Diff(tc.before, beforeResult); diff != "" {
			t.Errorf("Test case %d (%v): %s", i, tc.del, diff)
		}

		if got := m.Delete(tc.del); got != tc.expected {
			t.Errorf("Test case %d (%v): Unexpected return. Expected %t, Got: %t",
				i, tc.del, tc.expected, got)
		}

		afterResult := make(map[int]int, len(tc.after))
		m.Visit(tc.visit, accumulator(afterResult))
		if diff := test.Diff(tc.after, afterResult); diff != "" {
			t.Errorf("Test case %d (%v): %s", i, tc.del, diff)
		}
	}
}

func TestMapVisitPrefixes(t *testing.T) {
	m := Map{}
	m.Set(Path{}, 0)
	m.Set(Path{key.New("foo")}, 1)
	m.Set(Path{key.New("foo"), key.New("bar")}, 2)
	m.Set(Path{key.New("foo"), key.New("bar"), key.New("baz")}, 3)
	m.Set(Path{key.New("foo"), key.New("bar"), key.New("baz"), key.New("quux")}, 4)
	m.Set(Path{key.New("quux"), key.New("bar")}, 5)
	m.Set(Path{key.New("foo"), key.New("quux")}, 6)
	m.Set(Path{Wildcard}, 7)
	m.Set(Path{key.New("foo"), Wildcard}, 8)
	m.Set(Path{Wildcard, key.New("bar")}, 9)
	m.Set(Path{Wildcard, key.New("quux")}, 10)
	m.Set(Path{key.New("quux"), key.New("quux"), key.New("quux"), key.New("quux")}, 11)

	testCases := []struct {
		path     Path
		expected map[int]int
	}{{
		path:     Path{key.New("foo"), key.New("bar"), key.New("baz")},
		expected: map[int]int{0: 1, 1: 1, 2: 1, 3: 1, 7: 1, 8: 1, 9: 1},
	}, {
		path:     Path{key.New("zip"), key.New("zap")},
		expected: map[int]int{0: 1, 7: 1},
	}, {
		path:     Path{key.New("foo"), key.New("zap")},
		expected: map[int]int{0: 1, 1: 1, 8: 1, 7: 1},
	}, {
		path:     Path{key.New("quux"), key.New("quux"), key.New("quux")},
		expected: map[int]int{0: 1, 7: 1, 10: 1},
	}}

	for _, tc := range testCases {
		result := make(map[int]int, len(tc.expected))
		m.VisitPrefixes(tc.path, accumulator(result))
		if diff := test.Diff(tc.expected, result); diff != "" {
			t.Errorf("Test case %v: %s", tc.path, diff)
		}
	}
}

func TestMapVisitPrefixed(t *testing.T) {
	m := Map{}
	m.Set(Path{}, 0)
	m.Set(Path{key.New("qux")}, 1)
	m.Set(Path{key.New("foo")}, 2)
	m.Set(Path{key.New("foo"), key.New("qux")}, 3)
	m.Set(Path{key.New("foo"), key.New("bar")}, 4)
	m.Set(Path{Wildcard, key.New("bar")}, 5)
	m.Set(Path{key.New("foo"), Wildcard}, 6)
	m.Set(Path{key.New("qux"), key.New("foo"), key.New("bar")}, 7)

	testCases := []struct {
		in  Path
		out map[int]int
	}{{
		in:  Path{},
		out: map[int]int{0: 1, 1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 6: 1, 7: 1},
	}, {
		in:  Path{key.New("qux")},
		out: map[int]int{1: 1, 5: 1, 7: 1},
	}, {
		in:  Path{key.New("foo")},
		out: map[int]int{2: 1, 3: 1, 4: 1, 5: 1, 6: 1},
	}, {
		in:  Path{key.New("foo"), key.New("qux")},
		out: map[int]int{3: 1, 6: 1},
	}, {
		in:  Path{key.New("foo"), key.New("bar")},
		out: map[int]int{4: 1, 5: 1, 6: 1},
	}, {
		in:  Path{key.New(int64(0))},
		out: map[int]int{5: 1},
	}, {
		in:  Path{Wildcard},
		out: map[int]int{5: 1},
	}, {
		in:  Path{Wildcard, Wildcard},
		out: map[int]int{},
	}}

	for _, tc := range testCases {
		out := make(map[int]int, len(tc.out))
		m.VisitPrefixed(tc.in, accumulator(out))
		if diff := test.Diff(tc.out, out); diff != "" {
			t.Errorf("Test case %v: %s", tc.out, diff)
		}
	}
}

func TestMapString(t *testing.T) {
	m := Map{}
	m.Set(Path{}, 0)
	m.Set(Path{key.New("foo"), key.New("bar")}, 1)
	m.Set(Path{key.New("foo"), key.New("quux")}, 2)
	m.Set(Path{key.New("foo"), Wildcard}, 3)

	expected := `Val: 0
Child "foo":
  Child "*":
    Val: 3
  Child "bar":
    Val: 1
  Child "quux":
    Val: 2
`
	got := fmt.Sprint(&m)

	if expected != got {
		t.Errorf("Unexpected string. Expected:\n\n%s\n\nGot:\n\n%s", expected, got)
	}
}

func genWords(count, wordLength int) Path {
	chars := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	if count+wordLength > len(chars) {
		panic("need more chars")
	}
	result := make(Path, count)
	for i := 0; i < count; i++ {
		result[i] = key.New(string(chars[i : i+wordLength]))
	}
	return result
}

func benchmarkPathMap(pathLength, pathDepth int, b *testing.B) {
	// Push pathDepth paths, each of length pathLength
	path := genWords(pathLength, 10)
	words := genWords(pathDepth, 10)
	m := &Map{}
	for _, element := range path {
		m.children = map[key.Key]*Map{}
		for _, word := range words {
			m.children[word] = &Map{}
		}
		m = m.children[element]
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Visit(path, func(v interface{}) error { return nil })
	}
}

func BenchmarkPathMap1x25(b *testing.B)  { benchmarkPathMap(1, 25, b) }
func BenchmarkPathMap10x50(b *testing.B) { benchmarkPathMap(10, 25, b) }
func BenchmarkPathMap20x50(b *testing.B) { benchmarkPathMap(20, 25, b) }
