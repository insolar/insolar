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

// HostNetwork holds configuration for HostNetwork
type HostNetwork struct {
	Address        string
	PublicAddress  string
	BootstrapHosts []string
	UseStun        bool   // use stun to get public address
	IsRelay        bool   // set if node must be relay explicit
	Transport      string // transport type UTP or KCP
}

// NewHostNetwork creates new default HostNetwork configuration
func NewHostNetwork() HostNetwork {
	return HostNetwork{
		Address:   "0.0.0.0:17000",
		UseStun:   true,
		IsRelay:   false,
		Transport: "UTP",
	}
}
