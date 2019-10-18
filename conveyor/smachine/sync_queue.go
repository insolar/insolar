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

import "sync/atomic"

type DependencyQueueController interface {
	GetName() string
	IsReleaseOnStepping(link SlotLink, flags SlotDependencyFlags) bool
	IsReleaseOnWorking(link SlotLink, flags SlotDependencyFlags) bool
	Release(link SlotLink, flags SlotDependencyFlags, removeFn func()) []StepLink
}

type DependencyQueueHead struct {
	controller DependencyQueueController
	head       dependencyQueueEntry
	count      int
}

func (p *DependencyQueueHead) AddSlot(link SlotLink, flags SlotDependencyFlags) *dependencyQueueEntry {
	if !link.IsValid() {
		panic("illegal value")
	}
	entry := &dependencyQueueEntry{link: link, slotFlags: uint32(flags << flagsOffset)}
	p.AddLast(entry)
	return entry
}

func (p *DependencyQueueHead) AddFirst(entry *dependencyQueueEntry) {
	p.initEmpty()
	entry.ensureNotInQueue()

	p.head.nextInQueue._addQueuePrev(entry, entry)
	entry.setQueue(p)
	p.count++
}

func (p *DependencyQueueHead) AddLast(entry *dependencyQueueEntry) {
	p.initEmpty()
	entry.ensureNotInQueue()

	p.head._addQueuePrev(entry, entry)
	entry.setQueue(p)
	p.count++
}

func (p *DependencyQueueHead) Count() int {
	return p.count
}

func (p *DependencyQueueHead) FirstValid() (*dependencyQueueEntry, StepLink) {
	for {
		f := p.head.QueueNext()
		if f == nil {
			return f, StepLink{}
		}
		if step, ok := f.link.GetStepLink(); ok {
			return f, step
		}
		f.removeFromQueue()
	}
}

func (p *DependencyQueueHead) First() *dependencyQueueEntry {
	return p.head.QueueNext()
}

func (p *DependencyQueueHead) Last() *dependencyQueueEntry {
	return p.head.QueuePrev()
}

func (p *DependencyQueueHead) IsZero() bool {
	return p.head.nextInQueue == nil
}

func (p *DependencyQueueHead) IsEmpty() bool {
	return p.head.nextInQueue == nil || p.head.nextInQueue.isQueueHead()
}

func (p *DependencyQueueHead) initEmpty() {
	if p.head.queue == nil {
		p.head.nextInQueue = &p.head
		p.head.prevInQueue = &p.head
		p.head.queue = p
	}
}

func (p *DependencyQueueHead) CutHeadOut(fn func(*dependencyQueueEntry) bool) {
	for {
		entry := p.First()
		if entry == nil {
			return
		}
		entry.removeFromQueue()

		if !fn(entry) {
			return
		}
	}
}

func (p *DependencyQueueHead) CutTailOut(fn func(*dependencyQueueEntry) bool) {
	for {
		entry := p.Last()
		if entry == nil {
			return
		}
		entry.removeFromQueue()

		if !fn(entry) {
			return
		}
	}
}

func (p *DependencyQueueHead) FlushAllAsLinks() []StepLink {
	if p.count == 0 {
		return nil
	}

	deps := make([]StepLink, 0, p.count)
	for {
		entry := p.First()
		if entry == nil {
			break
		}
		entry.removeFromQueue()

		if step, ok := entry.link.GetStepLink(); ok {
			deps = append(deps, step)
		}
	}
	return deps
}

const flagsOffset = 1
const atomicInQueue = 1 << (flagsOffset - 1)

var _ SlotDependency = &dependencyQueueEntry{}

type dependencyQueueEntry struct {
	queue                    *DependencyQueueHead
	nextInQueue, prevInQueue *dependencyQueueEntry
	slotFlags                uint32
	link                     SlotLink
}

func (p *dependencyQueueEntry) getFlags() (bool, SlotDependencyFlags) {
	v := atomic.LoadUint32(&p.slotFlags)
	return v&atomicInQueue != 0, SlotDependencyFlags(v >> 1)
}

func (p *dependencyQueueEntry) IsReleaseOnStepping() bool {
	if inQueue, flags := p.getFlags(); inQueue {
		return p.queue.controller.IsReleaseOnStepping(p.link, flags)
	}
	return true
}

func (p *dependencyQueueEntry) IsReleaseOnWorking() bool {
	if inQueue, flags := p.getFlags(); inQueue {
		return p.queue.controller.IsReleaseOnWorking(p.link, flags)
	}
	return true
}

func (p *dependencyQueueEntry) Release() []StepLink {
	if inQueue, flags := p.getFlags(); inQueue {
		return p.queue.controller.Release(p.link, flags, p.removeFromQueue)
	}
	return nil
}

func (p *dependencyQueueEntry) _addQueuePrev(chainHead, chainTail *dependencyQueueEntry) {
	p.ensureInQueue()

	prev := p.prevInQueue

	chainHead.prevInQueue = prev
	chainTail.nextInQueue = p

	p.prevInQueue = chainTail
	prev.nextInQueue = chainHead
}

func (p *dependencyQueueEntry) QueueNext() *dependencyQueueEntry {
	next := p.nextInQueue
	if next == nil || next.isQueueHead() {
		return nil
	}
	return next
}

func (p *dependencyQueueEntry) QueuePrev() *dependencyQueueEntry {
	prev := p.prevInQueue
	if prev == nil || prev.isQueueHead() {
		return nil
	}
	return prev
}

func (p *dependencyQueueEntry) removeFromQueue() {
	if p.isQueueHead() {
		panic("illegal state")
	}
	p.ensureInQueue()

	next := p.nextInQueue
	prev := p.prevInQueue

	next.prevInQueue = prev
	prev.nextInQueue = next

	p.queue.count--
	p.setQueue(nil)
	p.nextInQueue = nil
	p.prevInQueue = nil
}

func (p *dependencyQueueEntry) isQueueHead() bool {
	return p == &p.queue.head
}

func (p *dependencyQueueEntry) ensureNotInQueue() {
	if p.isInQueue() {
		panic("illegal state")
	}
}

func (p *dependencyQueueEntry) ensureInQueue() {
	if !p.isInQueue() {
		panic("illegal state")
	}
}

func (p *dependencyQueueEntry) isInQueue() bool {
	return p.queue != nil || p.nextInQueue != nil || p.prevInQueue != nil
}

func (p *dependencyQueueEntry) setQueue(head *DependencyQueueHead) {
	p.queue = head
	for {
		v := atomic.LoadUint32(&p.slotFlags)
		vv := v
		if head == nil {
			vv &^= atomicInQueue
		} else {
			vv |= atomicInQueue
		}
		if v == vv || atomic.CompareAndSwapUint32(&p.slotFlags, v, vv) {
			return
		}
	}
}

func (p *dependencyQueueEntry) IsCompatibleWith(flags SlotDependencyFlags) bool {
	_, f := p.getFlags()
	return f&flags == flags
}
