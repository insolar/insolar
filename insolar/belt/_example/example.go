package _example

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/belt"
	"github.com/insolar/insolar/insolar/belt/handler"
)

type CheckPermissions struct {
	node string

	allowedToSave bool
}

func (a *CheckPermissions) Adapt(context.Context) {
	panic("implement me")
}

type GetObjectFromDB struct {
	hash string

	id     int
	exists bool
}

func (a *GetObjectFromDB) Adapt(context.Context) {
	panic("implement me")
}

type SaveObjectToDB struct {
	hash string

	id  int
	err error
}

func (a *SaveObjectToDB) Adapt(context.Context) {
	panic("implement me")
}

type SendReply struct {
	message string
}

func (a *SendReply) Adapt(context.Context) {
	panic("implement me")
}

type Redirect struct {
	toNode string
}

func (a *Redirect) Adapt(context.Context) {
	panic("implement me")
}

// =====================================================================================================================

// SaveObject describes handling "save object" message flow.
type SaveObject struct {
	// message message.Message
	perms  *CheckPermissions
	object *GetObjectFromDB
}

func (s *SaveObject) Future(ctx context.Context, FLOW belt.Flow) {
	FLOW.Wait(s.Present)
}

func (s *SaveObject) Present(ctx context.Context, FLOW belt.Flow) {
	s.perms = &CheckPermissions{}
	s.object = &GetObjectFromDB{hash: ""}
	FLOW.YieldAll(s.migrateToPast, s.perms, s.object)

	if !s.perms.allowedToSave {
		FLOW.YieldNone(nil, &SendReply{message: "You shall not pass!"})
		return
	}

	if s.object.exists {
		FLOW.YieldNone(nil, &SendReply{message: "Object already exists"})
		return
	}

	saved := &SaveObjectToDB{hash: ""}
	FLOW.YieldAll(s.migrateToPast, saved)

	if saved.err != nil {
		FLOW.YieldNone(nil, &SendReply{message: "Failed to save object"})
	} else {
		FLOW.YieldNone(nil, &SendReply{message: fmt.Sprintf("Object saved. ID: %d", saved.id)})
	}
}

func (s *SaveObject) Past(ctx context.Context, FLOW belt.Flow) {
	FLOW.YieldNone(nil, &SendReply{message: "Too late to save object"})
}

func (s *SaveObject) migrateToPast(ctx context.Context, FLOW belt.Flow) {
	if s.perms.allowedToSave {
		FLOW.YieldNone(nil, &SendReply{message: "You shall not pass!"})
	} else {
		FLOW.YieldNone(nil, &Redirect{toNode: "node that saves objects now"})
	}
}

// =====================================================================================================================

func bootstrapExample() {
	state := &SaveObject{}
	hand := handler.NewHandler(
		insolar.Pulse{PulseNumber: 42},
		state.Past,
		state.Present,
		state.Future,
	)

	// This should be called by an event bus.
	_, _ = hand.HandleMessage(&message.Message{
		Payload: []byte(`{"hash": "clear is better than clever"}`),
	})
}
