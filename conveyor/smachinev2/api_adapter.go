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

import "context"

type AdapterID string

func (v AdapterID) IsEmpty() bool {
	return len(v) == 0
}

/* This is a helper interface to facilitate implementation of service adapters */
type ExecutionAdapter interface {
	GetAdapterID() AdapterID
	PrepareSync(ctx ExecutionContext, fn AdapterCallFunc) SyncCallRequester
	PrepareAsync(ctx ExecutionContext, fn AdapterCallFunc) AsyncCallRequester
}

type SyncCallRequester interface {
	// TODO WithNestedHandler
	/* Returns true when the call was successful. Will return false when worker has a signal / interrupt */
	TryCall() bool
	/* Panics when it wasn't possible to perform a sync call */
	Call()
}

type AsyncCallRequester interface {
	/* Allocates and provides cancellation function. Repeated calls return the same. */
	WithCancel(*context.CancelFunc) AsyncCallRequester

	/* When true will automatically cancel this call after the step is changed */
	WithAutoCancelOnStep(attach bool) AsyncCallRequester

	// TODO WithNestedHandler
	// TODO With(mode AsyncCallFlags) AsyncCallRequester

	/* Starts async call  */
	Start()

	/* Creates an update that can be returned as a new state and will ONLY be executed if returned as a new state */
	DelayedStart() CallConditionalBuilder
}

type AsyncCallFlags uint8

const (
	CallBoundToStep AsyncCallFlags = iota << 1
	WakeUpBoundToStep
	AutoWakeUp
)

//const WakeUpDisabled AsyncCallFlags = 0

type NestedEventContext interface {
	NewChild()
	InitChild()
}

type NestedEventFunc func(precursor StepLink, ctx NestedEventContext)

type AdapterCallbackHandler interface {
	SendResult(AsyncResultFunc, error)
	// can be called ONLY until SendResult
	// TODO send from adapter?
	SendNestedEvent(NestedEventFunc)
}

type AdapterCallbackFunc func(AsyncResultFunc, error)
type AdapterCallFunc func() AsyncResultFunc
type AdapterNestedEventFunc func(precursor StepLink, eventPayload interface{}, requireCancel bool) context.CancelFunc

//type AdapterNestedEvent interface {
//	SendNestedEvent(precursor StepLink, eventPayload interface{}, requireCancel bool) context.CancelFunc
//}

/* Provided by adapter's internals */
type AdapterExecutor interface {
	/*
		Schedules asynchronous execution, MAY return native cancellation function if supported.
		When callback == nil then fn() must return nil as well.
		Panics are handled by caller.
	*/
	StartCall(stepLink StepLink, fn AdapterCallFunc, callback AdapterCallbackFunc, requireCancel bool) context.CancelFunc

	/*
		    Performs sync call ONLY if *natively* supported by the adapter, otherwise must return (false, nil)
			Panics are handled by caller.
	*/
	TrySyncCall(fn AdapterCallFunc) (bool, AsyncResultFunc)
	//Migrate(slotMachineState SlotMachineState, migrationCount uint16)
}

//type SharedStateAdapter interface {
//	PrepareUpdate(ctx ExecutionContext, fn func()) SharedUpdateRequester
//	TryCancel(ctx ExecutionContext)
//}
//
//type SharedUpdateRequester interface {
//	TryApply() (isValid, isApplied bool)
//	Apply()
//}
