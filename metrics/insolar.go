package metrics

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

func GetInsolarRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()

	// badger metrics
	registry.MustRegister(badgerCollector(insolarNamespace))
	// default system collectors
	registry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), insolarNamespace))
	registry.MustRegister(prometheus.NewGoCollector())
	// insolar collectors
	registry.MustRegister(NetworkMessageSentTotal)
	registry.MustRegister(NetworkFutures)
	registry.MustRegister(NetworkPacketSentTotal)
	registry.MustRegister(NetworkPacketReceivedTotal)

	registry.MustRegister(ParcelsSentTotal)
	registry.MustRegister(ParcelsSentSizeBytes)
	registry.MustRegister(ParcelsReplySizeBytes)
	registry.MustRegister(LocallyDeliveredParcelsTotal)

	registry.MustRegister(GopluginContractExecutionTime)

	return registry
}
