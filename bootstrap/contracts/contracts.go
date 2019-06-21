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
	"strconv"

	"github.com/insolar/insolar/application/contract/costcenter"
	"github.com/insolar/insolar/application/contract/member"
	"github.com/insolar/insolar/application/contract/nodedomain"
	"github.com/insolar/insolar/application/contract/rootdomain"
	"github.com/insolar/insolar/application/contract/tariff"
	"github.com/insolar/insolar/application/contract/wallet"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/genesisrefs"
)

// GenesisContractsStates returns list contract configs for genesis.
//
// Hint: order matters, because of dependency contracts on each other.
func GenesisContractsStates(cfg insolar.GenesisContractsConfig) []insolar.GenesisContractState {
	result := []insolar.GenesisContractState{
		rootDomain(),
		nodeDomain(),
		getMemberGenesisContractState(cfg.RootPublicKey, insolar.GenesisNameRootMember, insolar.GenesisNameRootDomain),
		getWalletGenesisContractState(cfg.RootBalance, insolar.GenesisNameRootWallet, insolar.GenesisNameRootMember),
		getMemberGenesisContractState(cfg.MigrationAdminPublicKey, insolar.GenesisNameMigrationAdminMember, insolar.GenesisNameRootDomain),
		getWalletGenesisContractState(cfg.MDBalance, insolar.GenesisNameMigrationWallet, insolar.GenesisNameMigrationAdminMember),
		getCostCenterGenesisContractState(),
		getTariffGenesisContractState(),
	}

	for i, key := range cfg.MigrationDaemonPublicKeys {
		result = append(result, getMemberGenesisContractState(key, insolar.GenesisNameMigrationDaemonMembers[i], insolar.GenesisNameRootDomain))
	}

	return result
}

func rootDomain() insolar.GenesisContractState {
	if len(genesisrefs.ContractMigrationDaemonMembers) < insolar.GenesisAmountActiveMigrationDaemonMembers {
		panic("need at least '" + strconv.Itoa(insolar.GenesisAmountActiveMigrationDaemonMembers) + "' migration daemons")
	}
	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameRootDomain,
		Prototype:  insolar.GenesisNameRootDomain,
		ParentName: "",

		Memory: mustGenMemory(&rootdomain.RootDomain{
			RootMember:             genesisrefs.ContractRootMember,
			MigrationDaemonMembers: genesisrefs.ContractMigrationDaemonMembers[:insolar.GenesisAmountActiveMigrationDaemonMembers],
			MigrationAdminMember:   genesisrefs.ContractMigrationAdminMember,
			MigrationWallet:        genesisrefs.ContractMigrationWallet,
			CostCenter:             genesisrefs.ContractCostCenter,
			CommissionWallet:       genesisrefs.ContractCommissionWallet,
			BurnAddressMap:         map[string]insolar.Reference{},
			PublicKeyMap:           map[string]insolar.Reference{},
			FreeBurnAddresses:      []string{},
			NodeDomain:             genesisrefs.ContractNodeDomain,
		}),
	}
}

func nodeDomain() insolar.GenesisContractState {
	nd, _ := nodedomain.NewNodeDomain()
	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameNodeDomain,
		Prototype:  insolar.GenesisNameNodeDomain,
		ParentName: insolar.GenesisNameRootDomain,
		Memory:     mustGenMemory(nd),
	}
}

func getMemberGenesisContractState(publicKey string, name string, parent string) insolar.GenesisContractState {
	m, err := member.New(name, publicKey)
	if err != nil {
		panic(fmt.Sprintf("'%s' member constructor failed", name))
	}

	return insolar.GenesisContractState{
		Name:       name,
		Prototype:  insolar.GenesisNameMember,
		ParentName: parent,
		Memory:     mustGenMemory(m),
	}
}

func getWalletGenesisContractState(balance string, name string, parrent string) insolar.GenesisContractState {
	w, err := wallet.New(balance)
	if err != nil {
		panic("failed to create ` " + name + "` wallet instance")
	}

	return insolar.GenesisContractState{
		Name:       name,
		Prototype:  insolar.GenesisNameWallet,
		ParentName: parrent,
		Delegate:   true,
		Memory:     mustGenMemory(w),
	}
}

func getCostCenterGenesisContractState() insolar.GenesisContractState {
	cc, err := costcenter.New(genesisrefs.ContractCommissionWallet, genesisrefs.ContractTariff)
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

func getTariffGenesisContractState() insolar.GenesisContractState {
	t, err := tariff.New("5000000")
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
