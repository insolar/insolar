/*
 *    Copyright 2018 INS Ecosystem
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
	BuiltIn  BuiltIn
	GoPlugin Goplugin
}

type BuiltIn struct {
}

type Goplugin struct {
	RunnerListen   string // ip:port of ginsider connectivity socket
	RunnerCodePath string // path where ginsider caches code
	MainListen     string // ip:port of main system connectivity socket
	MainCodePath   string // path of main system code cache
}

func NewLogicRunner() LogicRunner {
	return LogicRunner{
		BuiltIn: BuiltIn{},
		GoPlugin: Goplugin{
			RunnerListen:   "127.0.0.1:7777",
			RunnerCodePath: ".../TEMPDIR_CHANGE ME/...",
			MainListen:     "127.0.0.1:7778",
			MainCodePath:   "./testplugins/",
		},
	}
}
