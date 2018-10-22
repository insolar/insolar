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

package core

import "io"

// Logger is the interface for loggers used in the Insolar components.
type Logger interface {
	// SetLevel sets log level.
	SetLevel(string) error
	// GetLevel gets log level.
	GetLevel() string

	// Debug logs a message at level Debug.
	Debug(...interface{})
	// Debugln logs a message at level Debug.
	Debugln(...interface{})
	// Debugf formatted logs a message at level Debug.
	Debugf(string, ...interface{})

	// Info logs a message at level Info.
	Info(...interface{})
	// Infoln logs a message at level Info.
	Infoln(...interface{})
	// Infof formatted logs a message at level Info.
	Infof(string, ...interface{})

	// Warn logs a message at level Warn.
	Warn(...interface{})
	// Warnln logs a message at level Warn.
	Warnln(...interface{})
	// Warnf formatted logs a message at level Warn.
	Warnf(string, ...interface{})

	// Error logs a message at level Error.
	Error(...interface{})
	// Errorln logs a message at level Error.
	Errorln(...interface{})
	// Errorf formatted logs a message at level Error.
	Errorf(string, ...interface{})

	// Fatal logs a message at level Fatal and than call os.exit().
	Fatal(...interface{})
	// Fatalln logs a message at level Fatal and than call os.exit().
	Fatalln(...interface{})
	// Fatalf formatted logs a message at level Fatal and than call os.exit().
	Fatalf(string, ...interface{})

	// Panic logs a message at level Panic and than call panic().
	Panic(...interface{})
	// Panicln logs a message at level Panic and than call panic().
	Panicln(...interface{})
	// Panicf formatted logs a message at level Panic and than call panic().
	Panicf(string, ...interface{})

	// SetOutput sets the output destination for the logger.
	SetOutput(w io.Writer)
	// WithFields return copy of Logger with predefined fields.
	WithFields(map[string]interface{}) Logger
	// WithField return copy of Logger with predefined single field.
	WithField(string, interface{}) Logger
}
