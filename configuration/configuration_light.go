// Copyright 2020 Insolar Technologies GmbH
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

package configuration

import (
	"github.com/insolar/insconfig"
)

// ConfigLight contains configuration params for Light
type ConfigLight struct {
	GenericConfiguration `mapstructure:",squash" yaml:",inline"`
	Ledger               LedgerLight
}

func (c ConfigLight) GetConfig() interface{} {
	return &c
}

// NewConfigurationLight creates new default
func NewConfigurationLight() ConfigLight {
	return ConfigLight{
		Ledger:               NewLedgerLight(),
		GenericConfiguration: NewGenericConfiguration(),
	}
}

// HolderLight provides methods to manage light configuration
type HolderLight struct {
	Configuration ConfigLight
	Params        insconfig.Params
}

func (h HolderLight) GetGenericConfig() GenericConfiguration {
	return h.Configuration.GenericConfiguration
}
func (h HolderLight) GetNodeConfig() interface{} {
	return h.Configuration
}

// NewHolderLight creates new HolderLight with config path
func NewHolderLight(path string) *HolderLight {
	params := insconfig.Params{
		ConfigStruct:     ConfigLight{},
		EnvPrefix:        InsolarEnvPrefix,
		ConfigPathGetter: &stringPathGetter{Path: path},
		FileRequired:     false,
	}
	return &HolderLight{Configuration: ConfigLight{}, Params: params}
}

// Load method reads configuration from params file path
func (h *HolderLight) Load() error {
	insConfigurator := insconfig.NewInsConfigurator(h.Params)
	parsedConf, err := insConfigurator.Load()
	if err != nil {
		return err
	}
	cfg := parsedConf.(*ConfigLight)
	h.Configuration = *cfg
	return nil
}

// MustLoad wrapper around Load function which panics on error.
func (h *HolderLight) MustLoad() *HolderLight {
	err := h.Load()
	if err != nil {
		panic(err)
	}
	return h
}
