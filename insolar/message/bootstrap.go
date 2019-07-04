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

// GenesisRequest is used for genesis records generation.
// this is fake message that never passed to messageBus
// it implements Message Interface for ability to be converted
// to Parcel and than be passed to RegisterIncomingRequest method
type GenesisRequest struct {
	// Name should be unique for each genesis record.
	Name string
}

// AllowedSenderObjectAndRole implements interface method
func (*GenesisRequest) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	panic("never use GenesisRequest as message for messageBus, see comment on type declaration")
}

// DefaultRole returns role for this event
func (*GenesisRequest) DefaultRole() insolar.DynamicRole {
	panic("never use GenesisRequest as message for messageBus, see comment on type declaration")
}

// DefaultTarget returns of target of this event.
func (gr *GenesisRequest) DefaultTarget() *insolar.Reference {
	return &insolar.Reference{}
}

// Type implementation for genesis request.
func (*GenesisRequest) Type() insolar.MessageType {
	return insolar.TypeGenesisRequest
}

// GetCaller implementation for genesis request.
func (*GenesisRequest) GetCaller() *insolar.Reference {
	panic("never use GenesisRequest as message for messageBus, see comment on type declaration")
}
