package _example

import (
	"bytes"
	"context"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/flow/handler"
)

type GetObject struct {
	DBConnection *bytes.Buffer
}

func (s *GetObject) Past(context.Context, flow.Flow) error    { /* ... */ return nil }
func (s *GetObject) Present(context.Context, flow.Flow) error { /* ... */ return nil }
func (s *GetObject) Future(context.Context, flow.Flow) error  { /* ... */ return nil }

// =====================================================================================================================

func bootstrapExample() { // nolint
	DBConnection := bytes.NewBuffer(nil)

	hand := handler.NewHandler(
		// These functions can provide any variables via closure.
		// IMPORTANT: they must create NEW handle instances on every call.
		func(msg bus.Message) flow.Handle {
			s := GetObject{
				DBConnection: DBConnection,
			}
			return s.Present
		},
	)

	// Use handler to handle incoming messages.
	_ = hand
}
