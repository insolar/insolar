// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

import (
	"github.com/insolar/insconfig"
)

// ConfigLight contains configuration params for Light
type ConfigLight struct {
	GenericConfiguration `mapstructure:",squash" yaml:",inline"`
	Ledger               LedgerLight
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
	Configuration *ConfigLight
	Params        insconfig.Params
}

func (h *HolderLight) GetGenericConfig() GenericConfiguration {
	return h.Configuration.GenericConfiguration
}
func (h *HolderLight) GetNodeConfig() interface{} {
	return h.Configuration
}

// NewHolderLight creates new HolderLight with config path
func NewHolderLight(path string) *HolderLight {
	params := insconfig.Params{
		EnvPrefix:        InsolarEnvPrefix,
		ConfigPathGetter: &stringPathGetter{Path: path},
	}
	return &HolderLight{Configuration: &ConfigLight{}, Params: params}
}

// Load method reads configuration from params file path
func (h *HolderLight) Load() error {
	insConfigurator := insconfig.New(h.Params)
	if err := insConfigurator.Load(h.Configuration); err != nil {
		return err
	}
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
