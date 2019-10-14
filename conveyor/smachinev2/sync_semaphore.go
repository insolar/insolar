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

func NewSemaphore(initialCount int, name string) SemaphoreLink {
	ctl := &semaSync{}
	ctl.controller.Init(name)
	deps, _ := ctl.AdjustLimit(initialCount)
	if len(deps) != 0 {
		panic("illegal state")
	}
	return SemaphoreLink{ctl}
}

type SemaphoreLink struct {
	ctl *semaSync
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

func (v SemaphoreLink) NewBoolValue(isOpen bool) SyncAdjustment {
	if isOpen {
		return v.NewValue(0)
	}
	return v.NewValue(1)
}

func (v SemaphoreLink) SyncLink() SyncLink {
	return NewSyncLink(v.ctl)
}

type semaSync struct {
	controller holdingQueueController
}

func (p *semaSync) CheckState() Decision {
	if p.controller.IsOpen() {
		return Passed
	}
	return NotPassed
}

func (p *semaSync) CheckDependency(dep SlotDependency) Decision {
	if entry, ok := dep.(*DependencyQueueEntry); ok {
		switch {
		case !entry.link.IsValid(): // just to make sure
			return Impossible
		case !p.controller.Contains(entry):
			return Impossible
		case p.controller.IsOpen():
			return Passed
		default:
			return NotPassed
		}
	}
	return Impossible
}

func (p *semaSync) UseDependency(dep SlotDependency, oneStep bool) Decision {
	if entry, ok := dep.(*DependencyQueueEntry); ok {
		switch {
		case !entry.link.IsValid(): // just to make sure
			return Impossible
		case !oneStep && (entry.slotFlags&syncForOneStep != 0):
			return Impossible
		case !p.controller.Contains(entry):
			return Impossible
		case p.controller.IsOpen():
			return Passed
		default:
			return NotPassed
		}
	}
	return Impossible
}

func (p *semaSync) CreateDependency(slot *Slot, oneStep bool) (Decision, SlotDependency) {
	flags := DependencyQueueEntryFlags(0)
	if oneStep {
		flags |= syncForOneStep
	}
	if p.controller.IsOpen() {
		return Passed, nil
	}
	return NotPassed, p.controller.queue.AddSlot(slot.NewLink(), flags)
}

func (p *semaSync) GetLimit() (limit int, isAdjustable bool) {
	return p.controller.state, true
}

func (p *semaSync) AdjustLimit(limit int) ([]SlotLink, bool) {
	p.controller.state = limit
	if !p.controller.IsOpen() {
		return nil, false
	}
	return p.controller.queue.FlushAllAsLinks(), true
}

func (p *semaSync) GetWaitingCount() int {
	return p.controller.queue.Count()
}

func (p *semaSync) GetName() string {
	return p.controller.GetName()
}

type holdingQueueController struct {
	waitingQueueController
	state int
}

func (p *holdingQueueController) IsOpen() bool {
	return p.state <= 0
}

func (p *holdingQueueController) Release(_ SlotLink, _ DependencyQueueEntryFlags, removeFn func(), activateFn func(SlotLink)) {
	removeFn()
	if p.IsOpen() && p.queue.Count() > 0 {
		panic("illegal state")
	}
}

func (p *holdingQueueController) IsReleaseOnWorking(SlotLink, DependencyQueueEntryFlags) bool {
	return p.IsOpen()
}

func (p *holdingQueueController) IsReleaseOnStepping(_ SlotLink, flags DependencyQueueEntryFlags) bool {
	return flags&syncForOneStep != 0 || p.IsOpen()
}
