// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"fmt"
	"testing"

	"github.com/insolar/insolar/application/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils"

	"github.com/stretchr/testify/require"
)

func TestMemberCreate(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	result, err := testutils.SignedRequest(t, launchnet.TestRPCUrlPublic, member, "member.create", nil)
	require.NoError(t, err)
	output, ok := result.(map[string]interface{})
	require.True(t, ok)
	require.NotEqual(t, "", output["reference"])
}

func TestMemberCreateWithBadKey(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.PubKey = "fake"
	_, err = testutils.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member, "member.create", nil)
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, fmt.Sprintf("problems with parsing. Key - %s", member.PubKey))
}

func TestMemberCreateWithSamePublicKey(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)

	_, err = testutils.SignedRequest(t, launchnet.TestRPCUrlPublic, member, "member.create", nil)
	require.NoError(t, err)

	_, err = testutils.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member, "member.create", nil)
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "can't set reference because this key already exists")
}
