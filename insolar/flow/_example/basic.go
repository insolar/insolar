package _example

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/flow"
)

type CheckPermissions struct {
	Node string

	// We can group return parameters in build-in struct for clarity.
	Result struct {
		AllowedToSave bool
	}
}

func (a *CheckPermissions) Proceed(context.Context) {
	// Check for node permissions.
	a.Result.AllowedToSave = true
}

type GetObjectFromDB struct {
	DBConnection *bytes.Buffer

	Hash string

	Result struct {
		ID     int
		Exists bool
	}
}

func (a *GetObjectFromDB) Proceed(context.Context) {
	b := make([]byte, 10)
	id, err := a.DBConnection.Read(b)
	if err != nil {
		return
	}
	a.Result.Exists = true
	a.Result.ID = id
}

type SaveObjectToDB struct {
	DBConnection *bytes.Buffer

	Hash string

	Result struct {
		ID  int
		Err error
	}
}

func (a *SaveObjectToDB) Proceed(context.Context) {
	id, err := a.DBConnection.Write([]byte(a.Hash))
	if err != nil {
		a.Result.Err = err
		return
	}
	a.Result.ID = id
}

type SendReply struct {
	Message string
}

func (a *SendReply) Proceed(context.Context) {
	// Send reply over network.
}

type Redirect struct {
	ToNode string
}

func (a *Redirect) Proceed(context.Context) {
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

func (s *SaveObject) Future(ctx context.Context, FLOW flow.Flow) {
	FLOW.Procedure(s.Present, nil)
}

func (s *SaveObject) Present(ctx context.Context, FLOW flow.Flow) {
	s.perms = &CheckPermissions{Node: s.Message.Metadata["node"]}
	s.object = &GetObjectFromDB{Hash: string(s.Message.Payload)}
	FLOW.Procedure(s.migrateToPast, s.perms)
	FLOW.Procedure(s.migrateToPast, s.object)

	if !s.perms.Result.AllowedToSave {
		FLOW.Procedure(nil, &SendReply{Message: "You shall not pass!"})
		return
	}

	if s.object.Result.Exists {
		FLOW.Procedure(nil, &SendReply{Message: "Object already exists"})
		return
	}

	saved := &SaveObjectToDB{Hash: string(s.Message.Payload)}
	FLOW.Procedure(s.migrateToPast, saved)

	if saved.Result.Err != nil {
		FLOW.Procedure(nil, &SendReply{Message: "Failed to save object"})
	} else {
		FLOW.Procedure(nil, &SendReply{Message: fmt.Sprintf("Object saved. ID: %d", saved.Result.ID)})
	}
}

func (s *SaveObject) Past(ctx context.Context, FLOW flow.Flow) {
	FLOW.Procedure(nil, &SendReply{Message: "Too late to save object"})
}

func (s *SaveObject) migrateToPast(ctx context.Context, FLOW flow.Flow) {
	if s.perms.Result.AllowedToSave {
		FLOW.Procedure(nil, &SendReply{Message: "You shall not pass!"})
	} else {
		FLOW.Procedure(nil, &Redirect{ToNode: "node that saves objects now"})
	}
}
