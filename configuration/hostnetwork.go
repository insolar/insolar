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

package configuration

// Transport holds transport protocol configuration for HostNetwork
type Transport struct {
	// protocol type
	Protocol string
	// Address to listen
	Address string
	// if not empty - this should be public address of instance (to connect from the "other" side to)
	FixedPublicAddress string
}

// HostNetwork holds configuration for HostNetwork
type HostNetwork struct {
	Transport           Transport
	MinTimeout          int   // bootstrap timeout min
	MaxTimeout          int   // bootstrap timeout max
	TimeoutMult         int   // bootstrap timout multiplier
	SignMessages        bool  // signing a messages if true
	HandshakeSessionTTL int32 // ms
}

// NewHostNetwork creates new default HostNetwork configuration
func NewHostNetwork() HostNetwork {
	// IP address should not be 0.0.0.0!!!
	transport := Transport{Protocol: "TCP", Address: "127.0.0.1:0"}

	return HostNetwork{
		Transport:           transport,
		MinTimeout:          10,
		MaxTimeout:          2000,
		TimeoutMult:         2,
		SignMessages:        false,
		HandshakeSessionTTL: 5000,
	}
}
