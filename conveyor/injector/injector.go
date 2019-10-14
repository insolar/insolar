package injector

import (
	"fmt"
	"reflect"
	"strings"
)

func NewDependencyInjector(target interface{}, parent DependencyRegistry, copyParent DependencyRegistryFunc) DependencyInjector {
	if target == nil {
		panic("illegal value")
	}
	return DependencyInjector{target: target, parent: parent, copyParent: copyParent}
}

var _ DependencyRegistry = &DependencyInjector{}

type DependencyInjector struct {
	parent     DependencyRegistry
	copyParent DependencyRegistryFunc
	target     interface{}
	resolved   map[string]interface{}
}

func (p *DependencyInjector) IsZero() bool {
	return p.target == nil
}

func (p *DependencyInjector) IsEmpty() bool {
	return len(p.resolved) == 0
}

func (p *DependencyInjector) Count() int {
	return len(p.resolved)
}

func (p *DependencyInjector) ResolveAndPut(overrides map[string]interface{}) {
	for id, v := range overrides {
		if id == "" {
			panic("illegal value")
		}
		if dp, ok := v.(DependencyProviderFunc); ok {
			v = dp(p.target, id, p)
		}
		p.putResolved(id, v)
	}
}

func (p *DependencyInjector) FindDependency(id string) (interface{}, bool) {
	if id == "" {
		panic("illegal value")
	}
	if v, ok := p.resolved[id]; ok { // allows nil values
		return v, true
	}

	v, ok, _ := p.getParent(id)
	return v, ok
}

func (p *DependencyInjector) GetResolvedDependency(id string) (interface{}, bool) {
	if id == "" {
		panic("illegal value")
	}
	return p.getResolved(id)
}

func (p *DependencyInjector) PutResolvedDependency(id string, v interface{}) {
	if id == "" {
		panic("illegal value")
	}
	if _, ok := v.(DependencyProviderFunc); ok {
		panic("illegal value")
	}
	p.putResolved(id, v)
}

func (p *DependencyInjector) getParent(id string) (interface{}, bool, bool) {
	if p.copyParent != nil {
		if v, ok := p.copyParent(id); ok {
			return v, true, true
		}
	}
	if p.parent != nil {
		if v, ok := p.parent.FindDependency(id); ok {
			return v, true, false
		}
	}
	return nil, false, false
}

func (p *DependencyInjector) getResolved(id string) (interface{}, bool) {
	if v, ok := p.resolved[id]; ok { // allows nil values
		return v, true
	}

	if v, ok, cp := p.getParent(id); ok {
		if dp, ok := v.(DependencyProviderFunc); ok {
			v = dp(p.target, id, p)
			p.putResolved(id, v)
		} else if cp {
			p.putResolved(id, v)
		}
		return v, true
	}

	return nil, false
}

func (p *DependencyInjector) putResolved(id string, v interface{}) {
	if p.resolved == nil {
		p.resolved = make(map[string]interface{})
	}
	p.resolved[id] = v
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
	t := reflect.Indirect(reflect.ValueOf(p.target))
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
		id, ok := quickTag(`inject:"`, sf.Tag)
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
			typeName = strings.TrimLeft(tt.String(), "*")
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

func quickTag(prefix string, tag reflect.StructTag) (string, bool) {
	if len(tag) < len(prefix)+1 {
		return "", false
	}
	if strings.HasPrefix(string(tag), prefix) && strings.HasSuffix(string(tag), `"`) {
		return string(tag[len(prefix) : len(tag)-1]), true
	}
	return "", false
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
	if isSet {
		return fmt.Errorf("dependency is set: id=%s expectedType=%v", id, vt)
	}

	if id != "" {
		if p.resolveNameAndSet(id, v, vt, isNillable) {
			return nil
		}
	} else if p.resolveTypeAndSet(
		strings.TrimLeft(vt.String(), "*"),
		"", v, vt, isNillable) {
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

	switch val, ok := p.getResolved(n); {
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

func (p *DependencyInjector) CopyAsRegistryWithParent() ReadOnlyContainer {
	if p.parent == nil {
		return p.CopyAsRegistryNoParent()
	}
	return NewRegistryWithParent(p.parent.FindDependency, p.resolved)
}

func (p *DependencyInjector) CopyAsRegistryNoParent() ReadOnlyContainer {
	return NewRegistry(p.resolved)
}
