// Copyright 2020 Insolar Network Ltd.
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

package logmetrics

import (
	"context"
	"time"

	"github.com/insolar/insolar/insolar"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var levelContexts = initLevelContexts()

func initLevelContexts() (contexts []context.Context) {
	contexts = make([]context.Context, insolar.LogLevelCount)
	level := insolar.LogLevel(0)
	for i := range contexts {
		var err error
		contexts[i], err = tag.New(context.Background(), tag.Insert(tagLevel, level.String()))
		if err != nil {
			panic(err)
		}
		level++
	}
	return
}

func GetLogLevelContext(level insolar.LogLevel) context.Context {
	if int(level) >= len(levelContexts) {
		return context.Background()
	}
	return levelContexts[level]
}

func mustTagKey(key string) tag.Key {
	k, err := tag.NewKey(key)
	if err != nil {
		panic(err)
	}
	return k
}

var (
	tagLevel = mustTagKey("level")
)

var (
	statLogCalls = stats.Int64(
		"log_calls",
		"number of log calls",
		stats.UnitDimensionless,
	)
	statLogWrites = stats.Int64(
		"log_written",
		"number of log actually written",
		stats.UnitDimensionless,
	)
	statLogSkips = stats.Int64(
		"log_skipped",
		"number of log entries skipped due to overflow",
		stats.UnitDimensionless,
	)
	statLogWriteDelays = stats.Int64(
		"log_write_delays",
		"duration of log writes",
		"ns",
	)
)

func init() {
	tags := []tag.Key{tagLevel}
	err := view.Register(
		&view.View{
			Name:        statLogCalls.Name(),
			Description: statLogCalls.Description(),
			Measure:     statLogCalls,
			Aggregation: view.Count(),
			TagKeys:     tags,
		},
		&view.View{
			Name:        statLogWrites.Name(),
			Description: statLogWrites.Description(),
			Measure:     statLogWrites,
			Aggregation: view.Count(),
			TagKeys:     tags,
		},
		&view.View{
			Name:        statLogSkips.Name(),
			Description: statLogSkips.Description(),
			Measure:     statLogSkips,
			Aggregation: view.Count(),
			TagKeys:     tags,
		},
		&view.View{
			Name:        statLogWriteDelays.Name(),
			Description: statLogWriteDelays.Description(),
			Measure:     statLogWriteDelays,
			Aggregation: view.Distribution(0.0, float64(time.Second)),
			TagKeys:     tags,
		},
	)
	if err != nil {
		panic(err)
	}
}
