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
	"sync/atomic"
)

func NewInfiniteLock(name string) SyncLink {
	return NewSyncLink(&infiniteLock{name: name})
}

type infiniteLock struct {
	name  string
	count int32 //atomic
}

func (p *infiniteLock) CheckState() BoolDecision {
	return false
}

func (p *infiniteLock) CheckDependency(dep SlotDependency) Decision {
	if entry, ok := dep.(*infiniteLockEntry); ok && entry.ctl == p {
		return NotPassed
	}
	return Impossible
}

func (p *infiniteLock) UseDependency(dep SlotDependency, flags SlotDependencyFlags) Decision {
	if entry, ok := dep.(*infiniteLockEntry); ok {
		switch {
		case !entry.IsCompatibleWith(flags):
			return Impossible
		case entry.ctl == p:
			return NotPassed
		}
	}
	return Impossible
}

func (p *infiniteLock) CreateDependency(holder SlotLink, flags SlotDependencyFlags) (BoolDecision, SlotDependency) {
	atomic.AddInt32(&p.count, 1)
	return false, &infiniteLockEntry{p, flags}
}

func (p *infiniteLock) GetLimit() (limit int, isAdjustable bool) {
	return 0, false
}

func (p *infiniteLock) AdjustLimit(limit int, absolute bool) ([]StepLink, bool) {
	panic("illegal state")
}

func (p *infiniteLock) GetCounts() (active, inactive int) {
	return 0, int(p.count)
}

func (p *infiniteLock) GetName() string {
	return p.name
}

func (p *infiniteLock) EnumQueues(fn EnumQueueFunc) bool {
	return false
}

var _ SlotDependency = &infiniteLockEntry{}

type infiniteLockEntry struct {
	ctl       *infiniteLock
	slotFlags SlotDependencyFlags
}

func (v infiniteLockEntry) IsReleaseOnStepping() bool {
	return v.slotFlags&syncForOneStep != 0
}

func (infiniteLockEntry) IsReleaseOnWorking() bool {
	return false
}

func (v infiniteLockEntry) Release() (SlotDependency, []PostponedDependency, []StepLink) {
	v.ReleaseAll()
	return nil, nil, nil
}

func (v infiniteLockEntry) ReleaseAll() ([]PostponedDependency, []StepLink) {
	atomic.AddInt32(&v.ctl.count, -1)
	return nil, nil
}

func (v infiniteLockEntry) IsCompatibleWith(requiredFlags SlotDependencyFlags) bool {
	return v.slotFlags.isCompatibleWith(requiredFlags)
}
