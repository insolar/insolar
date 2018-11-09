/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package noderecord

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

const TestPubKey = "test"
const TestIP = "127.0.0.1"

var TestRole = "virtual"

func TestNewNodeRecord(t *testing.T) {

	r := core.GetRoleFromString(TestRole)
	assert.NotEqual(t, core.RoleUnknown, r)
	record, err := NewNodeRecord(TestPubKey, TestRole, TestIP)
	assert.NoError(t, err)
	assert.Equal(t, r, record.Record.Role)
	assert.Equal(t, TestPubKey, record.Record.PublicKey)
}

func TestFromString(t *testing.T) {
	role := core.GetRoleFromString("ZZZ")
	assert.Equal(t, core.RoleUnknown, role)
}

func TestNodeRecord_GetPublicKey(t *testing.T) {
	record, err := NewNodeRecord(TestPubKey, TestRole, TestIP)
	assert.NoError(t, err)
	pk, err := record.GetPublicKey()
	assert.NoError(t, err)
	assert.Equal(t, TestPubKey, pk)
}

func TestNodeRecord_GetNodeInfo(t *testing.T) {
	record, err := NewNodeRecord(TestPubKey, TestRole, TestIP)
	assert.NoError(t, err)
	info, err := record.GetNodeInfo()
	assert.NoError(t, err)
	assert.Equal(t, TestPubKey, info.PublicKey)
	r := core.GetRoleFromString(TestRole)
	assert.Equal(t, r, info.Role)
	assert.Equal(t, TestIP, info.IP)
}

func TestNodeRecord_GetRole(t *testing.T) {
	record, err := NewNodeRecord(TestPubKey, TestRole, TestIP)
	assert.NoError(t, err)
	role, err := record.GetRole()
	assert.NoError(t, err)
	r := core.GetRoleFromString(TestRole)
	assert.Equal(t, r, role)
}
