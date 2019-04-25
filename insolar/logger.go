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
	return NoLevel, fmt.Errorf("Unknown Level String: '%s', defaulting to NoLevel", levelStr)
}

// Logger is the interface for loggers used in the Insolar components.
type Logger interface {
	// SetLevel sets log level.
	WithLevel(string) (Logger, error)
	// SetLevelNumber set log level with number
	WithLevelNumber(level LogLevel) (Logger, error)

	// WithCaller switch on/off 'caller' field computation.
	WithCaller(flag bool) Logger
	// WithSkipFrameCount configures skipFrameCount for 'caller' field computation.
	WithSkipFrameCount(skipFrameCount int) Logger
	// WithFuncName switch on/off 'func' field computation.
	WithFuncName(flag bool) Logger

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

	// SetOutput sets the output destination for the logger.
	WithOutput(w io.Writer) Logger
	// WithFields return copy of Logger with predefined fields.
	WithFields(map[string]interface{}) Logger
	// WithField return copy of Logger with predefined single field.
	WithField(string, interface{}) Logger
}
