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
)

const (
	GenesisPrototypeSuffix = "_proto"
)

var PredefinedPrototypes = map[string]insolar.Reference{
	insolar.GenesisNameRootDomain + GenesisPrototypeSuffix: *GenerateReference("prototype", insolar.GenesisNameRootDomain, 0),
	insolar.GenesisNameNodeDomain + GenesisPrototypeSuffix: *GenerateReference("prototype", insolar.GenesisNameNodeDomain, 0),
	insolar.GenesisNameNodeRecord + GenesisPrototypeSuffix: *GenerateReference("prototype", insolar.GenesisNameNodeRecord, 0),
	insolar.GenesisNameRootMember + GenesisPrototypeSuffix: *GenerateReference("prototype", insolar.GenesisNameMember, 0),
	insolar.GenesisNameRootWallet + GenesisPrototypeSuffix: *GenerateReference("prototype", insolar.GenesisNameWallet, 0),
	insolar.GenesisNameCostCenter + GenesisPrototypeSuffix: *GenerateReference("prototype", insolar.GenesisNameCostCenter, 0),
	insolar.GenesisNameFeeWallet + GenesisPrototypeSuffix:  *GenerateReference("prototype", insolar.GenesisNameWallet, 0),
	insolar.GenesisNameDeposit + GenesisPrototypeSuffix:    *GenerateReference("prototype", insolar.GenesisNameDeposit, 0),

	insolar.GenesisNameMember + GenesisPrototypeSuffix:               *GenerateReference("prototype", insolar.GenesisNameMember, 0),
	insolar.GenesisNameMigrationAdminMember + GenesisPrototypeSuffix: *GenerateReference("prototype", insolar.GenesisNameMember, 0),
	insolar.GenesisNameMigrationWallet + GenesisPrototypeSuffix:      *GenerateReference("prototype", insolar.GenesisNameWallet, 0),
	insolar.GenesisNameStandardTariff + GenesisPrototypeSuffix:       *GenerateReference("prototype", insolar.GenesisNameTariff, 0),
	insolar.GenesisNameTariff + GenesisPrototypeSuffix:               *GenerateReference("prototype", insolar.GenesisNameTariff, 0),
	insolar.GenesisNameWallet + GenesisPrototypeSuffix:               *GenerateReference("prototype", insolar.GenesisNameWallet, 0),
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
	// ContractMigrationAdminMember is the migration admin member contract reference.
	ContractMigrationAdminMember = GenesisRef(insolar.GenesisNameMigrationAdminMember)
	// ContractMigrationWallet is the migration wallet contract reference.
	ContractMigrationWallet = GenesisRef(insolar.GenesisNameMigrationWallet)
	// ContractDeposit is the deposit contract reference.
	ContractDeposit = GenesisRef(insolar.GenesisNameDeposit)
	// ContractStandardTariff is the tariff contract reference.
	ContractStandardTariff = GenesisRef(insolar.GenesisNameStandardTariff)
	// ContractCostCenter is the cost center contract reference.
	ContractCostCenter = GenesisRef(insolar.GenesisNameCostCenter)
	// ContractFeeWallet is the commission wallet contract reference.
	ContractFeeWallet = GenesisRef(insolar.GenesisNameFeeWallet)

	// ContractMigrationDaemonMembers is the migration daemon members contracts references.
	ContractMigrationDaemonMembers = func() (result [insolar.GenesisAmountMigrationDaemonMembers]insolar.Reference) {
		for i, name := range insolar.GenesisNameMigrationDaemonMembers {
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

func GenerateTextReference(pulse insolar.PulseNumber, code []byte) *insolar.Reference {
	hasher := platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher()
	codeHash := hasher.Hash(code)
	return insolar.NewReference(*insolar.NewID(pulse, codeHash))
}
func GenerateReference(tp string, name string, version int) *insolar.Reference {
	contractID := fmt.Sprintf("%s::%s::v%02d", tp, name, version)
	return GenerateTextReference(insolar.BuiltinContractPulseNumber, []byte(contractID))
}

// GenesisRef returns reference to any genesis records based on the root domain.
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
	id := insolar.NewID(insolar.FirstPulseNumber, hash)
	return *insolar.NewReference(*id)
}
