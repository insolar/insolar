/*
 *    Copyright 2019 Insolar
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

import "github.com/prometheus/client_golang/prometheus"

// ConsensusPacketsSent is current consensus packets sent count metric
var ConsensusPacketsSent = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name:      "sent_count",
	Help:      "Current consensus transport packets sent",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
}, []string{"phase"})

// ConsensusPacketsRecv is current consensus packets received count metric
var ConsensusPacketsRecv = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name:      "recv_count",
	Help:      "Current consensus transport packets recv",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
}, []string{"phase"})

// ConsensusDeclinedClaims is current consensus declined claims count metric
var ConsensusDeclinedClaims = prometheus.NewCounter(prometheus.CounterOpts{
	Name:      "declined_claims_count",
	Help:      "Consensus claims declined",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
})

// ConsensusSentSize is current consensus recv packets size count metric
var ConsensusSentSize = prometheus.NewCounter(prometheus.CounterOpts{
	Name:      "sent_bytes",
	Help:      "Consensus received packets size",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
})

// ConsensusRecvSize is current consensus recv packets size count metric
var ConsensusRecvSize = prometheus.NewCounter(prometheus.CounterOpts{
	Name:      "received_bytes",
	Help:      "Consensus received packets size",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
})

// ConsensusFailedCheckProof is current consensus recv packets size count metric
var ConsensusFailedCheckProof = prometheus.NewCounter(prometheus.CounterOpts{
	Name:      "failed_proof_count",
	Help:      "Consensus validate proof fails",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
})

// ConsensusPhase2TimedOuts is a current consensus phase 2 timed out nodes count metric
var ConsensusPhase2TimedOuts = prometheus.NewCounter(prometheus.CounterOpts{
	Name:      "phase2_exec_count",
	Help:      "Timed out nodes on phase 2",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
})

// ConsensusPhase3Exec is current consensus phase 3 execution count metric
var ConsensusPhase3Exec = prometheus.NewCounter(prometheus.CounterOpts{
	Name:      "phase3_exec_count",
	Help:      "Phase 3 execution counter",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
})

// ConsensusPhase21Exec is current consensus phase 21 execution count metric
var ConsensusPhase21Exec = prometheus.NewCounter(prometheus.CounterOpts{
	Name:      "phase21_exec_count",
	Help:      "Phase 21 execution counter",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
})

// ConsensusActiveNodes is active nodes count after consensus
var ConsensusActiveNodes = prometheus.NewGauge(prometheus.GaugeOpts{
	Name:      "active_nodes_count",
	Help:      "Active nodes count after consensus",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
})
