package _example

import (
	"bytes"
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/belt/handler"
	"github.com/insolar/insolar/insolar/flow"
)

type GetObject struct {
	DBConnection *bytes.Buffer
}

func (s *GetObject) Past(context.Context, flow.Flow)    { /* ... */ }
func (s *GetObject) Present(context.Context, flow.Flow) { /* ... */ }
func (s *GetObject) Future(context.Context, flow.Flow)  { /* ... */ }

// =====================================================================================================================

func bootstrapExample() { // nolint
	DBConnection := bytes.NewBuffer(nil)

	hand := handler.NewHandler(
		// These functions can provide any variables via closure.
		// IMPORTANT: they must create NEW handle instances on every call.
		func(msg *message.Message) flow.Handle {
			// We can select required state for message here.
			s := GetObject{
				DBConnection: DBConnection,
			}
			return s.Past
		},
		func(msg *message.Message) flow.Handle {
			s := GetObject{
				DBConnection: DBConnection,
			}
			return s.Present
		},
		func(msg *message.Message) flow.Handle {
			s := GetObject{
				DBConnection: DBConnection,
			}
			return s.Future
		},
	)

	// This should be called by an event bus.
	_, _ = hand.HandleMessage(&message.Message{
		Payload: []byte(`{"hash": "clear is better than clever"}`),
	})
}
