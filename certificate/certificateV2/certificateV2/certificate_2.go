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

package certificateV2

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/insolar/insolar/core"
	ecdsahelper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/pkg/errors"
)

type BootstrapNode struct {
	PublicKey string `json:"public_key"`
	Host      string `json:"host"`
}

type Certificate struct {
	MajorityRule   int             `json:"majority_rule"`
	PublicKey      string          `json:"public_key"`
	Reference      string          `json:"reference"`
	Roles          []string        `json:"roles"`
	BootstrapNodes []BootstrapNode `json:"bootstrap_nodes"`

	privateKey *ecdsa.PrivateKey
}

// GetPublicKey returns public key as string
func (cert *Certificate) GetPublicKey() (string, error) {
	return ecdsahelper.ExportPublicKey(&cert.privateKey.PublicKey)
}

// GetPrivateKey returns private key as string
func (cert *Certificate) GetPrivateKey() (string, error) {
	return ecdsahelper.ExportPrivateKey(cert.privateKey)
}

// GetEcdsaPrivateKey returns private key in ecdsa format
func (cert *Certificate) GetEcdsaPrivateKey() *ecdsa.PrivateKey {
	return cert.privateKey
}

// Start is method from Component interface and it do nothing
func (cert *Certificate) Start(components core.Components) error {
	return nil
}

// Stop is method from Component interface and it do nothing
func (cert *Certificate) Stop() error {
	return nil
}

func readKeys(keysPath string, certPublicKey string) (*ecdsa.PrivateKey, error) {
	data, err := ioutil.ReadFile(filepath.Clean(keysPath))
	if err != nil {
		return nil, errors.Wrap(err, "[ readKeys ] couldn't read keys from: "+keysPath)
	}
	var keys map[string]string
	err = json.Unmarshal(data, &keys)
	if err != nil {
		return nil, errors.Wrap(err, "[ readKeys ] failed to parse json.")
	}

	private, err := ecdsahelper.ImportPrivateKey(keys["private_key"])
	if err != nil {
		return nil, errors.Wrap(err, "[ readKeys ] Failed to import private key.")
	}

	if keys["public_key"] != certPublicKey {
		return nil, errors.New("[ readKeys ] Public keys in certificate and keypath file are not the same")
	}

	valid, err := isValidPublicKey(keys["public_key"], private)
	if !valid {
		return nil, errors.Wrap(err, "[ readKeys ] public key is not valid")
	}

	return private, nil
}

// NewCertificate constructor creates new Certificate component
func NewCertificate(keysPath string, certPath string) (*Certificate, error) {
	data, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return nil, errors.New("[ NewCertificate ] couldn't read certificate from: " + certPath)
	}
	cert := Certificate{}
	err = json.Unmarshal(data, &cert)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewCertificate ] failed to parse certificate json")
	}

	private, err := readKeys(keysPath, cert.PublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewCertificate ] failed to read keys")
	}
	cert.privateKey = private
	return &cert, nil
}

func isValidPublicKey(publicKey string, privateKey *ecdsa.PrivateKey) (bool, error) {
	validPublicKeyString, err := ecdsahelper.ExportPublicKey(&privateKey.PublicKey)
	if err != nil {
		return false, errors.Wrap(err, "[ isValidPublicKey ]")
	} else if validPublicKeyString != publicKey {
		return false, errors.New("[ isValidPublicKey ] invalid public key in config")
	}
	return true, nil
}
