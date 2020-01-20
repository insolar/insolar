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

package requester

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/insolar/x-crypto"

	"github.com/insolar/insolar/insolar/secrets"
	"github.com/pkg/errors"
)

// UserConfigJSON holds info about user
type UserConfigJSON struct {
	PrivateKey       string `json:"private_key"`
	PublicKey        string `json:"public_key"`
	Caller           string `json:"caller"`
	privateKeyObject crypto.PrivateKey
}

func readFile(path string, configType interface{}) error {
	var rawConf []byte
	var err error
	if path == "-" {
		rawConf, err = ioutil.ReadAll(os.Stdin)
	} else {
		rawConf, err = ioutil.ReadFile(filepath.Clean(path))
	}
	if err != nil {
		return errors.Wrap(err, "[ readFile ] Problem with reading config")
	}

	err = json.Unmarshal(rawConf, &configType)
	if err != nil {
		return errors.Wrap(err, "[ readFile ] Problem with unmarshaling config")
	}

	return nil
}

// ReadUserConfigFromFile read user config from file
func ReadUserConfigFromFile(file string) (*UserConfigJSON, error) {
	cfgJSON := &UserConfigJSON{}
	err := readFile(file, cfgJSON)
	if err != nil {
		return nil, errors.Wrap(err, "[ readUserConfigFromFile ] ")
	}

	if cfgJSON.PrivateKey == "" {
		privKey, err := secrets.GeneratePrivateKey256k()
		if err != nil {
			return nil, errors.Wrap(err, "[ readUserConfigFromFile ] ")
		}
		privKeyStr, err := secrets.ExportPrivateKeyPEM(privKey)
		if err != nil {
			return nil, errors.Wrap(err, "[ readUserConfigFromFile ] ")
		}
		cfgJSON.PrivateKey = string(privKeyStr)
	}

	cfgJSON.privateKeyObject, err = secrets.ImportPrivateKeyPEM([]byte(cfgJSON.PrivateKey))
	if err != nil {
		return nil, errors.Wrap(err, "[ readUserConfigFromFile ] Problem with reading private key")
	}

	return cfgJSON, nil
}

// ReadRequestParamsFromFile read request config from file
func ReadRequestParamsFromFile(path string) (*Params, error) {
	rParams := &Params{}
	err := readFile(path, rParams)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read call params from file")
	}

	return rParams, nil
}

// CreateUserConfig creates user config from arguments
func CreateUserConfig(caller string, privKey string, publicKey string) (*UserConfigJSON, error) {
	userConfig := UserConfigJSON{PrivateKey: privKey, Caller: caller, PublicKey: publicKey}
	var err error

	userConfig.privateKeyObject, err = secrets.ImportPrivateKeyPEM([]byte(privKey))
	return &userConfig, err
}
