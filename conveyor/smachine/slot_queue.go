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

	AnotherSlotQueue

	UnusedSlots
	ActiveSlots
	PollingSlots
)

const NoQueue QueueType = -1

func NewSlotQueue(t QueueType) SlotQueue {
	if t <= InvalidQueue {
		panic("illegal value")
	}
	qh := SlotQueue{QueueHead: QueueHead{queueType: t}}
	qh.head = &qh.slot
	return qh
}

type SlotQueue struct {
	QueueHead
	slot Slot
}

func (p *SlotQueue) AppendAll(anotherQueue *SlotQueue) {
	p.QueueHead.AppendAll(&anotherQueue.QueueHead)
}

type QueueHead struct {
	head      *Slot
	queueType QueueType
	count     int
}

func (p *QueueHead) QueueType() QueueType {
	p.initEmpty()
	return p.queueType
}

func (p *QueueHead) Count() int {
	return p.count
}

func (p *QueueHead) First() *Slot {
	return p.head.QueueNext()
}

func (p *QueueHead) Last() *Slot {
	return p.head.QueuePrev()
}

func (p *QueueHead) IsZero() bool {
	return p.head.nextInQueue == nil
}

func (p *QueueHead) IsEmpty() bool {
	return p.head.nextInQueue == nil || p.head.nextInQueue.isQueueHead()
}

func (p *QueueHead) initEmpty() {
	if p.head.queue == nil {
		p.head.nextInQueue = p.head
		p.head.prevInQueue = p.head
		p.head.queue = p
	}
}

func (p *QueueHead) AddFirst(slot *Slot) {
	p.initEmpty()
	slot.ensureNotInQueue()

	p.head.nextInQueue._addQueuePrev(slot, slot)
	slot.queue = p
	p.count++
}

func (p *QueueHead) AddLast(slot *Slot) {
	p.initEmpty()
	slot.ensureNotInQueue()

	p.head._addQueuePrev(slot, slot)
	slot.queue = p
	p.count++
}

func (p *QueueHead) AppendAll(anotherQueue *QueueHead) {
	p.initEmpty()
	if anotherQueue.IsEmpty() {
		return
	}

	next := anotherQueue.head.nextInQueue
	prev := anotherQueue.head.prevInQueue

	c := anotherQueue.count

	anotherQueue.count = 0
	anotherQueue.head.nextInQueue = anotherQueue.head
	anotherQueue.head.prevInQueue = anotherQueue.head

	for n := next; n != anotherQueue.head; n = n.nextInQueue {
		n.queue = p
	}

	p.head._addQueuePrev(next, prev)
	p.count += c
}

func (p *QueueHead) RemoveAll() {
	p.initEmpty()
	if p.IsEmpty() {
		return
	}

	next := p.head.nextInQueue
	p.count = 0

	p.head.nextInQueue = p.head
	p.head.prevInQueue = p.head

	for next != p.head {
		prev := next
		next = next.nextInQueue

		prev.nextInQueue = nil
		prev.prevInQueue = nil
		prev.queue = nil

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

func (s *Slot) isQueueHead() bool {
	return s.queue != nil && s == s.queue.head
}

func (s *Slot) vacateQueueHead() {
	if s.queue == nil || s.queue.head != s || s.nextInQueue != s || s.prevInQueue != s {
		panic("illegal state")
	}

	s.queue.head = nil
	s.queue = nil
	s.nextInQueue = nil
	s.prevInQueue = nil
}

func (s *Slot) makeQueueHead() {
	s.ensureNotInQueue()

	s.queue = &QueueHead{head: s, queueType: AnotherSlotQueue}
	s.nextInQueue = s
	s.prevInQueue = s
}

func (s *Slot) _addQueuePrev(chainHead, chainTail *Slot) {
	s.ensureInQueue()

	prev := s.prevInQueue

	chainHead.prevInQueue = prev
	chainTail.nextInQueue = s

	s.prevInQueue = chainTail
	prev.nextInQueue = chainHead
}

func (s *Slot) QueueType() QueueType {
	if s.queue == nil {
		return NoQueue
	}
	return s.queue.queueType
}

func (s *Slot) QueueNext() *Slot {
	next := s.nextInQueue
	if next == nil || next.isQueueHead() {
		return nil
	}
	return next
}

func (s *Slot) QueuePrev() *Slot {
	prev := s.prevInQueue
	if prev == nil || prev.isQueueHead() {
		return nil
	}
	return prev
}

func (s *Slot) removeFromQueue() {
	if s.queue == nil {
		s.ensureNotInQueue()
		return
	}
	if s.isQueueHead() {
		panic("illegal state")
	}
	s.ensureInQueue()

	next := s.nextInQueue
	prev := s.prevInQueue

	next.prevInQueue = prev
	prev.nextInQueue = next

	s.queue.count--
	s.queue = nil
	s.nextInQueue = nil
	s.prevInQueue = nil
}
