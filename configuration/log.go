// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

// Log holds configuration for logging
type Log struct {
	// Default level for logger
	Level string
	//// Default level for global filtering
	//GlobalLevel string
	// Logging adapter - only zerolog by now
	Adapter string
	// Log output format - e.g. json or text
	Formatter string
	// Log output type - e.g. stderr, syslog
	OutputType string
	// Write-parallel limit for the output
	OutputParallelLimit string
	// Parameter for output - depends on OutputType
	OutputParams string

	// Number of regular log events that can be buffered, =0 to disable
	BufferSize int
	// Number of low-latency log events that can be buffered, =-1 to disable, =0 - default size
	LLBufferSize int
}

// NewLog creates new default configuration for logging
func NewLog() Log {
	return Log{
		Level:      "Info",
		Adapter:    "zerolog",
		Formatter:  "json",
		OutputType: "stderr",
		BufferSize: 0,
	}
}
