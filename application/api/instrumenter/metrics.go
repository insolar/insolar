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

package instrumenter

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	tagError     = insmetrics.MustTagKey("error")
	tagMethod    = insmetrics.MustTagKey("method")
	tagSubMethod = insmetrics.MustTagKey("subMethod")
)

var (
	incomingRequests = stats.Int64("api_incoming", "Count of incoming requests", stats.UnitDimensionless)
	statLatency      = stats.Int64("api_time", "The latency in milliseconds per API call", stats.UnitMilliseconds)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statLatency.Name(),
			Description: statLatency.Description(),
			Measure:     statLatency,
			Aggregation: view.Distribution(25, 500, 1000, 5000, 10000, 15000, 24800),
			TagKeys:     []tag.Key{tagMethod, tagSubMethod, tagError},
		},
		&view.View{
			Name:        incomingRequests.Name(),
			Description: incomingRequests.Description(),
			TagKeys:     []tag.Key{tagMethod},
			Measure:     incomingRequests,
			Aggregation: view.Count(),
		},
	)

	if err != nil {
		panic(err)
	}
}
