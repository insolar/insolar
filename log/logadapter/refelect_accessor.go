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

type fieldValueGetterFunc func(value reflect.Value) interface{}

var fieldValueGetters = map[reflect.Kind]func(unexported bool, t reflect.Type, nilZero bool) (bool, fieldValueGetterFunc){
	// ======== Simple values, are safe to read from unexported fields ============
	reflect.Bool: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && !value.Bool() {
				return nil
			}
			return value.Bool()
		}
	},
	reflect.Int: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Int() == 0 {
				return nil
			}
			return int(value.Int())
		}
	},
	reflect.Int8: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Int() == 0 {
				return nil
			}
			return int8(value.Int())
		}
	},
	reflect.Int16: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Int() == 0 {
				return nil
			}
			return int16(value.Int())
		}
	},
	reflect.Int32: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Int() == 0 {
				return nil
			}
			return int32(value.Int())
		}
	},
	reflect.Int64: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Int() == 0 {
				return nil
			}
			return value.Int()
		}
	},
	reflect.Uint: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Uint() == 0 {
				return nil
			}
			return uint(value.Uint())
		}
	},
	reflect.Uint8: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Uint() == 0 {
				return nil
			}
			return uint8(value.Uint())
		}
	},
	reflect.Uint16: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Uint() == 0 {
				return nil
			}
			return uint16(value.Uint())
		}
	},
	reflect.Uint32: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Uint() == 0 {
				return nil
			}
			return uint32(value.Uint())
		}
	},
	reflect.Uint64: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Uint() == 0 {
				return nil
			}
			return value.Uint()
		}
	},
	reflect.Uintptr: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Uint() == 0 {
				return nil
			}
			return uintptr(value.Uint())
		}
	},
	reflect.Float32: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Float() == 0.0 {
				return nil
			}
			return float32(value.Float())
		}
	},
	reflect.Float64: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Float() == 0.0 {
				return nil
			}
			return value.Float()
		}
	},
	reflect.Complex64: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Complex() == 0.0 {
				return nil
			}
			return complex64(value.Complex())
		}
	},
	reflect.Complex128: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Complex() == 0.0 {
				return nil
			}
			return value.Complex()
		}
	},
	reflect.String: func(_ bool, _ reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			if nilZero && value.Len() == 0 {
				return nil
			}
			return value.String()
		}
	},

	// ============ Special handling for unexported fields ===========

	reflect.Ptr: func(unexported bool, t reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		te := t.Elem()
		if te.Kind() == reflect.String {
			return false, func(value reflect.Value) interface{} {
				if value.IsNil() || nilZero && value.Len() == 0 {
					return nil
				}
				return value.Elem().String()
			}
		}
		return defaultObjFieldGetterFactory(unexported, t, nilZero)
	},

	reflect.Func: func(unexported bool, t reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
		if t.NumIn() == 0 && t.NumOut() == 1 && t.Out(0).Kind() == reflect.String {
			return unexported, func(value reflect.Value) interface{} {
				if value.IsNil() {
					return nil
				}
				fn := value.Interface().(func() string)
				return fn()
			}
		}
		return unexported, reflect.Value.Interface
	},

	reflect.Interface: func(unexported bool, _ reflect.Type, _ bool) (b bool, getterFunc fieldValueGetterFunc) {
		return unexported, func(value reflect.Value) interface{} {
			if value.IsNil() {
				return nil
			}
			iv := value.Interface()
			switch vv := iv.(type) {
			case func() string:
				return vv()
			default:
				vv, _ = tryDefaultValuePrepare(vv)
				return vv
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

var prepareObjTypes = []struct {
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

func defaultObjFieldGetterFactory(unexported bool, t reflect.Type, nilZero bool) (bool, fieldValueGetterFunc) {
	for _, f := range prepareObjTypes {
		if !t.Implements(f.t) {
			continue
		}

		fn := f.fn
		if t.Kind() == reflect.Struct {
			if nilZero {
				zeroValue := reflect.Zero(t).Interface()
				return unexported, func(value reflect.Value) interface{} {
					iv := value.Interface()
					if iv == zeroValue {
						return nil
					}
					vv, _ := fn(iv)
					return vv
				}
			}

			return unexported, func(value reflect.Value) interface{} {
				vv, _ := fn(value.Interface())
				return vv
			}
		}

		return unexported, func(value reflect.Value) interface{} {
			if value.IsNil() {
				return nil
			}
			vv, _ := fn(value.Interface())
			return vv
		}
	}
	return unexported, reflect.Value.Interface
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
