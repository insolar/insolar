package _example

import (
	"bytes"
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/handler"
)

type GetObject struct {
	DBConnection *bytes.Buffer
}

func (s *GetObject) Past(context.Context, belt.Flow)    { /* ... */ }
func (s *GetObject) Present(context.Context, belt.Flow) { /* ... */ }
func (s *GetObject) Future(context.Context, belt.Flow)  { /* ... */ }

// =====================================================================================================================

func bootstrapExample() { // nolint
	DBConnection := bytes.NewBuffer(nil)

	hand := handler.NewHandler(
		// These functions can provide any variables via closure.
		// IMPORTANT: they must create NEW handle instances on every call.
		func(msg *message.Message) belt.Handle {
			// We can select required state for message here.
			s := GetObject{
				DBConnection: DBConnection,
			}
			return s.Past
		},
		func(msg *message.Message) belt.Handle {
			s := GetObject{
				DBConnection: DBConnection,
			}
			return s.Present
		},
		func(msg *message.Message) belt.Handle {
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
