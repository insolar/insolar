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

// HeavyPayload carries Key/Value records and pulse number
// that replicates to Heavy Material node.
type HeavyPayload struct {
	JetID    insolar.JetID
	PulseNum insolar.PulseNumber
	Indexes  map[insolar.ID][]byte
	Drop     []byte
	Blobs    [][]byte
	Records  [][]byte
}

// AllowedSenderObjectAndRole implements interface method
func (*HeavyPayload) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return nil, 0
}

// DefaultRole returns role for this event
func (*HeavyPayload) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleHeavyExecutor
}

// DefaultTarget returns of target of this event.
func (hp *HeavyPayload) DefaultTarget() *insolar.Reference {
	return &insolar.Reference{}
}

// GetCaller implementation of Message interface.
func (HeavyPayload) GetCaller() *insolar.Reference {
	return nil
}

// Type implementation of Message interface.
func (hp *HeavyPayload) Type() insolar.MessageType {
	return insolar.TypeHeavyPayload
}
