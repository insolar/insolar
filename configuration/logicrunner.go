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

// LogicRunner configuration
type LogicRunner struct {
	// RPCListen - address logic runner binds RPC API to
	RPCListen string
	// RPCProtoco - protocol (network) of above address,
	// e.g. "tcp", "unix"... see `net.Dial`
	RPCProtocol string
	// BuiltIn - configuration of builtin executor
	BuiltIn *BuiltIn
	// GoPlugin - configuration of executor based on Go plugins
	GoPlugin *GoPlugin
	// PulseLRUSize - configuration of size of a pulse's cache
	PulseLRUSize int
}

// BuiltIn configuration, no options at the moment
type BuiltIn struct{}

// GoPlugin configuration
type GoPlugin struct {
	// RunnerListen - address Go plugins executor listens to
	RunnerListen string
	// RunnerProtocol - protocol (network) of above address,
	// e.g. "tcp", "unix"... see `net.Dial`
	RunnerProtocol string
}

// NewLogicRunner - returns default config of the logic runner
func NewLogicRunner() LogicRunner {
	return LogicRunner{
		RPCListen:   "127.0.0.1:7778",
		RPCProtocol: "tcp",
		BuiltIn:     &BuiltIn{},
		GoPlugin: &GoPlugin{
			RunnerListen:   "127.0.0.1:7777",
			RunnerProtocol: "tcp",
		},
		PulseLRUSize: 10,
	}
}
