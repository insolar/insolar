// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

import (
	"fmt"

	"github.com/insolar/insconfig"
	"gopkg.in/yaml.v2"
)

const InsolarEnvPrefix string = "insolar"

var GlobalHolder *Holder = nil

// Configuration contains configuration params for all Insolar components
type Configuration struct {
	Host                HostNetwork
	Service             ServiceNetwork
	Ledger              Ledger
	Log                 Log
	Metrics             Metrics
	LogicRunner         LogicRunner
	APIRunner           APIRunner
	AdminAPIRunner      APIRunner
	AvailabilityChecker AvailabilityChecker
	KeysPath            string
	CertificatePath     string
	Tracer              Tracer
	Introspection       Introspection
	Exporter            Exporter
	Bus                 Bus
}

func (c Configuration) GetConfig() interface{} {
	return &c
}

// PulsarConfiguration contains configuration params for the pulsar node
type PulsarConfiguration struct {
	Log      Log
	Pulsar   Pulsar
	Tracer   Tracer
	KeysPath string
	Metrics  Metrics
}

// Holder provides methods to manage configuration
type Holder struct {
	Configuration Configuration
	Params        insconfig.Params
}

// NewConfiguration creates new default configuration
func NewConfiguration() Configuration {
	cfg := Configuration{
		Host:                NewHostNetwork(),
		Service:             NewServiceNetwork(),
		Ledger:              NewLedger(),
		Log:                 NewLog(),
		Metrics:             NewMetrics(),
		LogicRunner:         NewLogicRunner(),
		APIRunner:           NewAPIRunner(false),
		AdminAPIRunner:      NewAPIRunner(true),
		AvailabilityChecker: NewAvailabilityChecker(),
		KeysPath:            "./",
		CertificatePath:     "",
		Tracer:              NewTracer(),
		Introspection:       NewIntrospection(),
		Exporter:            NewExporter(),
		Bus:                 NewBus(),
	}

	return cfg
}

// NewPulsarConfiguration creates new default configuration for pulsar
func NewPulsarConfiguration() PulsarConfiguration {
	return PulsarConfiguration{
		Log:      NewLog(),
		Pulsar:   NewPulsar(),
		Tracer:   NewTracer(),
		KeysPath: "./",
		Metrics:  NewMetrics(),
	}
}

type stringPathGetter struct {
	Path string
}

func (g *stringPathGetter) GetConfigPath() string {
	return g.Path
}

// MustLoad wrapper around Load function which panics on error.
func (h *Holder) MustLoad() *Holder {
	err := h.Load()
	if err != nil {
		panic(err)
	}
	return h
}

// NewHolder creates new Holder with config path
func NewHolder(path string) *Holder {
	params := insconfig.Params{
		ConfigStruct:     Configuration{},
		EnvPrefix:        InsolarEnvPrefix,
		ConfigPathGetter: &stringPathGetter{Path: path},
		FileRequired:     false,
	}
	holder := &Holder{Configuration: Configuration{}, Params: params}

	GlobalHolder = holder
	return holder
}

// Returns global holder if exists
func NewHolderGlobal(path string) *Holder {
	// todo refactor this
	if GlobalHolder != nil {
		return GlobalHolder
	}
	return NewHolder(path)
}

// Load method reads configuration from params file path
func (h *Holder) Load() error {
	insConfigurator := insconfig.NewInsConfigurator(h.Params)
	parsedConf, err := insConfigurator.Load()
	if err != nil {
		return err
	}
	cfg := parsedConf.(*Configuration)
	h.Configuration = *cfg
	return nil
}

// Deprecated
// ToString converts any configuration struct to yaml string
func ToString(in interface{}) string {
	d, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(d)
}
