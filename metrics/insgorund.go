// Copyright 2020 Insolar Network Ltd.
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

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var InsgorundCallsTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name:      "call_contract_method_total",
	Help:      "Total number of calls contracts methods",
	Namespace: insgorundNamespace,
})

var InsgorundContractExecutionTime = prometheus.NewSummaryVec(prometheus.SummaryOpts{
	Name:       "contract_execution_time",
	Help:       "Time spent on execution contract",
	Namespace:  insgorundNamespace,
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
}, []string{"method"})

func GetInsgorundRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()

	registerer := prometheus.WrapRegistererWith(prometheus.Labels{"role": "virtual"}, registry)

	registerer.MustRegister(InsgorundCallsTotal)
	registerer.MustRegister(InsgorundContractExecutionTime)
	// default system collectors
	registerer.MustRegister(prometheus.NewProcessCollector(
		prometheus.ProcessCollectorOpts{Namespace: insgorundNamespace},
	))
	registerer.MustRegister(prometheus.NewGoCollector())

	return registry
}
