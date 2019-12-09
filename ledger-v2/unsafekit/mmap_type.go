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
	"reflect"

	"github.com/insolar/insolar/longbits"
)

type MemoryMapModel uint8

const (
	MemoryMapIncompatible MemoryMapModel = iota
	MemoryModelDepended
	MemoryModelIndependent
)

type Unwrapper interface {
	Unwrap(longbits.ByteString) interface{}
	MemoryModelIndependent() bool
}

type MMapType struct {
	t                reflect.Type
	modelIndependent bool
}

func NewMMapType(t reflect.Type) (MMapType, bool) {
	if mm := MemoryModelDependencyOf(t); mm == MemoryMapIncompatible {
		return MMapType{}, false
	} else {
		return MMapType{t, mm == MemoryModelIndependent}, true
	}
}

func MustMMapType(t reflect.Type, mustBeIndependent bool) MMapType {
	switch mt, ok := NewMMapType(t); {
	case !ok:
		panic("illegal value - type must be memory-mappable")
	case !mustBeIndependent || mt.MemoryModelIndependent():
		return mt
	default:
		panic("illegal value - type must be memory-mappable and memory-model independent")
	}
}

type MMapSliceType struct {
	t                reflect.Type
	modelIndependent bool
}

func NewMMapSliceType(t reflect.Type) (MMapSliceType, bool) {
	if t.Kind() != reflect.Slice {
		panic("illegal value")
	}
	if mm := MemoryModelDependencyOf(t.Elem()); mm == MemoryMapIncompatible {
		return MMapSliceType{}, false
	} else {
		return MMapSliceType{t, mm == MemoryModelIndependent}, true
	}
}

func MustMMapSliceType(t reflect.Type, mustBeIndependent bool) MMapSliceType {
	switch mt, ok := NewMMapSliceType(t); {
	case !ok:
		panic("illegal value - type must be memory-mappable")
	case !mustBeIndependent || mt.MemoryModelIndependent():
		return mt
	default:
		panic("illegal value - type must be memory-mappable and memory-model independent")
	}
}

func MemoryModelDependencyOf(t reflect.Type) MemoryMapModel {
	switch k := t.Kind(); k {
	case reflect.Bool, reflect.Int8, reflect.Uint8:
		return MemoryModelIndependent

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return MemoryModelDepended

	case reflect.Complex64, reflect.Complex128:
		// there is no information of memory representation of complex types
		return MemoryMapIncompatible

	case reflect.Array:
		return MemoryModelDependencyOf(t.Elem())

	case reflect.Uintptr: // same as Pointer
		return MemoryMapIncompatible
	default:
		return MemoryMapIncompatible

	case reflect.Struct:
		//
	}

	fieldCount := t.NumField()
	endOfPrev := uintptr(0)
	for i := 0; i < fieldCount; i++ {
		f := t.Field(i)
		if f.Offset != endOfPrev {
			return MemoryModelDepended
		}

		switch mt := MemoryModelDependencyOf(f.Type); mt {
		case MemoryModelIndependent:
			//
		case MemoryModelDepended, MemoryMapIncompatible:
			return mt
		default:
			panic("unexpected")
		}

		endOfPrev += f.Type.Size()
	}
	if endOfPrev != t.Size() {
		return MemoryModelDepended
	}
	return MemoryModelIndependent
}

func (v MMapType) IsZero() bool {
	return v.t == nil
}

func (v MMapType) ReflectType() reflect.Type {
	if v.t == nil {
		panic("illegal state")
	}
	return v.t
}

func (v MMapType) MemoryModelIndependent() bool {
	return v.modelIndependent
}

func (v MMapType) Size() int {
	return int(v.t.Size())
}

func (v MMapType) SliceOf() MMapSliceType {
	return MMapSliceType{reflect.SliceOf(v.t), v.modelIndependent}
}

func (v MMapType) Unwrap(s longbits.ByteString) interface{} {
	return UnwrapAs(s, v)
}

func (v MMapSliceType) IsZero() bool {
	return v.t == nil
}

func (v MMapSliceType) ReflectType() reflect.Type {
	if v.t == nil {
		panic("illegal state")
	}
	return v.t
}

func (v MMapSliceType) MemoryModelIndependent() bool {
	return v.modelIndependent
}

func (v MMapSliceType) Elem() MMapType {
	return MMapType{v.t.Elem(), v.modelIndependent}
}

func (v MMapSliceType) ElemReflectType() reflect.Type {
	return v.t.Elem()
}

func (v MMapSliceType) Unwrap(s longbits.ByteString) interface{} {
	return UnwrapAsSlice(s, v)
}
