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

// Transport holds transport protocol configuration for HostNetwork
type Transport struct {
	// protocol type UTP or KCP
	Protocol string
	// Address to listen
	Address string
	// if true transport will use network traversal technique(like STUN) to get PublicAddress
	BehindNAT bool
}

// HostNetwork holds configuration for HostNetwork
type HostNetwork struct {
	Transport           Transport
	IsRelay             bool  // set if node must be relay explicit
	InfinityBootstrap   bool  // set true for infinity tries to bootstrap
	MinTimeout          int   // bootstrap timeout min
	MaxTimeout          int   // bootstrap timeout max
	TimeoutMult         int   // bootstrap timout multiplier
	SignMessages        bool  // signing a messages if true
	HandshakeSessionTTL int32 // ms
}

// NewHostNetwork creates new default HostNetwork configuration
func NewHostNetwork() HostNetwork {
	// IP address should not be 0.0.0.0!!!
	transport := Transport{Protocol: "UTP", Address: "127.0.0.1:0", BehindNAT: false}

	return HostNetwork{
		Transport:           transport,
		IsRelay:             false,
		MinTimeout:          1,
		MaxTimeout:          60,
		TimeoutMult:         2,
		InfinityBootstrap:   false,
		SignMessages:        false,
		HandshakeSessionTTL: 5000,
	}
}
