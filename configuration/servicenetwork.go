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

// ServiceNetwork is configuration for ServiceNetwork.
type ServiceNetwork struct {
	Skip           int // magic number that indicates what delta after last ignored pulse we should wait
	CacheDirectory string
	Consensus      Consensus
}

type Consensus struct {
	// timeouts for all phases measured in fractions of pulse duration
	Phase1Timeout  float64
	Phase2Timeout  float64
	Phase21Timeout float64
	Phase3Timeout  float64
}

// NewServiceNetwork creates a new ServiceNetwork configuration.
func NewServiceNetwork() ServiceNetwork {
	return ServiceNetwork{
		Skip:           10,
		CacheDirectory: "network_cache",
		Consensus:      NewConsensus(),
	}
}

func NewConsensus() Consensus {
	return Consensus{
		Phase1Timeout:  0.3,
		Phase2Timeout:  0.05,
		Phase21Timeout: 0.05,
		Phase3Timeout:  0.05,
	}
}
