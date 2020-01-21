// Copyright 2020 Insolar Network Ltd.
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

// +build functest

package functest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application/testutils/launchnet"
)

func TestGetFreeAddressCount(t *testing.T) {
	migrationShardsMap := getAddressCount(t, 0)

	for _, m := range migrationShardsMap {
		require.True(t, m > 0)
	}
}

func TestGetFreeAddressCount_WithIndex_NotAllRange(t *testing.T) {
	numLeftShards := 2
	numShards, err := launchnet.GetNumShards()
	require.NoError(t, err)
	var migrationShards = getAddressCount(t, numShards-numLeftShards)
	require.Len(t, migrationShards, numLeftShards)
}

func TestGetFreeAddressCount_StartIndexTooBig(t *testing.T) {
	numShards, err := launchnet.GetNumShards()
	require.NoError(t, err)
	_, _, err = makeSignedRequest(launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.getAddressCount",
		map[string]interface{}{"startWithIndex": numShards + 2})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "incorrect start shard index")
}

func TestGetFreeAddressCount_IncorrectIndexType(t *testing.T) {
	_, _, err := makeSignedRequest(launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.getAddressCount",
		map[string]interface{}{"startWithIndex": "0"})
	data := checkConvertRequesterError(t, err).Data
	expectedError(t, data.Trace, "doesn't match the schema")
}

func TestGetFreeAddressCount_FromMember(t *testing.T) {
	member := createMember(t)
	_, _, err := makeSignedRequest(launchnet.TestRPCUrl, member, "migration.getAddressCount",
		map[string]interface{}{"startWithIndex": 0})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "only migration daemon admin can call this method")
}
