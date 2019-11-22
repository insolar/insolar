///
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
///

package logadapter

import (
	"fmt"
	"reflect"
	"unsafe"
)

type fieldValueGetterFunc func(value reflect.Value) (interface{}, bool)
type fieldValueGetterFactoryFunc func(unexported bool, t reflect.Type, checkZero bool) (bool, fieldValueGetterFunc)

var fieldValueGetters = map[reflect.Kind]fieldValueGetterFactoryFunc{
	// ======== Simple values, are safe to read from unexported fields ============
	reflect.Bool: func(_ bool, _ reflect.Type, _ bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			return value.Bool(), !value.Bool()
		}
	},
	reflect.Int: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Int()
			return int(v), checkZero && v == 0
		}
	},
	reflect.Int8: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Int()
			return int8(v), checkZero && v == 0
		}
	},
	reflect.Int16: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Int()
			return int16(v), checkZero && v == 0
		}
	},
	reflect.Int32: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Int()
			return int32(v), checkZero && v == 0
		}
	},
	reflect.Int64: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Int()
			return v, checkZero && v == 0
		}
	},
	reflect.Uint: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Uint()
			return uint(v), checkZero && v == 0
		}
	},
	reflect.Uint8: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Uint()
			return uint8(v), checkZero && v == 0
		}
	},
	reflect.Uint16: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Uint()
			return uint16(v), checkZero && v == 0
		}
	},
	reflect.Uint32: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Uint()
			return uint32(v), checkZero && v == 0
		}
	},
	reflect.Uint64: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Uint()
			return v, checkZero && v == 0
		}
	},
	reflect.Uintptr: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Uint()
			return uintptr(v), checkZero && v == 0
		}
	},
	reflect.Float32: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Float()
			return float32(v), checkZero && v == 0.0
		}
	},
	reflect.Float64: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Float()
			return v, checkZero && v == 0.0
		}
	},
	reflect.Complex64: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Complex()
			return complex64(v), checkZero && v == 0.0
		}
	},
	reflect.Complex128: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.Complex()
			return v, checkZero && v == 0.0
		}
	},
	reflect.String: func(_ bool, _ reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) (interface{}, bool) {
			v := value.String()
			return v, checkZero && len(v) == 0
		}
	},

	// ============ Special handling for unexported fields ===========

	reflect.Ptr: func(unexported bool, t reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		te := t.Elem()
		if te.Kind() == reflect.String {
			return false, func(value reflect.Value) (interface{}, bool) {
				if value.IsNil() {
					return nil, checkZero
				}
				v := value.Elem().String()
				return v, checkZero && len(v) == 0
			}
		}
		return defaultObjFieldGetterFactory(unexported, t, checkZero)
	},

	reflect.Func: func(unexported bool, t reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
		if t.NumIn() == 0 && t.NumOut() == 1 && t.Out(0).Kind() == reflect.String {
			return unexported, func(value reflect.Value) (interface{}, bool) {
				if value.IsNil() {
					return nil, checkZero
				}
				fn := value.Interface().(func() string)
				v := fn()
				return v, checkZero && len(v) == 0
			}
		}
		return unexported, func(value reflect.Value) (interface{}, bool) {
			if value.IsNil() { // avoids interface nil
				return nil, checkZero
			}
			return value.Interface(), false
		}
	},

	reflect.Interface: func(unexported bool, _ reflect.Type, checkZero bool) (b bool, getterFunc fieldValueGetterFunc) {
		return unexported, func(value reflect.Value) (interface{}, bool) {
			if value.IsNil() { // avoids interface nil
				return nil, checkZero
			}
			iv := value.Interface()
			switch vv := iv.(type) {
			case func() string:
				v := vv()
				return v, checkZero && len(v) == 0
			default:
				vv, _ = tryDefaultValuePrepare(vv)
				return vv, checkZero && vv == nil
			}
		}
	},

	reflect.Struct: defaultObjFieldGetterFactory,
	reflect.Array:  defaultObjFieldGetterFactory,
	reflect.Map:    defaultObjFieldGetterFactory,
	reflect.Slice:  defaultObjFieldGetterFactory,
	reflect.Chan:   defaultObjFieldGetterFactory,

	// ============ Excluded ===================
	//reflect.UnsafePointer
}

var prepareObjTypes = []struct { // MUST be ordered, not map
	t  reflect.Type
	fn func(interface{}) (interface{}, bool)
}{
	{reflect.TypeOf((*LogStringer)(nil)).Elem(), func(value interface{}) (interface{}, bool) {
		if vv, ok := value.(LogStringer); ok {
			return vv.LogString(), true
		}
		return value, false
	}},
	{reflect.TypeOf((*fmt.Stringer)(nil)).Elem(), func(value interface{}) (interface{}, bool) {
		if vv, ok := value.(fmt.Stringer); ok {
			return vv.String(), true
		}
		return value, false
	}},
}

func defaultStrValuePrepare(iv interface{}) (string, bool) {
	switch vv := iv.(type) {
	case string:
		return vv, true
	case *string:
		if vv == nil {
			return "", false
		}
		return *vv, true
	case func() string:
		return vv(), true
	}
	if vv, ok := tryDefaultValuePrepare(iv); ok {
		return vv.(string), true
	}
	return "", false
}

func defaultObjFieldGetterFactory(unexported bool, t reflect.Type, checkZero bool) (bool, fieldValueGetterFunc) {
	for _, f := range prepareObjTypes {
		if !t.Implements(f.t) {
			continue
		}
		return unexported, _defaultObjFieldGetterFactory(t, checkZero, f.fn)
	}

	return unexported, _defaultObjFieldGetterFactoryNoConv(t, checkZero)
}

func _defaultObjFieldGetterFactory(t reflect.Type, checkZero bool, fn func(interface{}) (interface{}, bool)) fieldValueGetterFunc {
	if t.Kind() == reflect.Struct {
		if checkZero {
			zeroValue := reflect.Zero(t).Interface()
			return func(value reflect.Value) (interface{}, bool) {
				iv := value.Interface()
				if iv == zeroValue {
					return iv, true
				}
				vv, _ := fn(iv)
				return vv, vv == nil
			}
		}

		return func(value reflect.Value) (interface{}, bool) {
			vv, _ := fn(value.Interface())
			return vv, false
		}
	}

	return func(value reflect.Value) (interface{}, bool) {
		if value.IsNil() { // avoids interface nil
			return nil, checkZero
		}
		vv, _ := fn(value.Interface())
		return vv, checkZero && vv == nil
	}
}

func _defaultObjFieldGetterFactoryNoConv(t reflect.Type, checkZero bool) fieldValueGetterFunc {
	if t.Kind() == reflect.Struct {
		if checkZero {
			zeroValue := reflect.Zero(t).Interface()
			return func(value reflect.Value) (interface{}, bool) {
				iv := value.Interface()
				if iv == zeroValue {
					return iv, true
				}
				return iv, iv == nil
			}
		}

		return func(value reflect.Value) (interface{}, bool) {
			return value.Interface(), false
		}
	}

	return func(value reflect.Value) (interface{}, bool) {
		if value.IsNil() { // avoids interface nil
			return nil, checkZero
		}
		vv := value.Interface()
		return vv, checkZero && vv == nil
	}
}

func tryDefaultValuePrepare(iv interface{}) (interface{}, bool) {
	for _, f := range prepareObjTypes {
		if vv, ok := f.fn(iv); ok {
			return vv, true
		}
	}
	return iv, false
}

func getFieldGetter(index int, fd reflect.StructField, useAddr bool, baseOffset uintptr) func(reflect.Value) reflect.Value {
	if !useAddr {
		return func(value reflect.Value) reflect.Value {
			return value.Field(index)
		}
	}

	fieldOffset := fd.Offset + baseOffset
	fieldType := fd.Type

	return func(value reflect.Value) reflect.Value {
		return offsetFieldGetter(value, fieldOffset, fieldType)
	}
}

func offsetFieldGetter(v reflect.Value, fieldOffset uintptr, fieldType reflect.Type) reflect.Value {
	return reflect.NewAt(fieldType, unsafe.Pointer(v.UnsafeAddr()+fieldOffset)).Elem()
}
