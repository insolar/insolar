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

package certificate

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/insolar/insolar/core"
	ecdsahelper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/pkg/errors"
)

// NewCertificate constructor creates new Certificate component
func NewCertificate(keysPath string) (*Certificate, error) {
	path := filepath.Clean(keysPath)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("couldn't read keys from: " + path)
	}
	var keys map[string]string
	err = json.Unmarshal(data, &keys)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse json.")
	}

	private, err := ecdsahelper.ImportPrivateKey(keys["private_key"])
	if err != nil {
		return nil, errors.Wrap(err, "Failed to import private key.")
	}

	valid, err := isValidPublicKey(keys["public_key"], private)
	if !valid {
		return nil, err
	}

	return &Certificate{privateKey: private}, nil
}

func isValidPublicKey(publicKey string, privateKey *ecdsa.PrivateKey) (bool, error) {
	validPublicKeyString, err := ecdsahelper.ExportPublicKey(&privateKey.PublicKey)
	if err != nil {
		return false, err
	} else if validPublicKeyString != publicKey {
		return false, errors.New("invalid public key in config")
	}
	return true, nil
}

// Certificate component
type Certificate struct {
	privateKey *ecdsa.PrivateKey
}

// Start is method from Component interface and it do nothing
func (c *Certificate) Start(components core.Components) error {
	return nil
}

// Stop is method from Component interface and it do nothing
func (c *Certificate) Stop() error {
	return nil
}

// GetPublicKey returns public key as string
func (c *Certificate) GetPublicKey() (string, error) {
	return ecdsahelper.ExportPublicKey(&c.privateKey.PublicKey)
}

// GetPrivateKey returns private key as string
func (c *Certificate) GetPrivateKey() (string, error) {
	return ecdsahelper.ExportPrivateKey(c.privateKey)
}

// GetEcdsaPrivateKey returns private key in ecdsa format
func (c *Certificate) GetEcdsaPrivateKey() *ecdsa.PrivateKey {
	return c.privateKey
}

// GenerateKeys generates certificate keys
func (c *Certificate) GenerateKeys() error {
	key, err := ecdsahelper.GeneratePrivateKey()
	if err != nil {
		return errors.Wrap(err, "Failed to generate certificate keys.")
	}

	c.privateKey = key
	return nil
}
