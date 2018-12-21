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

package genesis

import (
	"context"
	"crypto"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/preprocessor"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

var pathToContracts = "application/contract/"

func serializeInstance(contractInstance interface{}) ([]byte, error) {
	var instanceData []byte

	instanceData, err := core.Serialize(contractInstance)
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

func getContractsMap() (map[string]*preprocessor.ParsedFile, error) {
	contracts := make(map[string]*preprocessor.ParsedFile)
	for _, name := range contractNames {
		contractPath, err := getContractPath(name)
		if err != nil {
			return nil, errors.Wrap(err, "[ contractsMap ] couldn't get path to contracts: ")
		}
		parsed, err := preprocessor.ParseFile(contractPath)
		if err != nil {
			return nil, errors.Wrap(err, "[ contractsMap ] couldn't read contract: ")
		}
		contracts[name] = parsed
	}
	return contracts, nil
}

func getKeysFromFile(ctx context.Context, file string) (crypto.PrivateKey, string, error) {
	absPath, err := filepath.Abs(file)
	if err != nil {
		return nil, "", errors.Wrap(err, "[ getKeyFromFile ] couldn't get abs path")
	}
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, "", errors.Wrap(err, "[ getKeyFromFile ] couldn't read keys file "+absPath)
	}
	var keys map[string]string
	err = json.Unmarshal(data, &keys)
	if err != nil {
		return nil, "", errors.Wrapf(err, "[ getKeyFromFile ] couldn't unmarshal data from %s", absPath)
	}
	if keys["private_key"] == "" {
		return nil, "", errors.New("[ getKeyFromFile ] empty private key")
	}
	if keys["public_key"] == "" {
		return nil, "", errors.New("[ getKeyFromFile ] empty public key")
	}
	kp := platformpolicy.NewKeyProcessor()
	key, err := kp.ImportPrivateKeyPEM([]byte(keys["private_key"]))
	if err != nil {
		return nil, "", errors.Wrapf(err, "[ getKeyFromFile ] couldn't import private key")
	}
	return key, keys["public_key"], nil
}
