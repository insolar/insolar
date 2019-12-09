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
	"unsafe"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/reference"
)

// WARNING! You MUST make sure that (v) stays alive while the resulting longbits.ByteString is in use.
func WrapLocalRef(v *reference.Local) longbits.ByteString {
	if v == nil {
		return ""
	}
	return WrapPtr(uintptr((unsafe.Pointer)(v)), unsafe.Sizeof(*v))
}

// WARNING! You MUST make sure that (v) stays alive while the resulting longbits.ByteString is in use.
func WrapGlobalRef(v *reference.Global) longbits.ByteString {
	if v == nil {
		return ""
	}
	return WrapPtr(uintptr((unsafe.Pointer)(v)), unsafe.Sizeof(*v))
}

// WARNING! This function has different guarantees based on (s) origin:
// 1) When (s) is made by wrapping another type - it satisfies Unsafe Rule (1) Conversion of a *T1 to Pointer to *T2.
//    You are safe.
//
// 2) When (s) is made by wrapping []byte or string - it violates Unsafe Rule (6) Conversion of SliceHeader/StringHeader
//    And YOU MUST make sure that the origin stays alive while the result is in use.
//
func UnwrapAsLocalRef(s longbits.ByteString) *reference.Local {
	switch len(s) {
	case 0:
		return nil
	case reference.LocalBinarySize:
		return (*reference.Local)(_unwrapUnsafe(s))
	default:
		panic("illegal value")
	}
}

// WARNING! This function has different guarantees based on (s) origin:
// 1) When (s) is made by wrapping another type - it satisfies Unsafe Rule (1) Conversion of a *T1 to Pointer to *T2.
//    You are safe.
//
// 2) When (s) is made by wrapping []byte or string - it violates Unsafe Rule (6) Conversion of SliceHeader/StringHeader
//    And YOU MUST make sure that the origin stays alive while the result is in use.
//
func UnwrapAsGlobalRef(s longbits.ByteString) *reference.Global {
	switch len(s) {
	case 0:
		return nil
	case reference.GlobalBinarySize:
		return (*reference.Global)(_unwrapUnsafe(s))
	default:
		panic("illegal value")
	}
}
