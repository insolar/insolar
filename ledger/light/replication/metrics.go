/*
 *    Copyright 2019 Insolar Technologies
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

package replication

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statHeavyPayloadCount = stats.Int64(
		"lightsyncer/heavypayload/count",
		"How many heavy-payload messages were sent to a heavy node",
		stats.UnitDimensionless,
	)
	statErrHeavyPayloadCount = stats.Int64(
		"lightsyncer/failedheavypayload/count",
		"How many heavy-payload messages were failed",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statHeavyPayloadCount.Name(),
			Description: statHeavyPayloadCount.Description(),
			Measure:     statHeavyPayloadCount,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        statErrHeavyPayloadCount.Name(),
			Description: statErrHeavyPayloadCount.Description(),
			Measure:     statErrHeavyPayloadCount,
			Aggregation: view.Count(),
		},
	)
	if err != nil {
		panic(err)
	}
}
