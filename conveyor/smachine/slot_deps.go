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

package smachine

type DependencyHead struct {
	head Slot
}

func (p *DependencyHead) isHead(slot *Slot) bool {
	if p == nil {
		panic("illegal state")
	}
	return slot == &p.head
}

/*
-----------------------------------
Slot methods to support linked list
-----------------------------------
*/

//func (s *Slot) insertAsNext(slot *Slot) {
//	slot.ensureNotInList()
//	s._insertAllAsNext(slot, slot)
//	slot.headDependency = s.headDependency
//	s.headDependency.depCount++
//}
//
//func (s *Slot) insertAsPrev(slot *Slot) {
//	s.prevDependency.insertAsNext(slot)
//}
//
//func (s *Slot) _insertAllAsNext(chainHead, chainTail *Slot) {
//	s.ensureInList()
//
//	chainTail.next = s.next
//	chainHead.prevDependency = s.next.prevDependency
//
//	s.next.prevDependency = chainTail
//	s.next = chainHead
//}
//
//func _updateHeads(chainHead, chainTail, newHead *Slot) {
//	for {
//		chainHead.headDependency = newHead
//		next := chainHead.next
//		if next == chainTail || next == chainHead {
//			return
//		}
//		chainHead = next
//	}
//}
//
//func _updateHeadsAndCounts(chainHead, chainTail, newHead *Slot) {
//	for {
//		if chainHead.headDependency != nil {
//			chainHead.headDependency.depCount--
//		}
//		chainHead.headDependency = newHead
//		newHead.depCount++
//
//		next := chainHead.next
//		if next == chainTail || next == chainHead {
//			return
//		}
//		chainHead = next
//	}
//}
//
//func (s *Slot) remove() {
//	if s.headDependency == nil || s.headDependency == s {
//		return
//	}
//
//	next := s.next
//	prev := s.prevDependency
//	s.headDependency.depCount--
//
//	next.prevDependency = prev
//	prev.next = next
//
//	s.headDependency = nil
//	s.next = nil
//	s.prevDependency = nil
//}
//

func (s *Slot) NextDependency() *Slot {
	next := s.nextDependency
	if next == nil || s.headDependency.isHead(next) {
		return nil
	}
	return next
}

func (s *Slot) PrevDependency() *Slot {
	prev := s.prevDependency
	if prev == nil || s.headDependency.isHead(prev) {
		return nil
	}
	return prev
}
