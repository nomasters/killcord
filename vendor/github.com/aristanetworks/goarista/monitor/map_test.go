// Copyright (c) 2017 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// This code was forked from the Go project, here's the original copyright header:

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package monitor

import (
	"encoding/json"
	"expvar"
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
)

func TestMapCounter(t *testing.T) {
	colors := NewMap("colors-in-french")

	red := expvar.NewString("red-in-french")
	red.Set("rouge")
	colors.Set("red", red)
	blue := expvar.NewString("blue-in-french")
	blue.Set("bleu")
	colors.Set("blue", blue)
	green := expvar.NewString("green-in-french")
	green.Set("vert")
	colors.Set("green", green)
	colors.Delete("green")
	if x := colors.Get("red").(*expvar.String).Value(); x != "rouge" {
		t.Errorf(`colors.m["red"] = %v, want "rouge"`, x)
	}
	if x := colors.Get("blue").(*expvar.String).Value(); x != "bleu" {
		t.Errorf(`colors.m["blue"] = %v, want "bleu"`, x)
	}

	// colors.String() should be `{"red":"rouge", "blue":"bleu"}`,
	// though the order of red and blue could vary.
	s := colors.String()
	var j interface{}
	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		t.Fatalf("colors.String() isn't valid JSON: %v", err)
	}
	m, ok := j.(map[string]interface{})
	if !ok {
		t.Error("colors.String() didn't produce a map.")
	}
	if len(m) != 2 {
		t.Error("Should've been only 2 entries in", m)
	}

	if rouge, ok := m["red"].(string); !ok || rouge != "rouge" {
		t.Error("bad value for red:", m)
	}
}

func BenchmarkMapSet(b *testing.B) {
	m := new(Map)

	v := new(expvar.Int)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Set("red", v)
		}
	})
}

func BenchmarkMapSetDifferent(b *testing.B) {
	procKeys := make([][]string, runtime.GOMAXPROCS(0))
	for i := range procKeys {
		keys := make([]string, 4)
		for j := range keys {
			keys[j] = fmt.Sprint(i, j)
		}
		procKeys[i] = keys
	}

	m := new(Map)
	v := new(expvar.Int)
	b.ResetTimer()

	var n int32
	b.RunParallel(func(pb *testing.PB) {
		i := int(atomic.AddInt32(&n, 1)-1) % len(procKeys)
		keys := procKeys[i]

		for pb.Next() {
			for _, k := range keys {
				m.Set(k, v)
			}
		}
	})
}

func BenchmarkMapSetString(b *testing.B) {
	m := new(Map)

	v := new(expvar.String)
	v.Set("Hello, !")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Set("red", v)
		}
	})
}
