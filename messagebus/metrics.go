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

package messagebus

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	tagMessageType = insmetrics.MustTagKey("messageType")
)

var (
	statParcelsSentTotal = stats.Int64(
		"messagebus/parcels/sent/count",
		"number of parcels sent",
		stats.UnitDimensionless,
	)
	statLocallyDeliveredParcelsTotal = stats.Int64(
		"messagebus/parcels/locally/delivered/count",
		"total number of parcels delivered to the same machine",
		stats.UnitDimensionless,
	)
	statParcelsTime = stats.Float64(
		"messagebus/parcels/time",
		"time spent on sending parcels",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Measure:     statParcelsSentTotal,
			Aggregation: view.Sum(),
			TagKeys:     []tag.Key{tagMessageType},
		},
		&view.View{
			Measure:     statLocallyDeliveredParcelsTotal,
			Aggregation: view.Count(),
			TagKeys:     []tag.Key{tagMessageType},
		},
		&view.View{
			Measure:     statParcelsTime,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000, 10000, 20000),
			TagKeys:     []tag.Key{tagMessageType},
		},
	)
	if err != nil {
		panic(err)
	}
}
