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

type VersionManager struct {
	VersionTable  map[string]*Feature
	AgreedVersion string
	viper         *viper.Viper
}

var instance *VersionManager

func GetVM() *VersionManager {
	if instance == nil {
		instance = newVersionManager()
	}
	return instance
}

func (vm *VersionManager) Verify(key string) bool {
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

func newVersionManager() *VersionManager {
	versionTable := make(map[string]*Feature)
	vm := &VersionManager{
		versionTable,
		InsBaseVersion,
		viper.New(),
	}
	vm.viper.SetDefault("versiontable", vm.VersionTable)
	vm.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vm.viper.SetEnvPrefix("insolar")
	vm.viper.SetConfigType("yml")
	return vm
}

func (vm *VersionManager) Load() error {
	err := vm.viper.ReadInConfig()
	if err != nil {
		return err
	}
	return vm.viper.UnmarshalKey("versiontable", &vm.VersionTable)
}

func (vm *VersionManager) LoadFromVariable() error {
	err := vm.viper.ReadConfig(strings.NewReader(InsVersionTable))
	if err != nil {
		return err
	}
	return vm.viper.UnmarshalKey("versiontable", &vm.VersionTable)
}

// SaveAs method writes configuration to particular file path
func (vm *VersionManager) SaveAs(path string) error {
	return vm.viper.WriteConfigAs(path)
}

// LoadFromFile method reads configuration from particular file path
func (vm *VersionManager) LoadFromFile(path string) error {
	vm.viper.SetConfigFile(path)
	return vm.Load()
}

// Save method writes configuration to default file path
func (vm *VersionManager) Save() error {
	vm.viper.Set("versiontable", vm.VersionTable)
	return vm.viper.WriteConfig()
}

func (vm *VersionManager) Add(key string, startVersion string, description string) (*Feature, error) {

	key = strings.ToLower(key)
	if vm.Get(key) != nil {
		return nil, errors.New("Feature already exists")
	}
	feature, err := NewFeature(key, startVersion, description)
	if err != nil {
		return nil, err
	}
	vm.VersionTable[key] = feature
	return feature, nil
}

func (vm *VersionManager) Get(key string) *Feature {
	key = strings.ToLower(key)
	if feature, ok := vm.VersionTable[key]; ok {
		return feature
	}
	return nil
}

func (vm *VersionManager) Remove(key string) {
	key = strings.ToLower(key)
	if _, ok := vm.VersionTable[key]; ok {
		delete(vm.VersionTable, key)
	}
}
