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

package message

import (
	"github.com/insolar/insolar/insolar"
)

type NodeSignPayloadInt interface {
	insolar.Message
	GetNodeRef() *insolar.RecordRef
}

type NodeSignPayload struct {
	NodeRef *insolar.RecordRef
}

// AllowedSenderObjectAndRole implements interface method
func (nsp *NodeSignPayload) AllowedSenderObjectAndRole() (*insolar.RecordRef, insolar.DynamicRole) {
	return nil, insolar.DynamicRoleUndefined
}

// DefaultRole returns role for this event
func (nsp *NodeSignPayload) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleUndefined
}

// DefaultTarget returns of target of this event.
func (nsp *NodeSignPayload) DefaultTarget() *insolar.RecordRef {
	return nsp.NodeRef
}

// GetCaller implementation of Message interface.
func (NodeSignPayload) GetCaller() *insolar.RecordRef {
	return nil
}

// Type implementation of Message interface.
func (nsp *NodeSignPayload) Type() insolar.MessageType {
	return insolar.TypeNodeSignRequest
}

func (nsp *NodeSignPayload) GetNodeRef() *insolar.RecordRef {
	return nsp.NodeRef
}
