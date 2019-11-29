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

package insolar

import (
	"fmt"
	"strings"
)

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

	LogLevelCount = int(iota)
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

func WrapEmbeddedLogger(embedded EmbeddedLogger) Logger {
	if embedded == nil {
		panic("illegal value")
	}
	return &loggerStruct{embedded}
}

type Logger = *loggerStruct // TODO get rid of pointer?

type loggerStruct struct {
	embedded EmbeddedLogger
}

func (z loggerStruct) Is(level LogLevel) bool {
	return z.embedded.Is(level)
}

func (z loggerStruct) Copy() LoggerBuilder {
	return z.embedded.Copy()
}

// Deprecated: do not use, or use Builder
func (z loggerStruct) Level(lvl LogLevel) Logger {
	if logger, err := z.embedded.Copy().WithLevel(lvl).Build(); err != nil {
		panic(err)
	} else {
		return logger
	}
}

func (z loggerStruct) WithFields(fields map[string]interface{}) Logger {
	if assist, ok := z.embedded.(EmbeddedLoggerAssistant); ok {
		return assist.WithFields(fields)
	}
	if logger, err := z.embedded.Copy().WithFields(fields).Build(); err != nil {
		panic(err)
	} else {
		return logger
	}
}

func (z loggerStruct) WithField(name string, value interface{}) Logger {
	if assist, ok := z.embedded.(EmbeddedLoggerAssistant); ok {
		return assist.WithField(name, value)
	}
	if logger, err := z.embedded.Copy().WithField(name, value).Build(); err != nil {
		panic(err)
	} else {
		return logger
	}
}

func (z loggerStruct) Embeddable() EmbeddedLogger {
	return z.embedded
}

func (z loggerStruct) Event(level LogLevel, args ...interface{}) {
	if fn := z.embedded.NewEvent(level); fn != nil {
		fn(args)
	}
}

func (z loggerStruct) Eventf(level LogLevel, fmt string, args ...interface{}) {
	if fn := z.embedded.NewEventFmt(level); fn != nil {
		fn(fmt, args)
	}
}

func (z loggerStruct) Debug(args ...interface{}) {
	if fn := z.embedded.NewEvent(DebugLevel); fn != nil {
		fn(args)
	}
}

func (z loggerStruct) Debugf(fmt string, args ...interface{}) {
	if fn := z.embedded.NewEventFmt(DebugLevel); fn != nil {
		fn(fmt, args)
	}
}

func (z loggerStruct) Info(args ...interface{}) {
	if fn := z.embedded.NewEvent(InfoLevel); fn != nil {
		fn(args)
	}
}

func (z loggerStruct) Infof(fmt string, args ...interface{}) {
	if fn := z.embedded.NewEventFmt(InfoLevel); fn != nil {
		fn(fmt, args)
	}
}

func (z loggerStruct) Warn(args ...interface{}) {
	if fn := z.embedded.NewEvent(WarnLevel); fn != nil {
		fn(args)
	}
}

func (z loggerStruct) Warnf(fmt string, args ...interface{}) {
	if fn := z.embedded.NewEventFmt(WarnLevel); fn != nil {
		fn(fmt, args)
	}
}

func (z loggerStruct) Error(args ...interface{}) {
	if fn := z.embedded.NewEvent(ErrorLevel); fn != nil {
		fn(args)
	}
}

func (z loggerStruct) Errorf(fmt string, args ...interface{}) {
	if fn := z.embedded.NewEventFmt(ErrorLevel); fn != nil {
		fn(fmt, args)
	}
}

func (z loggerStruct) Fatal(args ...interface{}) {
	if fn := z.embedded.NewEvent(FatalLevel); fn != nil {
		fn(args)
	}
}

func (z loggerStruct) Fatalf(fmt string, args ...interface{}) {
	if fn := z.embedded.NewEventFmt(FatalLevel); fn != nil {
		fn(fmt, args)
	}
}

func (z loggerStruct) Panic(args ...interface{}) {
	if fn := z.embedded.NewEvent(PanicLevel); fn != nil {
		fn(args)
	}
}

func (z loggerStruct) Panicf(fmt string, args ...interface{}) {
	if fn := z.embedded.NewEventFmt(PanicLevel); fn != nil {
		fn(fmt, args)
	}
}

/*
	This interface provides methods with -1 call levels.
	DO NOT USE directly, otherwise WithCaller() functionality will be broken.
*/
type EmbeddedLogger interface {
	NewEventStruct(level LogLevel) func(interface{})
	NewEvent(level LogLevel) func(args []interface{})
	NewEventFmt(level LogLevel) func(fmt string, args []interface{})

	// Does flushing of an underlying buffer. Implementation and factual output may vary.
	EmbeddedFlush(msg string)

	Is(LogLevel) bool

	Copy() LoggerBuilder
}

type EmbeddedLoggerAssistant interface {
	WithFields(fields map[string]interface{}) Logger
	WithField(name string, value interface{}) Logger
}
