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

type AdapterID string

func (v AdapterID) IsEmpty() bool {
	return len(v) == 0
}

type AdapterCallbackFunc func(AsyncResultFunc, error)
type AdapterCallFunc func() AsyncResultFunc
type AdapterNotifyFunc func()

/* This is a helper interface to facilitate implementation of service adapters */
type ExecutionAdapter interface {
	GetAdapterID() AdapterID
	PrepareSync(ctx ExecutionContext, fn AdapterCallFunc) SyncCallRequester
	PrepareAsync(ctx ExecutionContext, fn AdapterCallFunc) AsyncCallRequester
	// TODO PrepareNotify(ctx ExecutionContext, fn AdapterNotifyFunc) NotifyRequester
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

	/*
		When true will automatically cancel this call after the step is changed
		NB! This cancel functionality is PASSIVE and requires state to be checked by an executor of request.
	*/
	WithAutoCancelOnStep(attach bool) AsyncCallRequester

	// TODO WithNestedHandler
	// TODO With(mode AsyncCallFlags) AsyncCallRequester

	/* Starts async call  */
	Start()

	/* Creates an update that can be returned as a new state and will ONLY be executed if returned as a new state */
	DelayedStart() CallConditionalBuilder
}

type NotifyRequester interface {
	/* Sends notify */
	Send()

	/* Creates an update that can be returned as a new state and will ONLY be executed if returned as a new state */
	DelayedSend() CallConditionalBuilder
}

type AsyncCallFlags uint8

const (
	CallBoundToStep AsyncCallFlags = iota << 1
	WakeUpBoundToStep
	AutoWakeUp
)

//const WakeUpDisabled AsyncCallFlags = 0

/* Provided by adapter's internals */
type AdapterExecutor interface {
	/*
		Schedules asynchronous execution, MAY return native cancellation function if supported.
		Panics are handled by caller's function.
	*/
	StartCall(fn AdapterCallFunc, callback *AdapterCallback, requireCancel bool) context.CancelFunc

	SendNotify(AdapterCallFunc)

	/*
		    Performs sync call ONLY if *natively* supported by the adapter, otherwise must return (false, nil)
			Panics are handled by caller.
	*/
	TrySyncCall(AdapterCallFunc) (bool, AsyncResultFunc)
	//Migrate(slotMachineState SlotMachineState, migrationCount uint16)
}

type CreateFactoryFunc func(eventPayload interface{}) CreateFunc
