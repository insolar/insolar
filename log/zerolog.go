/*
 *    Copyright 2019 Insolar
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
	"os"
	"strings"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type zerologAdapter struct {
	skipCallNumber int
	logger         zerolog.Logger
	logLevel       string
}

func newZerologAdapter(cfg configuration.Log) (*zerologAdapter, error) {

	var output io.Writer
	switch strings.ToLower(cfg.Formatter) {
	case "text":
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: timestampFormat}
	case "json":
		output = os.Stdout
	default:
		return nil, errors.New("unknown formatter " + cfg.Formatter)
	}

	zerolog.CallerSkipFrameCount = 4
	return &zerologAdapter{logger: zerolog.New(output).With().Timestamp().Caller().Logger()}, nil
}

// WithFields return copy of adapter with predefined fields.
func (z zerologAdapter) WithFields(fields map[string]interface{}) core.Logger {
	panic("not supported")
}

// WithField return copy of adapter with predefined single field.
func (z zerologAdapter) WithField(key string, value interface{}) core.Logger {
	panic("not supported")
}

// Debug logs a message at level Debug on the stdout.
func (z zerologAdapter) Debug(args ...interface{}) {
	z.logger.Debug().Msg(args[0].(string))
}

// Debugln logs a message at level Debug on the stdout.
func (z zerologAdapter) Debugln(args ...interface{}) {
	z.logger.Debug().Msg(args[0].(string))
}

// Debugf formatted logs a message at level Debug on the stdout.
func (z zerologAdapter) Debugf(format string, args ...interface{}) {
	z.logger.Debug().Msgf(format, args...)
}

// Info logs a message at level Info on the stdout.
func (z zerologAdapter) Info(args ...interface{}) {
	z.logger.Info().Msg(args[0].(string))
}

// Infoln logs a message at level Info on the stdout.
func (z zerologAdapter) Infoln(args ...interface{}) {
	z.logger.Info().Msg(args[0].(string))
}

// Infof formatted logs a message at level Info on the stdout.
func (z zerologAdapter) Infof(format string, args ...interface{}) {
	z.logger.Info().Msgf(format, args...)
}

// Warn logs a message at level Warn on the stdout.
func (z zerologAdapter) Warn(args ...interface{}) {
	z.logger.Warn().Msg(args[0].(string))
}

// Warnln logs a message at level Warn on the stdout.
func (z zerologAdapter) Warnln(args ...interface{}) {
	z.logger.Warn().Msg(args[0].(string))
}

// Warnf formatted logs a message at level Warn on the stdout.
func (z zerologAdapter) Warnf(format string, args ...interface{}) {
	z.logger.Warn().Msgf(format, args...)
}

// Error logs a message at level Error on the stdout.
func (z zerologAdapter) Error(args ...interface{}) {
	z.logger.Error().Msg(args[0].(string))
}

// Errorln logs a message at level Error on the stdout.
func (z zerologAdapter) Errorln(args ...interface{}) {
	z.logger.Error().Msg(args[0].(string))
}

// Errorf formatted logs a message at level Error on the stdout.
func (z zerologAdapter) Errorf(format string, args ...interface{}) {
	z.logger.Error().Msgf(format, args...)
}

// Fatal logs a message at level Fatal on the stdout.
func (z zerologAdapter) Fatal(args ...interface{}) {
	z.logger.Fatal().Msg(args[0].(string))
}

// Fatalln logs a message at level Fatal on the stdout.
func (z zerologAdapter) Fatalln(args ...interface{}) {
	z.logger.Fatal().Msg(args[0].(string))
}

// Fatalf formatted logs a message at level Fatal on the stdout.
func (z zerologAdapter) Fatalf(format string, args ...interface{}) {
	z.logger.Fatal().Msg(args[0].(string))
}

// Panic logs a message at level Panic on the stdout.
func (z zerologAdapter) Panic(args ...interface{}) {
	z.logger.Panic().Msg(args[0].(string))
}

// Panicln logs a message at level Panic on the stdout.
func (z zerologAdapter) Panicln(args ...interface{}) {
	z.logger.Panic().Msg(args[0].(string))
}

// Panicf formatted logs a message at level Panic on the stdout.
func (z zerologAdapter) Panicf(format string, args ...interface{}) {
	z.logger.Panic().Msgf(format, args...)
}

// SetLevel sets log level
func (z zerologAdapter) SetLevel(level string) error {
	l, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		return errors.Wrap(err, "Failed to parse log level")
	}

	z.logLevel = level
	z.logger = z.logger.Level(l)
	return nil
}

// GetLevel returns log level
func (z zerologAdapter) GetLevel() string {
	return z.logLevel
}

// SetOutput sets the output destination for the logger.
func (z zerologAdapter) SetOutput(w io.Writer) {
	z.logger = z.logger.Output(w)
}
