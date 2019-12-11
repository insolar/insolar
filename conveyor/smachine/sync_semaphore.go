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
	"sync"
)

// Semaphore allows Acquire() call to pass through for a number of workers within the limit.
func NewFixedSemaphore(limit int, name string) SyncLink {
	return NewFixedSemaphoreWithFlags(limit, name, 0)
}

func NewFixedSemaphoreWithFlags(limit int, name string, flags DependencyQueueFlags) SyncLink {
	if limit < 0 {
		panic("illegal value")
	}
	switch limit {
	case 0:
		return NewInfiniteLock(name)
	case 1:
		return NewExclusiveWithFlags(name, flags)
	default:
		return NewSyncLink(newSemaphore(limit, false, name, flags))
	}
}

// Semaphore allows Acquire() call to pass through for a number of workers within the limit.
// Negative and zero values are not passable.
// The limit can be changed with adjustments. Overflows are capped by min/max int.
func NewSemaphore(initialValue int, name string) SemaphoreLink {
	return NewSemaphoreWithFlags(initialValue, name, 0)
}

func NewSemaphoreWithFlags(initialValue int, name string, flags DependencyQueueFlags) SemaphoreLink {
	return SemaphoreLink{newSemaphore(initialValue, true, name, flags)}
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

func (v SemaphoreLink) NewChild(childValue int, name string) SyncLink {
	if childValue <= 0 {
		return NewInfiniteLock(name)
	}
	return NewSyncLink(newSemaphoreChild(v.ctl, childValue, name))
}

func (v SemaphoreLink) SyncLink() SyncLink {
	return NewSyncLink(v.ctl)
}

func newSemaphore(initialLimit int, isAdjustable bool, name string, flags DependencyQueueFlags) *semaphoreSync {
	ctl := &semaphoreSync{isAdjustable: true}
	ctl.controller.awaiters.queue.flags = flags
	ctl.controller.Init(name, &ctl.mutex, &ctl.controller)

	deps, _ := ctl.AdjustLimit(initialLimit, false)
	if len(deps) != 0 {
		panic("illegal state")
	}
	ctl.isAdjustable = isAdjustable
	return ctl
}

type semaphoreSync struct {
	mutex        sync.RWMutex
	controller   workingQueueController
	isAdjustable bool
}

func (p *semaphoreSync) CheckState() BoolDecision {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.checkState()
}

func (p *semaphoreSync) checkState() BoolDecision {
	return BoolDecision(p.controller.canPassThrough())
}

func (p *semaphoreSync) CheckDependency(dep SlotDependency) Decision {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if entry, ok := dep.(*dependencyQueueEntry); ok {
		return p.checkDependency(entry)
	}
	return Impossible
}

func (p *semaphoreSync) checkDependency(entry *dependencyQueueEntry) Decision {
	switch {
	case !entry.link.IsValid(): // just to make sure
		return Impossible
	case p.controller.contains(entry):
		return Passed
	case p.controller.containsInAwaiters(entry):
		return NotPassed
	}
	return Impossible
}

func (p *semaphoreSync) UseDependency(dep SlotDependency, flags SlotDependencyFlags) Decision {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if entry, ok := dep.(*dependencyQueueEntry); ok {
		if d := p.checkDependency(entry); d.IsValid() && entry.IsCompatibleWith(flags) {
			return d
		}
	}
	return Impossible
}

func (p *semaphoreSync) createDependency(holder SlotLink, flags SlotDependencyFlags) (BoolDecision, *dependencyQueueEntry) {
	if p.controller.canPassThrough() {
		return true, p.controller.queue.AddSlot(holder, flags)
	}
	return false, p.controller.awaiters.queue.AddSlot(holder, flags)
}

func (p *semaphoreSync) CreateDependency(holder SlotLink, flags SlotDependencyFlags) (BoolDecision, SlotDependency) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.createDependency(holder, flags)
}

func (p *semaphoreSync) GetLimit() (limit int, isAdjustable bool) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.controller.workerLimit, p.isAdjustable
}

func (p *semaphoreSync) AdjustLimit(limit int, absolute bool) ([]StepLink, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

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
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.controller.queue.Count(), p.controller.awaiters.queue.Count()
}

func (p *semaphoreSync) GetName() string {
	return p.controller.GetName()
}

func (p *semaphoreSync) EnumQueues(fn EnumQueueFunc) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.controller.enum(1, fn)
}

type waitingQueueController struct {
	mutex *sync.RWMutex
	queueControllerTemplate
}

func (p *waitingQueueController) Init(name string, mutex *sync.RWMutex, controller DependencyQueueController) {
	p.queueControllerTemplate.Init(name, mutex, controller)
	p.mutex = mutex
}

func (p *waitingQueueController) IsOpen(SlotDependency) bool {
	return false
}

func (p *waitingQueueController) Release(link SlotLink, flags SlotDependencyFlags, removeFn func()) ([]PostponedDependency, []StepLink) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	removeFn()
	return nil, nil
}

func (p *waitingQueueController) moveToInactive(n int, q *DependencyQueueHead, stacker *dependencyStackEntry) int {
	if n <= 0 {
		return 0
	}

	count := 0
	for n > 0 {
		if f, _ := p.queue.FirstValid(); f == nil {
			break
		} else {
			f.removeFromQueue()
			switch {
			case stacker == nil:
			case f.stacker != nil:
				panic("illegal state")
			default:
				f.stacker = stacker
			}
			q.AddLast(f)
			count++
			n--
		}
	}
	return count
}

func (p *waitingQueueController) moveToActive(n int, q *DependencyQueueHead) ([]PostponedDependency, []StepLink) {
	if n <= 0 {
		return nil, nil
	}

	var postponed []PostponedDependency
	links := make([]StepLink, 0, n)
	for n > 0 {
		if f, step := p.queue.FirstValid(); f == nil {
			break
		} else {
			f.removeFromQueue()
			q.AddLast(f)
			if pp := f.stacker.ActivateStack(f, step); pp != nil {
				postponed = append(postponed, pp)
			} else {
				links = append(links, step)
			}
			n--
		}
	}
	return postponed, links
}

type workingQueueController struct {
	queueControllerTemplate
	workerLimit int
	awaiters    waitingQueueController
}

func (p *workingQueueController) Init(name string, mutex *sync.RWMutex, controller DependencyQueueController) {
	p.queueControllerTemplate.Init(name, mutex, controller)
	p.awaiters.Init(name, mutex, &p.awaiters)
}

func (p *workingQueueController) IsOpen(SlotDependency) bool {
	return true
}

func (p *workingQueueController) canPassThrough() bool {
	return p.queue.Count() < p.workerLimit
}

func (p *workingQueueController) Release(link SlotLink, flags SlotDependencyFlags, removeFn func()) ([]PostponedDependency, []StepLink) {
	p.awaiters.mutex.Lock()
	defer p.awaiters.mutex.Unlock()

	removeFn()
	// p.queue.FirstValid() // check for stale items
	return p.awaiters.moveToActive(p.workerLimit-p.queue.Count(), &p.queue)
}

func (p *workingQueueController) containsInAwaiters(entry *dependencyQueueEntry) bool {
	return p.awaiters.contains(entry)
}

func (p *workingQueueController) enum(qId int, fn EnumQueueFunc) bool {
	if p.queueControllerTemplate.enum(qId, fn) {
		return true
	}
	return p.awaiters.enum(qId-1, fn)
}
