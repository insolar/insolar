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

package injector

import (
	"fmt"
	"reflect"
	"strings"
)

func GetDefaultInjectionId(v interface{}) string {
	return GetDefaultInjectionIdByType(reflect.TypeOf(v))
}

func GetDefaultInjectionIdByType(vt reflect.Type) string {
	return strings.TrimLeft(vt.String(), "*")
}

func NewDependencyInjector(target interface{}, globalParent DependencyRegistry, localParent DependencyRegistryFunc) DependencyInjector {
	resolver := NewDependencyResolver(target, globalParent, localParent, nil)
	return NewDependencyInjectorFor(&resolver)
}

func NewDependencyInjectorFor(resolver *DependencyResolver) DependencyInjector {
	if resolver == nil || resolver.IsZero() {
		panic("illegal value")
	}
	return DependencyInjector{resolver}
}

type DependencyInjector struct {
	resolver *DependencyResolver
}

func (p *DependencyInjector) IsZero() bool {
	return p.resolver.IsZero()
}

func (p *DependencyInjector) IsEmpty() bool {
	return p.resolver.IsEmpty()
}

func (p *DependencyInjector) MustInject(varRef interface{}) {
	if err := p.tryInjectVar("", varRef); err != nil {
		panic(err)
	}
}

func (p *DependencyInjector) MustInjectById(id string, varRef interface{}) {
	if err := p.tryInjectVar(id, varRef); err != nil {
		panic(err)
	}
}

func (p *DependencyInjector) Inject(varRef interface{}) error {
	return p.tryInjectVar("", varRef)
}

func (p *DependencyInjector) InjectById(id string, varRef interface{}) error {
	if id == "" {
		panic("illegal value")
	}
	return p.tryInjectVar(id, varRef)
}

func (p *DependencyInjector) InjectAll() error {
	t := reflect.Indirect(reflect.ValueOf(p.resolver.Target()))
	if t.Kind() != reflect.Struct {
		panic("illegal value")
	}
	if !t.CanSet() {
		panic("illegal value: readonly")
	}
	tt := t.Type()
	typeName := ""

	for i := 0; i < tt.NumField(); i++ {
		sf := tt.Field(i)
		id, ok := sf.Tag.Lookup("inject")
		if !ok {
			continue
		}

		fv := t.Field(i)
		switch isNillable, isSet := p.check(fv, sf.Type); {
		case isSet:
			return fmt.Errorf("dependency is set: target=%v field=%s", tt.String(), sf.Name)
		case id != "":
			if p.resolveNameAndSet(id, fv, sf.Type, isNillable) {
				continue
			}
		case typeName == "":
			typeName = GetDefaultInjectionIdByType(tt)
			fallthrough
		default:
			if p.resolveTypeAndSet(typeName, sf.Name, fv, sf.Type, isNillable) {
				continue
			}
		}

		return fmt.Errorf("dependency is missing: target=%v field=%s id=%s expectedType=%v", tt, sf.Name, id, sf.Type)
	}

	return nil
}

func (p *DependencyInjector) tryInjectVar(id string, varRef interface{}) error {
	if varRef == nil {
		panic("illegal value")
	}

	v := reflect.ValueOf(varRef)
	switch {
	case v.Kind() != reflect.Ptr:
		panic("illegal value: not a reference")
	case v.IsNil():
		panic("illegal value: nil reference")
	case v.CanSet() || v.CanAddr():
		panic("illegal value: must be a literal reference")
	}
	v = v.Elem()

	vt := v.Type()
	isNillable, isSet := p.check(v, vt)

	switch {
	case isSet:
		return fmt.Errorf("dependency is set: id=%s expectedType=%v", id, vt)
	case id != "":
		if p.resolveNameAndSet(id, v, vt, isNillable) {
			return nil
		}
	case p.resolveTypeAndSet(GetDefaultInjectionIdByType(vt), "", v, vt, isNillable):
		return nil
	}

	return fmt.Errorf("dependency is missing: id=%s expectedType=%v", id, vt)
}

func (p *DependencyInjector) check(v reflect.Value, vt reflect.Type) (bool, bool) {
	if !v.CanSet() {
		panic("illegal value: readonly")
	}

	switch vt.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true, !v.IsNil()
	default:
		zeroValue := reflect.Zero(vt).Interface()
		return false, v.Interface() != zeroValue
	}
}

func (p *DependencyInjector) resolveTypeAndSet(typeName, fieldName string, v reflect.Value, vt reflect.Type, nillable bool) bool {
	if p.resolveNameAndSet(typeName, v, vt, nillable) {
		return true
	}
	idx := strings.LastIndexByte(typeName, '.')
	if idx >= 0 {
		if p.resolveNameAndSet(typeName[idx+1:], v, vt, nillable) {
			return true
		}
	}

	if fieldName == "" {
		return false
	}
	typeName = typeName + "." + fieldName

	if p.resolveNameAndSet(typeName, v, vt, nillable) {
		return true
	}
	if idx >= 0 {
		if p.resolveNameAndSet(typeName[idx+1:], v, vt, nillable) {
			return true
		}
	}
	return false
}

func (p *DependencyInjector) resolveNameAndSet(n string, v reflect.Value, vt reflect.Type, nillable bool) bool {
	if len(n) == 0 {
		return false
	}

	switch val, ok := p.resolver.getResolved(n); {
	case !ok:
		return false
	case nillable && val == nil:
		return true
	default:
		dv := reflect.ValueOf(val)
		dt := dv.Type()
		if !dt.AssignableTo(vt) {
			return false // fmt.Errorf("dependency type mismatch: id=%s expected=%v provided=%v", id, vt, dt)
		}
		v.Set(dv)
		return true
	}
}
