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

import "github.com/insolar/insolar/network/consensus/common/rwlock"

type SynchronizationContext interface {
	Check(SyncLink) Decision
	AcquireForThisStep(SyncLink) BoolDecision
	Acquire(SyncLink) BoolDecision
	Release(SyncLink) bool
	ApplyAdjustment(SyncAdjustment) bool
}

func NewSyncLink(controller DependencyController) SyncLink {
	if controller == nil {
		panic("illegal value")
	}
	return SyncLink{&syncMutexWrapper{inner: controller}}
}

func NewSyncLinkNoLock(controller DependencyController) SyncLink {
	if controller == nil {
		panic("illegal value")
	}
	return SyncLink{controller}
}

type SyncLink struct {
	controller DependencyController
}

func (v SyncLink) IsZero() bool {
	return v.controller == nil
}

type SyncAdjustment struct {
	controller DependencyController
	adjustment int
	isAbsolute bool
}

func (v SyncAdjustment) IsZero() bool {
	return v.controller == nil
}

func (v SyncAdjustment) IsEmpty() bool {
	return v.controller == nil || !v.isAbsolute && v.adjustment == 0
}

func (v SyncLink) GetCounts() (active, inactive int) {
	return v.controller.GetCounts()
}

type SlotDependencyFlags uint32

const (
	syncForOneStep SlotDependencyFlags = 1 << iota
)

type DependencyController interface {
	CheckState() Decision
	CheckDependency(dep SlotDependency) Decision
	UseDependency(dep SlotDependency, flags SlotDependencyFlags) Decision
	CreateDependency(slot *Slot, flags SlotDependencyFlags, syncer rwlock.RWLocker) (BoolDecision, SlotDependency)

	GetLimit() (limit int, isAdjustable bool)
	AdjustLimit(limit int) (deps []StepLink, activate bool)

	GetCounts() (active, inactive int)
	GetName() string
}
