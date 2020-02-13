// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
