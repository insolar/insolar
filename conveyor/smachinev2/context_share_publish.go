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

import "reflect"

// this structure provides isolation of shared data to avoid SM being retained via SharedDataLink
type uniqueAlias struct {
	valueType reflect.Type
}

func (p *slotContext) Share(data interface{}, flags ShareDataFlags) SharedDataLink {
	p.ensureAtLeast(updCtxInit)
	switch data.(type) {
	case nil:
		panic("illegal value")
	case dependencyKey, *slotAliases, *uniqueAlias:
		panic("illegal value")
	case SharedDataLink, *SharedDataLink:
		panic("illegal value - SharedDataLink can't be shared")
	}

	switch {
	case flags&ShareDataUnbound != 0: // ShareDataDirect is irrelevant
		return SharedDataLink{SlotLink{}, data, flags}
	case flags&ShareDataDirect != 0:
		return SharedDataLink{p.s.NewLink(), data, flags}
	default:
		alias := &uniqueAlias{reflect.TypeOf(data)}
		if !p.s.registerBoundAlias(alias, data) {
			panic("impossible")
		}
		return SharedDataLink{p.s.NewLink(), alias, flags}
	}
}

func (p *slotContext) Publish(key, data interface{}) bool {
	p.ensureAtLeast(updCtxInit)
	switch key.(type) {
	case nil:
		panic("illegal value")
	case dependencyKey, *slotAliases, *uniqueAlias:
		panic("illegal value")
	}

	switch data.(type) {
	case nil:
		panic("illegal value")
	case dependencyKey, *slotAliases, *uniqueAlias:
		panic("illegal value")
	}
	return p.s.registerBoundAlias(key, data)
}

func (p *slotContext) Unpublish(key interface{}) bool {
	p.ensureAtLeast(updCtxInit)
	return p.s.unregisterBoundAlias(key)
}

func (p *slotContext) GetPublished(key interface{}) interface{} {
	p.ensureAtLeast(updCtxInit)
	if v, ok := p.s.machine.getPublished(key); ok {
		return v
	}
	return nil
}

func (p *machineCallContext) GetPublished(key interface{}) interface{} {
	p.ensureValid()
	if v, ok := p.m.getPublished(key); ok {
		return v
	}
	return nil
}

func (m *SlotMachine) getPublished(key interface{}) (interface{}, bool) {
	switch key.(type) {
	case nil:
		return nil, false
	case dependencyKey, *slotAliases, *uniqueAlias:
		return nil, false
	}
	return m.localRegistry.Load(key)
}

// Provides external access to published data.
// But SharedDataLink can only be accessed when is unbound.
// The bool value indicates presence of a valid key, but value can be nil when access is not allowed.
func (m *SlotMachine) GetPublished(key interface{}) (interface{}, bool) {
	if v, ok := m.getPublished(key); ok {
		// unwrap unbound values
		// but slot-bound values can NOT be accessed outside of a slot machine
		switch sdl := v.(type) {
		case dependencyKey, *slotAliases, *uniqueAlias:
			return nil, false
		case SharedDataLink:
			if sdl.IsUnbound() {
				return sdl.getData(), true
			}
			return nil, true
		case *SharedDataLink:
			if sdl.IsUnbound() {
				return sdl.getData(), true
			}
			return nil, true
		default:
			return v, true
		}
	}
	return nil, false
}

func _asSharedDataLink(v interface{}) SharedDataLink {
	switch d := v.(type) {
	case SharedDataLink:
		return d
	case *SharedDataLink:
		return *d
	default:
		return SharedDataLink{}
	}
}

func (p *slotContext) GetPublishedLink(key interface{}) SharedDataLink {
	return _asSharedDataLink(p.GetPublished(key))
}

func (p *machineCallContext) GetPublishedLink(key interface{}) SharedDataLink {
	return _asSharedDataLink(p.GetPublished(key))
}
