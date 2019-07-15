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
	"testing"

	"github.com/insolar/insolar/testutils"

	"github.com/stretchr/testify/require"
)

func TestMemberGet(t *testing.T) {
	member1 := *createMember(t)
	member2 := member1
	member2.ref = root.ref
	res, err := signedRequest(&member2, "member.get", nil)
	require.Nil(t, err)
	require.Equal(t, member1.ref, res.(map[string]interface{})["reference"].(string))
}

func TestMigrationMemberGet(t *testing.T) {
	member1, _ := newUserWithKeys()
	member1.ref = root.ref

	ba := testutils.RandomString()
	_, _ = signedRequest(&migrationAdmin, "migration.addBurnAddresses", map[string]interface{}{"burnAddresses": []string{ba}})

	res1, err := retryableMemberMigrationCreate(member1, true)

	member2 := *member1
	member2.ref = root.ref
	res2, err := signedRequest(&member2, "member.get", nil)
	require.Nil(t, err)
	require.Equal(t, res1.(map[string]interface{})["reference"].(string), res2.(map[string]interface{})["reference"].(string))
	require.Equal(t, ba, res2.(map[string]interface{})["migrationAddress"].(string))
}

func TestMemberGetWrongPublicKey(t *testing.T) {
	member1, _ := newUserWithKeys()
	member1.ref = root.ref
	_, err := signedRequest(member1, "member.get", nil)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "member for this public key does not exist")
}
