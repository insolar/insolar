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

// +build functest

package functest

import (
	"fmt"
	"testing"

	"github.com/insolar/insolar/testutils"

	"github.com/stretchr/testify/require"
)

func TestMemberGet(t *testing.T) {
	member1 := *createMember(t)
	member2, _ := newUserWithKeys()
	member2.pubKey = member1.pubKey
	member2.privKey = member1.privKey
	res, err := signedRequest(member2, "member.get", nil)
	require.Nil(t, err)
	require.Equal(t, member1.ref, res.(map[string]interface{})["reference"].(string))
}

func TestMigrationMemberGet(t *testing.T) {
	member1, _ := newUserWithKeys()

	ba := testutils.RandomString()
	_, _ = signedRequest(&migrationAdmin, "migration.addBurnAddresses", map[string]interface{}{"burnAddresses": []string{ba}})

	res1, err := retryableMemberMigrationCreate(member1, true)
	require.Nil(t, err)

	decodedRes1, ok := res1.(map[string]interface{})
	require.True(t, ok, fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res1))

	res2, err := signedRequest(member1, "member.get", nil)
	require.Nil(t, err)

	decodedRes2, ok := res2.(map[string]interface{})
	require.True(t, ok, fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))

	require.Equal(t, decodedRes1["reference"].(string), decodedRes2["reference"].(string))
	require.Equal(t, ba, res2.(map[string]interface{})["migrationAddress"].(string))
}

func TestMemberGetWrongPublicKey(t *testing.T) {
	member1, _ := newUserWithKeys()
	_, err := signedRequest(member1, "member.get", nil)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "member for this public key does not exist")
}
