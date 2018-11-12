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
	"crypto"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/insolar/insolar/platformpolicy"
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
}

// NewCertificate constructor creates new Certificate component
func NewCertificate(publicKey crypto.PublicKey, certPath string) (*Certificate, error) {
	data, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return nil, errors.New("[ NewCertificate ] couldn't read certificate from: " + certPath)
	}
	cert := Certificate{}
	err = json.Unmarshal(data, &cert)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewCertificate ] failed to parse certificate json")
	}

	keyProcessor := platformpolicy.NewKeyProcessor()
	pub, err := keyProcessor.ExportPublicKey(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewCertificate ] failed to retrieve public key from node private key")
	}

	if cert.PublicKey != string(pub) {
		return nil, errors.Wrap(err, "[ NewCertificate ] Different public keys. Cert path: "+certPath+".")
	}

	return &cert, nil
}

func (cert *Certificate) reset() {
	cert.PublicKey = ""
	cert.BootstrapNodes = []BootstrapNode{}
	cert.Reference = ""
	cert.MajorityRule = 0
	cert.Roles = []string{}
}

// NewCertificatesWithKeys generate certificate from given keys
func NewCertificatesWithKeys(keysPath string) (*Certificate, error) {
	cert := Certificate{}
	cert.reset()

	cert.Reference = testutils.RandomRef().String()

	return &cert, nil
}

// GenerateKeys generates certificate keys
func (cert *Certificate) GenerateKeys() error {
	// keyProcessor := platformpolicy.NewKeyProcessor()
	// privateKey, err := keyProcessor.GeneratePrivateKey()
	// if err != nil {
	// 	return errors.Wrap(err, "[ GenerateKeys ] Failed to generate private key.")
	// }
	//
	// err = cert.setKeys(privateKey)
	// if err != nil {
	// 	return errors.Wrap(err, "[ GenerateKeys ] Problem with setting keys.")
	// }

	return nil
}

func (cert *Certificate) Dump() (string, error) {
	result, err := json.MarshalIndent(cert, "", "    ")
	if err != nil {
		return "", errors.Wrap(err, "[ Certificate::Dump ]")
	}

	return string(result), nil
}
