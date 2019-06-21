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

package inslogger

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/utils"
	logger "github.com/insolar/insolar/log"
)

type loggerKey struct{}
type loggerLevelKey struct{}

func TraceID(ctx context.Context) string {
	return utils.TraceID(ctx)
}

// FromContext returns logger from context.
func FromContext(ctx context.Context) insolar.Logger {
	return getLogger(ctx)
}

// SetLogger returns context with provided insolar.Logger,
func SetLogger(ctx context.Context, l insolar.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

// SetLoggerLevel returns context with provided insolar.LogLevel and set logLevel on logger,
func WithLoggerLevel(ctx context.Context, logLevel insolar.LogLevel) context.Context {
	if logLevel != insolar.NoLevel {
		oldLogger := FromContext(ctx)
		logCopy, err := oldLogger.WithLevelNumber(logLevel)
		if err != nil {
			oldLogger.Error("failed to set log level: ", err.Error())
			return ctx
		}
		ctx = SetLogger(ctx, logCopy)
	}
	return context.WithValue(ctx, loggerLevelKey{}, logLevel)
}

// WithField returns context with logger initialized with provided field's key value and logger itself.
func WithField(ctx context.Context, key string, value string) (context.Context, insolar.Logger) {
	l := getLogger(ctx).WithField(key, value)
	return SetLogger(ctx, l), l
}

// WithFields returns context with logger initialized with provided fields map.
func WithFields(ctx context.Context, fields map[string]interface{}) (context.Context, insolar.Logger) {
	l := getLogger(ctx).WithFields(fields)
	return SetLogger(ctx, l), l
}

// WithTraceField returns context with logger initialized with provided traceid value and logger itself.
func WithTraceField(ctx context.Context, traceid string) (context.Context, insolar.Logger) {
	ctx, err := utils.SetInsTraceID(ctx, traceid)
	if err != nil {
		getLogger(ctx).Error(err)
	}
	return WithField(ctx, "traceid", traceid)
}

// ContextWithTrace returns only context with logger initialized with provided traceid.
func ContextWithTrace(ctx context.Context, traceid string) context.Context {
	ctx, _ = WithTraceField(ctx, traceid)
	return ctx
}

func getLogger(ctx context.Context) insolar.Logger {
	val := ctx.Value(loggerKey{})
	if val == nil {
		return logger.GlobalLogger.WithSkipFrameCount(1)
	}
	l := val.(insolar.Logger)
	ln := GetLoggerLevel(ctx)
	if ln != insolar.NoLevel {
		l, _ = l.WithLevelNumber(ln)
	}
	return l
}

// TestContext returns context with initalized log field "testname" equal t.Name() value.
func TestContext(t *testing.T) context.Context {
	ctx, _ := WithField(context.Background(), "testname", t.Name())
	return ctx
}

func GetLoggerLevel(ctx context.Context) insolar.LogLevel {
	logLevel := ctx.Value(loggerLevelKey{})
	if logLevel == nil {
		return insolar.NoLevel
	}

	logLevelValue, ok := logLevel.(insolar.LogLevel)
	if !ok {
		return insolar.NoLevel
	}

	return logLevelValue
}
