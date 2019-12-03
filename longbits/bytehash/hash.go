//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package bytehash

import (
	"reflect"
	"runtime"
	"unsafe"
)

// Hash hashes the given string using the algorithm used by Go's hash tables
func HashStr(s string) uint32 {
	return HashStrWithSeed(s, 0)
}

func HashStrWithSeed(s string, seed uint) uint32 {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return hash(sh.Data, sh.Len, seed)
}

// Hash hashes the given string using the algorithm used by Go's hash tables
func HashSlice(b []byte) uint32 {
	return HashSliceWithSeed(b, 0)
}

func HashSliceWithSeed(b []byte, seed uint) uint32 {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return hash(sh.Data, sh.Len, seed)
}

func hash(data uintptr, len int, seed uint) uint32 {
	b := dataHeader{data, len}
	r := uint32(aeshashstr(noescape(unsafe.Pointer(&b)), uintptr(seed)))
	runtime.KeepAlive(b)
	return r
}

type dataHeader struct {
	Data uintptr
	Len  int
}

func aeshashstr(p unsafe.Pointer /* *dataHeader */, seed uintptr) uintptr

// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input.  noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}
