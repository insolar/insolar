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

import (
	"sync"
)

// SlotPool by default recycles deallocated pages to mitigate possible memory leak through SlotLink references
// When flow of slots varies a lot and there is no long-living links then deallocateOnCleanup can be enabled.
func newSlotPool(pageSize uint16, deallocateOnCleanup bool) SlotPool {
	//if locker == nil {
	//	panic("illegal value")
	//}
	if pageSize < 1 {
		panic("illegal value")
	}
	return SlotPool{
		slotPages:  [][]Slot{make([]Slot, pageSize)},
		deallocate: deallocateOnCleanup,
	}
}

type SlotPool struct {
	mutex sync.RWMutex

	unusedSlots SlotQueue
	slotPages   [][]Slot
	emptyPages  [][]Slot
	slotPgPos   uint16
	deallocate  bool
}

func (p *SlotPool) initSlotPool() {
	if p.slotPages == nil {
		panic("illegal nil")
	}
	p.unusedSlots.initSlotQueue(UnusedSlots)
}

func (p *SlotPool) Count() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	n := len(p.slotPages)
	if n == 0 {
		return 0
	}
	return (n-1)*int(len(p.slotPages[0])) + int(p.slotPgPos) - p.unusedSlots.Count()
}

func (p *SlotPool) Capacity() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	n := len(p.slotPages)
	if n == 0 {
		return 0
	}
	return len(p.slotPages) * len(p.slotPages[0])
}

func (p *SlotPool) IsEmpty() bool {
	return p.Count() == 0
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
	case p.slotPages == nil:
		panic("illegal state")
	default:
		lenSlots := len(p.slotPages[0])
		if int(p.slotPgPos) == lenSlots {
			p.slotPages = append(p.slotPages, p.slotPages[0])
			p.slotPages[0] = p.allocatePage(lenSlots)
			p.slotPgPos = 0
		}
		slot = &p.slotPages[0][p.slotPgPos]
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

type SlotPageScanFunc func([]Slot, FixedSlotWorker) (isPageEmptyOrWeak, hasWeakSlots bool)
type SlotDisposeFunc func(*Slot, FixedSlotWorker)

func (p *SlotPool) ScanAndCleanup(cleanupWeak bool, w FixedSlotWorker,
	disposeWeakFn SlotDisposeFunc, scanPageFn SlotPageScanFunc,
) bool {
	if len(p.slotPages) == 0 || len(p.slotPages) == 1 && p.slotPgPos == 0 {
		return true
	}

	isAllEmptyOrWeak, hasSomeWeakSlots := scanPageFn(p.slotPages[0][:p.slotPgPos], w)

	nextSlotPageNo := 1
	for i, slotPage := range p.slotPages[1:] {
		isPageEmptyOrWeak, hasWeakSlots := scanPageFn(slotPage, w)
		switch {
		case !isPageEmptyOrWeak:
			isAllEmptyOrWeak = false
		case !hasWeakSlots:
			cleanupEmptyPage(slotPage)
			p.recyclePage(slotPage)
			p.slotPages[i+1] = nil
			continue
		default:
			hasSomeWeakSlots = true
		}

		if nextSlotPageNo != i+1 {
			p.slotPages[nextSlotPageNo] = slotPage
			p.slotPages[i+1] = nil
		}
		nextSlotPageNo++
	}

	if isAllEmptyOrWeak && (cleanupWeak || !hasSomeWeakSlots) {
		for i, slotPage := range p.slotPages {
			if slotPage == nil {
				break
			}
			for j := range slotPage {
				slot := &slotPage[j]
				if !slot.isEmpty() {
					disposeWeakFn(slot, w)
				}
				qt := slot.QueueType()
				if qt == UnusedSlots {
					slot.removeFromQueue()
					continue
				}
				if qt == NoQueue && i == 0 {
					break
				}
				panic("illegal state")
			}
		}
		if p.unusedSlots.Count() != 0 {
			panic("illegal state")
		}
		p.slotPages = p.slotPages[:1]
		p.slotPgPos = 0
		return true
	}

	if len(p.slotPages) > nextSlotPageNo {
		p.slotPages = p.slotPages[:nextSlotPageNo]
	}
	return false
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

func (p *SlotPool) allocatePage(lenSlots int) []Slot {
	n := len(p.emptyPages)
	if n > 0 {
		n--
		pg := p.emptyPages[n]
		p.emptyPages[n] = nil
		p.emptyPages = p.emptyPages[:n]

		if len(pg) == lenSlots {
			return pg
		}
	}
	return make([]Slot, lenSlots)
}

func (p *SlotPool) recyclePage(pg []Slot) {
	if !p.deallocate {
		p.emptyPages = append(p.emptyPages, pg)
	}
}
