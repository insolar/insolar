// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

// LogicRunner configuration
type LogicRunner struct {
	// PulseLRUSize - configuration of size of a pulse's cache
	PulseLRUSize int
}

// NewLogicRunner - returns default config of the logic runner
func NewLogicRunner() LogicRunner {
	return LogicRunner{
		PulseLRUSize: 100,
	}
}
