// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

import (
	"github.com/insolar/insconfig"
)

// LightConfig contains configuration params for Light
type LightConfig struct {
	GenericConfiguration `mapstructure:",squash" yaml:",inline"`
	Ledger               LedgerLight
}

// NewLightConfig creates new default
func NewLightConfig() LightConfig {
	return LightConfig{
		Ledger:               NewLedgerLight(),
		GenericConfiguration: NewGenericConfiguration(),
	}
}

// LightHolder provides methods to manage light configuration
type LightHolder struct {
	Configuration *LightConfig
	Params        insconfig.Params
}

func (h *LightHolder) GetGenericConfig() GenericConfiguration {
	return h.Configuration.GenericConfiguration
}
func (h *LightHolder) GetNodeConfig() interface{} {
	return h.Configuration
}

// NewLightHolder creates new LightHolder with config path
func NewLightHolder(path string) *LightHolder {
	params := insconfig.Params{
		EnvPrefix:        InsolarEnvPrefix,
		ConfigPathGetter: &stringPathGetter{Path: path},
	}
	return &LightHolder{Configuration: &LightConfig{}, Params: params}
}

// Load method reads configuration from params file path
func (h *LightHolder) Load() error {
	insConfigurator := insconfig.New(h.Params)
	if err := insConfigurator.Load(h.Configuration); err != nil {
		return err
	}
	return nil
}

// MustLoad wrapper around Load function which panics on error.
func (h *LightHolder) MustLoad() *LightHolder {
	err := h.Load()
	if err != nil {
		panic(err)
	}
	return h
}
