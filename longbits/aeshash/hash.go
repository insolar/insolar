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

package aeshash

import (
	"reflect"
	"runtime"
	"unsafe"
)

// Hash hashes the given string using the algorithm used by Go's hash tables
func Str(s string) uint {
	return StrWithSeed(s, 0)
}

func StrWithSeed(s string, seed uint) uint {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return hash(s, sh.Data, sh.Len, seed)
}

// Hash hashes the given slice using the algorithm used by Go's hash tables
func Slice(b []byte) uint {
	return SliceWithSeed(b, 0)
}

func SliceWithSeed(b []byte, seed uint) uint {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return hash(b, sh.Data, sh.Len, seed)
}

func hash(keepAlive interface{}, data uintptr, len int, seed uint) uint {
	r := aeshash(data, uintptr(seed), uintptr(len))
	runtime.KeepAlive(keepAlive)
	return uint(r)
}

func aeshash(pData, hSeed, sLen uintptr) uintptr
