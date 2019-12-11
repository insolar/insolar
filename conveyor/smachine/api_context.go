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
	"context"
	"time"
)

type InitFunc func(ctx InitializationContext) StateUpdate
type StateFunc func(ctx ExecutionContext) StateUpdate
type CreateFunc func(ctx ConstructionContext) StateMachine
type MigrateFunc func(ctx MigrationContext) StateUpdate
type AsyncResultFunc func(ctx AsyncResultContext)
type ErrorHandlerFunc func(ctx FailureContext)
type ErrorHandlerAction uint8

const (
	ErrorHandlerDefault ErrorHandlerAction = iota
	ErrorHandlerMute
	ErrorHandlerRecover
	ErrorHandlerRecoverAndWakeUp
	errorHandlerRecoverDenied
)

type BasicContext interface {
	SlotLink() SlotLink
	// can be empty - no parent
	ParentLink() SlotLink
	// not nil
	GetContext() context.Context
}

/*------------------  Contexts for in-order steps -----------------------*/

type DependencyInheritanceMode uint8

const (
	// Only local overrides and dependencies provided by SlotMachine will be used for injections.
	// And only factually injected dependencies will be inherited by children.
	IgnoreInheritedDependencies DependencyInheritanceMode = 0

	// for internal use
	copyAllDependencies DependencyInheritanceMode = 1

	// Injection will also use resolved dependencies from creator (NB! not from parent)
	// And only factually injected dependencies will be inherited by children.
	InheritResolvedDependencies DependencyInheritanceMode = 2

	// Injection will also use resolved dependencies from creator (NB! not from parent)
	// And all injected dependencies will be inherited by children.
	InheritAllDependencies DependencyInheritanceMode = 3

	// Can be combined with other modes to prevent inheritance by immediate children.
	// This does NOT affect dependencies provided by SlotMachine.
	DiscardResolvedDependencies DependencyInheritanceMode = 4
)

/* During construction SlotLink() will have correct SlotID, but MAY have INVALID status, as slot was not yet created */
type ConstructionContext interface {
	BasicContext

	// Puts a dependency for injector. Value can be nil
	OverrideDependency(id string, v interface{})
	// Precedence of dependencies (from the highest):
	// 1) overrides in this context
	// 2) provided via CreateDefaultValues
	// 3) inherited from creator (when allowed by inheritance mode)
	// 4) slot machine
	// 5) provided to slot machine
	SetDependencyInheritanceMode(DependencyInheritanceMode)

	SetContext(context.Context)
	SetParentLink(SlotLink)

	// Sets a special termination handler that will be invoked AFTER termination of SM.
	// This handler is invoked with either (1) GetDefaultTerminationResult() after Stop() or (2) with error after Error() or panic.
	// Behavior for error can be modified with a custom ErrorHandlerFunc.
	// This handler is not directly accessible to SM.
	// WARNING! This handler is UNSAFE to access another SM. Use BargeIn() to create a necessary handler.
	SetTerminationHandler(TerminationHandlerFunc)
	// Sets a default value to be passed to TerminationHandlerFunc when the slot stops.
	SetDefaultTerminationResult(interface{})

	// Sets tracing mode for the slot. Actual impact depends on implementation of a logger.
	SetLogTracing(bool)

	// Sets tracer id for the slot. This can't be directly changed after construction. See UpdateDefaultStepLogger()
	SetTracerId(TracerId)
}

type StepLoggerUpdateFunc func(StepLogger, StepLoggerFactoryFunc) StepLogger

/* A context parent for all regular step contexts */
type InOrderStepContext interface {
	BasicContext
	SynchronizationContext

	// Handler for migrations. Is applied when current SlotStep has no migration handler.
	// MUST be fast as it blocks whole SlotMachine and can't be detached.
	SetDefaultMigration(fn MigrateFunc)
	// Handler for errors and panics. Is applied when current SlotStep has no error handler.
	// MUST be fast as it blocks whole SlotMachine and can't be detached.
	SetDefaultErrorHandler(fn ErrorHandlerFunc)
	// Default flags are merged when SlotStep is set.
	SetDefaultFlags(StepFlags)
	// Sets a default value to be passed to TerminationHandlerFunc when the slot stops.
	SetDefaultTerminationResult(interface{})
	// Gets a value from the last SetDefaultTerminationResult().
	GetDefaultTerminationResult() interface{}

	// Marks this slot for priority on scheduling and sync queues. Overrides automatic boost.
	SetDynamicBoost(bool)

	// Returns a slot logger for this context. It is only valid while this context is valid.
	Log() Logger
	// Sets tracing mode for the slot. Actual impact depends on implementation of a logger.
	SetLogTracing(bool)
	// Overrides default step logger. Current logger is provided as argument. Update func can return nil.
	UpdateDefaultStepLogger(updateFn StepLoggerUpdateFunc)

	// Go to the next step. Flags, migrate and error handlers are provided by SetDefaultXXX()
	Jump(StateFunc) StateUpdate
	// Go to the next step with flags, migrate and error handlers.
	// Flags are merged with SetDefaultFlags() unless StepResetAllFlags is included.
	// Transition must not be nil, other handlers will use SetDefaultXXX() when nil
	JumpExt(SlotStep) StateUpdate

	// Creates a lazy link to the provided data. Link is invalidated when this SM is stopped.
	// This SM is always has a safe access when active. The shared data is guaranteed to be accessed by only one SM.
	// Access to the data is via ExecutionContext.UseShared().
	// Can be used across different SlotMachines.
	//
	// Do NOT share a reference to a field of SM with ShareDataDirect flag to avoid accidental memory leak.
	// It is recommended to use typed wrappers to access the data.
	Share(data interface{}, flags ShareDataFlags) SharedDataLink

	// Makes the data to be directly accessible via GetPublished().
	// Data is unpublished when this SM is stopped.
	// Visibility of key/data is limited by the SlotMachine running this SM.
	//
	// WARNING! There are NO safety guarantees. Publish only immutable data, e.g. publish SharedDataLink.
	// Returns false when key is in use.
	// It is recommended to use typed wrappers to access the data.
	Publish(key, data interface{}) bool
	// Returns false when key is not in use or the key was published by a different SM.
	// Is always able to unpublish link with ShareDataUnbound flag.
	Unpublish(key interface{}) bool
	// Removes all keys published by this SM.
	UnpublishAll()

	// Gets data shared by Publish().
	// Visibility of key/data is limited by the SlotMachine running this SM.
	// Returns nil when key is unknown or data is invalidated.
	// It is recommended to use typed wrappers to access the data.
	GetPublished(key interface{}) interface{}
	// Convenience wrapper for GetPublished(). Use SharedDataLink.IsXXX() to check availability.
	// It is recommended to use typed wrappers to access the data.
	GetPublishedLink(key interface{}) SharedDataLink

	// Publish this Slot globally under the given (key).
	// Published aliases will be unpublished on terminations of SM.
	// Returns false when key is in use.
	PublishGlobalAlias(key interface{}) bool
	// Unpublish the given (key)
	// Returns false when (key) is not published or is published by another slot.
	UnpublishGlobalAlias(key interface{}) bool
	// Get SlotLink for the given alias (key).
	// When (key) is unknown, then zero/empty SlotLink is returned.
	GetPublishedGlobalAlias(key interface{}) SlotLink

	// Slot will be terminated by calling an error handler.
	Error(error) StateUpdate
	Errorf(msg string, a ...interface{}) StateUpdate
	// Slot will be terminated.
	Stop() StateUpdate

	// Creates a barge-in function that can be used to signal or interrupt SM from outside.
	//
	// Provided BargeInParamFunc sends an async signal to the SM and will be ignored when SM has stopped.
	// When the signal is received by SM the BargeInApplyFunc is invoked. BargeInApplyFunc is safe to access SM.
	// BargeInParamFunc returns false when SM was stopped at the moment of the call.
	BargeInWithParam(BargeInApplyFunc) BargeInParamFunc

	// Provides a builder for a simple barge-in.
	BargeIn() BargeInBuilder
}

type ShareDataFlags uint32

const (
	// SM that called Share() will be woken up after each use of the shared data.
	ShareDataWakesUpAfterUse = 1 << iota

	// WARNING! Can ONLY be used for concurrency-safe data. Must NOT keep references to SM or its fields.
	// Data is immediately accessible. Data is not bound to SM and will never be invalidated.
	// Keeping such SharedDataLink will retain the data in memory.
	ShareDataUnbound

	// WARNING! Must NOT keep references to SM or its fields.
	// Data is bound to SM and will invalidated after stop.
	// But keeping such SharedDataLink will retain the data in memory even when invalidated.
	ShareDataDirect
)

type InitializationContext interface {
	InOrderStepContext
}

type PostInitStepContext interface {
	InOrderStepContext

	// Provides a builder for a simple barge-in. The barge-in function will be ignored if the step has changed.
	BargeInThisStepOnly() BargeInBuilder
}

type ExecutionContext interface {
	PostInitStepContext

	StepLink() StepLink
	GetPendingCallCount() int

	// WARNING! AVOID this method unless really needed.
	// The method forces detachment of this slot from SlotMachine's worker to allow slow processing and/or multiple sync calls.
	// Can only be called once per step. Detachment remains until end of the step.
	// Detached step will PREVENT access to any bound data shared by this SM.
	// To avoid doubt - detached step, like a normal step, will NOT receive async results, it can only receive result of sync calls.
	//
	// WARNING! SM with a detached step will NOT receive migrations until the detached step is finished.
	// Hence, SM may become inconsistent with other shared objects and injections that could be updated by migrations.
	//
	// Will panic when: (1) not supported by current worker, (2) detachment limit exceeded, (3) called repeatedly.
	InitiateLongRun(LongRunFlags)

	// Returns slot logger usable in async calls on adapter side. It will be able to report until the slot is invalidated.
	// This logger reports all events as at the step the logger was created at.
	LogAsync() Logger

	// Immediately allocates a new slot and constructs SM. And schedules initialization.
	// It is guaranteed that:
	// 1) the child will start at the same migration state as the creator (caller of this function)
	// 2) initialization of the new slot will happen before any migration
	NewChild(CreateFunc) SlotLink
	NewChildExt(CreateFunc, CreateDefaultValues) SlotLink

	// Same as NewChild, but also grantees that child's initialization will be completed before return.
	// Please prefer NewChild() to avoid unnecessary dependency.
	InitChild(CreateFunc) SlotLink
	InitChildExt(CreateFunc, CreateDefaultValues) SlotLink

	// After completion of the current step, SM will be stopped and the new SM created/started.
	// The new SM will by default inherit from this SM: parent, context, termination handler/result and injected dependencies.
	// When Replace() is successful, then stopping of this SM will not fire the termination handler.
	// WARNING! Use of SetTerminationHandler() inside CreateFunc will replace the current handler, so it will never fire then.
	Replace(CreateFunc) StateUpdate
	ReplaceExt(CreateFunc, CreateDefaultValues) StateUpdate
	// See Replace()
	ReplaceWith(StateMachine) StateUpdate

	// Applies the accessor produced by a SharedDataLink.
	// SharedDataLink can be used across different SlotMachines.
	UseShared(SharedDataAccessor) SharedAccessReport

	// Repeats current step (it is not considered as change of step).
	// The param limitPerCycle defines how many times this step will be repeated without switching to other slots unless interrupted.
	Repeat(limitPerCycle int) StateUpdate

	// SM will apply an action chosen by the builder and wait till next work cycle.
	Yield() StateConditionalBuilder
	// SM will apply an action chosen by the builder and wait for a poll interval (configured on SlotMachine).
	Poll() StateConditionalBuilder

	// EXPERIMENTAL! SM will apply an action chosen by the builder and wait for activation or stop of the given slot.
	// TODO WaitActivation(SlotLink) StateConditionalBuilder
	// TODO WaitStepChange(StepLink, tolerance uint32) StateConditionalBuilder

	// SM will apply an action chosen by the builder and wait for availability of the SharedDataLink.
	WaitShared(SharedDataLink) StateConditionalBuilder
	// SM will apply an action chosen by the builder and wait for any event (even for one irrelevant to this SM).
	WaitAny() StateConditionalBuilder
	// SM will apply an action chosen by the builder and wait for any event (even for one irrelevant to this SM), but not later than the given time.
	WaitAnyUntil(time.Time) StateConditionalBuilder

	// SM will apply an action chosen by the builder and wait for an explicit activation of this slot, e.g. any WakeUp() action.
	Sleep() StateConditionalBuilder
}

type LongRunFlags uint8

const (
	manualDetach LongRunFlags = 1 << iota
	IgnoreSignal
)

type MigrationContext interface {
	PostInitStepContext

	/* A step this SM is at during migration */
	AffectedStep() SlotStep

	// Indicates that multiple pending migrations can be skipped / do not need to be applied individually
	SkipMultipleMigrations()

	/* Keeps the last state */
	Stay() StateUpdate
	/* Makes SM active if it was waiting or polling */
	WakeUp() StateUpdate
}

type StateConditionalBuilder interface {
	ConditionalBuilder
	/* Returns information if the condition is already met */
	Decider
	// When the conditional requires wait, then return (Repeat(), true), otherwise returns ({}, false)
	ThenRepeatOrElse() (StateUpdate, bool)
	// When the conditional requires wait, then returns Repeat(), otherwise Jump()
	ThenRepeatOrJump(StateFunc) StateUpdate
	// When the conditional requires wait, then returns Repeat(), otherwise JumpExt()
	ThenRepeatOrJumpExt(SlotStep) StateUpdate
}

type CallConditionalBuilder interface {
	ConditionalBuilder
	Sleep() ConditionalBuilder
	Poll() ConditionalBuilder
	WaitAny() ConditionalBuilder
}

type ConditionalBuilder interface {
	ThenJump(StateFunc) StateUpdate
	ThenJumpExt(SlotStep) StateUpdate
	ThenRepeat() StateUpdate
}

/*------------------  Contexts for out-of-order steps -----------------------*/

type AsyncResultContext interface {
	BasicContext

	// Returns a slot logger for this context. It is only valid while this context is valid.
	Log() Logger

	// Makes SM active if it was waiting or polling
	WakeUp()
}

type BargeInApplyFunc func(BargeInContext) StateUpdate
type BargeInParamFunc func(interface{}) bool
type BargeInFunc func() bool

type BargeInBuilder interface {
	// BargeIn will change SM's step and wake it up
	WithJumpExt(SlotStep) BargeInFunc
	// BargeIn will change SM's step and wake it up
	WithJump(StateFunc) BargeInFunc
	// BargeIn will wake up SM at its current step
	WithWakeUp() BargeInFunc
	// BargeIn will stop SM
	WithStop() BargeInFunc
	// BargeIn will stop SM with the given error
	WithError(error) BargeInFunc
}

type BargeInContext interface {
	BasicContext

	// Returns a slot logger for this context. It is only valid while this context is valid.
	Log() Logger

	AffectedStep() SlotStep

	BargeInParam() interface{}

	// Returns true when SM step didn't change since barge-in creation
	IsAtOriginalStep() bool

	JumpExt(SlotStep) StateUpdate
	Jump(StateFunc) StateUpdate

	// Slot will be terminated by calling an error handler.
	Error(error) StateUpdate

	/* Keeps the last state */
	Stay() StateUpdate
	/* Makes active if was waiting or polling */
	WakeUp() StateUpdate

	Stop() StateUpdate
}

type FailureContext interface {
	BasicContext

	/* A step the slot is at */
	AffectedStep() SlotStep

	// Reason of the failure
	GetError() error

	// False when the error was initiated by ctx.Error(). When true, then GetError() should be SlotPanicError
	IsPanic() bool

	// True when this error can be recovered by SetAction(ErrorHandlerRecover).
	// A panic inside async call / callback can be recovered.
	CanRecover() bool

	// Gets a last value set by SetDefaultTerminationResult()
	GetDefaultTerminationResult() interface{}

	// Sets a value to be passed to TerminationHandlerFunc.
	// By default - termination result on error will be GetError()
	SetTerminationResult(interface{})

	// Choose an action to be applied.
	// Recovery actions will be ignored when CanRecover() is false.
	SetAction(action ErrorHandlerAction)

	// See ExecutionContext
	NewChild(CreateFunc) SlotLink
	// See ExecutionContext
	InitChild(CreateFunc) SlotLink
}
