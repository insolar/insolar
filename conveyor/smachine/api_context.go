///
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
