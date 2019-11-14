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

// full copy of zerolog functions to work with logging level
// needed to support logging level in packet
type LogLevel uint8

// NoLevel means it should be ignored
const (
	NoLevel LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	maxLogLevel
)

const LogLevelCount = int(maxLogLevel)

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
	//JournalDOutput
)
const DefaultLogOutput = StdErrOutput

const DefaultOutputParallelLimit = 5

//go:generate minimock -i github.com/insolar/insolar/insolar.Logger -o ./ -s _mock.go -g

// Logger is the interface for loggers used in the Insolar components.
type Logger interface {
	// Debug logs a message at level Debug.
	Debug(...interface{})
	// Debugf formatted logs a message at level Debug.
	Debugf(string, ...interface{})

	// Info logs a message at level Info.
	Info(...interface{})
	// Infof formatted logs a message at level Info.
	Infof(string, ...interface{})

	// Warn logs a message at level Warn.
	Warn(...interface{})
	// Warnf logs a message at level Warn.
	Warnf(string, ...interface{})

	// Error logs a message at level Error.
	Error(...interface{})
	// Errorf logs a message at level Error.
	Errorf(string, ...interface{})

	// Fatal logs a message at level Fatal and than call os.exit().
	Fatal(...interface{})
	// Fatalf formatted logs a message at level Fatal and than call os.exit().
	Fatalf(string, ...interface{})

	// Panic logs a message at level Panic and than call panic().
	Panic(...interface{})
	// Panicf formatted logs a message at level Panic and than call panic().
	Panicf(string, ...interface{})

	// Event logs a message with the given level.
	Event(level LogLevel, args ...interface{})
	// Eventf formats and logs a message with the given level.
	Eventf(level LogLevel, fmt string, args ...interface{})

	// Is() returns true when a message of the given level will get to output. Considers the global log filter.
	Is(level LogLevel) bool

	// WithFields return copy of Logger with the given fields added. Fields are not deduplicated.
	WithFields(map[string]interface{}) Logger
	// WithField return copy of Logger with the given field added. Fields are not deduplicated.
	WithField(string, interface{}) Logger

	// Provides a builder based on configuration of this logger.
	Copy() LoggerBuilder
	// Provides a copy of this logger with a filter set to lvl.
	Level(lvl LogLevel) Logger

	// DO NOT USE directly. Provides access to an embeddable methods of this logger.
	Embeddable() EmbeddedLogger
}

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

	// WithLevel sets log level. Cancels WithDynamicLevel()
	WithLevel(level LogLevel) LoggerBuilder
	// WithDynamicLevel sets a dynamic log level. Nil value will panic. Resets WithLevel()
	WithDynamicLevel(level LogLevelGetter) LoggerBuilder

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

	// Adds a dynamically-evaluated field. Fields are deduplicated within a single builder only. When func=nil or func()=nil then the field is omitted.
	// NB! Dynamically-evaluated fields are not inherited by derived loggers.
	WithDynamicField(string, func() interface{}) LoggerBuilder

	// Creates a logger.
	Build() (Logger, error)
	// Creates a logger with no write delays.
	BuildLowLatency() (Logger, error)
}

/*
	This interface provides methods with -1 call levels.
	DO NOT USE directly, otherwise WithCaller() functionality will be broken.
*/
type EmbeddedLogger interface {
	// Event logs a message with the given level. DO NOT USE directly.
	EmbeddedEvent(level LogLevel, args ...interface{})
	// Eventf formats and logs a message with the given level. DO NOT USE directly.
	EmbeddedEventf(level LogLevel, fmt string, args ...interface{})
	// Does flushing of an underlying buffer. Implementation and factual output may vary.
	EmbeddedFlush(msg string)
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

type LogLevelGetter interface {
	GetLogLevel() LogLevel
}

type LogObjectWriter interface {
	AddFieldMap(map[string]interface{})
	AddField(key string, v interface{})
	AddRawJSON(key string, b []byte)
}

type LogObjectMarshaller interface {
	MarshalLogObject(LogObjectWriter) string
}

type LogLevelWriter interface {
	io.WriteCloser
	LogLevelWrite(LogLevel, []byte) (int, error)
	Flush() error
}

func (l LogLevel) Equal(other LogLevel) bool {
	return l == other
}

func (l LogLevel) String() string {
	switch l {
	case NoLevel:
		return ""
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}
	return ""
}

func ParseLevel(levelStr string) (LogLevel, error) {
	switch strings.ToLower(levelStr) {
	case NoLevel.String():
		return NoLevel, nil
	case DebugLevel.String():
		return DebugLevel, nil
	case InfoLevel.String():
		return InfoLevel, nil
	case WarnLevel.String():
		return WarnLevel, nil
	case ErrorLevel.String():
		return ErrorLevel, nil
	case FatalLevel.String():
		return FatalLevel, nil
	case PanicLevel.String():
		return PanicLevel, nil
	}
	return NoLevel, fmt.Errorf("unknown Level String: '%s', defaulting to NoLevel", levelStr)
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
		//case JournalDOutput.String():
		//	return JournalDOutput, nil
	}
	return defValue, fmt.Errorf("unknown Output: '%s', replaced with '%s'", outputStr, defValue)
}

func (l LogOutput) String() string {
	switch l {
	case StdErrOutput:
		return "stderr"
	case SysLogOutput:
		return "syslog"
		//case JournalDOutput:
		//	return "journald"
	}
	return string(l)
}
