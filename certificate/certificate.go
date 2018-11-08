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

	ecdsahelper "github.com/insolar/insolar/cryptoproviders/ecdsa"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
)

// BootstrapNode holds info about bootstrap nodes
type BootstrapNode struct {
	PublicKey string `json:"public_key"`
	Host      string `json:"host"`
}

// Certificate holds info about certificate
type Certificate struct {
	MajorityRule   int             `json:"majority_rule"`
	PublicKey      string          `json:"public_key"`
	Reference      string          `json:"reference"`
	Roles          []string        `json:"roles"`
	BootstrapNodes []BootstrapNode `json:"bootstrap_nodes"`

	privateKey *ecdsa.PrivateKey
}

func AreKeysTheSame(privateKey *ecdsa.PrivateKey, certPubKey string) error {
	pubKeyString, err := ecdsahelper.ExportPublicKey(&privateKey.PublicKey)
	if err != nil {
		return errors.Wrap(err, "[ AreKeysTheSame ]")
	}
	if pubKeyString != certPubKey {
		msg := "[ AreKeysTheSame ] Public keys in certificate and keypath file are not the same: " +
			"pubKeyString = " + pubKeyString +
			"\ncertPubKey = " + certPubKey
		return errors.New(msg)
	}
	return nil
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

	private, err := readPrivateKey(keysPath)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewCertificate ] failed to read keys")
	}

	err = AreKeysTheSame(private, cert.PublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewCertificate ] Different public keys. Cert path: "+certPath+". Key path: "+keysPath)
	}

	cert.privateKey = private
	return &cert, nil
}

func (cert *Certificate) reset() {
	cert.PublicKey = ""
	cert.BootstrapNodes = []BootstrapNode{}
	cert.privateKey = nil
	cert.Reference = ""
	cert.MajorityRule = 0
	cert.Roles = []string{}
}

// NewCertificatesWithKeys generate certificate from given keys
func NewCertificatesWithKeys(keysPath string) (*Certificate, error) {
	cert := Certificate{}
	cert.reset()
	private, err := readPrivateKey(keysPath)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewCertificatesWithKeys ] failed to read keys")
	}

	err = cert.setKeys(private)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewCertificatesWithKeys ] Problem with setting keys.")
	}

	cert.Reference = testutils.RandomRef().String()

	return &cert, nil
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

func readPrivateKey(keysPath string) (*ecdsa.PrivateKey, error) {
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

	err = isValidPublicKey(keys["public_key"], private)
	if err != nil {
		return nil, errors.Wrap(err, "[ readKeys ] public key is not valid")
	}

	return private, nil
}

func isValidPublicKey(publicKey string, privateKey *ecdsa.PrivateKey) error {
	validPublicKeyString, err := ecdsahelper.ExportPublicKey(&privateKey.PublicKey)
	if err != nil {
		return errors.Wrap(err, "[ isValidPublicKey ]")
	}
	if validPublicKeyString != publicKey {
		return errors.New("[ isValidPublicKey ] invalid public key in config")
	}
	return nil
}

func (cert *Certificate) setKeys(privateKey *ecdsa.PrivateKey) error {
	expPubKey, err := ecdsahelper.ExportPublicKey(&privateKey.PublicKey)
	if err != nil {
		return errors.Wrap(err, "[ GenerateKeys ] Failed to export public key.")
	}

	cert.PublicKey = expPubKey
	cert.privateKey = privateKey

	return nil
}

// GenerateKeys generates certificate keys
func (cert *Certificate) GenerateKeys() error {
	privateKey, err := ecdsahelper.GeneratePrivateKey()
	if err != nil {
		return errors.Wrap(err, "[ GenerateKeys ] Failed to generate private key.")
	}

	err = cert.setKeys(privateKey)
	if err != nil {
		return errors.Wrap(err, "[ GenerateKeys ] Problem with setting keys.")
	}

	return nil
}

func (cert *Certificate) Dump() (string, error) {
	result, err := json.MarshalIndent(cert, "", "    ")
	if err != nil {
		return "", errors.Wrap(err, "[ Certificate::Dump ]")
	}

	return string(result), nil
}
