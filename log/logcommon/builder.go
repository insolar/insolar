//
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
//

package logcommon

import (
	"io"
	"time"
)

type DynFieldFunc func() interface{}
type DynFieldMap map[string]DynFieldFunc

type LoggerBuilder interface {
	// Returns the current output
	GetOutput() io.Writer
	// Returns the current log level
	GetLogLevel() LogLevel

	// Sets the output destination for the logger.
	WithOutput(w io.Writer) LoggerBuilder
	// WithFormat sets logger output format.
	WithFormat(format LogFormat) LoggerBuilder
	// Set buffer size and applicability of the buffer. Will be IGNORED when a reused output is already buffered.
	WithBuffer(bufferSize int, bufferForAll bool) LoggerBuilder

	// WithLevel sets log level.
	WithLevel(level LogLevel) LoggerBuilder

	//// Sets level for active Trace() operations. Parameter can only be Info, Warn or NoLevel (ignores any Trace).
	//WithTracingLevel(LogLevel) LoggerBuilder
	//// Enables remapping of Trace() to the level set by WithTracingLevel
	//WithTracing(bool) LoggerBuilder

	// Controls 'func' and 'caller' field computation. See also WithSkipFrameCount().
	WithCaller(mode CallerFieldMode) LoggerBuilder
	// WithSkipFrameCount changes skipFrameCount to the absolute value. But the value can be negative, and it is applied to a baseline. Value exceeding int8 will panic
	WithSkipFrameCount(skipFrameCount int) LoggerBuilder

	// Controls collection of metrics. Required flags are ADDED to the current flags. Include specify LogMetricsResetMode to replace flags.
	WithMetrics(mode LogMetricsMode) LoggerBuilder
	//Sets an custom recorder for metric collection.
	WithMetricsRecorder(recorder LogMetricsRecorder) LoggerBuilder

	// WithFields adds fields for to-be-built logger. Fields are deduplicated within a single builder only.
	WithFields(map[string]interface{}) LoggerBuilder
	// WithField add a fields for to-be-built logger. Fields are deduplicated within a single builder only.
	WithField(string, interface{}) LoggerBuilder

	// Clears out inherited fields (dynamic or not)
	WithoutInheritedFields() LoggerBuilder
	// Clears out inherited dynamic fields only
	WithoutInheritedDynFields() LoggerBuilder

	// Adds a dynamically-evaluated field. Fields are deduplicated within a single builder only. When func=nil or func()=nil then the field is omitted.
	// NB! Dynamically-evaluated fields are not inherited by derived loggers.
	WithDynamicField(string, DynFieldFunc) LoggerBuilder

	// Creates a logger.
	Build() (Logger, error)
	// Creates a logger with no write delays.
	BuildLowLatency() (Logger, error)
}

type LogMetricsRecorder interface {
	RecordLogEvent(level LogLevel)
	RecordLogWrite(level LogLevel)
	RecordLogDelay(level LogLevel, d time.Duration)
}

type CallerFieldMode uint8

const (
	NoCallerField CallerFieldMode = iota
	CallerField
	CallerFieldWithFuncName
)

type LogMetricsMode uint8

const NoLogMetrics LogMetricsMode = 0
const (
	// Logger will report every event to metrics
	LogMetricsEventCount LogMetricsMode = 1 << iota
	// Logger will report to metrics a write duration (time since an event was created till it was directed to the output)
	LogMetricsWriteDelayReport
	// Logger will add a write duration field into to the output
	LogMetricsWriteDelayField
	// No effect on logger. Indicates that WithMetrics should replace the mode, instead of adding it.
	LogMetricsResetMode
)

type LogFormat uint8

const (
	TextFormat LogFormat = iota
	JSONFormat
)

func (l LogFormat) String() string {
	switch l {
	case TextFormat:
		return "text"
	case JSONFormat:
		return "json"
	}
	return string(l)
}
