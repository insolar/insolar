/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

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
	registry.MustRegister(NetworkFutures)
	registry.MustRegister(NetworkConnections)
	registry.MustRegister(NetworkPacketSentTotal)
	registry.MustRegister(NetworkPacketTimeoutTotal)
	registry.MustRegister(NetworkPacketReceivedTotal)
	registry.MustRegister(NetworkParcelReceivedTotal)
	registry.MustRegister(NetworkComplete)
	registry.MustRegister(NetworkSentSize)
	registry.MustRegister(NetworkRecvSize)

	registry.MustRegister(GopluginContractExecutionTime)

	registry.MustRegister(APIContractExecutionTime)

	// consensus metrics
	registry.MustRegister(ConsensusPacketsSent)
	registry.MustRegister(ConsensusPacketsRecv)
	registry.MustRegister(ConsensusDeclinedClaims)
	registry.MustRegister(ConsensusSentSize)
	registry.MustRegister(ConsensusRecvSize)
	registry.MustRegister(ConsensusFailedCheckProof)
	registry.MustRegister(ConsensusPhase2TimedOuts)
	registry.MustRegister(ConsensusPhase3Exec)
	registry.MustRegister(ConsensusPhase21Exec)
	registry.MustRegister(ConsensusActiveNodes)

	return registry
}
