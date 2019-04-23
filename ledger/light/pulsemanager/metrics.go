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

package pulsemanager

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statCleanLatencyTotal = stats.Int64("lightcleanup/latency/total", "Light storage cleanup time in milliseconds", stats.UnitMilliseconds)
	statHotObjectsSent    = stats.Int64("hotdata/objects/total", "Amount of hot objects sent to the next executor", stats.UnitDimensionless)
	statPendingSent       = stats.Int64("hotdata/pending/total", "Amount of pending requests sent to the next executor", stats.UnitDimensionless)
)

func init() {
	err := view.Register(

		&view.View{
			Name:        statCleanLatencyTotal.Name(),
			Description: statCleanLatencyTotal.Description(),
			Measure:     statCleanLatencyTotal,
			Aggregation: view.Distribution(100, 500, 1000, 5000, 10000),
		},

		&view.View{
			Name:        statHotObjectsSent.Name(),
			Description: statHotObjectsSent.Description(),
			Measure:     statHotObjectsSent,
			Aggregation: view.Sum(),
		},

		&view.View{
			Name:        statPendingSent.Name(),
			Description: statPendingSent.Description(),
			Measure:     statPendingSent,
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		panic(err)
	}
}
