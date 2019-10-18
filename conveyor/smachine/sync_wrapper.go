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
	"github.com/insolar/insolar/network/consensus/common/rwlock"
	"sync"
)

var _ DependencyController = &syncMutexWrapper{}

type syncMutexWrapper struct {
	mutex sync.RWMutex
	inner DependencyController
}

func (w *syncMutexWrapper) CheckState() Decision {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.inner.CheckState()
}

func (w *syncMutexWrapper) CheckDependency(dep SlotDependency) Decision {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.inner.CheckDependency(dep)
}

func (w *syncMutexWrapper) UseDependency(dep SlotDependency, flags SlotDependencyFlags) Decision {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.inner.UseDependency(dep, flags)
}

func (w *syncMutexWrapper) CreateDependency(slot *Slot, flags SlotDependencyFlags, syncer rwlock.RWLocker) (BoolDecision, SlotDependency) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.inner.CreateDependency(slot, flags, &w.mutex)
}

func (w *syncMutexWrapper) GetLimit() (limit int, isAdjustable bool) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.inner.GetLimit()
}

func (w *syncMutexWrapper) AdjustLimit(limit int, absolute bool) (deps []StepLink, activate bool) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	return w.inner.AdjustLimit(limit, absolute)
}

func (w *syncMutexWrapper) GetCounts() (active, inactive int) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.inner.GetCounts()
}

func (w *syncMutexWrapper) GetName() string {
	return w.inner.GetName()
}
