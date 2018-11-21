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

/*
Package inslogger contains context helpers for log

Examples:

	// initialize base context with default logger with provided trace id
	ctx, inslog := inslogger.WithTraceField(context.Background(), "TraceID")
	inslog.Warn("warn")

	// get logger from context
	inslog := inslogger.FromContext(ctx)

	// initalize logger (SomeNewLogger() should return core.Logger)
	inslogger.SetLogger(ctx, SomeNewLogger())

Hints:

	Use environment variables for log level setup:

	INSOLAR_LOG_LEVEL=debug go test ./yourpackage/...
*/
package inslogger
