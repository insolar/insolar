// Copyright 2020 Insolar Technologies GmbH
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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/application/appfoundation"
	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/insolar/insolar/application/genesisrefs/contracts"
	"github.com/insolar/insolar/applicationbase/genesis"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/pulse"
	"github.com/pkg/errors"
)

const (
	XNS                        = "XNS"
	MigrationDaemonUnholdDate  = 1596456000 // 03.08.2020 12-00-00
	MigrationDaemonVesting     = 0
	MigrationDaemonVestingStep = 0

	NetworkIncentivesUnholdStartDate = 1583020800 // 01.03.2020 00-00-00
	NetworkIncentivesVesting         = 0
	NetworkIncentivesVestingStep     = 0

	ApplicationIncentivesUnholdStartDate = 1609459200 // 01.01.2021 00-00-00
	ApplicationIncentivesVesting         = 0
	ApplicationIncentivesVestingStep     = 0

	FoundationUnholdStartDate = 1609459200 // 01.01.2021 00-00-00
	FoundationVestingPeriod   = 0
	FoundationVestingStep     = 0
)

func initStates(configPath, genesisConfigPath string) ([]genesis.GenesisContractState, error) {
	cfgHolder := configuration.NewHolder()
	var err error
	if len(configPath) != 0 {
		err = cfgHolder.LoadFromFile(configPath)
	} else {
		err = cfgHolder.Load()
	}
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err.Error())
	}

	b, err := ioutil.ReadFile(genesisConfigPath)
	if err != nil {
		log.Fatalf("failed to load genesis configuration from file: %v", genesisConfigPath)
	}
	var config struct {
		ContractsConfig application.GenesisContractsConfig
	}
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("failed to pares genesis configuration from file: %v", genesisConfigPath)
	}

	contractsConfig := config.ContractsConfig

	migrationAccounts := make(foundation.StableMap)
	migrationAccounts[XNS] = genesisrefs.ContractMigrationAccount.String()

	migrationDeposits := make(foundation.StableMap)
	migrationDeposits[genesisrefs.FundsDepositName] = genesisrefs.ContractMigrationDeposit.String()

	contracts.ContractMigrationAddressShardRefs(contractsConfig.MAShardCount)
	contracts.ContractPublicKeyShardRefs(contractsConfig.PKShardCount)

	// Hint: order matters, because of dependency contracts on each other.
	states := []genesis.GenesisContractState{
		contracts.RootDomain(contractsConfig.PKShardCount),
		contracts.GetMemberGenesisContractState(contractsConfig.RootPublicKey, application.GenesisNameRootMember, application.GenesisNameRootDomain, genesisrefs.ContractRootWallet),
		contracts.GetMemberGenesisContractState(contractsConfig.MigrationAdminPublicKey, application.GenesisNameMigrationAdminMember, application.GenesisNameRootDomain, genesisrefs.ContractMigrationWallet),
		contracts.GetMemberGenesisContractState(contractsConfig.FeePublicKey, application.GenesisNameFeeMember, application.GenesisNameRootDomain, genesisrefs.ContractFeeWallet),

		contracts.GetWalletGenesisContractState(application.GenesisNameRootWallet, application.GenesisNameRootDomain, genesisrefs.ContractRootAccount),
		contracts.GetPreWalletGenesisContractState(application.GenesisNameMigrationAdminWallet, application.GenesisNameRootDomain, migrationAccounts, migrationDeposits),
		contracts.GetWalletGenesisContractState(application.GenesisNameFeeWallet, application.GenesisNameRootDomain, genesisrefs.ContractFeeAccount),

		contracts.GetAccountGenesisContractState(contractsConfig.RootBalance, application.GenesisNameRootAccount, application.GenesisNameRootDomain),
		contracts.GetAccountGenesisContractState("0", application.GenesisNameMigrationAdminAccount, application.GenesisNameRootDomain),
		contracts.GetAccountGenesisContractState("0", application.GenesisNameFeeAccount, application.GenesisNameRootDomain),

		contracts.GetDepositGenesisContractState(
			contractsConfig.MDBalance,
			MigrationDaemonVesting,
			MigrationDaemonVestingStep,
			appfoundation.Vesting2,
			pulse.OfUnixTime(MigrationDaemonUnholdDate), // Unhold date
			application.GenesisNameMigrationAdminDeposit,
			application.GenesisNameRootDomain,
		),
		contracts.GetMigrationAdminGenesisContractState(contractsConfig.LockupPeriodInPulses, contractsConfig.VestingPeriodInPulses, contractsConfig.VestingStepInPulses, contractsConfig.MAShardCount),
		contracts.GetCostCenterGenesisContractState(),
	}

	for i, key := range contractsConfig.MigrationDaemonPublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, application.GenesisNameMigrationDaemonMembers[i], application.GenesisNameRootDomain, *insolar.NewEmptyReference()))
		states = append(states, contracts.GetMigrationDaemonGenesisContractState(i))
	}

	for i, key := range contractsConfig.ApplicationIncentivesPublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, application.GenesisNameApplicationIncentivesMembers[i], application.GenesisNameRootDomain, genesisrefs.ContractApplicationIncentivesWallets[i]))

		states = append(states, contracts.GetAccountGenesisContractState("0", application.GenesisNameApplicationIncentivesAccounts[i], application.GenesisNameRootDomain))

		unholdWithMonth := time.Unix(ApplicationIncentivesUnholdStartDate, 0).AddDate(0, i, 0).Unix()

		states = append(states, contracts.GetDepositGenesisContractState(
			application.AppIncentivesDistributionAmount,
			ApplicationIncentivesVesting,
			ApplicationIncentivesVestingStep,
			appfoundation.Vesting2,
			pulse.OfUnixTime(unholdWithMonth),
			application.GenesisNameApplicationIncentivesDeposits[i],
			application.GenesisNameRootDomain,
		))

		membersAccounts := make(foundation.StableMap)
		membersAccounts[XNS] = genesisrefs.ContractApplicationIncentivesAccounts[i].String()

		membersDeposits := make(foundation.StableMap)
		membersDeposits[genesisrefs.FundsDepositName] = genesisrefs.ContractApplicationIncentivesDeposits[i].String()

		states = append(states, contracts.GetPreWalletGenesisContractState(
			application.GenesisNameApplicationIncentivesWallets[i],
			application.GenesisNameRootDomain,
			membersAccounts,
			membersDeposits,
		))
	}

	for i, key := range contractsConfig.NetworkIncentivesPublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, application.GenesisNameNetworkIncentivesMembers[i], application.GenesisNameRootDomain, genesisrefs.ContractNetworkIncentivesWallets[i]))
		states = append(states, contracts.GetAccountGenesisContractState("0", application.GenesisNameNetworkIncentivesAccounts[i], application.GenesisNameRootDomain))

		unholdWithMonth := time.Unix(NetworkIncentivesUnholdStartDate, 0).AddDate(0, i, 0).Unix()

		states = append(states, contracts.GetDepositGenesisContractState(
			application.NetworkIncentivesDistributionAmount,
			NetworkIncentivesVesting,
			NetworkIncentivesVestingStep,
			appfoundation.Vesting2,
			pulse.OfUnixTime(unholdWithMonth),
			application.GenesisNameNetworkIncentivesDeposits[i],
			application.GenesisNameRootDomain,
		))

		membersAccounts := make(foundation.StableMap)
		membersAccounts[XNS] = genesisrefs.ContractNetworkIncentivesAccounts[i].String()

		membersDeposits := make(foundation.StableMap)
		membersDeposits[genesisrefs.FundsDepositName] = genesisrefs.ContractNetworkIncentivesDeposits[i].String()

		states = append(states, contracts.GetPreWalletGenesisContractState(
			application.GenesisNameNetworkIncentivesWallets[i],
			application.GenesisNameRootDomain,
			membersAccounts,
			membersDeposits,
		))
	}

	for i, key := range contractsConfig.FoundationPublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, application.GenesisNameFoundationMembers[i], application.GenesisNameRootDomain, genesisrefs.ContractFoundationWallets[i]))
		states = append(states, contracts.GetAccountGenesisContractState("0", application.GenesisNameFoundationAccounts[i], application.GenesisNameRootDomain))

		unholdWithMonth := time.Unix(FoundationUnholdStartDate, 0).AddDate(0, i, 0).Unix()

		states = append(states, contracts.GetDepositGenesisContractState(
			application.FoundationDistributionAmount,
			FoundationVestingPeriod,
			FoundationVestingStep,
			appfoundation.Vesting2,
			pulse.OfUnixTime(unholdWithMonth),
			application.GenesisNameFoundationDeposits[i],
			application.GenesisNameRootDomain,
		))

		membersAccounts := make(foundation.StableMap)
		membersAccounts[XNS] = genesisrefs.ContractFoundationAccounts[i].String()

		membersDeposits := make(foundation.StableMap)
		membersDeposits[genesisrefs.FundsDepositName] = genesisrefs.ContractFoundationDeposits[i].String()

		states = append(states, contracts.GetPreWalletGenesisContractState(
			application.GenesisNameFoundationWallets[i],
			application.GenesisNameRootDomain,
			membersAccounts,
			membersDeposits,
		))
	}

	for i, key := range contractsConfig.EnterprisePublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, application.GenesisNameEnterpriseMembers[i], application.GenesisNameRootDomain, genesisrefs.ContractEnterpriseWallets[i]))
		states = append(states, contracts.GetAccountGenesisContractState(
			application.EnterpriseDistributionAmount,
			application.GenesisNameEnterpriseAccounts[i],
			application.GenesisNameRootDomain,
		))

		membersAccounts := make(foundation.StableMap)
		membersAccounts[XNS] = genesisrefs.ContractEnterpriseAccounts[i].String()

		membersDeposits := make(foundation.StableMap)

		states = append(states, contracts.GetPreWalletGenesisContractState(
			application.GenesisNameEnterpriseWallets[i],
			application.GenesisNameRootDomain,
			membersAccounts,
			membersDeposits,
		))
	}

	if contractsConfig.PKShardCount <= 0 {
		panic(fmt.Sprintf("[genesis] store contracts failed: setup pk_shard_count parameter, current value %v", contractsConfig.PKShardCount))
	}
	if contractsConfig.VestingStepInPulses > 0 && contractsConfig.VestingPeriodInPulses%contractsConfig.VestingStepInPulses != 0 {
		panic(fmt.Sprintf("[genesis] store contracts failed: vesting_pulse_period (%d) is not a multiple of vesting_pulse_step (%d)", contractsConfig.VestingPeriodInPulses, contractsConfig.VestingStepInPulses))
	}

	// Split genesis members by PK shards
	var membersByPKShards []foundation.StableMap
	for i := 0; i < contractsConfig.PKShardCount; i++ {
		membersByPKShards = append(membersByPKShards, make(foundation.StableMap))
	}
	trimmedRootPublicKey, err := foundation.ExtractCanonicalPublicKey(contractsConfig.RootPublicKey)
	if err != nil {
		panic(errors.Wrapf(err, "[genesis] extracting canonical pk failed, current value %v", contractsConfig.RootPublicKey))
	}
	index := foundation.GetShardIndex(trimmedRootPublicKey, contractsConfig.PKShardCount)
	membersByPKShards[index][trimmedRootPublicKey] = genesisrefs.ContractRootMember.String()

	trimmedMigrationAdminPublicKey, err := foundation.ExtractCanonicalPublicKey(contractsConfig.MigrationAdminPublicKey)
	if err != nil {
		panic(errors.Wrapf(err, "[genesis] extracting canonical pk failed, current value %v", contractsConfig.MigrationAdminPublicKey))
	}
	index = foundation.GetShardIndex(trimmedMigrationAdminPublicKey, contractsConfig.PKShardCount)
	membersByPKShards[index][trimmedMigrationAdminPublicKey] = genesisrefs.ContractMigrationAdminMember.String()

	trimmedFeeMemberPublicKey, err := foundation.ExtractCanonicalPublicKey(contractsConfig.FeePublicKey)
	if err != nil {
		panic(errors.Wrapf(err, "[genesis] extracting canonical pk failed, current value %v", contractsConfig.FeePublicKey))
	}
	index = foundation.GetShardIndex(trimmedFeeMemberPublicKey, contractsConfig.PKShardCount)
	membersByPKShards[index][trimmedFeeMemberPublicKey] = genesisrefs.ContractFeeMember.String()

	for i, key := range contractsConfig.MigrationDaemonPublicKeys {
		trimmedMigrationDaemonPublicKey, err := foundation.ExtractCanonicalPublicKey(key)
		if err != nil {
			panic(errors.Wrapf(err, "[genesis] extracting canonical pk failed, current value %v", key))
		}
		index := foundation.GetShardIndex(trimmedMigrationDaemonPublicKey, contractsConfig.PKShardCount)
		membersByPKShards[index][trimmedMigrationDaemonPublicKey] = genesisrefs.ContractMigrationDaemonMembers[i].String()
	}

	for i, key := range contractsConfig.NetworkIncentivesPublicKeys {
		trimmedNetworkIncentivesPublicKey, err := foundation.ExtractCanonicalPublicKey(key)
		if err != nil {
			panic(errors.Wrapf(err, "[genesis] extracting canonical pk failed, current value %v", key))
		}
		index := foundation.GetShardIndex(trimmedNetworkIncentivesPublicKey, contractsConfig.PKShardCount)
		membersByPKShards[index][trimmedNetworkIncentivesPublicKey] = genesisrefs.ContractNetworkIncentivesMembers[i].String()
	}

	for i, key := range contractsConfig.ApplicationIncentivesPublicKeys {
		trimmedApplicationIncentivesPublicKey, err := foundation.ExtractCanonicalPublicKey(key)
		if err != nil {
			panic(errors.Wrapf(err, "[genesis] extracting canonical pk failed, current value %v", key))
		}
		index := foundation.GetShardIndex(trimmedApplicationIncentivesPublicKey, contractsConfig.PKShardCount)
		membersByPKShards[index][trimmedApplicationIncentivesPublicKey] = genesisrefs.ContractApplicationIncentivesMembers[i].String()
	}

	for i, key := range contractsConfig.FoundationPublicKeys {
		trimmedFoundationPublicKey, err := foundation.ExtractCanonicalPublicKey(key)
		if err != nil {
			panic(errors.Wrapf(err, "[genesis] extracting canonical pk failed, current value %v", key))
		}
		index := foundation.GetShardIndex(trimmedFoundationPublicKey, contractsConfig.PKShardCount)
		membersByPKShards[index][trimmedFoundationPublicKey] = genesisrefs.ContractFoundationMembers[i].String()
	}

	for i, key := range contractsConfig.EnterprisePublicKeys {
		trimmedEnterprisePublicKey, err := foundation.ExtractCanonicalPublicKey(key)
		if err != nil {
			panic(errors.Wrapf(err, "[genesis] extracting canonical pk failed, current value %v", key))
		}
		index := foundation.GetShardIndex(trimmedEnterprisePublicKey, contractsConfig.PKShardCount)
		membersByPKShards[index][trimmedEnterprisePublicKey] = genesisrefs.ContractEnterpriseMembers[i].String()
	}

	// Append states for shards
	for i, name := range genesisrefs.ContractPublicKeyNameShards(contractsConfig.PKShardCount) {
		states = append(states, contracts.GetPKShardGenesisContractState(name, membersByPKShards[i]))
	}
	for i, name := range genesisrefs.ContractMigrationAddressNameShards(contractsConfig.MAShardCount) {
		states = append(states, contracts.GetMigrationShardGenesisContractState(name, contractsConfig.MigrationAddresses[i]))
	}

	return states, nil
}
