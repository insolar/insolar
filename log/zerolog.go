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
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
)

type zerologAdapter struct {
	logger zerolog.Logger
}

type loglevelChangeHandler struct {
}

func NewLoglevelChangeHandler() *loglevelChangeHandler {
	handler := &loglevelChangeHandler{}
	return handler
}

// ServeHTTP is an HTTP handler that changes the global minimum log level
func (h *loglevelChangeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	levelStr := "(nil)"
	if values["level"] != nil {
		levelStr = values["level"][0]
	}
	level, err := insolar.ParseLevel(levelStr)
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprintf(w, "Invalid level '%v': %v\n", levelStr, err)
		return
	}

	zlevel, err := InternalLevelToZerologLevel(level)
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprintf(w, "Invalid level '%v': %v\n", levelStr, err)
		return
	}

	zerolog.SetGlobalLevel(zlevel)

	w.WriteHeader(200)
	_, _ = fmt.Fprintf(w, "New log level: '%v'\n", levelStr)
}

func InternalLevelToZerologLevel(level insolar.LogLevel) (zerolog.Level, error) {
	switch level {
	case insolar.DebugLevel:
		return zerolog.DebugLevel, nil
	case insolar.InfoLevel:
		return zerolog.InfoLevel, nil
	case insolar.WarnLevel:
		return zerolog.WarnLevel, nil
	case insolar.ErrorLevel:
		return zerolog.ErrorLevel, nil
	case insolar.FatalLevel:
		return zerolog.FatalLevel, nil
	case insolar.PanicLevel:
		return zerolog.PanicLevel, nil
	}
	return zerolog.NoLevel, errors.New("Unknown internal level")
}

func newZerologAdapter(cfg configuration.Log) (*zerologAdapter, error) {
	var output io.Writer
	switch strings.ToLower(cfg.Formatter) {
	case "text":
		output = zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true, TimeFormat: timestampFormat}
	case "json":
		output = os.Stderr
	default:
		return nil, errors.New("unknown formatter " + cfg.Formatter)
	}

	zerolog.CallerSkipFrameCount = 3
	logger := zerolog.New(output).Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()
	return &zerologAdapter{logger}, nil
}

// WithFields return copy of adapter with predefined fields.
func (z *zerologAdapter) WithFields(fields map[string]interface{}) insolar.Logger {
	w := z.logger.With()
	for key, value := range fields {
		w = w.Interface(key, value)
	}
	return &zerologAdapter{w.Logger()}
}

// WithField return copy of adapter with predefined single field.
func (z *zerologAdapter) WithField(key string, value interface{}) insolar.Logger {
	return &zerologAdapter{z.logger.With().Interface(key, value).Logger()}
}

// Debug logs a message at level Debug on the stdout.
func (z *zerologAdapter) Debug(args ...interface{}) {
	z.logger.Debug().Msg(fmt.Sprint(args...))
}

// Debugf formatted logs a message at level Debug on the stdout.
func (z *zerologAdapter) Debugf(format string, args ...interface{}) {
	z.logger.Debug().Msgf(format, args...)
}

// Info logs a message at level Info on the stdout.
func (z *zerologAdapter) Info(args ...interface{}) {
	z.logger.Info().Msg(fmt.Sprint(args...))
}

// Infof formatted logs a message at level Info on the stdout.
func (z *zerologAdapter) Infof(format string, args ...interface{}) {
	z.logger.Info().Msgf(format, args...)
}

// Warn logs a message at level Warn on the stdout.
func (z *zerologAdapter) Warn(args ...interface{}) {
	z.logger.Warn().Msg(fmt.Sprint(args...))
}

// Warnf formatted logs a message at level Warn on the stdout.
func (z *zerologAdapter) Warnf(format string, args ...interface{}) {
	z.logger.Warn().Msgf(format, args...)
}

// Error logs a message at level Error on the stdout.
func (z *zerologAdapter) Error(args ...interface{}) {
	z.logger.Error().Msg(fmt.Sprint(args...))
}

// Errorf formatted logs a message at level Error on the stdout.
func (z *zerologAdapter) Errorf(format string, args ...interface{}) {
	z.logger.Error().Msgf(format, args...)
}

// Fatal logs a message at level Fatal on the stdout.
func (z *zerologAdapter) Fatal(args ...interface{}) {
	z.logger.Fatal().Msg(fmt.Sprint(args...))
}

// Fatalf formatted logs a message at level Fatal on the stdout.
func (z *zerologAdapter) Fatalf(format string, args ...interface{}) {
	z.logger.Fatal().Msgf(format, args...)
}

// Panic logs a message at level Panic on the stdout.
func (z *zerologAdapter) Panic(args ...interface{}) {
	z.logger.Panic().Msg(fmt.Sprint(args...))
}

// Panicf formatted logs a message at level Panic on the stdout.
func (z zerologAdapter) Panicf(format string, args ...interface{}) {
	z.logger.Panic().Msgf(format, args...)
}

// WithLevel sets log level
func (z *zerologAdapter) WithLevel(level string) (insolar.Logger, error) {
	levelNumber, err := insolar.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	return z.WithLevelNumber(levelNumber)
}

func (z *zerologAdapter) WithLevelNumber(level insolar.LogLevel) (insolar.Logger, error) {
	if level == insolar.NoLevel {
		return z, nil
	}
	zerologLevel, err := InternalLevelToZerologLevel(level)
	if err != nil {
		return nil, err
	}
	zCopy := *z
	zCopy.logger = z.logger.Level(zerologLevel)
	return &zCopy, nil
}

// SetOutput sets the output destination for the logger.
func (z *zerologAdapter) WithOutput(w io.Writer) insolar.Logger {
	zCopy := *z
	zCopy.logger = z.logger.Output(w)
	return &zCopy
}
