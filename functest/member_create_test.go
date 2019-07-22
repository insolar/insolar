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

	"github.com/stretchr/testify/require"
)

func TestMemberCreate(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	result, err := retryableMemberCreate(member, true)
	require.NoError(t, err)
	output, ok := result.(map[string]interface{})
	require.True(t, ok, fmt.Sprintf("failed to decode result: expected map[string]interface{}, got %T", result))
	require.NotEqual(t, "", output["reference"])
}

func TestMemberCreateWithBadKey(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.pubKey = "fake"
	_, err = retryableMemberCreate(member, false)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("problems with decoding. Key - %s", member.pubKey))
}

func TestMemberCreateWithSamePublicKey(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)

	_, err = retryableMemberCreate(member, true)
	require.NoError(t, err)

	_, err = signedRequest(member, "member.create", map[string]interface{}{})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "member for this publicKey already exist")
}
