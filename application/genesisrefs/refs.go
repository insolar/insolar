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

package genesisrefs

import (
	"fmt"
	"strconv"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulse"
)

const (
	PrototypeType    = "prototype"
	PrototypeSuffix  = "_proto"
	FundsDepositName = "genesis_deposit"
)

var PredefinedPrototypes = map[string]insolar.Reference{
	application.GenesisNameRootDomain + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameRootDomain, 0),
	application.GenesisNameNodeDomain + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameNodeDomain, 0),
	application.GenesisNameNodeRecord + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameNodeRecord, 0),
	application.GenesisNameRootMember + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameMember, 0),
	application.GenesisNameRootWallet + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameWallet, 0),
	application.GenesisNameRootAccount + PrototypeSuffix:           *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameAccount, 0),
	application.GenesisNameCostCenter + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameCostCenter, 0),
	application.GenesisNameFeeMember + PrototypeSuffix:             *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameMember, 0),
	application.GenesisNameFeeWallet + PrototypeSuffix:             *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameWallet, 0),
	application.GenesisNameFeeAccount + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameAccount, 0),
	application.GenesisNameDeposit + PrototypeSuffix:               *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameDeposit, 0),
	application.GenesisNameMember + PrototypeSuffix:                *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameMember, 0),
	application.GenesisNameMigrationAdminMember + PrototypeSuffix:  *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameMember, 0),
	application.GenesisNameMigrationAdmin + PrototypeSuffix:        *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameMigrationAdmin, 0),
	application.GenesisNameMigrationAdminWallet + PrototypeSuffix:  *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameWallet, 0),
	application.GenesisNameMigrationAdminAccount + PrototypeSuffix: *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameAccount, 0),
	application.GenesisNameMigrationAdminDeposit + PrototypeSuffix: *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameDeposit, 0),
	application.GenesisNameWallet + PrototypeSuffix:                *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameWallet, 0),
}

func init() {
	for _, el := range application.GenesisNameMigrationDaemonMembers {
		PredefinedPrototypes[el+PrototypeSuffix] = *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameMember, 0)
	}

	for _, el := range application.GenesisNameMigrationDaemons {
		PredefinedPrototypes[el+PrototypeSuffix] = *GenerateProtoReferenceFromContractID(PrototypeType, application.GenesisNameMigrationDaemon, 0)
	}
}

var (
	// ContractRootDomain is the root domain contract reference.
	ContractRootDomain = GenesisRef(application.GenesisNameRootDomain)
	// ContractNodeDomain is the node domain contract reference.
	ContractNodeDomain = GenesisRef(application.GenesisNameNodeDomain)
	// ContractNodeRecord is the node contract reference.
	ContractNodeRecord = GenesisRef(application.GenesisNameNodeRecord)
	// ContractRootMember is the root member contract reference.
	ContractRootMember = GenesisRef(application.GenesisNameRootMember)
	// ContractRootWallet is the root wallet contract reference.
	ContractRootWallet = GenesisRef(application.GenesisNameRootWallet)
	// ContractRootAccount is the root account contract reference.
	ContractRootAccount = GenesisRef(application.GenesisNameRootAccount)
	// ContractMigrationAdminMember is the migration admin member contract reference.
	ContractMigrationAdminMember = GenesisRef(application.GenesisNameMigrationAdminMember)
	// ContractMigrationAdmin is the migration wallet contract reference.
	ContractMigrationAdmin = GenesisRef(application.GenesisNameMigrationAdmin)
	// ContractMigrationWallet is the migration wallet contract reference.
	ContractMigrationWallet = GenesisRef(application.GenesisNameMigrationAdminWallet)
	// ContractMigrationAccount is the migration account contract reference.
	ContractMigrationAccount = GenesisRef(application.GenesisNameMigrationAdminAccount)
	// ContractMigrationDeposit is the migration deposit contract reference.
	ContractMigrationDeposit = GenesisRef(application.GenesisNameMigrationAdminDeposit)
	// ContractDeposit is the deposit contract reference.
	ContractDeposit = GenesisRef(application.GenesisNameDeposit)
	// ContractCostCenter is the cost center contract reference.
	ContractCostCenter = GenesisRef(application.GenesisNameCostCenter)
	// ContractFeeMember is the fee member contract reference.
	ContractFeeMember = GenesisRef(application.GenesisNameFeeMember)
	// ContractFeeWallet is the fee wallet contract reference.
	ContractFeeWallet = GenesisRef(application.GenesisNameFeeWallet)
	// ContractFeeAccount is the fee account contract reference.
	ContractFeeAccount = GenesisRef(application.GenesisNameFeeAccount)

	// ContractMigrationDaemonMembers is the migration daemon members contracts references.
	ContractMigrationDaemonMembers = func() (result [application.GenesisAmountMigrationDaemonMembers]insolar.Reference) {
		for i, name := range application.GenesisNameMigrationDaemonMembers {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractMigrationMap where key is migration daemon member  references and value related migration daemon contract
	ContractMigrationMap = func() (result map[insolar.Reference]insolar.Reference) {
		result = make(map[insolar.Reference]insolar.Reference)
		for i := 0; i < application.GenesisAmountMigrationDaemonMembers; i++ {
			result[GenesisRef(application.GenesisNameMigrationDaemonMembers[i])] = GenesisRef(application.GenesisNameMigrationDaemons[i])
		}
		return
	}()

	// ContractNetworkIncentivesMembers is the network incentives members contracts references.
	ContractNetworkIncentivesMembers = func() (result [application.GenesisAmountNetworkIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameNetworkIncentivesMembers {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractApplicationIncentivesMembers is the application incentives members contracts references.
	ContractApplicationIncentivesMembers = func() (result [application.GenesisAmountApplicationIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameApplicationIncentivesMembers {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFoundationMembers is the foundation members contracts references.
	ContractFoundationMembers = func() (result [application.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range application.GenesisNameFoundationMembers {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFundsMembers is the funds members contracts references.
	ContractFundsMembers = func() (result [application.GenesisAmountFundsMembers]insolar.Reference) {
		for i, name := range application.GenesisNameFundsMembers {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractEnterpriseMembers is the enterprise members contracts references.
	ContractEnterpriseMembers = func() (result [application.GenesisAmountEnterpriseMembers]insolar.Reference) {
		for i, name := range application.GenesisNameEnterpriseMembers {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractNetworkIncentivesWallets is the network incentives members contracts references.
	ContractNetworkIncentivesWallets = func() (result [application.GenesisAmountNetworkIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameNetworkIncentivesWallets {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractApplicationIncentivesWallets is the application incentives members contracts references.
	ContractApplicationIncentivesWallets = func() (result [application.GenesisAmountApplicationIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameApplicationIncentivesWallets {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFoundationWallets is the foundation members contracts references.
	ContractFoundationWallets = func() (result [application.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range application.GenesisNameFoundationWallets {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFundsWallets is the funds members contracts references.
	ContractFundsWallets = func() (result [application.GenesisAmountFundsMembers]insolar.Reference) {
		for i, name := range application.GenesisNameFundsWallets {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractEnterpriseWallets is the enterprise members contracts references.
	ContractEnterpriseWallets = func() (result [application.GenesisAmountEnterpriseMembers]insolar.Reference) {
		for i, name := range application.GenesisNameEnterpriseWallets {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractNetworkIncentivesDeposits is the network incentives deposits contracts references.
	ContractNetworkIncentivesDeposits = func() (result [application.GenesisAmountNetworkIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameNetworkIncentivesDeposits {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractApplicationIncentivesDeposits is the application incentives deposits contracts references.
	ContractApplicationIncentivesDeposits = func() (result [application.GenesisAmountApplicationIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameApplicationIncentivesDeposits {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFoundationDeposits is the foundation deposits contracts references.
	ContractFoundationDeposits = func() (result [application.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range application.GenesisNameFoundationDeposits {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFundsDeposits is the foundation deposits contracts references.
	ContractFundsDeposits = func() (result [application.GenesisAmountFundsMembers]insolar.Reference) {
		for i, name := range application.GenesisNameFundsDeposits {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractEnterpriseDeposits is the foundation deposits contracts references.
	ContractEnterpriseDeposits = func() (result [application.GenesisAmountEnterpriseMembers]insolar.Reference) {
		for i, name := range application.GenesisNameEnterpriseDeposits {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractNetworkIncentivesAccounts is the network incentives accounts contracts references.
	ContractNetworkIncentivesAccounts = func() (result [application.GenesisAmountNetworkIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameNetworkIncentivesAccounts {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractApplicationIncentivesAccounts is the application incentives accounts contracts references.
	ContractApplicationIncentivesAccounts = func() (result [application.GenesisAmountApplicationIncentivesMembers]insolar.Reference) {
		for i, name := range application.GenesisNameApplicationIncentivesAccounts {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFoundationAccounts is the foundation accounts contracts references.
	ContractFoundationAccounts = func() (result [application.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range application.GenesisNameFoundationAccounts {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFundsAccounts is the funds accounts contracts references.
	ContractFundsAccounts = func() (result [application.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range application.GenesisNameFundsAccounts {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractEnterpriseAccounts is the enterprise accounts contracts references.
	ContractEnterpriseAccounts = func() (result [application.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range application.GenesisNameEnterpriseAccounts {
			result[i] = GenesisRef(name)
		}
		return
	}()
)

// Generate reference from hash code.
func GenerateProtoReferenceFromCode(pulse insolar.PulseNumber, code []byte) *insolar.Reference {
	hasher := platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher()
	codeHash := hasher.Hash(code)
	id := insolar.NewID(pulse, codeHash)
	return insolar.NewReference(*id)
}

// Generate prototype reference from contract id.
func GenerateProtoReferenceFromContractID(typeContractID string, name string, version int) *insolar.Reference {
	contractID := fmt.Sprintf("%s::%s::v%02d", typeContractID, name, version)
	return GenerateProtoReferenceFromCode(pulse.BuiltinContract, []byte(contractID))
}

// Generate contract reference from contract id.
func GenerateCodeReferenceFromContractID(typeContractID string, name string, version int) *insolar.Reference {
	contractID := fmt.Sprintf("%s::%s::v%02d", typeContractID, name, version)
	hasher := platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher()
	codeHash := hasher.Hash([]byte(contractID))
	id := insolar.NewID(pulse.BuiltinContract, codeHash)
	return insolar.NewRecordReference(*id)
}

// GenesisRef returns reference to any genesis records.
func GenesisRef(name string) insolar.Reference {
	if ref, ok := PredefinedPrototypes[name]; ok {
		return ref
	}
	pcs := platformpolicy.NewPlatformCryptographyScheme()
	req := record.IncomingRequest{
		CallType: record.CTGenesis,
		Method:   name,
	}
	virtRec := record.Wrap(&req)
	hash := record.HashVirtual(pcs.ReferenceHasher(), virtRec)
	id := insolar.NewID(pulse.MinTimePulse, hash)
	return *insolar.NewReference(*id)
}

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
		result[i] = GenesisRef(name)
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
		result[i] = GenesisRef(name)
	}
	return result
}
