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

package pulsar

import (
	"testing"

	ecdsa_helper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/stretchr/testify/assert"
)

func TestSingAndVerify(t *testing.T) {
	assertObj := assert.New(t)
	privateKey, err := ecdsa_helper.GeneratePrivateKey()
	assert.NoError(t, err)
	publicKey, err := ecdsa_helper.ExportPublicKey(&privateKey.PublicKey)
	assert.NoError(t, err)

	signature, err := singData(privateKey, "This is the message to be signed and verified!")
	assertObj.NoError(err)

	checkSignature, err := checkPayloadSignature(&Payload{PublicKey: publicKey, Signature: signature, Body: "This is the message to be signed and verified!"})

	assertObj.Equal(true, checkSignature)
}
