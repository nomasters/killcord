// Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package sizeof

import (
	"fmt"
	"strconv"
	"testing"
	"unsafe"

	"github.com/aristanetworks/goarista/test"
)

type yolo struct {
	i int32
	a [3]int8
	p unsafe.Pointer
}

func (y yolo) String() string {
	return "Yolo"
}

func TestDeepSizeof(t *testing.T) {
	ptrSize := uintptr(unsafe.Sizeof(unsafe.Pointer(t)))
	// hmapStructSize represent the size of struct hmap defined in
	// file /go/src/runtime/hashmap.go
	hmapStructSize := uintptr(unsafe.Sizeof(int(0)) + 2*1 + 2 + 4 +
		2*ptrSize + ptrSize + ptrSize)
	var alignement uintptr = 4
	if ptrSize == 4 {
		alignement = 0
	}
	strHdrSize := unsafe.Sizeof("") // int + ptr to data
	sliceHdrSize := 3 * ptrSize     // ptr to data + 2 * int

	// struct hchan is defined in /go/src/runtime/chan.go
	chanHdrSize := 2*ptrSize + ptrSize + 2 + 2 /* padding */ + 4 + ptrSize + 2*ptrSize +
		2*(2*ptrSize) + ptrSize
	yoloSize := unsafe.Sizeof(yolo{})
	interfaceSize := 2 * ptrSize
	topHashSize := uintptr(8)
	tests := map[string]struct {
		getStruct    func() interface{}
		expectedSize uintptr
	}{
		"bool": {
			getStruct: func() interface{} {
				var test bool
				return &test
			},
			expectedSize: 1,
		},
		"int8": {
			getStruct: func() interface{} {
				test := int8(4)
				return &test
			},
			expectedSize: 1,
		},
		"int16": {
			getStruct: func() interface{} {
				test := int16(4)
				return &test
			},
			expectedSize: 2,
		},
		"int32": {
			getStruct: func() interface{} {
				test := int32(4)
				return &test
			},
			expectedSize: 4,
		},
		"int64": {
			getStruct: func() interface{} {
				test := int64(4)
				return &test
			},
			expectedSize: 8,
		},
		"uint": {
			getStruct: func() interface{} {
				test := uint(4)
				return &test
			},
			expectedSize: ptrSize,
		},
		"uint8": {
			getStruct: func() interface{} {
				test := uint8(4)
				return &test
			},
			expectedSize: 1,
		},
		"uint16": {
			getStruct: func() interface{} {
				test := uint16(4)
				return &test
			},
			expectedSize: 2,
		},
		"uint32": {
			getStruct: func() interface{} {
				test := uint32(4)
				return &test
			},
			expectedSize: 4,
		},
		"uint64": {
			getStruct: func() interface{} {
				test := uint64(4)
				return &test
			},
			expectedSize: 8,
		},
		"uintptr": {
			getStruct: func() interface{} {
				test := uintptr(4)
				return &test
			},
			expectedSize: ptrSize,
		},
		"float32": {
			getStruct: func() interface{} {
				test := float32(4)
				return &test
			},
			expectedSize: 4,
		},
		"float64": {
			getStruct: func() interface{} {
				test := float64(4)
				return &test
			},
			expectedSize: 8,
		},
		"complex64": {
			getStruct: func() interface{} {
				test := complex64(4 + 1i)
				return &test
			},
			expectedSize: 8,
		},
		"complex128": {
			getStruct: func() interface{} {
				test := complex128(4 + 1i)
				return &test
			},
			expectedSize: 16,
		},
		"string": {
			getStruct: func() interface{} {
				test := "Hello Dolly!"
				return &test
			},
			expectedSize: strHdrSize + 12,
		},
		"unsafe_Pointer": {
			getStruct: func() interface{} {
				tmp := uint64(54)
				var test unsafe.Pointer
				test = unsafe.Pointer(&tmp)
				return &test
			},
			expectedSize: ptrSize,
		}, "rune": {
			getStruct: func() interface{} {
				test := rune('A')
				return &test
			},
			expectedSize: 4,
		}, "intPtr": {
			getStruct: func() interface{} {
				test := int(4)
				return &test
			},
			expectedSize: ptrSize,
		}, "FuncPtr": {
			getStruct: func() interface{} {
				test := TestDeepSizeof
				return &test
			},
			expectedSize: ptrSize,
		}, "struct_1": {
			getStruct: func() interface{} {
				v := struct {
					a uint32
					b *uint32
					c struct {
						e [8]byte
						d string
					}
					f string
				}{
					a: 10,
					c: struct {
						e [8]byte
						d string
					}{
						e: [8]byte{0, 1, 2, 3, 4, 5, 6, 7},
						d: "Hello Test!",
					},
					f: "Hello Test!",
				}
				a := uint32(47)
				v.b = &a
				return &v
			},
			expectedSize: 4 + alignement + ptrSize + 8 + strHdrSize*2 +
				11 /* "Hello Test!" */ + 4, /* uint32(47) */
		}, "struct_2": {
			getStruct: func() interface{} {
				v := struct {
					a []byte
					b []byte
					c []byte
				}{
					c: make([]byte, 32, 64),
				}
				v.a = v.c[20:32]
				v.b = v.c[10:20]
				return &v
			},
			expectedSize: 3*sliceHdrSize + 64, /*slice capacity*/
		}, "struct_3": {
			getStruct: func() interface{} {
				type test struct {
					a *byte
					c []byte
				}
				tmp := make([]byte, 64, 128)
				v := (*test)(unsafe.Pointer(&tmp[16]))
				v.c = tmp
				v.a = (*byte)(unsafe.Pointer(&tmp[5]))
				return v
			},
			// we expect to see 128 bytes as struct test is part of the bytes slice
			// and field c point to it.
			expectedSize: 128,
		}, "map_string_interface": {
			getStruct: func() interface{} {
				return &map[string]interface{}{}
			},
			expectedSize: ptrSize + topHashSize + hmapStructSize +
				(8*(strHdrSize+interfaceSize) + ptrSize),
		}, "map_interface_interface": {
			getStruct: func() interface{} {
				// All the map will use only one bucket because there is less than 8
				// entries in each map. Also for amd64 and i386 the bucket size is
				// computed like in function bucketOf in /go/src/reflect/type.go:
				return &map[interface{}]interface{}{
					// 2 + (8 + 4) 386
					// 2 + (16 + 4) amd64
					uint16(123): "yolo",
					// (4 + 8) + (4 + 28 + 140) + (4 for SWAG) 386
					// (4 + 16) + (8 + 48 + 272) + (4 for SWAG) amd64
					"meow": map[string]string{"SWAG": "yolo"},
					// (12) + (4 + 12) 386
					// (16) + (8 + 16) amd64
					yolo{i: 523}: &yolo{i: 126},
					// (12) + (12) 386
					// (16) + (16) amd64
					fmt.Stringer(yolo{i: 123}): yolo{i: 234},
				}
			},
			// Total
			// 386: (4 + 28 + 140) + 2 + (8 + 4) + (4 + 8) + (4 + 28 + 140) + 4 + (12) +
			// (4 + 12) + 12 + 12
			// amd64: (8 + 48 + 272) + 2 + (16 + 4) + (4 + 16) + (8 + 48 + 272) + (4) +
			// (16) + (8 + 16) + (16) + (16)
			expectedSize: (ptrSize + topHashSize + hmapStructSize +
				(8*(2*interfaceSize) + ptrSize)) +
				(unsafe.Sizeof(uint16(123)) + strHdrSize + 4 /* "yolo" */) +
				(strHdrSize + 4 /* "meow" */ +
					(ptrSize + hmapStructSize + topHashSize +
						(8*(2*strHdrSize) + ptrSize)) /*map[string]string*/ +
					4 /* "SWAG" */) +
				(yoloSize /* obj: */ + (ptrSize + yoloSize) /* &obj */) +
				yoloSize*2,
		}, "struct_4": {
			getStruct: func() interface{} {
				return &struct {
					a map[interface{}]interface{}
					c string
					d []string
				}{
					a: map[interface{}]interface{}{
						uint16(123):                "yolo",
						"meow":                     map[string]string{"SWAG": "yolo"},
						yolo{i: 127}:               &yolo{i: 124},
						fmt.Stringer(yolo{i: 123}): yolo{i: 234},
					}, // 4 (386) or 8 (amd64)
					c: "Hello",                              // 8 (386) or 16 (amd64)
					d: []string{"Bonjour", "Hello", "Hola"}, // 12 (386) or 24 (amd64)
				}
			},
			// Total
			// 386: sizeof(tmp map) + 8(test.c) + 12(test.d) +
			// 3 * 8 (strSlice) + 16(len("Bonjour") + len("Hello")...)
			// amd64: sizeof(tmp map) + 8 (test.b) + 16(test.c) + 24(test.d) +
			// 3 * 16 (strSlice) + 16(len("Bonjour") + len("Hello")...)
			expectedSize: (ptrSize + strHdrSize + sliceHdrSize) + (hmapStructSize +
				topHashSize + (8*(2*2*ptrSize /* interface size */) + ptrSize) +
				unsafe.Sizeof(uint16(123)) +
				strHdrSize + 4 /* "yolo" */ + strHdrSize + 4 /* "meow" */ +
				+(ptrSize + hmapStructSize + topHashSize + (8*(2*strHdrSize) + ptrSize) +
					4 /* "SWAG" */) + yoloSize /* obj */ + (ptrSize + yoloSize) /* &obj */ +
				yoloSize*2) + 5 /* "Hello" */ +
				3*strHdrSize /*strings in strSlice*/ + 11, /* "Bonjour" + "Hola" */
		}, "chan_int": {
			getStruct: func() interface{} {
				test := make(chan int)
				return &test
			},
			// The expected size should be equal to the size of the struct hchan
			// defined in /go/src/runtime/chan.go
			expectedSize: ptrSize + chanHdrSize,
		}, "chan_int_16": {
			getStruct: func() interface{} {
				test := make(chan int, 16)
				return &test
			},
			expectedSize: ptrSize + chanHdrSize + 16*ptrSize,
		}, "chan_yoloPtr_16": {
			getStruct: func() interface{} {
				test := make(chan *yolo, 16)
				for i := 0; i < 16; i++ {
					tmp := &yolo{
						i: int32(i),
					}
					tmp.p = unsafe.Pointer(&tmp.i)
					test <- tmp
				}
				return &test
			},
			expectedSize: ptrSize + chanHdrSize + 16*(ptrSize+yoloSize),
		}, "struct_5": {
			getStruct: func() interface{} {
				tmp := make([]byte, 32)
				test := struct {
					a []byte
					b **uint32
				}{
					a: tmp,
				}
				bob := uint32(42)
				ptrInt := (*uintptr)(unsafe.Pointer(&tmp[0]))
				*ptrInt = uintptr(unsafe.Pointer(&bob))
				test.b = (**uint32)(unsafe.Pointer(&tmp[0]))
				return &test
			},
			expectedSize: sliceHdrSize + ptrSize + 32 + 4,
		}, "struct_6": {
			getStruct: func() interface{} {
				type A struct {
					a uintptr
					b *yolo
				}
				type B struct {
					a *A
					b uintptr
				}
				tmp := make([]byte, 32)
				test := struct {
					a []byte
					b *B
				}{
					a: tmp,
				}
				y := yolo{i: 42}
				test.b = (*B)(unsafe.Pointer(&tmp[0]))
				test.b.a = (*A)(unsafe.Pointer(&tmp[0]))
				test.b.a.b = &y
				return &test
			},
			expectedSize: sliceHdrSize + ptrSize + 32 + yoloSize,
		}, "chan_chan_int_16": {
			getStruct: func() interface{} {
				test := make(chan chan int, 16)
				for i := 0; i < 16; i++ {
					tmp := make(chan int)
					test <- tmp
				}
				return &test
			},
			expectedSize: ptrSize + chanHdrSize*17 + 16*ptrSize,
		}, "chan_yolo_16": {
			getStruct: func() interface{} {
				test := make(chan yolo, 16)
				for i := 0; i < 16; i++ {
					tmp := yolo{
						i: int32(i),
					}
					test <- tmp
				}
				return &test
			},
			expectedSize: ptrSize + chanHdrSize + 16*yoloSize,
		}, "chan_map_string_interface_16)": {
			getStruct: func() interface{} {
				test := make(chan map[string]interface{}, 16)
				for i := 0; i < 16; i++ {
					tmp := make(map[string]interface{})
					test <- tmp
				}
				return &test
			},
			expectedSize: ptrSize + chanHdrSize + 16*(ptrSize+hmapStructSize+
				(8*(1+strHdrSize+interfaceSize)+ptrSize)),
		}, "chan_unsafe_Pointer_16": {
			getStruct: func() interface{} {
				test := make(chan unsafe.Pointer, 16)
				for i := 0; i < 16; i++ {
					var a int
					ptrToA := (unsafe.Pointer)(unsafe.Pointer(&a))
					test <- ptrToA
				}
				return &test
			},
			expectedSize: ptrSize + chanHdrSize + 16*ptrSize,
		}, "chan_[]int_16": {
			getStruct: func() interface{} {
				test := make(chan []int, 16)
				for i := 0; i < 8; i++ {
					intSlice := make([]int, 16)
					test <- intSlice
				}
				return &test
			},
			expectedSize: ptrSize + chanHdrSize + 16*sliceHdrSize + 8*16*ptrSize,
		}, "chan_func": {
			getStruct: func() interface{} {
				test := make(chan func(), 16)
				f := func() {
					fmt.Printf("Hello!")
				}
				for i := 0; i < 8; i++ {
					test <- f
				}
				return &test
			},
			expectedSize: ptrSize + chanHdrSize + 16*ptrSize,
		},
	}

	for key, tcase := range tests {
		t.Run(key, func(t *testing.T) {
			v := tcase.getStruct()
			m, err := DeepSizeof(v)
			if err != nil {
				t.Fatal(err)
			}
			var totalSize uintptr
			for _, size := range m {
				totalSize += size
			}
			expectedSize := tcase.expectedSize
			if totalSize != expectedSize {
				t.Fatalf("Expected size: %v, but got %v", expectedSize, totalSize)
			}
		})
	}
}

func TestUpdateSeenAreas(t *testing.T) {
	tests := []struct {
		seen         []block
		expectedSeen []block
		expectedSize uintptr
		update       block
	}{{
		seen: []block{
			{start: 0x100000, end: 0x100050},
		},
		expectedSeen: []block{
			{start: 0x100000, end: 0x100050},
			{start: 0x100100, end: 0x100150},
		},
		expectedSize: 0x50,
		update:       block{start: 0x100100, end: 0x100150},
	}, {
		seen: []block{
			{start: 0x100000, end: 0x100050},
		},
		expectedSeen: []block{
			{start: 0x100, end: 0x150},
			{start: 0x100000, end: 0x100050},
		},
		expectedSize: 0x50,
		update:       block{start: 0x100, end: 0x150},
	}, {
		seen: []block{
			{start: 0x100000, end: 0x100500},
		},
		expectedSeen: []block{
			{start: 0x100000, end: 0x100750},
		},
		expectedSize: 0x250,
		update:       block{start: 0x100250, end: 0x100750},
	}, {
		seen: []block{
			{start: 0x100250, end: 0x100750},
		},
		expectedSeen: []block{
			{start: 0x100000, end: 0x100750},
		},
		expectedSize: 0x250,
		update:       block{start: 0x100000, end: 0x100500},
	}, {
		seen: []block{
			{start: 0x1000, end: 0x1250},
			{start: 0x1500, end: 0x1750},
		},
		expectedSeen: []block{
			{start: 0x1000, end: 0x1750},
		},
		expectedSize: 0x2B0,
		update:       block{start: 0x1200, end: 0x1700},
	}, {
		seen: []block{
			{start: 0x1000, end: 0x1250},
			{start: 0x1500, end: 0x1750},
			{start: 0x1F50, end: 0x21A0},
		},
		expectedSeen: []block{
			{start: 0xF00, end: 0x1F00},
			{start: 0x1F50, end: 0x21A0},
		},
		expectedSize: 0xB60,
		update:       block{start: 0xF00, end: 0x1F00},
	}, {
		seen: []block{
			{start: 0x1000, end: 0x1250},
			{start: 0x1500, end: 0x1750},
			{start: 0x1F00, end: 0x2150},
		},
		expectedSeen: []block{
			{start: 0xF00, end: 0x2150},
		},
		expectedSize: 0xB60,
		update:       block{start: 0xF00, end: 0x1F00},
	}, {
		seen: []block{
			{start: 0x1000, end: 0x1250},
			{start: 0x1500, end: 0x1750},
			{start: 0x1F00, end: 0x2150},
		},
		expectedSeen: []block{
			{start: 0x1000, end: 0x1750},
			{start: 0x1F00, end: 0x2150},
		},
		expectedSize: 0x2B0,
		update:       block{start: 0x1250, end: 0x1500},
	}}

	for i, tcase := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			seen, size := updateSeenBlocks(tcase.update, tcase.seen)
			if !test.DeepEqual(seen, tcase.expectedSeen) {
				t.Fatalf("seen blocks %x for iterration %v are different than the "+
					"one expected:\n %x", seen, i, tcase.expectedSeen)
			}
			if size != tcase.expectedSize {
				t.Fatalf("Size does not match, expected 0x%x got 0x%x",
					tcase.expectedSize, size)
			}
		})
	}
}
