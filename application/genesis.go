// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package application

import (
	"strconv"
)

const (
	// GenesisNameRootDomain is the name of root domain contract for genesis record.
	GenesisNameRootDomain = "rootdomain"
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

	GenesisAmountEnterpriseMembers            = 1
	GenesisAmountNetworkIncentivesMembers     = 20
	GenesisAmountApplicationIncentivesMembers = 20
	GenesisAmountFoundationMembers            = 20

	EnterpriseDistributionAmount        = "2000000000000000000"
	FoundationDistributionAmount        = "50000000000000000"
	AppIncentivesDistributionAmount     = "50000000000000000"
	NetworkIncentivesDistributionAmount = "50000000000000000"
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

var GenesisNameEnterpriseAccounts = func() (result [GenesisAmountEnterpriseMembers]string) {
	for i := 0; i < GenesisAmountEnterpriseMembers; i++ {
		result[i] = "enterprise_" + strconv.Itoa(i) + "_" + GenesisNameAccount
	}
	return
}()

// GenesisContractsConfig carries data required for contract object initialization via genesis.
type GenesisContractsConfig struct {
	// RootBalance is a balance of Root Member.
	RootBalance string
	// MDBalance is a balance of Migration Daemon.
	MDBalance string
	// RootPublicKey is public key of Root Member.
	RootPublicKey string
	// FeePublicKey is public key of Fee Member.
	FeePublicKey string
	// MigrationAdminPublicKey is public key of Migration Admin.
	MigrationAdminPublicKey string
	// MigrationDaemonPublicKeys is a public keys array of Migration Daemon members.
	MigrationDaemonPublicKeys []string
	// VestingPeriodInPulses is a vesting period measured with pulses.
	VestingPeriodInPulses int64
	// LockupPeriodInPulses is a lockup period before vesting measured with pulses.
	LockupPeriodInPulses int64
	// VestingStepInPulses is a one vesting step measured with pulses.
	VestingStepInPulses int64
	// MigrationAddresses are migration addresses array per shard where index array is a shard index.
	MigrationAddresses [][]string
	// FundsPublicKeys is a public keys array of Funds members.
	FundsPublicKeys []string
	// EnterprisePublicKeys is a public keys array of Enterprise members.
	EnterprisePublicKeys []string
	// NetworkIncentivesPublicKeys is a public keys array of Network Incentives members.
	NetworkIncentivesPublicKeys []string
	// ApplicationIncentivesPublicKeys is a public keys array of Application Incentives members.
	ApplicationIncentivesPublicKeys []string
	// FoundationPublicKeys is a public keys array of Foundation members.
	FoundationPublicKeys []string
	// PKShardCount is a primary keys shards count.
	PKShardCount int
	// MAShardCount is a migration addresses shards count.
	MAShardCount int
}
