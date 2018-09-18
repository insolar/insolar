package log

import (
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

// globalLogger creates global logger with correct skipCallNumber
var globalLogger, _ = func() (core.Logger, error) {
	logger := newLogrusAdapter()
	logger.skipCallNumber = defaultSkipCallNumber + 1
	return logger, logger.SetLevel(configuration.NewLog().Level)
}()

// Debug logs a message at level Debug to the global logger.
func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

// Debugln logs a message at level Debug to the global logger.
func Debugln(args ...interface{}) {
	globalLogger.Debugln(args...)
}

// Debugf logs a message at level Debug to the global logger.
func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

// Info logs a message at level Info to the global logger.
func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

// Infoln logs a message at level Info to the global logger.
func Infoln(args ...interface{}) {
	globalLogger.Infoln(args...)
}

// Infof logs a message at level Info to the global logger.
func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

// Warn logs a message at level Warn to the global logger.
func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

// Warnln logs a message at level Warn to the global logger.
func Warnln(args ...interface{}) {
	globalLogger.Warnln(args...)
}

// Warnf logs a message at level Warn to the global logger.
func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

// Error logs a message at level Error to the global logger.
func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

// Errorln logs a message at level Error to the global logger.
func Errorln(args ...interface{}) {
	globalLogger.Errorln(args...)
}

// Errorf logs a message at level Error to the global logger.
func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

// Fatal logs a message at level Fatal to the global logger.
func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

// Fatalln logs a message at level Fatal to the global logger.
func Fatalln(args ...interface{}) {
	globalLogger.Fatalln(args...)
}

// Fatalf logs a message at level Fatal to the global logger.
func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

// Panic logs a message at level Panic to the global logger.
func Panic(args ...interface{}) {
	globalLogger.Panic(args...)
}

// Panicln logs a message at level Panic to the global logger.
func Panicln(args ...interface{}) {
	globalLogger.Panicln(args...)
}

// Panicf logs a message at level Panic to the global logger.
func Panicf(format string, args ...interface{}) {
	globalLogger.Panicf(format, args...)
}
