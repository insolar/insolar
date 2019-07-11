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
)

func TestMemberMigrationCreate(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.ref = root.ref
	ba := testutils.RandomString()
	_, _ = signedRequest(&migrationAdmin, "migration.addBurnAddresses", map[string]interface{}{"burnAddresses": []string{ba}})
	result, err := retryableMemberMigrationCreate(member, true)
	require.NoError(t, err)
	output, ok := result.(map[string]interface{})
	require.True(t, ok)
	require.NotEqual(t, "", output["reference"])
	require.Equal(t, ba, output["migrationAddress"])
}

func TestMemberMigrationCreateWhenNoBurnAddressesLeft(t *testing.T) {
	member1, err := newUserWithKeys()
	require.NoError(t, err)
	member1.ref = root.ref
	addBurnAddress(t)
	_, err = retryableMemberMigrationCreate(member1, true)
	require.Nil(t, err)

	member2, err := newUserWithKeys()
	require.NoError(t, err)
	member2.ref = root.ref

	_, err = retryableMemberMigrationCreate(member2, true)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "no more burn addresses left")
}

func TestMemberMigrationCreateWithBadKey(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.ref = root.ref
	member.pubKey = "fake"
	_, err = retryableMemberMigrationCreate(member, false)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("problems with decoding. Key - %s", member.pubKey))
}

func TestMemberMigrationCreateWithSamePublicKey(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.ref = root.ref

	addBurnAddress(t)

	_, err = retryableMemberMigrationCreate(member, true)
	require.NoError(t, err)

	addBurnAddress(t)

	_, err = signedRequest(member, "member.migrationCreate", map[string]interface{}{})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "member for this publicKey already exist")

	memberForBurn, err := newUserWithKeys()
	require.NoError(t, err)
	memberForBurn.ref = root.ref

	_, err = retryableMemberMigrationCreate(memberForBurn, true)
}

func TestMemberMigrationCreateWithSameBurnAddress(t *testing.T) {
	member1, err := newUserWithKeys()
	require.NoError(t, err)
	member1.ref = root.ref

	ba := testutils.RandomString()
	_, _ = signedRequest(&migrationAdmin, "migration.addBurnAddresses", map[string]interface{}{"burnAddresses": []string{ba, ba}})

	_, err = retryableMemberMigrationCreate(member1, true)
	require.NoError(t, err)

	member2, err := newUserWithKeys()
	require.NoError(t, err)
	member2.ref = root.ref

	_, err = retryableMemberMigrationCreate(member2, true)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "member for this burnAddress already exist")
}
