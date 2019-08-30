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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulse"
)

const (
	PrototypeType   = "prototype"
	PrototypeSuffix = "_proto"
)

var PredefinedPrototypes = map[string]insolar.Reference{
	insolar.GenesisNameRootDomain + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameRootDomain, 0),
	insolar.GenesisNameNodeDomain + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameNodeDomain, 0),
	insolar.GenesisNameNodeRecord + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameNodeRecord, 0),
	insolar.GenesisNameRootMember + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameMember, 0),
	insolar.GenesisNameRootWallet + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameWallet, 0),
	insolar.GenesisNameRootAccount + PrototypeSuffix:           *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameAccount, 0),
	insolar.GenesisNameCostCenter + PrototypeSuffix:            *GenerateProtoReferenceFromContractID(PrototypeType, insolar.GenesisNameCostCenter, 0),
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
	// ContractDeposit is the deposit contract reference.
	ContractDeposit = GenesisRef(insolar.GenesisNameDeposit)
	// ContractCostCenter is the cost center contract reference.
	ContractCostCenter = GenesisRef(insolar.GenesisNameCostCenter)
	// ContractFeeAccount is the commission account contract reference.
	ContractFeeAccount = GenesisRef(insolar.GenesisNameFeeAccount)

	// ContractMigrationDaemonMembers is the migration daemon members contracts references.
	ContractMigrationDaemonMembers = func() (result [insolar.GenesisAmountMigrationDaemonMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameMigrationDaemonMembers {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractMigrationDaemons is the migration daemon contracts references,which is associated with MigrationDaemonMember.
	ContractMigrationDaemons = func() (result [insolar.GenesisAmountMigrationDaemonMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameMigrationDaemons {
			result[i] = GenesisRef(name)
		}
		return
	}()

	// ContractPublicKeyShards is the public key shards contracts references.
	ContractPublicKeyShards = func() (result [insolar.GenesisAmountPublicKeyShards]insolar.Reference) {
		for i, name := range insolar.GenesisNamePublicKeyShards {
			result[i] = GenesisRef(name)
		}
		return
	}()
	// ContractMigrationAddressShards is the migration address shards contracts references.
	ContractMigrationAddressShards = func() (result [insolar.GenesisAmountMigrationAddressShards]insolar.Reference) {
		for i, name := range insolar.GenesisNameMigrationAddressShards {
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
