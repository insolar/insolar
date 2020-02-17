// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

/*
Package inslogger contains context helpers for log

Examples:

	// initialize base context with default logger with provided trace id
	ctx, inslog := inslogger.WithTraceField(context.Background(), "TraceID")
	inslog.Warn("warn")

	// get logger from context
	inslog := inslogger.FromContext(ctx)

	// initalize logger (SomeNewLogger() should return insolar.Logger)
	inslogger.SetLogger(ctx, SomeNewLogger())

Hints:

	Use environment variables for log level setup:

	INSOLAR_LOG_LEVEL=debug INSOLAR_LOG_FORMATTER=text go test ./yourpackage/...
*/
package inslogger
