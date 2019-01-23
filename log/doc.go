/*
 *    Copyright 2019 Insolar
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

/*
Package log contains adapter for third-party loggers

Example:

package main

	import (
		"github.com/insolar/insolar/configuration"
		"github.com/insolar/insolar/log"
	)

	func main() {
		// global logger
		log.SetLevel("Debug")
		log.Debugln("debug log message")

		// local logger
		logger, _ := log.NewLog(configuration.Log{Level: "Warning", Adapter: "logrus"})
		logger.Warnln("warning log message")
	}

*/
package log
