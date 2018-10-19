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

package manager

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

func newActiveNode(ver string) *core.ActiveNode {
	// key, _ := ecdsa.GeneratePrivateKey()
	return &core.ActiveNode{
		NodeID:   core.RecordRef{255},
		PulseNum: core.PulseNumber(0),
		State:    core.NodeActive,
		Roles:    []core.NodeRole{core.RoleUnknown},
		// PublicKey: &key.PublicKey,
	}
}

func TestGetMapOfVersions(t *testing.T) {
	nodes := []*core.ActiveNode{
		newActiveNode("v0.5.0"),
		newActiveNode("v0.5.0"),
		newActiveNode("v0.5.1"),
		newActiveNode("v0.5.1"),
	}
	mapOfVersions := getMapOfVersion(nodes)
	mapOfVersions2 := make(map[string]int)
	mapOfVersions2["v0.5.0"] = 4
	assert.NotNil(t, mapOfVersions)
	assert.NotNil(t, mapOfVersions2)
	assert.Equal(t, *mapOfVersions, mapOfVersions2)
}

func TestProcessVersionConsensus(t *testing.T) {
	nodes := []*core.ActiveNode{
		newActiveNode("v0.5.0"),
		newActiveNode("v0.5.0"),
		newActiveNode("v0.5.1"),
		newActiveNode("v0.5.1"),
	}
	assert.Error(t, ProcessVersionConsensus([]*core.ActiveNode{}))
	assert.NoError(t, ProcessVersionConsensus(nodes))
}

func TestGetMaxVersion(t *testing.T) {

	mapOfVersions := make(map[string]int)
	mapOfVersions["v0.3.0"] = 3
	mapOfVersions["v0.3.1"] = 4
	mapOfVersions["v0.3.2"] = 1
	res := getMaxVersion(5, &mapOfVersions)
	assert.Nil(t, res)
	mapOfVersions = make(map[string]int)
	mapOfVersions["v0.3.0"] = 3
	mapOfVersions["v0.3.1"] = 4
	res = getMaxVersion(4, &mapOfVersions)
	assert.Equal(t, *res, "v0.3.1")
	mapOfVersions = make(map[string]int)
	mapOfVersions["v0.3.0"] = 5
	mapOfVersions["v0.3.1"] = 4
	res = getMaxVersion(5, &mapOfVersions)
	assert.Equal(t, *res, "v0.3.0")
}

func TestGetRequired(t *testing.T) {
	assert.Equal(t, getRequired(5), 3)
	assert.Equal(t, getRequired(4), 3)
	assert.Equal(t, getRequired(7), 4)
	assert.Equal(t, getRequired(1), 1)
}

func TestVerify(t *testing.T) {
	vm := GetVM()
	feature, err := vm.Add("INSOLAR4", "v1.1.1", "Version manager for Insolar platform test")
	assert.NoError(t, err)
	assert.NotNil(t, feature)
	feature, err = vm.Add("INSOLAR5", "v1.1.2", "Version manager for Insolar platform test")
	assert.NoError(t, err)
	assert.NotNil(t, feature)
	assert.Equal(t, vm, GetVM())
	vm.AgreedVersion = "v1.1.1"
	assert.Equal(t, Verify("InsoLar4"), true)
	vm.AgreedVersion = "v1.1.0"
	assert.Equal(t, Verify("InsoLar4"), false)
	feature, err = vm.Add("INSOLAR6", "", "Version manager for Insolar platform test")
	assert.Error(t, err)
	assert.Nil(t, feature)

}
