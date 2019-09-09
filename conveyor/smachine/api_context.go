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
	"context"
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
	GetContext() context.Context
}

type ConstructionContext interface {
	BasicContext
	SetContext(context.Context)
}

type stepContext interface {
	BasicContext

	SlotLink() SlotLink

	SetDefaultMigration(fn MigrateFunc)
	SetDefaultFlags(StepFlags)

	JumpExt(StateFunc, MigrateFunc, StepFlags) StateUpdate
	Jump(StateFunc) StateUpdate

	Stop() StateUpdate
}

type InitializationContext interface {
	stepContext

	BargeInWithParam(BargeInApplyFunc) BargeInParamFunc
	BargeIn() BargeInRequester
}

type BargeInApplyFunc func(BargeInContext) StateUpdate
type BargeInParamFunc func(interface{}) bool
type BargeInFunc func() bool

type BargeInRequester interface {
	WithJumpExt(StateFunc, MigrateFunc, StepFlags) BargeInFunc
	WithJump(StateFunc) BargeInFunc
	WithWakeUp() BargeInFunc
}

type MigrationContext interface {
	stepContext

	Replace(CreateFunc) StateUpdate
	/* Keeps the last state */
	Stay() StateUpdate
	/* Makes active if was waiting or polling */
	WakeUp() StateUpdate
}

type ExecutionContext interface {
	stepContext

	StepLink() StepLink
	GetPendingCallCount() int

	//ListenBroadcast(key string, broadcastFn BroadcastReceiveFunc)
	//SyncOneStep(key string, weight int32, broadcastFn BroadcastReceiveFunc) Syncronizer
	//SyncManySteps(key string)

	NewChild(context.Context, CreateFunc) SlotLink

	BargeInWithParam(BargeInApplyFunc) BargeInParamFunc
	BargeInThisStepOnly() BargeInRequester
	BargeIn() BargeInRequester

	Replace(CreateFunc) StateUpdate
	Repeat(limit int) StateUpdate

	Yield() StateConditionalUpdate
	Poll() StateConditionalUpdate
	WaitForActive(SlotLink) StateConditionalUpdate
	//	WaitForInput() StateConditionalUpdate
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
	ThenJumpExt(StateFunc, MigrateFunc, StepFlags) StateUpdate
	ThenRepeat() StateUpdate
}

//type Syncronizer interface {
//	IsFirst() bool
//	Broadcast(payload interface{}) (total, accepted int)
//	ReleaseAll()
//
//	Wait() StateUpdate
//	WaitOrDeadline(d time.Time) StateUpdate
//}

type AsyncResultContext interface {
	BasicContext

	WakeUp()
}

type BargeInContext interface {
	BasicContext

	GetBargeInParam() interface{}
	IsAtOriginalStep() bool

	/* A step the target slot is at */
	AffectedStep() SlotStep

	JumpExt(StateFunc, MigrateFunc, StepFlags) StateUpdate
	Jump(StateFunc) StateUpdate

	/* Keeps the last state */
	Stay() StateUpdate
	/* Makes active if was waiting or polling */
	WakeUp() StateUpdate
}
