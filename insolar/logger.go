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
)

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

type LogFormat uint8

const (
	TextFormat LogFormat = iota
	JSONFormat
)

const DefaultLogFormat = TextFormat

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

type LogOutput uint8

const (
	StdErrOutput LogOutput = iota
	SysLogOutput
	JournalDOutput
)

const DefaultLogOutput = StdErrOutput

func ParseOutput(outputStr string, defValue LogOutput) (LogOutput, error) {
	switch strings.ToLower(outputStr) {
	case "", "default":
		return defValue, nil
	case StdErrOutput.String():
		return StdErrOutput, nil
	case SysLogOutput.String():
		return SysLogOutput, nil
	case JournalDOutput.String():
		return JournalDOutput, nil
	}
	return defValue, fmt.Errorf("unknown Output: '%s', replaced with '%s'", outputStr, defValue)
}

func (l LogOutput) String() string {
	switch l {
	case StdErrOutput:
		return "stderr"
	case SysLogOutput:
		return "syslog"
	case JournalDOutput:
		return "journald"
	}
	return string(l)
}

type LogLevelReader interface {
	GetLogLevel() LogLevel
}

type LogLevelWriter interface {
	io.Writer
	LogLevelWrite(LogLevel, []byte) (int, error)
}

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

	// Event logs a message with the given level
	Event(level LogLevel, args ...interface{})
	// Eventf formats and logs a message with the given level
	Eventf(level LogLevel, fmt string, args ...interface{})

	// Is() returns true when a message of the given level will get to output
	Is(level LogLevel) bool

	// WithFields return copy of Logger with predefined fields.
	WithFields(map[string]interface{}) Logger
	// WithField return copy of Logger with predefined single field.
	WithField(string, interface{}) Logger

	Copy() LoggerBuilder

	Level(lvl LogLevel) Logger
}

type CallerFieldMode uint8

const (
	NoCallerField CallerFieldMode = iota
	CallerField
	CallerFieldWithFuncName
)

type LoggerBuilder interface {
	GetOutput() io.Writer

	// SetOutput sets the output destination for the logger.
	WithOutput(w io.Writer) LoggerBuilder
	// WithLevel sets log level.
	WithLevel(level LogLevel) LoggerBuilder
	// WithLevel sets log level.
	WithDynamicLevel(level LogLevelReader) LoggerBuilder
	// WithFormat sets logger output format
	WithFormat(format LogFormat) LoggerBuilder

	// WithCaller switch on/off 'caller' field computation.
	WithCaller(mode CallerFieldMode) LoggerBuilder

	// WithSkipFrameCountDelta changes skipFrameCount by delta value (it can be negative).
	WithSkipFrameCount(skipFrameCount int) LoggerBuilder

	//BuildForTimeCritical(bufSize int) (Logger, error)
	Build() (Logger, error)
}
