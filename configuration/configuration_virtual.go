// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
	Configuration *ConfigVirtual
	Params        insconfig.Params
}

func (h *HolderVirtual) GetGenericConfig() GenericConfiguration {
	return h.Configuration.GenericConfiguration
}
func (h *HolderVirtual) GetNodeConfig() interface{} {
	return h.Configuration
}

// NewHolderVirtual creates new HolderVirtual with config path
func NewHolderVirtual(path string) *HolderVirtual {
	params := insconfig.Params{
		EnvPrefix:        InsolarEnvPrefix,
		ConfigPathGetter: &stringPathGetter{Path: path},
	}
	return &HolderVirtual{Configuration: &ConfigVirtual{}, Params: params}
}

// Load method reads configuration from params file path
func (h *HolderVirtual) Load() error {
	insConfigurator := insconfig.New(h.Params)
	if err := insConfigurator.Load(h.Configuration); err != nil {
		return err
	}
	return nil
}
