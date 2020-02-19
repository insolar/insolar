// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
