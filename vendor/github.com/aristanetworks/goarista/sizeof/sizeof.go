// Copyright (c) 2017 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package sizeof

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/aristanetworks/goarista/areflect"
)

// blocks are used to keep track of which memory areas were already
// been visited.
type block struct {
	start uintptr
	end   uintptr
}

func (b block) size() uintptr {
	return b.end - b.start
}

// DeepSizeof returns total memory occupied by each type for val.
// The value passed in argument must be a pointer.
func DeepSizeof(val interface{}) (map[string]uintptr, error) {
	value := reflect.ValueOf(val)
	// We want to force val to be a pointer to the original value, because if we get a copy, we
	// can get some pointers that will point back to our original value.
	if value.Kind() != reflect.Ptr {
		return nil, errors.New("cannot get the deep size of a non-pointer value")
	}
	m := make(map[string]uintptr)
	ptrsTypes := make(map[uintptr]map[string]struct{})
	sizeof(value.Elem(), m, ptrsTypes, false, block{start: uintptr(value.Pointer())}, nil)
	return m, nil
}

// Check if curBlock overlap tmpBlock
func isOverlapping(curBlock, tmpBlock block) bool {
	return curBlock.start <= tmpBlock.end && tmpBlock.start <= curBlock.end
}

func getOverlappingBlocks(curBlock block, seen []block) ([]block, int) {
	var tmp []block
	for idx, a := range seen {
		if a.start > curBlock.end {
			return tmp, idx
		}
		if isOverlapping(curBlock, a) {
			tmp = append(tmp, a)
		}
	}
	return tmp, len(seen)
}

func insertBlock(curBlock block, idxToInsert int, seen []block) []block {
	seen = append(seen, block{})
	copy(seen[idxToInsert+1:], seen[idxToInsert:])
	seen[idxToInsert] = curBlock
	return seen
}

// get the size of our current block that is not overlapping other blocks.
func getUnseenSizeOfCurrentBlock(curBlock block, overlappingBlocks []block) uintptr {
	var size uintptr
	for idx, a := range overlappingBlocks {
		if idx == 0 && curBlock.start < a.start {
			size += a.start - curBlock.start
		}
		if idx == len(overlappingBlocks)-1 {
			if curBlock.end > a.end {
				size += curBlock.end - a.end
			}
		} else {
			size += overlappingBlocks[idx+1].start - a.end
		}
	}
	return size
}

func updateSeenBlocks(curBlock block, seen []block) ([]block, uintptr) {
	if len(seen) == 0 {
		return []block{curBlock}, curBlock.size()
	}
	overlappingBlocks, idx := getOverlappingBlocks(curBlock, seen)
	if len(overlappingBlocks) == 0 {
		// No overlap, so we will insert our new block in our list.
		return insertBlock(curBlock, idx, seen), curBlock.size()
	}
	unseenSize := getUnseenSizeOfCurrentBlock(curBlock, overlappingBlocks)
	idxFirstOverlappingBlock := idx - len(overlappingBlocks)
	firstOverlappingBlock := &seen[idxFirstOverlappingBlock]
	lastOverlappingBlock := seen[idx-1]
	if firstOverlappingBlock.start > curBlock.start {
		firstOverlappingBlock.start = curBlock.start
	}
	if lastOverlappingBlock.end < curBlock.end {
		firstOverlappingBlock.end = curBlock.end
	} else {
		firstOverlappingBlock.end = lastOverlappingBlock.end
	}
	tailLen := len(seen[idx:])
	copy(seen[idxFirstOverlappingBlock+1:], seen[idx:])
	return seen[:idxFirstOverlappingBlock+1+tailLen], unseenSize
}

// Check if this current block is already fully contained in our list of seen blocks
func isKnownBlock(curBlock block, seen []block) bool {
	for _, a := range seen {
		if a.start <= curBlock.start &&
			a.end >= curBlock.end {
			// curBlock is fully contained in an other block
			// that we already know
			return true
		}
		if a.start > curBlock.start {
			// Our slice of seens block is order by pointer address.
			// That means, if curBlock was not contained in a previous known
			// block, there is no need to continue.
			return false
		}
	}
	return false
}

func sizeof(v reflect.Value, m map[string]uintptr, ptrsTypes map[uintptr]map[string]struct{},
	counted bool, curBlock block, seen []block) []block {
	if !v.IsValid() {
		return seen
	}
	vn := v.Type().String()
	vs := v.Type().Size()
	curBlock.end = vs + curBlock.start
	if counted {
		// already accounted for the size (field in struct, in array, etc)
		vs = 0
	}
	if curBlock.start != 0 {
		// A pointer can point to the same memory area than a previous pointer,
		// but its type should be different (see tests struct_5 and struct_6).
		if types, ok := ptrsTypes[curBlock.start]; ok {
			if _, ok := types[vn]; ok {
				return seen
			}
			types[vn] = struct{}{}
		} else {
			ptrsTypes[curBlock.start] = make(map[string]struct{})
		}
		if isKnownBlock(curBlock, seen) {
			// we don't want to count this size if we have a known block
			vs = 0
		} else {
			var tmpVs uintptr
			seen, tmpVs = updateSeenBlocks(curBlock, seen)
			if !counted {
				vs = tmpVs
			}
		}
	}
	switch v.Kind() {
	case reflect.Interface:
		seen = sizeof(v.Elem(), m, ptrsTypes, false, block{}, seen)
	case reflect.Ptr:
		if v.IsNil() {
			break
		}
		seen = sizeof(v.Elem(), m, ptrsTypes, false, block{start: uintptr(v.Pointer())}, seen)
	case reflect.Array:
		// get size of all elements in the array in case there are pointers
		l := v.Len()
		for i := 0; i < l; i++ {
			seen = sizeof(v.Index(i), m, ptrsTypes, true, block{}, seen)
		}
	case reflect.Slice:
		// get size of all elements in the slice in case there are pointers
		// TODO: count elements that are not accessible after reslicing
		l := v.Len()
		vLen := v.Type().Elem().Size()
		for i := 0; i < l; i++ {
			e := v.Index(i)
			eStart := uintptr(e.UnsafeAddr())
			eBlock := block{
				start: eStart,
				end:   eStart + vLen,
			}
			if !isKnownBlock(eBlock, seen) {
				vs += vLen
				seen = sizeof(e, m, ptrsTypes, true, eBlock, seen)
			}
		}
		capStart := uintptr(v.Pointer()) + (v.Type().Elem().Size() * uintptr(v.Len()))
		capEnd := uintptr(v.Pointer()) + (v.Type().Elem().Size() * uintptr(v.Cap()))
		capBlock := block{start: capStart, end: capEnd}
		if isKnownBlock(capBlock, seen) {
			break
		}
		var tmpSize uintptr
		seen, tmpSize = updateSeenBlocks(capBlock, seen)
		vs += tmpSize
	case reflect.Map:
		if v.IsNil() {
			break
		}
		var tmpSize uintptr
		if tmpSize, seen = sizeofmap(v, seen); tmpSize == 0 {
			// we saw this map
			break
		}
		vs += tmpSize
		for _, k := range v.MapKeys() {
			kv := v.MapIndex(k)
			seen = sizeof(k, m, ptrsTypes, true, block{}, seen)
			seen = sizeof(kv, m, ptrsTypes, true, block{}, seen)
		}
	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			vf := areflect.ForceExport(v.Field(i))
			seen = sizeof(vf, m, ptrsTypes, true, block{}, seen)
		}
	case reflect.String:
		str := v.String()
		strHdr := (*reflect.StringHeader)(unsafe.Pointer(&str))
		tmpSize := uintptr(strHdr.Len)
		strBlock := block{start: strHdr.Data, end: strHdr.Data + tmpSize}
		if isKnownBlock(strBlock, seen) {
			break
		}
		seen, tmpSize = updateSeenBlocks(strBlock, seen)
		vs += tmpSize
	case reflect.Chan:
		var tmpSize uintptr
		tmpSize, seen = sizeofChan(v, m, ptrsTypes, seen)
		vs += tmpSize
	}
	if vs != 0 {
		m[vn] += vs
	}
	return seen
}

//go:linkname typesByString reflect.typesByString
func typesByString(s string) []unsafe.Pointer

func sizeofmap(v reflect.Value, seen []block) (uintptr, []block) {
	// get field typ *rtype of our Value v and store it in an interface
	var ti interface{} = v.Type()
	tp := (*unsafe.Pointer)(unsafe.Pointer(&ti))
	// we know that this pointer rtype point at the begining of struct
	// mapType defined in /go/src/reflect/type.go, so we can change the underlying
	// type of the interface to be a pointer to runtime.maptype because it as the
	// exact same definition as reflect.mapType.
	*tp = typesByString("*runtime.maptype")[0]
	maptypev := reflect.ValueOf(ti)
	maptypev = reflect.Indirect(maptypev)
	// now we can access field bucketsize in struct maptype
	bucketsize := maptypev.FieldByName("bucketsize").Uint()
	// get hmap
	var m interface{} = v.Interface()
	hmap := (*unsafe.Pointer)(unsafe.Pointer(&m))
	*hmap = typesByString("*runtime.hmap")[0]

	hmapv := reflect.ValueOf(m)
	// account for the size of the hmap, buckets and oldbuckets
	hmapv = reflect.Indirect(hmapv)
	mapBlock := block{
		start: hmapv.UnsafeAddr(),
		end:   hmapv.UnsafeAddr() + hmapv.Type().Size(),
	}
	// is it a map we already saw?
	if isKnownBlock(mapBlock, seen) {
		return 0, seen
	}
	seen, _ = updateSeenBlocks(mapBlock, seen)
	B := hmapv.FieldByName("B").Uint()
	oldbuckets := hmapv.FieldByName("oldbuckets").Pointer()
	buckets := hmapv.FieldByName("buckets").Pointer()
	noverflow := int16(hmapv.FieldByName("noverflow").Uint())
	nb := 2
	if B == 0 {
		nb = 1
	}
	size := uint64((nb << B)) * bucketsize
	if noverflow != 0 {
		size += uint64(noverflow) * bucketsize
	}
	seen, _ = updateSeenBlocks(block{start: buckets, end: buckets + uintptr(size)},
		seen)
	// As defined in /go/src/runtime/hashmap.go in struct hmap, oldbuckets is the
	// previous bucket array that is half the size of the current one. We need to
	// also take that in consideration since there is still a pointer to this previous bucket.
	if oldbuckets != 0 {
		tmp := (2 << (B - 1)) * bucketsize
		size += tmp
		seen, _ = updateSeenBlocks(block{
			start: oldbuckets,
			end:   oldbuckets + uintptr(tmp),
		}, seen)
	}
	return hmapv.Type().Size() + uintptr(size), seen
}

func getSliceToChanBuffer(buff unsafe.Pointer, buffLen uint, dataSize uint) []byte {
	var slice []byte
	sliceHdr := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	sliceHdr.Len = int(buffLen * dataSize)
	sliceHdr.Cap = sliceHdr.Len
	sliceHdr.Data = uintptr(buff)
	return slice
}

func sizeofChan(v reflect.Value, m map[string]uintptr, ptrsTypes map[uintptr]map[string]struct{},
	seen []block) (uintptr, []block) {
	var c interface{} = v.Interface()
	hchan := (*unsafe.Pointer)(unsafe.Pointer(&c))
	*hchan = typesByString("*runtime.hchan")[0]

	hchanv := reflect.ValueOf(c)
	hchanv = reflect.Indirect(hchanv)
	chanBlock := block{
		start: hchanv.UnsafeAddr(),
		end:   hchanv.UnsafeAddr() + hchanv.Type().Size(),
	}
	// is it a chan we already saw?
	if isKnownBlock(chanBlock, seen) {
		return 0, seen
	}
	seen, _ = updateSeenBlocks(chanBlock, seen)
	elemType := unsafe.Pointer(hchanv.FieldByName("elemtype").Pointer())
	buff := unsafe.Pointer(hchanv.FieldByName("buf").Pointer())
	buffLen := hchanv.FieldByName("dataqsiz").Uint()
	elemSize := uint16(hchanv.FieldByName("elemsize").Uint())
	seen, _ = updateSeenBlocks(block{
		start: uintptr(buff),
		end:   uintptr(buff) + uintptr(buffLen*uint64(elemSize)),
	}, seen)

	buffSlice := getSliceToChanBuffer(buff, uint(buffLen), uint(elemSize))
	recvx := hchanv.FieldByName("recvx").Uint()
	qcount := hchanv.FieldByName("qcount").Uint()

	var tmp interface{}
	eface := (*struct {
		typ unsafe.Pointer
		ptr unsafe.Pointer
	})(unsafe.Pointer(&tmp))
	eface.typ = elemType
	for i := uint64(0); buffLen > 0 && i < qcount; i++ {
		idx := (recvx + i) % buffLen
		// get the pointer to the data inside the chan buffer.
		elem := unsafe.Pointer(&buffSlice[uint64(elemSize)*idx])
		eface.ptr = elem
		ev := reflect.ValueOf(tmp)
		var blk block
		k := ev.Kind()
		if k == reflect.Ptr || k == reflect.Chan || k == reflect.Map || k == reflect.Func {
			// let's say our chan is a chan *whatEver, or chan chan whatEver or
			// chan map[whatEver]whatEver. In this case elemType will
			// be either of type *whatEver, chan whatEver or map[whatEver]whatEver
			// but what we set eface.ptr = elem above, we make it point to a pointer
			// to where the data is sotred in the buffer of our channel.
			// So the interface tmp would look like:
			// chan *whatEver -> (type=*whatEver, ptr=**whatEver)
			// chan chan whatEver -> (type= chan whatEver, ptr=*chan whatEver)
			// chan map[whatEver]whatEver -> (type= map[whatEver]whatEver{},
			// ptr=*map[whatEver]whatEver)
			// So we need to take the ptr which is stored into the buffer and replace
			// the ptr to the data of our interface tmp.
			ptr := (*unsafe.Pointer)(elem)
			eface.ptr = *ptr
			ev = reflect.ValueOf(tmp)
			ev = reflect.Indirect(ev)
			blk.start = uintptr(*ptr)
		}
		// It seems that when the chan is of type chan *whatEver, the type in eface
		// will be whatEver and not *whatEver, but ev.Kind() will be a reflect.ptr.
		// So if k is a reflect.Ptr (i.e. a pointer) to a struct, then we want to take
		// the size of the struct into account because
		// vs := v.Type().Size() will return us the size of the struct and not the size
		// of the pointer that is in the channel's buffer.
		seen = sizeof(ev, m, ptrsTypes, true && k != reflect.Ptr, blk, seen)
	}
	return hchanv.Type().Size() + uintptr(uint64(elemSize)*buffLen), seen
}
