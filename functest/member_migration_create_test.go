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
	"github.com/insolar/insolar/api/requester"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/testutils/launchnet"
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
	require.Error(t, err)
	require.IsType(t, &requester.Error{}, err)
	data := err.(*requester.Error).Data
	require.Contains(t, data.Trace, fmt.Sprintf("problems with decoding. Key - %s", member.PubKey))
}

func TestMemberMigrationCreateWithSamePublicKey(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)

	_, err = signedRequest(t, launchnet.TestRPCUrlPublic, member, "member.migrationCreate", nil)
	require.NoError(t, err)

	_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member, "member.migrationCreate", map[string]interface{}{})
	require.Error(t, err)
	require.IsType(t, &requester.Error{}, err)
	data := err.(*requester.Error).Data
	require.Contains(t, data.Trace, "can't set reference because this key already exists")
}
