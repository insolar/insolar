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
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	stdlog "log"
	"os"
	"strings"
	"sync"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
)

const timestampFormat = "2006-01-02T15:04:05.000000000Z07:00"

const defaultCallerSkipFrameCount = 1

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

// NewLog creates logger instance with particular configuration
func NewLog(cfg configuration.Log) (insolar.Logger, error) {
	return NewLogExt(cfg, 0)
}

// NewLog creates logger instance with particular configuration
func NewLogExt(cfg configuration.Log, skipFrameBaselineDelta int) (insolar.Logger, error) {
	pCfg, err := parseLogConfig(cfg)
	if err == nil {
		var logger insolar.Logger

		pCfg.SkipFrameBaselineDelta = skipFrameBaselineDelta
		switch strings.ToLower(cfg.Adapter) {
		case "zerolog":
			logger, err = newZerologAdapter(cfg, pCfg)
		default:
			err = errors.New("unknown adapter")
		}

		if err == nil {
			if logger != nil {
				return logger, nil
			}
			return nil, errors.New("logger was not initialized")
		}
	}
	return nil, errors.Wrap(err, "invalid logger config")
}

var globalLogger = struct {
	mutex  sync.RWMutex
	output ProxyWriter
	logger insolar.Logger
}{}

func GlobalLogger() insolar.Logger {
	globalLogger.mutex.RLock()
	l := globalLogger.logger
	globalLogger.mutex.RUnlock()

	if l != nil {
		return l
	}

	globalLogger.mutex.Lock()
	defer globalLogger.mutex.Unlock()
	if globalLogger.logger == nil {
		createGlobalLogger()
	}

	return globalLogger.logger
}

func SaveGlobalLogger() func() {
	globalLogger.mutex.RLock()
	defer globalLogger.mutex.RUnlock()

	loggerCopy := globalLogger.logger
	outputCopy := globalLogger.output.getTarget()

	return func() {
		globalLogger.mutex.Lock()
		defer globalLogger.mutex.Unlock()

		globalLogger.logger = loggerCopy
		globalLogger.output.setTarget(outputCopy)
	}
}

const globalLoggerCallerSkipFrameDelta = 1

// GlobalLogger creates global logger with correct skipCallNumber
func createGlobalLogger() {
	holder := configuration.NewHolder().MustInit(false)
	logCfg := holder.Configuration.Log
	logCfg.BufferSize = 0 // enforce buffer-less for a non-configured logger

	logger, err := NewLog(logCfg)

	if err == nil {
		err = setGlobalLogger(logger, true)
	}

	if err != nil || logger == nil {
		stdlog.Println("warning: ", err)
		panic("unable to initialize global logger with default config")
	}
}

func setGlobalLogger(logger insolar.Logger, isDefault bool) error {
	b := logger.Copy()
	globalLogger.output.setTarget(b.GetOutput())
	b = b.WithOutput(&globalLogger.output)

	b = b.WithSkipFrameCount(globalLoggerCallerSkipFrameDelta)
	if isDefault {
		b = b.WithCaller(insolar.CallerField)
	}

	var err error
	logger, err = b.Build()
	if err != nil {
		return err
	}
	if isDefault {
		logger = logger.WithField("loginstance", "global_default")
	} else {
		logger = logger.WithField("loginstance", "global")
	}

	globalLogger.logger = logger
	return nil
}

func SetGlobalLogger(logger insolar.Logger) {
	globalLogger.mutex.Lock()
	defer globalLogger.mutex.Unlock()

	if globalLogger.logger == logger {
		return
	}

	err := setGlobalLogger(logger, false)

	if err != nil || logger == nil {
		stdlog.Println("warning: ", err)
		panic("unable to update global logger")
	}
}

// SetLevel lets log level for global logger
func SetLevel(level string) error {
	lvl, err := insolar.ParseLevel(level)
	if err != nil {
		return err
	}

	globalLogger.mutex.Lock()
	defer globalLogger.mutex.Unlock()

	if globalLogger.logger == nil {
		createGlobalLogger()
	}

	globalLogger.logger = globalLogger.logger.Level(lvl)
	return nil
}

// Debug logs a message at level Debug to the global logger.
func Debug(args ...interface{}) {
	GlobalLogger().Debug(args...)
}

// Debugf logs a message at level Debug to the global logger.
func Debugf(format string, args ...interface{}) {
	GlobalLogger().Debugf(format, args...)
}

// Info logs a message at level Info to the global logger.
func Info(args ...interface{}) {
	GlobalLogger().Info(args...)
}

// Infof logs a message at level Info to the global logger.
func Infof(format string, args ...interface{}) {
	GlobalLogger().Infof(format, args...)
}

// Warn logs a message at level Warn to the global logger.
func Warn(args ...interface{}) {
	GlobalLogger().Warn(args...)
}

// Warnf logs a message at level Warn to the global logger.
func Warnf(format string, args ...interface{}) {
	GlobalLogger().Warnf(format, args...)
}

// Error logs a message at level Error to the global logger.
func Error(args ...interface{}) {
	GlobalLogger().Error(args...)
}

// Errorf logs a message at level Error to the global logger.
func Errorf(format string, args ...interface{}) {
	GlobalLogger().Errorf(format, args...)
}

// Fatal logs a message at level Fatal to the global logger.
func Fatal(args ...interface{}) {
	GlobalLogger().Fatal(args...)
}

// Fatalf logs a message at level Fatal to the global logger.
func Fatalf(format string, args ...interface{}) {
	GlobalLogger().Fatalf(format, args...)
}

// Panic logs a message at level Panic to the global logger.
func Panic(args ...interface{}) {
	GlobalLogger().Panic(args...)
}

// Panicf logs a message at level Panic to the global logger.
func Panicf(format string, args ...interface{}) {
	GlobalLogger().Panicf(format, args...)
}
