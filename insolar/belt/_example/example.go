package _example

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/handler"
)

// Adapters are tasks that can execute themselves.
// IMPORTANT: "Adapt" function MUST wait for all spawned goroutines to finish.

type CheckPermissions struct {
	node string

	// We can group return parameters in build-in struct for clarity.
	result struct {
		allowedToSave bool
	}
}

func (a *CheckPermissions) Adapt(context.Context) {
	// Check for node permissions.
	a.result.allowedToSave = true
}

type GetObjectFromDB struct {
	DBConnection *bytes.Buffer

	hash string

	result struct {
		id     int
		exists bool
	}
}

func (a *GetObjectFromDB) Adapt(context.Context) {
	b := make([]byte, 10)
	id, err := a.DBConnection.Read(b)
	if err != nil {
		return
	}
	a.result.exists = true
	a.result.id = id
}

type SaveObjectToDB struct {
	DBConnection *bytes.Buffer

	hash string

	result struct {
		id  int
		err error
	}
}

func (a *SaveObjectToDB) Adapt(context.Context) {
	id, err := a.DBConnection.Write([]byte(a.hash))
	if err != nil {
		a.result.err = err
		return
	}
	a.result.id = id
}

type SendReply struct {
	message string
}

func (a *SendReply) Adapt(context.Context) {
	// Send reply over network.
}

type Redirect struct {
	toNode string
}

func (a *Redirect) Adapt(context.Context) {
	// Redirect to other node.
}

// =====================================================================================================================

// SaveObject describes handling "save object" message flow.
type SaveObject struct {
	DBConnection *bytes.Buffer
	Message      *message.Message

	// Keep internal state unexported.
	perms  *CheckPermissions
	object *GetObjectFromDB
}

// These functions represent message handling in different slots.
// It's probably a good idea to use UPPER CASE to highlight FLOW control so it will be harder to miss.
// IMPORTANT: Asynchronous code is NOT ALLOWED here.

func (s *SaveObject) Future(ctx context.Context, FLOW belt.Flow) {
	FLOW.Wait(s.Present)
}

func (s *SaveObject) Present(ctx context.Context, FLOW belt.Flow) {
	s.perms = &CheckPermissions{node: s.Message.Metadata["node"]}
	s.object = &GetObjectFromDB{hash: string(s.Message.Payload)}
	FLOW.YieldAll(s.migrateToPast, s.perms, s.object)

	if !s.perms.result.allowedToSave {
		FLOW.YieldNone(nil, &SendReply{message: "You shall not pass!"})
		return
	}

	if s.object.result.exists {
		FLOW.YieldNone(nil, &SendReply{message: "Object already exists"})
		return
	}

	saved := &SaveObjectToDB{hash: string(s.Message.Payload)}
	FLOW.YieldAll(s.migrateToPast, saved)

	if saved.result.err != nil {
		FLOW.YieldNone(nil, &SendReply{message: "Failed to save object"})
	} else {
		FLOW.YieldNone(nil, &SendReply{message: fmt.Sprintf("Object saved. ID: %d", saved.result.id)})
	}
}

func (s *SaveObject) Past(ctx context.Context, FLOW belt.Flow) {
	FLOW.YieldNone(nil, &SendReply{message: "Too late to save object"})
}

func (s *SaveObject) migrateToPast(ctx context.Context, FLOW belt.Flow) {
	if s.perms.result.allowedToSave {
		FLOW.YieldNone(nil, &SendReply{message: "You shall not pass!"})
	} else {
		FLOW.YieldNone(nil, &Redirect{toNode: "node that saves objects now"})
	}
}

// =====================================================================================================================

func bootstrapExample() { // nolint
	DBConnection := bytes.NewBuffer(nil)

	hand := handler.NewHandler(
		insolar.Pulse{PulseNumber: 42},
		// These functions can provide any variables via closure.
		// IMPORTANT: they must create NEW handle instances on every call.
		func(msg *message.Message) belt.Handle {
			// We can select required state for message here.
			s := SaveObject{
				DBConnection: DBConnection,
				Message:      msg,
			}
			return s.Past
		},
		func(msg *message.Message) belt.Handle {
			s := SaveObject{
				DBConnection: DBConnection,
				Message:      msg,
			}
			return s.Present
		},
		func(msg *message.Message) belt.Handle {
			s := SaveObject{
				DBConnection: DBConnection,
				Message:      msg,
			}
			return s.Future
		},
	)

	// This should be called by an event bus.
	_, _ = hand.HandleMessage(&message.Message{
		Payload: []byte(`{"hash": "clear is better than clever"}`),
	})
}
