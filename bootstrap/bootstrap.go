/*
 *    Copyright 2018 Insolar
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
	wallet     = "wallet"
	member     = "member"
	allowance  = "allowance"
)

var contractNames = []string{wallet, member, allowance, rootDomain, nodeDomain, nodeRecord}

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
		return "", errors.Wrap(nil, "[ getContractPath ] couldn't find info about current file")
	}
	rootDir := filepath.Dir(filepath.Dir(currentFile))
	contractDir := filepath.Join(rootDir, pathToContracts)
	contractFile := name + ".go"
	return filepath.Join(contractDir, name, contractFile), nil
}

func getContractsMap() (map[string]string, error) {
	contracts := make(map[string]string)
	for _, name := range contractNames {
		contractPath, err := getContractPath(name)
		if err != nil {
			return nil, errors.Wrap(err, "[ contractsMap ] couldn't get path to contracts: ")
		}
		code, err := ioutil.ReadFile(filepath.Clean(contractPath))
		if err != nil {
			return nil, errors.Wrap(err, "[ contractsMap ] couldn't read contract: ")
		}
		contracts[name] = string(code)
	}
	return contracts, nil
}

func isLightExecutor(c core.Components) (bool, error) {
	am := c.Ledger.GetArtifactManager()
	jc := c.Ledger.GetJetCoordinator()
	pm := c.Ledger.GetPulseManager()
	currentPulse, err := pm.Current()
	if err != nil {
		return false, errors.Wrap(err, "[ isLightExecutor ] couldn't get current pulse")
	}

	network := c.Network
	nodeID := network.GetNodeID()

	isLightExecutor, err := jc.IsAuthorized(core.RoleLightExecutor, *am.RootRef(), currentPulse.PulseNumber, nodeID)
	if err != nil {
		return false, errors.Wrap(err, "[ isLightExecutor ] couldn't authorized node")
	}
	if !isLightExecutor {
		log.Info("[ isLightExecutor ] Is not light executor. Don't build contracts")
		return false, nil
	}
	return true, nil
}

func getRootDomainRef(c core.Components) (*core.RecordRef, error) {
	am := c.Ledger.GetArtifactManager()
	rootObj, err := am.GetObject(*am.RootRef(), nil)
	rootRefChildren := rootObj.Children()
	if err != nil {
		return nil, errors.Wrap(err, "[ getRootDomainRef ] couldn't get children of RootRef object")
	}
	if rootRefChildren.HasNext() {
		rootDomainRef, err := rootRefChildren.Next()
		if err != nil {
			return nil, errors.Wrap(err, "[ getRootDomainRef ] couldn't get next child of RootRef object")
		}
		return &rootDomainRef, nil
	}
	return nil, nil
}

func buildSmartContracts(cb *testutil.ContractsBuilder) error {
	log.Info("[ buildSmartContracts ] building contracts:", contractNames)
	contracts, err := getContractsMap()
	if err != nil {
		return errors.Wrap(err, "[ buildSmartContracts ] couldn't build contracts")
	}

	log.Info("[ buildSmartContracts ] Start building contracts ...")
	err = cb.Build(contracts)
	if err != nil {
		return errors.Wrap(err, "[ buildSmartContracts ] couldn't build contracts")
	}
	log.Info("[ buildSmartContracts ] Stop building contracts ...")

	return nil
}

func serializeInstance(contractInstance interface{}) ([]byte, error) {
	var instanceData []byte

	ch := new(codec.CborHandle)
	err := codec.NewEncoderBytes(&instanceData, ch).Encode(
		contractInstance,
	)
	if err != nil {
		return nil, errors.Wrap(err, "[ serializeInstance ] Problem with CBORing")
	}

	return instanceData, nil
}

func (b *Bootstrapper) activateRootDomain(am core.ArtifactManager, cb *testutil.ContractsBuilder) error {
	instanceData, err := serializeInstance(rootdomain.NewRootDomain())
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootDomain ]")
	}

	contract, err := am.ActivateObject(
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

func (b *Bootstrapper) activateNodeDomain(am core.ArtifactManager, cb *testutil.ContractsBuilder) error {
	instanceData, err := serializeInstance(nodedomain.NewNodeDomain())
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ]")
	}

	contract, err := am.ActivateObject(
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

func (b *Bootstrapper) activateSmartContracts(am core.ArtifactManager, cb *testutil.ContractsBuilder) error {
	err := b.activateRootDomain(am, cb)
	errMsg := "[ ActivateSmartContracts ]"
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = b.activateNodeDomain(am, cb)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

// Start creates types and RootDomain instance
func (b *Bootstrapper) Start(c core.Components) error {
	log.Info("[ Bootstrapper ] Starting Bootstrap ...")

	rootDomainRef, err := getRootDomainRef(c)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't get ref of rootDomain")
	}
	if rootDomainRef != nil {
		b.rootDomainRef = rootDomainRef
		log.Info("[ Bootstrapper ] RootDomain was found in ledger. Don't do bootstrap")
		return nil
	}

	isLightExecutor, err := isLightExecutor(c)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't check if node is light executor")
	}
	if !isLightExecutor {
		log.Info("[ Bootstrapper ] Node is not light executor. Don't do bootstrap")
		return nil
	}

	_, insgocc, err := testutil.Build()
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't build insgocc")
	}

	am := c.Ledger.GetArtifactManager()
	cb := testutil.NewContractBuilder(am, insgocc)
	defer cb.Clean()

	err = buildSmartContracts(cb)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't build contracts")
	}

	err = b.activateSmartContracts(am, cb)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ]")
	}

	return nil
}

// Stop implements core.Component method
func (b *Bootstrapper) Stop() error {
	return nil
}
