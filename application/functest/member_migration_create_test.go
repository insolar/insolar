// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application/testutils/launchnet"
)

func TestMemberMigrationCreate(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	result, err := signedRequest(t, launchnet.TestRPCUrlPublic, member, "member.migrationCreate", nil)
	require.NoError(t, err)
	output, ok := result.(map[string]interface{})
	require.True(t, ok)
	require.NotEqual(t, "", output["reference"])
	require.NotEqual(t, "", output["migrationAddress"])
}

func TestMemberMigrationCreateWithBadKey(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.PubKey = "fake"
	_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member, "member.migrationCreate", nil)
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, fmt.Sprintf("problems with parsing. Key - %s", member.PubKey))
}

func TestMemberMigrationCreateWithSamePublicKey(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)

	_, err = signedRequest(t, launchnet.TestRPCUrlPublic, member, "member.migrationCreate", nil)
	require.NoError(t, err)

	_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member, "member.migrationCreate", map[string]interface{}{})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "can't set reference because this key already exists")
}
