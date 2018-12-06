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
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/stretchr/testify/assert"
)

func newActiveNode(ver string) core.Node {
	return nodenetwork.NewNode(core.RecordRef{255}, core.StaticRoleUnknown, nil, "", ver)
}

func TestGetMapOfVersions(t *testing.T) {
	nodes := []core.Node{
		newActiveNode("v0.5.0"),
		newActiveNode("v0.5.0"),
		newActiveNode("v0.5.1"),
		newActiveNode("v0.5.1"),
	}
	mapOfVersions := getMapOfVersion(nodes)
	mapOfVersions2 := make(map[string]int)
	mapOfVersions2["v0.5.0"] = 2
	mapOfVersions2["v0.5.1"] = 2
	assert.NotNil(t, mapOfVersions)
	assert.NotNil(t, mapOfVersions2)
	assert.Equal(t, *mapOfVersions, mapOfVersions2)
}

func TestProcessVersionConsensus(t *testing.T) {
	nodes := []core.Node{
		newActiveNode("v0.5.0"),
		newActiveNode("v0.5.0"),
		newActiveNode("v0.5.1"),
		newActiveNode("v0.5.1"),
	}
	assert.Error(t, ProcessVersionConsensus([]core.Node{}))
	assert.NoError(t, ProcessVersionConsensus(nodes))
}

func TestGetMaxVersion(t *testing.T) {

	mapOfVersions := make(map[string]int)
	mapOfVersions["v0.3.0"] = 3
	mapOfVersions["v0.3.1"] = 4
	mapOfVersions["v0.3.2"] = 1
	res, err := getMaxVersion(5, &mapOfVersions)
	assert.NoError(t, err)
	assert.Equal(t, StringVersion(res), "v0.3.1")
	mapOfVersions = make(map[string]int)
	mapOfVersions["v0.3.0"] = 3
	mapOfVersions["v0.3.1"] = 4
	res, err = getMaxVersion(4, &mapOfVersions)
	assert.NoError(t, err)
	assert.Equal(t, StringVersion(res), "v0.3.1")
	mapOfVersions = make(map[string]int)
	mapOfVersions["v0.3.0"] = 5
	mapOfVersions["v0.3.1"] = 4
	res, err = getMaxVersion(5, &mapOfVersions)
	assert.NoError(t, err)
	assert.Equal(t, StringVersion(res), "v0.3.0")

	mapOfVersions = make(map[string]int)
	res, err = getMaxVersion(5, &mapOfVersions)
	assert.Error(t, err)
	mapOfVersions["error"] = 1
	res, err = getMaxVersion(5, &mapOfVersions)
	assert.Error(t, err)
	mapOfVersions["error"] = 6
	res, err = getMaxVersion(5, &mapOfVersions)
	assert.Error(t, err)
}

func TestGetRequired(t *testing.T) {
	assert.Equal(t, getRequired(5), 3)
	assert.Equal(t, getRequired(4), 3)
	assert.Equal(t, getRequired(7), 4)
	assert.Equal(t, getRequired(1), 1)
}

func TestVerify(t *testing.T) {
	vm, err := GetVersionManager()
	assert.NoError(t, err)
	feature, err := vm.Add("INSOLAR4", "v1.1.1", "Version manager for Insolar platform test")
	assert.NoError(t, err)
	assert.NotNil(t, feature)
	feature, err = vm.Add("INSOLAR5", "v1.1.2", "Version manager for Insolar platform test")
	assert.NoError(t, err)
	assert.NotNil(t, feature)
	vm2, err := GetVersionManager()
	assert.NoError(t, err)
	assert.Equal(t, vm, vm2)
	vm.AgreedVersion, err = ParseVersion("v1.1.1")
	assert.NoError(t, err)
	assert.Equal(t, Verify("InsoLar4"), true)
	vm.AgreedVersion, err = ParseVersion("v1.1.0")
	assert.NoError(t, err)
	assert.Equal(t, Verify("InsoLar4"), false)
	feature, err = vm.Add("INSOLAR6", "", "Version manager for Insolar platform test")
	assert.Error(t, err)
	assert.Nil(t, feature)

	ver, err := ParseVersion("abc")
	assert.Error(t, err)
	assert.Nil(t, ver)
}
