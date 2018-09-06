/*
 *    Copyright 2018 INS Ecosystem
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
	Transport      Transport
	BootstrapHosts []string
	IsRelay        bool // set if node must be relay explicit
}

// NewHostNetwork creates new default HostNetwork configuration
func NewHostNetwork() HostNetwork {
	// IP address should not be 0.0.0.0!!!
	transport := Transport{Protocol: "UTP", Address: "0.0.0.0:0", BehindNAT: true}
	return HostNetwork{
		Transport:      transport,
		IsRelay:        false,
		BootstrapHosts: make([]string, 0),
	}
}
