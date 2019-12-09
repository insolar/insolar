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

const PtrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const

type MemoryModelSupported uint8

const (
	_ MemoryModelSupported = iota
	LittleEndianSupported
	BigEndianSupported
	EndianIndependent // LittleEndianSupported | BigEndianSupported
)

func IsCompatibleMemoryModel(v MemoryModelSupported) bool {
	switch v {
	case EndianIndependent:
		return true
	case LittleEndianSupported:
		return !BigEndian
	case BigEndianSupported:
		return BigEndian
	default:
		panic("illegal value")
	}
}
