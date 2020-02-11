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

package genesisrefs

import (
	"strconv"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/applicationbase/genesisrefs"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/preprocessor"
)

const (
	FundsDepositName = "genesis_deposit"
)

var applicationPrototypes = map[string]insolar.Reference{
	application.GenesisNameRootDomain + genesisrefs.PrototypeSuffix:            *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameRootDomain, 0),
	application.GenesisNameRootMember + genesisrefs.PrototypeSuffix:            *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameMember, 0),
	application.GenesisNameRootWallet + genesisrefs.PrototypeSuffix:            *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameWallet, 0),
	application.GenesisNameRootAccount + genesisrefs.PrototypeSuffix:           *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameAccount, 0),
	application.GenesisNameCostCenter + genesisrefs.PrototypeSuffix:            *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameCostCenter, 0),
	application.GenesisNameFeeMember + genesisrefs.PrototypeSuffix:             *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameMember, 0),
	application.GenesisNameFeeWallet + genesisrefs.PrototypeSuffix:             *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameWallet, 0),
	application.GenesisNameFeeAccount + genesisrefs.PrototypeSuffix:            *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameAccount, 0),
	application.GenesisNameDeposit + genesisrefs.PrototypeSuffix:               *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameDeposit, 0),
	application.GenesisNameMember + genesisrefs.PrototypeSuffix:                *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameMember, 0),
	application.GenesisNameMigrationAdminMember + genesisrefs.PrototypeSuffix:  *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameMember, 0),
	application.GenesisNameMigrationAdmin + genesisrefs.PrototypeSuffix:        *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameMigrationAdmin, 0),
	application.GenesisNameMigrationAdminWallet + genesisrefs.PrototypeSuffix:  *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameWallet, 0),
	application.GenesisNameMigrationAdminAccount + genesisrefs.PrototypeSuffix: *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameAccount, 0),
	application.GenesisNameMigrationAdminDeposit + genesisrefs.PrototypeSuffix: *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameDeposit, 0),
	application.GenesisNameWallet + genesisrefs.PrototypeSuffix:                *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameWallet, 0),
}

func init() {
	for i, val := range applicationPrototypes {
		genesisrefs.PredefinedPrototypes[i] = val
	}

	for _, el := range application.GenesisNameMigrationDaemonMembers {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameMember, 0)
	}

	for _, el := range application.GenesisNameMigrationDaemons {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameMigrationDaemon, 0)
	}

	// Incentives Application
	for _, el := range application.GenesisNameApplicationIncentivesMembers {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameMember, 0)
	}
	for _, el := range application.GenesisNameApplicationIncentivesWallets {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameWallet, 0)
	}
	for _, el := range application.GenesisNameApplicationIncentivesAccounts {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameAccount, 0)
	}
	for _, el := range application.GenesisNameApplicationIncentivesDeposits {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameDeposit, 0)
	}

	// Network Incentives
	for _, el := range application.GenesisNameNetworkIncentivesMembers {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameMember, 0)
	}
	for _, el := range application.GenesisNameNetworkIncentivesWallets {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameWallet, 0)
	}
	for _, el := range application.GenesisNameNetworkIncentivesAccounts {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameAccount, 0)
	}
	for _, el := range application.GenesisNameNetworkIncentivesDeposits {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameDeposit, 0)
	}

	// Foundation
	for _, el := range application.GenesisNameFoundationMembers {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameMember, 0)
	}
	for _, el := range application.GenesisNameFoundationWallets {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameWallet, 0)
	}
	for _, el := range application.GenesisNameFoundationAccounts {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameAccount, 0)
	}
	for _, el := range application.GenesisNameFoundationDeposits {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameDeposit, 0)
	}

	// Enterprise
	for _, el := range application.GenesisNameEnterpriseMembers {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameMember, 0)
	}
	for _, el := range application.GenesisNameEnterpriseWallets {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameWallet, 0)
	}
	for _, el := range application.GenesisNameEnterpriseAccounts {
		genesisrefs.PredefinedPrototypes[el+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(genesisrefs.PrototypeType, application.GenesisNameAccount, 0)
	}
}

var (
	// ContractRootDomain is the root domain contract reference.
	ContractRootDomain = genesisrefs.GenesisRef(application.GenesisNameRootDomain)
	// ContractRootMember is the root member contract reference.
	ContractRootMember = genesisrefs.GenesisRef(application.GenesisNameRootMember)
	// ContractRootWallet is the root wallet contract reference.
	ContractRootWallet = genesisrefs.GenesisRef(application.GenesisNameRootWallet)
	// ContractRootAccount is the root account contract reference.
	ContractRootAccount = genesisrefs.GenesisRef(application.GenesisNameRootAccount)
	// ContractMigrationAdminMember is the migration admin member contract reference.
	ContractMigrationAdminMember = genesisrefs.GenesisRef(application.GenesisNameMigrationAdminMember)
	// ContractMigrationAdmin is the migration wallet contract reference.
	ContractMigrationAdmin = genesisrefs.GenesisRef(application.GenesisNameMigrationAdmin)
	// ContractMigrationWallet is the migration wallet contract reference.
	ContractMigrationWallet = genesisrefs.GenesisRef(application.GenesisNameMigrationAdminWallet)
	// ContractMigrationAccount is the migration account contract reference.
	ContractMigrationAccount = genesisrefs.GenesisRef(application.GenesisNameMigrationAdminAccount)
	// ContractMigrationDeposit is the migration deposit contract reference.
	ContractMigrationDeposit = genesisrefs.GenesisRef(application.GenesisNameMigrationAdminDeposit)
	// ContractDeposit is the deposit contract reference.
	ContractDeposit = genesisrefs.GenesisRef(application.GenesisNameDeposit)
	// ContractCostCenter is the cost center contract reference.
	ContractCostCenter = genesisrefs.GenesisRef(application.GenesisNameCostCenter)
	// ContractFeeMember is the fee member contract reference.
	ContractFeeMember = genesisrefs.GenesisRef(application.GenesisNameFeeMember)
	// ContractFeeWallet is the fee wallet contract reference.
	ContractFeeWallet = genesisrefs.GenesisRef(application.GenesisNameFeeWallet)
	// ContractFeeAccount is the fee account contract reference.
	ContractFeeAccount = genesisrefs.GenesisRef(application.GenesisNameFeeAccount)

	// ContractMigrationDaemonMembers is the migration daemon members contracts references.
	ContractMigrationDaemonMembers = func() (result [application.GenesisAmountMigrationDaemonMembers]insolar.Reference) {
		for i, name := range application.GenesisNameMigrationDaemonMembers {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractMigrationMap where key is migration daemon member  references and value related migration daemon contract
	ContractMigrationMap = func() (result map[insolar.Reference]insolar.Reference) {
		result = make(map[insolar.Reference]insolar.Reference)
		for i := 0; i < application.GenesisAmountMigrationDaemonMembers; i++ {
			result[genesisrefs.GenesisRef(application.GenesisNameMigrationDaemonMembers[i])] = genesisrefs.GenesisRef(application.GenesisNameMigrationDaemons[i])
		}
		return
	}()

	// ContractNetworkIncentivesMembers is the network incentives members contracts references.
	ContractNetworkIncentivesMembers = func() (result [application.GenesisAmountNetworkIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameNetworkIncentivesMembers {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractApplicationIncentivesMembers is the application incentives members contracts references.
	ContractApplicationIncentivesMembers = func() (result [application.GenesisAmountApplicationIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameApplicationIncentivesMembers {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractFoundationMembers is the foundation members contracts references.
	ContractFoundationMembers = func() (result [application.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range application.GenesisNameFoundationMembers {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractEnterpriseMembers is the enterprise members contracts references.
	ContractEnterpriseMembers = func() (result [application.GenesisAmountEnterpriseMembers]insolar.Reference) {
		for i, name := range application.GenesisNameEnterpriseMembers {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractNetworkIncentivesWallets is the network incentives members contracts references.
	ContractNetworkIncentivesWallets = func() (result [application.GenesisAmountNetworkIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameNetworkIncentivesWallets {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractApplicationIncentivesWallets is the application incentives members contracts references.
	ContractApplicationIncentivesWallets = func() (result [application.GenesisAmountApplicationIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameApplicationIncentivesWallets {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractFoundationWallets is the foundation members contracts references.
	ContractFoundationWallets = func() (result [application.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range application.GenesisNameFoundationWallets {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractEnterpriseWallets is the enterprise members contracts references.
	ContractEnterpriseWallets = func() (result [application.GenesisAmountEnterpriseMembers]insolar.Reference) {
		for i, name := range application.GenesisNameEnterpriseWallets {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractNetworkIncentivesDeposits is the network incentives deposits contracts references.
	ContractNetworkIncentivesDeposits = func() (result [application.GenesisAmountNetworkIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameNetworkIncentivesDeposits {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractApplicationIncentivesDeposits is the application incentives deposits contracts references.
	ContractApplicationIncentivesDeposits = func() (result [application.GenesisAmountApplicationIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameApplicationIncentivesDeposits {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractFoundationDeposits is the foundation deposits contracts references.
	ContractFoundationDeposits = func() (result [application.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range application.GenesisNameFoundationDeposits {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractNetworkIncentivesAccounts is the network incentives accounts contracts references.
	ContractNetworkIncentivesAccounts = func() (result [application.GenesisAmountNetworkIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameNetworkIncentivesAccounts {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractApplicationIncentivesAccounts is the application incentives accounts contracts references.
	ContractApplicationIncentivesAccounts = func() (result [application.GenesisAmountApplicationIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameApplicationIncentivesAccounts {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractFoundationAccounts is the foundation accounts contracts references.
	ContractFoundationAccounts = func() (result [application.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range application.GenesisNameFoundationAccounts {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()

	// ContractEnterpriseAccounts is the enterprise accounts contracts references.
	ContractEnterpriseAccounts = func() (result [application.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range application.GenesisNameEnterpriseAccounts {
			result[i] = genesisrefs.GenesisRef(name)
		}
		return
	}()
)

// ContractPublicKeyNameShards is the public key shards contracts names.
func ContractPublicKeyNameShards(pkShardCount int) []string {
	result := make([]string, pkShardCount)
	for i := 0; i < pkShardCount; i++ {
		name := application.GenesisNamePKShard + strconv.Itoa(i)
		result[i] = name
	}
	return result
}

// ContractPublicKeyShards is the public key shards contracts references.
func ContractPublicKeyShards(pkShardCount int) []insolar.Reference {
	result := make([]insolar.Reference, pkShardCount)
	for i, name := range ContractPublicKeyNameShards(pkShardCount) {
		result[i] = genesisrefs.GenesisRef(name)
	}
	return result
}

// ContractMigrationAddressNameShards is the migration address shards contracts names.
func ContractMigrationAddressNameShards(maShardCount int) []string {
	result := make([]string, maShardCount)
	for i := 0; i < maShardCount; i++ {
		name := application.GenesisNameMigrationShard + strconv.Itoa(i)
		result[i] = name
	}
	return result
}

// ContractMigrationAddressShards is the migration address shards contracts references.
func ContractMigrationAddressShards(maShardCount int) []insolar.Reference {
	result := make([]insolar.Reference, maShardCount)
	for i, name := range ContractMigrationAddressNameShards(maShardCount) {
		result[i] = genesisrefs.GenesisRef(name)
	}
	return result
}

func ContractPublicKeyShardRefs(pkShardCount int) {
	for _, name := range ContractPublicKeyNameShards(pkShardCount) {
		genesisrefs.PredefinedPrototypes[name+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(preprocessor.PrototypeType, application.GenesisNamePKShard, 0)
	}
}

func ContractMigrationAddressShardRefs(maShardCount int) {
	for _, name := range ContractMigrationAddressNameShards(maShardCount) {
		genesisrefs.PredefinedPrototypes[name+genesisrefs.PrototypeSuffix] = *genesisrefs.GenerateProtoReferenceFromContractID(preprocessor.PrototypeType, application.GenesisNameMigrationShard, 0)
	}
}
