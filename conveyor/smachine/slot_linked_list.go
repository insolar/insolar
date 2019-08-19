///
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
///

package smachine

type DependencyType int8

const (
	AnotherSlot DependencyType = iota

	UnusedList
	ActiveList
	PollingList
	RecoveryList
)

const NoDependency DependencyType = -1

func NewListHead(t DependencyType) ListHead {
	if t <= AnotherSlot {
		panic("illegal value")
	}
	return ListHead{Slot{slotType: t}}
}

type ListHead struct {
	head Slot
}

func (p *ListHead) Category() DependencyType {
	p.initEmpty()
	return p.head.slotType
}

func (p *ListHead) Count() int {
	return p.head.depCount
}

func (p *ListHead) First() *Slot {
	return p.head.Next()
}

func (p *ListHead) Last() *Slot {
	return p.head.Prev()
}

func (p *ListHead) IsEmpty() bool {
	return p.head.next == nil || p.head.next == &p.head
}

func (p *ListHead) initEmpty() {
	if p.head.head == nil {
		if p.head.slotType <= AnotherSlot {
			panic("illegal state")
		}
		p.head.head = &p.head
		p.head.next = &p.head
		p.head.prev = &p.head
	}
}

func (p *ListHead) AddFirst(slot *Slot) {
	p.initEmpty()
	p.head.insertAsNext(slot)
	slot.head = &p.head
}

func (p *ListHead) AddLast(slot *Slot) {
	p.initEmpty()
	p.head.insertAsPrev(slot)
	slot.head = &p.head
}

func (p *ListHead) AppendAll(anotherList *ListHead) {
	p.initEmpty()
	if anotherList.IsEmpty() {
		return
	}

	anotherHead := &anotherList.head
	next := anotherHead.next
	prev := anotherHead.prev
	c := anotherHead.depCount

	anotherHead.depCount = 0
	anotherHead.next = anotherHead
	anotherHead.prev = anotherHead

	p.head.prev._insertAllAsNext(next, prev)
	_updateHeads(next, prev, &p.head)
	p.head.depCount += c
}

func (p *ListHead) RemoveAll() {
	p.initEmpty()
	if p.IsEmpty() {
		return
	}

	next := p.head.next

	p.head.depCount = 0
	p.head.next = &p.head
	p.head.prev = &p.head

	for next != &p.head {
		prev := next
		next = next.next
		prev.prev = nil
		prev.prev = nil
		prev.head = nil

		if prev == next {
			break
		}
	}
}

/*
-----------------------------------
Slot methods to support linked list
-----------------------------------
*/

func (s *Slot) insertAsNext(slot *Slot) {
	slot.ensureNotInList()
	s._insertAllAsNext(slot, slot)
	slot.head = s.head
	s.head.depCount++
}

func (s *Slot) insertAsPrev(slot *Slot) {
	s.prev.insertAsNext(slot)
}

func (s *Slot) _insertAllAsNext(chainHead, chainTail *Slot) {
	s.ensureInList()

	chainTail.next = s.next
	chainHead.prev = s.next.prev

	s.next.prev = chainTail
	s.next = chainHead
}

func _updateHeads(chainHead, chainTail, newHead *Slot) {
	for {
		chainHead.head = newHead
		next := chainHead.next
		if next == chainTail || next == chainHead {
			return
		}
		chainHead = next
	}
}

func _updateHeadsAndCounts(chainHead, chainTail, newHead *Slot) {
	for {
		if chainHead.head != nil {
			chainHead.head.depCount--
		}
		chainHead.head = newHead
		newHead.depCount++

		next := chainHead.next
		if next == chainTail || next == chainHead {
			return
		}
		chainHead = next
	}
}

func (s *Slot) remove() {
	if s.head == nil || s.head == s {
		return
	}

	next := s.next
	prev := s.prev
	s.head.depCount--

	next.prev = prev
	prev.next = next

	s.head = nil
	s.next = nil
	s.prev = nil
}

func (s *Slot) DependencyType() DependencyType {
	if s.head == nil {
		return NoDependency
	}
	return s.head.slotType
}

func (s *Slot) Next() *Slot {
	next := s.next
	if next == nil || next == s.head || next.head == next {
		return nil
	}
	return next
}

func (s *Slot) Prev() *Slot {
	prev := s.prev
	if prev == nil || prev == s.head || prev.head == prev {
		return nil
	}
	return prev
}
