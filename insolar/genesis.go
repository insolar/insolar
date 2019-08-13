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
	// GenesisNameAccount is the name of wallet contract for genesis record.
	GenesisNameAccount = "account"
	// GenesisNameDeposit is the name of deposit contract for genesis record.
	GenesisNameDeposit = "deposit"
	// GenesisNameCostCenter is the name of cost center contract for genesis record.
	GenesisNameCostCenter = "costcenter"
	// GenesisNameMigrationAdmin is the name of migration admin contract for genesis record.
	GenesisNameMigrationAdmin = "migrationadmin"
	// GenesisNamePKShard is the name of public key shard contract for genesis record.
	GenesisNamePKShard = "pkshard"
	// GenesisNameMigrationShard is the name of migration address shard contract for genesis record.
	GenesisNameMigrationShard = "migrationshard"

	GenesisNameRootMember            = "root" + GenesisNameMember
	GenesisNameRootWallet            = "root" + GenesisNameWallet
	GenesisNameRootAccount           = "root" + GenesisNameAccount
	GenesisNameMigrationAdminMember  = "migration" + GenesisNameMember
	GenesisNameMigrationAdminWallet  = "migration" + GenesisNameWallet
	GenesisNameMigrationAdminAccount = "migration" + GenesisNameAccount
	GenesisNameFeeWallet             = "fee" + GenesisNameWallet
	GenesisNameFeeAccount            = "fee" + GenesisNameAccount

	GenesisAmountMigrationDaemonMembers       = 10
	GenesisAmountActiveMigrationDaemonMembers = 3

	GenesisAmountPublicKeyShards        = 10
	GenesisAmountMigrationAddressShards = 10
)

var GenesisNameMigrationDaemonMembers = func() (result [GenesisAmountMigrationDaemonMembers]string) {
	for i := 0; i < GenesisAmountMigrationDaemonMembers; i++ {
		result[i] = "migration_daemon_" + strconv.Itoa(i) + "_" + GenesisNameMember
	}
	return
}()

var GenesisNameMigrationAddressShards = func() (result [GenesisAmountMigrationAddressShards]string) {
	for i := 0; i < GenesisAmountMigrationAddressShards; i++ {
		result[i] = GenesisNameMigrationShard + "_" + strconv.Itoa(i)
	}
	return
}()
var GenesisNamePublicKeyShards = func() (result [GenesisAmountPublicKeyShards]string) {
	for i := 0; i < GenesisAmountPublicKeyShards; i++ {
		result[i] = GenesisNamePKShard + "_" + strconv.Itoa(i)
	}
	return
}()

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
	Memory     []byte
}

// GenesisContractsConfig carries data required for contract object initialization via genesis.
type GenesisContractsConfig struct {
	RootBalance               string
	MDBalance                 string
	RootPublicKey             string
	MigrationAdminPublicKey   string
	MigrationDaemonPublicKeys []string
}

// GenesisHeavyConfig carries data required for initial genesis on heavy node.
type GenesisHeavyConfig struct {
	// DiscoveryNodes is the list with discovery node info.
	DiscoveryNodes  []DiscoveryNodeRegister
	ContractsConfig GenesisContractsConfig
	// Skip is flag for skipping genesis on heavy node. Useful for some test cases.
	Skip bool
}
