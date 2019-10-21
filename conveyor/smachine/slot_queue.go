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

type QueueType int8

const (
	InvalidQueue QueueType = iota

	ActivationOfSlot

	UnusedSlots

	PollingSlots

	ActiveSlots
	WorkingSlots
)

const NoQueue QueueType = -1

func (v QueueType) IsActiveOrPolling() bool {
	return v >= PollingSlots
}

func (v QueueType) IsActive() bool {
	return v >= ActiveSlots
}

type SlotQueue struct {
	QueueHead
	slot Slot
}

func (p *SlotQueue) initSlotQueue(queueType QueueType) {
	p.head = &p.slot
	p.initQueueHead(queueType)
}

func (p *SlotQueue) AppendAll(anotherQueue *SlotQueue) {
	p.QueueHead.AppendAll(&anotherQueue.QueueHead)
}

func (p *SlotQueue) PrependAll(anotherQueue *SlotQueue) {
	p.QueueHead.PrependAll(&anotherQueue.QueueHead)
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

func (p *QueueHead) initQueueHead(queueType QueueType) {
	if queueType <= InvalidQueue {
		panic("illegal value")
	}
	p.queueType = queueType
	p.head.nextInQueue = p.head
	p.head.prevInQueue = p.head
	p.head.queue = p
}

func (p *QueueHead) initEmpty() {
	//if p.head.queue == nil {
	//	p.initQueueHead()
	//	//p.head.nextInQueue = p.head
	//	//p.head.prevInQueue = p.head
	//	//p.head.queue = p
	//}
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

func (p *QueueHead) extractAll(targetQueue *QueueHead) (head, tail *Slot, count int) {
	c := p.count
	if c == 0 {
		return nil, nil, 0
	}

	next := p.head.nextInQueue
	prev := p.head.prevInQueue
	prev.nextInQueue = nil

	p.count = 0
	p.head.nextInQueue = p.head
	p.head.prevInQueue = p.head

	for n := next; n != nil; n = n.nextInQueue {
		n.queue = targetQueue
	}

	return next, prev, c
}

func (p *QueueHead) AppendAll(anotherQueue *QueueHead) {
	p.initEmpty()
	if anotherQueue.IsEmpty() {
		return
	}
	head, tail, count := anotherQueue.extractAll(p)
	p.head._addQueuePrev(head, tail)
	p.count += count
}

func (p *QueueHead) PrependAll(anotherQueue *QueueHead) {
	p.initEmpty()
	if anotherQueue.IsEmpty() {
		return
	}
	head, tail, count := anotherQueue.extractAll(p)
	p.head.nextInQueue._addQueuePrev(head, tail)
	p.count += count
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

	s.queue.initQueueHead(ActivationOfSlot)
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

func (s *Slot) _cutNext() *Slot {
	next := s.nextInQueue
	s.nextInQueue = nil
	s.prevInQueue = nil
	s.queue = nil
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
