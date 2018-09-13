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

func (b *Bootstrapper) Start(c core.Components) error {
	am := c["core.Ledger"].(core.Ledger).GetManager()
	iccDir := "cmd/insgocc"
	cb := testutil.NewContractBuilder(am, iccDir+"/insgocc")
	var contractNames = []string{"wallet", "member", "allowance", "rootdomain"}
	contracts := make(map[string]string)
	for _, name := range contractNames {
		contractPath, _ := getContractPath(name)
		code, err := ioutil.ReadFile(contractPath)
		if err != nil {
			return err
		}
		contracts[name] = string(code)
	}
	err := cb.Build(contracts)
	if err != nil {
		return err
	}
	var data []byte
	ch := new(codec.CborHandle)
	err = codec.NewEncoderBytes(&data, ch).Encode(
		&struct{}{},
	)
	if err != nil {
		return err
	}
	contract, err := am.ActivateObj(
		core.RecordRef{}, core.RecordRef{},
		*cb.Classes["rootdomain"],
		*am.RootRef(),
		data,
	)
	if contract == nil {
		return err
	}
	b.rootDomainRef = contract

	return nil
}

func (b *Bootstrapper) Stop() error {
	return nil
}
