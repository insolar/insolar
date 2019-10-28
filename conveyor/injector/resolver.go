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

func NewDependencyResolver(target interface{}, globalParent DependencyRegistry, localParent DependencyRegistryFunc,
	onPutFn func(id string, v interface{}, from DependencyOrigin)) DependencyResolver {
	if target == nil {
		panic("illegal value")
	}
	return DependencyResolver{target: target, globalParent: globalParent, localParent: localParent, onPutFn: onPutFn}
}

type DependencyOrigin uint8

const (
	DependencyFromLocal DependencyOrigin = 1 << iota
	DependencyFromProvider
)

type DependencyResolver struct {
	globalParent DependencyRegistry
	localParent  DependencyRegistryFunc
	target       interface{}
	resolved     map[string]interface{}
	onPutFn      func(id string, v interface{}, from DependencyOrigin)
}

func (p *DependencyResolver) IsZero() bool {
	return p.target == nil
}

func (p *DependencyResolver) IsEmpty() bool {
	return len(p.resolved) == 0
}

func (p *DependencyResolver) Count() int {
	return len(p.resolved)
}

func (p *DependencyResolver) Target() interface{} {
	return p.target
}

func (p *DependencyResolver) ResolveAndReplace(overrides map[string]interface{}) {
	for id, v := range overrides {
		if id == "" {
			panic("illegal value")
		}
		p.resolveAndPut(id, v, DependencyFromLocal)
	}
}

func (p *DependencyResolver) ResolveAndMerge(values map[string]interface{}) {
	for id, v := range values {
		if id == "" {
			panic("illegal value")
		}
		if _, ok := p.resolved[id]; ok {
			continue
		}
		p.resolveAndPut(id, v, DependencyFromLocal)
	}
}

func (p *DependencyResolver) FindDependency(id string) (interface{}, bool) {
	if id == "" {
		panic("illegal value")
	}
	if v, ok := p.resolved[id]; ok { // allows nil values
		return v, true
	}

	v, ok, _ := p.getFromParent(id)
	return v, ok
}

func (p *DependencyResolver) GetResolvedDependency(id string) (interface{}, bool) {
	if id == "" {
		panic("illegal value")
	}
	return p.getResolved(id)
}

//func (p *DependencyResolver) PutResolvedDependency(id string, v interface{}) {
//	if id == "" {
//		panic("illegal value")
//	}
//	if _, ok := v.(DependencyProviderFunc); ok {
//		panic("illegal value")
//	}
//	p.putResolved(id, v)
//}

func (p *DependencyResolver) getFromParent(id string) (interface{}, bool, DependencyOrigin) {
	if p.localParent != nil {
		if v, ok := p.localParent(id); ok {
			return v, true, DependencyFromLocal
		}
	}
	if p.globalParent != nil {
		if v, ok := p.globalParent.FindDependency(id); ok {
			return v, true, 0
		}
	}
	return nil, false, 0
}

func (p *DependencyResolver) getResolved(id string) (interface{}, bool) {
	if v, ok := p.resolved[id]; ok { // allows nil values
		return v, true
	}
	if v, ok, from := p.getFromParent(id); ok {
		return p.resolveAndPut(id, v, from), true
	}

	return nil, false
}

func (p *DependencyResolver) resolveAndPut(id string, v interface{}, from DependencyOrigin) interface{} {
	if dp, ok := v.(DependencyProviderFunc); ok {
		p._putResolved(id, nil) // guard for resolve loop
		v = dp(p.target, id, p.GetResolvedDependency)
		p.putResolved(id, v, from|DependencyFromProvider)
	} else if from|DependencyFromLocal != 0 {
		p.putResolved(id, v, from)
	}
	return v
}

func (p *DependencyResolver) putResolved(id string, v interface{}, from DependencyOrigin) {
	p._putResolved(id, v)
	if p.onPutFn != nil {
		p.onPutFn(id, v, from)
	}
}

func (p *DependencyResolver) _putResolved(id string, v interface{}) {
	if p.resolved == nil {
		p.resolved = make(map[string]interface{})
	}
	p.resolved[id] = v
}

func (p *DependencyResolver) CopyAsRegistryWithParent() ReadOnlyContainer {
	if p.globalParent == nil {
		return p.CopyAsRegistryNoParent()
	}
	return NewRegistryWithParent(p.globalParent.FindDependency, p.resolved)
}

func (p *DependencyResolver) CopyAsRegistryNoParent() ReadOnlyContainer {
	return NewRegistry(p.resolved)
}

func (p *DependencyResolver) CopyResolved() map[string]interface{} {
	n := len(p.resolved)
	if n == 0 {
		return nil
	}
	result := make(map[string]interface{}, n)

	for id, v := range p.resolved {
		result[id] = v
	}

	return result
}

func (p *DependencyResolver) Flush() map[string]interface{} {
	n := len(p.resolved)
	if n == 0 {
		return nil
	}
	m := p.resolved
	p.resolved = nil
	return m
}
