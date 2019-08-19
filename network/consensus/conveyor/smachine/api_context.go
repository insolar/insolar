///
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
///

package smachine

import (
	"context"
	"time"
)

type InitFunc func(ctx InitializationContext) StateUpdate
type StateFunc func(ctx ExecutionContext) StateUpdate
type MigrateFunc func(ctx MigrationContext) StateUpdate
type CreateFunc func(ctx ConstructionContext) StateMachine
type AsyncResultFunc func(ctx AsyncResultContext)

type BasicContext interface {
	GetSlotID() SlotID
	GetParent() SlotLink
}

type ConstructionContext interface {
	BasicContext
}

type syncContext interface {
	BasicContext

	GetSelf() SlotLink

	//SetMigrationForStep(fn MigrateFunc)
	SetMigration(fn MigrateFunc)

	//NextWithMigrate(StateFunc, MigrateFunc) StateUpdate
	Next(StateFunc) StateUpdate
	Stop() StateUpdate
}

type InitializationContext interface {
	syncContext
}

type MigrationContext interface {
	syncContext

	Replace(CreateFunc) StateUpdate
	Same() StateUpdate
}

type ExecutionContext interface {
	syncContext

	NewChild(CreateFunc) SlotLink

	/* In-state call to adapter */
	AdapterSyncCall(a ExecutionAdapter, fn AdapterCallFunc) bool
	AdapterAsyncCall(a ExecutionAdapter, fn AdapterCallFunc) context.CancelFunc
	NextAdapterCall(a ExecutionAdapter, fn AdapterCallFunc, resultState StateFunc) (StateUpdate, context.CancelFunc)

	Replace(CreateFunc) StateUpdate
	WaitAny() StateUpdate
	Yield() StateUpdate
	Repeat(limit int) StateUpdate

	//Before(d time.Time) StateUpdate
}

type AsyncResultContext interface {
	BasicContext

	WakeUp()
}

const UnknownSlotID SlotID = 0

type SlotID uint32

func (id SlotID) IsUnknown() bool {
	return id == UnknownSlotID
}

type stateUpdateFlags uint8

const (
	stateUpdateNoChange stateUpdateFlags = iota
	stateUpdateRepeat
	stateUpdateNext
	stateUpdateReplace
	stateUpdateStop
	stateUpdateColdWait
	stateUpdateHotWait
	stateUpdateFailed

	stateUpdateHasAsync stateUpdateFlags = 1 << 6
	stateUpdateYield    stateUpdateFlags = 1 << 7
	stateUpdateMask     stateUpdateFlags = 0x0F
)

type SlotStep struct {
	transition StateFunc
	migration  MigrateFunc
}

type StateUpdate struct {
	marker     *struct{}
	nextStep   SlotStep
	flags      stateUpdateFlags
	wakeupTime time.Time
	param      interface{}
	//prepare    func()
	//nextCreate CreateFunc
	//param0     int
}

func (u StateUpdate) getCreateFn() CreateFunc {
	return u.param.(CreateFunc)
}

func (u StateUpdate) getInt() int {
	if u.param == nil {
		return 0
	}
	return u.param.(int)
}

func (u StateUpdate) getPrepare() func() {
	return u.param.(func())
}

func (u StateUpdate) isEmpty() bool {
	return u.marker == nil
}

func (u *StateUpdate) setYield() {
	u.flags |= stateUpdateYield
}

func (u StateUpdate) getMode() stateUpdateFlags {
	return u.flags & stateUpdateMask
}

func (u StateUpdate) ensureContext(p *struct{}) StateUpdate {
	if u.marker != p {
		panic("illegal value")
	}
	return u
}

func (u StateUpdate) hasAny(flag stateUpdateFlags) bool {
	return u.flags&flag != 0
}
