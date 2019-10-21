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

type slotAliases struct {
	//owner *Slot
	keys []interface{}
}

// ONLY to be used by a holder of a slot
func (s *Slot) registerBoundAlias(k, v interface{}) bool {
	if k == nil {
		panic("illegal value")
	}

	m := &s.machine.localRegistry
	if _, loaded := m.LoadOrStore(k, v); loaded {
		return false
	}
	isa, _ := m.LoadOrStore(s.GetSlotID(), &slotAliases{ /* owner:s */ })
	sa := isa.(*slotAliases)
	sa.keys = append(sa.keys, k)
	s.slotFlags |= slotHasAliases

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
func (s *Slot) unregisterBoundAliases() {
	m := &s.machine.localRegistry // SAFE for concurrent use
	var key interface{} = s.GetSlotID()

	if isa, ok := m.Load(key); ok {
		sa := isa.(*slotAliases)
		m.Delete(key)

		for _, k := range sa.keys {
			m.Delete(k)
		}
	}
}
