///
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
///

// +build functest

package functest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/launchnet"
)

func TestMemberMigrationCreate(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	ba := testutils.RandomString()
	_, err = signedRequest(t, launchnet.TestRPCUrl, &launchnet.MigrationAdmin,
		"migration.addAddresses", map[string]interface{}{"migrationAddresses": []string{ba}})
	require.NoError(t, err)
	result, err := signedRequest(t, launchnet.TestRPCUrlPublic, member, "member.migrationCreate", nil)
	require.NoError(t, err)
	output, ok := result.(map[string]interface{})
	require.True(t, ok)
	require.NotEqual(t, "", output["reference"])
	require.Equal(t, ba, output["migrationAddress"])
}

func TestMemberMigrationCreateWhenNomigrationAddressesLeft(t *testing.T) {
	member1, err := newUserWithKeys()
	require.NoError(t, err)
	addMigrationAddress(t)
	_, err = signedRequest(t, launchnet.TestRPCUrlPublic, member1, "member.migrationCreate", nil)
	require.Nil(t, err)

	member2, err := newUserWithKeys()
	require.NoError(t, err)

	_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member2, "member.migrationCreate", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no more migration addresses left in any shard")
}

func TestMemberMigrationCreateWithBadKey(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.PubKey = "fake"
	_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member, "member.migrationCreate", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("problems with decoding. Key - %s", member.PubKey))
}

func TestMemberMigrationCreateWithSamePublicKey(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)

	addMigrationAddress(t)

	_, err = signedRequest(t, launchnet.TestRPCUrlPublic, member, "member.migrationCreate", nil)
	require.NoError(t, err)

	addMigrationAddress(t)

	_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member, "member.migrationCreate", map[string]interface{}{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to set reference in public key shard: can't set reference because this key already exists")

	memberForBurn, err := newUserWithKeys()
	require.NoError(t, err)

	_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, memberForBurn, "member.migrationCreate", nil)
}

func TestMemberMigrationCreateWithSameMigrationAddress(t *testing.T) {
	member1, err := newUserWithKeys()
	require.NoError(t, err)

	ba := testutils.RandomString()
	_, _ = signedRequest(t, launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.addAddresses",
		map[string]interface{}{"migrationAddresses": []string{ba, ba}})

	_, err = signedRequest(t, launchnet.TestRPCUrlPublic, member1, "member.migrationCreate", nil)
	require.NoError(t, err)

	member2, err := newUserWithKeys()
	require.NoError(t, err)

	_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member2, "member.migrationCreate", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to set reference in migration address shard: can't set reference because this key already exists")
}
