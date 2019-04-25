//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package flow

import (
	"context"

	"github.com/insolar/insolar/insolar/flow/bus"
)

// Handle is a one-function synchronous process that can call routines to do long processing.
// IMPORTANT: Asynchronous code is NOT ALLOWED here.
// To create a new Handle of a given message use handler.NewHandler procedure.
// After creating a Handle you can register it in MessageBus like this:
// `h.Bus.MustRegister(insolar.TypeGetObject, createdHandle.WrapBusHandle)`
// You can find an example in insolar/ladger/artifactmanager/handler.go
type Handle func(context.Context, Flow) error

// MakeHandle is a function that constructs new Handle.
type MakeHandle func(bus.Message) Handle

//go:generate minimock -i github.com/insolar/insolar/insolar/flow.Procedure -o . -s _mock.go

// Procedure is a task that can execute itself.
// Please note that the Procedure is marked as canceled if a pulse happens during it's execution. This means that it
// continues to execute in the background, though it's return value will be discarded.
// Thus if you have multiple steps that can be executed in different pulses split them into separate Procedures.
// Otherwise join the steps into a single Procedure.
// It's a good idea to keep Procedures in a separate package to hide internal state from Handle.
type Procedure interface {
	// Proceed is called when Procedure is given control. When it returns, control will be given back to Handle.
	Proceed(context.Context) error
}

//go:generate minimock -i github.com/insolar/insolar/insolar/flow.Flow -o . -s _mock.go

// Flow will be pasted to all Handles to control execution.
// This is very important not to blow this interface. Keep it minimal.
type Flow interface {
	// Handle gives control to another handle and waits for its return. Consider it "calling" another handler.
	Handle(context.Context, Handle) error

	// Procedure starts a routine and blocks Handle execution until cancellation happens or routine returns.
	// If cancellation happens first, ErrCancelled will immediately be returned to the Handle. The Procedure
	// continues to execute in the background, but it's state must be discarded by the Handle as invalid.
	// If Routine returns first, Procedure error (if any) will be returned.
	// Procedure can figure out whether it's execution was canceled and there is no point to continue
	// the execution by reading from context.Done()
	Procedure(ctx context.Context, proc Procedure, cancelable bool) error

	// Migrate blocks caller execution until cancellation happens then runs provided Handle in a new flow.
	// Note that this method can be called after cancellation. Use it to migrate processing after Handle or Procedure
	// returned ErrCancelled.
	//
	// IMPORTANT: Migrate can be called only once per flow. Calling it the second time will result in error.
	Migrate(context.Context, Handle) error

	// Continue blocks caller execution until cancellation happens then updates 'cancel' and returns control to caller.
	// It might be called multiple times, but each time it will wait for cancellation.
	// Might be used to continue processing in Handle after Procedure returns ErrCancelled
	Continue(context.Context)
}
