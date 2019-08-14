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

package bus

import (
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	tagMessageType = insmetrics.MustTagKey("message_type")
)

var (
	statSent = stats.Int64(
		"bus/sent",
		"messages stats",
		stats.UnitDimensionless,
	)
	statSentTime = stats.Float64(
		"bus/sent/time",
		"time spent on sending parcels",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(

		&view.View{
			Name:        statSent.Name() + "/count",
			Description: statSent.Description(),
			Measure:     statSent,
			Aggregation: view.Count(),
			TagKeys:     []tag.Key{tagMessageType},
		},

		&view.View{
			Name:        statSent.Name() + "/bytes",
			Description: statSent.Description(),
			Measure:     statSent,
			Aggregation: view.Sum(),
			TagKeys:     []tag.Key{tagMessageType},
		},

		&view.View{
			Measure:     statSentTime,
			Aggregation: view.Distribution(1, 10, 100, 1000, 5000, 10000, 20000),
			TagKeys:     []tag.Key{tagMessageType},
		},
	)
	if err != nil {
		panic(err)
	}
}
