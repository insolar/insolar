package metrics

import "github.com/prometheus/client_golang/prometheus"

var InsgorundContractExecutionTime = prometheus.NewSummary(prometheus.SummaryOpts{
	Name:       "contract_execution_time",
	Help:       "Time spent on execution contract, measured in goplugin",
	Namespace:  insgorundNamespace,
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.05, 0.95: 0.05, 0.99: 0.05},
})
