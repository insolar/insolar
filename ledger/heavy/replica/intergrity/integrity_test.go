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

package intergrity

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar/secrets"
	"github.com/insolar/insolar/ledger/heavy/sequence"
)

func TestIntegrity_WrapUnwrap(t *testing.T) {
	var (
		expected []sequence.Item
	)
	f := fuzz.New().Funcs(func(item *sequence.Item, c fuzz.Continue) {
		c.Fuzz(&item.Key)
		c.Fuzz(&item.Value)
	})
	f.NumElements(3, 10).Fuzz(&expected)
	keys, err := secrets.GenerateKeyPair()
	require.NoError(t, err)
	cs := cryptography.NewKeyBoundCryptographyService(keys.Private)
	parentPubKey, err := cs.GetPublicKey()
	provider := NewProvider(cs)
	validator := NewValidator(cs, parentPubKey)

	packet := Wrap(expected)
	actual := UnwrapAndValidate(packet)
	require.Equal(t, expected, actual)
}
