//
// Copyright 2019 Insolar Technologies GmbH
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
//

package configuration

// Pulsar holds configuration for pulsar node.
type Pulsar struct {
	PulseTime int32 // ms

	NumberDelta uint32

	DistributionTransport Transport
	PulseDistributor      PulseDistributor
}

type PulseDistributor struct {
	BootstrapHosts      []string
	PulseRequestTimeout int32 // ms
}

type PulsarNodeAddress struct {
	Address   string
	PublicKey string
}

// NewPulsar creates new default configuration for pulsar node.
func NewPulsar() Pulsar {
	return Pulsar{
		PulseTime: 10000,

		NumberDelta: 10,
		DistributionTransport: Transport{
			Protocol: "TCP",
			Address:  "0.0.0.0:18091",
		},
		PulseDistributor: PulseDistributor{
			BootstrapHosts:      []string{"localhost:53837"},
			PulseRequestTimeout: 1000,
		},
	}
}
