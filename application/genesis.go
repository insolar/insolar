///
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
///

package application

import (
	"strconv"

	"github.com/insolar/insolar/insolar"
)

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
	// GenesisNameMigrationAdmin is the name of migration daemon contract,which is associated with MigrationDaemonMember  for genesis record.
	GenesisNameMigrationDaemon = "migrationdaemon"
	// GenesisNamePKShard is the name of public key shard contract for genesis record.
	GenesisNamePKShard = "pkshard"
	// GenesisNameMigrationShard is the name of migration address shard contract for genesis record.
	GenesisNameMigrationShard        = "migrationshard"
	GenesisNameRootMember            = "root" + GenesisNameMember
	GenesisNameRootWallet            = "root" + GenesisNameWallet
	GenesisNameRootAccount           = "root" + GenesisNameAccount
	GenesisNameMigrationAdminMember  = "migration" + GenesisNameMember
	GenesisNameMigrationAdminWallet  = "migration" + GenesisNameWallet
	GenesisNameMigrationAdminAccount = "migration" + GenesisNameAccount
	GenesisNameMigrationAdminDeposit = "migration" + GenesisNameDeposit
	GenesisNameFeeMember             = "fee" + GenesisNameMember
	GenesisNameFeeWallet             = "fee" + GenesisNameWallet
	GenesisNameFeeAccount            = "fee" + GenesisNameAccount

	GenesisAmountMigrationDaemonMembers       = 10
	GenesisAmountActiveMigrationDaemonMembers = 2

	GenesisAmountFundsMembers                 = 1
	GenesisAmountEnterpriseMembers            = 4
	GenesisAmountNetworkIncentivesMembers     = 30
	GenesisAmountApplicationIncentivesMembers = 30
	GenesisAmountFoundationMembers            = 30

	DefaultDistributionAmount = "1000000000000000000"
)

var GenesisNameMigrationDaemonMembers = func() (result [GenesisAmountMigrationDaemonMembers]string) {
	for i := 0; i < GenesisAmountMigrationDaemonMembers; i++ {
		result[i] = "migration_daemon_" + strconv.Itoa(i) + "_" + GenesisNameMember
	}
	return
}()

var GenesisNameMigrationDaemons = func() (result [GenesisAmountMigrationDaemonMembers]string) {
	for i := 0; i < GenesisAmountMigrationDaemonMembers; i++ {
		result[i] = GenesisNameMigrationDaemon + "_" + strconv.Itoa(i)
	}
	return
}()

var GenesisNameNetworkIncentivesMembers = func() (result [GenesisAmountNetworkIncentivesMembers]string) {
	for i := 0; i < GenesisAmountNetworkIncentivesMembers; i++ {
		result[i] = "network_incentives_" + strconv.Itoa(i) + "_" + GenesisNameMember
	}
	return
}()

var GenesisNameApplicationIncentivesMembers = func() (result [GenesisAmountApplicationIncentivesMembers]string) {
	for i := 0; i < GenesisAmountApplicationIncentivesMembers; i++ {
		result[i] = "application_incentives_" + strconv.Itoa(i) + "_" + GenesisNameMember
	}
	return
}()

var GenesisNameFoundationMembers = func() (result [GenesisAmountFoundationMembers]string) {
	for i := 0; i < GenesisAmountFoundationMembers; i++ {
		result[i] = "foundation_" + strconv.Itoa(i) + "_" + GenesisNameMember
	}
	return
}()

var GenesisNameFundsMembers = func() (result [GenesisAmountFundsMembers]string) {
	for i := 0; i < GenesisAmountFundsMembers; i++ {
		result[i] = "funds_" + strconv.Itoa(i) + "_" + GenesisNameMember
	}
	return
}()

var GenesisNameEnterpriseMembers = func() (result [GenesisAmountEnterpriseMembers]string) {
	for i := 0; i < GenesisAmountEnterpriseMembers; i++ {
		result[i] = "enterprise_" + strconv.Itoa(i) + "_" + GenesisNameMember
	}
	return
}()

var GenesisNameNetworkIncentivesDeposits = func() (result [GenesisAmountNetworkIncentivesMembers]string) {
	for i := 0; i < GenesisAmountNetworkIncentivesMembers; i++ {
		result[i] = "network_incentives_" + strconv.Itoa(i) + "_" + GenesisNameDeposit
	}
	return
}()

var GenesisNameApplicationIncentivesDeposits = func() (result [GenesisAmountApplicationIncentivesMembers]string) {
	for i := 0; i < GenesisAmountApplicationIncentivesMembers; i++ {
		result[i] = "application_incentives_" + strconv.Itoa(i) + "_" + GenesisNameDeposit
	}
	return
}()

var GenesisNameFoundationDeposits = func() (result [GenesisAmountFoundationMembers]string) {
	for i := 0; i < GenesisAmountFoundationMembers; i++ {
		result[i] = "foundation_" + strconv.Itoa(i) + "_" + GenesisNameDeposit
	}
	return
}()

var GenesisNameFundsDeposits = func() (result [GenesisAmountFundsMembers]string) {
	for i := 0; i < GenesisAmountFundsMembers; i++ {
		result[i] = "funds_" + strconv.Itoa(i) + "_" + GenesisNameDeposit
	}
	return
}()

var GenesisNameEnterpriseDeposits = func() (result [GenesisAmountEnterpriseMembers]string) {
	for i := 0; i < GenesisAmountEnterpriseMembers; i++ {
		result[i] = "enterprise_" + strconv.Itoa(i) + "_" + GenesisNameDeposit
	}
	return
}()

var GenesisNameNetworkIncentivesWallets = func() (result [GenesisAmountNetworkIncentivesMembers]string) {
	for i := 0; i < GenesisAmountNetworkIncentivesMembers; i++ {
		result[i] = "network_incentives_" + strconv.Itoa(i) + "_" + GenesisNameWallet
	}
	return
}()

var GenesisNameApplicationIncentivesWallets = func() (result [GenesisAmountApplicationIncentivesMembers]string) {
	for i := 0; i < GenesisAmountApplicationIncentivesMembers; i++ {
		result[i] = "application_incentives_" + strconv.Itoa(i) + "_" + GenesisNameWallet
	}
	return
}()

var GenesisNameFoundationWallets = func() (result [GenesisAmountFoundationMembers]string) {
	for i := 0; i < GenesisAmountFoundationMembers; i++ {
		result[i] = "foundation_" + strconv.Itoa(i) + "_" + GenesisNameWallet
	}
	return
}()

var GenesisNameFundsWallets = func() (result [GenesisAmountFundsMembers]string) {
	for i := 0; i < GenesisAmountFundsMembers; i++ {
		result[i] = "funds_" + strconv.Itoa(i) + "_" + GenesisNameWallet
	}
	return
}()

var GenesisNameEnterpriseWallets = func() (result [GenesisAmountEnterpriseMembers]string) {
	for i := 0; i < GenesisAmountEnterpriseMembers; i++ {
		result[i] = "enterprise_" + strconv.Itoa(i) + "_" + GenesisNameWallet
	}
	return
}()

var GenesisNameNetworkIncentivesAccounts = func() (result [GenesisAmountNetworkIncentivesMembers]string) {
	for i := 0; i < GenesisAmountNetworkIncentivesMembers; i++ {
		result[i] = "network_incentives_" + strconv.Itoa(i) + "_" + GenesisNameAccount
	}
	return
}()

var GenesisNameApplicationIncentivesAccounts = func() (result [GenesisAmountApplicationIncentivesMembers]string) {
	for i := 0; i < GenesisAmountApplicationIncentivesMembers; i++ {
		result[i] = "application_incentives_" + strconv.Itoa(i) + "_" + GenesisNameAccount
	}
	return
}()

var GenesisNameFoundationAccounts = func() (result [GenesisAmountFoundationMembers]string) {
	for i := 0; i < GenesisAmountFoundationMembers; i++ {
		result[i] = "foundation_" + strconv.Itoa(i) + "_" + GenesisNameAccount
	}
	return
}()

var GenesisNameFundsAccounts = func() (result [GenesisAmountFundsMembers]string) {
	for i := 0; i < GenesisAmountFundsMembers; i++ {
		result[i] = "funds_" + strconv.Itoa(i) + "_" + GenesisNameAccount
	}
	return
}()
var GenesisNameEnterpriseAccounts = func() (result [GenesisAmountEnterpriseMembers]string) {
	for i := 0; i < GenesisAmountEnterpriseMembers; i++ {
		result[i] = "enterprise_" + strconv.Itoa(i) + "_" + GenesisNameAccount
	}
	return
}()

type genesisBinary []byte

// GenesisRecord is initial chain record.
var GenesisRecord genesisBinary = []byte{0xAC}

// ID returns genesis record id.
func (r genesisBinary) ID() insolar.ID {
	return *insolar.NewID(insolar.GenesisPulse.PulseNumber, r)
}

// Ref returns genesis record reference.
func (r genesisBinary) Ref() insolar.Reference {
	return *insolar.NewReference(r.ID())
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
	RootBalance                     string
	Fee                             string
	MDBalance                       string
	RootPublicKey                   string
	FeePublicKey                    string
	MigrationAdminPublicKey         string
	MigrationDaemonPublicKeys       []string
	VestingPeriodInPulses           int64
	LockupPeriodInPulses            int64
	VestingStepInPulses             int64
	MigrationAddresses              [][]string
	FundsPublicKeys                 []string
	EnterprisePublicKeys            []string
	NetworkIncentivesPublicKeys     []string
	ApplicationIncentivesPublicKeys []string
	FoundationPublicKeys            []string
	PKShardCount                    int
	MAShardCount                    int
}

// GenesisHeavyConfig carries data required for initial genesis on heavy node.
type GenesisHeavyConfig struct {
	// DiscoveryNodes is the list with discovery node info.
	DiscoveryNodes  []DiscoveryNodeRegister
	ContractsConfig GenesisContractsConfig
	// Skip is flag for skipping genesis on heavy node. Useful for some test cases.
	Skip bool
}
