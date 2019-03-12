/*
 *    Copyright 2019 Insolar Technologies
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
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
)

type logrusAdapter struct {
	skipCallNumber int
	entry          *logrus.Entry
}

func newLogrusAdapter(cfg configuration.Log) (*logrusAdapter, error) {
	log := logrus.New()

	var formatter logrus.Formatter

	switch strings.ToLower(cfg.Formatter) {
	case "text":
		formatter = &logrus.TextFormatter{TimestampFormat: timestampFormat}
	case "json":
		formatter = &logrus.JSONFormatter{TimestampFormat: timestampFormat}
	default:
		return nil, errors.New("unknown formatter " + cfg.Formatter)
	}

	log.SetFormatter(formatter)
	return &logrusAdapter{entry: logrus.NewEntry(log), skipCallNumber: defaultSkipCallNumber}, nil
}

// sourced adds a source info fields that contains
// the package, func, file name and line where the logging happened.
func (l logrusAdapter) sourced() *logrus.Entry {
	info := getCallInfo(l.skipCallNumber)
	return l.entry.WithFields(logrus.Fields{
		"package": info.packageName,
		"func":    info.funcName,
		"file":    fmt.Sprintf("%s:%d", info.fileName, info.line),
	})
}

// WithFields return copy of adapter with predefined fields.
func (l logrusAdapter) WithFields(fields map[string]interface{}) core.Logger {
	lcopy := l
	lcopy.entry = l.entry.WithFields(logrus.Fields(fields))
	return lcopy
}

// WithField return copy of adapter with predefined single field.
func (l logrusAdapter) WithField(key string, value interface{}) core.Logger {
	lcopy := l
	lcopy.entry = l.entry.WithField(key, value)
	return lcopy
}

// Debug logs a message at level Debug on the stdout.
func (l logrusAdapter) Debug(args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.DebugLevel) {
		l.sourced().Debug(args...)
	}
}

// Debugf formatted logs a message at level Debug on the stdout.
func (l logrusAdapter) Debugf(format string, args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.DebugLevel) {
		l.sourced().Debugf(format, args...)
	}
}

// Info logs a message at level Info on the stdout.
func (l logrusAdapter) Info(args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.InfoLevel) {
		l.sourced().Info(args...)
	}
}

// Infof formatted logs a message at level Info on the stdout.
func (l logrusAdapter) Infof(format string, args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.InfoLevel) {
		l.sourced().Infof(format, args...)
	}
}

// Warn logs a message at level Warn on the stdout.
func (l logrusAdapter) Warn(args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.WarnLevel) {
		l.sourced().Warn(args...)
	}
}

// Warnf formatted logs a message at level Warn on the stdout.
func (l logrusAdapter) Warnf(format string, args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.WarnLevel) {
		l.sourced().Warnf(format, args...)
	}
}

// Error logs a message at level Error on the stdout.
func (l logrusAdapter) Error(args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.ErrorLevel) {
		l.sourced().Error(args...)
	}
}

// Errorln logs a message at level Error on the stdout.
func (l logrusAdapter) Errorln(args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.ErrorLevel) {
		l.sourced().Errorln(args...)
	}
}

// Errorf formatted logs a message at level Error on the stdout.
func (l logrusAdapter) Errorf(format string, args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.ErrorLevel) {
		l.sourced().Errorf(format, args...)
	}
}

// Fatal logs a message at level Fatal on the stdout.
func (l logrusAdapter) Fatal(args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.FatalLevel) {
		l.sourced().Fatal(args...)
	}
}

// Fatalln logs a message at level Fatal on the stdout.
func (l logrusAdapter) Fatalln(args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.FatalLevel) {
		l.sourced().Fatalln(args...)
	}
}

// Fatalf formatted logs a message at level Fatal on the stdout.
func (l logrusAdapter) Fatalf(format string, args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.FatalLevel) {
		l.sourced().Fatalf(format, args...)
	}
}

// Panic logs a message at level Panic on the stdout.
func (l logrusAdapter) Panic(args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.PanicLevel) {
		l.sourced().Panic(args...)
	}
}

// Panicln logs a message at level Panic on the stdout.
func (l logrusAdapter) Panicln(args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.PanicLevel) {
		l.sourced().Panicln(args...)
	}
}

// Panicf formatted logs a message at level Panic on the stdout.
func (l logrusAdapter) Panicf(format string, args ...interface{}) {
	if l.entry.Logger.IsLevelEnabled(logrus.PanicLevel) {
		l.sourced().Panicf(format, args...)
	}
}

// SetLevel sets log level
func (l logrusAdapter) SetLevel(level string) error {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	l.entry.Logger.Level = lvl
	return nil
}

// SetOutput sets the output destination for the logger.
func (l logrusAdapter) SetOutput(w io.Writer) {
	l.entry.Logger.SetOutput(w)
}

// WithSkipDelta changes current skip stack frames value for underlying logrus adapter
// on delta value. More about skip value is here https://golang.org/pkg/runtime/#Caller.
//
// This is useful than logger methods called not from place they should report,
// like helper functions.
func WithSkipDelta(cl core.Logger, delta int) core.Logger {
	l, ok := cl.(logrusAdapter)
	if !ok {
		return cl
	}
	newskip := l.skipCallNumber + delta
	out := l
	if newskip < 0 {
		newskip = 0
	}
	out.skipCallNumber = newskip
	return out
}
