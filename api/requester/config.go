// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
		privKey, err := secrets.GeneratePrivateKeyEthereum()
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
