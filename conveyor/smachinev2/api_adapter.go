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

type AdapterCallFunc func() AsyncResultFunc
type AdapterNotifyFunc func()
type CreateFactoryFunc func(eventPayload interface{}) CreateFunc

type AsyncCallRequester interface {
	/* Allocates and provides cancellation function. Repeated call returns same. */
	WithCancel(*context.CancelFunc) AsyncCallRequester
	// Sets a handler to map nested calls from the target adapter to new SMs
	// If this handler is nil or returns nil, then a default handler of the adapter will be in use.
	// To block a nested event - return non-nil CreateFunc, and then return nil from CreateFunc.
	WithNested(CreateFactoryFunc) AsyncCallRequester
	// See AsyncCallFlags
	WithFlags(flags AsyncCallFlags) AsyncCallRequester

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
	/*
		Call stays valid for this step (where the call is made) and for a next step.
		When SM will went further, the call or its result will be cancelled / ignored.
		NB! This cancel functionality is PASSIVE, an adapter should check this status explicitly.
	*/
	CallBoundToStep AsyncCallFlags = iota << 1
	/*
		When set, a wakeup from call's result will be valid for this step (where the call is made) and for a next step.
	*/
	WakeUpBoundToStep
	/*
		When set, receiving of call's successful result will wake up the slot without WakeUp().
		Behavior of this flag is also affected by WakeUpBoundToStep.
	*/
	AutoWakeUp
)

type SyncCallRequester interface {
	// Sets a handler to map nested calls from the target adapter to new SMs.
	// See AsyncCallRequester.WithNested() for details.
	WithNested(CreateFactoryFunc) AsyncCallRequester

	/* Returns true when the call was successful. May return false on a signal - depends on context mode */
	TryCall() bool
	/* May panic on migrate - depends on context mode */
	Call()
}

/* Provided by adapter's internals */
type AdapterExecutor interface {
	/*
		Schedules asynchronous execution, MAY return native cancellation function if supported.
		Panics are handled by caller.
	*/
	StartCall(fn AdapterCallFunc, callback *AdapterCallback, requireCancel bool) context.CancelFunc

	/*
		Schedules asynchronous, fire-and-forget execution.
		Panics are handled by caller.
	*/
	SendNotify(AdapterNotifyFunc)

	/*
		    Performs sync call ONLY if *natively* supported by the adapter, otherwise must return (false, nil)
			Panics are handled by caller.
	*/
	TrySyncCall(AdapterCallFunc) (bool, AsyncResultFunc)
}

/* This is interface of a helper to facilitate implementation of service adapters. */
type ExecutionAdapter interface {
	GetAdapterID() AdapterID
	PrepareSync(ctx ExecutionContext, fn AdapterCallFunc) SyncCallRequester
	PrepareAsync(ctx ExecutionContext, fn AdapterCallFunc) AsyncCallRequester
	PrepareNotify(ctx ExecutionContext, fn AdapterNotifyFunc) NotifyRequester
}
