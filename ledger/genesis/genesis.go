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

	"github.com/insolar/insolar/bootstrap/contracts"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/genesisrefs"
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
	FundsDepositName           = "genesis_deposit"
	MigrationDaemonLockup      = 1604188800
	MigrationDaemonVesting     = 1735689600
	MigrationDaemonVestingStep = 2629746
	MigrationDaemonMaturePulse = 1735689600

	FundsAndEnterpriseLockup      = 0
	FundsAndEnterpriseVesting     = 0
	FundsAndEnterpriseVestingStep = 1
	FundsAndEnterpriseMaturePulse = 0

	NetworkIncentivesLockup      = 1672531200
	NetworkIncentivesVesting     = 1635724800
	NetworkIncentivesVestingStep = 1
	NetworkIncentivesMaturePulse = 1735689600

	ApplicationIncentivesLockup      = 1672531200
	ApplicationIncentivesVesting     = 1735689600
	ApplicationIncentivesVestingStep = 2629746
	ApplicationIncentivesMaturePulse = 1735689600

	FoundationLockup      = 1672531200
	FoundationVesting     = 1735689600
	FoundationVestingStep = 2629746
	FoundationMaturePulse = 1735689600
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

	genesisID := insolar.GenesisRecord.ID()
	genesisRecord := record.Genesis{Hash: insolar.GenesisRecord}
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
	return br.DB.Set(Key{}, insolar.GenesisRecord.Ref().Bytes())
}

// Genesis holds data and objects required for genesis on heavy node.
type Genesis struct {
	ArtifactManager artifact.Manager
	BaseRecord      *BaseRecord

	DiscoveryNodes  []insolar.DiscoveryNodeRegister
	PluginsDir      string
	ContractsConfig insolar.GenesisContractsConfig
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

	// Hint: order matters, because of dependency contracts on each other.
	states := []insolar.GenesisContractState{
		contracts.RootDomain(),
		contracts.NodeDomain(),
		contracts.GetMemberGenesisContractState(g.ContractsConfig.RootPublicKey, insolar.GenesisNameRootMember, insolar.GenesisNameRootDomain, genesisrefs.ContractRootWallet),
		contracts.GetMemberGenesisContractState(g.ContractsConfig.MigrationAdminPublicKey, insolar.GenesisNameMigrationAdminMember, insolar.GenesisNameRootDomain, genesisrefs.ContractMigrationWallet),
		contracts.GetMemberGenesisContractState(g.ContractsConfig.FeePublicKey, insolar.GenesisNameFeeMember, insolar.GenesisNameRootDomain, genesisrefs.ContractFeeWallet),
		contracts.GetMemberGenesisContractState(g.ContractsConfig.FundsAndEnterprisePublicKey, insolar.GenesisNameEnterpriseMember, insolar.GenesisNameRootDomain, genesisrefs.ContractEnterpriseWallet),

		contracts.GetWalletGenesisContractState(insolar.GenesisNameRootWallet, insolar.GenesisNameRootDomain, genesisrefs.ContractRootAccount),
		contracts.GetWalletGenesisContractState(insolar.GenesisNameMigrationAdminWallet, insolar.GenesisNameRootDomain, genesisrefs.ContractMigrationAccount),
		contracts.GetWalletGenesisContractState(insolar.GenesisNameFeeWallet, insolar.GenesisNameRootDomain, genesisrefs.ContractFeeAccount),
		contracts.GetWalletGenesisContractState(insolar.GenesisNameEnterpriseWallet, insolar.GenesisNameRootDomain, genesisrefs.ContractEnterpriseAccount),

		contracts.GetAccountGenesisContractState(g.ContractsConfig.RootBalance, insolar.GenesisNameRootAccount, insolar.GenesisNameRootDomain),
		contracts.GetAccountGenesisContractState(g.ContractsConfig.MDBalance, insolar.GenesisNameMigrationAdminAccount, insolar.GenesisNameRootDomain),
		contracts.GetAccountGenesisContractState("0", insolar.GenesisNameFeeAccount, insolar.GenesisNameRootDomain),
		contracts.GetAccountGenesisContractState(insolar.DefaultDistributionAmount, insolar.GenesisNameEnterpriseAccount, insolar.GenesisNameRootDomain),

		contracts.GetDepositGenesisContractState(
			insolar.DefaultDistributionAmount,
			int64(pulse.OfUnixTime(FundsAndEnterpriseLockup)),
			int64(pulse.OfUnixTime(FundsAndEnterpriseVesting)),
			FundsAndEnterpriseVestingStep,
			foundation.Vesting1,
			FundsAndEnterpriseMaturePulse,
			0,
			insolar.GenesisNameEnterpriseDeposit,
			insolar.GenesisNameRootDomain,
		),
		contracts.GetDepositGenesisContractState(
			"0",
			int64(pulse.OfUnixTime(MigrationDaemonLockup)),
			int64(pulse.OfUnixTime(MigrationDaemonVesting)),
			MigrationDaemonVestingStep,
			foundation.Vesting2,
			pulse.OfUnixTime(MigrationDaemonMaturePulse),
			0,
			insolar.GenesisNameMigrationAdminDeposit,
			insolar.GenesisNameRootDomain,
		),
		contracts.GetMigrationAdminGenesisContractState(g.ContractsConfig.LoсkupPeriodInPulses, g.ContractsConfig.VestingPeriodInPulses, g.ContractsConfig.VestingStepInPulses),
		contracts.GetCostCenterGenesisContractState(),
	}

	for i, key := range g.ContractsConfig.MigrationDaemonPublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, insolar.GenesisNameMigrationDaemonMembers[i], insolar.GenesisNameRootDomain, *insolar.NewEmptyReference()))
		states = append(states, contracts.GetMigrationDaemonGenesisContractState(i))
	}

	for i, key := range g.ContractsConfig.ApplicationIncentivesPublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, insolar.GenesisNameApplicationIncentivesMembers[i], insolar.GenesisNameRootDomain, genesisrefs.ContractApplicationIncentivesMembers[i]))
	}

	for i, key := range g.ContractsConfig.NetworkIncentivesPublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, insolar.GenesisNameNetworkIncentivesMembers[i], insolar.GenesisNameRootDomain, genesisrefs.ContractNetworkIncentivesMembers[i]))
	}

	for i, key := range g.ContractsConfig.FoundationPublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, insolar.GenesisNameFoundationMembers[i], insolar.GenesisNameRootDomain, genesisrefs.ContractFoundationMembers[i]))
	}

	for i := range g.ContractsConfig.ApplicationIncentivesPublicKeys {
		states = append(states, contracts.GetAccountGenesisContractState("0", insolar.GenesisNameApplicationIncentivesAccounts[i], insolar.GenesisNameRootDomain))
	}

	for i := range g.ContractsConfig.NetworkIncentivesPublicKeys {
		states = append(states, contracts.GetAccountGenesisContractState("0", insolar.GenesisNameNetworkIncentivesAccounts[i], insolar.GenesisNameRootDomain))
	}

	for i := range g.ContractsConfig.FoundationPublicKeys {
		states = append(states, contracts.GetAccountGenesisContractState("0", insolar.GenesisNameFoundationAccounts[i], insolar.GenesisNameRootDomain))
	}

	for i := range g.ContractsConfig.NetworkIncentivesPublicKeys {
		states = append(states, contracts.GetDepositGenesisContractState(
			insolar.DefaultDistributionAmount,
			int64(pulse.OfUnixTime(NetworkIncentivesLockup)),
			int64(pulse.OfUnixTime(NetworkIncentivesVesting)),
			NetworkIncentivesVestingStep,
			foundation.Vesting2,
			pulse.OfUnixTime(NetworkIncentivesMaturePulse),
			0,
			insolar.GenesisNameNetworkIncentivesDeposits[i],
			insolar.GenesisNameRootDomain,
		))
	}

	for i := range g.ContractsConfig.ApplicationIncentivesPublicKeys {
		states = append(states, contracts.GetDepositGenesisContractState(
			insolar.DefaultDistributionAmount,
			int64(pulse.OfUnixTime(ApplicationIncentivesLockup)),
			int64(pulse.OfUnixTime(ApplicationIncentivesVesting)),
			ApplicationIncentivesVestingStep,
			foundation.Vesting2,
			pulse.OfUnixTime(ApplicationIncentivesMaturePulse),
			0,
			insolar.GenesisNameApplicationIncentivesDeposits[i],
			insolar.GenesisNameRootDomain,
		))
	}

	for i := range g.ContractsConfig.FoundationPublicKeys {
		states = append(states, contracts.GetDepositGenesisContractState(
			insolar.DefaultDistributionAmount,
			int64(pulse.OfUnixTime(FoundationLockup)),
			int64(pulse.OfUnixTime(FoundationVesting)),
			FoundationVestingStep,
			foundation.Vesting2,
			pulse.OfUnixTime(FoundationMaturePulse),
			0,
			insolar.GenesisNameFoundationDeposits[i],
			insolar.GenesisNameRootDomain,
		))
	}

	for i := range g.ContractsConfig.ApplicationIncentivesPublicKeys {
		membersAccounts := make(foundation.StableMap)
		membersAccounts[XNS] = genesisrefs.ContractApplicationIncentivesAccounts[i].String()

		membersDeposits := make(foundation.StableMap)
		membersDeposits[FundsDepositName] = genesisrefs.ContractApplicationIncentivesDeposits[i].String()

		states = append(states, contracts.GetPreWalletGenesisContractState(
			insolar.GenesisNameApplicationIncentivesWallets[i],
			insolar.GenesisNameRootDomain,
			genesisrefs.ContractApplicationIncentivesWallets[i],
			membersAccounts,
			membersDeposits,
		))
	}

	for i := range g.ContractsConfig.NetworkIncentivesPublicKeys {
		membersAccounts := make(foundation.StableMap)
		membersAccounts[XNS] = genesisrefs.ContractNetworkIncentivesAccounts[i].String()

		membersDeposits := make(foundation.StableMap)
		membersDeposits[FundsDepositName] = genesisrefs.ContractNetworkIncentivesDeposits[i].String()

		states = append(states, contracts.GetPreWalletGenesisContractState(
			insolar.GenesisNameNetworkIncentivesWallets[i],
			insolar.GenesisNameRootDomain,
			genesisrefs.ContractNetworkIncentivesWallets[i],
			membersAccounts,
			membersDeposits,
		))
	}

	for i := range g.ContractsConfig.FoundationPublicKeys {
		membersAccounts := make(foundation.StableMap)
		membersAccounts[XNS] = genesisrefs.ContractFoundationAccounts[i].String()

		membersDeposits := make(foundation.StableMap)
		membersDeposits[FundsDepositName] = genesisrefs.ContractFoundationDeposits[i].String()

		states = append(states, contracts.GetPreWalletGenesisContractState(
			insolar.GenesisNameFoundationWallets[i],
			insolar.GenesisNameRootDomain,
			genesisrefs.ContractFoundationWallets[i],
			membersAccounts,
			membersDeposits,
		))
	}

	// Split genesis members by PK shards
	var MembersByPKShards [insolar.GenesisAmountPublicKeyShards]foundation.StableMap
	for i := 0; i < insolar.GenesisAmountPublicKeyShards; i++ {
		MembersByPKShards[i] = make(foundation.StableMap)
	}
	trimmedRootPublicKey := foundation.TrimPublicKey(g.ContractsConfig.RootPublicKey)
	index := foundation.GetShardIndex(trimmedRootPublicKey, insolar.GenesisAmountPublicKeyShards)
	MembersByPKShards[index][trimmedRootPublicKey] = genesisrefs.ContractRootMember.String()

	trimmedMigrationAdminPublicKey := foundation.TrimPublicKey(g.ContractsConfig.MigrationAdminPublicKey)
	index = foundation.GetShardIndex(trimmedMigrationAdminPublicKey, insolar.GenesisAmountPublicKeyShards)
	MembersByPKShards[index][trimmedMigrationAdminPublicKey] = genesisrefs.ContractMigrationAdminMember.String()

	trimmedFeeAdminPublicKey := foundation.TrimPublicKey(g.ContractsConfig.FeePublicKey)
	index = foundation.GetShardIndex(trimmedFeeAdminPublicKey, insolar.GenesisAmountPublicKeyShards)
	MembersByPKShards[index][trimmedFeeAdminPublicKey] = genesisrefs.ContractFeeMember.String()

	trimmedFundsAndEnterprisePublicKey := foundation.TrimPublicKey(g.ContractsConfig.FundsAndEnterprisePublicKey)
	index = foundation.GetShardIndex(trimmedFundsAndEnterprisePublicKey, insolar.GenesisAmountPublicKeyShards)
	MembersByPKShards[index][trimmedFundsAndEnterprisePublicKey] = genesisrefs.ContractEnterpriseMember.String()

	for i, key := range g.ContractsConfig.MigrationDaemonPublicKeys {
		trimmedMigrationDaemonPublicKey := foundation.TrimPublicKey(key)
		index := foundation.GetShardIndex(trimmedMigrationDaemonPublicKey, insolar.GenesisAmountPublicKeyShards)
		MembersByPKShards[index][trimmedMigrationDaemonPublicKey] = genesisrefs.ContractMigrationDaemonMembers[i].String()
	}

	for i, key := range g.ContractsConfig.NetworkIncentivesPublicKeys {
		trimmedNetworkIncentivesPublicKey := foundation.TrimPublicKey(key)
		index := foundation.GetShardIndex(trimmedNetworkIncentivesPublicKey, insolar.GenesisAmountPublicKeyShards)
		MembersByPKShards[index][trimmedNetworkIncentivesPublicKey] = genesisrefs.ContractNetworkIncentivesMembers[i].String()
	}

	for i, key := range g.ContractsConfig.ApplicationIncentivesPublicKeys {
		trimmedApplicationIncentivesPublicKey := foundation.TrimPublicKey(key)
		index := foundation.GetShardIndex(trimmedApplicationIncentivesPublicKey, insolar.GenesisAmountPublicKeyShards)
		MembersByPKShards[index][trimmedApplicationIncentivesPublicKey] = genesisrefs.ContractApplicationIncentivesMembers[i].String()
	}

	for i, key := range g.ContractsConfig.FoundationPublicKeys {
		trimmedFoundationPublicKey := foundation.TrimPublicKey(key)
		index := foundation.GetShardIndex(trimmedFoundationPublicKey, insolar.GenesisAmountPublicKeyShards)
		MembersByPKShards[index][trimmedFoundationPublicKey] = genesisrefs.ContractFoundationMembers[i].String()
	}

	// Append states for shards
	for i, name := range insolar.GenesisNamePublicKeyShards {
		states = append(states, contracts.GetPKShardGenesisContractState(name, MembersByPKShards[i]))
	}
	for i, name := range insolar.GenesisNameMigrationAddressShards {
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

func (g *Genesis) activateContract(ctx context.Context, state insolar.GenesisContractState) (*insolar.Reference, error) {
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

	parentRef := insolar.GenesisRecord.Ref()
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
