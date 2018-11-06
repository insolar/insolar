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
	"io/ioutil"
	"os"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/assert"
)

func TestNewVersionManager(t *testing.T) {
	vm, err := GetVersionManager()
	assert.NoError(t, err)
	feature, err := vm.Add("INSOLAR", "v1.1.1", "Version manager for Insolar platform test")
	assert.NoError(t, err)
	assert.NotNil(t, feature)

	feature, err = vm.Add("InsoLAR", "v1.1.1", "Version manager for Insolar platform test")
	assert.Error(t, err)
	assert.Nil(t, feature)

	feature, err = vm.Add("INSOLAR2", "v1.1.2", "Version manager for Insolar platform test")
	assert.NoError(t, err)
	assert.NotNil(t, feature)
	vm2, err := GetVersionManager()
	assert.NoError(t, err)
	assert.Equal(t, vm, vm2)

	vm.AgreedVersion, err = ParseVersion("v1.1.1")
	assert.NoError(t, err)
	assert.Equal(t, vm.IsAvailable("InsoLar"), true)
	vm.AgreedVersion, err = ParseVersion("v1.1.0")
	assert.NoError(t, err)
	assert.Equal(t, vm.IsAvailable("InsoLar"), false)
	vm.AgreedVersion, err = ParseVersion("v1.1.1")
	assert.NoError(t, err)
	assert.Equal(t, vm.IsAvailable("InsoLar10"), false)
}

func TestLoadSaveVersionManager(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)
	vm, err := NewVersionManager(configuration.VersionManager{"v0.3.0"})
	assert.NoError(t, err)
	feature, err := vm.Add("insolar", "v1.1.1", "Version manager for Insolar platform test")
	assert.NoError(t, err)
	assert.NotNil(t, feature)
	feature2, err := vm.Add("INSOLAR2", "v1.1.1", "Version manager for Insolar platform test")
	assert.NoError(t, err)
	assert.NotNil(t, feature2)
	err = vm.SaveAs(dir + "versiontable.yml")
	assert.NoError(t, err)
	feature, err = vm.Add("insolar3", "v1.1.2", "Version manager for Insolar platform test")
	assert.NoError(t, err)
	assert.NotNil(t, feature)
	vm2, err := NewVersionManager(configuration.VersionManager{"v0.3.0"})
	assert.NoError(t, err)
	err = vm2.LoadFromFile(dir + "versiontable.yml")
	assert.NoError(t, err)
	feature = vm2.Get("Insolar2")
	assert.NotNil(t, feature)
	assert.Equal(t, feature, feature2)
	vm2.Remove("insolar2")
	feature = vm2.Get("Insolar2")
	assert.Nil(t, feature)
	vm, err = NewVersionManager(configuration.VersionManager{"error"})
	assert.Error(t, err)
}
