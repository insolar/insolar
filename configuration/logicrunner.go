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

package configuration

type LogicRunner struct {
	RPCListen string // ip:port of main system connectivity socket
	BuiltIn   *BuiltIn
	GoPlugin  *GoPlugin
}

type BuiltIn struct{}

type GoPlugin struct {
	RunnerListen   string // ip:port of ginsider connectivity socket
	RunnerPath     string // path to gincider executable
	RunnerCodePath string // path where ginsider caches code
}

func NewLogicRunner() LogicRunner {
	return LogicRunner{
		RPCListen: "127.0.0.1:7778",
		BuiltIn:   &BuiltIn{},
		GoPlugin: &GoPlugin{
			RunnerPath:     "testdata/logicrunner/insgorund",
			RunnerListen:   "127.0.0.1:7777",
			RunnerCodePath: "",
		},
	}
}
