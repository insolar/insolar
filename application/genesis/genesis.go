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

package genesis

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/application/appfoundation"
	"github.com/insolar/insolar/application/bootstrap/contracts"
	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/insolar/insolar/insolar"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/artifact"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/pulse"
)

// BaseRecord provides methods for genesis base record manipulation.
type BaseRecord struct {
	DB             store.DB
	DropModifier   drop.Modifier
	PulseAppender  insolarPulse.Appender
	PulseAccessor  insolarPulse.Accessor
	RecordModifier object.RecordModifier
	IndexModifier  object.IndexModifier
}

// Key is genesis key.
type Key struct{}

func (Key) ID() []byte {
	return insolar.GenesisPulse.PulseNumber.Bytes()
}

func (Key) Scope() store.Scope {
	return store.ScopeGenesis
}

const (
	XNS                        = "XNS"
	MigrationDaemonUnholdDate  = 1610668800 // 15.01.2021
	MigrationDaemonVesting     = 31536000   // 365 days
	MigrationDaemonVestingStep = 2629746    // 1 month

	EnterpriseUnholdDate  = pulse.UnixTimeOfMinTimePulse
	EnterpriseVesting     = 10
	EnterpriseVestingStep = 10

	FundsUnholdDate  = pulse.UnixTimeOfMinTimePulse
	FundsVesting     = 10
	FundsVestingStep = 10

	NetworkIncentivesUnholdDate  = 1602720000 // 15.10.2020
	NetworkIncentivesVesting     = 315569520  // 10 years
	NetworkIncentivesVestingStep = 2629746    // 1 month

	ApplicationIncentivesUnholdDate  = 1602720000 // 15.10.2020
	ApplicationIncentivesVesting     = 315569520  // 10 years
	ApplicationIncentivesVestingStep = 2629746    // 1 month

	FoundationUnholdDate  = 1673740800 // 15.01.2023
	FoundationVesting     = 10
	FoundationVestingStep = 10
)

// IsGenesisRequired checks if genesis record already exists.
func (br *BaseRecord) IsGenesisRequired(ctx context.Context) (bool, error) {
	b, err := br.DB.Get(Key{})
	if err != nil {
		if err == store.ErrNotFound {
			return true, nil
		}
		return false, errors.Wrap(err, "genesis record fetch failed")
	}

	if len(b) == 0 {
		return false, errors.New("genesis record is empty (genesis hasn't properly finished)")
	}
	return false, nil
}

// Create creates new base genesis record if needed.
func (br *BaseRecord) Create(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("start storage bootstrap")

	err := br.PulseAppender.Append(
		ctx,
		insolar.Pulse{
			PulseNumber: insolar.GenesisPulse.PulseNumber,
			Entropy:     insolar.GenesisPulse.Entropy,
		},
	)
	if err != nil {
		return errors.Wrap(err, "fail to set genesis pulse")
	}
	// Add initial drop
	err = br.DropModifier.Set(ctx, drop.Drop{
		Pulse: insolar.GenesisPulse.PulseNumber,
		JetID: insolar.ZeroJetID,
	})
	if err != nil {
		return errors.Wrap(err, "fail to set initial drop")
	}

	lastPulse, err := br.PulseAccessor.Latest(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to get last pulse")
	}
	if lastPulse.PulseNumber != insolar.GenesisPulse.PulseNumber {
		return fmt.Errorf(
			"last pulse number %v is not equal to genesis special value %v",
			lastPulse.PulseNumber,
			insolar.GenesisPulse.PulseNumber,
		)
	}

	genesisID := application.GenesisRecord.ID()
	genesisRecord := record.Genesis{Hash: application.GenesisRecord}
	virtRec := record.Wrap(&genesisRecord)
	rec := record.Material{
		Virtual: virtRec,
		ID:      genesisID,
		JetID:   insolar.ZeroJetID,
	}
	err = br.RecordModifier.Set(ctx, rec)
	if err != nil {
		return errors.Wrap(err, "can't save genesis record into storage")
	}

	err = br.IndexModifier.SetIndex(
		ctx,
		pulse.MinTimePulse,
		record.Index{
			ObjID: genesisID,
			Lifeline: record.Lifeline{
				LatestState: &genesisID,
			},
			PendingRecords: []insolar.ID{},
		},
	)
	if err != nil {
		return errors.Wrap(err, "fail to set genesis index")
	}

	return br.DB.Set(Key{}, nil)
}

// Done saves genesis value. Should be called when all genesis steps finished properly.
func (br *BaseRecord) Done(ctx context.Context) error {
	return br.DB.Set(Key{}, application.GenesisRecord.Ref().Bytes())
}

// Genesis holds data and objects required for genesis on heavy node.
type Genesis struct {
	ArtifactManager artifact.Manager
	BaseRecord      *BaseRecord

	DiscoveryNodes  []application.DiscoveryNodeRegister
	PluginsDir      string
	ContractsConfig application.GenesisContractsConfig
}

// Start implements components.Starter.
func (g *Genesis) Start(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)

	isRequired, err := g.BaseRecord.IsGenesisRequired(ctx)
	inslogger.FromContext(ctx).Infof("[genesis] required=%v", isRequired)
	if err != nil {
		panic(err.Error())
	}

	if !isRequired {
		inslog.Info("[genesis] base genesis record exists, skip genesis")
		return nil
	}

	inslogger.FromContext(ctx).Info("[genesis] start...")

	inslog.Info("[genesis] create genesis record")
	err = g.BaseRecord.Create(ctx)
	if err != nil {
		return err
	}
	contracts.ContractMigrationAddressShardRefs(g.ContractsConfig.MAShardCount)
	contracts.ContractPublicKeyShardRefs(g.ContractsConfig.PKShardCount)

	inslog.Info("[genesis] store contracts")
	err = g.storeContracts(ctx)
	if err != nil {
		panic(fmt.Sprintf("[genesis] store contracts failed: %v", err))
	}

	inslog.Info("[genesis] store discovery nodes")
	discoveryNodeManager := NewDiscoveryNodeManager(g.ArtifactManager)
	err = discoveryNodeManager.StoreDiscoveryNodes(ctx, g.DiscoveryNodes)
	if err != nil {
		panic(fmt.Sprintf("[genesis] store discovery nodes failed: %v", err))
	}

	if err := g.BaseRecord.IndexModifier.UpdateLastKnownPulse(ctx, pulse.MinTimePulse); err != nil {
		panic("can't update last known pulse on genesis")
	}

	inslog.Info("[genesis] finalize genesis record")
	err = g.BaseRecord.Done(ctx)
	if err != nil {
		panic(fmt.Sprintf("[genesis] finalize genesis record failed: %v", err))
	}

	return nil
}

func (g *Genesis) storeContracts(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)

	migrationAccounts := make(foundation.StableMap)
	migrationAccounts[XNS] = genesisrefs.ContractMigrationAccount.String()

	migrationDeposits := make(foundation.StableMap)
	migrationDeposits[genesisrefs.FundsDepositName] = genesisrefs.ContractMigrationDeposit.String()

	// Hint: order matters, because of dependency contracts on each other.
	states := []application.GenesisContractState{
		contracts.RootDomain(g.ContractsConfig.PKShardCount),
		contracts.NodeDomain(),
		contracts.GetMemberGenesisContractState(g.ContractsConfig.RootPublicKey, application.GenesisNameRootMember, application.GenesisNameRootDomain, genesisrefs.ContractRootWallet),
		contracts.GetMemberGenesisContractState(g.ContractsConfig.MigrationAdminPublicKey, application.GenesisNameMigrationAdminMember, application.GenesisNameRootDomain, genesisrefs.ContractMigrationWallet),
		contracts.GetMemberGenesisContractState(g.ContractsConfig.FeePublicKey, application.GenesisNameFeeMember, application.GenesisNameRootDomain, genesisrefs.ContractFeeWallet),

		contracts.GetWalletGenesisContractState(application.GenesisNameRootWallet, application.GenesisNameRootDomain, genesisrefs.ContractRootAccount),
		contracts.GetPreWalletGenesisContractState(application.GenesisNameMigrationAdminWallet, application.GenesisNameRootDomain, migrationAccounts, migrationDeposits),
		contracts.GetWalletGenesisContractState(application.GenesisNameFeeWallet, application.GenesisNameRootDomain, genesisrefs.ContractFeeAccount),

		contracts.GetAccountGenesisContractState(g.ContractsConfig.RootBalance, application.GenesisNameRootAccount, application.GenesisNameRootDomain),
		contracts.GetAccountGenesisContractState("0", application.GenesisNameMigrationAdminAccount, application.GenesisNameRootDomain),
		contracts.GetAccountGenesisContractState("0", application.GenesisNameFeeAccount, application.GenesisNameRootDomain),

		contracts.GetDepositGenesisContractState(
			g.ContractsConfig.MDBalance,
			int64(pulse.OfUnixTime(MigrationDaemonUnholdDate)-pulse.MinTimePulse), // Lockup
			MigrationDaemonVesting,
			MigrationDaemonVestingStep,
			appfoundation.Vesting2,
			pulse.OfUnixTime(MigrationDaemonUnholdDate), // Unhold date
			application.GenesisNameMigrationAdminDeposit,
			application.GenesisNameRootDomain,
		),
		contracts.GetMigrationAdminGenesisContractState(g.ContractsConfig.LockupPeriodInPulses, g.ContractsConfig.VestingPeriodInPulses, g.ContractsConfig.VestingStepInPulses, g.ContractsConfig.MAShardCount),
		contracts.GetCostCenterGenesisContractState(g.ContractsConfig.Fee),
	}

	for i, key := range g.ContractsConfig.MigrationDaemonPublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, application.GenesisNameMigrationDaemonMembers[i], application.GenesisNameRootDomain, *insolar.NewEmptyReference()))
		states = append(states, contracts.GetMigrationDaemonGenesisContractState(i))
	}

	for i, key := range g.ContractsConfig.ApplicationIncentivesPublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, application.GenesisNameApplicationIncentivesMembers[i], application.GenesisNameRootDomain, genesisrefs.ContractApplicationIncentivesWallets[i]))

		states = append(states, contracts.GetAccountGenesisContractState("0", application.GenesisNameApplicationIncentivesAccounts[i], application.GenesisNameRootDomain))

		states = append(states, contracts.GetDepositGenesisContractState(
			application.DefaultDistributionAmount,
			int64(pulse.OfUnixTime(ApplicationIncentivesUnholdDate)-pulse.MinTimePulse),
			ApplicationIncentivesVesting,
			ApplicationIncentivesVestingStep,
			appfoundation.Vesting2,
			pulse.OfUnixTime(ApplicationIncentivesUnholdDate),
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

	for i, key := range g.ContractsConfig.NetworkIncentivesPublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, application.GenesisNameNetworkIncentivesMembers[i], application.GenesisNameRootDomain, genesisrefs.ContractNetworkIncentivesWallets[i]))
		states = append(states, contracts.GetAccountGenesisContractState("0", application.GenesisNameNetworkIncentivesAccounts[i], application.GenesisNameRootDomain))
		states = append(states, contracts.GetDepositGenesisContractState(
			application.DefaultDistributionAmount,
			int64(pulse.OfUnixTime(NetworkIncentivesUnholdDate)-pulse.MinTimePulse),
			NetworkIncentivesVesting,
			NetworkIncentivesVestingStep,
			appfoundation.Vesting2,
			pulse.OfUnixTime(NetworkIncentivesUnholdDate),
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

	for i, key := range g.ContractsConfig.FoundationPublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, application.GenesisNameFoundationMembers[i], application.GenesisNameRootDomain, genesisrefs.ContractFoundationWallets[i]))
		states = append(states, contracts.GetAccountGenesisContractState("0", application.GenesisNameFoundationAccounts[i], application.GenesisNameRootDomain))
		states = append(states, contracts.GetDepositGenesisContractState(
			application.DefaultDistributionAmount,
			int64(pulse.OfUnixTime(FoundationUnholdDate)-pulse.MinTimePulse),
			FoundationVesting,
			FoundationVestingStep,
			appfoundation.Vesting2,
			pulse.OfUnixTime(FoundationUnholdDate),
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

	for i, key := range g.ContractsConfig.FundsPublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, application.GenesisNameFundsMembers[i], application.GenesisNameRootDomain, genesisrefs.ContractFundsWallets[i]))
		states = append(states, contracts.GetAccountGenesisContractState("0", application.GenesisNameFundsAccounts[i], application.GenesisNameRootDomain))
		states = append(states, contracts.GetDepositGenesisContractState(
			application.DefaultDistributionAmount,
			int64(pulse.OfUnixTime(FundsUnholdDate)-pulse.MinTimePulse),
			FundsVesting,
			FundsVestingStep,
			appfoundation.Vesting2,
			pulse.OfUnixTime(FundsUnholdDate),
			application.GenesisNameFundsDeposits[i],
			application.GenesisNameRootDomain,
		))

		membersAccounts := make(foundation.StableMap)
		membersAccounts[XNS] = genesisrefs.ContractFundsAccounts[i].String()

		membersDeposits := make(foundation.StableMap)
		membersDeposits[genesisrefs.FundsDepositName] = genesisrefs.ContractFundsDeposits[i].String()

		states = append(states, contracts.GetPreWalletGenesisContractState(
			application.GenesisNameFundsWallets[i],
			application.GenesisNameRootDomain,
			membersAccounts,
			membersDeposits,
		))
	}

	for i, key := range g.ContractsConfig.EnterprisePublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, application.GenesisNameEnterpriseMembers[i], application.GenesisNameRootDomain, genesisrefs.ContractEnterpriseWallets[i]))
		states = append(states, contracts.GetAccountGenesisContractState("0", application.GenesisNameEnterpriseAccounts[i], application.GenesisNameRootDomain))
		states = append(states, contracts.GetDepositGenesisContractState(
			application.DefaultDistributionAmount,
			int64(pulse.OfUnixTime(EnterpriseUnholdDate)-pulse.MinTimePulse),
			EnterpriseVesting,
			EnterpriseVestingStep,
			appfoundation.Vesting2,
			pulse.OfUnixTime(EnterpriseUnholdDate),
			application.GenesisNameEnterpriseDeposits[i],
			application.GenesisNameRootDomain,
		))

		membersAccounts := make(foundation.StableMap)
		membersAccounts[XNS] = genesisrefs.ContractEnterpriseAccounts[i].String()

		membersDeposits := make(foundation.StableMap)
		membersDeposits[genesisrefs.FundsDepositName] = genesisrefs.ContractEnterpriseDeposits[i].String()

		states = append(states, contracts.GetPreWalletGenesisContractState(
			application.GenesisNameEnterpriseWallets[i],
			application.GenesisNameRootDomain,
			membersAccounts,
			membersDeposits,
		))
	}

	if g.ContractsConfig.PKShardCount <= 0 {
		panic(fmt.Sprintf("[genesis] store contracts failed: setup pk_shard_count parameter, current value %v", g.ContractsConfig.PKShardCount))
	}

	// Split genesis members by PK shards
	var membersByPKShards []foundation.StableMap
	for i := 0; i < g.ContractsConfig.PKShardCount; i++ {
		membersByPKShards = append(membersByPKShards, make(foundation.StableMap))
	}
	trimmedRootPublicKey := foundation.TrimPublicKey(g.ContractsConfig.RootPublicKey)
	index := foundation.GetShardIndex(trimmedRootPublicKey, g.ContractsConfig.PKShardCount)
	membersByPKShards[index][trimmedRootPublicKey] = genesisrefs.ContractRootMember.String()

	trimmedMigrationAdminPublicKey := foundation.TrimPublicKey(g.ContractsConfig.MigrationAdminPublicKey)
	index = foundation.GetShardIndex(trimmedMigrationAdminPublicKey, g.ContractsConfig.PKShardCount)
	membersByPKShards[index][trimmedMigrationAdminPublicKey] = genesisrefs.ContractMigrationAdminMember.String()

	trimmedFeeMemberPublicKey := foundation.TrimPublicKey(g.ContractsConfig.FeePublicKey)
	index = foundation.GetShardIndex(trimmedFeeMemberPublicKey, g.ContractsConfig.PKShardCount)
	membersByPKShards[index][trimmedFeeMemberPublicKey] = genesisrefs.ContractFeeMember.String()

	for i, key := range g.ContractsConfig.MigrationDaemonPublicKeys {
		trimmedMigrationDaemonPublicKey := foundation.TrimPublicKey(key)
		index := foundation.GetShardIndex(trimmedMigrationDaemonPublicKey, g.ContractsConfig.PKShardCount)
		membersByPKShards[index][trimmedMigrationDaemonPublicKey] = genesisrefs.ContractMigrationDaemonMembers[i].String()
	}

	for i, key := range g.ContractsConfig.NetworkIncentivesPublicKeys {
		trimmedNetworkIncentivesPublicKey := foundation.TrimPublicKey(key)
		index := foundation.GetShardIndex(trimmedNetworkIncentivesPublicKey, g.ContractsConfig.PKShardCount)
		membersByPKShards[index][trimmedNetworkIncentivesPublicKey] = genesisrefs.ContractNetworkIncentivesMembers[i].String()
	}

	for i, key := range g.ContractsConfig.ApplicationIncentivesPublicKeys {
		trimmedApplicationIncentivesPublicKey := foundation.TrimPublicKey(key)
		index := foundation.GetShardIndex(trimmedApplicationIncentivesPublicKey, g.ContractsConfig.PKShardCount)
		membersByPKShards[index][trimmedApplicationIncentivesPublicKey] = genesisrefs.ContractApplicationIncentivesMembers[i].String()
	}

	for i, key := range g.ContractsConfig.FoundationPublicKeys {
		trimmedFoundationPublicKey := foundation.TrimPublicKey(key)
		index := foundation.GetShardIndex(trimmedFoundationPublicKey, g.ContractsConfig.PKShardCount)
		membersByPKShards[index][trimmedFoundationPublicKey] = genesisrefs.ContractFoundationMembers[i].String()
	}

	for i, key := range g.ContractsConfig.FundsPublicKeys {
		trimmedFundsPublicKey := foundation.TrimPublicKey(key)
		index := foundation.GetShardIndex(trimmedFundsPublicKey, g.ContractsConfig.PKShardCount)
		membersByPKShards[index][trimmedFundsPublicKey] = genesisrefs.ContractFundsMembers[i].String()
	}

	for i, key := range g.ContractsConfig.EnterprisePublicKeys {
		trimmedEnterprisePublicKey := foundation.TrimPublicKey(key)
		index := foundation.GetShardIndex(trimmedEnterprisePublicKey, g.ContractsConfig.PKShardCount)
		membersByPKShards[index][trimmedEnterprisePublicKey] = genesisrefs.ContractEnterpriseMembers[i].String()
	}

	// Append states for shards
	for i, name := range genesisrefs.ContractPublicKeyNameShards(g.ContractsConfig.PKShardCount) {
		states = append(states, contracts.GetPKShardGenesisContractState(name, membersByPKShards[i]))
	}
	for i, name := range genesisrefs.ContractMigrationAddressNameShards(g.ContractsConfig.MAShardCount) {
		states = append(states, contracts.GetMigrationShardGenesisContractState(name, g.ContractsConfig.MigrationAddresses[i]))
	}
	for _, conf := range states {
		_, err := g.activateContract(ctx, conf)
		if err != nil {
			return errors.Wrapf(err, "failed to activate contract %v", conf.Name)
		}
		inslog.Infof("[genesis] activate contract %v", conf.Name)
	}
	return nil
}

func (g *Genesis) activateContract(ctx context.Context, state application.GenesisContractState) (*insolar.Reference, error) {
	name := state.Name
	objRef := genesisrefs.GenesisRef(name)

	protoName := name + genesisrefs.PrototypeSuffix
	protoRef := genesisrefs.GenesisRef(protoName)

	reqID, err := g.ArtifactManager.RegisterRequest(
		ctx,
		record.IncomingRequest{
			CallType: record.CTGenesis,
			Method:   name,
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to register '%v' contract", name)
	}

	parentRef := application.GenesisRecord.Ref()
	if state.ParentName != "" {
		parentRef = genesisrefs.GenesisRef(state.ParentName)
	}

	err = g.ArtifactManager.ActivateObject(
		ctx,
		*insolar.NewEmptyReference(),
		objRef,
		parentRef,
		protoRef,
		state.Memory,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to activate object for '%v'", name)
	}

	_, err = g.ArtifactManager.RegisterResult(ctx, genesisrefs.ContractRootDomain, objRef, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to register result for '%v'", name)
	}

	return insolar.NewReference(*reqID), nil
}
