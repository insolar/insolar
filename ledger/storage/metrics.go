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

package storage

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	recordType = insmetrics.MustTagKey("rectype")
)

var (
	statCleanScanned = stats.Int64("lightcleanup/scanned", "How many records have been scanned on LM cleanup", stats.UnitDimensionless)
	statCleanRemoved = stats.Int64("lightcleanup/removed", "How many records have been removed on LM cleanup", stats.UnitDimensionless)
	statCleanFailed  = stats.Int64("lightcleanup/rmfailed", "How many records have not been removed because of error", stats.UnitDimensionless)

	statPulseDeleted = stats.Int64("lightcleanup/pulses/removed/total", "How many pulses deleted from pulseTracker on LM cleanup", stats.UnitDimensionless)
	statPulseAdded   = stats.Int64("lightcleanup/pulses/added/total", "How many pulses added to pulseTracker", stats.UnitDimensionless)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statCleanScanned.Name(),
			Description: statCleanScanned.Description(),
			Measure:     statCleanScanned,
			Aggregation: view.Sum(),
			TagKeys:     []tag.Key{recordType},
		},
		&view.View{
			Name:        statCleanRemoved.Name(),
			Description: statCleanRemoved.Description(),
			Measure:     statCleanRemoved,
			Aggregation: view.Sum(),
			TagKeys:     []tag.Key{recordType},
		},
		&view.View{
			Name:        statCleanFailed.Name(),
			Description: statCleanFailed.Description(),
			Measure:     statCleanFailed,
			Aggregation: view.Sum(),
			TagKeys:     []tag.Key{recordType},
		},
		&view.View{
			Name:        statPulseDeleted.Name(),
			Description: statPulseDeleted.Description(),
			Measure:     statPulseDeleted,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        statPulseAdded.Name(),
			Description: statPulseAdded.Description(),
			Measure:     statPulseAdded,
			Aggregation: view.Count(),
		},
	)
	if err != nil {
		panic(err)
	}
}
