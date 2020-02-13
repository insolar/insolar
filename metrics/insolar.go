// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// GetInsolarRegistry creates and registers Insolar global metrics
func GetInsolarRegistry(nodeRole string) *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registerer := prometheus.WrapRegistererWith(prometheus.Labels{"role": nodeRole}, registry)

	// badger metrics
	registerer.MustRegister(badgerCollector(insolarNamespace))
	// default system collectors
	registerer.MustRegister(prometheus.NewProcessCollector(
		prometheus.ProcessCollectorOpts{Namespace: insolarNamespace},
	))
	registerer.MustRegister(prometheus.NewGoCollector())
	// insolar collectors
	registerer.MustRegister(NetworkFutures)
	registerer.MustRegister(NetworkConnections)
	registerer.MustRegister(NetworkPacketTimeoutTotal)
	registerer.MustRegister(NetworkPacketReceivedTotal)
	registerer.MustRegister(NetworkSentSize)
	registerer.MustRegister(NetworkRecvSize)

	registerer.MustRegister(APIContractExecutionTime)

	return registry
}
