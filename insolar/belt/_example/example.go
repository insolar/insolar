package _example

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/handler"
)

// Adapters are tasks that can execute themselves.
// IMPORTANT: "Adapt" function MUST wait for all spawned goroutines to finish.

type CallGorund struct {
	Request string
	Token   string

	externalCall chan string
	done         chan struct{}
	// We can group return parameters in build-in struct for clarity.
	result struct {
		externalCall string
		result       string
	}
}

func (a *CallGorund) Adapt(context.Context) bool {
	if a.done == nil {
		a.done = make(chan struct{})

		// This is a simulation of external call. It will not be in the client code.
		a.externalCall = make(chan string)
		go func() {
			for range time.NewTicker(100 * time.Millisecond).C {
				a.externalCall <- "hello!"
			}
		}()

		// Long operation simulation
		go func() {
			<-time.After(time.Second * 2)
			close(a.done)
		}()
		return false
	}

	select {
	case <-a.done:
		return true
	case a.result.externalCall = <-a.externalCall:
		return false
	}
}

type GetToken struct {
	result struct {
		Token string
	}
}

func (a *GetToken) Adapt(context.Context) bool {
	a.result.Token = "I allow!"
	return true
}

type CheckPermissions struct {
	node string

	// We can group return parameters in build-in struct for clarity.
	result struct {
		allowedToSave bool
	}
}

func (a *CheckPermissions) Adapt(context.Context) bool {
	// Check for node permissions.
	a.result.allowedToSave = true
	return true
}

type GetObjectFromDB struct {
	DBConnection *bytes.Buffer

	hash string

	result struct {
		id     int
		exists bool
	}
}

func (a *GetObjectFromDB) Adapt(context.Context) bool {
	b := make([]byte, 10)
	id, err := a.DBConnection.Read(b)
	if err != nil {
		return true
	}
	a.result.exists = true
	a.result.id = id

	return true
}

type SaveObjectToDB struct {
	DBConnection *bytes.Buffer

	hash string

	result struct {
		id  int
		err error
	}
}

func (a *SaveObjectToDB) Adapt(context.Context) bool {
	id, err := a.DBConnection.Write([]byte(a.hash))
	if err != nil {
		a.result.err = err
		return true
	}
	a.result.id = id
	return true
}

type SendReply struct {
	message string
}

func (a *SendReply) Adapt(context.Context) bool {
	// Send reply over network.
	return true
}

type Redirect struct {
	toNode string
}

func (a *Redirect) Adapt(context.Context) bool {
	// Redirect to other node.
	return true
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
	FLOW.Yield(s.migrateToPast, s.perms)
	FLOW.Yield(s.migrateToPast, s.object)

	if !s.perms.result.allowedToSave {
		FLOW.Yield(nil, &SendReply{message: "You shall not pass!"})
		return
	}

	if s.object.result.exists {
		FLOW.Yield(nil, &SendReply{message: "Object already exists"})
		return
	}

	saved := &SaveObjectToDB{hash: string(s.Message.Payload)}
	FLOW.Yield(s.migrateToPast, saved)

	if saved.result.err != nil {
		FLOW.Yield(nil, &SendReply{message: "Failed to save object"})
	} else {
		FLOW.Yield(nil, &SendReply{message: fmt.Sprintf("Object saved. ID: %d", saved.result.id)})
	}
}

func (s *SaveObject) Past(ctx context.Context, FLOW belt.Flow) {
	FLOW.Yield(nil, &SendReply{message: "Too late to save object"})
}

func (s *SaveObject) migrateToPast(ctx context.Context, FLOW belt.Flow) {
	if s.perms.result.allowedToSave {
		FLOW.Yield(nil, &SendReply{message: "You shall not pass!"})
	} else {
		FLOW.Yield(nil, &Redirect{toNode: "node that saves objects now"})
	}
}

type CallMethod struct {
	call  *CallGorund
	calls []string
}

func (s *CallMethod) Present(ctx context.Context, FLOW belt.Flow) {
	s.call = &CallGorund{Request: "where's the money Lebowski?!"}
	for FLOW.Yield(s.MigrateToPast, s.call) {
		s.calls = append(s.calls, s.call.result.externalCall)
	}
}

func (s *CallMethod) MigrateToPast(ctx context.Context, FLOW belt.Flow) {
	token := &GetToken{}
	FLOW.Yield(nil, token)
	s.call.Token = token.result.Token
	for FLOW.Yield(nil, s.call) {
		s.calls = append(s.calls, s.call.result.externalCall)
	}

	FLOW.Yield(nil, &SendReply{message: "My calls: " + strings.Join(s.calls, ",")})
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
