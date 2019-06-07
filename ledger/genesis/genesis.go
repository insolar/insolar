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

	"github.com/insolar/insolar/bootstrap/contracts"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/genesisrefs"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/rootdomain"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/artifact"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
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
	err = br.DropModifier.Set(ctx, drop.Drop{JetID: insolar.ZeroJetID})
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
	virtRec := record.Wrap(genesisRecord)
	rec := record.Material{
		Virtual: &virtRec,
		JetID:   insolar.ZeroJetID,
	}
	err = br.RecordModifier.Set(ctx, genesisID, rec)
	if err != nil {
		return errors.Wrap(err, "can't save genesis record into storage")
	}

	err = br.IndexLifelineModifier.Set(
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

	plugins, err := readPluginsDir(g.PluginsDir)
	if err != nil {
		return errors.Wrapf(err, "failed to read plugins dir")
	}
	inslog.Infof("[genesis] found %v plugins in %v", len(plugins), g.PluginsDir)

	for name, file := range plugins {
		err = g.prepareContractPrototype(ctx, name, file)
		if err != nil {
			return errors.Wrapf(err, "failed to prepare plugin's prototype %v for contract %v", file, name)
		}
		inslog.Infof("[genesis] code and prototype are activated for %v", name)
	}

	states := contracts.GenesisContractsStates(g.ContractsConfig)
	for _, conf := range states {
		err = g.activateContract(ctx, conf)
		if err != nil {
			return errors.Wrapf(err, "failed to activate contract %v", conf.Name)
		}
		inslog.Infof("[genesis] activate contract %v", conf.Name)
	}
	return nil
}

func (g *Genesis) activateContract(ctx context.Context, state insolar.GenesisContractState) error {
	name := state.Name

	objRef := rootdomain.GenesisRef(name)
	_, err := g.ArtifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   name,
		},
	)
	if err != nil {
		return errors.Wrapf(err, "failed to register '%v' contract", name)
	}

	parentRef := insolar.GenesisRecord.Ref()
	if state.ParentName != "" {
		parentRef = rootdomain.GenesisRef(state.ParentName)
	}

	_, err = g.ArtifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		objRef,
		parentRef,
		rootdomain.GenesisRef(name+"_proto"),
		state.Delegate,
		state.Memory,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to activate object for '%v'", name)
	}

	_, err = g.ArtifactManager.RegisterResult(ctx, genesisrefs.ContractRootDomain, objRef, nil)
	return errors.Wrapf(err, "failed to register result for '%v'", name)
}

func (g *Genesis) prepareContractPrototype(ctx context.Context, name string, binFile string) error {
	rootDomainRef := rootdomain.RootDomain.Ref()

	protoID, err := g.ArtifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   name + "_proto",
		},
	)
	if err != nil {
		return errors.Wrapf(err, "can't register request for prototype '%v'", name)
	}
	protoRef := insolar.NewReference(*protoID)
	assertGenesisRef(*protoRef, name+"_proto")

	pluginBinary, err := ioutil.ReadFile(binFile)
	if err != nil {
		return errors.Wrap(err, "failed to read plugin file")
	}
	codeReq, err := g.ArtifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   name + "_code",
		},
	)
	if err != nil {
		return errors.Wrapf(err, "failed to register request for code '%v'", name)
	}

	codeID, err := g.ArtifactManager.DeployCode(
		ctx,
		rootDomainRef,
		*insolar.NewReference(*codeReq),
		pluginBinary,
		insolar.MachineTypeGoPlugin,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to deploy code for code '%v", name)
	}

	codeRef := insolar.NewReference(*codeID)

	_, err = g.ArtifactManager.RegisterResult(ctx, rootDomainRef, *codeRef, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to register code result for '%v'", name)
	}

	_, err = g.ArtifactManager.ActivatePrototype(
		ctx,
		rootDomainRef,
		*protoRef,
		insolar.GenesisRecord.Ref(),
		*codeRef,
		nil,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to activate prototype for '%v'", name)
	}

	_, err = g.ArtifactManager.RegisterResult(ctx, rootDomainRef, *protoRef, nil)
	return errors.Wrapf(err, "failed to register request for '%v'", name)
}

func readPluginsDir(dir string) (map[string]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "open plugin dir %v failed", dir)
	}

	result := map[string]string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		soName := file.Name()
		if strings.HasSuffix(soName, ".so") {
			name := soName[:len(soName)-len(".so")]
			result[name] = filepath.Join(dir, soName)
			continue
		}
	}

	return result, nil
}

func assertGenesisRef(gotRef insolar.Reference, name string) {
	expectRef := rootdomain.GenesisRef(name)
	// check just in case
	if gotRef != expectRef {
		panic(fmt.Sprintf(
			"mismatch actual reference and expected for name %v (got=%v; expected=%v)",
			name, gotRef, expectRef,
		))
	}
}
