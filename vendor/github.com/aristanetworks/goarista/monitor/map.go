// Copyright (c) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// This code was forked from the Go project, here's the original copyright header:

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package monitor

import (
	"bytes"
	"expvar"
	"fmt"
	"sync"
)

// Map is a string-to-Var map variable that satisfies the Var interface.
// This a streamlined, more efficient version of expvar.Map, that also
// supports deletion.
type Map struct {
	m sync.Map // map[string]expvar.Var
}

func (v *Map) String() string {
	var b bytes.Buffer
	b.WriteByte('{')
	first := true
	v.m.Range(func(k, value interface{}) bool {
		if !first {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%q: %v", k, value)
		first = false
		return true
	})
	b.WriteByte('}')
	return b.String()
}

// Get atomically returns the Var for the given key or nil.
func (v *Map) Get(key string) expvar.Var {
	i, _ := v.m.Load(key)
	av, _ := i.(expvar.Var)
	return av
}

// Set atomically associates the given Var to the given key.
func (v *Map) Set(key string, av expvar.Var) {
	// Before we store the value, check to see whether the key is new. Try a Load
	// before LoadOrStore: LoadOrStore causes the key interface to escape even on
	// the Load path.
	if _, ok := v.m.Load(key); !ok {
		if _, dup := v.m.LoadOrStore(key, av); !dup {
			return
		}
	}

	v.m.Store(key, av)
}

// Delete atomically deletes the given key if it exists.
func (v *Map) Delete(key string) {
	v.m.Delete(key)
}

// Do calls f for each entry in the map.
// The map is locked during the iteration,
// but existing entries may be concurrently updated.
func (v *Map) Do(f func(expvar.KeyValue)) {
	v.m.Range(func(k, value interface{}) bool {
		f(expvar.KeyValue{Key: k.(string), Value: value.(expvar.Var)})
		return true
	})
}

// NewMap creates a new Map and publishes it with the given name.
func NewMap(name string) *Map {
	v := new(Map)
	expvar.Publish(name, v)
	return v
}
