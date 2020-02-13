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

package contracts

import (
	"fmt"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/insolar/insolar/applicationbase/genesis"

	"github.com/insolar/insolar/application/builtin/contract/account"
	"github.com/insolar/insolar/application/builtin/contract/costcenter"
	"github.com/insolar/insolar/application/builtin/contract/deposit"
	"github.com/insolar/insolar/application/builtin/contract/member"
	"github.com/insolar/insolar/application/builtin/contract/migrationadmin"
	"github.com/insolar/insolar/application/builtin/contract/migrationdaemon"
	"github.com/insolar/insolar/application/builtin/contract/migrationshard"
	"github.com/insolar/insolar/application/builtin/contract/pkshard"
	"github.com/insolar/insolar/application/builtin/contract/rootdomain"
	"github.com/insolar/insolar/application/builtin/contract/wallet"

	"github.com/insolar/insolar/application/appfoundation"
	maProxy "github.com/insolar/insolar/application/builtin/proxy/migrationshard"
	pkProxy "github.com/insolar/insolar/application/builtin/proxy/pkshard"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/pulse"
)

func RootDomain(pkShardCount int) genesis.ContractState {

	return genesis.ContractState{
		Name:       application.GenesisNameRootDomain,
		Prototype:  application.GenesisNameRootDomain,
		ParentName: "",

		Memory: genesis.MustGenMemory(&rootdomain.RootDomain{
			PublicKeyShards: genesisrefs.ContractPublicKeyShards(pkShardCount),
		}),
	}
}

func GetMemberGenesisContractState(publicKey string, name string, parent string, walletRef insolar.Reference) genesis.ContractState {
	m, err := member.New(publicKey, "", *insolar.NewEmptyReference())
	if err != nil {
		panic(fmt.Sprintf("'%s' member constructor failed", name))
	}

	m.Wallet = walletRef

	return genesis.ContractState{
		Name:       name,
		Prototype:  application.GenesisNameMember,
		ParentName: parent,
		Memory:     genesis.MustGenMemory(m),
	}
}

func GetWalletGenesisContractState(name string, parent string, accountRef insolar.Reference) genesis.ContractState {
	w, err := wallet.New(accountRef)
	if err != nil {
		panic("failed to create ` " + name + "` wallet instance")
	}

	return genesis.ContractState{
		Name:       name,
		Prototype:  application.GenesisNameWallet,
		ParentName: parent,
		Memory:     genesis.MustGenMemory(w),
	}
}

func GetPreWalletGenesisContractState(name string, parent string, accounts foundation.StableMap, deposits foundation.StableMap) genesis.ContractState {
	return genesis.ContractState{
		Name:       name,
		Prototype:  application.GenesisNameWallet,
		ParentName: parent,
		Memory: genesis.MustGenMemory(&wallet.Wallet{
			Accounts: accounts,
			Deposits: deposits,
		}),
	}
}

func GetAccountGenesisContractState(balance string, name string, parent string) genesis.ContractState {
	w, err := account.New(balance)
	if err != nil {
		panic("failed to create ` " + name + "` account instance")
	}

	return genesis.ContractState{
		Name:       name,
		Prototype:  application.GenesisNameAccount,
		ParentName: parent,
		Memory:     genesis.MustGenMemory(w),
	}
}

func GetCostCenterGenesisContractState() genesis.ContractState {
	cc, err := costcenter.New(&genesisrefs.ContractFeeMember)
	if err != nil {
		panic("failed to create cost center instance")
	}

	return genesis.ContractState{
		Name:       application.GenesisNameCostCenter,
		Prototype:  application.GenesisNameCostCenter,
		ParentName: application.GenesisNameRootDomain,
		Memory:     genesis.MustGenMemory(cc),
	}
}

func GetPKShardGenesisContractState(name string, members foundation.StableMap) genesis.ContractState {
	s, err := pkshard.New(members)
	if err != nil {
		panic(fmt.Sprintf("'%s' shard constructor failed", name))
	}

	return genesis.ContractState{
		Name:       name,
		Prototype:  application.GenesisNamePKShard,
		ParentName: application.GenesisNameRootDomain,
		Memory:     genesis.MustGenMemory(s),
	}
}

func GetMigrationShardGenesisContractState(name string, migrationAddresses []string) genesis.ContractState {
	s, err := migrationshard.New(migrationAddresses)
	if err != nil {
		panic(fmt.Sprintf("'%s' shard constructor failed", name))
	}

	return genesis.ContractState{
		Name:       name,
		Prototype:  application.GenesisNameMigrationShard,
		ParentName: application.GenesisNameRootDomain,
		Memory:     genesis.MustGenMemory(s),
	}
}

func GetMigrationAdminGenesisContractState(lockup int64, vesting int64, vestingStep int64, maShardCount int) genesis.ContractState {
	return genesis.ContractState{
		Name:       application.GenesisNameMigrationAdmin,
		Prototype:  application.GenesisNameMigrationAdmin,
		ParentName: application.GenesisNameRootDomain,
		Memory: genesis.MustGenMemory(&migrationadmin.MigrationAdmin{
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
	vesting int64,
	vestingStep int64,
	vestingType appfoundation.VestingType,
	pulseDepositUnHold insolar.PulseNumber,
	name string, parent string,
) genesis.ContractState {
	return genesis.ContractState{
		Name:       name,
		Prototype:  application.GenesisNameDeposit,
		ParentName: parent,
		Memory: genesis.MustGenMemory(&deposit.Deposit{
			Balance:            amount,
			Amount:             amount,
			PulseDepositUnHold: pulseDepositUnHold,
			VestingType:        vestingType,
			TxHash:             genesisrefs.FundsDepositName,
			Lockup:             int64(pulseDepositUnHold - pulse.MinTimePulse),
			Vesting:            vesting,
			VestingStep:        vestingStep,
			IsConfirmed:        true,
		}),
	}
}

func GetMigrationDaemonGenesisContractState(numberMigrationDaemon int) genesis.ContractState {

	return genesis.ContractState{
		Name:       application.GenesisNameMigrationDaemons[numberMigrationDaemon],
		Prototype:  application.GenesisNameMigrationDaemon,
		ParentName: application.GenesisNameRootDomain,
		Memory: genesis.MustGenMemory(&migrationdaemon.MigrationDaemon{
			IsActive:              false,
			MigrationDaemonMember: genesisrefs.ContractMigrationDaemonMembers[numberMigrationDaemon],
		}),
	}
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
