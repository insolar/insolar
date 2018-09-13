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

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
	"github.com/ugorji/go/codec"
)

type Bootstrapper struct {
	rootDomainRef *core.RecordRef
}

func (b *Bootstrapper) GetRootDomainRef() *core.RecordRef {
	return b.rootDomainRef
}

func NewBootstrapper(cfg configuration.Configuration) (*Bootstrapper, error) {
	bootstrapper := &Bootstrapper{}
	return bootstrapper, nil
}

var pathWithContracts = "genesis/experiment/"

func contractPath(name string) string {
	return pathWithContracts + name + "/" + name + ".insgoc"
}

func (b *Bootstrapper) Start(c core.Components) error {
	am := c["core.Ledger"].(core.Ledger).GetManager()
	iccDir := "cmd/insgocc"
	cb := testutil.NewContractBuilder(am, iccDir+"/insgocc")
	var contractNames = []string{"wallet", "member", "allowance", "rootdomain"}
	contracts := make(map[string]string)
	for _, name := range contractNames {
		code, err := ioutil.ReadFile(contractPath(name))
		if err != nil {
			return err
		}
		contracts[name] = string(code)
	}
	err := cb.Build(contracts)
	if err != nil {
		return err
	}
	// Create code for member
	// Set code for member on ledger
	// Create code for allowance
	// Set code for allowance on ledger
	// Create proxy for member
	// Create proxy for allowance
	// Create code for wallet
	// Set code for wallet on ledger
	// Create proxy for wallet
	// Create code for rootDomain
	// Set code for rootDomain on ledger
	// Create instance of rootDomain with ArtifactManager.RootRef as parent
	// Ref to rootDomain instance return to user

	// This is just for showing idea, will be remove
	var data []byte
	err = codec.NewEncoderBytes(&data, ch).Encode(
		&struct{}{},
	)
	if err != nil {
		return err
	}
	b.rootDomainRef = am.RootRef()
	contract, err := am.ActivateObj(
		core.RecordRef{}, core.RecordRef{},
		*cb.Classes["rootdomain"],
		*am.RootRef(),
		data,
	)
	if contract == nil {
		return err
	}

	return nil
}

func (b *Bootstrapper) Stop() error {
	return nil
}
