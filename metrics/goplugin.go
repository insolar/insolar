package metrics

import "github.com/prometheus/client_golang/prometheus"

var GopluginContractExecutionTime = prometheus.NewSummary(prometheus.SummaryOpts{
	Name:       "goplugin_contract_execution_time",
	Help:       "Time spent on execution contract, measured in goplugin",
	Namespace:  insolarNamespace,
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
})
