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

// ConfigVirtual contains configuration params for Virtual node
type ConfigVirtual struct {
	GenericConfiguration `mapstructure:",squash" yaml:",inline"`
	LogicRunner          LogicRunner
}

func (c ConfigVirtual) GetConfig() interface{} {
	return &c
}

// NewConfigurationVirtual creates new default configuration
func NewConfigurationVirtual() ConfigVirtual {
	return ConfigVirtual{
		LogicRunner:          NewLogicRunner(),
		GenericConfiguration: NewGenericConfiguration(),
	}
}

// HolderVirtual provides methods to manage virtual node configuration
type HolderVirtual struct {
	Configuration ConfigVirtual
	Params        insconfig.Params
}

func (h HolderVirtual) GetGenericConfig() GenericConfiguration {
	return h.Configuration.GenericConfiguration
}
func (h HolderVirtual) GetNodeConfig() interface{} {
	return h.Configuration
}

// NewHolderVirtual creates new HolderVirtual with config path
func NewHolderVirtual(path string) *HolderVirtual {
	params := insconfig.Params{
		ConfigStruct:     ConfigVirtual{},
		EnvPrefix:        InsolarEnvPrefix,
		ConfigPathGetter: &stringPathGetter{Path: path},
		FileRequired:     false,
	}
	return &HolderVirtual{Configuration: ConfigVirtual{}, Params: params}
}

// Load method reads configuration from params file path
func (h *HolderVirtual) Load() error {
	insConfigurator := insconfig.NewInsConfigurator(h.Params)
	parsedConf, err := insConfigurator.Load()
	if err != nil {
		return err
	}
	cfg := parsedConf.(*ConfigVirtual)
	h.Configuration = *cfg
	return nil
}
