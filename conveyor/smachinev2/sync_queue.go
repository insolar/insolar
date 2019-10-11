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

type DependencyQueueController interface {
	GetName() string
	IsReleaseOnStepping(link SlotLink, flags DependencyQueueEntryFlags) bool
	IsReleaseOnWorking(link SlotLink, flags DependencyQueueEntryFlags) bool
	Release(link SlotLink, flags DependencyQueueEntryFlags, removeFn func(), activateFn func(SlotLink))
	Dispose(link SlotLink, flags DependencyQueueEntryFlags, removeFn func(), activateFn func(SlotLink))
}

type DependencyQueueHead struct {
	controller DependencyQueueController
	head       DependencyQueueEntry
	count      int
}

func (p *DependencyQueueHead) AddSlot(link SlotLink, flags DependencyQueueEntryFlags) *DependencyQueueEntry {
	if !link.IsValid() {
		panic("illegal value")
	}
	entry := &DependencyQueueEntry{link: link, slotFlags: flags}
	p.AddLast(entry)
	return entry
}

func (p *DependencyQueueHead) AddFirst(entry *DependencyQueueEntry) {
	p.initEmpty()
	entry.ensureNotInQueue()

	p.head.nextInQueue._addQueuePrev(entry, entry)
	entry.queue = p
	p.count++
}

func (p *DependencyQueueHead) AddLast(entry *DependencyQueueEntry) {
	p.initEmpty()
	entry.ensureNotInQueue()

	p.head._addQueuePrev(entry, entry)
	entry.queue = p
	p.count++
}

func (p *DependencyQueueHead) Count() int {
	return p.count
}

func (p *DependencyQueueHead) FirstValid() *DependencyQueueEntry {
	for {
		f := p.head.QueueNext()
		if f == nil || f.link.IsValid() {
			return f
		}
		f.removeFromQueue()
	}
}

func (p *DependencyQueueHead) First() *DependencyQueueEntry {
	return p.head.QueueNext()
}

func (p *DependencyQueueHead) Last() *DependencyQueueEntry {
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

func (p *DependencyQueueHead) FlushOut(limit int, cutHead bool) []*DependencyQueueEntry {

	if limit > p.count {
		limit = p.count
	}
	if limit <= 0 {
		return nil
	}

	deps := make([]*DependencyQueueEntry, 0, limit)
	for limit > 0 {
		var entry *DependencyQueueEntry
		if cutHead {
			entry = p.First()
		} else {
			entry = p.Last()
		}
		if entry == nil {
			break
		}
		entry.removeFromQueue()

		if entry.link.IsValid() {
			deps = append(deps, entry)
			limit--
		}
	}
	return deps
}

type DependencyQueueEntryFlags uint32

var _ SlotDependency = &DependencyQueueEntry{}

type DependencyQueueEntry struct {
	queue                    *DependencyQueueHead
	nextInQueue, prevInQueue *DependencyQueueEntry
	slotFlags                DependencyQueueEntryFlags
	link                     SlotLink
}

func (p *DependencyQueueEntry) IsReleaseOnStepping() bool {
	if !p.isInQueue() {
		return true
	}
	return p.queue.controller.IsReleaseOnStepping(p.link, p.slotFlags)
}

func (p *DependencyQueueEntry) IsReleaseOnWorking() bool {
	if !p.isInQueue() {
		return true
	}
	return p.queue.controller.IsReleaseOnWorking(p.link, p.slotFlags)
}

func (p *DependencyQueueEntry) Release(activateFn func(SlotLink)) {
	if !p.isInQueue() {
		return
	}
	p.queue.controller.Release(p.link, p.slotFlags, p.removeFromQueue, activateFn)
}

func (p *DependencyQueueEntry) ReleaseOnDisposed(activateFn func(SlotLink)) {
	if !p.isInQueue() {
		return
	}
	p.queue.controller.Dispose(p.link, p.slotFlags, p.removeFromQueue, activateFn)
}

func (p *DependencyQueueEntry) _addQueuePrev(chainHead, chainTail *DependencyQueueEntry) {
	p.ensureInQueue()

	prev := p.prevInQueue

	chainHead.prevInQueue = prev
	chainTail.nextInQueue = p

	p.prevInQueue = chainTail
	prev.nextInQueue = chainHead
}

func (p *DependencyQueueEntry) QueueNext() *DependencyQueueEntry {
	next := p.nextInQueue
	if next == nil || next.isQueueHead() {
		return nil
	}
	return next
}

func (p *DependencyQueueEntry) QueuePrev() *DependencyQueueEntry {
	prev := p.prevInQueue
	if prev == nil || prev.isQueueHead() {
		return nil
	}
	return prev
}

func (p *DependencyQueueEntry) removeFromQueue() {
	if p.isQueueHead() {
		panic("illegal state")
	}
	p.ensureInQueue()

	next := p.nextInQueue
	prev := p.prevInQueue

	next.prevInQueue = prev
	prev.nextInQueue = next

	p.queue.count--
	p.queue = nil
	p.nextInQueue = nil
	p.prevInQueue = nil
}

func (p *DependencyQueueEntry) isQueueHead() bool {
	return p == &p.queue.head
}

func (p *DependencyQueueEntry) ensureNotInQueue() {
	if p.isInQueue() {
		panic("illegal state")
	}
}

func (p *DependencyQueueEntry) ensureInQueue() {
	if !p.isInQueue() {
		panic("illegal state")
	}
}

func (p *DependencyQueueEntry) isInQueue() bool {
	return p.queue != nil || p.nextInQueue != nil || p.prevInQueue != nil
}
