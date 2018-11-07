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

package artifactmanager

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	tagMethod = insmetrics.MustTagKey("method")
	tagResult = insmetrics.MustTagKey("result")
)

var (
	statCalls   = stats.Int64("artifactmanager/calls", "The number of AM method calls", stats.UnitDimensionless)
	statLatency = stats.Int64("artifactmanager/latency", "The latency in milliseconds per AM call", stats.UnitMilliseconds)
)

func init() {
	commontags := []tag.Key{tagMethod, tagResult}
	err := view.Register(
		&view.View{
			Name:        statCalls.Name(),
			Description: statCalls.Description(),
			Measure:     statCalls,
			Aggregation: view.Count(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        "artifactmanager_latency",
			Description: statLatency.Description(),
			Measure:     statLatency,
			Aggregation: view.Distribution(0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000),
			TagKeys:     commontags,
		},
	)
	if err != nil {
		panic(err)
	}
}
