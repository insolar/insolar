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

// ConsensusPacketsSent is current network transport futures count metric
var ConsensusPacketsSent = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name:      "packets sent",
	Help:      "Current consensus transport packets sent",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
}, []string{"packetType"})

// Consensus1PhasePacketsRecv is current network transport futures count metric
var Consensus1PhasePacketsRecv = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name:      "1 phase packets recv",
	Help:      "Current consensus transport packets recv",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
}, []string{"packetType"})

// Consensus2PhasePacketsRecv is current network transport futures count metric
var Consensus2PhasePacketsRecv = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name:      "2 phase packets recv",
	Help:      "Current consensus transport packets recv",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
}, []string{"packetType"})

// Consensus21PhasePacketsRecv is current network transport futures count metric
var Consensus21PhasePacketsRecv = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name:      "21 phase packets recv",
	Help:      "Current consensus transport packets recv",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
}, []string{"packetType"})

// Consensus3PhasePacketsRecv is current network transport futures count metric
var Consensus3PhasePacketsRecv = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name:      "3 phase packets recv",
	Help:      "Current consensus transport packets recv",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
}, []string{"packetType"})

// ConsensusDeclinedClaims is current network transport futures count metric
var ConsensusDeclinedClaims = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name:      "declined claims",
	Help:      "Claim signature check failed",
	Namespace: insolarNamespace,
	Subsystem: "consensus",
}, []string{"packetType"})
