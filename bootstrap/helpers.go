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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/inscontext"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

var pathToContracts = "application/contract/"

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

func getAbsolutePath(relativePath string) (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.Wrap(nil, "[ getFullPath ] couldn't find info about current file")
	}
	rootDir := filepath.Dir(filepath.Dir(currentFile))
	return filepath.Join(rootDir, relativePath), nil
}

func getContractPath(name string) (string, error) {
	contractDir, err := getAbsolutePath(pathToContracts)
	if err != nil {
		return "", errors.Wrap(nil, "[ getContractPath ] couldn't get absolute path to contracts")
	}
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

	isLightExecutor, err := jc.IsAuthorized(core.RoleLightExecutor, *am.GenesisRef(), currentPulse.PulseNumber, nodeID)
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
	ctx := inscontext.TODO()
	rootObj, err := am.GetObject(ctx, *am.GenesisRef(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "[ getRootDomainRef ] couldn't get children of GenesisRef object")
	}
	rootRefChildren, err := rootObj.Children(nil)
	if err != nil {
		return nil, err
	}
	if rootRefChildren.HasNext() {
		rootDomainRef, err := rootRefChildren.Next()
		if err != nil {
			return nil, errors.Wrap(err, "[ getRootDomainRef ] couldn't get next child of GenesisRef object")
		}
		return rootDomainRef, nil
	}
	return nil, nil
}
