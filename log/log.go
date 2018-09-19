/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package log

import (
	"io"
	stdlog "log"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

const defaultSkipCallNumber = 3

// NewLog creates logger instance with particular configuration
func NewLog(cfg configuration.Log) (core.Logger, error) {
	var logger core.Logger
	switch cfg.Adapter {
	case "logrus":
		logger = newLogrusAdapter()
	default:
		return nil, errors.New("invalid logger config")
	}

	err := logger.SetLevel(cfg.Level)
	if err != nil {
		return nil, errors.Wrap(err, "invalid logger config")
	}

	return logger, nil
}

// SetLevel lets log level for global logger
func SetLevel(level string) error {
	return globalLogger.SetLevel(level)
}

// GetLevel lets log level for global logger
func GetLevel() string {
	return globalLogger.GetLevel()
}

// globalLogger creates global logger with correct skipCallNumber
var globalLogger = func() core.Logger {
	logger := newLogrusAdapter()
	logger.skipCallNumber = defaultSkipCallNumber + 1
	holder := configuration.NewHolder().MustInit(false)
	if err := logger.SetLevel(holder.Configuration.Log.Level); err != nil {
		stdlog.Println("warning:", err.Error())
	}
	return logger
}()

// Debug logs a event at level Debug to the global logger.
func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

// Debugln logs a event at level Debug to the global logger.
func Debugln(args ...interface{}) {
	globalLogger.Debugln(args...)
}

// Debugf logs a event at level Debug to the global logger.
func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

// Info logs a event at level Info to the global logger.
func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

// Infoln logs a event at level Info to the global logger.
func Infoln(args ...interface{}) {
	globalLogger.Infoln(args...)
}

// Infof logs a event at level Info to the global logger.
func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

// Warn logs a event at level Warn to the global logger.
func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

// Warnln logs a event at level Warn to the global logger.
func Warnln(args ...interface{}) {
	globalLogger.Warnln(args...)
}

// Warnf logs a event at level Warn to the global logger.
func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

// Error logs a event at level Error to the global logger.
func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

// Errorln logs a event at level Error to the global logger.
func Errorln(args ...interface{}) {
	globalLogger.Errorln(args...)
}

// Errorf logs a event at level Error to the global logger.
func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

// Fatal logs a event at level Fatal to the global logger.
func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

// Fatalln logs a event at level Fatal to the global logger.
func Fatalln(args ...interface{}) {
	globalLogger.Fatalln(args...)
}

// Fatalf logs a event at level Fatal to the global logger.
func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

// Panic logs a event at level Panic to the global logger.
func Panic(args ...interface{}) {
	globalLogger.Panic(args...)
}

// Panicln logs a event at level Panic to the global logger.
func Panicln(args ...interface{}) {
	globalLogger.Panicln(args...)
}

// Panicf logs a event at level Panic to the global logger.
func Panicf(format string, args ...interface{}) {
	globalLogger.Panicf(format, args...)
}

// SetOutput sets the output destination for the logger.
func SetOutput(w io.Writer) {
	globalLogger.SetOutput(w)
}
