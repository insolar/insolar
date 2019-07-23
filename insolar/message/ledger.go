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

// FIXME: @andreyromancev. 21.12.18. Remove this and create 'LogicRunnerMessage' interface to get rid of 'GetCaller' in ledger.
type ledgerMessage struct {
}

// GetCaller implementation of Message interface.
func (ledgerMessage) GetCaller() *insolar.Reference {
	return nil
}

// GetDelegate retrieves object represented as provided type.
type GetDelegate struct {
	ledgerMessage
	Head   insolar.Reference
	AsType insolar.Reference
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetDelegate) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.Head, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetDelegate) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetDelegate) DefaultTarget() *insolar.Reference {
	return &m.Head
}

// Type implementation of Message interface.
func (*GetDelegate) Type() insolar.MessageType {
	return insolar.TypeGetDelegate
}

// RegisterChild amends object.
type RegisterChild struct {
	ledgerMessage
	Record []byte
	Parent insolar.Reference
	Child  insolar.Reference
	AsType *insolar.Reference // If not nil, considered as delegate.
}

// AllowedSenderObjectAndRole implements interface method
func (m *RegisterChild) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.Child, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*RegisterChild) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *RegisterChild) DefaultTarget() *insolar.Reference {
	return &m.Parent
}

// Type implementation of Message interface.
func (*RegisterChild) Type() insolar.MessageType {
	return insolar.TypeRegisterChild
}

// GetChildren retrieves a chunk of children references.
type GetChildren struct {
	ledgerMessage
	Parent    insolar.Reference
	FromChild *insolar.ID
	FromPulse *insolar.PulseNumber
	Amount    int
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetChildren) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.Parent, insolar.DynamicRoleVirtualExecutor
}

// DefaultRole returns role for this event
func (*GetChildren) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleLightExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetChildren) DefaultTarget() *insolar.Reference {
	return &m.Parent
}

// Type implementation of Message interface.
func (*GetChildren) Type() insolar.MessageType {
	return insolar.TypeGetChildren
}

// GetObjectIndex fetches objects index.
type GetObjectIndex struct {
	ledgerMessage

	Object insolar.Reference
}

// AllowedSenderObjectAndRole implements interface method
func (m *GetObjectIndex) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &m.Object, insolar.DynamicRoleLightExecutor
}

// DefaultRole returns role for this event
func (*GetObjectIndex) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleHeavyExecutor
}

// DefaultTarget returns of target of this event.
func (m *GetObjectIndex) DefaultTarget() *insolar.Reference {
	return &m.Object
}

// Type implementation of Message interface.
func (*GetObjectIndex) Type() insolar.MessageType {
	return insolar.TypeGetObjectIndex
}
