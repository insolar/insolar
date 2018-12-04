package metrics

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

// NetworkMessageSentTotal is total number of sent messages metric
var InsgorundCallsTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name:      "call_contract_method_total",
	Help:      "Total number of calls contracts methods",
	Namespace: insgorundNamespace,
})

func GetInsgorundRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()

	registry.MustRegister(InsgorundCallsTotal)
	// default system collectors
	registry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), insgorundNamespace))
	registry.MustRegister(prometheus.NewGoCollector())

	return registry
}
