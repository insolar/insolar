/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
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
)
