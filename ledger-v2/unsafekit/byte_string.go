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

package unsafekit

import (
	"hash"
	"io"
	"reflect"
	"unsafe"

	"github.com/insolar/insolar/longbits"
)

// WARNING! The given array MUST be immutable
// WARNING! This method violates unsafe pointer-conversion rules.
// You MUST make sure that (b) stays alive while the resulting ByteString is in use.
func WrapBytes(b []byte) longbits.ByteString {
	if len(b) == 0 {
		return longbits.EmptyByteString
	}
	return wrapUnsafe(b)
}

// WARNING! The given struct MUST be immutable. Expects struct ptr.
// WARNING! This method violates unsafe pointer-conversion rules.
// You MUST make sure that (v) stays alive while the resulting ByteString is in use.
func WrapStruct(v interface{}) longbits.ByteString {
	vt := reflect.ValueOf(v)
	if vt.Kind() != reflect.Ptr {
		panic("illegal value")
	}
	vt = vt.Elem()
	if vt.Kind() != reflect.Struct {
		panic("illegal value")
	}
	if !vt.CanAddr() {
		panic("illegal value")
	}
	return wrapUnsafePtr(vt.Pointer(), vt.Type().Size())
}

// WARNING! The given address MUST be of an immutable object.
// WARNING! This method violates unsafe pointer-conversion rules.
// You MUST make sure that (p) stays alive while the resulting ByteString is in use.
func WrapPtr(p uintptr, size uintptr) longbits.ByteString {
	return wrapUnsafePtr(p, size)
}

func UnwrapAs(v longbits.ByteString, vt reflect.Type) interface{} {
	if int(vt.Size()) != len(v) {
		panic("illegal value")
	}
	return reflect.NewAt(vt, _unwrapUnsafe(v)).Interface()
}

func Hash(v longbits.ByteString, h hash.Hash) hash.Hash {
	unwrapUnsafe(v, func(b []byte) uintptr {
		_, _ = h.Write(b)
		return 0
	})
	return h
}

func WriteTo(v longbits.ByteString, w io.Writer) (n int64, err error) {
	unwrapUnsafe(v, func(b []byte) uintptr {
		nn := 0
		nn, err = w.Write(b)
		n = int64(nn)
		return 0
	})
	return
}

// WARNING! This method violates unsafe pointer-conversion rules.
// You MUST make sure that (b) stays alive while the resulting ByteString is in use.
func wrapUnsafe(b []byte) longbits.ByteString {
	pSlice := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	var res longbits.ByteString
	pString := (*reflect.StringHeader)(unsafe.Pointer(&res))

	pString.Data = pSlice.Data
	pString.Len = pSlice.Len

	return res
}

// WARNING! This method violates unsafe pointer-conversion rules.
// You MUST make sure that (p) stays alive while the resulting ByteString is in use.
func wrapUnsafePtr(p uintptr, size uintptr) longbits.ByteString {
	var res longbits.ByteString
	pString := (*reflect.StringHeader)(unsafe.Pointer(&res))

	pString.Data = p
	pString.Len = int(size)

	return res
}

// WARNING! This method violates unsafe pointer-conversion rules.
// You MUST make sure that (p) stays alive while the resulting Pointer is in use.
func _unwrapUnsafe(s longbits.ByteString) unsafe.Pointer {
	return unsafe.Pointer(_unwrapUnsafeUintptr(s))
}

func _unwrapUnsafeUintptr(s longbits.ByteString) uintptr {
	if len(s) == 0 {
		return 0
	}
	return ((*reflect.StringHeader)(unsafe.Pointer(&s))).Data
}

func unwrapUnsafe(s longbits.ByteString, fn func([]byte) uintptr) uintptr {
	return KeepAliveWhile(unsafe.Pointer(&s), func(p unsafe.Pointer) uintptr {
		pString := (*reflect.StringHeader)(p)

		var b []byte
		pSlice := (*reflect.SliceHeader)(unsafe.Pointer(&b))

		pSlice.Data = pString.Data
		pSlice.Len = pString.Len
		pSlice.Cap = pString.Len

		r := fn(b)
		//*pSlice = reflect.SliceHeader{}

		return r
	})
}
