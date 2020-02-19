// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statJets           = stats.Int64("heavy_jets", "jets counter", stats.UnitDimensionless)
	statFinalizedPulse = stats.Int64(
		"heavy_finalized_pulse",
		"last pulse with fully finalized data",
		stats.UnitDimensionless,
	)

	statAbandonedRequests = stats.Int64(
		"requests_abandoned",
		"Amount of abandoned requests on heavy",
		stats.UnitDimensionless,
	)

	statBackupTime = stats.Int64(
		"backup_time",
		"duration backup process",
		"ns",
	)

	statBadgerValueGCTime = stats.Int64(
		"badger_value_gc_time",
		"duration of badger's value GC",
		"s",
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statFinalizedPulse.Name(),
			Description: statFinalizedPulse.Description(),
			Measure:     statFinalizedPulse,
			Aggregation: view.LastValue(),
		},
		&view.View{
			Name:        "heavy_jets_counter",
			Description: "how many jets on start of pulse",
			Measure:     statJets,
			Aggregation: view.LastValue(),
		},
		&view.View{
			Name:        statAbandonedRequests.Name(),
			Description: statAbandonedRequests.Description(),
			Measure:     statAbandonedRequests,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        statBackupTime.Name(),
			Description: statBackupTime.Description(),
			Measure:     statBackupTime,
			Aggregation: view.Distribution(0.0, float64(time.Minute)),
		},
		&view.View{
			Name:        statBadgerValueGCTime.Name(),
			Description: statBadgerValueGCTime.Description(),
			Measure:     statBadgerValueGCTime,
			Aggregation: view.Distribution(0.0, float64(time.Minute)*2),
		},
	)
	if err != nil {
		panic(err)
	}
}
