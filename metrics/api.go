package metrics

import "github.com/prometheus/client_golang/prometheus"

var APIContractExecutionTime = prometheus.NewSummaryVec(prometheus.SummaryOpts{
	Name:       "contract_execution_time",
	Help:       "Time spent on execution contract, measured from API",
	Namespace:  insolarNamespace,
	Subsystem:  "API",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
}, []string{"method", "success"})
