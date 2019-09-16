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

import "github.com/insolar/insolar/network/consensus/common/rwlock"

func NewSlotPool(locker rwlock.RWLocker, pageSize uint16) SlotPool {
	if locker == nil {
		panic("illegal value")
	}
	if pageSize < 1 {
		panic("illegal value")
	}
	return SlotPool{
		mutex:       locker,
		slots:       [][]Slot{make([]Slot, pageSize)},
		unusedSlots: NewSlotQueue(UnusedSlots),
	}
}

type SlotPool struct {
	mutex rwlock.RWLocker

	slots     [][]Slot
	slotPgPos uint16

	unusedSlots SlotQueue
}

func (p *SlotPool) Count() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	n := len(p.slots)
	if n == 0 {
		return 0
	}
	return (n-1)*int(len(p.slots[0])) + int(p.slotPgPos) - p.unusedSlots.Count()
}

func (p *SlotPool) Capacity() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	n := len(p.slots)
	if n == 0 {
		return 0
	}
	return len(p.slots) * len(p.slots[0])
}

func (p *SlotPool) IsEmpty() bool {
	return p.Count() == 0
}

func (p *SlotPool) IsZero() bool {
	return p.mutex == nil
}

/* creates or reuse a slot, and marks it as BUSY */
func (p *SlotPool) AllocateSlot(m *SlotMachine, id SlotID) (slot *Slot) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	switch {
	case !p.unusedSlots.IsEmpty():
		slot = p.unusedSlots.First()
		slot.removeFromQueue()
		if slot.machine != m {
			panic("illegal state")
		}
	case p.slots == nil:
		panic("illegal state")
	default:
		lenSlots := len(p.slots[0])
		if int(p.slotPgPos) == lenSlots {
			p.slots = append(p.slots, p.slots[0])
			p.slots[0] = make([]Slot, lenSlots)
			p.slotPgPos = 0
		}
		slot = &p.slots[0][p.slotPgPos]
		slot.machine = m
		p.slotPgPos++
	}
	slot._slotAllocated(id)

	return slot
}

func (p *SlotPool) RecycleSlot(slot *Slot) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.unusedSlots.AddFirst(slot)
}

type SlotPageScanFunc func([]Slot, SlotWorker) (isPageEmptyOrWeak, hasWeakSlots bool)
type SlotDisposeFunc func(*Slot, SlotWorker)

func (p *SlotPool) ScanAndCleanup(cleanupWeak bool, w SlotWorker,
	disposeFn SlotDisposeFunc,
	scanPageFn SlotPageScanFunc,
) {
	if len(p.slots) == 0 || len(p.slots) == 1 && p.slotPgPos == 0 {
		return
	}

	isAllEmptyOrWeak, hasSomeWeakSlots := scanPageFn(p.slots[0][:p.slotPgPos], w)

	j := 1
	for i, slotPage := range p.slots[1:] {
		isPageEmptyOrWeak, hasWeakSlots := scanPageFn(slotPage, w)
		switch {
		case !isPageEmptyOrWeak:
			isAllEmptyOrWeak = false
		case !hasWeakSlots:
			cleanupEmptyPage(slotPage)
			p.slots[i+1] = nil
			continue
		default:
			hasSomeWeakSlots = true
		}

		if j != i+1 {
			p.slots[j] = slotPage
			p.slots[i+1] = nil
		}
		j++
	}

	if isAllEmptyOrWeak && (cleanupWeak || !hasSomeWeakSlots) {
		for _, slotPage := range p.slots {
			for i := range slotPage {
				slot := &slotPage[i]
				if slot.isEmpty() {
					continue
				}
				disposeFn(slot, w)
			}
		}
		p.slots = p.slots[:1]
		p.slotPgPos = 0
		return
	}

	if len(p.slots) > j {
		p.slots = p.slots[:j]
	}
}

func cleanupEmptyPage(slotPage []Slot) {
	for i := range slotPage {
		slot := &slotPage[i]
		if slot.QueueType() != UnusedSlots {
			panic("illegal state")
		}
		slot.removeFromQueue()
	}
}
