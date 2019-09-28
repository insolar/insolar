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

package proc

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statHotsAbandoned = stats.Int64(
		"requests_abandoned",
		"How many abandoned requests in hot data",
		stats.UnitDimensionless,
	)
	statAbandonedRequestAge = stats.Int64(
		"oldest_abandoned_request_age",
		"How many pulses passed from last abandoned request creation",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statHotsAbandoned.Name(),
			Description: statHotsAbandoned.Description(),
			Measure:     statHotsAbandoned,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        statAbandonedRequestAge.Name(),
			Description: statAbandonedRequestAge.Description(),
			Measure:     statAbandonedRequestAge,
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}
