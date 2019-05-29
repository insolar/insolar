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
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/insolar/insolar/bootstrap"
	"github.com/insolar/insolar/bootstrap/contracts"
	"github.com/insolar/insolar/bootstrap/rootdomain"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/artifact"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// BaseRecord provides methods for genesis base record manipulation.
type BaseRecord struct {
	DB                    store.DB
	DropModifier          drop.Modifier
	PulseAppender         pulse.Appender
	PulseAccessor         pulse.Accessor
	RecordModifier        object.RecordModifier
	IndexLifelineModifier object.LifelineModifier
}

// Key is genesis key.
type Key struct{}

func (Key) ID() []byte {
	return []byte{0x01}
}

func (Key) Scope() store.Scope {
	return store.ScopeGenesis
}

// CreateIfNeeded creates new base genesis record if needed.
// Returns reference of genesis record and flag if base record have been created.
func (gi *BaseRecord) CreateIfNeeded(ctx context.Context) (bool, error) {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("start storage bootstrap")

	getGenesisRef := func() (*insolar.Reference, error) {
		buff, err := gi.DB.Get(Key{})
		if err != nil {
			return nil, err
		}
		var genesisRef insolar.Reference
		copy(genesisRef[:], buff)
		return &genesisRef, nil
	}

	createGenesisRecord := func() error {
		err := gi.PulseAppender.Append(
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
		err = gi.DropModifier.Set(ctx, drop.Drop{JetID: insolar.ZeroJetID})
		if err != nil {
			return errors.Wrap(err, "fail to set initial drop")
		}

		lastPulse, err := gi.PulseAccessor.Latest(ctx)
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
		virtRec := record.Wrap(genesisRecord)
		rec := record.Material{
			Virtual: &virtRec,
			JetID:   insolar.ZeroJetID,
		}
		err = gi.RecordModifier.Set(ctx, genesisID, rec)
		if err != nil {
			return errors.Wrap(err, "can't save genesis record into storage")
		}

		err = gi.IndexLifelineModifier.Set(
			ctx,
			insolar.FirstPulseNumber,
			genesisID,
			object.Lifeline{
				LatestState:         &genesisID,
				LatestStateApproved: &genesisID,
				JetID:               insolar.ZeroJetID,
			},
		)
		if err != nil {
			return errors.Wrap(err, "fail to set genesis index")
		}

		return gi.DB.Set(Key{}, insolar.GenesisRecord.Ref().Bytes())
	}

	_, err := getGenesisRef()
	if err == nil {
		return false, nil
	}
	if err != store.ErrNotFound {
		return false, errors.Wrap(err, "genesis bootstrap failed")
	}

	err = createGenesisRecord()
	if err != nil {
		return true, err
	}

	return true, nil
}

// DiscoveryNodesStore provides interface for persisting discovery nodes.
type DiscoveryNodesStore interface {
	// StoreDiscoveryNodes saves discovery nodes on ledger, adds index with them to Node Domain object.
	//
	// If Node Domain object already has index - skip all saves.
	StoreDiscoveryNodes(context.Context, []insolar.DiscoveryNodeRegister) error
}

type Genesis struct {
	ArtifactManager artifact.Manager
	BaseRecord      *BaseRecord

	DiscoveryNodes  []insolar.DiscoveryNodeRegister
	PluginsDir      string
	ContractsConfig insolar.GenesisContractsConfig
}

// implements components.Starter
func (g *Genesis) Start(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)

	inslog.Info("CALL CreateIfNeeded")
	isInit, err := g.BaseRecord.CreateIfNeeded(ctx)
	if err != nil {
		return err
	}
	inslog.Infof("CALL CreateIfNeeded result:", isInit)
	if !isInit {
		return nil
	}
	inslogger.FromContext(ctx).Info("START Genesis.Init(START Genesis.Init()")

	err = g.StoreContracts(ctx)

	discoveryNodeManager := NewDiscoveryNodeManager(g.ArtifactManager)
	inslog.Info("CALL StoreDiscoveryNodes")
	return discoveryNodeManager.StoreDiscoveryNodes(ctx, g.DiscoveryNodes)
}

func (g *Genesis) StoreContracts(ctx context.Context) error {
	inslogger.FromContext(ctx).Info("CALL Genesis.StoreContracts()")

	plugins, err := readPluginsDir(g.PluginsDir)
	if err != nil {
		return err
	}

	for name, file := range plugins {
		err = g.prepareContractPrototype(ctx, name, file)
		if err != nil {
			return errors.Wrapf(err, "failed to prepare plugin's prototype %v for contract %v",
				file, name)
		}
	}

	for name, conf := range contracts.GenesisContractsStates(g.ContractsConfig) {

		err = g.activateContract(ctx, name, conf)
		if err != nil {
			return errors.Wrapf(err, "failed to activate contract %v", conf.Name)
		}
	}

	return nil
}

func (g *Genesis) activateContract(ctx context.Context, name string, state insolar.GenesisContractState) error {
	objRef := rootdomain.GenesisRef(name)
	id, err := g.ArtifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   name,
		},
	)
	if err != nil {
		return errors.Wrapf(err, "failed to register contract by name %v", name)
	}

	// just in case
	actualRef := insolar.NewReference(rootdomain.RootDomain.ID(), *id)
	if objRef != *actualRef {
		return errors.Errorf(
			"mismatch actual reference and expected for contract %v (actual=%v) (expected=%v)",
			name, *actualRef, objRef,
		)
	}

	_, err = g.ArtifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		objRef,
		rootdomain.GenesisRef(state.ParentName),
		rootdomain.GenesisRef(name+"_proto"),
		state.Delegate,
		state.Memory,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to activate %v contract", name)
	}

	_, err = g.ArtifactManager.RegisterResult(ctx, bootstrap.ContractRootDomain, objRef, nil)
	return errors.Wrapf(err, "failed to register %v contract instance", name)
}

func (g *Genesis) prepareContractPrototype(ctx context.Context, name string, binFile string) error {
	rootDomainID := rootdomain.RootDomain.ID()
	rootDomainRef := rootdomain.RootDomain.Ref()

	protoID, err := g.ArtifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   name + "_proto",
		},
	)
	if err != nil {
		return err
	}
	protoRef := insolar.NewReference(rootDomainID, *protoID)

	pluginBinary, err := ioutil.ReadFile(binFile)
	if err != nil {
		return errors.Wrap(err, "[ buildPrototypes ] Can't ReadFile")
	}
	codeReq, err := g.ArtifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   name + "_code",
		},
	)
	if err != nil {
		return errors.Wrapf(err, "[ buildPrototypes ] Can't RegisterRequest for code '%v'", name)
	}
	inslogger.FromContext(ctx).Debugf("%v code ID=%v\n", name, codeReq)

	log.Debugf("Deploying code for contract %q", name)
	codeID, err := g.ArtifactManager.DeployCode(
		ctx,
		rootDomainRef,
		*insolar.NewReference(rootDomainID, *codeReq),
		pluginBinary,
		insolar.MachineTypeGoPlugin,
	)
	if err != nil {
		return errors.Wrapf(err, "[ buildPrototypes ] Can't DeployCode for code '%v", name)
	}

	codeRef := insolar.NewReference(rootDomainID, *codeID)
	_, err = g.ArtifactManager.RegisterResult(ctx, rootDomainRef, *codeRef, nil)
	if err != nil {
		return errors.Wrapf(err, "[ buildPrototypes ] Can't SetRecord for code '%v'", name)
	}

	log.Infof("Deployed code %q for contract %q", codeRef.String(), name)

	_, err = g.ArtifactManager.ActivatePrototype(
		ctx,
		rootDomainRef,
		*protoRef,
		insolar.GenesisRecord.Ref(),
		*codeRef,
		nil,
	)
	if err != nil {
		return errors.Wrapf(err, "[ buildPrototypes ] Can't ActivatePrototypef for code '%v'", name)
	}

	_, err = g.ArtifactManager.RegisterResult(ctx, rootDomainRef, *protoRef, nil)
	return errors.Wrap(err, "failed register request")
}

func readPluginsDir(dir string) (map[string]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(), ".so") {
			result[file.Name()] = filepath.Join(dir, file.Name())
			continue
		}
	}

	return result, nil
}
