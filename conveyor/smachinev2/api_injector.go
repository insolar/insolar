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
	"fmt"
	"reflect"
	"strings"
)

type InjectProviderFunc func(id string) interface{}

func tryInject(provider InjectProviderFunc, id string, v reflect.Value) error {
	if provider == nil {
		panic("illegal value: provider")
	}
	k := v.Kind()
	if k != reflect.Ptr {
		panic("illegal value: not pointer")
	}
	v = v.Elem()
	if !v.CanSet() {
		panic("illegal value: readonly")
	}

	vt := v.Type()
	switch vt.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if !v.IsNil() {
			return nil
		}
	default:
		if !v.CanInterface() {
			break
		}

		zeroValue := reflect.Zero(vt).Interface()
		if v.Interface() != zeroValue {
			return nil
		}
	}

	val := provider(id)
	if val == nil {
		return fmt.Errorf("dependency is missing: id=%s", id)
	}

	dv := reflect.ValueOf(val)
	switch dv.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if !dv.IsNil() {
			break
		}
		fallthrough
	case reflect.Invalid:
		return fmt.Errorf("dependency is missing: id=%s", id)
	}

	dt := dv.Type()
	if !dt.AssignableTo(vt) {
		return fmt.Errorf("dependency type mismatch: id=%s expected=%v provided=%v", id, vt, dt)
	}
	v.Set(dv)
	return nil
}

func TryInject(provider InjectProviderFunc, id string, varRef interface{}) error {
	if varRef == nil {
		panic("illegal value: value")
	}
	if id == "" {
		panic("illegal value: id")
	}
	return tryInject(provider, id, reflect.ValueOf(varRef))
}

func InjectById(provider InjectProviderFunc, id string, varRef interface{}) {
	err := TryInject(provider, id, varRef)
	if err != nil {
		panic(err)
	}
}

func Inject(provider InjectProviderFunc, varRef interface{}) {
	if varRef == nil {
		panic("illegal value: value")
	}
	v := reflect.ValueOf(varRef)
	s := strings.TrimLeft(v.Type().String(), "*")
	err := tryInject(provider, s, v)
	if err != nil {
		panic(err)
	}
}
