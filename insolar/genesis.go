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

import "strconv"

const (
	// GenesisNameRootDomain is the name of root domain contract for genesis record.
	GenesisNameRootDomain = "rootdomain"
	// GenesisNameNodeDomain is the name of node domain contract for genesis record.
	GenesisNameNodeDomain = "nodedomain"
	// GenesisNameNodeRecord is the name of node contract for genesis record.
	GenesisNameNodeRecord = "noderecord"
	// GenesisNameMember is the name of member contract for genesis record.
	GenesisNameMember = "member"
	// GenesisNameWallet is the name of wallet contract for genesis record.
	GenesisNameWallet = "wallet"
	// GenesisNameDeposit is the name of deposit contract for genesis record.
	GenesisNameDeposit = "deposit"
	// GenesisNameTariff is the name of tariff contract for genesis record.
	GenesisNameTariff = "tariff"
	// GenesisNameCostCenter is the name of cost center contract for genesis record.
	GenesisNameCostCenter = "costcenter"
)

func GetGenesisNameRootDomain() string {
	return GenesisNameRootDomain
}
func GetGenesisNameNodeDomain() string {
	return GenesisNameNodeDomain
}
func GetGenesisNameNodeRecord() string {
	return GenesisNameNodeRecord
}
func GetGenesisNameRootMember() string {
	return "root" + GenesisNameMember
}
func GetGenesisNameRootWallet() string {
	return "root" + GenesisNameWallet
}
func GetGenesisNameMigrationAdminMember() string {
	return "mdadmin" + GenesisNameMember
}
func GetGenesisNameMigrationWallet() string {
	return "md" + GenesisNameWallet
}
func GetGenesisNameMigrationDamonMembers() (result [10]string) {
	for i := 0; i < 10; i++ {
		result[i] = "migration_damon_" + strconv.Itoa(i) + GenesisNameMember
	}
	return
}
func GetGenesisNameDeposit() string {
	return GenesisNameDeposit
}
func GetGenesisNameTariff() string {
	return GenesisNameTariff
}
func GetGenesisNameCostCenter() string {
	return GenesisNameCostCenter
}

type genesisBinary []byte

// GenesisRecord is initial chain record.
var GenesisRecord genesisBinary = []byte{0xAC}

// ID returns genesis record id.
func (r genesisBinary) ID() ID {
	return *NewID(GenesisPulse.PulseNumber, r)
}

// Ref returns genesis record reference.
func (r genesisBinary) Ref() Reference {
	return *NewReference(r.ID())
}

// DiscoveryNodeRegister carries data required for registering discovery node via genesis.
type DiscoveryNodeRegister struct {
	Role      string
	PublicKey string
}

// GenesisContractState carries data required for contract object creation via genesis.
type GenesisContractState struct {
	Name       string
	Prototype  string
	ParentName string
	Delegate   bool
	Memory     []byte
}

// GenesisContractsConfig carries data required for contract object initialization via genesis.
type GenesisContractsConfig struct {
	RootBalance             string
	MDBalance               string
	RootPublicKey           string
	MigrationAdminPublicKey string
	MigrationDamonPublicKey [10]string
}

// GenesisHeavyConfig carries data required for initial genesis on heavy node.
type GenesisHeavyConfig struct {
	// DiscoveryNodes is the list with discovery node info.
	DiscoveryNodes []DiscoveryNodeRegister
	// ContractsDir is the directory with contracts plugins and memory files.
	PluginsDir      string
	ContractsConfig GenesisContractsConfig
	// Skip is flag for skipping genesis on heavy node. Useful for some test cases.
	Skip bool
}
