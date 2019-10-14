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

func NewFixedLimiter(initialLimit int, name string) SyncLink {
	if initialLimit < 0 {
		panic("illegal value")
	}
	switch initialLimit {
	case 0:
		return NewInfiniteLock(name)
	case 1:
		return NewExclusive(name)
	default:
		return NewSyncLink(newLimiter(initialLimit, false, name))
	}
}

func NewLimiter(initialLimit int, name string) LimiterLink {
	return LimiterLink{newLimiter(initialLimit, true, name)}
}

type LimiterLink struct {
	ctl *limiterSync
}

func (v LimiterLink) IsZero() bool {
	return v.ctl == nil
}

func (v LimiterLink) NewDelta(delta int) SyncAdjustment {
	if v.ctl == nil {
		panic("illegal state")
	}
	return SyncAdjustment{controller: v.ctl, adjustment: delta, isAbsolute: false}
}

func (v LimiterLink) NewValue(value int) SyncAdjustment {
	if v.ctl == nil {
		panic("illegal state")
	}
	return SyncAdjustment{controller: v.ctl, adjustment: value, isAbsolute: true}
}

func (v LimiterLink) SyncLink() SyncLink {
	return NewSyncLink(v.ctl)
}

func newLimiter(initialLimit int, isAdjustable bool, name string) *limiterSync {
	ctl := &limiterSync{isAdjustable: true}
	ctl.controller.Init(name)
	deps, _ := ctl.AdjustLimit(initialLimit)
	if len(deps) != 0 {
		panic("illegal state")
	}
	ctl.isAdjustable = isAdjustable
	return ctl
}

type limiterSync struct {
	controller   workingQueueController
	isAdjustable bool
}

func (p *limiterSync) CheckState() Decision {
	if p.controller.IsOpen() {
		return Passed
	}
	return NotPassed
}

func (p *limiterSync) CheckDependency(dep SlotDependency) Decision {
	if entry, ok := dep.(*DependencyQueueEntry); ok {
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

func (p *limiterSync) UseDependency(dep SlotDependency, oneStep bool) Decision {
	if entry, ok := dep.(*DependencyQueueEntry); ok {
		switch {
		case !entry.link.IsValid(): // just to make sure
			return Impossible
		case !oneStep && (entry.slotFlags&syncForOneStep != 0):
			return Impossible
		case p.controller.Contains(entry):
			return Passed
		case p.controller.ContainsInAwaiters(entry):
			return NotPassed
		}
	}
	return Impossible
}

func (p *limiterSync) CreateDependency(slot *Slot, oneStep bool) (Decision, SlotDependency) {
	flags := DependencyQueueEntryFlags(0)
	if oneStep {
		flags |= syncForOneStep
	}
	if p.controller.IsOpen() {
		return Passed, p.controller.queue.AddSlot(slot.NewLink(), flags)
	}
	return NotPassed, p.controller.awaiters.queue.AddSlot(slot.NewLink(), flags)
}

func (p *limiterSync) GetLimit() (limit int, isAdjustable bool) {
	return p.controller.workerLimit, p.isAdjustable
}

func (p *limiterSync) AdjustLimit(limit int) ([]SlotLink, bool) {
	if p.controller.workerLimit == limit {
		return nil, false
	}
	if !p.isAdjustable {
		panic("illegal value")
	}

	delta := limit - p.controller.workerLimit
	p.controller.workerLimit = limit

	if delta > 0 {
		toBeActivated := p.controller.awaiters.queue.FlushOut(delta, true)

		links := make([]SlotLink, len(toBeActivated))
		for i, entry := range toBeActivated {
			links[i] = entry.link
			p.controller.queue.AddLast(entry)
		}
		return links, true
	}

	// sequence is reversed!
	toBeDeactivated := p.controller.queue.FlushOut(-delta, false)

	links := make([]SlotLink, len(toBeDeactivated))
	for i, entry := range toBeDeactivated {
		// keep the original sequence by reversing
		links[len(toBeDeactivated)-i] = entry.link
		p.controller.queue.AddFirst(entry)
	}
	return links, false
}

func (p *limiterSync) GetWaitingCount() int {
	return p.controller.queue.Count()
}

func (p *limiterSync) GetName() string {
	return p.controller.GetName()
}

type waitingQueueController struct {
	exclusiveQueueController
}

func (p *waitingQueueController) Release(_ SlotLink, _ DependencyQueueEntryFlags, removeFn func(), _ func(SlotLink)) {
	removeFn()
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

func (p *workingQueueController) IsOpen() bool {
	return p.queue.Count() < p.workerLimit
}

func (p *workingQueueController) Release(link SlotLink, flags DependencyQueueEntryFlags, removeFn func(), activateFn func(SlotLink)) {
	removeFn()

	for {
		// early check to provide some cleanup
		f := p.awaiters.queue.FirstValid()
		if f == nil {
			return
		}
		if !p.IsOpen() {
			return
		}
		f.removeFromQueue()
		p.queue.AddLast(f)
		activateFn(f.link)
	}
}

func (p *workingQueueController) ContainsInAwaiters(entry *DependencyQueueEntry) bool {
	return p.awaiters.Contains(entry)
}
