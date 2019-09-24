//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
