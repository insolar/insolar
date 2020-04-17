// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

const InsolarEnvPrefix string = "insolar"

// GenericConfiguration contains configuration params for all Insolar components
type GenericConfiguration struct {
	Host                HostNetwork
	Service             ServiceNetwork
	Log                 Log
	Metrics             Metrics
	APIRunner           APIRunner
	AdminAPIRunner      APIRunner
	AvailabilityChecker AvailabilityChecker
	KeysPath            string
	CertificatePath     string
	Tracer              Tracer
	Introspection       Introspection
	Bus                 Bus

	// LightChainLimit is maximum pulse difference (NOT number of pulses)
	// between current and the latest replicated on heavy.
	//
	// IMPORTANT: It should be the same on ALL nodes.
	LightChainLimit int
}

func (c GenericConfiguration) GetConfig() interface{} {
	return &c
}

// Holds GenericConfiguration + node specific config
type ConfigHolder interface {
	// Returns Generic Config struct
	GetGenericConfig() GenericConfiguration
	// Returns Node specific config struct
	GetNodeConfig() interface{}
}

// PulsarConfiguration contains configuration params for the pulsar node
type PulsarConfiguration struct {
	Log      Log
	Pulsar   Pulsar
	Tracer   Tracer
	KeysPath string
	Metrics  Metrics
}

// NewGenericConfiguration creates new default configuration
func NewGenericConfiguration() GenericConfiguration {
	cfg := GenericConfiguration{
		Host:                NewHostNetwork(),
		Service:             NewServiceNetwork(),
		Log:                 NewLog(),
		Metrics:             NewMetrics(),
		APIRunner:           NewAPIRunner(false),
		AdminAPIRunner:      NewAPIRunner(true),
		AvailabilityChecker: NewAvailabilityChecker(),
		KeysPath:            "./",
		CertificatePath:     "",
		Tracer:              NewTracer(),
		Introspection:       NewIntrospection(),
		Bus:                 NewBus(),
		LightChainLimit:     5, // 5 pulses
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

// ToString converts any configuration struct to yaml string
func ToString(in interface{}) string {
	d, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(d)
}
