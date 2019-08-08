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
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/rootdomain"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/artifact"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
)

// BaseRecord provides methods for genesis base record manipulation.
type BaseRecord struct {
	DB             store.DB
	DropModifier   drop.Modifier
	PulseAppender  pulse.Appender
	PulseAccessor  pulse.Accessor
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
		insolar.FirstPulseNumber,
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
	rootDomain := contracts.RootDomain()
	_, err := g.activateContract(ctx, rootDomain)
	if err != nil {
		return errors.Wrapf(err, "failed to activate contract %v", rootDomain.Name)
	}
	inslog.Infof("[genesis] activate contract %v", rootDomain.Name)

	nodeDomain := contracts.NodeDomain()
	_, err = g.activateContract(ctx, nodeDomain)
	if err != nil {
		return errors.Wrapf(err, "failed to activate contract %v", nodeDomain.Name)
	}
	inslog.Infof("[genesis] activate contract %v", nodeDomain.Name)

	rootWallet := contracts.GetWalletGenesisContractState(g.ContractsConfig.RootBalance, insolar.GenesisNameRootWallet, insolar.GenesisNameRootDomain)
	rwRef, err := g.activateContract(ctx, rootWallet)
	if err != nil {
		return errors.Wrapf(err, "failed to activate contract %v", rootWallet.Name)
	}
	inslog.Infof("[genesis] activate contract %v", rootWallet.Name)

	migrationAdminWallet := contracts.GetWalletGenesisContractState(g.ContractsConfig.MDBalance, insolar.GenesisNameMigrationWallet, insolar.GenesisNameRootDomain)
	mawRef, err := g.activateContract(ctx, migrationAdminWallet)
	if err != nil {
		return errors.Wrapf(err, "failed to activate contract %v", migrationAdminWallet.Name)
	}
	inslog.Infof("[genesis] activate contract %v", migrationAdminWallet.Name)

	states := []insolar.GenesisContractState{
		contracts.GetMemberGenesisContractState(g.ContractsConfig.RootPublicKey, insolar.GenesisNameRootMember, insolar.GenesisNameRootDomain, *rwRef),
		contracts.GetMemberGenesisContractState(g.ContractsConfig.MigrationAdminPublicKey, insolar.GenesisNameMigrationAdminMember, insolar.GenesisNameRootDomain, *mawRef),
		contracts.GetWalletGenesisContractState("0", insolar.GenesisNameFeeWallet, insolar.GenesisNameRootDomain),
		contracts.GetCostCenterGenesisContractState(),
	}

	for i, key := range g.ContractsConfig.MigrationDaemonPublicKeys {
		states = append(states, contracts.GetMemberGenesisContractState(key, insolar.GenesisNameMigrationDaemonMembers[i], insolar.GenesisNameRootDomain, insolar.Reference{}))
	}
	for _, name := range insolar.GenesisNamePublicKeyShards {
		states = append(states, contracts.GetPKShardGenesisContractState(name))
	}
	for _, name := range insolar.GenesisNameMigrationAddressShards {
		states = append(states, contracts.GetMigrationShardGenesisContractState(name))
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
	objRef := rootdomain.GenesisRef(name)

	protoName := name + rootdomain.GenesisPrototypeSuffix
	protoRef := rootdomain.GenesisRef(protoName)

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
		parentRef = rootdomain.GenesisRef(state.ParentName)
	}

	err = g.ArtifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
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
