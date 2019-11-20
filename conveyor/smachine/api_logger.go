///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package smachine

type Logger interface {

	// This level has a special behavior - by default it is treated as Debug, but can be raised as Info or Warn
	Trace(interface{})
	IsTracing() bool
	GetTracerId() string

	Meter(interface{})

	Debug(interface{})
	Info(interface{})
	Warn(interface{})
	Error(interface{})
	Fatal(interface{})
	Panic(interface{})
}
