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

package contracts

import (
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/genesisrefs"
	"github.com/insolar/insolar/logicrunner/builtin/contract/costcenter"
	"github.com/insolar/insolar/logicrunner/builtin/contract/member"
	"github.com/insolar/insolar/logicrunner/builtin/contract/nodedomain"
	"github.com/insolar/insolar/logicrunner/builtin/contract/rootdomain"
	"github.com/insolar/insolar/logicrunner/builtin/contract/tariff"
	"github.com/insolar/insolar/logicrunner/builtin/contract/wallet"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

func RootDomain() insolar.GenesisContractState {
	var activeMigrationDaemonMembers [insolar.GenesisAmountActiveMigrationDaemonMembers]insolar.Reference
	for i := 0; i < insolar.GenesisAmountActiveMigrationDaemonMembers; i++ {
		activeMigrationDaemonMembers[i] = genesisrefs.ContractMigrationDaemonMembers[i]
	}

	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameRootDomain,
		Prototype:  insolar.GenesisNameRootDomain,
		ParentName: "",

		Memory: mustGenMemory(&rootdomain.RootDomain{
			RootMember:             genesisrefs.ContractRootMember,
			MigrationDaemonMembers: activeMigrationDaemonMembers,
			MigrationAdminMember:   genesisrefs.ContractMigrationAdminMember,
			MigrationWallet:        genesisrefs.ContractMigrationWallet,
			CostCenter:             genesisrefs.ContractCostCenter,
			FeeWallet:              genesisrefs.ContractFeeWallet,
			BurnAddressMap:         make(foundation.StableMap),
			PublicKeyMap:           make(foundation.StableMap),
			FreeBurnAddresses:      []string{},
			NodeDomain:             genesisrefs.ContractNodeDomain,
		}),
	}
}

func NodeDomain() insolar.GenesisContractState {
	nd, _ := nodedomain.NewNodeDomain()
	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameNodeDomain,
		Prototype:  insolar.GenesisNameNodeDomain,
		ParentName: insolar.GenesisNameRootDomain,
		Memory:     mustGenMemory(nd),
	}
}

func GetMemberGenesisContractState(publicKey string, name string, parent string, walletRef insolar.Reference) insolar.GenesisContractState {
	m, err := member.New(genesisrefs.ContractRootDomain, name, publicKey, "", insolar.Reference{})
	if err != nil {
		panic(fmt.Sprintf("'%s' member constructor failed", name))
	}

	m.Wallet = walletRef

	return insolar.GenesisContractState{
		Name:       name,
		Prototype:  insolar.GenesisNameMember,
		ParentName: parent,
		Memory:     mustGenMemory(m),
	}
}

func GetWalletGenesisContractState(balance string, name string, parent string) insolar.GenesisContractState {
	w, err := wallet.New(balance)
	if err != nil {
		panic("failed to create ` " + name + "` wallet instance")
	}

	return insolar.GenesisContractState{
		Name:       name,
		Prototype:  insolar.GenesisNameWallet,
		ParentName: parent,
		Delegate:   true,
		Memory:     mustGenMemory(w),
	}
}

func GetCostCenterGenesisContractState() insolar.GenesisContractState {
	cc, err := costcenter.New(genesisrefs.ContractFeeWallet, genesisrefs.ContractStandardTariff)
	if err != nil {
		panic("failed to create cost center instance")
	}

	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameCostCenter,
		Prototype:  insolar.GenesisNameCostCenter,
		ParentName: insolar.GenesisNameRootDomain,
		Delegate:   true,
		Memory:     mustGenMemory(cc),
	}
}

func GetTariffGenesisContractState() insolar.GenesisContractState {
	t, err := tariff.New()
	if err != nil {
		panic("failed to create tariff instance")
	}

	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameStandardTariff,
		Prototype:  insolar.GenesisNameTariff,
		ParentName: insolar.GenesisNameCostCenter,
		Delegate:   true,
		Memory:     mustGenMemory(t),
	}
}

func mustGenMemory(data interface{}) []byte {
	b, err := insolar.Serialize(data)
	if err != nil {
		panic("failed to serialize contract data")
	}
	return b
}
