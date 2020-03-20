// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

import (
	"github.com/insolar/insconfig"
)

// VirtualConfig contains configuration params for Virtual node
type VirtualConfig struct {
	GenericConfiguration `mapstructure:",squash" yaml:",inline"`
	LogicRunner          LogicRunner
}

func (c VirtualConfig) GetConfig() interface{} {
	return &c
}

// NewVirtualConfig creates new default configuration
func NewVirtualConfig() VirtualConfig {
	return VirtualConfig{
		LogicRunner:          NewLogicRunner(),
		GenericConfiguration: NewGenericConfiguration(),
	}
}

// VirtualHolder provides methods to manage virtual node configuration
type VirtualHolder struct {
	Configuration *VirtualConfig
	Params        insconfig.Params
}

func (h *VirtualHolder) GetGenericConfig() GenericConfiguration {
	return h.Configuration.GenericConfiguration
}
func (h *VirtualHolder) GetNodeConfig() interface{} {
	return h.Configuration
}

// NewVirtualHolder creates new VirtualHolder with config path
func NewVirtualHolder(path string) *VirtualHolder {
	params := insconfig.Params{
		EnvPrefix:        InsolarEnvPrefix,
		ConfigPathGetter: &stringPathGetter{Path: path},
	}
	return &VirtualHolder{Configuration: &VirtualConfig{}, Params: params}
}

// Load method reads configuration from params file path
func (h *VirtualHolder) Load() error {
	insConfigurator := insconfig.New(h.Params)
	if err := insConfigurator.Load(h.Configuration); err != nil {
		return err
	}
	return nil
}
