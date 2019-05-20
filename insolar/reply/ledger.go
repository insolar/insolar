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

package reply

import (
	"github.com/insolar/insolar/insolar"
)

// Code is code from storage.
type Code struct {
	Code        []byte
	MachineType insolar.MachineType
}

// Type implementation of Reply interface.
func (e *Code) Type() insolar.ReplyType {
	return TypeCode
}

// Object is object from storage.
type Object struct {
	Head         insolar.Reference
	State        insolar.ID
	Prototype    *insolar.Reference
	IsPrototype  bool
	ChildPointer *insolar.ID
	Memory       []byte
	Parent       insolar.Reference
}

// Type implementation of Reply interface.
func (e *Object) Type() insolar.ReplyType {
	return TypeObject
}

// Delegate is delegate reference from storage.
type Delegate struct {
	Head insolar.Reference
}

// Type implementation of Reply interface.
func (e *Delegate) Type() insolar.ReplyType {
	return TypeDelegate
}

// ID is common reaction for methods returning id to lifeline states.
type ID struct {
	ID insolar.ID
}

// Type implementation of Reply interface.
func (e *ID) Type() insolar.ReplyType {
	return TypeID
}

// Children is common reaction for methods returning id to lifeline states.
type Children struct {
	Refs     []insolar.Reference
	NextFrom *insolar.ID
}

// Type implementation of Reply interface.
func (e *Children) Type() insolar.ReplyType {
	return TypeChildren
}

// ObjectIndex contains serialized object index. It can be stored in DB without processing.
type ObjectIndex struct {
	Index []byte
}

// Type implementation of Reply interface.
func (e *ObjectIndex) Type() insolar.ReplyType {
	return TypeObjectIndex
}

// JetMiss is returned for miscalculated jets due to incomplete jet tree.
type JetMiss struct {
	JetID insolar.ID
	Pulse insolar.PulseNumber
}

// Type implementation of Reply interface.
func (e *JetMiss) Type() insolar.ReplyType {
	return TypeJetMiss
}

// HasPendingRequests contains unclosed requests for an object.
type HasPendingRequests struct {
	Has bool
}

// Type implementation of Reply interface.
func (e *HasPendingRequests) Type() insolar.ReplyType {
	return TypePendingRequests
}

// HasPendingRequests contains unclosed requests for an object.
type PendingRequest struct {
	ID  insolar.ID
	Err Error
}

// Type implementation of Reply interface.
func (PendingRequest) Type() insolar.ReplyType {
	return TypePendingRequestID
}

// Jet contains jet.
type Jet struct {
	ID     insolar.ID
	Actual bool
}

// Type implementation of Reply interface.
func (r *Jet) Type() insolar.ReplyType {
	return TypeJet
}

// Request contains jet.
type Request struct {
	ID     insolar.ID
	Record []byte
}

// Type implementation of Reply interface.
func (r *Request) Type() insolar.ReplyType {
	return TypeRequest
}
