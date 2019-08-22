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
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

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
		"log/calls",
		"number of log calls",
		stats.UnitDimensionless,
	)
	statLogWrites = stats.Int64(
		"log/written",
		"number of log actually written",
		stats.UnitDimensionless,
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
		})
	if err != nil {
		panic(err)
	}
}
