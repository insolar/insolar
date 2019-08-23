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
	"github.com/insolar/insolar/logicrunner/builtin/contract/account"
	"github.com/insolar/insolar/logicrunner/builtin/contract/costcenter"
	"github.com/insolar/insolar/logicrunner/builtin/contract/member"
	"github.com/insolar/insolar/logicrunner/builtin/contract/migrationadmin"
	"github.com/insolar/insolar/logicrunner/builtin/contract/migrationshard"
	"github.com/insolar/insolar/logicrunner/builtin/contract/nodedomain"
	"github.com/insolar/insolar/logicrunner/builtin/contract/pkshard"
	"github.com/insolar/insolar/logicrunner/builtin/contract/rootdomain"
	"github.com/insolar/insolar/logicrunner/builtin/contract/wallet"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

func RootDomain() insolar.GenesisContractState {

	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameRootDomain,
		Prototype:  insolar.GenesisNameRootDomain,
		ParentName: "",

		Memory: mustGenMemory(&rootdomain.RootDomain{
			MigrationAddressShards: genesisrefs.ContractMigrationAddressShards,
			PublicKeyShards:        genesisrefs.ContractPublicKeyShards,
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
	m, err := member.New(genesisrefs.ContractRootDomain, name, publicKey, "", *insolar.NewEmptyReference())
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

func GetWalletGenesisContractState(name string, parent string, accountRef insolar.Reference) insolar.GenesisContractState {
	w, err := wallet.New(accountRef)
	if err != nil {
		panic("failed to create ` " + name + "` wallet instance")
	}

	return insolar.GenesisContractState{
		Name:       name,
		Prototype:  insolar.GenesisNameWallet,
		ParentName: parent,
		Memory:     mustGenMemory(w),
	}
}

func GetAccountGenesisContractState(balance string, name string, parent string) insolar.GenesisContractState {
	w, err := account.New(balance)
	if err != nil {
		panic("failed to create ` " + name + "` account instance")
	}

	return insolar.GenesisContractState{
		Name:       name,
		Prototype:  insolar.GenesisNameAccount,
		ParentName: parent,
		Memory:     mustGenMemory(w),
	}
}

func GetCostCenterGenesisContractState() insolar.GenesisContractState {
	cc, err := costcenter.New(genesisrefs.ContractFeeAccount)
	if err != nil {
		panic("failed to create cost center instance")
	}

	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameCostCenter,
		Prototype:  insolar.GenesisNameCostCenter,
		ParentName: insolar.GenesisNameRootDomain,
		Memory:     mustGenMemory(cc),
	}
}

func GetPKShardGenesisContractState(name string, members foundation.StableMap) insolar.GenesisContractState {
	s, err := pkshard.New(members)
	if err != nil {
		panic(fmt.Sprintf("'%s' shard constructor failed", name))
	}

	return insolar.GenesisContractState{
		Name:       name,
		Prototype:  insolar.GenesisNamePKShard,
		ParentName: insolar.GenesisNameRootDomain,
		Memory:     mustGenMemory(s),
	}
}

func GetMigrationShardGenesisContractState(name string, migrationAddresses []string) insolar.GenesisContractState {
	s, err := migrationshard.New(migrationAddresses)
	if err != nil {
		panic(fmt.Sprintf("'%s' shard constructor failed", name))
	}

	return insolar.GenesisContractState{
		Name:       name,
		Prototype:  insolar.GenesisNameMigrationShard,
		ParentName: insolar.GenesisNameRootDomain,
		Memory:     mustGenMemory(s),
	}
}

func GetMigrationAdminGenesisContractState(lokup int64, vesting int64) insolar.GenesisContractState {
	migrationDaemons := make(foundation.StableMap)
	for i := 0; i < insolar.GenesisAmountMigrationDaemonMembers; i++ {
		migrationDaemons[genesisrefs.ContractMigrationDaemonMembers[i].String()] = migrationadmin.StatusInactivate
	}

	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameMigrationAdmin,
		Prototype:  insolar.GenesisNameMigrationAdmin,
		ParentName: insolar.GenesisNameRootDomain,
		Memory: mustGenMemory(&migrationadmin.MigrationAdmin{
			MigrationDaemons:     migrationDaemons,
			MigrationAdminMember: genesisrefs.ContractMigrationAdminMember,
			Lokup:                lokup,
			Vesting:              vesting,
		}),
	}
}

func mustGenMemory(data interface{}) []byte {
	b, err := insolar.Serialize(data)
	if err != nil {
		panic("failed to serialize contract data")
	}
	return b
}
