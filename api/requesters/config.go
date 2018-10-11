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

package requesters

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	ecdsa_helper "github.com/insolar/insolar/cryptohelpers/ecdsa"

	"github.com/pkg/errors"
)

// UserConfigJSON holds info about user
type UserConfigJSON struct {
	PrivateKey       string `json:"private_key"`
	privateKeyObject *ecdsa.PrivateKey
}

// RequestConfigJSON holds info about request
type RequestConfigJSON struct {
	Params    []interface{} `json:"params"`
	Method    string        `json:"method"`
	Requester string        `json:"requester"`
	Target    string        `json:"target"`
}

func readFile(path string, configType interface{}) error {
	rawConf, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return errors.Wrap(err, "[ readFile ] Problem with reading config")
	}

	err = json.Unmarshal(rawConf, &configType)
	if err != nil {
		return errors.Wrap(err, "[ readFile ] Problem with unmarshaling config")
	}

	return nil
}

// ReadUserConfigFromFile read user confgi from file
func ReadUserConfigFromFile(path string) (*UserConfigJSON, error) {
	cfgJSON := &UserConfigJSON{}
	err := readFile(path, cfgJSON)
	if err != nil {
		return nil, errors.Wrap(err, "[ readUserConfigFromFile ] ")
	}

	cfgJSON.privateKeyObject, err = ecdsa_helper.ImportPrivateKey(cfgJSON.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ readUserConfigFromFile ] Problem with reading private key")
	}

	return cfgJSON, nil
}

// ReadRequesterConfigFromFile read request config from file
func ReadRequestConfigFromFile(path string) (*RequestConfigJSON, error) {
	rConfig := &RequestConfigJSON{}
	err := readFile(path, rConfig)
	if err != nil {
		return nil, errors.Wrap(err, "[ readRequesterConfigFromFile ] ")
	}

	return rConfig, nil
}
