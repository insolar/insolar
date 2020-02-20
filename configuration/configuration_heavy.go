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

// ConfigHeavyBadger contains configuration params for HeavyBadger
type ConfigHeavyBadger struct {
	GenericConfiguration `mapstructure:",squash"`
	Ledger               Ledger
	Exporter             Exporter
}

// ConfigHeavyPg contains configuration params for HeavyPg
type ConfigHeavyPg struct {
	GenericConfiguration `mapstructure:",squash"`
	Ledger               LedgerPg
	Exporter             Exporter
}

func (c ConfigHeavyBadger) GetConfig() interface{} {
	return &c
}

func (c ConfigHeavyPg) GetConfig() interface{} {
	return &c
}

// NewConfigurationHeavyBadger creates new default configuration
func NewConfigurationHeavyBadger() ConfigHeavyBadger {
	return ConfigHeavyBadger{
		Ledger:               NewLedger(),
		Exporter:             NewExporter(),
		GenericConfiguration: NewConfiguration(),
	}
}

// NewConfigurationHeavyBadger creates new default configuration
func NewConfigurationHeavyPg() ConfigHeavyPg {
	cfg := ConfigHeavyPg{
		Ledger:               NewLedgerPg(),
		Exporter:             NewExporter(),
		GenericConfiguration: NewConfiguration(),
	}
	return cfg
}

// HolderHeavyBadger provides methods to manage heavy configuration
type HolderHeavyBadger struct {
	Configuration ConfigHeavyBadger
	Params        insconfig.Params
}

func (h HolderHeavyBadger) GetGenericConfig() GenericConfiguration {
	return h.Configuration.GenericConfiguration
}
func (h HolderHeavyBadger) GetAllConfig() interface{} {
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
		ConfigStruct:     ConfigHeavyBadger{},
		EnvPrefix:        InsolarEnvPrefix,
		ConfigPathGetter: &stringPathGetter{Path: path},
		FileRequired:     false,
	}
	return &HolderHeavyBadger{Configuration: ConfigHeavyBadger{}, Params: params}
}

// Load method reads configuration from params file path
func (h *HolderHeavyBadger) Load() error {
	insConfigurator := insconfig.NewInsConfigurator(h.Params)
	parsedConf, err := insConfigurator.Load()
	if err != nil {
		return err
	}
	cfg := parsedConf.(*ConfigHeavyBadger)
	h.Configuration = *cfg
	return nil
}

// HolderHeavyPg provides methods to manage heavy configuration
type HolderHeavyPg struct {
	Configuration ConfigHeavyPg
	Params        insconfig.Params
}

func (h HolderHeavyPg) GetGenericConfig() GenericConfiguration {
	return h.Configuration.GenericConfiguration
}
func (h HolderHeavyPg) GetAllConfig() interface{} {
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
		ConfigStruct:     ConfigHeavyPg{},
		EnvPrefix:        InsolarEnvPrefix,
		ConfigPathGetter: &stringPathGetter{Path: path},
		FileRequired:     false,
	}
	return &HolderHeavyPg{Configuration: ConfigHeavyPg{}, Params: params}
}

// Load method reads configuration from params file path
func (h *HolderHeavyPg) Load() error {
	insConfigurator := insconfig.NewInsConfigurator(h.Params)
	parsedConf, err := insConfigurator.Load()
	if err != nil {
		return err
	}
	cfg := parsedConf.(*ConfigHeavyPg)
	h.Configuration = *cfg
	return nil
}
