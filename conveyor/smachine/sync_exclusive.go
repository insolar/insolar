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

func (p *exclusiveSync) UseDependency(dep SlotDependency, flags SlotDependencyFlags) (Decision, SlotDependency) {
	if entry, ok := dep.(*dependencyQueueEntry); ok {
		switch {
		case !entry.link.IsValid(): // just to make sure
			return Impossible, nil
		case !p.awaiters.Contains(entry):
			return Impossible, nil
		case !entry.IsCompatibleWith(flags):
			return Impossible, nil
		case p.awaiters.IsEmptyOrFirst(entry.link):
			return Passed, nil
		default:
			return NotPassed, nil
		}
	}
	return Impossible, nil
}

func (p *exclusiveSync) CreateDependency(holder SlotLink, flags SlotDependencyFlags) (BoolDecision, SlotDependency) {
	sd := p.awaiters.queue.AddSlot(holder, flags)
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

func (p *exclusiveSync) AdjustLimit(limit int, absolute bool) (deps []StepLink, activate bool) {
	panic("illegal state")
}

var _ DependencyQueueController = &exclusiveQueueController{}

type exclusiveQueueController struct {
	queueControllerTemplate
}

func (p *exclusiveQueueController) Init(name string) {
	p.queueControllerTemplate.Init(name, p)
}

func (p *exclusiveQueueController) IsOpen(sd SlotDependency) bool {
	return p.queue.First() == sd
}

func (p *exclusiveQueueController) Release(link SlotLink, flags SlotDependencyFlags, removeFn func()) ([]PostponedDependency, []StepLink) {
	if f, _ := p.queue.FirstValid(); f == nil || f.link != link {
		removeFn()
		return nil, nil
	}

	removeFn()
	switch f, step := p.queue.FirstValid(); {
	case f == nil:
		return nil, nil
	case f.childOf != nil:
		if postponed := f.childOf.ActivateStack(f, step); postponed != nil {
			return []PostponedDependency{postponed}, nil
		}
		fallthrough
	default:
		return nil, []StepLink{step}
	}
}
