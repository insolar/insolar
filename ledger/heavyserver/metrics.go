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

package heavyserver

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	tagJet = insmetrics.MustTagKey("jet")
	// tagPulse = insmetrics.MustTagKey("result")
)

var (
	statSyncedRecords = stats.Int64("heavyserver/synced/records", "The number synced records", stats.UnitDimensionless)
	statSyncedPulse   = stats.Int64("heavyserver/synced/pulse", "Last synced pulse", stats.UnitDimensionless)
	// statLatency = stats.Int64("artifactmanager/latency", "The latency in milliseconds per AM call", stats.UnitMilliseconds)
)

func init() {
	commontags := []tag.Key{tagJet}
	err := view.Register(
		&view.View{
			Name:        statSyncedRecords.Name(),
			Description: statSyncedRecords.Description(),
			Measure:     statSyncedRecords,
			Aggregation: view.Count(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        statSyncedPulse.Name(),
			Description: statSyncedPulse.Description(),
			Measure:     statSyncedPulse,
			Aggregation: view.LastValue(),
			TagKeys:     commontags,
		},
	)
	if err != nil {
		panic(err)
	}
}
