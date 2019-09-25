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
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log/critlog"
	"github.com/insolar/insolar/log/logadapter"
	"github.com/insolar/insolar/log/zlogadapter"
	"github.com/pkg/errors"
	stdlog "log"
	"strings"
	"sync"
	"time"
)

// NewLog creates logger instance with particular configuration
func NewLog(cfg configuration.Log) (insolar.Logger, error) {
	return NewLogExt(cfg, 0)
}

// NewLog creates logger instance with particular configuration
func NewLogExt(cfg configuration.Log, skipFrameBaselineAdjustment int8) (insolar.Logger, error) {
	pCfg, err := logadapter.ParseLogConfig(cfg)
	if err == nil {
		var logger insolar.Logger

		pCfg.SkipFrameBaselineAdjustment = skipFrameBaselineAdjustment

		msgFmt := logadapter.GetDefaultLogMsgFormatter()

		switch strings.ToLower(cfg.Adapter) {
		case "zerolog":
			logger, err = zlogadapter.NewZerologAdapter(pCfg, msgFmt)
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
	mutex   sync.RWMutex
	output  critlog.ProxyLoggerOutput
	logger  insolar.Logger
	adapter insolar.GlobalLogAdapter
}{}

func g() insolar.EmbeddedLogger {
	return GlobalLogger().Embeddable()
}

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

func CopyGlobalLoggerForContext() insolar.Logger {
	return GlobalLogger()
}

func SaveGlobalLogger() func() {
	return SaveGlobalLoggerAndFilter(false)
}

func SaveGlobalLoggerAndFilter(includeFilter bool) func() {
	globalLogger.mutex.RLock()
	defer globalLogger.mutex.RUnlock()

	loggerCopy := globalLogger.logger
	outputCopy := globalLogger.output.GetTarget()
	hasLoggerFilter, loggerFilter := getGlobalLevelFilter()

	return func() {
		globalLogger.mutex.Lock()
		defer globalLogger.mutex.Unlock()

		globalLogger.logger = loggerCopy
		globalLogger.output.SetTarget(outputCopy)

		if includeFilter && hasLoggerFilter {
			_ = setGlobalLevelFilter(loggerFilter)
		}
	}
}

var globalTickerOnce sync.Once

func InitTicker() {
	globalTickerOnce.Do(func() {
		// as we use GlobalLogger() - the copy will follow any redirection made on the GlobalLogger()
		tickLogger, err := GlobalLogger().Copy().WithCaller(insolar.NoCallerField).Build()
		if err != nil {
			panic(err)
		}

		go func() {
			for {
				// Tick between seconds
				time.Sleep(time.Second - time.Since(time.Now().Truncate(time.Second)))
				tickLogger.Debug("Logger tick")
			}
		}()
	})
}

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

func getGlobalLogAdapter(b insolar.LoggerBuilder) insolar.GlobalLogAdapter {
	if f, ok := b.(insolar.GlobalLogAdapterFactory); ok {
		return f.CreateGlobalLogAdapter()
	}
	return nil
}

func setGlobalLogger(logger insolar.Logger, isDefault bool) error {
	b := logger.Copy()

	output := b.(insolar.LoggerOutputGetter).GetLoggerOutput()
	b = b.WithOutput(&globalLogger.output)

	if isDefault {
		// TODO move to logger construction configuration?
		b = b.WithCaller(insolar.CallerField)
	}

	adapter := getGlobalLogAdapter(b)

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

	if globalLogger.adapter != nil && adapter != nil && globalLogger.adapter != adapter {
		adapter.SetGlobalLoggerFilter(globalLogger.adapter.GetGlobalLoggerFilter())
	}
	globalLogger.adapter = adapter
	globalLogger.logger = logger
	globalLogger.output.SetTarget(output)
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

	SetLogLevel(lvl)
	return nil
}

func SetLogLevel(level insolar.LogLevel) {
	globalLogger.mutex.Lock()
	defer globalLogger.mutex.Unlock()

	if globalLogger.logger == nil {
		createGlobalLogger()
	}

	globalLogger.logger = globalLogger.logger.Level(level)
}

func SetGlobalLevelFilter(level insolar.LogLevel) error {
	globalLogger.mutex.RLock()
	defer globalLogger.mutex.RUnlock()

	return setGlobalLevelFilter(level)
}

func setGlobalLevelFilter(level insolar.LogLevel) error {
	if globalLogger.adapter != nil {
		globalLogger.adapter.SetGlobalLoggerFilter(level)
		return nil
	}
	return errors.New("not supported")
}

func GetGlobalLevelFilter() insolar.LogLevel {
	globalLogger.mutex.RLock()
	defer globalLogger.mutex.RUnlock()

	_, l := getGlobalLevelFilter()
	return l
}

func getGlobalLevelFilter() (bool, insolar.LogLevel) {
	if globalLogger.adapter != nil {
		return true, globalLogger.adapter.GetGlobalLoggerFilter()
	}
	return false, insolar.NoLevel
}

/*
We use EmbeddedEvent functions here to avoid SkipStackFrame corrections
*/

func Debug(args ...interface{}) {
	g().EmbeddedEvent(insolar.DebugLevel, args...)
}

func Debugf(format string, args ...interface{}) {
	g().EmbeddedEventf(insolar.DebugLevel, format, args...)
}

func Info(args ...interface{}) {
	g().EmbeddedEvent(insolar.InfoLevel, args...)
}

func Infof(format string, args ...interface{}) {
	g().EmbeddedEventf(insolar.InfoLevel, format, args...)
}

func Warn(args ...interface{}) {
	g().EmbeddedEvent(insolar.WarnLevel, args...)
}

func Warnf(format string, args ...interface{}) {
	g().EmbeddedEventf(insolar.WarnLevel, format, args...)
}

func Error(args ...interface{}) {
	g().EmbeddedEvent(insolar.ErrorLevel, args...)
}

func Errorf(format string, args ...interface{}) {
	g().EmbeddedEventf(insolar.ErrorLevel, format, args...)
}

func Fatal(args ...interface{}) {
	g().EmbeddedEvent(insolar.FatalLevel, args...)
}

func Fatalf(format string, args ...interface{}) {
	g().EmbeddedEventf(insolar.FatalLevel, format, args...)
}

func Panic(args ...interface{}) {
	g().EmbeddedEvent(insolar.PanicLevel, args...)
}

func Panicf(format string, args ...interface{}) {
	g().EmbeddedEventf(insolar.PanicLevel, format, args...)
}

func Event(level insolar.LogLevel, args ...interface{}) {
	g().EmbeddedEvent(level, args...)
}

func Eventf(level insolar.LogLevel, format string, args ...interface{}) {
	g().EmbeddedEventf(level, format, args...)
}
