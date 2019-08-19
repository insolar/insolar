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

type QueueType int8

const (
	InvalidQueue QueueType = iota

	UnusedList
	ActiveList
	PollingList
)

const NoQueue QueueType = -1

func NewQueueHead(t QueueType) QueueHead {
	if t <= InvalidQueue {
		panic("illegal value")
	}
	return QueueHead{Slot{queueType: t}}
}

type QueueHead struct {
	head Slot
}

func (p *QueueHead) QueueType() QueueType {
	p.initEmpty()
	return p.head.queueType
}

func (p *QueueHead) Count() int {
	return p.head.listCount
}

func (p *QueueHead) First() *Slot {
	return p.head.Next()
}

func (p *QueueHead) Last() *Slot {
	return p.head.Prev()
}

func (p *QueueHead) IsEmpty() bool {
	return p.head.next == nil || p.head.next == &p.head
}

func (p *QueueHead) initEmpty() {
	if p.head.headDependency == nil {
		if p.head.slotType <= AnotherSlot {
			panic("illegal state")
		}
		p.head.headDependency = &p.head
		p.head.next = &p.head
		p.head.prevDependency = &p.head
	}
}

func (p *QueueHead) AddFirst(slot *Slot) {
	p.initEmpty()
	p.head.insertAsNext(slot)
	slot.headDependency = &p.head
}

func (p *QueueHead) AddLast(slot *Slot) {
	p.initEmpty()
	p.head.insertAsPrev(slot)
	slot.headDependency = &p.head
}

func (p *QueueHead) AppendAll(anotherList *QueueHead) {
	p.initEmpty()
	if anotherList.IsEmpty() {
		return
	}

	anotherHead := &anotherList.head
	next := anotherHead.next
	prev := anotherHead.prevDependency
	c := anotherHead.depCount

	anotherHead.depCount = 0
	anotherHead.next = anotherHead
	anotherHead.prevDependency = anotherHead

	p.head.prevDependency._insertAllAsNext(next, prev)
	_updateHeads(next, prev, &p.head)
	p.head.depCount += c
}

func (p *QueueHead) RemoveAll() {
	p.initEmpty()
	if p.IsEmpty() {
		return
	}

	next := p.head.next

	p.head.depCount = 0
	p.head.next = &p.head
	p.head.prevDependency = &p.head

	for next != &p.head {
		prev := next
		next = next.next
		prev.prevDependency = nil
		prev.prevDependency = nil
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
	slot.headDependency = s.headDependency
	s.headDependency.depCount++
}

func (s *Slot) insertAsPrev(slot *Slot) {
	s.prevDependency.insertAsNext(slot)
}

func (s *Slot) _insertAllAsNext(chainHead, chainTail *Slot) {
	s.ensureInList()

	chainTail.next = s.next
	chainHead.prevDependency = s.next.prevDependency

	s.next.prevDependency = chainTail
	s.next = chainHead
}

func _updateHeads(chainHead, chainTail, newHead *Slot) {
	for {
		chainHead.headDependency = newHead
		next := chainHead.next
		if next == chainTail || next == chainHead {
			return
		}
		chainHead = next
	}
}

func _updateHeadsAndCounts(chainHead, chainTail, newHead *Slot) {
	for {
		if chainHead.headDependency != nil {
			chainHead.headDependency.depCount--
		}
		chainHead.headDependency = newHead
		newHead.depCount++

		next := chainHead.next
		if next == chainTail || next == chainHead {
			return
		}
		chainHead = next
	}
}

func (s *Slot) ensureInQueue(inQueue bool) {
	next := s.nextInQueue
	prev := s.prevInQueue
	if (next == nil) != (prev == nil) {
		panic("illegal state - inconsistent")
	}
	if (next != nil) != inQueue {
		panic("illegal state")
	}
}

func (s *Slot) removeFromQueue() {
	next := s.nextInQueue
	prev := s.prevInQueue
	if (next == nil) != (prev == nil) {
		panic("illegal state - inconsistent")
	}
	if s.queueType != NoQueue {
		panic("illegal state")
	}
	if prev == nil && next == nil {
		return
	}

	s.headDependency.listCount--

	next.prevDependency = prev
	prev.next = next

	s.headDependency = nil
	s.next = nil
	s.prevDependency = nil
}

func (s *Slot) NextInQueue() *Slot {
	next := s.nextInQueue
	if next == nil || next.queueType != InvalidQueue {
		return nil
	}
	return next
}

func (s *Slot) PrevInQueue() *Slot {
	prev := s.prevInQueue
	if prev == nil || prev.queueType != InvalidQueue {
		return nil
	}
	return prev
}
