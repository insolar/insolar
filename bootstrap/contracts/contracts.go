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
	"github.com/insolar/insolar/logicrunner/builtin/contract/deposit"
	"github.com/insolar/insolar/logicrunner/builtin/contract/member"
	"github.com/insolar/insolar/logicrunner/builtin/contract/migrationadmin"
	"github.com/insolar/insolar/logicrunner/builtin/contract/migrationdaemon"
	"github.com/insolar/insolar/logicrunner/builtin/contract/migrationshard"
	"github.com/insolar/insolar/logicrunner/builtin/contract/nodedomain"
	"github.com/insolar/insolar/logicrunner/builtin/contract/pkshard"
	"github.com/insolar/insolar/logicrunner/builtin/contract/rootdomain"
	"github.com/insolar/insolar/logicrunner/builtin/contract/wallet"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	maProxy "github.com/insolar/insolar/logicrunner/builtin/proxy/migrationshard"
	pkProxy "github.com/insolar/insolar/logicrunner/builtin/proxy/pkshard"
)

func RootDomain(pkShardCount int) insolar.GenesisContractState {

	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameRootDomain,
		Prototype:  insolar.GenesisNameRootDomain,
		ParentName: "",

		Memory: mustGenMemory(&rootdomain.RootDomain{
			PublicKeyShards: genesisrefs.ContractPublicKeyShards(pkShardCount),
			NodeDomain:      genesisrefs.ContractNodeDomain,
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

func GetPreWalletGenesisContractState(name string, parent string, accounts foundation.StableMap, deposits foundation.StableMap) insolar.GenesisContractState {
	return insolar.GenesisContractState{
		Name:       name,
		Prototype:  insolar.GenesisNameWallet,
		ParentName: parent,
		Memory: mustGenMemory(&wallet.Wallet{
			Accounts: accounts,
			Deposits: deposits,
		}),
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

func GetCostCenterGenesisContractState(fee string) insolar.GenesisContractState {
	cc, err := costcenter.New(&genesisrefs.ContractFeeMember, fee)
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

func GetMigrationAdminGenesisContractState(lockup int64, vesting int64, vestingStep int64, maShardCount int) insolar.GenesisContractState {
	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameMigrationAdmin,
		Prototype:  insolar.GenesisNameMigrationAdmin,
		ParentName: insolar.GenesisNameRootDomain,
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
	vestingType foundation.VestingType,
	maturePulse insolar.PulseNumber,
	pulseDepositUnHold insolar.PulseNumber,
	name string, parent string,
) insolar.GenesisContractState {
	return insolar.GenesisContractState{
		Name:       name,
		Prototype:  insolar.GenesisNameDeposit,
		ParentName: parent,
		Memory: mustGenMemory(&deposit.Deposit{
			Balance:            amount,
			Amount:             amount,
			PulseDepositUnHold: pulseDepositUnHold,
			VestingType:        vestingType,
			MaturePulse:        maturePulse,
			TxHash:             genesisrefs.FundsDepositName,
			Lockup:             lockup,
			Vesting:            vesting,
			VestingStep:        vestingStep,
		}),
	}
}

func GetMigrationDaemonGenesisContractState(numberMigrationDaemon int) insolar.GenesisContractState {

	return insolar.GenesisContractState{
		Name:       insolar.GenesisNameMigrationDaemons[numberMigrationDaemon],
		Prototype:  insolar.GenesisNameMigrationDaemon,
		ParentName: insolar.GenesisNameRootDomain,
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
