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
	"context"

	"github.com/rs/zerolog"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

type metricsHook struct{}

// Run implements zerolog.Hook.
func (ch *metricsHook) Run(_ *zerolog.Event, level zerolog.Level, _ string) {
	stats.Record(contextWithLogLevel(level), statLogWrites.M(1))
}

// cache contexts with tag in map per log level to avoid context creation
// on every log metrics measurement
var metricContextByLevel = func() map[zerolog.Level]context.Context {
	var levels = []zerolog.Level{
		zerolog.DebugLevel,
		zerolog.InfoLevel,
		zerolog.WarnLevel,
		zerolog.ErrorLevel,
		zerolog.FatalLevel,
		zerolog.PanicLevel,
		zerolog.NoLevel,
	}
	m := map[zerolog.Level]context.Context{}
	for _, l := range levels {
		ctx, err := tag.New(context.Background(), tag.Insert(tagLevel, l.String()))
		if err != nil {
			panic(err)
		}
		m[l] = ctx
	}
	return m
}()

func contextWithLogLevel(level zerolog.Level) context.Context {
	ctx, ok := metricContextByLevel[level]
	if !ok {
		return context.Background()
	}
	return ctx
}
