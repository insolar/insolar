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

// HeavyBadgerConfig contains configuration params for HeavyBadger
type HeavyBadgerConfig struct {
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

// NewHeavyBadgerConfig creates new default configuration
func NewHeavyBadgerConfig() HeavyBadgerConfig {
	return HeavyBadgerConfig{
		DatabaseType:         DbTypeBadger,
		Ledger:               NewLedger(),
		Exporter:             NewExporter(),
		GenericConfiguration: NewGenericConfiguration(),
	}
}

// NewHeavyPgConfig creates new default configuration
func NewHeavyPgConfig() ConfigHeavyPg {
	cfg := ConfigHeavyPg{
		DatabaseType:         DbTypePg,
		Ledger:               NewLedgerPg(),
		Exporter:             NewExporter(),
		GenericConfiguration: NewGenericConfiguration(),
	}
	return cfg
}

// HeavyBadgerHolder provides methods to manage heavy configuration
type HeavyBadgerHolder struct {
	Configuration *HeavyBadgerConfig
	Params        insconfig.Params
}

func (h *HeavyBadgerHolder) GetGenericConfig() GenericConfiguration {
	return h.Configuration.GenericConfiguration
}
func (h *HeavyBadgerHolder) GetNodeConfig() interface{} {
	return h.Configuration
}

// MustLoad wrapper around Load function which panics on error.
func (h *HeavyBadgerHolder) MustLoad() *HeavyBadgerHolder {
	err := h.Load()
	if err != nil {
		panic(err)
	}
	return h
}

// NewHeavyBadgerHolder creates new HeavyBadgerHolder with config path
func NewHeavyBadgerHolder(path string) *HeavyBadgerHolder {
	params := insconfig.Params{
		EnvPrefix:        InsolarEnvPrefix,
		ConfigPathGetter: &stringPathGetter{Path: path},
	}
	return &HeavyBadgerHolder{Configuration: &HeavyBadgerConfig{}, Params: params}
}

// Load method reads configuration from params file path
func (h *HeavyBadgerHolder) Load() error {
	insConfigurator := insconfig.New(h.Params)
	if err := insConfigurator.Load(h.Configuration); err != nil {
		return err
	}
	return nil
}

// HeavyPgHolder provides methods to manage heavy configuration
type HeavyPgHolder struct {
	Configuration *ConfigHeavyPg
	Params        insconfig.Params
}

func (h *HeavyPgHolder) GetGenericConfig() GenericConfiguration {
	return h.Configuration.GenericConfiguration
}
func (h *HeavyPgHolder) GetNodeConfig() interface{} {
	return h.Configuration
}

// MustLoad wrapper around Load function which panics on error.
func (h *HeavyPgHolder) MustLoad() *HeavyPgHolder {
	err := h.Load()
	if err != nil {
		panic(err)
	}
	return h
}

// NewHeavyPgHolder creates new HeavyPgHolder with config path
func NewHeavyPgHolder(path string) *HeavyPgHolder {
	params := insconfig.Params{
		EnvPrefix:        InsolarEnvPrefix,
		ConfigPathGetter: &stringPathGetter{Path: path},
	}
	return &HeavyPgHolder{Configuration: &ConfigHeavyPg{}, Params: params}
}

// Load method reads configuration from params file path
func (h *HeavyPgHolder) Load() error {
	insConfigurator := insconfig.New(h.Params)
	if err := insConfigurator.Load(h.Configuration); err != nil {
		return err
	}
	return nil
}
