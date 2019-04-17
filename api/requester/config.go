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
	"crypto"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/insolar/insolar/platformpolicy"

	"github.com/pkg/errors"
)

// UserConfigJSON holds info about user
type UserConfigJSON struct {
	PrivateKey       string `json:"private_key"`
	Caller           string `json:"caller"`
	privateKeyObject crypto.PrivateKey
}

// RequestConfigJSON holds info about request
type RequestConfigJSON struct {
	Params   []interface{} `json:"params"`
	Method   string        `json:"method"`
	LogLevel interface{}   `json:"logLevel,omitempty"`
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

	ks := platformpolicy.NewKeyProcessor()

	if cfgJSON.PrivateKey == "" {
		privKey, err := ks.GeneratePrivateKey()
		if err != nil {
			return nil, errors.Wrap(err, "[ readUserConfigFromFile ] ")
		}
		privKeyStr, err := ks.ExportPrivateKeyPEM(privKey)
		if err != nil {
			return nil, errors.Wrap(err, "[ readUserConfigFromFile ] ")
		}
		cfgJSON.PrivateKey = string(privKeyStr)
	}

	cfgJSON.privateKeyObject, err = ks.ImportPrivateKeyPEM([]byte(cfgJSON.PrivateKey))
	if err != nil {
		return nil, errors.Wrap(err, "[ readUserConfigFromFile ] Problem with reading private key")
	}

	return cfgJSON, nil
}

// ReadRequestConfigFromFile read request config from file
func ReadRequestConfigFromFile(path string) (*RequestConfigJSON, error) {
	rConfig := &RequestConfigJSON{}
	err := readFile(path, rConfig)
	if err != nil {
		return nil, errors.Wrap(err, "[ readRequesterConfigFromFile ] ")
	}

	return rConfig, nil
}

// CreateUserConfig creates user config from arguments
func CreateUserConfig(caller string, privKey string) (*UserConfigJSON, error) {
	userConfig := UserConfigJSON{PrivateKey: privKey, Caller: caller}
	var err error

	ks := platformpolicy.NewKeyProcessor()
	userConfig.privateKeyObject, err = ks.ImportPrivateKeyPEM([]byte(privKey))
	return &userConfig, err
}
