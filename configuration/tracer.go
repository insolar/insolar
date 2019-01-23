/*
 *    Copyright 2019 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

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
