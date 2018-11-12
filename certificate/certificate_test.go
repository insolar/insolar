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

//
// import (
// 	"testing"
//
// 	ecdsahelper "github.com/insolar/insolar/cryptohelpers/ecdsa"
// 	"github.com/stretchr/testify/assert"
// )
//
// const TEST_CERT = "testdata/cert.json"
// const TEST_BAD_CERT = "testdata/bad_cert.json"
//
// const TEST_KEYS = "testdata/keys.json"
// const TEST_BAD_KEYS = "testdata/bad_keys.json"
//
// func TestAreKeysTheSame(t *testing.T) {
// 	privateKey, err := ecdsahelper.GeneratePrivateKey()
// 	assert.NoError(t, err)
// 	pubKey, err := ecdsahelper.ExportPublicKey(&privateKey.PublicKey)
// 	assert.NoError(t, err)
// 	assert.NoError(t, AreKeysTheSame(privateKey, pubKey))
// }
//
// func TestAreKeysTheSame_NotTheSame(t *testing.T) {
// 	privateKey, err := ecdsahelper.GeneratePrivateKey()
// 	assert.NoError(t, err)
// 	err = AreKeysTheSame(privateKey, "")
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "Public keys in certificate and keypath file are not the same")
// }
//
// func TestNewCertificate_NoCert(t *testing.T) {
// 	_, err := ReadCertificate("", "")
// 	assert.EqualError(t, err, "[ ReadCertificate ] couldn't read certificate from: ")
// }
//
// func TestNewCertificate_BadCert(t *testing.T) {
// 	_, err := ReadCertificate("", TEST_BAD_CERT)
// 	assert.Contains(t, err.Error(), "failed to parse certificate json")
// }
//
// func TestNewCertificate_NoKeys(t *testing.T) {
// 	_, err := ReadCertificate("", TEST_CERT)
// 	assert.Contains(t, err.Error(), "failed to read keys")
// }
//
// func checkKeys(cert *Certificate, t *testing.T) {
// 	pubKey, err := ecdsahelper.ExportPublicKey(&cert.privateKey.PublicKey)
// 	assert.NoError(t, err)
// 	assert.Equal(t, pubKey, cert.PublicKey)
// }
//
// func TestNewCertificate(t *testing.T) {
// 	cert, err := ReadCertificate(TEST_KEYS, TEST_CERT)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, cert.PublicKey)
// 	assert.NotEmpty(t, cert.Reference)
//
// 	checkKeys(cert, t)
// }
//
// func TestCertificate_GenerateKeys(t *testing.T) {
// 	cert := Certificate{}
// 	assert.Nil(t, cert.privateKey)
// 	assert.Empty(t, cert.PublicKey)
//
// 	assert.NoError(t, cert.GenerateKeys())
//
// 	assert.NotNil(t, cert.privateKey)
// 	assert.NotEmpty(t, cert.PublicKey)
// }
//
// func TestNewCertificatesWithKeys(t *testing.T) {
// 	cert, err := NewCertificatesWithKeys(TEST_KEYS)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, cert.Reference)
// 	checkKeys(cert, t)
// }
//
// func TestNewCertificatesWithKeys_NoFile(t *testing.T) {
// 	_, err := NewCertificatesWithKeys("")
// 	assert.Contains(t, err.Error(), "failed to read keys: [ readKeys ] couldn't read keys from")
// }
//
// func TestReadPrivateKey(t *testing.T) {
// 	_, err := readPrivateKey("")
// 	assert.Contains(t, err.Error(), "couldn't read keys from")
// }
//
// func TestReadPrivateKey_BadJson(t *testing.T) {
// 	_, err := readPrivateKey(TEST_BAD_CERT)
// 	assert.Contains(t, err.Error(), "failed to parse json")
// }
//
// func TestReadPrivateKey_BadPrivateKey(t *testing.T) {
// 	_, err := readPrivateKey(TEST_BAD_KEYS)
// 	assert.Contains(t, err.Error(), "Failed to import private key")
// }
//
// func TestReadPrivateKey_BadKeyPair(t *testing.T) {
// 	_, err := readPrivateKey("testdata/different_keys.json")
// 	assert.Contains(t, err.Error(), "public key is not valid")
// }
//
// func TestIsPublicKeyValid(t *testing.T) {
// 	privateKey, err := ecdsahelper.GeneratePrivateKey()
// 	assert.NoError(t, err)
// 	pubKey, err := ecdsahelper.ExportPublicKey(&privateKey.PublicKey)
// 	assert.NoError(t, err)
//
// 	assert.Nil(t, isValidPublicKey(pubKey, privateKey))
// }
//
// func TestIsPublicKeyValid_BadKeyPair(t *testing.T) {
// 	privateKey, err := ecdsahelper.GeneratePrivateKey()
// 	assert.NoError(t, err)
// 	pubKey, err := ecdsahelper.ExportPublicKey(&privateKey.PublicKey)
// 	assert.NoError(t, err)
//
// 	anotherPrivateKey, err := ecdsahelper.GeneratePrivateKey()
// 	assert.NoError(t, err)
//
// 	err = isValidPublicKey(pubKey, anotherPrivateKey)
// 	assert.Contains(t, err.Error(), "[ isValidPublicKey ] invalid public key in config")
//
// }
