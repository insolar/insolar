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
	"time"
)

type InitFunc func(ctx InitializationContext) StateUpdate
type StateFunc func(ctx ExecutionContext) StateUpdate
type MigrateFunc func(ctx MigrationContext) StateUpdate
type CreateFunc func(ctx ConstructionContext) StateMachine
type AsyncResultFunc func(ctx AsyncResultContext)
type BroadcastReceiveFunc func(ctx AsyncResultContext, payload interface{}) bool

type BasicContext interface {
	GetSlotID() SlotID
	GetParent() SlotLink
}

type ConstructionContext interface {
	BasicContext
}

type stepContext interface {
	BasicContext

	GetSelf() SlotLink

	SetDefaultMigration(fn MigrateFunc)

	JumpWithMigrate(StateFunc, MigrateFunc) StateUpdate
	Jump(StateFunc) StateUpdate
	Stop() StateUpdate
}

type InitializationContext interface {
	stepContext
}

type MigrationContext interface {
	stepContext

	Replace(CreateFunc) StateUpdate
	Same() StateUpdate
}

type ExecutionContext interface {
	stepContext

	//ListenBroadcast(key string, broadcastFn BroadcastReceiveFunc)
	SyncOneStep(key string, weight int32, broadcastFn BroadcastReceiveFunc) Syncronizer
	//SyncManySteps(key string)

	NewChild(CreateFunc) SlotLink

	Replace(CreateFunc) StateUpdate
	Repeat(limit int) StateUpdate
	Yield() StateUpdate

	WaitAny() StateConditionalUpdate
}

type conditionalUpdateAction interface {
	WakeUp(enable bool) ConditionalUpdate
	WakeUpAlways() ConditionalUpdate
}

type StateConditionalUpdate interface {
	CallConditionalUpdate
	AsyncJump(enable bool) CallConditionalUpdate
}

type CallConditionalUpdate interface {
	conditionalUpdateAction
	Poll() ConditionalUpdate
	Active(slot SlotLink) ConditionalUpdate
}

type ConditionalUpdate interface {
	conditionalUpdateAction

	ThenJump(StateFunc) StateUpdate
	ThenJumpWithMigrate(StateFunc, MigrateFunc) StateUpdate
	ThenRepeat() StateUpdate
}

type Syncronizer interface {
	IsFirst() bool
	Broadcast(payload interface{}) (total, accepted int)
	ReleaseAll()

	Wait() StateUpdate
	WaitOrDeadline(d time.Time) StateUpdate
}

type AsyncResultContext interface {
	BasicContext

	// caller will execute its current step
	WakeUp()
	//WakeUpAndJump()
}

const UnknownSlotID SlotID = 0

type SlotID uint32

func (id SlotID) IsUnknown() bool {
	return id == UnknownSlotID
}

type StateUpdate struct {
	marker *struct{}
	step   SlotStep
	link   SlotLink
	param  interface{}
}

func (u StateUpdate) IsZero() bool {
	return u.marker == nil && u.step.paramTmp == 0
}

func NewStateUpdate(marker *struct{}, updType uint16, slotStep SlotStep, param interface{}) StateUpdate {
	slotStep.paramTmp = updType
	return StateUpdate{
		marker: marker,
		param:  param,
		step:   slotStep,
	}
}

func NewStateUpdateLink(marker *struct{}, updType uint16, link SlotLink, slotStep SlotStep, param interface{}) StateUpdate {
	slotStep.paramTmp = updType
	return StateUpdate{
		marker: marker,
		param:  param,
		link:   link,
		step:   slotStep,
	}
}

func EnsureUpdateContext(p *struct{}, u StateUpdate) StateUpdate {
	if u.marker != p {
		panic("illegal value")
	}
	return u
}

func ExtractStateUpdate(u StateUpdate) (updType uint16, slotStep SlotStep, param interface{}) {
	t := u.step.paramTmp
	u.step.paramTmp = 0
	return t, u.step, u.param
}

func ExtractStateUpdateParam(u StateUpdate) (updType uint16, param interface{}) {
	return u.step.paramTmp, u.param
}

type SlotStep struct {
	Transition StateFunc
	Migration  MigrateFunc
	StepFlags  uint16
	paramTmp   uint16
}

func (s *SlotStep) IsZero() bool {
	return s.Transition == nil && s.StepFlags == 0 && s.paramTmp == 0
}

func (s *SlotStep) HasTransition() bool {
	return s.Transition != nil
}
