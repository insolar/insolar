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

/* ------- Slot-dependant aliases and mappings ------------- */

type slotIdKey SlotID

type slotAliasesValue struct {
	//owner *Slot
	keys []interface{}
}

//type slotAliasesMap struct {
//	//owner *Slot
//	keys map[interface{}
//}

// ONLY to be used by a holder of a slot
func (s *Slot) registerBoundAlias(k, v interface{}) bool {
	if k == nil {
		panic("illegal value")
	}

	m := &s.machine.localRegistry
	if _, loaded := m.LoadOrStore(k, v); loaded {
		return false
	}

	var key interface{} = slotIdKey(s.GetSlotID())

	switch isa, ok := m.Load(key); {
	case !ok:
		isa, _ = m.LoadOrStore(key, &slotAliasesValue{ /* owner:s */ })
		fallthrough
	default:
		sa := isa.(*slotAliasesValue)
		sa.keys = append(sa.keys, k)
		s.slotFlags |= slotHasAliases
	}

	if sar := s.machine.config.SlotAliasRegistry; sar != nil {
		if ga, ok := k.(globalAliasKey); ok {
			return sar.PublishAlias(ga.key, s.NewLink())
		}
	}

	return true
}

// ONLY to be used by a holder of a slot
func (s *Slot) unregisterBoundAlias(k interface{}) bool {
	if k == nil {
		panic("illegal value")
	}

	switch keyExists, wasUnpublished, _ := s.machine.unpublishUnbound(k); {
	case !keyExists:
		return false
	case wasUnpublished:
		return true
	}
	return s.machine._unregisterSlotBoundAlias(s.GetSlotID(), k)
}

// ONLY to be used by a holder of a slot
func (m *SlotMachine) unregisterBoundAliases(id SlotID) {
	mm := &m.localRegistry // SAFE for concurrent use
	var key interface{} = slotIdKey(id)

	if isa, ok := mm.Load(key); ok {
		sa := isa.(*slotAliasesValue)
		mm.Delete(key)

		sar := m.config.SlotAliasRegistry
		for _, k := range sa.keys {
			mm.Delete(k)

			if sar != nil {
				if ga, ok := k.(globalAliasKey); ok {
					sar.UnpublishAlias(ga.key)
				}
			}
		}
	}
}

// ONLY to be used by a holder of a slot
func (m *SlotMachine) _unregisterSlotBoundAlias(slotID SlotID, k interface{}) bool {
	var key interface{} = slotIdKey(slotID)

	if isa, loaded := m.localRegistry.Load(key); loaded {
		sa := isa.(*slotAliasesValue)

		for i, kk := range sa.keys {
			if k == kk {
				m.localRegistry.Delete(k)
				if sar := m.config.SlotAliasRegistry; sar != nil {
					if ga, ok := k.(globalAliasKey); ok {
						sar.UnpublishAlias(ga.key)
					}
				}

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
