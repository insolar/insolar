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
