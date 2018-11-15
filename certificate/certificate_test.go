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
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestCert = "testdata/cert.json"
const TestBadCert = "testdata/bad_cert.json"
const TestInvalidFileCert = "testdata/bad_cert11111.json"

const TestKeys = "testdata/keys.json"
const TestDifferentKeys = "testdata/different_keys.json"

func TestNewCertificate_NoCert(t *testing.T) {
	_, err := ReadCertificate(nil, nil, TestInvalidFileCert)
	assert.EqualError(t, err, "[ ReadCertificate ] failed to read certificate from: "+TestInvalidFileCert)
}

func TestNewCertificate_BadCert(t *testing.T) {
	_, err := ReadCertificate(nil, nil, TestBadCert)
	assert.Contains(t, err.Error(), "failed to parse certificate json")
}

func checkKeys(cert *Certificate, cs core.CryptographyService, t *testing.T) {
	kp := platformpolicy.NewKeyProcessor()

	pubKey, err := cs.GetPublicKey()
	assert.NoError(t, err)

	pubKeyString, err := kp.ExportPublicKey(pubKey)
	assert.NoError(t, err)

	assert.Equal(t, string(pubKeyString), cert.PublicKey)
}

func TestReadCertificate(t *testing.T) {
	cs, _ := cryptography.NewStorageBoundCryptographyService(TestKeys)
	kp := platformpolicy.NewKeyProcessor()
	pk, _ := cs.GetPublicKey()

	cert, err := ReadCertificate(pk, kp, TestCert)
	require.NoError(t, err)
	require.NotEmpty(t, cert.PublicKey)
	require.Equal(t, "virtual", cert.Role)
	require.Equal(t, "2prKtCG51YhseciDY5EnnHapPskNHvrvhSc3HrCvYLKKxXn4K3kFQtiz3QLVD1acpQmaDBHUG2Q988xjSFhswJLs",
		cert.Reference)
	require.Equal(t, 7, cert.MajorityRule)

	bootstrapNodes := []BootstrapNode{
		BootstrapNode{
			PublicKey: "PUBKEY_1",
			Host:      "localhost:22001",
		},
		BootstrapNode{
			PublicKey: "PUBKEY_2",
			Host:      "localhost:22002",
		},
		BootstrapNode{
			PublicKey: "PUBKEY_3",
			Host:      "localhost:22003",
		},
	}
	require.Equal(t, bootstrapNodes, cert.BootstrapNodes)

	checkKeys(cert, cs, t)
}

func TestReadPrivateKey_BadJson(t *testing.T) {
	keyProcessor := platformpolicy.NewKeyProcessor()
	_, err := ReadCertificate(nil, keyProcessor, TestBadCert)
	assert.Contains(t, err.Error(), "failed to parse certificate json")
}

func TestReadPrivateKey_BadKeyPair(t *testing.T) {
	cs, _ := cryptography.NewStorageBoundCryptographyService(TestDifferentKeys)
	kp := platformpolicy.NewKeyProcessor()
	pk, _ := cs.GetPublicKey()

	_, err := ReadCertificate(pk, kp, TestCert)
	assert.Contains(t, err.Error(), "Different public keys")
}
