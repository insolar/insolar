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

package internal

import (
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/platformpolicy/internal/hash"
	"github.com/insolar/insolar/platformpolicy/internal/sign"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func TestEcdsaMarshalUnmarshal(t *testing.T) {
	count := 10000
	data := testutils.RandomRef()

	kp := platformpolicy.NewKeyProcessor()
	provider := sign.NewECDSAProvider()

	cm := component.Manager{}
	cm.Inject(provider, hash.NewSHA3Provider())

	for i := 0; i < count; i++ {
		privateKey, err := kp.GeneratePrivateKey()
		assert.NoError(t, err)

		signer := provider.Sign(privateKey)
		verifier := provider.Verify(kp.ExtractPublicKey(privateKey))

		signature, err := signer.Sign(data.Bytes())
		assert.NoError(t, err)

		assert.True(t, verifier.Verify(*signature, data.Bytes()))
	}
}
