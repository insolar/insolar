// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

// Tracer configures tracer.
type Tracer struct {
	Jaeger JaegerConfig
	// TODO: add SamplingRules configuration
	SamplingRules struct{}
}

// JaegerConfig holds Jaeger settings.
type JaegerConfig struct {
	CollectorEndpoint string
	AgentEndpoint     string
	ProbabilityRate   float64
}

// NewTracer creates new default Tracer configuration.
func NewTracer() Tracer {
	return Tracer{
		Jaeger: JaegerConfig{
			AgentEndpoint:   "",
			ProbabilityRate: 1,
		},
	}
}
