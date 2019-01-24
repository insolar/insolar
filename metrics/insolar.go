package metrics

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

// GetInsolarRegistry creates and registers Insolar global metrics
func GetInsolarRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()

	// badger metrics
	registry.MustRegister(badgerCollector(insolarNamespace))
	// default system collectors
	registry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), insolarNamespace))
	registry.MustRegister(prometheus.NewGoCollector())
	// insolar collectors
	registry.MustRegister(NetworkParcelSentTotal)
	registry.MustRegister(NetworkFutures)
	registry.MustRegister(NetworkConnections)
	registry.MustRegister(NetworkPacketSentTotal)
	registry.MustRegister(NetworkPacketTimeoutTotal)
	registry.MustRegister(NetworkPacketReceivedTotal)
	registry.MustRegister(NetworkParcelReceivedTotal)
	registry.MustRegister(NetworkComplete)

	registry.MustRegister(ParcelsSentTotal)
	registry.MustRegister(ParcelsTime)
	registry.MustRegister(ParcelsSentSizeBytes)
	registry.MustRegister(ParcelsReplySizeBytes)
	registry.MustRegister(LocallyDeliveredParcelsTotal)

	registry.MustRegister(GopluginContractExecutionTime)

	registry.MustRegister(APIContractExecutionTime)

	return registry
}
