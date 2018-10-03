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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCertificateFromFile(t *testing.T) {
	cert, err := NewCertificateFromFile("testdata/cert.json")
	assert.NoError(t, err)
	ok, err := cert.Validate()
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestNewCertificateFromFile_WrongSignature(t *testing.T) {
	_, err := NewCertificateFromFile("testdata/cert_wrong_signature.json")
	assert.EqualError(t, err, "[ NewCertificateFromFile ]: [ Validate ] invalid signature: 0")
}

func TestNewCertificateFromFields(t *testing.T) {
	NewCertificateFromFields(nil, nil)
}
