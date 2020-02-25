// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
