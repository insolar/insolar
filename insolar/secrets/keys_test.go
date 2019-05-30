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

package secrets

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeys_GetKeysFromFile(t *testing.T) {
	pair, err := ReadKeysFile("testdata/keypair.json")
	require.NoError(t, err, "read keys from json")
	assert.Equal(t, fmt.Sprintf("%T", pair.Private), "*ecdsa.PrivateKey", "private key has proper type")
	assert.Equal(t, fmt.Sprintf("%T", pair.Public), "*ecdsa.PublicKey", "public key has proper type")
}
