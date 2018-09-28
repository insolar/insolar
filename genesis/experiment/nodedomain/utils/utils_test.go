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

package utils

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeDeserializePublicKey(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(GetCurve(), rand.Reader)
	if err != nil {
		panic(err)
	}

	serPubKey, err := SerializePublicKey(privateKey.PublicKey)
	assert.NoError(t, err)

	newPK, err := DeserializePublicKey(serPubKey)
	assert.NoError(t, err)
	assert.Equal(t, newPK, &privateKey.PublicKey)
}
