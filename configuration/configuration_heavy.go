// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

import (
	"github.com/insolar/insconfig"
)

const DbTypeBadger = "badger"
const DbTypePg = "postgres"

// ConfigHeavyBadger contains configuration params for HeavyBadger
type ConfigHeavyBadger struct {
	GenericConfiguration `mapstructure:",squash" yaml:",inline"`
	DatabaseType         string
	Ledger               Ledger
	Exporter             Exporter
}

// ConfigHeavyPg contains configuration params for HeavyPg
type ConfigHeavyPg struct {
	GenericConfiguration `mapstructure:",squash" yaml:",inline"`
	DatabaseType         string
	Ledger               LedgerPg
	Exporter             Exporter
}

// NewConfigurationHeavyBadger creates new default configuration
func NewConfigurationHeavyBadger() ConfigHeavyBadger {
	return ConfigHeavyBadger{
		DatabaseType:         DbTypeBadger,
		Ledger:               NewLedger(),
		Exporter:             NewExporter(),
		GenericConfiguration: NewGenericConfiguration(),
	}
}

// NewConfigurationHeavyBadger creates new default configuration
func NewConfigurationHeavyPg() ConfigHeavyPg {
	cfg := ConfigHeavyPg{
		DatabaseType:         DbTypePg,
		Ledger:               NewLedgerPg(),
		Exporter:             NewExporter(),
		GenericConfiguration: NewGenericConfiguration(),
	}
	return cfg
}

// HolderHeavyBadger provides methods to manage heavy configuration
type HolderHeavyBadger struct {
	Configuration *ConfigHeavyBadger
	Params        insconfig.Params
}

func (h *HolderHeavyBadger) GetGenericConfig() GenericConfiguration {
	return h.Configuration.GenericConfiguration
}
func (h *HolderHeavyBadger) GetNodeConfig() interface{} {
	return h.Configuration
}

// MustLoad wrapper around Load function which panics on error.
func (h *HolderHeavyBadger) MustLoad() *HolderHeavyBadger {
	err := h.Load()
	if err != nil {
		panic(err)
	}
	return h
}

// NewHolderHeavyBadger creates new HolderHeavyBadger with config path
func NewHolderHeavyBadger(path string) *HolderHeavyBadger {
	params := insconfig.Params{
		EnvPrefix:        InsolarEnvPrefix,
		ConfigPathGetter: &stringPathGetter{Path: path},
	}
	return &HolderHeavyBadger{Configuration: &ConfigHeavyBadger{}, Params: params}
}

// Load method reads configuration from params file path
func (h *HolderHeavyBadger) Load() error {
	insConfigurator := insconfig.New(h.Params)
	if err := insConfigurator.Load(h.Configuration); err != nil {
		return err
	}
	return nil
}

// HolderHeavyPg provides methods to manage heavy configuration
type HolderHeavyPg struct {
	Configuration *ConfigHeavyPg
	Params        insconfig.Params
}

func (h *HolderHeavyPg) GetGenericConfig() GenericConfiguration {
	return h.Configuration.GenericConfiguration
}
func (h *HolderHeavyPg) GetNodeConfig() interface{} {
	return h.Configuration
}

// MustLoad wrapper around Load function which panics on error.
func (h *HolderHeavyPg) MustLoad() *HolderHeavyPg {
	err := h.Load()
	if err != nil {
		panic(err)
	}
	return h
}

// NewHolderHeavyPg creates new HolderHeavyPg with config path
func NewHolderHeavyPg(path string) *HolderHeavyPg {
	params := insconfig.Params{
		EnvPrefix:        InsolarEnvPrefix,
		ConfigPathGetter: &stringPathGetter{Path: path},
	}
	return &HolderHeavyPg{Configuration: &ConfigHeavyPg{}, Params: params}
}

// Load method reads configuration from params file path
func (h *HolderHeavyPg) Load() error {
	insConfigurator := insconfig.New(h.Params)
	if err := insConfigurator.Load(h.Configuration); err != nil {
		return err
	}
	return nil
}
