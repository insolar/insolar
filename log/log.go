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

package log

import (
	"fmt"
	stdlog "log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
)

const timestampFormat = "2006-01-02T15:04:05.000000000Z07:00"

const defaultCallerSkipFrameCount = 3

var fieldsOrder = []string{
	zerolog.TimestampFieldName,
	zerolog.LevelFieldName,
	zerolog.MessageFieldName,
	zerolog.CallerFieldName,
}

var cwd string

func init() {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		cwd = ""
		fmt.Println("couldn't get current working directory: ", err.Error())
	}
}

func formatCaller() zerolog.Formatter {
	return func(i interface{}) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			if len(cwd) > 0 {
				c = strings.TrimPrefix(c, cwd)
				c = strings.TrimPrefix(c, "/")
			}
			c = "file=" + c
		}
		return c
	}
}

// NewLog creates logger instance with particular configuration
func NewLog(cfg configuration.Log) (insolar.Logger, error) {
	var logger insolar.Logger
	var err error

	switch strings.ToLower(cfg.Adapter) {
	case "zerolog":
		logger, err = newZerologAdapter(cfg)
	default:
		err = errors.New("unknown adapter")
	}

	if err == nil {
		logger, err = logger.WithLevel(cfg.Level)
	}

	if err != nil {
		return nil, errors.Wrap(err, "invalid logger config")
	}

	return logger, nil
}

// GlobalLogger creates global logger with correct skipCallNumber
// TODO: make it private again
var GlobalLogger = func() insolar.Logger {
	holder := configuration.NewHolder().MustInit(false)

	logger, err := NewLog(holder.Configuration.Log)
	if err != nil {
		stdlog.Println("warning:", err.Error())
	}
	if logger == nil {
		panic("couldn't initialize global logger with default config")
	}

	return logger.WithCaller(true).WithSkipFrameCount(1).WithField("loginstance", "global_default")
}()

func SetGlobalLogger(logger insolar.Logger) {
	GlobalLogger = logger.WithSkipFrameCount(1).WithField("loginstance", "global")
}

// SetLevel lets log level for global logger
func SetLevel(level string) error {
	newGlobalLogger, err := GlobalLogger.WithLevel(level)
	if err != nil {
		return err
	}
	GlobalLogger = newGlobalLogger
	return nil
}

// Debug logs a message at level Debug to the global logger.
func Debug(args ...interface{}) {
	GlobalLogger.Debug(args...)
}

// Debugf logs a message at level Debug to the global logger.
func Debugf(format string, args ...interface{}) {
	GlobalLogger.Debugf(format, args...)
}

// Info logs a message at level Info to the global logger.
func Info(args ...interface{}) {
	GlobalLogger.Info(args...)
}

// Infof logs a message at level Info to the global logger.
func Infof(format string, args ...interface{}) {
	GlobalLogger.Infof(format, args...)
}

// Warn logs a message at level Warn to the global logger.
func Warn(args ...interface{}) {
	GlobalLogger.Warn(args...)
}

// Warnf logs a message at level Warn to the global logger.
func Warnf(format string, args ...interface{}) {
	GlobalLogger.Warnf(format, args...)
}

// Error logs a message at level Error to the global logger.
func Error(args ...interface{}) {
	GlobalLogger.Error(args...)
}

// Errorf logs a message at level Error to the global logger.
func Errorf(format string, args ...interface{}) {
	GlobalLogger.Errorf(format, args...)
}

// Fatal logs a message at level Fatal to the global logger.
func Fatal(args ...interface{}) {
	GlobalLogger.Fatal(args...)
}

// Fatalf logs a message at level Fatal to the global logger.
func Fatalf(format string, args ...interface{}) {
	GlobalLogger.Fatalf(format, args...)
}

// Panic logs a message at level Panic to the global logger.
func Panic(args ...interface{}) {
	GlobalLogger.Panic(args...)
}

// Panicf logs a message at level Panic to the global logger.
func Panicf(format string, args ...interface{}) {
	GlobalLogger.Panicf(format, args...)
}
