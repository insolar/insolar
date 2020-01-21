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
