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
	insolar.GenesisNameRootDomain + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameRootDomain, 0),
	insolar.GenesisNameNodeDomain + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameNodeDomain, 0),
	insolar.GenesisNameNodeRecord + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameNodeRecord, 0),
	insolar.GenesisNameRootMember + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameMember, 0),
	insolar.GenesisNameRootWallet + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameWallet, 0),
	insolar.GenesisNameRootAccount + PrototypeSuffix:           *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameAccount, 0),
	insolar.GenesisNameCostCenter + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameCostCenter, 0),
	insolar.GenesisNameFeeMember + PrototypeSuffix:             *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameMember, 0),
	insolar.GenesisNameFeeWallet + PrototypeSuffix:             *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameWallet, 0),
	insolar.GenesisNameFeeAccount + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameAccount, 0),
	insolar.GenesisNameDeposit + PrototypeSuffix:               *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameDeposit, 0),
	insolar.GenesisNameMember + PrototypeSuffix:                *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameMember, 0),
	insolar.GenesisNameMigrationAdminMember + PrototypeSuffix:  *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameMember, 0),
	insolar.GenesisNameMigrationAdmin + PrototypeSuffix:        *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameMigrationAdmin, 0),
	insolar.GenesisNameMigrationAdminWallet + PrototypeSuffix:  *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameWallet, 0),
	insolar.GenesisNameMigrationAdminAccount + PrototypeSuffix: *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameAccount, 0),
	insolar.GenesisNameWallet + PrototypeSuffix:                *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameWallet, 0),
}

func init() {
	for _, el := range insolar.GenesisNameMigrationDaemonMembers {
		PredefinedPrototypes[el+PrototypeSuffix] = *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameMember, 0)
	}

	for _, el := range insolar.GenesisNameMigrationDaemons {
		PredefinedPrototypes[el+PrototypeSuffix] = *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameMigrationDaemon, 0)
	}
}

var (
	// ContractRootDomain is the root domain contract reference.
	ContractRootDomain = GenesisRef(insolar.GenesisNameRootDomain)
	// ContractNodeDomain is the node domain contract reference.
	ContractNodeDomain = GenesisRef(insolar.GenesisNameNodeDomain)
	// ContractNodeRecord is the node contract reference.
	ContractNodeRecord = GenesisRef(insolar.GenesisNameNodeRecord)
	// ContractRootMember is the root member contract reference.
	ContractRootMember = GenesisRef(insolar.GenesisNameRootMember)
	// ContractRootWallet is the root wallet contract reference.
	ContractRootWallet = GenesisRef(insolar.GenesisNameRootWallet)
	// ContractRootAccount is the root account contract reference.
	ContractRootAccount = GenesisRef(insolar.GenesisNameRootAccount)
	// ContractMigrationAdminMember is the migration admin member contract reference.
	ContractMigrationAdminMember = GenesisRef(insolar.GenesisNameMigrationAdminMember)
	// ContractMigrationWallet is the migration wallet contract reference.
	ContractMigrationAdmin = GenesisRef(insolar.GenesisNameMigrationAdmin)
	// ContractMigrationWallet is the migration wallet contract reference.
	ContractMigrationWallet = GenesisRef(insolar.GenesisNameMigrationAdminWallet)
	// ContractMigrationAccount is the migration account contract reference.
	ContractMigrationAccount = GenesisRef(insolar.GenesisNameMigrationAdminAccount)
	// ContractMigrationDeposit is the migration deposit contract reference.
	ContractMigrationDeposit = GenesisRef(insolar.GenesisNameMigrationAdminDeposit)
	// ContractDeposit is the deposit contract reference.
	ContractDeposit = GenesisRef(insolar.GenesisNameDeposit)
	// ContractCostCenter is the cost center contract reference.
	ContractCostCenter = GenesisRef(insolar.GenesisNameCostCenter)
	// ContractFeeMember is the fee member contract reference.
	ContractFeeMember = GenesisRef(insolar.GenesisNameFeeMember)
	// ContractFeeWallet is the fee wallet contract reference.
	ContractFeeWallet = GenesisRef(insolar.GenesisNameFeeWallet)
	// ContractFeeAccount is the fee account contract reference.
	ContractFeeAccount = GenesisRef(insolar.GenesisNameFeeAccount)

	// ContractMigrationDaemonMembers is the migration daemon members contracts references.
	ContractMigrationDaemonMembers = func() (result [insolar.GenesisAmountMigrationDaemonMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameMigrationDaemonMembers {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractMigrationMap where key is migration daemon member  references and value related migration daemon contract
	ContractMigrationMap = func() (result map[insolar.Reference]insolar.Reference) {
		result = make(map[insolar.Reference]insolar.Reference)
		for i := 0; i < insolar.GenesisAmountMigrationDaemonMembers; i++ {
			result[GenesisRef(insolar.GenesisNameMigrationDaemonMembers[i])] = GenesisRef(insolar.GenesisNameMigrationDaemons[i])
		}
		return
	}()

	// ContractNetworkIncentivesMembers is the network incentives members contracts references.
	ContractNetworkIncentivesMembers = func() (result [insolar.GenesisAmountNetworkIncentivesMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameNetworkIncentivesMembers {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractApplicationIncentivesMembers is the application incentives members contracts references.
	ContractApplicationIncentivesMembers = func() (result [insolar.GenesisAmountApplicationIncentivesMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameApplicationIncentivesMembers {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFoundationMembers is the foundation members contracts references.
	ContractFoundationMembers = func() (result [insolar.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameFoundationMembers {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFundsMembers is the funds members contracts references.
	ContractFundsMembers = func() (result [insolar.GenesisAmountFundsMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameFundsMembers {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractEnterpriseMembers is the enterprise members contracts references.
	ContractEnterpriseMembers = func() (result [insolar.GenesisAmountEnterpriseMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameEnterpriseMembers {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractNetworkIncentivesWallets is the network incentives members contracts references.
	ContractNetworkIncentivesWallets = func() (result [insolar.GenesisAmountNetworkIncentivesMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameNetworkIncentivesWallets {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractApplicationIncentivesWallets is the application incentives members contracts references.
	ContractApplicationIncentivesWallets = func() (result [insolar.GenesisAmountApplicationIncentivesMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameApplicationIncentivesWallets {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFoundationWallets is the foundation members contracts references.
	ContractFoundationWallets = func() (result [insolar.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameFoundationWallets {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFundsWallets is the funds members contracts references.
	ContractFundsWallets = func() (result [insolar.GenesisAmountFundsMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameFundsWallets {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractEnterpriseWallets is the enterprise members contracts references.
	ContractEnterpriseWallets = func() (result [insolar.GenesisAmountEnterpriseMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameEnterpriseWallets {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractNetworkIncentivesDeposits is the network incentives deposits contracts references.
	ContractNetworkIncentivesDeposits = func() (result [insolar.GenesisAmountNetworkIncentivesMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameNetworkIncentivesDeposits {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractApplicationIncentivesDeposits is the application incentives deposits contracts references.
	ContractApplicationIncentivesDeposits = func() (result [insolar.GenesisAmountApplicationIncentivesMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameApplicationIncentivesDeposits {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFoundationDeposits is the foundation deposits contracts references.
	ContractFoundationDeposits = func() (result [insolar.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameFoundationDeposits {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFundsDeposits is the foundation deposits contracts references.
	ContractFundsDeposits = func() (result [insolar.GenesisAmountFundsMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameFundsDeposits {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractEnterpriseDeposits is the foundation deposits contracts references.
	ContractEnterpriseDeposits = func() (result [insolar.GenesisAmountEnterpriseMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameEnterpriseDeposits {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractNetworkIncentivesAccounts is the network incentives accounts contracts references.
	ContractNetworkIncentivesAccounts = func() (result [insolar.GenesisAmountNetworkIncentivesMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameNetworkIncentivesAccounts {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractApplicationIncentivesAccounts is the application incentives accounts contracts references.
	ContractApplicationIncentivesAccounts = func() (result [insolar.GenesisAmountApplicationIncentivesMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameApplicationIncentivesAccounts {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFoundationAccounts is the foundation accounts contracts references.
	ContractFoundationAccounts = func() (result [insolar.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameFoundationAccounts {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractFundsAccounts is the funds accounts contracts references.
	ContractFundsAccounts = func() (result [insolar.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameFundsAccounts {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractEnterpriseAccounts is the enterprise accounts contracts references.
	ContractEnterpriseAccounts = func() (result [insolar.GenesisAmountFoundationMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameEnterpriseAccounts {
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
		name := insolar.GenesisNamePKShard + strconv.Itoa(i)
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
		name := insolar.GenesisNameMigrationShard + strconv.Itoa(i)
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
