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

package inslogger

import (
	"context"

	"github.com/insolar/insolar/core"
	logger "github.com/insolar/insolar/log"
)

type loggerKey struct{}

// FromContext returns logger from context.
func FromContext(ctx context.Context) core.Logger {
	return getLogger(ctx)
}

// SetLogger returns context with provided core.Logger,
func SetLogger(ctx context.Context, l core.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

// WithField returns logger with provided field's key value and context with new logger.
func WithField(ctx context.Context, key string, value string) (context.Context, core.Logger) {
	l := getLogger(ctx).WithField(key, value)
	return SetLogger(ctx, l), l
}

func getLogger(ctx context.Context) core.Logger {
	l := ctx.Value(loggerKey{})
	if l == nil {
		return logger.GlobalLogger
	}
	return l.(core.Logger)
}
