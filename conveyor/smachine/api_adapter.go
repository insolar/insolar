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

type AdapterID uint32

type ExecutionAdapter interface {
	GetAdapterID() AdapterID
	PrepareSync(ctx ExecutionContext, fn AdapterCallFunc) SyncCallContext
	PrepareAsync(ctx ExecutionContext, fn AdapterCallFunc) CallContext
}

type SyncCallContext interface {
	TryCall() bool
	Call()
}

type CallContext interface {
	/* Allocates and provides cancellation function. Repeated calls return the same. */
	GetCancel(*context.CancelFunc) CallContext

	/* Will automatically cancel this call when step is changed */
	CancelOnStep(attach bool) CallContext

	/* Starts async call  */
	Start()

	/* Start async call that will try to do Jump after the result is returned and applied */
	Callback(fn StateFunc)
	CallbackWithMigrate(fn StateFunc, mf MigrateFunc)

	/* Creates an update that can be returned as a new state and will ONLY be executed if returned as a new state */
	Wait() CallConditionalUpdate
}

type AdapterCallbackFunc func(AsyncResultFunc)
type AdapterCallFunc func() AsyncResultFunc

type ExecutionAdapterSink interface {
	CallAsync(stepLink StepLink, fn AdapterCallFunc, callback AdapterCallbackFunc)
	CallAsyncWithCancel(stepLink StepLink, fn AdapterCallFunc, callback AdapterCallbackFunc) context.CancelFunc
}
