// Copyright (c) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package path

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/aristanetworks/goarista/key"
)

// Wildcard is a special key representing any possible path.
var Wildcard = key.New(WildcardType{})

// WildcardType type represents a wildcard element to represent any possible element
// in a path.
// This type must implement value.Value interface so it can be used as a key.
type WildcardType struct{}

func (w WildcardType) String() string {
	return "*"
}

// Equal implements the value.Value interface
func (w WildcardType) Equal(other interface{}) bool {
	_, ok := other.(WildcardType)
	return ok
}

// ToBuiltin implements the value.Value interface
func (w WildcardType) ToBuiltin() interface{} {
	return WildcardType{}
}

// MarshalJSON implements the value.Value interface
func (w WildcardType) MarshalJSON() ([]byte, error) {
	return []byte(`{"_wildcard":{}}`), nil
}

// Map associates paths to values. It allows wildcards. A Map
// is primarily used to register handlers with paths that can
// be easily looked up each time a path is updated.
type Map struct {
	val      interface{}
	ok       bool
	wildcard *Map
	children map[key.Key]*Map
}

// VisitorFunc is a function that handles the value associated
// with a path in a Map. Note that only the value is passed in
// as an argument since the path can be stored inside the value
// if needed.
type VisitorFunc func(v interface{}) error

// Visit calls a function fn for every value in the Map
// that is registered with a match of a path p. In the
// general case, time complexity is linear with respect
// to the length of p but it can be as bad as O(2^len(p))
// if there are a lot of paths with wildcards registered.
//
// Example:
//
// a := path.New("foo", "bar", "baz")
// b := path.New("foo", path.Wildcard, "baz")
// c := path.New(path.Wildcard, "bar", "baz")
// d := path.New("foo", "bar", path.Wildcard)
// e := path.New(path.Wildcard, path.Wildcard, "baz")
// f := path.New(path.Wildcard, "bar", path.Wildcard)
// g := path.New("foo", path.Wildcard, path.Wildcard)
// h := path.New(path.Wildcard, path.Wildcard, path.Wildcard)
//
// m.Set(a, 1)
// m.Set(b, 2)
// m.Set(c, 3)
// m.Set(d, 4)
// m.Set(e, 5)
// m.Set(f, 6)
// m.Set(g, 7)
// m.Set(h, 8)
//
// p := path.New("foo", "bar", "baz")
//
// m.Visit(p, fn)
//
// Result: fn(1), fn(2), fn(3), fn(4), fn(5), fn(6), fn(7) and fn(8)
func (m *Map) Visit(p Path, fn VisitorFunc) error {
	for i, element := range p {
		if m.wildcard != nil {
			if err := m.wildcard.Visit(p[i+1:], fn); err != nil {
				return err
			}
		}
		next, ok := m.children[element]
		if !ok {
			return nil
		}
		m = next
	}
	if !m.ok {
		return nil
	}
	return fn(m.val)
}

// VisitPrefixes calls a function fn for every value in the
// Map that is registered with a prefix of a path p.
//
// Example:
//
// a := path.New()
// b := path.New("foo")
// c := path.New("foo", "bar")
// d := path.New("foo", "baz")
// e := path.New(path.Wildcard, "bar")
//
// m.Set(a, 1)
// m.Set(b, 2)
// m.Set(c, 3)
// m.Set(d, 4)
// m.Set(e, 5)
//
// p := path.New("foo", "bar", "baz")
//
// m.VisitPrefixes(p, fn)
//
// Result: fn(1), fn(2), fn(3), fn(5)
func (m *Map) VisitPrefixes(p Path, fn VisitorFunc) error {
	for i, element := range p {
		if m.ok {
			if err := fn(m.val); err != nil {
				return err
			}
		}
		if m.wildcard != nil {
			if err := m.wildcard.VisitPrefixes(p[i+1:], fn); err != nil {
				return err
			}
		}
		next, ok := m.children[element]
		if !ok {
			return nil
		}
		m = next
	}
	if !m.ok {
		return nil
	}
	return fn(m.val)
}

// VisitPrefixed calls fn for every value in the map that is
// registerd with a path that is prefixed by p. This method
// can be used to visit every registered path if p is the
// empty path (or root path) which prefixes all paths.
//
// Example:
//
// a := path.New("foo")
// b := path.New("foo", "bar")
// c := path.New("foo", "bar", "baz")
// d := path.New("foo", path.Wildcard)
//
// m.Set(a, 1)
// m.Set(b, 2)
// m.Set(c, 3)
// m.Set(d, 4)
//
// p := path.New("foo", "bar")
//
// m.VisitPrefixed(p, fn)
//
// Result: fn(2), fn(3), fn(4)
func (m *Map) VisitPrefixed(p Path, fn VisitorFunc) error {
	for i, element := range p {
		if m.wildcard != nil {
			if err := m.wildcard.VisitPrefixed(p[i+1:], fn); err != nil {
				return err
			}
		}
		next, ok := m.children[element]
		if !ok {
			return nil
		}
		m = next
	}
	return m.visitSubtree(fn)
}

func (m *Map) visitSubtree(fn VisitorFunc) error {
	if m.ok {
		if err := fn(m.val); err != nil {
			return err
		}
	}
	if m.wildcard != nil {
		if err := m.wildcard.visitSubtree(fn); err != nil {
			return err
		}
	}
	for _, next := range m.children {
		if err := next.visitSubtree(fn); err != nil {
			return err
		}
	}
	return nil
}

// Get returns the value registered with an exact match of a
// path p. If there is no exact match for p, Get returns nil
// and false. If p has an exact match and it is set to true,
// Get returns nil and true.
//
// Example:
//
// m.Set(path.New("foo", "bar"), 1)
// m.Set(path.New("baz", "qux"), nil)
//
// a := m.Get(path.New("foo", "bar"))
// b := m.Get(path.New("foo", path.Wildcard))
// c, ok := m.Get(path.New("baz", "qux"))
//
// Result: a == 1, b == nil, c == nil and ok == true
func (m *Map) Get(p Path) (interface{}, bool) {
	for _, element := range p {
		if element.Equal(Wildcard) {
			if m.wildcard == nil {
				return nil, false
			}
			m = m.wildcard
			continue
		}
		next, ok := m.children[element]
		if !ok {
			return nil, false
		}
		m = next
	}
	return m.val, m.ok
}

// Set registers a path p with a value. Any previous value that
// was registered with p is overwritten.
//
// Example:
//
// p := path.New("foo", "bar")
//
// m.Set(p, 0)
// m.Set(p, 1)
//
// v := m.Get(p)
//
// Result: v == 1
func (m *Map) Set(p Path, v interface{}) {
	for _, element := range p {
		if element.Equal(Wildcard) {
			if m.wildcard == nil {
				m.wildcard = &Map{}
			}
			m = m.wildcard
			continue
		}
		if m.children == nil {
			m.children = map[key.Key]*Map{}
		}
		next, ok := m.children[element]
		if !ok {
			next = &Map{}
			m.children[element] = next
		}
		m = next
	}
	m.val, m.ok = v, true
}

// Delete unregisters the value registered with a path. It
// returns true if a value was deleted and false otherwise.
//
// Example:
//
// p := path.New("foo", "bar")
//
// m.Set(p, 0)
//
// a := m.Delete(p)
// b := m.Delete(p)
//
// Result: a == true and b == false
func (m *Map) Delete(p Path) bool {
	maps := make([]*Map, len(p)+1)
	for i, element := range p {
		maps[i] = m
		if element.Equal(Wildcard) {
			if m.wildcard == nil {
				return false
			}
			m = m.wildcard
			continue
		}
		next, ok := m.children[element]
		if !ok {
			return false
		}
		m = next
	}
	m.val, m.ok = nil, false
	maps[len(p)] = m

	// Remove any empty maps.
	for i := len(p); i > 0; i-- {
		m = maps[i]
		if m.ok || m.wildcard != nil || len(m.children) > 0 {
			break
		}
		parent := maps[i-1]
		element := p[i-1]
		if element.Equal(Wildcard) {
			parent.wildcard = nil
		} else {
			delete(parent.children, element)
		}
	}
	return true
}

func (m *Map) String() string {
	var b bytes.Buffer
	m.write(&b, "")
	return b.String()
}

func (m *Map) write(b *bytes.Buffer, indent string) {
	if m.ok {
		b.WriteString(indent)
		fmt.Fprintf(b, "Val: %v", m.val)
		b.WriteString("\n")
	}
	if m.wildcard != nil {
		b.WriteString(indent)
		fmt.Fprintf(b, "Child %q:\n", Wildcard)
		m.wildcard.write(b, indent+"  ")
	}
	children := make([]key.Key, 0, len(m.children))
	for key := range m.children {
		children = append(children, key)
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].String() < children[j].String()
	})

	for _, key := range children {
		child := m.children[key]
		b.WriteString(indent)
		fmt.Fprintf(b, "Child %q:\n", key.String())
		child.write(b, indent+"  ")
	}
}
