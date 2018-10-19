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
	"errors"
	"strings"

	"github.com/spf13/viper"
)

type versionManager struct {
	VersionTable  map[string]*Feature
	AgreedVersion string
	viper         *viper.Viper
}

var instance *versionManager

func GetVM() *versionManager {
	if instance == nil {
		instance = newVersionManager()
	}
	return instance
}

func (vm *versionManager) Verify(key string) bool {
	key = strings.ToLower(key)
	feature := vm.Get(key)
	if feature == nil {
		return false
	}
	if feature.StartVersion <= vm.AgreedVersion {
		return true
	}
	return false
}

func newVersionManager() *versionManager {
	versionTable := make(map[string]*Feature)
	vm := &versionManager{
		versionTable,
		"v0.0.0",
		viper.New(),
	}
	vm.viper.SetConfigName("versiontable")
	vm.viper.AddConfigPath(".")
	vm.viper.SetConfigType("yml")

	vm.viper.SetDefault("versiontable", versionTable)

	vm.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vm.viper.SetEnvPrefix("insolar")
	return vm
}

func (vm *versionManager) Load() error {
	err := vm.viper.ReadInConfig()
	if err != nil {
		return err
	}
	return vm.viper.UnmarshalKey("versiontable", &vm.VersionTable)
}

// SaveAs method writes configuration to particular file path
func (vm *versionManager) SaveAs(path string) error {
	return vm.viper.WriteConfigAs(path)
}

// LoadFromFile method reads configuration from particular file path
func (vm *versionManager) LoadFromFile(path string) error {
	vm.viper.SetConfigFile(path)
	return vm.Load()
}

// Save method writes configuration to default file path
func (vm *versionManager) Save() error {
	vm.viper.Set("agreedversion", vm.AgreedVersion)
	vm.viper.Set("versiontable", vm.VersionTable)
	return vm.viper.WriteConfig()
}

func (vm *versionManager) Add(key string, startVersion string, description string) (*Feature, error) {

	key = strings.ToLower(key)
	if vm.Get(key) != nil {
		return nil, errors.New("Feature already exists")
	}
	feature, err := NewFeature(key, startVersion, description)
	if err != nil {
		return nil, err
	}
	vm.VersionTable[key] = feature
	vm.Save()
	return feature, nil
}

func (vm *versionManager) Get(key string) *Feature {
	key = strings.ToLower(key)
	if feature, ok := vm.VersionTable[key]; ok {
		return feature
	}
	return nil
}

func (vm *versionManager) Remove(key string) {
	key = strings.ToLower(key)
	if _, ok := vm.VersionTable[key]; ok {
		delete(vm.VersionTable, key)
	}
}
