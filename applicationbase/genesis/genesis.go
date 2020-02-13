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

package genesis

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/artifact"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/pulse"
	"github.com/pkg/errors"
)

type BaseRecord interface {
	IsGenesisRequired(ctx context.Context) (bool, error)
	Create(ctx context.Context) error
	Done(ctx context.Context) error
}

type Options struct {
	States           []ContractState
	NodeDomainParent string
}

// Genesis holds data and objects required for genesis on heavy node.
type Genesis struct {
	ArtifactManager artifact.Manager
	IndexModifier   object.IndexModifier
	BaseRecord      BaseRecord

	DiscoveryNodes []DiscoveryNodeRegister
	GenesisOptions Options
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
	err = g.storeContracts(ctx, g.GenesisOptions.States)
	if err != nil {
		panic(fmt.Sprintf("[genesis] store contracts failed: %v", err))
	}

	inslog.Info("[genesis] store discovery nodes")
	discoveryNodeManager := NewDiscoveryNodeManager(g.ArtifactManager)
	err = discoveryNodeManager.StoreDiscoveryNodes(ctx, g.DiscoveryNodes)
	if err != nil {
		panic(fmt.Sprintf("[genesis] store discovery nodes failed: %v", err))
	}

	if err := g.IndexModifier.UpdateLastKnownPulse(ctx, pulse.MinTimePulse); err != nil {
		panic("can't update last known pulse on genesis")
	}

	inslog.Info("[genesis] finalize genesis record")
	err = g.BaseRecord.Done(ctx)
	if err != nil {
		panic(fmt.Sprintf("[genesis] finalize genesis record failed: %v", err))
	}

	return nil
}

func (g *Genesis) storeContracts(ctx context.Context, states []ContractState) error {
	inslog := inslogger.FromContext(ctx)

	states = append(states, NodeDomain(g.GenesisOptions.NodeDomainParent))

	for _, conf := range states {
		_, err := g.activateContract(ctx, conf)
		if err != nil {
			return errors.Wrapf(err, "failed to activate contract %v", conf.Name)
		}
		inslog.Infof("[genesis] activate contract %v", conf.Name)
	}
	return nil
}

func (g *Genesis) activateContract(ctx context.Context, state ContractState) (*insolar.Reference, error) {
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

	parentRef := Record.Ref()
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
