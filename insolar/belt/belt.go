package belt

import (
	"context"

	"github.com/insolar/insolar/insolar/belt/bus"
)

// Handle is a one-function synchronous process that can call routines to do long processing.
// It's probably a good idea to use UPPER CASE in variable name to highlight FLOW control so it will be harder to miss.
//
// IMPORTANT: Asynchronous code is NOT ALLOWED here.
type Handle func(ctx context.Context, FLOW Flow)

// MakeHandle is a function that constructs new Handle.
type MakeHandle func(msg bus.Message) Handle

// Procedure is a task that can execute itself.
// It's a good idea to keep Procedures in a separate package to hide internal state from Handle.
type Procedure interface {
	// Proceed is called every time Procedure is given control. When it returns, control will be given back to Handle.
	// If Procedure requires Handle to make decision and give control back, it should return true.
	// Handle will make modifications in Procedure state and Proceed will be called again to continue procedure.
	// If false is returned, Procedure is considered complete.
	Proceed(context.Context)
}

// Flow will be pasted to all Handles to control execution.
// This is very important not to blow this interface. Keep it minimal.
type Flow interface {

	// Jump instantly changes Handle. It MUST be used instead of calling other Handles directly.
	Jump(to Handle)

	// Handle gives control to another handle and waits for its return. Consider it "calling" another handler.
	// If cancellation happens during execution of that handle, the callie will be migrated and the caller will be
	// interrupted.
	Handle(ctx context.Context, handle Handle)

	// Yield starts routine and blocks Handle execution until cancellation happens or routine returns.
	//
	// If cancellation happens first, Handle will be migrated.
	// If Routine returns first, Handle execution will continue.
	//
	// If Routine is nil, execution blocks until cancellation and migrates Handle. If Handle is nil, execution
	// interrupts immediately.
	//
	// Returns true if there is still work to do. Handle can decide to call Yield again to receive more results.
	Yield(migrate Handle, routine Procedure)
}
