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
	"fmt"

	"github.com/sirupsen/logrus"
)

type logrusAdapter struct {
	skipCallNumber int
	entry          *logrus.Entry
}

func newLogrusAdapter() logrusAdapter {
	return logrusAdapter{entry: logrus.NewEntry(logrus.New()), skipCallNumber: defaultSkipCallNumber}
}

// sourced adds a source info fields that contains
// the package, func, file name and line where the logging happened.
func (l logrusAdapter) sourced() *logrus.Entry {
	info := getCallInfo(l.skipCallNumber)
	return l.entry.WithFields(logrus.Fields{"package": info.packageName,
		"func": info.funcName,
		"file": fmt.Sprintf("%s:%d", info.fileName, info.line)})
}

// Debug logs a event at level Debug on the stdout.
func (l logrusAdapter) Debug(args ...interface{}) {
	l.sourced().Debug(args...)
}

// Debugln logs a event at level Debug on the stdout.
func (l logrusAdapter) Debugln(args ...interface{}) {
	l.sourced().Debugln(args...)
}

// Debugf formatted logs a event at level Debug on the stdout.
func (l logrusAdapter) Debugf(format string, args ...interface{}) {
	l.sourced().Debugf(format, args...)
}

// Info logs a event at level Info on the stdout.
func (l logrusAdapter) Info(args ...interface{}) {
	l.sourced().Info(args...)
}

// Infoln logs a event at level Info on the stdout.
func (l logrusAdapter) Infoln(args ...interface{}) {
	l.sourced().Infoln(args...)
}

// Infof formatted logs a event at level Info on the stdout.
func (l logrusAdapter) Infof(format string, args ...interface{}) {
	l.sourced().Infof(format, args...)
}

// Warn logs a event at level Warn on the stdout.
func (l logrusAdapter) Warn(args ...interface{}) {
	l.sourced().Warn(args...)
}

// Warnln logs a event at level Warn on the stdout.
func (l logrusAdapter) Warnln(args ...interface{}) {
	l.sourced().Warnln(args...)
}

// Warnf formatted logs a event at level Warn on the stdout.
func (l logrusAdapter) Warnf(format string, args ...interface{}) {
	l.sourced().Warnf(format, args...)
}

// Error logs a event at level Error on the stdout.
func (l logrusAdapter) Error(args ...interface{}) {
	l.sourced().Error(args...)
}

// Errorln logs a event at level Error on the stdout.
func (l logrusAdapter) Errorln(args ...interface{}) {
	l.sourced().Errorln(args...)
}

// Errorf formatted logs a event at level Error on the stdout.
func (l logrusAdapter) Errorf(format string, args ...interface{}) {
	l.sourced().Errorf(format, args...)
}

// Fatal logs a event at level Fatal on the stdout.
func (l logrusAdapter) Fatal(args ...interface{}) {
	l.sourced().Fatal(args...)
}

// Fatalln logs a event at level Fatal on the stdout.
func (l logrusAdapter) Fatalln(args ...interface{}) {
	l.sourced().Fatalln(args...)
}

// Fatalf formatted logs a event at level Fatal on the stdout.
func (l logrusAdapter) Fatalf(format string, args ...interface{}) {
	l.sourced().Fatalf(format, args...)
}

// Panic logs a event at level Panic on the stdout.
func (l logrusAdapter) Panic(args ...interface{}) {
	l.sourced().Panic(args...)
}

// Panicln logs a event at level Panic on the stdout.
func (l logrusAdapter) Panicln(args ...interface{}) {
	l.sourced().Panicln(args...)
}

// Panicf formatted logs a event at level Panic on the stdout.
func (l logrusAdapter) Panicf(format string, args ...interface{}) {
	l.sourced().Panicf(format, args...)
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
