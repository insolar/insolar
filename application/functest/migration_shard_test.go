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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application/api/requester"
	"github.com/insolar/insolar/application/testutils/launchnet"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

func TestGetFreeAddressCount(t *testing.T) {
	migrationShardsMap := getAddressCount(t, 0)

	for _, m := range migrationShardsMap {
		require.True(t, m > 0)
	}
}

const numShards = 1000

func TestGetFreeAddressCount_ChangesAfterMigration(t *testing.T) {

	member, err := newUserWithKeys()
	require.NoError(t, err)

	trimmedPublicKey := foundation.TrimPublicKey(member.PubKey)
	shardIndex := foundation.GetShardIndex(trimmedPublicKey, numShards)

	var migrationShardsMapBefore = getAddressCount(t, shardIndex)

	result, err := signedRequest(t, launchnet.TestRPCUrlPublic, member, "member.migrationCreate", nil)
	require.NoError(t, err)
	output, ok := result.(map[string]interface{})
	require.True(t, ok)
	require.NotEqual(t, "", output["reference"])
	require.NotEqual(t, "", output["migrationAddress"])

	result, err = signedRequest(t, launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.getAddressCount",
		map[string]interface{}{"startWithIndex": 0})
	require.NoError(t, err)

	var migrationShardsMapAfter = getAddressCount(t, shardIndex)

	isFound := false
	for i, countBefore := range migrationShardsMapBefore {
		countAfter := migrationShardsMapAfter[i]
		if countBefore == countAfter {
			continue
		}
		if (countBefore-countAfter) == 1 && !isFound {
			isFound = true
			continue
		}
		t.Errorf("Wrong count of free migration addresses: for shard %d, "+
			"count before one migration is %d, "+
			"after %d (migration was already found - %t)", i, countBefore, countAfter, isFound)

	}
	require.True(t, isFound)
}

func TestGetFreeAddressCount_WithIndex_NotAllRange(t *testing.T) {
	numLeftShards := 2
	var migrationShards = getAddressCount(t, numShards-numLeftShards)
	require.Len(t, migrationShards, numLeftShards)
}

func TestGetFreeAddressCount_StartIndexTooBig(t *testing.T) {
	_, _, err := makeSignedRequest(launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.getAddressCount",
		map[string]interface{}{"startWithIndex": numShards + 2})
	require.Error(t, err)
	require.IsType(t, &requester.Error{}, err)
	data := err.(*requester.Error).Data
	require.Contains(t, data.Trace, "incorrect start shard index")
}

func TestGetFreeAddressCount_IncorrectIndexType(t *testing.T) {
	_, _, err := makeSignedRequest(launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.getAddressCount",
		map[string]interface{}{"startWithIndex": "0"})
	require.Error(t, err)
	require.IsType(t, &requester.Error{}, err)
	data := err.(*requester.Error).Data
	require.Contains(t, data.Trace, "failed to get 'startWithIndex' param")
}

func TestGetFreeAddressCount_FromMember(t *testing.T) {
	member := createMember(t)
	_, _, err := makeSignedRequest(launchnet.TestRPCUrl, member, "migration.getAddressCount",
		map[string]interface{}{"startWithIndex": 0})
	require.Error(t, err)
	require.IsType(t, &requester.Error{}, err)
	data := err.(*requester.Error).Data
	require.Contains(t, data.Trace, "only migration daemon admin can call this method")
}
