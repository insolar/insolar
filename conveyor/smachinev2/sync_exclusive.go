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
	if entry, ok := dep.(*DependencyQueueEntry); ok {
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

func (p *exclusiveSync) UseDependency(dep SlotDependency, oneStep bool) Decision {
	if entry, ok := dep.(*DependencyQueueEntry); ok {
		switch {
		case !entry.link.IsValid(): // just to make sure
			return Impossible
		case !p.awaiters.Contains(entry):
			return Impossible
		case !oneStep && (entry.slotFlags&syncForOneStep != 0):
			return Impossible
		case p.awaiters.IsEmptyOrFirst(entry.link):
			return Passed
		default:
			return NotPassed
		}
	}
	return Impossible
}

func (p *exclusiveSync) CreateDependency(slot *Slot, oneStep bool) (Decision, SlotDependency) {
	flags := DependencyQueueEntryFlags(0)
	if oneStep {
		flags |= syncForOneStep
	}
	sd := p.awaiters.queue.AddSlot(slot.NewLink(), flags)
	if p.awaiters.queue.FirstValid() == sd {
		return Passed, sd
	}
	return NotPassed, sd
}

func (p *exclusiveSync) GetWaitingCount() int {
	return p.awaiters.queue.Count()
}

func (p *exclusiveSync) GetName() string {
	return p.awaiters.GetName()
}

func (p *exclusiveSync) GetLimit() (limit int, isAdjustable bool) {
	return 1, false
}

func (p *exclusiveSync) AdjustLimit(limit int) (deps []SlotLink, activate bool) {
	if limit != 1 {
		panic("illegal value")
	}
	return nil, false
}

var _ DependencyQueueController = &exclusiveQueueController{}

const (
	syncForOneStep DependencyQueueEntryFlags = 1 << iota
)

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

func (p *exclusiveQueueController) IsReleaseOnWorking(SlotLink, DependencyQueueEntryFlags) bool {
	return false
}

func (p *exclusiveQueueController) IsReleaseOnStepping(_ SlotLink, flags DependencyQueueEntryFlags) bool {
	return flags&syncForOneStep != 0
}

func (p *exclusiveQueueController) Release(link SlotLink, flags DependencyQueueEntryFlags, removeFn func(), activateFn func(SlotLink)) {
	f := p.queue.FirstValid()
	isFirst := f != nil && f.link == link
	removeFn()
	if !isFirst {
		return
	}
	f = p.queue.FirstValid()
	if f != nil {
		activateFn(f.link)
	}
}

func (p *exclusiveQueueController) Dispose(link SlotLink, flags DependencyQueueEntryFlags, removeFn func(), activateFn func(SlotLink)) {
	p.Release(link, flags, removeFn, activateFn)
}

func (p *exclusiveQueueController) Contains(entry *DependencyQueueEntry) bool {
	return entry.queue == &p.queue
}
