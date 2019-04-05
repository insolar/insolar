package _example

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/belt"
)

type CheckPermissions struct {
	Node string

	// We can group return parameters in build-in struct for clarity.
	Result struct {
		AllowedToSave bool
	}
}

func (a *CheckPermissions) Proceed(context.Context) bool {
	// Check for node permissions.
	a.Result.AllowedToSave = true
	return false
}

type GetObjectFromDB struct {
	DBConnection *bytes.Buffer

	Hash string

	Result struct {
		ID     int
		Exists bool
	}
}

func (a *GetObjectFromDB) Proceed(context.Context) bool {
	b := make([]byte, 10)
	id, err := a.DBConnection.Read(b)
	if err != nil {
		return false
	}
	a.Result.Exists = true
	a.Result.ID = id

	return false
}

type SaveObjectToDB struct {
	DBConnection *bytes.Buffer

	Hash string

	Result struct {
		ID  int
		Err error
	}
}

func (a *SaveObjectToDB) Proceed(context.Context) bool {
	id, err := a.DBConnection.Write([]byte(a.Hash))
	if err != nil {
		a.Result.Err = err
		return false
	}
	a.Result.ID = id
	return false
}

type SendReply struct {
	Message string
}

func (a *SendReply) Proceed(context.Context) bool {
	// Send reply over network.
	return false
}

type Redirect struct {
	ToNode string
}

func (a *Redirect) Proceed(context.Context) bool {
	// Redirect to other node.
	return false
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

func (s *SaveObject) Future(ctx context.Context, FLOW belt.Flow) {
	FLOW.Yield(s.Present, nil)
}

func (s *SaveObject) Present(ctx context.Context, FLOW belt.Flow) {
	s.perms = &CheckPermissions{Node: s.Message.Metadata["node"]}
	s.object = &GetObjectFromDB{Hash: string(s.Message.Payload)}
	FLOW.Yield(s.migrateToPast, s.perms)
	FLOW.Yield(s.migrateToPast, s.object)

	if !s.perms.Result.AllowedToSave {
		FLOW.Yield(nil, &SendReply{Message: "You shall not pass!"})
		return
	}

	if s.object.Result.Exists {
		FLOW.Yield(nil, &SendReply{Message: "Object already exists"})
		return
	}

	saved := &SaveObjectToDB{Hash: string(s.Message.Payload)}
	FLOW.Yield(s.migrateToPast, saved)

	if saved.Result.Err != nil {
		FLOW.Yield(nil, &SendReply{Message: "Failed to save object"})
	} else {
		FLOW.Yield(nil, &SendReply{Message: fmt.Sprintf("Object saved. ID: %d", saved.Result.ID)})
	}
}

func (s *SaveObject) Past(ctx context.Context, FLOW belt.Flow) {
	FLOW.Yield(nil, &SendReply{Message: "Too late to save object"})
}

func (s *SaveObject) migrateToPast(ctx context.Context, FLOW belt.Flow) {
	if s.perms.Result.AllowedToSave {
		FLOW.Yield(nil, &SendReply{Message: "You shall not pass!"})
	} else {
		FLOW.Yield(nil, &Redirect{ToNode: "node that saves objects now"})
	}
}
