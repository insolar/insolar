// Copyright 2020 Insolar Network Ltd.
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

package genesis

import (
	"github.com/insolar/insolar/insolar"
)

type binary []byte

// Record is initial chain record.
var Record binary = []byte{0xAC}

// ID returns genesis record id.
func (r binary) ID() insolar.ID {
	return *insolar.NewID(insolar.GenesisPulse.PulseNumber, r)
}

// Ref returns genesis record reference.
func (r binary) Ref() insolar.Reference {
	return *insolar.NewReference(r.ID())
}

// DiscoveryNodeRegister carries data required for registering discovery node via genesis.
type DiscoveryNodeRegister struct {
	Role      string
	PublicKey string
}

// ContractState carries data required for contract object creation via genesis.
type ContractState struct {
	Name       string
	Prototype  string
	ParentName string
	Memory     []byte
}

// HeavyConfig carries data required for initial genesis on heavy node.
type HeavyConfig struct {
	// DiscoveryNodes is the list with discovery node info.
	DiscoveryNodes []DiscoveryNodeRegister
	// Skip is flag for skipping genesis on heavy node. Useful for some test cases.
	Skip bool
}
