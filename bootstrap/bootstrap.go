/*
 *    Copyright 2018 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package bootstrap

import (
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/genesis/experiment/nodedomain"
	"github.com/insolar/insolar/genesis/experiment/rootdomain"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

const (
	nodeDomain = "nodedomain"
	nodeRecord = "noderecord"
	rootDomain = "rootdomain"
)

// Bootstrapper is a component for precreation core contracts types and RootDomain instance
type Bootstrapper struct {
	rootDomainRef *core.RecordRef
}

// GetRootDomainRef returns reference to RootDomain instance
func (b *Bootstrapper) GetRootDomainRef() *core.RecordRef {
	return b.rootDomainRef
}

// NewBootstrapper creates new Bootstrapper
func NewBootstrapper(cfg configuration.Configuration) (*Bootstrapper, error) {
	bootstrapper := &Bootstrapper{}
	bootstrapper.rootDomainRef = &core.RecordRef{}
	return bootstrapper, nil
}

var pathToContracts = "genesis/experiment/"

func getContractPath(name string) (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.Wrap(nil, "[Bootstrapper] couldn't find info about current file")
	}
	rootDir := filepath.Dir(filepath.Dir(currentFile))
	contractDir := filepath.Join(rootDir, pathToContracts)
	contractFile := name + ".go"
	return filepath.Join(contractDir, name, contractFile), nil
}

func CborInstance(contractIntance interface{}) ([]byte, error) {
	var instanceData []byte

	ch := new(codec.CborHandle)
	err := codec.NewEncoderBytes(&instanceData, ch).Encode(
		rootdomain.NewRootDomain(),
	)
	if err != nil {
		return nil, errors.Wrap(err, "[ CborInstance ] Problem with CBORing")
	}

	return instanceData, nil
}

func (b *Bootstrapper) ActivateRootDomain(am core.ArtifactManager, cb *testutil.ContractsBuilder) error {
	instanceData, err := CborInstance(rootdomain.NewRootDomain())
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootDomain ]")
	}

	contract, err := am.ActivateObj(
		core.RecordRef{}, core.RecordRef{},
		*cb.Classes[rootDomain],
		*am.RootRef(),
		instanceData,
	)
	if contract == nil {
		return errors.Wrap(err, "[ ActivateRootDomain ] Couldn't create rootdomain instance")
	}
	b.rootDomainRef = contract

	return nil
}

func (b *Bootstrapper) ActivateNodeDomain(am core.ArtifactManager, cb *testutil.ContractsBuilder) error {
	instanceData, err := CborInstance(nodedomain.NewNodeDomain())
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ]")
	}

	contract, err := am.ActivateObj(
		core.RecordRef{}, core.RecordRef{},
		*cb.Classes[nodeDomain],
		*b.rootDomainRef,
		instanceData,
	)
	if contract == nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}

	return nil
}

func (b *Bootstrapper) ActivateSmartContracts(am core.ArtifactManager, cb *testutil.ContractsBuilder) error {
	err := b.ActivateRootDomain(am, cb)
	errMsg := "[ ActivateSmartContracts ]"
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = b.ActivateNodeDomain(am, cb)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

// Start creates types and RootDomain instance
func (b *Bootstrapper) Start(c core.Components) error {
	am := c["core.Ledger"].(core.Ledger).GetArtifactManager()

	rootRefChildren, err := am.GetObjChildren(*am.RootRef())
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't get children of RootRef object")
	}
	if rootRefChildren.HasNext() {
		rootDomainRef, err := rootRefChildren.Next()
		if err != nil {
			return errors.Wrap(err, "[ Bootstrapper ] couldn't get next child of RootRef object")
		}
		b.rootDomainRef = &rootDomainRef
		return nil
	}

	jc := c["core.Ledger"].(core.Ledger).GetJetCoordinator()
	pm := c["core.Ledger"].(core.Ledger).GetPulseManager()
	currentPulse, err := pm.Current()
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't get current pulse")
	}

	network := c["core.Network"].(core.Network)
	nodeID := network.GetNodeID()

	isLightExecutor, err := jc.IsAuthorized(core.RoleLightExecutor, *am.RootRef(), currentPulse.PulseNumber, nodeID)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't authorized node")
	}
	if !isLightExecutor {
		log.Info("[ Bootstrapper ] Is not light executor. Don't build contracts")
		return nil
	}

	_, insgocc, err := testutil.Build()
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't build insgocc")
	}

	cb := testutil.NewContractBuilder(am, insgocc)
	defer cb.Clean()
	var contractNames = []string{"wallet", "member", "allowance", rootDomain, nodeDomain, nodeRecord}
	log.Info("[Bootstrapper] building contracts:", contractNames)
	contracts := make(map[string]string)
	for _, name := range contractNames {
		contractPath, _ := getContractPath(name)
		code, err := ioutil.ReadFile(contractPath)
		if err != nil {
			return errors.Wrap(err, "[ Bootstrapper ] couldn't read contract: ")
		}
		contracts[name] = string(code)
	}
	err = cb.Build(contracts)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't build contracts")
	}

	err = b.ActivateSmartContracts(am, cb)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ]")
	}

	return nil
}

// Stop implements core.Component method
func (b *Bootstrapper) Stop() error {
	return nil
}
