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

package insolar

const (
	// GenesisNameRootDomain is the name of root domain contract for genesis record.
	GenesisNameRootDomain = "rootdomain"
	// GenesisNameNodeDomain is the name of node domain contract for genesis record.
	GenesisNameNodeDomain = "nodedomain"
	// GenesisNameNodeRecord is the name of node contract for genesis record.
	GenesisNameNodeRecord = "noderecord"
	// GenesisNameRootMember is the name of root member contract for genesis record.
	GenesisNameRootMember = "member"
	// GenesisNameRootWallet is the name of wallet contract for genesis record.
	GenesisNameRootWallet = "wallet"
	// GenesisNameAllowance is the name of allowance contract for genesis record.
	GenesisNameAllowance = "allowance"
)

type genesisBinary []byte

// GenesisRecord is initial chain record.
var GenesisRecord genesisBinary = []byte{0xAC}

// ID returns genesis record id.
func (r genesisBinary) ID() ID {
	return *NewID(GenesisPulse.PulseNumber, r)
}

// Ref returns genesis record reference.
func (r genesisBinary) Ref() Reference {
	id := r.ID()
	return *NewReference(id, id)
}

// DiscoveryNodeRegister carries data required for registering discovery node via genesis.
type DiscoveryNodeRegister struct {
	Role      string
	PublicKey string
}

// GenesisContractState carries data required for contract object creation via genesis.
type GenesisContractState struct {
	Name       string
	ParentName string
	Delegate   bool
	Memory     []byte
}

// GenesisContractsConfig carries data required for contract object initialization via genesis.
type GenesisContractsConfig struct {
	RootBalance   uint
	RootPublicKey string
}

// GenesisHeavyConfig carries data required for initial genesis on heavy node.
type GenesisHeavyConfig struct {
	// DiscoveryNodes is the list with discovery node info.
	DiscoveryNodes []DiscoveryNodeRegister
	// ContractsDir is the directory with contracts plugins and memory files.
	PluginsDir      string
	ContractsConfig GenesisContractsConfig
}
