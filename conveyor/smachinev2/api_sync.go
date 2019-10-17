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
	"fmt"

	"github.com/insolar/insolar/network/consensus/common/rwlock"
)

type SynchronizationContext interface {
	// Provides current state of a sync object.
	// When the sync was previously acquired, then this function returns SM's status of a sync object.
	// When the sync was not previously acquired, then this function returns status of
	// Panics on zero or incorrectly initialized value.
	Check(SyncLink) Decision

	// NB! This function RELEASES any previously acquired sync object before acquiring a new one.
	//
	// Acquires a holder of the sync object and returns status of the acquired holder:
	//
	// 1) Passed/true - SM can proceed to access resources controlled by this sync object.
	//    Passed holder MUST be released to ensure that other SM can also pass.
	//
	// 2) NotPassed/false - SM can't proceed to access resources controlled by this sync object.
	//    NotPassed holder remains valid and ensures that SM retains location an a queue of the sync object.
	//    NotPassed holder will at some moment converted into Passed holder and the relevant SM will be be woken up.
	//    NotPassed holder is MUST be released.
	//
	// Acquired holder will be released when SM is stopped.
	// Panics on zero or incorrectly initialized value.
	Acquire(SyncLink) BoolDecision

	// Similar to Acquire(), but the acquired holder will also be released when a step is changed.
	// To avoid doubt - Repeat(), WakeUp() and Stay() operations will not release.
	// Other operations, including Jump() to the same step will do RELEASE.
	// Panics on zero or incorrectly initialized value.
	AcquireForThisStep(SyncLink) BoolDecision

	// Releases a holder of this SM for the given sync object.
	// When there is no holder or the current holder belongs to a different sync object then operation is ignored and false is returned.
	// NB! Some sync objects (e.g. conditionals) may release a passed holder automatically, hence this function will return false as well.
	// Panics on zero or incorrectly initialized value.
	Release(SyncLink) bool

	// Releases a holder of this SM for any sync object if present.
	// Returns true when a holder of a sync object was released.
	// NB! Some sync objects (e.g. conditionals) may release a passed holder automatically, hence this function will return false as well.
	// Panics on zero or incorrectly initialized value.
	ReleaseAny() bool

	// Applies the given adjustment to a relevant sync object. SM doesn't need to acquire the relevant sync object.
	// Returns true when at least one holder of the sync object was affected.
	// Panics on zero or incorrectly initialized value.
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

// Represents a sync object.
type SyncLink struct {
	controller DependencyController
}

func (v SyncLink) IsZero() bool {
	return v.controller == nil
}

// Provides an implementation depended state of the sync object. Safe for concurrent use.
// Safe for concurrent use.
func (v SyncLink) GetCounts() (active, inactive int) {
	return v.controller.GetCounts()
}

// Provides an implementation depended state of the sync object
// Safe for concurrent use.
func (v SyncLink) GetLimit() (limit int, isAdjustable bool) {
	return v.controller.GetLimit()
}

func (v SyncLink) String() string {
	name := v.controller.GetName()
	if len(name) > 0 {
		return name
	}
	return fmt.Sprintf("sync-%p", v.controller)
}

/* ============================================== */

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

/* ============================================== */

type SlotDependencyFlags uint32

const (
	syncForOneStep SlotDependencyFlags = 1 << iota
)

// Internals of a sync object
type DependencyController interface {
	CheckState() Decision
	CheckDependency(dep SlotDependency) Decision
	UseDependency(dep SlotDependency, flags SlotDependencyFlags) Decision
	CreateDependency(slot *Slot, flags SlotDependencyFlags, syncer rwlock.RWLocker) (BoolDecision, SlotDependency)

	GetLimit() (limit int, isAdjustable bool)
	AdjustLimit(limit int, absolute bool) (deps []StepLink, activate bool)

	GetCounts() (active, inactive int)
	GetName() string
}
