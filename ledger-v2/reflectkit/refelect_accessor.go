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

package reflectkit

import (
	"reflect"
	"unsafe"
)

type TypedReceiver interface {
	ReceiveBool(reflect.Kind, bool)
	ReceiveInt(reflect.Kind, int64)
	ReceiveUint(reflect.Kind, uint64)
	ReceiveFloat(reflect.Kind, float64)
	ReceiveComplex(reflect.Kind, complex128)
	ReceiveString(reflect.Kind, string)
	ReceiveZero(reflect.Kind)
	ReceiveNil(reflect.Kind)
	ReceiveIface(reflect.Kind, interface{})
	ReceiveElse(t reflect.Kind, v interface{}, isZero bool)
}

type FieldGetterFunc func(reflect.Value) reflect.Value
type ValueToReceiverFunc func(value reflect.Value, out TypedReceiver)
type IfaceToReceiverFunc func(value interface{}, k reflect.Kind, out TypedReceiver)
type IfaceToReceiverFactoryFunc func(t reflect.Type, checkZero bool) IfaceToReceiverFunc
type ValueToReceiverFactoryFunc func(unexported bool, t reflect.Type, checkZero bool) (bool, ValueToReceiverFunc)

func ValueToReceiverFactory(k reflect.Kind, custom IfaceToReceiverFactoryFunc) ValueToReceiverFactoryFunc {
	switch k {
	// ======== Simple values, are safe to read from unexported fields ============
	case reflect.Bool:
		return func(_ bool, _ reflect.Type, checkZero bool) (bool, ValueToReceiverFunc) {
			return false, func(value reflect.Value, outFn TypedReceiver) {
				v := value.Bool()
				if !checkZero || v {
					outFn.ReceiveBool(value.Kind(), value.Bool())
				} else {
					outFn.ReceiveZero(value.Kind())
				}
			}
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(_ bool, _ reflect.Type, checkZero bool) (bool, ValueToReceiverFunc) {
			return false, func(value reflect.Value, outFn TypedReceiver) {
				v := value.Int()
				if !checkZero || v != 0 {
					outFn.ReceiveInt(value.Kind(), v)
				} else {
					outFn.ReceiveZero(value.Kind())
				}
			}
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return func(_ bool, _ reflect.Type, checkZero bool) (bool, ValueToReceiverFunc) {
			return false, func(value reflect.Value, outFn TypedReceiver) {
				v := value.Uint()
				if !checkZero || v != 0 {
					outFn.ReceiveUint(value.Kind(), v)
				} else {
					outFn.ReceiveZero(value.Kind())
				}
			}
		}

	case reflect.Float32, reflect.Float64:
		return func(_ bool, _ reflect.Type, checkZero bool) (bool, ValueToReceiverFunc) {
			return false, func(value reflect.Value, outFn TypedReceiver) {
				v := value.Float()
				if !checkZero || v != 0 {
					outFn.ReceiveFloat(value.Kind(), v)
				} else {
					outFn.ReceiveZero(value.Kind())
				}
			}
		}

	case reflect.Complex64, reflect.Complex128:
		return func(_ bool, _ reflect.Type, checkZero bool) (bool, ValueToReceiverFunc) {
			return false, func(value reflect.Value, outFn TypedReceiver) {
				v := value.Complex()
				if !checkZero || v != 0 {
					outFn.ReceiveComplex(value.Kind(), v)
				} else {
					outFn.ReceiveZero(value.Kind())
				}
			}
		}

	case reflect.String:
		return func(_ bool, _ reflect.Type, checkZero bool) (bool, ValueToReceiverFunc) {
			return false, func(value reflect.Value, outFn TypedReceiver) {
				v := value.String()
				if !checkZero || v != "" {
					outFn.ReceiveString(value.Kind(), v)
				} else {
					outFn.ReceiveZero(value.Kind())
				}
			}
		}

	// ============ Types that require special handling for unexported fields ===========

	// nillable types fully defined at compile time
	case reflect.Func, reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan: //, reflect.UnsafePointer:
		return func(unexported bool, t reflect.Type, checkZero bool) (bool, ValueToReceiverFunc) {
			if custom != nil {
				if fn := custom(t, checkZero); fn != nil {
					return unexported, func(value reflect.Value, outFn TypedReceiver) {
						if value.IsNil() {
							outFn.ReceiveNil(k)
						} else {
							fn(value.Interface(), value.Kind(), outFn)
						}
					}
				}
			}
			return unexported, func(value reflect.Value, outFn TypedReceiver) {
				if value.IsNil() { // avoids interface nil
					outFn.ReceiveNil(k)
				} else {
					outFn.ReceiveElse(value.Kind(), value.Interface(), false)
				}
			}
		}

	// nillable types fully undefined at compile time
	case reflect.Interface:
		return func(unexported bool, t reflect.Type, checkZero bool) (b bool, getterFunc ValueToReceiverFunc) {
			if custom != nil {
				if fn := custom(t, checkZero); fn != nil {
					return unexported, func(value reflect.Value, outFn TypedReceiver) {
						v := value.Interface()
						fn(v, value.Elem().Kind(), outFn)
					}
				}
			}

			return unexported, func(value reflect.Value, outFn TypedReceiver) {
				if value.IsNil() {
					outFn.ReceiveNil(value.Kind())
				} else {
					v := value.Interface()
					outFn.ReceiveElse(value.Kind(), v, false)
				}
			}
		}

	// non-nillable
	case reflect.Struct, reflect.Array:
		return func(unexported bool, t reflect.Type, checkZero bool) (b bool, getterFunc ValueToReceiverFunc) {
			valueFn := custom(t, checkZero)
			if checkZero {
				zeroValue := reflect.Zero(t).Interface()
				return unexported, func(value reflect.Value, outFn TypedReceiver) {
					if v := value.Interface(); valueFn == nil {
						outFn.ReceiveElse(value.Kind(), v, v == zeroValue)
					} else {
						valueFn(v, value.Kind(), outFn)
					}
				}
			}
			return unexported, func(value reflect.Value, outFn TypedReceiver) {
				if v := value.Interface(); valueFn == nil {
					outFn.ReceiveElse(value.Kind(), v, false)
				} else {
					valueFn(v, value.Kind(), outFn)
				}
			}
		}

	// ============ Excluded ===================
	//reflect.UnsafePointer
	default:
		return nil
	}
}

func MakeAddressable(value reflect.Value) reflect.Value {
	if value.CanAddr() {
		return value
	}
	valueCopy := reflect.New(value.Type()).Elem()
	valueCopy.Set(value)
	return valueCopy
}

func FieldValueGetter(index int, fd reflect.StructField, useAddr bool, baseOffset uintptr) FieldGetterFunc {
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
