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

type SynchronizationContext interface {
	Check(SyncLink) Decision
	AcquireForThisStep(SyncLink) Decision
	Acquire(SyncLink) Decision
	Release(SyncLink) bool
}

func NewSyncLink(controller DependencyController) SyncLink {
	return SyncLink{controller}
}

type SyncLink struct {
	controller DependencyController
}

func (v SyncLink) GetQueueCount() int {
	return v.controller.GetWaitingCount()
}

type DependencyController interface {
	CheckState() Decision
	CheckDependency(dep SlotDependency) Decision
	UseDependency(dep SlotDependency, oneStep bool) Decision
	CreateDependency(slot *Slot, oneStep bool) (Decision, SlotDependency)

	GetLimit() (limit int, isAdjustable bool)
	AdjustLimit(limit int) (deps []SlotLink, activate bool)

	GetWaitingCount() int
	GetName() string
}
