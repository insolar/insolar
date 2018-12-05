package metrics

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

var InsgorundCallsTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name:      "call_contract_method_total",
	Help:      "Total number of calls contracts methods",
	Namespace: insgorundNamespace,
})

var InsgorundContractExecutionTime = prometheus.NewSummary(prometheus.SummaryOpts{
	Name:       "contract_execution_time",
	Help:       "Time spent on execution contract",
	Namespace:  insgorundNamespace,
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
})

func GetInsgorundRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()

	registry.MustRegister(InsgorundCallsTotal)
	registry.MustRegister(InsgorundContractExecutionTime)
	// default system collectors
	registry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), insgorundNamespace))
	registry.MustRegister(prometheus.NewGoCollector())

	return registry
}
