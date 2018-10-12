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
	"os"
	"testing"

	"github.com/insolar/insolar/core"
	ecdsahelper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/stretchr/testify/assert"
)

type key struct {
	Private_key string `json:"private_key"`
	Public_key  string `json:"public_key"`
}

func TestNewCertificateFromFile(t *testing.T) {
	cert, err := NewCertificateFromFile("testdata/cert.json")
	assert.NoError(t, err)
	err = cert.Validate()
	assert.NoError(t, err)
}

func TestNewCertificateFromFields_DifferentFieldsLength(t *testing.T) {
	privKeys := []*ecdsa.PrivateKey{}
	privKeys = append(privKeys, nil)
	privKeys = append(privKeys, nil)

	cRecords := CertRecords{}

	_, err := NewCertificateFromFields(cRecords, privKeys)
	assert.EqualError(t, err, "[ NewCertificateFromFields ] params must be the same length")
}

func TestNewCertificateFromFields(t *testing.T) {
	privKeys := []*ecdsa.PrivateKey{}
	records := CertRecords{}
	for i := 0; i < 10; i++ {
		key, err := ecdsahelper.GeneratePrivateKey()
		assert.NoError(t, err)
		privKeys = append(privKeys, key)

		pubKey, err := ecdsahelper.ExportPublicKey(&key.PublicKey)
		assert.NoError(t, err)
		records = append(records, Record{NodeRef: core.RandomRef().String(), PublicKey: pubKey})
	}

	cert, err := NewCertificateFromFields(records, privKeys)
	assert.NoError(t, err)
	assert.NoError(t, cert.Validate())

}

func TestNewCertifiacteFromFile_BadFile(t *testing.T) {
	_, err := NewCertificateFromFile("______")
	assert.EqualError(t, err, "[ NewCertificateFromFile ]: open ______: no such file or directory")
}

func TestNewCertifiacteFromFile_BadFileData(t *testing.T) {
	_, err := NewCertificateFromFile("testdata/private_keys.json")
	assert.EqualError(t, err, "[ NewCertificateFromFile ]: json: cannot unmarshal array into Go value of type bootstrapcertificate.Certificate")
}

func TestNewCertificateFromFile_WrongSignature(t *testing.T) {
	_, err := NewCertificateFromFile("testdata/cert_wrong_signature.json")
	assert.EqualError(t, err, "[ NewCertificateFromFile ]: [ Validate ] invalid signature: 0")
}

func TestNewCertificateFromFields_EmptyLists(t *testing.T) {
	_, err := NewCertificateFromFields(nil, nil)
	assert.EqualError(t, err, "[ NewCertificateFromFields ] params must not be empty")
}

func readPrivateKeys() ([]*ecdsa.PrivateKey, error) {
	rawKeys, err := ioutil.ReadFile("testdata/private_keys.json")
	if err != nil {
		return nil, err
	}

	keysData := []key{}
	err = json.Unmarshal(rawKeys, &keysData)
	if err != nil {
		return nil, err
	}

	privateKeys := []*ecdsa.PrivateKey{}

	for i := 0; i < len(keysData); i++ {
		privKey, err := ecdsahelper.ImportPrivateKey(keysData[i].Private_key)
		if err != nil {
			return nil, err
		}
		privateKeys = append(privateKeys, privKey)
	}

	return privateKeys, nil
}

func TestComplexCheck(t *testing.T) {
	// Read from test dump
	cert, err := NewCertificateFromFile("testdata/cert.json")
	assert.NoError(t, err)

	// Dump to tmp file
	dumpCert, err := cert.Dump()
	assert.NoError(t, err)
	tmpDir, err := ioutil.TempDir("", "test-")
	defer os.RemoveAll(tmpDir)
	assert.NoError(t, err)
	tmpFile := tmpDir + "/test_cert.json"
	ioutil.WriteFile(tmpFile, []byte(dumpCert), 0644)

	// Read from tmp file
	newCert, err := NewCertificateFromFile(tmpFile)
	assert.NoError(t, err)
	assert.Equal(t, cert, newCert)

	// Construct from fields
	privateKeys, err := readPrivateKeys()
	assert.NoError(t, err)

	new2Cert, err := NewCertificateFromFields(cert.CertRecords, privateKeys)
	assert.NoError(t, err)
	assert.Equal(t, newCert.CertRecords, new2Cert.CertRecords)

	err = cert.Validate()
	assert.NoError(t, err)

}
