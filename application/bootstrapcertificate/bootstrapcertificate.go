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

package bootstrapcertificate

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strconv"

	ecdsahelper "github.com/insolar/insolar/cryptohelpers/ecdsa"

	"github.com/pkg/errors"
)

// Record contains info about node
type Record struct {
	NodeRef   string
	PublicKey string
}

// CertRecords is array od Records
type CertRecords = []Record

func serializeToJSON(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "    ")
}

type Certificate struct {
	CertRecords CertRecords `json:"nodes"`
	Signs       []string    `json:"signatures"`
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

func dumpRecords(cRecords CertRecords) ([]byte, error) {
	return serializeToJSON(cRecords)
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

func (cr *Certificate) Dump() (string, error) {
	result, err := serializeToJSON(cr)
	if err != nil {
		return "", errors.Wrap(err, "[ Certificate::Dump ]")
	}

	return string(result), nil
}
