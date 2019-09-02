//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package internal

import (
	"testing"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/platformpolicy/internal/hash"
	"github.com/insolar/insolar/platformpolicy/internal/sign"
	"github.com/stretchr/testify/assert"
)

func TestEcdsaMarshalUnmarshal(t *testing.T) {
	data := gen.Reference()

	kp := platformpolicy.NewKeyProcessor()
	provider := sign.NewECDSAProvider()

	hasher := hash.NewSHA3Provider().Hash512bits()

	privateKey, err := kp.GeneratePrivateKey()
	assert.NoError(t, err)

	signer := provider.DataSigner(privateKey, hasher)
	verifier := provider.DataVerifier(kp.ExtractPublicKey(privateKey), hasher)

	signature, err := signer.Sign(data.Bytes())
	assert.NoError(t, err)

	assert.True(t, verifier.Verify(*signature, data.Bytes()))
}
