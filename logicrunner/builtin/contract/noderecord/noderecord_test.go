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

package noderecord

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
)

const TestPubKey = "test"

var TestRole = "virtual"

func TestNewNodeRecord(t *testing.T) {

	r := insolar.GetStaticRoleFromString(TestRole)
	require.NotEqual(t, insolar.StaticRoleUnknown, r)
	record, err := NewNodeRecord(TestPubKey, TestRole)
	require.NoError(t, err)
	require.Equal(t, r, record.Record.Role)
	require.Equal(t, TestPubKey, record.Record.PublicKey)
}

func TestFromString(t *testing.T) {
	role := insolar.GetStaticRoleFromString("ZZZ")
	require.Equal(t, insolar.StaticRoleUnknown, role)
}

func TestNodeRecord_GetPublicKey(t *testing.T) {
	record, err := NewNodeRecord(TestPubKey, TestRole)
	require.NoError(t, err)
	pk, err := record.GetPublicKey()
	require.NoError(t, err)
	require.Equal(t, TestPubKey, pk)
}

func TestNodeRecord_GetNodeInfo(t *testing.T) {
	record, err := NewNodeRecord(TestPubKey, TestRole)
	require.NoError(t, err)
	info, err := record.GetNodeInfo()
	require.NoError(t, err)
	require.Equal(t, TestPubKey, info.PublicKey)
	r := insolar.GetStaticRoleFromString(TestRole)
	require.Equal(t, r, info.Role)
}

func TestNodeRecord_GetRole(t *testing.T) {
	record, err := NewNodeRecord(TestPubKey, TestRole)
	require.NoError(t, err)
	role, err := record.GetRole()
	require.NoError(t, err)
	r := insolar.GetStaticRoleFromString(TestRole)
	require.Equal(t, r, role)
}
