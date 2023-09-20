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
		logger, _ := log.NewLog(configuration.Log{Level: "Warning", Adapter: "zerolog"})
		logger.Warnln("warning log message")
	}

*/
package log
