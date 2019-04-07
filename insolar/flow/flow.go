package flow

import (
	"context"

	"github.com/insolar/insolar/insolar/flow/bus"
)

// Handle is a one-function synchronous process that can call routines to do long processing.
//
// IMPORTANT: Asynchronous code is NOT ALLOWED here.
type Handle func(context.Context, Flow) error

// MakeHandle is a function that constructs new Handle.
type MakeHandle func(bus.Message) Handle

// Procedure is a task that can execute itself.
// It's a good idea to keep Procedures in a separate package to hide internal state from Handle.
type Procedure interface {
	// Proceed is called when Procedure is given control. When it returns, control will be given back to Handle.
	Proceed(context.Context) error
}

// Flow will be pasted to all Handles to control execution.
// This is very important not to blow this interface. Keep it minimal.
type Flow interface {
	// Handle gives control to another handle and waits for its return. Consider it "calling" another handler.
	// If cancellation happens during Handle execution, ErrCancelled will be returned.
	Handle(context.Context, Handle) error

	// Procedure starts routine and blocks Handle execution until cancellation happens or routine returns.
	// If cancellation happens first, ErrCancelled will be returned.
	// If Routine returns first, Procedure error (if any) will be returned.
	Procedure(context.Context, Procedure) error

	// Migrate blocks caller execution until cancellation happens then runs provided Handle in a new flow.
	// Note that this method can be called after cancellation. Use it to migrate processing after Handle or Procedure
	// returned ErrCancelled.
	//
	// IMPORTANT: Migrate can be called only once per flow. Calling it the second time will result in error.
	Migrate(context.Context, Handle) error
}
