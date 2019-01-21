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

package heavyclient

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	tagJet = insmetrics.MustTagKey("jet")
)

var (
	statUnsyncedPulsesCount = stats.Int64("heavyclient/unsynced/count", "How many pulses unsynced", stats.UnitDimensionless)
	statFirstUnsyncedPulse  = stats.Int64("heavyclient/unsynced/firstpulse", "First unsynced pulse number", stats.UnitDimensionless)

	statSyncedPulsesCount = stats.Int64("heavyclient/synced/count", "How many pulses unsynced", stats.UnitDimensionless)
)

func init() {
	commontags := []tag.Key{tagJet}
	err := view.Register(
		&view.View{
			Name:        statUnsyncedPulsesCount.Name(),
			Description: statUnsyncedPulsesCount.Description(),
			Measure:     statUnsyncedPulsesCount,
			Aggregation: view.LastValue(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        statFirstUnsyncedPulse.Name(),
			Description: statFirstUnsyncedPulse.Description(),
			Measure:     statFirstUnsyncedPulse,
			Aggregation: view.LastValue(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        statSyncedPulsesCount.Name(),
			Description: statSyncedPulsesCount.Description(),
			Measure:     statSyncedPulsesCount,
			Aggregation: view.Count(),
			TagKeys:     commontags,
		},
	)
	if err != nil {
		panic(err)
	}
}
