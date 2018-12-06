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

// GlobalLogger creates global logger with correct skipCallNumber
// TODO: make it private again
var GlobalLogger = func() core.Logger {
	logger := newLogrusAdapter()
	logger.skipCallNumber = defaultSkipCallNumber + 1
	holder := configuration.NewHolder().MustInit(false)
	if err := logger.SetLevel(holder.Configuration.Log.Level); err != nil {
		stdlog.Println("warning:", err.Error())
	}
	return logger
}()

func SetGlobalLogger(logger core.Logger) {
	GlobalLogger = logger.WithField("GLOBAL_OLD_LOG( Kitty, please, use new one )", "")
}

// SetLevel lets log level for global logger
func SetLevel(level string) error {
	return GlobalLogger.SetLevel(level)
}

// GetLevel lets log level for global logger
func GetLevel() string {
	return GlobalLogger.GetLevel()
}

// Debug logs a message at level Debug to the global logger.
func Debug(args ...interface{}) {
	GlobalLogger.Debug(args...)
}

// Debugln logs a message at level Debug to the global logger.
func Debugln(args ...interface{}) {
	GlobalLogger.Debugln(args...)
}

// Debugf logs a message at level Debug to the global logger.
func Debugf(format string, args ...interface{}) {
	GlobalLogger.Debugf(format, args...)
}

// Info logs a message at level Info to the global logger.
func Info(args ...interface{}) {
	GlobalLogger.Info(args...)
}

// Infoln logs a message at level Info to the global logger.
func Infoln(args ...interface{}) {
	GlobalLogger.Infoln(args...)
}

// Infof logs a message at level Info to the global logger.
func Infof(format string, args ...interface{}) {
	GlobalLogger.Infof(format, args...)
}

// Warn logs a message at level Warn to the global logger.
func Warn(args ...interface{}) {
	GlobalLogger.Warn(args...)
}

// Warnln logs a message at level Warn to the global logger.
func Warnln(args ...interface{}) {
	GlobalLogger.Warnln(args...)
}

// Warnf logs a message at level Warn to the global logger.
func Warnf(format string, args ...interface{}) {
	GlobalLogger.Warnf(format, args...)
}

// Error logs a message at level Error to the global logger.
func Error(args ...interface{}) {
	GlobalLogger.Error(args...)
}

// Errorln logs a message at level Error to the global logger.
func Errorln(args ...interface{}) {
	GlobalLogger.Errorln(args...)
}

// Errorf logs a message at level Error to the global logger.
func Errorf(format string, args ...interface{}) {
	GlobalLogger.Errorf(format, args...)
}

// Fatal logs a message at level Fatal to the global logger.
func Fatal(args ...interface{}) {
	GlobalLogger.Fatal(args...)
}

// Fatalln logs a message at level Fatal to the global logger.
func Fatalln(args ...interface{}) {
	GlobalLogger.Fatalln(args...)
}

// Fatalf logs a message at level Fatal to the global logger.
func Fatalf(format string, args ...interface{}) {
	GlobalLogger.Fatalf(format, args...)
}

// Panic logs a message at level Panic to the global logger.
func Panic(args ...interface{}) {
	GlobalLogger.Panic(args...)
}

// Panicln logs a message at level Panic to the global logger.
func Panicln(args ...interface{}) {
	GlobalLogger.Panicln(args...)
}

// Panicf logs a message at level Panic to the global logger.
func Panicf(format string, args ...interface{}) {
	GlobalLogger.Panicf(format, args...)
}

// SetOutput sets the output destination for the logger.
func SetOutput(w io.Writer) {
	GlobalLogger.SetOutput(w)
}
