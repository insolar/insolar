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

// +build mips mips64 ppc64 s390x

package unsafekit

import "unsafe"

const BigEndian = true

func init() {
	bytes := [4]byte{1, 2, 3, 4}
	v := *((*uint32)((unsafe.Pointer)(&bytes)))
	if v != 0x01020304 {
		panic("FATAL - expected BigEndian memory architecture")
	}
}
