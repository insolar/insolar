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
	SetDefaultFlags(StepFlags)

	JumpOverride(StateFunc, MigrateFunc, StepFlags) StateUpdate
	Jump(StateFunc) StateUpdate

	Stop() StateUpdate
}

type InitializationContext interface {
	stepContext
}

type MigrationContext interface {
	stepContext

	Replace(CreateFunc) StateUpdate
	Stay() StateUpdate
}

type ExecutionContext interface {
	stepContext

	//ListenBroadcast(key string, broadcastFn BroadcastReceiveFunc)
	SyncOneStep(key string, weight int32, broadcastFn BroadcastReceiveFunc) Syncronizer
	//SyncManySteps(key string)

	NewChild(CreateFunc) SlotLink

	Replace(CreateFunc) StateUpdate
	Repeat(limit int) StateUpdate

	Yield() StateConditionalUpdate
	Poll() StateConditionalUpdate
	WaitForActive(SlotLink) StateConditionalUpdate
	Wait() StateConditionalUpdate
}

type StateConditionalUpdate interface {
	ConditionalUpdate
}

type CallConditionalUpdate interface {
	ConditionalUpdate
}

type ConditionalUpdate interface {
	ThenJump(StateFunc) StateUpdate
	ThenJumpOverride(StateFunc, MigrateFunc, StepFlags) StateUpdate
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

	WakeUp()
}

const UnknownSlotID SlotID = 0

type SlotID uint32

func (id SlotID) IsUnknown() bool {
	return id == UnknownSlotID
}

type StepFlags uint16

const (
	StepResetAllFlags StepFlags = 1 << iota
	StepIgnoreAsyncWakeup
	StepForceAsyncWakeup
	StepIgnoreAsyncPanic
	//StepAllowAsyncJump
)

type StateUpdate struct {
	marker  *struct{}
	link    SlotLink
	param   interface{}
	step    SlotStep
	updType uint16
}

func (u StateUpdate) IsZero() bool {
	return u.marker == nil && u.updType == 0
}

func NewStateUpdate(marker *struct{}, updType uint16, slotStep SlotStep, param interface{}) StateUpdate {
	return StateUpdate{
		marker:  marker,
		param:   param,
		step:    slotStep,
		updType: updType,
	}
}

func NewStateUpdateLink(marker *struct{}, updType uint16, link SlotLink, slotStep SlotStep, param interface{}) StateUpdate {
	return StateUpdate{
		marker:  marker,
		param:   param,
		link:    link,
		step:    slotStep,
		updType: updType,
	}
}

func EnsureUpdateContext(p *struct{}, u StateUpdate) StateUpdate {
	if u.marker != p {
		panic("illegal value")
	}
	return u
}

func ExtractStateUpdate(u StateUpdate) (updType uint16, slotStep SlotStep, param interface{}) {
	return u.updType, u.step, u.param
}

func ExtractStateUpdateParam(u StateUpdate) (updType uint16, param interface{}) {
	return u.updType, u.param
}

type SlotStep struct {
	Transition StateFunc
	Migration  MigrateFunc
	StepFlags  StepFlags
}

func (s *SlotStep) IsZero() bool {
	return s.Transition == nil && s.StepFlags == 0 && s.Migration == nil
}

func (s *SlotStep) HasTransition() bool {
	return s.Transition != nil
}
