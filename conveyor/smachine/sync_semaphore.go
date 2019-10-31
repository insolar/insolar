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

import (
	"math"
)

// Semaphore allows Acquire() call to pass through for a number of workers within the limit.
func NewFixedSemaphore(limit int, name string) SyncLink {
	if limit < 0 {
		panic("illegal value")
	}
	switch limit {
	case 0:
		return NewInfiniteLock(name)
	case 1:
		return NewExclusive(name)
	default:
		return NewSyncLink(newSemaphore(limit, false, name))
	}
}

// Semaphore allows Acquire() call to pass through for a number of workers within the limit.
// Negative and zero values are not passable.
// The limit can be changed with adjustments. Overflows are capped by min/max int.
func NewSemaphore(initialValue int, name string) SemaphoreLink {
	return SemaphoreLink{newSemaphore(initialValue, true, name)}
}

type SemaphoreLink struct {
	ctl *semaphoreSync
}

func (v SemaphoreLink) IsZero() bool {
	return v.ctl == nil
}

func (v SemaphoreLink) NewDelta(delta int) SyncAdjustment {
	if v.ctl == nil {
		panic("illegal state")
	}
	return SyncAdjustment{controller: v.ctl, adjustment: delta, isAbsolute: false}
}

func (v SemaphoreLink) NewValue(value int) SyncAdjustment {
	if v.ctl == nil {
		panic("illegal state")
	}
	return SyncAdjustment{controller: v.ctl, adjustment: value, isAbsolute: true}
}

func (v SemaphoreLink) SyncLink() SyncLink {
	return NewSyncLink(v.ctl)
}

func newSemaphore(initialLimit int, isAdjustable bool, name string) *semaphoreSync {
	ctl := &semaphoreSync{isAdjustable: true}
	ctl.controller.Init(name)
	deps, _ := ctl.AdjustLimit(initialLimit, false)
	if len(deps) != 0 {
		panic("illegal state")
	}
	ctl.isAdjustable = isAdjustable
	return ctl
}

type semaphoreSync struct {
	controller   workingQueueController
	isAdjustable bool
}

func (p *semaphoreSync) CheckState() Decision {
	if p.controller.isOpen() {
		return Passed
	}
	return NotPassed
}

func (p *semaphoreSync) CheckDependency(dep SlotDependency) Decision {
	if entry, ok := dep.(*dependencyQueueEntry); ok {
		switch {
		case !entry.link.IsValid(): // just to make sure
			return Impossible
		case p.controller.Contains(entry):
			return Passed
		case p.controller.ContainsInAwaiters(entry):
			return NotPassed
		}
	}
	return Impossible
}

func (p *semaphoreSync) UseDependency(dep SlotDependency, flags SlotDependencyFlags) (Decision, SlotDependency) {
	if entry, ok := dep.(*dependencyQueueEntry); ok {
		switch {
		case !entry.link.IsValid(): // just to make sure
			return Impossible, nil
		case !entry.IsCompatibleWith(flags):
			return Impossible, nil
		case p.controller.Contains(entry):
			return Passed, nil
		case p.controller.ContainsInAwaiters(entry):
			return NotPassed, nil
		}
	}
	return Impossible, nil
}

func (p *semaphoreSync) CreateDependency(holder SlotLink, flags SlotDependencyFlags) (BoolDecision, SlotDependency) {
	if p.controller.isOpen() {
		return true, p.controller.queue.AddSlot(holder, flags)
	}
	return false, p.controller.awaiters.queue.AddSlot(holder, flags)
}

func (p *semaphoreSync) GetLimit() (limit int, isAdjustable bool) {
	return p.controller.workerLimit, p.isAdjustable
}

func (p *semaphoreSync) AdjustLimit(limit int, absolute bool) ([]StepLink, bool) {
	if !p.isAdjustable {
		panic("illegal state")
	}

	if ok, newLimit := applyWrappedAdjustment(p.controller.workerLimit, limit, math.MinInt32, math.MaxInt32, absolute); ok {
		limit = newLimit
	} else {
		return nil, false
	}

	delta := limit - p.controller.workerLimit
	p.controller.workerLimit = limit

	if delta > 0 {
		links := make([]StepLink, delta)
		pos := 0
		p.controller.awaiters.queue.CutHeadOut(func(entry *dependencyQueueEntry) bool {
			if step, ok := entry.link.GetStepLink(); ok {
				p.controller.queue.AddLast(entry)
				links[pos] = step
				pos++
				return pos < delta
			}
			return true
		})
		return links[:pos], true
	}

	delta = -delta
	links := make([]StepLink, delta)

	// sequence is reversed!
	p.controller.queue.CutTailOut(func(entry *dependencyQueueEntry) bool {
		if step, ok := entry.link.GetStepLink(); ok {
			p.controller.awaiters.queue.AddFirst(entry)
			delta--
			links[delta] = step
			return delta > 0
		}
		return true
	})
	return links[delta:], false
}

func (p *semaphoreSync) GetCounts() (active, inactive int) {
	return p.controller.queue.Count(), p.controller.awaiters.queue.Count()
}

func (p *semaphoreSync) GetName() string {
	return p.controller.GetName()
}

type waitingQueueController struct {
	exclusiveQueueController
}

func (p *waitingQueueController) IsOpen(SlotDependency) bool {
	return false
}

func (p *waitingQueueController) Release(link SlotLink, flags SlotDependencyFlags, removeFn func()) ([]PostponedDependency, []StepLink) {
	removeFn()
	return nil, nil
}

type workingQueueController struct {
	exclusiveQueueController
	workerLimit int
	awaiters    waitingQueueController
}

func (p *workingQueueController) Init(name string) {
	if p.queue.controller != nil {
		panic("illegal state")
	}
	p.name = name
	p.awaiters.name = name
	p.queue.controller = p
	p.awaiters.queue.controller = &p.awaiters
}

func (p *workingQueueController) IsOpen(SlotDependency) bool {
	return p.isOpen()
}

func (p *workingQueueController) isOpen() bool {
	return p.queue.Count() < p.workerLimit
}

func (p *workingQueueController) Release(link SlotLink, flags SlotDependencyFlags, removeFn func()) ([]PostponedDependency, []StepLink) {
	removeFn()

	n := p.workerLimit - p.queue.Count()
	if n <= 0 {
		return nil, nil
	}

	var postponed []PostponedDependency
	links := make([]StepLink, 0, n)
	for n > 0 {
		if f, step := p.awaiters.queue.FirstValid(); f == nil {
			break
		} else {
			f.removeFromQueue()
			p.queue.AddLast(f)
			if pp := f.childOf.ActivateStack(f, step); pp != nil {
				postponed = append(postponed, pp)
			} else {
				links = append(links, step)
			}
			n--
		}
	}
	return postponed, links
}

func (p *workingQueueController) ContainsInAwaiters(entry *dependencyQueueEntry) bool {
	return p.awaiters.Contains(entry)
}
