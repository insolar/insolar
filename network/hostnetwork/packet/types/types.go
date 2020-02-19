// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package types

//go:generate stringer -type=PacketType
type PacketType int

const (
	Unknown PacketType = iota
	// RPC is packet type to execute RPC on a remote node.
	RPC
	// Pulse is packet type to receive Pulse from pulsard and resend it on remote nodes.
	Pulse
	// Bootstrap is packet type for the node bootstrap process.
	Bootstrap
	// Authorize is packet type to authorize bootstrapping node on discovery node.
	Authorize
	// Disconnect is packet type to gracefully disconnect from network.
	Disconnect
	// SignCert used to request signature of certificate from another node
	SignCert
	// UpdateSchedule used for fetching pulse history
	UpdateSchedule
	// Reconnect used to notify nodes to reconnect to the bigger network
	Reconnect
)

// RequestID is 64 bit unsigned int request id.
type RequestID uint64
