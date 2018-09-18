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
	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
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
		return "", errors.Wrap(nil, "couldn't find info about current file")
	}
	rootDir := filepath.Dir(filepath.Dir(currentFile))
	contractDir := filepath.Join(rootDir, pathToContracts)
	contractFile := name + ".insgoc"
	return filepath.Join(contractDir, name, contractFile), nil
}

// Start creates types and RootDomain instance
func (b *Bootstrapper) Start(c core.Components) error {
	am := c["core.Ledger"].(core.Ledger).GetArtifactManager()

	rootRefChildren, err := am.GetObjChildren(*am.RootRef())
	if err != nil {
		return errors.Wrap(err, "couldn't get children of RootRef object")
	}
	if rootRefChildren.HasNext() {
		rootDomainRef, err := rootRefChildren.Next()
		if err != nil {
			return errors.Wrap(err, "couldn't get next child of RootRef object")
		}
		b.rootDomainRef = &rootDomainRef
		return nil
	}

	jc := c["core.Ledger"].(core.Ledger).GetJetCoordinator()
	pm := c["core.Ledger"].(core.Ledger).GetPulseManager()
	currentPulse, err := pm.Current()
	if err != nil {
		return errors.Wrap(err, "couldn't get current pulse")
	}

	network := c["core.Network"].(core.Network)
	nodeID := network.GetNodeID()

	isLightExecutor, err := jc.IsAuthorized(core.RoleLightExecutor, *am.RootRef(), currentPulse.PulseNumber, nodeID)
	if err != nil {
		return errors.Wrap(err, "couldn't get children of RootRef object")
	}
	if !isLightExecutor {
		return nil
	}

	_, insgocc, err := testutil.Build()
	if err != nil {
		return errors.Wrap(err, "couldn't build insgocc")
	}
	cb, cleaner := testutil.NewContractBuilder(am, insgocc)
	defer cleaner()
	var contractNames = []string{"wallet", "member", "allowance", "rootdomain"}
	contracts := make(map[string]string)
	for _, name := range contractNames {
		contractPath, _ := getContractPath(name)
		code, err := ioutil.ReadFile(contractPath)
		if err != nil {
			return errors.Wrap(err, "couldn't read contract: ")
		}
		contracts[name] = string(code)
	}
	err = cb.Build(contracts)
	if err != nil {
		return errors.Wrap(err, "couldn't build contracts")
	}
	var data []byte

	ch := new(codec.CborHandle)
	err = codec.NewEncoderBytes(&data, ch).Encode(
		&struct{}{},
	)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper: Start ]")
	}

	contract, err := am.ActivateObj(
		core.RecordRef{}, core.RecordRef{},
		*cb.Classes["rootdomain"],
		*am.RootRef(),
		data,
	)
	if contract == nil {
		return errors.Wrap(err, "couldn't create rootdomain instance")
	}
	b.rootDomainRef = contract

	return nil
}

// Stop implements core.Component method
func (b *Bootstrapper) Stop() error {
	return nil
}
