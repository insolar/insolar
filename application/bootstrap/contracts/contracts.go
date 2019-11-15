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

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/application/appfoundation"
	"github.com/insolar/insolar/application/builtin/contract/account"
	"github.com/insolar/insolar/application/builtin/contract/costcenter"
	"github.com/insolar/insolar/application/builtin/contract/deposit"
	"github.com/insolar/insolar/application/builtin/contract/member"
	"github.com/insolar/insolar/application/builtin/contract/migrationadmin"
	"github.com/insolar/insolar/application/builtin/contract/migrationdaemon"
	"github.com/insolar/insolar/application/builtin/contract/migrationshard"
	"github.com/insolar/insolar/application/builtin/contract/nodedomain"
	"github.com/insolar/insolar/application/builtin/contract/pkshard"
	"github.com/insolar/insolar/application/builtin/contract/rootdomain"
	"github.com/insolar/insolar/application/builtin/contract/wallet"
	maProxy "github.com/insolar/insolar/application/builtin/proxy/migrationshard"
	pkProxy "github.com/insolar/insolar/application/builtin/proxy/pkshard"
	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

func RootDomain(pkShardCount int) application.GenesisContractState {

	return application.GenesisContractState{
		Name:       application.GenesisNameRootDomain,
		Prototype:  application.GenesisNameRootDomain,
		ParentName: "",

		Memory: mustGenMemory(&rootdomain.RootDomain{
			PublicKeyShards: genesisrefs.ContractPublicKeyShards(pkShardCount),
		}),
	}
}

func NodeDomain() application.GenesisContractState {
	nd, _ := nodedomain.NewNodeDomain()
	return application.GenesisContractState{
		Name:       application.GenesisNameNodeDomain,
		Prototype:  application.GenesisNameNodeDomain,
		ParentName: application.GenesisNameRootDomain,
		Memory:     mustGenMemory(nd),
	}
}

func GetMemberGenesisContractState(publicKey string, name string, parent string, walletRef insolar.Reference) application.GenesisContractState {
	m, err := member.New(publicKey, "", *insolar.NewEmptyReference())
	if err != nil {
		panic(fmt.Sprintf("'%s' member constructor failed", name))
	}

	m.Wallet = walletRef

	return application.GenesisContractState{
		Name:       name,
		Prototype:  application.GenesisNameMember,
		ParentName: parent,
		Memory:     mustGenMemory(m),
	}
}

func GetWalletGenesisContractState(name string, parent string, accountRef insolar.Reference) application.GenesisContractState {
	w, err := wallet.New(accountRef)
	if err != nil {
		panic("failed to create ` " + name + "` wallet instance")
	}

	return application.GenesisContractState{
		Name:       name,
		Prototype:  application.GenesisNameWallet,
		ParentName: parent,
		Memory:     mustGenMemory(w),
	}
}

func GetPreWalletGenesisContractState(name string, parent string, accounts foundation.StableMap, deposits foundation.StableMap) application.GenesisContractState {
	return application.GenesisContractState{
		Name:       name,
		Prototype:  application.GenesisNameWallet,
		ParentName: parent,
		Memory: mustGenMemory(&wallet.Wallet{
			Accounts: accounts,
			Deposits: deposits,
		}),
	}
}

func GetAccountGenesisContractState(balance string, name string, parent string) application.GenesisContractState {
	w, err := account.New(balance)
	if err != nil {
		panic("failed to create ` " + name + "` account instance")
	}

	return application.GenesisContractState{
		Name:       name,
		Prototype:  application.GenesisNameAccount,
		ParentName: parent,
		Memory:     mustGenMemory(w),
	}
}

func GetCostCenterGenesisContractState(fee string) application.GenesisContractState {
	cc, err := costcenter.New(&genesisrefs.ContractFeeMember, fee)
	if err != nil {
		panic("failed to create cost center instance")
	}

	return application.GenesisContractState{
		Name:       application.GenesisNameCostCenter,
		Prototype:  application.GenesisNameCostCenter,
		ParentName: application.GenesisNameRootDomain,
		Memory:     mustGenMemory(cc),
	}
}

func GetPKShardGenesisContractState(name string, members foundation.StableMap) application.GenesisContractState {
	s, err := pkshard.New(members)
	if err != nil {
		panic(fmt.Sprintf("'%s' shard constructor failed", name))
	}

	return application.GenesisContractState{
		Name:       name,
		Prototype:  application.GenesisNamePKShard,
		ParentName: application.GenesisNameRootDomain,
		Memory:     mustGenMemory(s),
	}
}

func GetMigrationShardGenesisContractState(name string, migrationAddresses []string) application.GenesisContractState {
	s, err := migrationshard.New(migrationAddresses)
	if err != nil {
		panic(fmt.Sprintf("'%s' shard constructor failed", name))
	}

	return application.GenesisContractState{
		Name:       name,
		Prototype:  application.GenesisNameMigrationShard,
		ParentName: application.GenesisNameRootDomain,
		Memory:     mustGenMemory(s),
	}
}

func GetMigrationAdminGenesisContractState(lockup int64, vesting int64, vestingStep int64, maShardCount int) application.GenesisContractState {
	return application.GenesisContractState{
		Name:       application.GenesisNameMigrationAdmin,
		Prototype:  application.GenesisNameMigrationAdmin,
		ParentName: application.GenesisNameRootDomain,
		Memory: mustGenMemory(&migrationadmin.MigrationAdmin{
			MigrationAddressShards: genesisrefs.ContractMigrationAddressShards(maShardCount),
			MigrationAdminMember:   genesisrefs.ContractMigrationAdminMember,
			VestingParams: &migrationadmin.VestingParams{
				Lockup:      lockup,
				Vesting:     vesting,
				VestingStep: vestingStep,
			},
		}),
	}
}

func GetDepositGenesisContractState(
	amount string,
	lockup int64,
	vesting int64,
	vestingStep int64,
	vestingType appfoundation.VestingType,
	pulseDepositUnHold insolar.PulseNumber,
	name string, parent string,
) application.GenesisContractState {
	return application.GenesisContractState{
		Name:       name,
		Prototype:  application.GenesisNameDeposit,
		ParentName: parent,
		Memory: mustGenMemory(&deposit.Deposit{
			Balance:            amount,
			Amount:             amount,
			PulseDepositUnHold: pulseDepositUnHold,
			VestingType:        vestingType,
			TxHash:             genesisrefs.FundsDepositName,
			Lockup:             lockup,
			Vesting:            vesting,
			VestingStep:        vestingStep,
			IsConfirmed:        true,
		}),
	}
}

func GetMigrationDaemonGenesisContractState(numberMigrationDaemon int) application.GenesisContractState {

	return application.GenesisContractState{
		Name:       application.GenesisNameMigrationDaemons[numberMigrationDaemon],
		Prototype:  application.GenesisNameMigrationDaemon,
		ParentName: application.GenesisNameRootDomain,
		Memory: mustGenMemory(&migrationdaemon.MigrationDaemon{
			IsActive:              false,
			MigrationDaemonMember: genesisrefs.ContractMigrationDaemonMembers[numberMigrationDaemon],
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

func ContractPublicKeyShardRefs(pkShardCount int) {
	for _, name := range genesisrefs.ContractPublicKeyNameShards(pkShardCount) {
		genesisrefs.PredefinedPrototypes[name+genesisrefs.PrototypeSuffix] = *pkProxy.PrototypeReference
	}
}

func ContractMigrationAddressShardRefs(maShardCount int) {
	for _, name := range genesisrefs.ContractMigrationAddressNameShards(maShardCount) {
		genesisrefs.PredefinedPrototypes[name+genesisrefs.PrototypeSuffix] = *maProxy.PrototypeReference
	}
}
