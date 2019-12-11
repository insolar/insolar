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

import "sync"

func NewExclusive(name string) SyncLink {
	return NewExclusiveWithFlags(name, 0)
}

func NewExclusiveWithFlags(name string, flags DependencyQueueFlags) SyncLink {
	ctl := &exclusiveSync{}
	ctl.awaiters.queue.flags = flags
	ctl.awaiters.Init(name, &ctl.mutex, &ctl.awaiters)
	return NewSyncLink(ctl)
}

type exclusiveSync struct {
	mutex    sync.RWMutex
	awaiters exclusiveQueueController
}

func (p *exclusiveSync) CheckState() BoolDecision {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return BoolDecision(p.awaiters.isEmpty())
}

func (p *exclusiveSync) CheckDependency(dep SlotDependency) Decision {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if entry, ok := dep.(*dependencyQueueEntry); ok {
		switch {
		case !entry.link.IsValid(): // just to make sure
			return Impossible
		case !p.awaiters.contains(entry):
			return Impossible
		case p.awaiters.isEmptyOrFirst(entry.link):
			return Passed
		default:
			return NotPassed
		}
	}
	return Impossible
}

func (p *exclusiveSync) UseDependency(dep SlotDependency, flags SlotDependencyFlags) Decision {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if entry, ok := dep.(*dependencyQueueEntry); ok {
		switch {
		case !entry.link.IsValid(): // just to make sure
			return Impossible
		case !p.awaiters.contains(entry):
			return Impossible
		case !entry.IsCompatibleWith(flags):
			return Impossible
		case p.awaiters.isEmptyOrFirst(entry.link):
			return Passed
		default:
			return NotPassed
		}
	}
	return Impossible
}

func (p *exclusiveSync) CreateDependency(holder SlotLink, flags SlotDependencyFlags) (BoolDecision, SlotDependency) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	sd := p.awaiters.queue.addSlotForExclusive(holder, flags)
	if f, _ := p.awaiters.queue.FirstValid(); f == sd {
		return true, sd
	}
	return false, sd
}

func (p *exclusiveSync) GetCounts() (active, inactive int) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

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

func (p *exclusiveSync) AdjustLimit(limit int, absolute bool) (deps []StepLink, activate bool) {
	panic("illegal state")
}

func (p *exclusiveSync) EnumQueues(fn EnumQueueFunc) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.awaiters.enum(1, fn)
}

var _ DependencyQueueController = &exclusiveQueueController{}

type exclusiveQueueController struct {
	mutex *sync.RWMutex
	queueControllerTemplate
}

func (p *exclusiveQueueController) Init(name string, mutex *sync.RWMutex, controller DependencyQueueController) {
	p.queueControllerTemplate.Init(name, mutex, controller)
	p.mutex = mutex
}

func (p *exclusiveQueueController) IsOpen(sd SlotDependency) bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.queue.First() == sd
}

func (p *exclusiveQueueController) Release(link SlotLink, flags SlotDependencyFlags, removeFn func()) ([]PostponedDependency, []StepLink) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if f := p.queue.First(); f == nil || f.link != link {
		removeFn()
		return nil, nil
	}

	removeFn()
	switch f, step := p.queue.FirstValid(); {
	case f == nil:
		return nil, nil
	case f.stacker != nil:
		if postponed := f.stacker.ActivateStack(f, step); postponed != nil {
			return []PostponedDependency{postponed}, nil
		}
		fallthrough
	default:
		return nil, []StepLink{step}
	}
}

func (p *exclusiveQueueController) enum(qId int, fn EnumQueueFunc) bool {
	item := p.queue.head.QueueNext()
	if item == nil {
		return false
	}

	_, flags := item.getFlags()
	if fn(qId, item.link, flags) {
		return true
	}
	qId--

	for item = item.QueueNext(); item != nil; item = item.QueueNext() {
		_, flags := item.getFlags()
		if fn(qId, item.link, flags) {
			return true
		}
	}
	return false
}
