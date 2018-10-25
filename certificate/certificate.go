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
	"strconv"

	"github.com/insolar/insolar/core"
	ecdsahelper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/pkg/errors"
)

// NewCertificate constructor creates new Certificate component
func NewCertificate(keysPath string) (*Certificate, error) {
	data, err := ioutil.ReadFile(filepath.Clean(keysPath))
	if err != nil {
		return nil, errors.New("couldn't read keys from: " + keysPath)
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

// Record contains info about node
type Record struct {
	NodeRef   string
	PublicKey string
}

// CertRecords is array od Records
type CertRecords = []Record

// Certificate component
type Certificate struct {
	CertRecords CertRecords `json:"nodes"`
	Signs       []string    `json:"signatures"`

	privateKey *ecdsa.PrivateKey
}

// Start is method from Component interface and it do nothing
func (c *Certificate) Start(ctx core.Context, components core.Components) error {
	return nil
}

// Stop is method from Component interface and it do nothing
func (c *Certificate) Stop(ctx core.Context) error {
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

func NewCertificateFromFile(path string) (*Certificate, error) {
	data, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrap(err, "[ NewCertificateFromFile ]")
	}
	cert := Certificate{}
	err = json.Unmarshal(data, &cert)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewCertificateFromFile ]")
	}

	err = cert.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "[ NewCertificateFromFile ]")
	}

	return &cert, nil
}

// NewCertificateFromFields creates new Certificate from prefilled fields
func NewCertificateFromFields(cRecords CertRecords, keys []*ecdsa.PrivateKey) (*Certificate, error) {
	if len(cRecords) != len(keys) {
		return nil, errors.New("[ NewCertificateFromFields ] params must be the same length")
	}
	if len(cRecords) == 0 {
		return nil, errors.New("[ NewCertificateFromFields ] params must not be empty")
	}
	certData, err := dumpRecords(cRecords)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewCertificateFromFields ]")
	}

	var signList []string
	for _, k := range keys {
		sign, err := ecdsahelper.Sign(certData, k)
		if err != nil {
			return nil, errors.Wrap(err, "[ NewCertificateFromFields ]")
		}

		signList = append(signList, ecdsahelper.ExportSignature(sign))
	}

	return &Certificate{
		CertRecords: cRecords,
		Signs:       signList,
	}, nil

}

func (cr *Certificate) Validate() error {
	if len(cr.Signs) != len(cr.CertRecords) {
		return errors.New("[ Validate ] Wrong number of nodes and signatures")
	}

	if len(cr.Signs) == 0 {
		return errors.New("[ Validate ] Empty fields")
	}

	size := len(cr.Signs)
	certData, err := dumpRecords(cr.CertRecords)
	if err != nil {
		return errors.Wrap(err, "[ Validate ]")
	}
	for i := 0; i < size; i++ {
		sign, err := ecdsahelper.ImportSignature(cr.Signs[i])
		if err != nil {
			return errors.Wrap(err, "[ Validate ]")
		}
		ok, err := ecdsahelper.Verify(certData, sign, cr.CertRecords[i].PublicKey)
		if err != nil {
			return errors.Wrap(err, "[ Validate ]")
		}
		if !ok {
			return errors.New("[ Validate ] invalid signature: " + strconv.Itoa(i))
		}
	}

	return nil
}

func serializeToJSON(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "    ")
}

func dumpRecords(cRecords CertRecords) ([]byte, error) {
	return serializeToJSON(cRecords)
}

func (cr *Certificate) Dump() (string, error) {
	result, err := serializeToJSON(cr)
	if err != nil {
		return "", errors.Wrap(err, "[ Certificate::Dump ]")
	}

	return string(result), nil
}
