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

package insolar

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type LogFormat uint8

const (
	TextFormat LogFormat = iota
	JSONFormat
)
const DefaultLogFormat = TextFormat

type LogOutput uint8

const (
	StdErrOutput LogOutput = iota
	SysLogOutput
)
const DefaultLogOutput = StdErrOutput

const DefaultOutputParallelLimit = 5

//go:generate minimock -i github.com/insolar/insolar/insolar.Logger -o ./ -s _mock.go -g

type DynFieldFunc func() interface{}

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

type GlobalLogAdapterFactory interface {
	CreateGlobalLogAdapter() GlobalLogAdapter
}

type GlobalLogAdapter interface {
	SetGlobalLoggerFilter(level LogLevel)
	GetGlobalLoggerFilter() LogLevel
}

type CallerFieldMode uint8

const (
	NoCallerField CallerFieldMode = iota
	CallerField
	CallerFieldWithFuncName
)

type LogMetricsRecorder interface {
	RecordLogEvent(level LogLevel)
	RecordLogWrite(level LogLevel)
	RecordLogDelay(level LogLevel, d time.Duration)
}

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

type LoggerOutputGetter interface {
	GetLoggerOutput() LoggerOutput
}

type LoggerOutput interface {
	LogLevelWriter
	LowLatencyWrite(LogLevel, []byte) (int, error)
	IsLowLatencySupported() bool
}

// Presence of this interface indicates that this object can be used as a log event
type LogObject interface {
	// should return nil to use default (external) marshaller
	GetLogObjectMarshaller() LogObjectMarshaller
}

var _ LogObject = &LogObjectTemplate{}

type LogObjectTemplate struct{}

func (*LogObjectTemplate) GetLogObjectMarshaller() LogObjectMarshaller {
	return nil
}

type LogObjectFields struct {
	Msg    string
	Fields map[string]interface{}
}

func (v LogObjectFields) MarshalLogObject(w LogObjectWriter, _ LogObjectMetricCollector) string {
	for k, v := range v.Fields {
		w.AddIntfField(k, v, LogFieldFormat{})
	}
	return v.Msg
}

type LogObjectMarshaller interface {
	MarshalLogObject(LogObjectWriter, LogObjectMetricCollector) string
}

type MutedLogObjectMarshaller interface {
	MarshalMutedLogObject(LogObjectMetricCollector)
}

type LogObjectMetricCollector interface {
	LogObjectMetricCollector()
	//ReportMetricSample(metricType uint32, reporterFieldName string, value interface{})
}

type LogFieldFormat struct {
	HasFmt bool
	Fmt    string
}

type LogObjectWriter interface {
	AddIntField(key string, v int64, fmt LogFieldFormat)
	AddUintField(key string, v uint64, fmt LogFieldFormat)
	AddFloatField(key string, v float64, fmt LogFieldFormat)
	AddStrField(key string, v string, fmt LogFieldFormat)
	AddIntfField(key string, v interface{}, fmt LogFieldFormat)
	AddRawJSONField(key string, b []byte)
}

type LogLevelWriter interface {
	io.WriteCloser
	LogLevelWrite(LogLevel, []byte) (int, error)
	Flush() error
}

func ParseFormat(formatStr string, defValue LogFormat) (LogFormat, error) {
	switch strings.ToLower(formatStr) {
	case "", "default":
		return defValue, nil
	case TextFormat.String():
		return TextFormat, nil
	case JSONFormat.String():
		return JSONFormat, nil
	}
	return defValue, fmt.Errorf("unknown Format: '%s', replaced with '%s'", formatStr, defValue)
}

func (l LogFormat) String() string {
	switch l {
	case TextFormat:
		return "text"
	case JSONFormat:
		return "json"
	}
	return string(l)
}

func ParseOutput(outputStr string, defValue LogOutput) (LogOutput, error) {
	switch strings.ToLower(outputStr) {
	case "", "default":
		return defValue, nil
	case StdErrOutput.String():
		return StdErrOutput, nil
	case SysLogOutput.String():
		return SysLogOutput, nil
	}
	return defValue, fmt.Errorf("unknown Output: '%s', replaced with '%s'", outputStr, defValue)
}

func (l LogOutput) String() string {
	switch l {
	case StdErrOutput:
		return "stderr"
	case SysLogOutput:
		return "syslog"
	}
	return string(l)
}
