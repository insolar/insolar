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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/testutils"
)

func TestMigrationGetAddressCount(t *testing.T) {
	ma := "1"
	_, err := signedRequest(&migrationAdmin, "migration.addBurnAddresses", map[string]interface{}{"burnAddresses": []string{ma}})
	require.NoError(t, err)

	member := createMember(t)

	result, err := signedRequest(member, "migration.getAddressCount", nil)
	require.NoError(t, err)
	output, ok := result.(map[string]interface{})
	require.True(t, ok)
	shardCounts, ok := output["shardCounts"].([]interface{})
	require.True(t, ok)
	require.Equal(t, float64(1), shardCounts[4])

	m1, err := newUserWithKeys()
	require.NoError(t, err)
	_, err = retryableMemberMigrationCreate(m1, true)
	require.NoError(t, err)

	result, err = signedRequest(member, "migration.getAddressCount", nil)
	require.NoError(t, err)
	output, ok = result.(map[string]interface{})
	require.True(t, ok)
	shardCounts, ok = output["shardCounts"].([]interface{})
	require.True(t, ok)
	require.Equal(t, float64(0), shardCounts[4])
}

func TestMigrationGetAddressCountWithManyAddresses(t *testing.T) {
	const maCount = 10
	maList := [maCount]string{}
	maAmountList := [insolar.GenesisAmountMigrationAddressShards]int{}

	for i := 0; i < maCount; i++ {
		ma := testutils.RandomString()
		maList[i] = ma
		index := foundation.GetShardIndex(ma, insolar.GenesisAmountMigrationAddressShards)
		maAmountList[index] = maAmountList[index] + 1
	}

	_, err := signedRequest(&migrationAdmin, "migration.addBurnAddresses", map[string]interface{}{"burnAddresses": maList[:]})
	require.NoError(t, err)

	member := createMember(t)

	result, err := signedRequest(member, "migration.getAddressCount", nil)
	require.NoError(t, err)
	output, ok := result.(map[string]interface{})
	require.True(t, ok)
	shardCounts, ok := output["shardCounts"].([]interface{})
	require.True(t, ok)

	for i, a := range maAmountList {
		require.Equal(t, float64(a), shardCounts[i])
	}

	for i := 0; i < maCount; i++ {
		m, err := newUserWithKeys()
		require.NoError(t, err)
		_, err = retryableMemberMigrationCreate(m, true)
		require.NoError(t, err)
	}

	result, err = signedRequest(member, "migration.getAddressCount", nil)
	require.NoError(t, err)
	output, ok = result.(map[string]interface{})
	require.True(t, ok)
	shardCounts, ok = output["shardCounts"].([]interface{})
	require.True(t, ok)

	for i := 0; i < insolar.GenesisAmountMigrationAddressShards; i++ {
		require.Equal(t, float64(0), shardCounts[i])
	}
}
