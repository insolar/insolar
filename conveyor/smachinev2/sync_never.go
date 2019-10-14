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

func NewInfiniteLock(name string) SyncLink {
	return NewSyncLink(&infiniteLock{name: name})
}

type infiniteLock struct {
	name  string
	count int
}

func (p *infiniteLock) CheckState() Decision {
	return NotPassed
}

func (p *infiniteLock) CheckDependency(dep SlotDependency) Decision {
	if entry, ok := dep.(*infiniteLockEntry); ok && entry.ctl == p {
		return NotPassed
	}
	return Impossible
}

func (p *infiniteLock) UseDependency(dep SlotDependency, oneStep bool) Decision {
	if entry, ok := dep.(*infiniteLockEntry); ok {
		switch {
		case !oneStep && (entry.slotFlags&syncForOneStep != 0):
			return Impossible
		case entry.ctl == p:
			return NotPassed
		}
	}
	return Impossible
}

func (p *infiniteLock) CreateDependency(slot *Slot, oneStep bool) (Decision, SlotDependency) {
	flags := DependencyQueueEntryFlags(0)
	if oneStep {
		flags |= syncForOneStep
	}
	p.count++
	return NotPassed, &infiniteLockEntry{p, flags}
}

func (p *infiniteLock) GetLimit() (limit int, isAdjustable bool) {
	return 0, false
}

func (p *infiniteLock) AdjustLimit(limit int) ([]SlotLink, bool) {
	panic("illegal state")
}

func (p *infiniteLock) GetWaitingCount() int {
	return p.count
}

func (p *infiniteLock) GetName() string {
	return p.name
}

var _ SlotDependency = &infiniteLockEntry{}

type infiniteLockEntry struct {
	ctl       *infiniteLock
	slotFlags DependencyQueueEntryFlags
}

func (v infiniteLockEntry) IsReleaseOnStepping() bool {
	return v.slotFlags&syncForOneStep != 0
}

func (infiniteLockEntry) IsReleaseOnWorking() bool {
	return false
}

func (v infiniteLockEntry) Release(_ func(SlotLink)) {
	v.ctl.count--
}

func (v infiniteLockEntry) ReleaseOnDisposed(_ func(SlotLink)) {
	v.ctl.count--
}
