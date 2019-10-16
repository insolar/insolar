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

func NewExclusive(name string) SyncLink {
	ctl := &exclusiveSync{}
	ctl.awaiters.Init(name)
	return NewSyncLink(ctl)
}

type exclusiveSync struct {
	awaiters exclusiveQueueController
}

func (p *exclusiveSync) CheckState() Decision {
	if p.awaiters.IsEmpty() {
		return Passed
	}
	return NotPassed
}

func (p *exclusiveSync) CheckDependency(dep SlotDependency) Decision {
	if entry, ok := dep.(*dependencyQueueEntry); ok {
		switch {
		case !entry.link.IsValid(): // just to make sure
			return Impossible
		case !p.awaiters.Contains(entry):
			return Impossible
		case p.awaiters.IsEmptyOrFirst(entry.link):
			return Passed
		default:
			return NotPassed
		}
	}
	return Impossible
}

func (p *exclusiveSync) UseDependency(dep SlotDependency, flags SlotDependencyFlags) Decision {
	if entry, ok := dep.(*dependencyQueueEntry); ok {
		switch {
		case !entry.link.IsValid(): // just to make sure
			return Impossible
		case !p.awaiters.Contains(entry):
			return Impossible
		case !entry.IsCompatibleWith(flags):
			return Impossible
		case p.awaiters.IsEmptyOrFirst(entry.link):
			return Passed
		default:
			return NotPassed
		}
	}
	return Impossible
}

func (p *exclusiveSync) CreateDependency(slot *Slot, flags SlotDependencyFlags, syncer rwlock.RWLocker) (BoolDecision, SlotDependency) {
	sd := p.awaiters.queue.AddSlot(slot.NewLink(), flags)
	if f, _ := p.awaiters.queue.FirstValid(); f == sd {
		return true, sd
	}
	return false, sd
}

func (p *exclusiveSync) GetCounts() (active, inactive int) {
	n := p.awaiters.queue.Count()
	if n <= 0 {
		return 0, n
	}
	return 1, n - 1
}

func (p *exclusiveSync) GetName() string {
	return p.awaiters.GetName()
}

func (p *exclusiveSync) GetLimit() (limit int, isAdjustable bool) {
	return 1, false
}

func (p *exclusiveSync) AdjustLimit(limit int) (deps []StepLink, activate bool) {
	if limit != 1 {
		panic("illegal value")
	}
	return nil, false
}

var _ DependencyQueueController = &exclusiveQueueController{}

type exclusiveQueueController struct {
	name  string
	queue DependencyQueueHead
}

func (p *exclusiveQueueController) Init(name string) {
	if p.queue.controller != nil {
		panic("illegal state")
	}
	p.name = name
	p.queue.controller = p
}

func (p *exclusiveQueueController) IsEmpty() bool {
	return p.queue.IsEmpty()
}

func (p *exclusiveQueueController) IsEmptyOrFirst(link SlotLink) bool {
	f := p.queue.First()
	return f == nil || f.link == link
}

func (p *exclusiveQueueController) GetName() string {
	return p.name
}

func (p *exclusiveQueueController) IsReleaseOnWorking(SlotLink, SlotDependencyFlags) bool {
	return false
}

func (p *exclusiveQueueController) IsReleaseOnStepping(_ SlotLink, flags SlotDependencyFlags) bool {
	return flags&syncForOneStep != 0
}

func (p *exclusiveQueueController) Release(link SlotLink, flags SlotDependencyFlags, removeFn func()) []StepLink {
	if f, _ := p.queue.FirstValid(); f == nil || f.link != link {
		removeFn()
		return nil
	}

	removeFn()
	if f, step := p.queue.FirstValid(); f != nil {
		return []StepLink{step}
	}
	return nil
}

func (p *exclusiveQueueController) Contains(entry *dependencyQueueEntry) bool {
	return entry.queue == &p.queue
}
