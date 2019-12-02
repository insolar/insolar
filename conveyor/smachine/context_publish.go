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

package smachine

import (
	"reflect"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/pulse"
)

type globalAliasKey struct {
	key interface{}
}

func isValidPublishValue(data interface{}) bool {
	switch data.(type) {
	case nil, dependencyKey, slotIdKey, *slotAliasesValue, *uniqueAliasKey, globalAliasKey:
		return false
	}
	return true
}

func isValidPublishKey(key interface{}) bool {
	switch key.(type) {
	case nil, dependencyKey, slotIdKey, *slotAliasesValue, *uniqueAliasKey, globalAliasKey:
		return false
	case bool, int8, int16, int32, int64, int, uint8, uint16, uint32, uint64, uint, uintptr:
		return true
	case float32, float64, complex64, complex128, string:
		return true
	case longbits.ByteString, longbits.Bits64, longbits.Bits128, longbits.Bits224, longbits.Bits256, longbits.Bits512:
		return true
	case pulse.Number, SlotID:
		return true
	default:
		// have to go for reflection
		switch tt := reflect.TypeOf(key).Kind(); {
		case tt <= reflect.Array: // literals
			return tt > reflect.Invalid
		case tt >= reflect.String: // String, Struct, UnsafePointer
			return tt <= reflect.UnsafePointer
		case tt == reflect.Ptr:
			return true
		default: // Chan, Func, Interface, Map, Slice
			return false
		}
	}
}

func ensurePublishValue(data interface{}) {
	if !isValidPublishValue(data) {
		panic("illegal value")
	}
}

func ensureShareValue(data interface{}) {
	if !isValidPublishValue(data) {
		panic("illegal value")
	}
	switch data.(type) {
	case SharedDataLink, *SharedDataLink:
		panic("illegal value - SharedDataLink can't be shared")
	}
}

func ensurePublishKey(key interface{}) {
	if !isValidPublishKey(key) {
		panic("illegal value")
	}
}
