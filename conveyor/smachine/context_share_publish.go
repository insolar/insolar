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

func (p *slotContext) UnpublishAll() {
	p.ensureAtLeast(updCtxInit)
	p.s.unregisterBoundAliases()
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
			if sdl != nil && sdl.IsUnbound() {
				return sdl.getData(), true
			}
			return nil, true
		default:
			return v, true
		}
	}
	return nil, false
}

func (m *SlotMachine) TryPublish(key, data interface{}) (interface{}, bool) {
	switch key.(type) {
	case nil:
		panic("illegal value")
	case dependencyKey, *slotAliases, *uniqueAlias:
		panic("illegal value")
	}

	switch sdl := data.(type) {
	case dependencyKey, *slotAliases, *uniqueAlias:
		panic("illegal value")
	case SharedDataLink:
		if !sdl.IsUnbound() {
			panic("illegal value")
		}
	case *SharedDataLink:
		if sdl == nil || !sdl.IsUnbound() {
			panic("illegal value")
		}
	}

	v, loaded := m.localRegistry.LoadOrStore(key, data)
	return v, !loaded
}

// WARNING! USE WITH CAUTION. Interfering with published names may be unexpected by SM.
// This method can unpublish keys published by SMs, but it is not able to always do it in a right ways.
// As a result - a hidden registry of names published by SM can become inconsistent and will
// cause an SM to remove key(s) if they were later published by another SM.
func (m *SlotMachine) TryUnsafeUnpublish(key interface{}) (keyExists, wasUnpublished bool) {
	switch key.(type) {
	case nil:
		panic("illegal value")
	case dependencyKey, *slotAliases, *uniqueAlias:
		panic("illegal value")
	}

	// Lets try to make it right

	switch keyExists, wasUnpublished, v := m.unpublishUnbound(key); {
	case !keyExists:
		return false, false
	case wasUnpublished:
		return true, true
	default:
		var valueOwner SlotLink

		switch sdl := v.(type) {
		case SharedDataLink:
			valueOwner = sdl.link
		case *SharedDataLink:
			if sdl != nil {
				valueOwner = sdl.link
			}
		}

		// This is the most likely case ... yet it doesn't cover all the cases
		if valueOwner.IsValid() && m._unregisterSlotBoundAlias(valueOwner.SlotID(), key) {
			return true, true
		}
	}

	// as there are no more options to do it right - then do it wrong
	m.localRegistry.Delete(key)
	return true, true
}

func (m *SlotMachine) unpublishUnbound(k interface{}) (keyExists, wasUnpublished bool, value interface{}) {
	if v, ok := m.localRegistry.Load(k); !ok {
		return false, false, nil
	} else {
		switch sdl := v.(type) {
		case SharedDataLink:
			if sdl.IsUnbound() {
				m.localRegistry.Delete(k)
				return true, true, v
			}
		case *SharedDataLink:
			if sdl != nil && sdl.IsUnbound() {
				m.localRegistry.Delete(k)
				return true, true, v
			}
		}
		return true, false, v
	}
}

// ONLY to be used by a holder of a slot
func (m *SlotMachine) _unregisterSlotBoundAlias(slotID SlotID, k interface{}) bool {
	var key interface{} = slotID

	if isa, loaded := m.localRegistry.Load(key); loaded {
		sa := isa.(*slotAliases)

		for i, kk := range sa.keys {
			if k == kk {
				m.localRegistry.Delete(k)
				switch last := len(sa.keys) - 1; {
				case last == 0:
					m.localRegistry.Delete(key)
				case i < last:
					copy(sa.keys[i:], sa.keys[i+1:])
					fallthrough
				default:
					sa.keys = sa.keys[:last]
				}
				return true
			}
		}
	}
	return false
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
