// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package network

import (
	"time"

	"github.com/insolar/insolar/configuration"
)

// Options contains configuration options for the local host.
type Options struct {
	// The maximum time to wait for a response to any packet.
	PacketTimeout time.Duration

	// The maximum time to wait for a response to ack packet.
	AckPacketTimeout time.Duration

	// Bootstrap ETA for join the Insolar network
	BootstrapTimeout time.Duration

	// Min bootstrap retry timeout
	MinTimeout time.Duration

	// Max bootstrap retry timeout
	MaxTimeout time.Duration

	// Multiplier for boostrap retry time
	TimeoutMult time.Duration

	// The maximum time to wait for a new pulse
	PulseWatchdogTimeout time.Duration
}

// ConfigureOptions convert daemon configuration to controller options
func ConfigureOptions(config configuration.HostNetwork) *Options {
	return &Options{
		TimeoutMult:          time.Duration(config.TimeoutMult) * time.Millisecond,
		MinTimeout:           time.Duration(config.MinTimeout) * time.Millisecond,
		MaxTimeout:           time.Duration(config.MaxTimeout) * time.Millisecond,
		PacketTimeout:        15 * time.Second,
		AckPacketTimeout:     5 * time.Second,
		BootstrapTimeout:     90 * time.Second,
		PulseWatchdogTimeout: 30 * time.Second,
	}
}
