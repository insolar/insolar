/*
 *    Copyright 2018 Insolar
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

package configuration

type ConnectionType string

const (
	TCP ConnectionType = "tcp"
)

func (ct ConnectionType) String() string {
	return string(ct)
}

// Pulsar holds configuration for pulsar node.
type Pulsar struct {
	ConnectionType      ConnectionType
	MainListenerAddress string
	Storage             Storage

	PulseTime                      int32 // ms
	ReceivingSignTimeout           int32 // ms
	ReceivingNumberTimeout         int32 // ms
	ReceivingVectorTimeout         int32 // ms
	ReceivingSignsForChosenTimeout int32 // ms

	Neighbours []PulsarNodeAddress

	NumberDelta uint32

	DistributionTransport Transport
	PulseDistributor      PulseDistributor
}

type PulseDistributor struct {
	BootstrapHosts            []string
	PingRequestTimeout        int32 // ms
	RandomHostsRequestTimeout int32 // ms
	PulseRequestTimeout       int32 // ms
	RandomNodesCount          int
}

type PulsarNodeAddress struct {
	Address        string
	ConnectionType ConnectionType
	PublicKey      string
}

// NewPulsar creates new default configuration for pulsar node.
func NewPulsar() Pulsar {
	return Pulsar{
		MainListenerAddress: "0.0.0.0:18090",

		ConnectionType: TCP,

		PulseTime:              10000,
		ReceivingSignTimeout:   1000,
		ReceivingNumberTimeout: 1000,
		ReceivingVectorTimeout: 1000,

		Neighbours: []PulsarNodeAddress{},
		Storage:    Storage{DataDirectory: "./data/pulsar"},

		NumberDelta: 10,
		DistributionTransport: Transport{
			Protocol:  "UTP",
			Address:   "0.0.0.0:18091",
			BehindNAT: false,
		},
		PulseDistributor: PulseDistributor{
			BootstrapHosts:            []string{"localhost:53837"},
			PingRequestTimeout:        1000,
			RandomHostsRequestTimeout: 1000,
			PulseRequestTimeout:       1000,
			RandomNodesCount:          5,
		},
	}
}
