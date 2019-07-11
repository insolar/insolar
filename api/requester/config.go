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

// Request is a representation of request struct to api
type Request struct {
	JSONRPC  string `json:"jsonrpc"`
	ID       int    `json:"id"`
	Method   string `json:"method"`
	Params   Params `json:"params"`
	LogLevel string `json:"logLevel,omitempty"`
	Test     string `json:"test,omitempty"`
}

type Params struct {
	Seed       string      `json:"seed"`
	CallSite   string      `json:"callSite"`
	CallParams interface{} `json:"callParams"`
	Reference  string      `json:"reference"`
	PublicKey  string      `json:"memberPublicKey"`
}

type ContractAnswer struct {
	JSONRPC string  `json:"jsonrpc"`
	ID      int     `json:"id"`
	Result  *Result `json:"result,omitempty"`
	Error   *Error  `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    Data   `json:"data,omitempty"`
}

type Data struct {
	TraceID string `json:"traceID,omitempty"`
}

type Result struct {
	ContractResult interface{} `json:"callResult,omitempty"`
	TraceID        string      `json:"traceID,omitempty"`
}

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
func ReadRequestConfigFromFile(path string) (*Request, error) {
	rConfig := &Request{}
	err := readFile(path, rConfig)
	if err != nil {
		return nil, errors.Wrap(err, "[ readRequesterConfigFromFile ] ")
	}

	return rConfig, nil
}

// CreateUserConfig creates user config from arguments
func CreateUserConfig(caller string, privKey string, publicKey string) (*UserConfigJSON, error) {
	userConfig := UserConfigJSON{PrivateKey: privKey, Caller: caller, PublicKey: publicKey}
	var err error

	ks := platformpolicy.NewKeyProcessor()
	userConfig.privateKeyObject, err = ks.ImportPrivateKeyPEM([]byte(privKey))
	return &userConfig, err
}
