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
	"github.com/insolar/insolar/testutils/launchnet"

	"github.com/stretchr/testify/require"
)

func TestMemberGet(t *testing.T) {
	member1 := *createMember(t)
	member2, _ := newUserWithKeys()
	member2.PubKey = member1.PubKey
	member2.PrivKey = member1.PrivKey
	res, err := signedRequest(t, launchnet.TestRPCUrlPublic, member2, "member.get", nil)
	require.Nil(t, err)
	require.Equal(t, member1.Ref, res.(map[string]interface{})["reference"].(string))
}

func TestMigrationMemberGet(t *testing.T) {
	member1, _ := newUserWithKeys()

	ba := testutils.RandomString()
	_, _ = signedRequest(t, launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.addAddresses", map[string]interface{}{"migrationAddresses": []string{ba}})

	res1, err := signedRequest(t, launchnet.TestRPCUrlPublic, member1, "member.migrationCreate", nil)
	require.Nil(t, err)

	decodedRes1, ok := res1.(map[string]interface{})
	require.True(t, ok, fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res1))

	res2, err := signedRequest(t, launchnet.TestRPCUrlPublic, member1, "member.get", nil)
	require.Nil(t, err)

	decodedRes2, ok := res2.(map[string]interface{})
	require.True(t, ok, fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))

	require.Equal(t, decodedRes1["reference"].(string), decodedRes2["reference"].(string))
	require.Equal(t, ba, res2.(map[string]interface{})["migrationAddress"].(string))
}

func TestMemberGetWrongPublicKey(t *testing.T) {
	member1, _ := newUserWithKeys()
	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member1, "member.get", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get reference by public key: failed to get reference in shard: failed to find reference by key")
}

func TestMemberGetGenesisMember(t *testing.T) {
	res, err := signedRequest(t, launchnet.TestRPCUrlPublic, &launchnet.MigrationAdmin, "member.get", nil)
	require.Nil(t, err)
	require.Equal(t, launchnet.MigrationAdmin.Ref, res.(map[string]interface{})["reference"].(string))
}
