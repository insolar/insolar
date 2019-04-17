//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package _example

import (
	"bytes"
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar/flow"
)

type CheckPermissions struct {
	Node string

	// We can group return parameters in build-in struct for clarity.
	Result struct {
		AllowedToSave bool
	}
}

func (a *CheckPermissions) Proceed(context.Context) error {
	// Check for node permissions.
	a.Result.AllowedToSave = true
	return nil
}

type GetObjectFromDB struct {
	Hash string

	Result struct {
		ID     int
		Exists bool
	}
	// We can group dependencies in build-in struct for clarity.
	Dep struct {
		DBConnection *bytes.Buffer
	}
}

func (a *GetObjectFromDB) Proceed(context.Context) error {
	b := make([]byte, 10)
	id, err := a.Dep.DBConnection.Read(b)
	if err != nil {
		return err
	}
	a.Result.Exists = true
	a.Result.ID = id
	return nil
}

type SaveObjectToDB struct {
	Hash string

	Result struct {
		ID int
	}
	Dep struct {
		DBConnection *bytes.Buffer
	}
}

func (a *SaveObjectToDB) Proceed(context.Context) error {
	id, err := a.Dep.DBConnection.Write([]byte(a.Hash))
	if err != nil {
		return err
	}
	a.Result.ID = id
	return nil
}

type SendReply struct {
	Message string
}

func (a *SendReply) Proceed(context.Context) error {
	// Send reply over network.
	return nil
}

type Redirect struct {
	ToNode string
}

func (a *Redirect) Proceed(context.Context) error {
	// Redirect to other node.
	return nil
}

// =====================================================================================================================

// SaveObject describes handling "save object" message flow.
type SaveObject struct {
	Message map[string]string

	// Keep internal state unexported.
	perms  *CheckPermissions
	object *GetObjectFromDB

	Dep struct {
		DBConnection *bytes.Buffer
	}
}

// These functions represent message handling in different time slots.
// Each function is a Handle.

func (s *SaveObject) Future(ctx context.Context, f flow.Flow) error {
	return f.Migrate(ctx, s.Present)
}

func (s *SaveObject) Present(ctx context.Context, f flow.Flow) error {
	s.perms = &CheckPermissions{Node: s.Message["node"]}
	if err := f.Procedure(ctx, s.perms); err != nil {
		return err
	}

	s.object = &GetObjectFromDB{Hash: string(s.Message["payload"])}
	if err := f.Procedure(ctx, s.object); err != nil {
		if err != flow.ErrCancelled {
			return err
		}
		return f.Migrate(ctx, s.migrate)
	}

	if !s.perms.Result.AllowedToSave {
		return f.Procedure(nil, &SendReply{Message: "You shall not pass!"})
	}

	if s.object.Result.Exists {
		return f.Procedure(nil, &SendReply{Message: "Object already exists"})
	}

	saved := &SaveObjectToDB{Hash: string(s.Message["payload"])}
	if err := f.Procedure(ctx, saved); err != nil {
		if err != flow.ErrCancelled {
			return f.Procedure(ctx, &SendReply{Message: "Failed to save object"})
		}
		return f.Migrate(ctx, s.migrate)
	}

	return f.Procedure(nil, &SendReply{Message: fmt.Sprintf("Object saved. ID: %d", saved.Result.ID)})
}

func (s *SaveObject) Past(ctx context.Context, f flow.Flow) error {
	return f.Procedure(nil, &SendReply{Message: "Too late to save object"})
}

func (s *SaveObject) migrate(ctx context.Context, f flow.Flow) error {
	if !s.perms.Result.AllowedToSave {
		return f.Procedure(nil, &SendReply{Message: "You shall not pass!"})
	}

	return f.Procedure(nil, &Redirect{ToNode: "node that saves objects now"})
}
