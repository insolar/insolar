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
