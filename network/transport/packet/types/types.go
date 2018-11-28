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

package types

//go:generate stringer -type=PacketType
type PacketType int

const (
	// Ping is packet type to ping remote node.
	Ping PacketType = iota + 1
	// RPC is packet type to execute RPC on a remote node.
	RPC
	// Cascade is packet type to send cascade message and execute RPC on each node of the cascade.
	Cascade
	// Pulse is packet type to receive Pulse from pulsard and resend it on remote nodes.
	Pulse
	// GetRandomHosts is packet type for pulsar daemon to get random hosts from insolar daemon.
	GetRandomHosts
	// Bootstrap is packet type for the node bootstrap process.
	Bootstrap
	// Authorize is packet type to authorize bootstrapping node on discovery node.
	Authorize
	// Register is packet type to connect node to discovery node and add join claim to consensus
	Register
	// Genesis is packet type for active list exchange between discovery nodes
	Genesis
	// Challenge1 is packet type for first phase of double challenge response
	Challenge1
	// Challenge2 is packet type for first phase of double challenge response
	Challenge2
	// Disconnect is packet type to gracefully disconnect from network.
	Disconnect

	// Phase1 is packet type for phase 1
	Phase1
	// Phase1 is packet type for phase 2
	Phase2
	// Phase3Pulse is packet type for phase 3 ( pulse )
	Phase3
)
