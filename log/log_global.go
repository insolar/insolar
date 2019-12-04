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
	stdlog "log"
	"sync"
	"time"

	"github.com/insolar/insolar/log/logcommon"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/log/critlog"
)

// Creates and sets global logger. It has a different effect than SetGlobalLogger(NewLog(...)) as it sets a global filter also.
func NewGlobalLogger(cfg configuration.Log) (logcommon.Logger, error) {
	logger, err := NewLog(cfg)
	if err != nil {
		return nil, err
	}

	b := logger.Copy()
	a := getGlobalLogAdapter(b)
	if a == nil {
		return nil, errors.New("Log adapter has no global filter")
	}

	globalLogger.mutex.Lock()
	defer globalLogger.mutex.Unlock()
	err = setGlobalLogger(logger, false, true)
	if err != nil {
		return nil, err
	}

	return globalLogger.logger, nil
}

var globalLogger = struct {
	mutex   sync.RWMutex
	output  critlog.ProxyLoggerOutput
	logger  logcommon.Logger
	adapter logcommon.GlobalLogAdapter
}{}

func GlobalLogger() logcommon.Logger {
	globalLogger.mutex.RLock()
	l := globalLogger.logger
	globalLogger.mutex.RUnlock()

	if l != nil {
		return l
	}

	globalLogger.mutex.Lock()
	defer globalLogger.mutex.Unlock()
	if globalLogger.logger == nil {
		createNoConfigGlobalLogger()
	}

	return globalLogger.logger
}

func CopyGlobalLoggerForContext() logcommon.Logger {
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
		tickLogger, err := GlobalLogger().Copy().WithCaller(logcommon.NoCallerField).Build()
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
func createNoConfigGlobalLogger() {
	holder := configuration.NewHolder().MustInit(false)
	logCfg := holder.Configuration.Log

	// enforce buffer-less for a non-configured logger
	logCfg.BufferSize = 0
	logCfg.LLBufferSize = -1

	logger, err := NewLog(logCfg)

	if err == nil {
		err = setGlobalLogger(logger, true, true)
	}

	if err != nil || logger == nil {
		stdlog.Println("warning: ", err)
		panic("unable to initialize global logger with default config")
	}
}

func getGlobalLogAdapter(b logcommon.LoggerBuilder) logcommon.GlobalLogAdapter {
	if f, ok := b.(logcommon.GlobalLogAdapterFactory); ok {
		return f.CreateGlobalLogAdapter()
	}
	return nil
}

func setGlobalLogger(logger logcommon.Logger, isDefault, isNewGlobal bool) error {
	b := logger.Copy()

	output := b.(logcommon.LoggerOutputGetter).GetLoggerOutput()
	b = b.WithOutput(&globalLogger.output)

	if isDefault {
		// TODO move to logger construction configuration?
		b = b.WithCaller(logcommon.CallerField)
		b = b.WithField("loginstance", "global_default")
	} else {
		b = b.WithField("loginstance", "global")
	}

	adapter := getGlobalLogAdapter(b)
	lvl := b.GetLogLevel()

	var err error
	logger, err = b.Build()
	switch {
	case err != nil:
		return err
	case adapter == nil:
		break
	case isNewGlobal:
		adapter.SetGlobalLoggerFilter(lvl)
		logger = logger.Level(logcommon.DebugLevel)
	case globalLogger.adapter != adapter && globalLogger.adapter != nil:
		adapter.SetGlobalLoggerFilter(globalLogger.adapter.GetGlobalLoggerFilter())
	}

	globalLogger.adapter = adapter
	globalLogger.logger = logger

	globalLogger.output.SetTarget(output)
	return nil
}

func SetGlobalLogger(logger logcommon.Logger) {
	globalLogger.mutex.Lock()
	defer globalLogger.mutex.Unlock()

	if globalLogger.logger == logger {
		return
	}

	err := setGlobalLogger(logger, false, false)

	if err != nil || logger == nil {
		stdlog.Println("warning: ", err)
		panic("unable to update global logger")
	}
}

// SetLevel lets log level for global logger
func SetLevel(level string) error {
	lvl, err := logcommon.ParseLevel(level)
	if err != nil {
		return err
	}

	SetLogLevel(lvl)
	return nil
}

func SetLogLevel(level logcommon.LogLevel) {
	globalLogger.mutex.Lock()
	defer globalLogger.mutex.Unlock()

	if globalLogger.logger == nil {
		createNoConfigGlobalLogger()
	}

	globalLogger.logger = globalLogger.logger.Level(level)
}

func SetGlobalLevelFilter(level logcommon.LogLevel) error {
	globalLogger.mutex.RLock()
	defer globalLogger.mutex.RUnlock()

	return setGlobalLevelFilter(level)
}

func setGlobalLevelFilter(level logcommon.LogLevel) error {
	if globalLogger.adapter != nil {
		globalLogger.adapter.SetGlobalLoggerFilter(level)
		return nil
	}
	return errors.New("not supported")
}

func GetGlobalLevelFilter() logcommon.LogLevel {
	globalLogger.mutex.RLock()
	defer globalLogger.mutex.RUnlock()

	_, l := getGlobalLevelFilter()
	return l
}

func getGlobalLevelFilter() (bool, logcommon.LogLevel) {
	if globalLogger.adapter != nil {
		return true, globalLogger.adapter.GetGlobalLoggerFilter()
	}
	return false, logcommon.NoLevel
}

/*
We use EmbeddedLog functions here to avoid SkipStackFrame corrections
*/

func g() logcommon.EmbeddedLogger {
	return GlobalLogger().Embeddable()
}

func Event(level logcommon.LogLevel, args ...interface{}) {
	if fn := g().NewEvent(level); fn != nil {
		fn(args)
	}
}

func Eventf(level logcommon.LogLevel, fmt string, args ...interface{}) {
	if fn := g().NewEventFmt(level); fn != nil {
		fn(fmt, args)
	}
}

func Debug(args ...interface{}) {
	if fn := g().NewEvent(logcommon.DebugLevel); fn != nil {
		fn(args)
	}
}

func Debugf(fmt string, args ...interface{}) {
	if fn := g().NewEventFmt(logcommon.DebugLevel); fn != nil {
		fn(fmt, args)
	}
}

func Info(args ...interface{}) {
	if fn := g().NewEvent(logcommon.InfoLevel); fn != nil {
		fn(args)
	}
}

func Infof(fmt string, args ...interface{}) {
	if fn := g().NewEventFmt(logcommon.InfoLevel); fn != nil {
		fn(fmt, args)
	}
}

func Warn(args ...interface{}) {
	if fn := g().NewEvent(logcommon.WarnLevel); fn != nil {
		fn(args)
	}
}

func Warnf(fmt string, args ...interface{}) {
	if fn := g().NewEventFmt(logcommon.WarnLevel); fn != nil {
		fn(fmt, args)
	}
}

func Error(args ...interface{}) {
	if fn := g().NewEvent(logcommon.ErrorLevel); fn != nil {
		fn(args)
	}
}

func Errorf(fmt string, args ...interface{}) {
	if fn := g().NewEventFmt(logcommon.ErrorLevel); fn != nil {
		fn(fmt, args)
	}
}

func Fatal(args ...interface{}) {
	if fn := g().NewEvent(logcommon.FatalLevel); fn != nil {
		fn(args)
	}
}

func Fatalf(fmt string, args ...interface{}) {
	if fn := g().NewEventFmt(logcommon.FatalLevel); fn != nil {
		fn(fmt, args)
	}
}

func Panic(args ...interface{}) {
	if fn := g().NewEvent(logcommon.PanicLevel); fn != nil {
		fn(args)
	}
}

func Panicf(fmt string, args ...interface{}) {
	if fn := g().NewEventFmt(logcommon.PanicLevel); fn != nil {
		fn(fmt, args)
	}
}

func Flush() {
	g().EmbeddedFlush("Global logger flush")
}
