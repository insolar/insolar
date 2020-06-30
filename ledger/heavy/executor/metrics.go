// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"time"

	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var TagJetID = insmetrics.MustTagKey("jet_id_heavy")

var (
	statJets = stats.Int64(
		"heavy_jets_counter",
		"how many jets on start of pulse",
		stats.UnitDimensionless,
	)

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

	statJetDropStoreTime = stats.Float64(
		"jet_drop_store_time",
		"duration backup process",
		stats.UnitMilliseconds,
	)

	statBadgerValueGCTime = stats.Int64(
		"badger_value_gc_time",
		"duration of badger's value GC",
		"s",
	)

	topSyncPulseTime = stats.Float64(
		"jetkeeper_topsync_time",
		"time spent on topSyncPulse",
		stats.UnitMilliseconds,
	)

	getTime = stats.Float64(
		"jetkeeper_get_time",
		"time spent on get",
		stats.UnitMilliseconds,
	)

	setTime = stats.Float64(
		"jetkeeper_set_time",
		"time spent on set",
		stats.UnitMilliseconds,
	)
	setRetries = stats.Int64(
		"jetkeeper_set_retries",
		"retries while jetkeeper set",
		stats.UnitDimensionless,
	)

	updateSyncPulseTime = stats.Float64(
		"jetkeeper_updatesyncpulse_time",
		"time spent on updateSyncPulse",
		stats.UnitMilliseconds,
	)
	updateSyncPulseRetries = stats.Int64(
		"jetkeeper_updatesyncpuls_retries",
		"retries while jetkeeper updatesyncpuls",
		stats.UnitDimensionless,
	)

	TruncateHeadTime = stats.Float64(
		"jetkeeper_truncate_head_time",
		"time spent on TruncateHead",
		stats.UnitMilliseconds,
	)
	TruncateHeadRetries = stats.Int64(
		"jetkeeper_truncate_head_retries",
		"retries while jetkeeper TruncateHead",
		stats.UnitDimensionless,
	)

	statRecordInDrop = stats.Int64(
		"record_in_drop_count",
		"How many record in drop were received from a light-node",
		stats.UnitDimensionless,
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
			Name:        statJets.Name(),
			Description: statJets.Description(),
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
			Name:        statJetDropStoreTime.Name(),
			Description: statJetDropStoreTime.Description(),
			Measure:     statJetDropStoreTime,
			Aggregation: view.Distribution(0.0, float64(time.Minute)),
		},
		&view.View{
			Name:        statBadgerValueGCTime.Name(),
			Description: statBadgerValueGCTime.Description(),
			Measure:     statBadgerValueGCTime,
			Aggregation: view.Distribution(0.0, float64(time.Minute)*2),
		},
		&view.View{
			Name:        topSyncPulseTime.Name(),
			Description: topSyncPulseTime.Description(),
			Measure:     topSyncPulseTime,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        getTime.Name(),
			Description: getTime.Description(),
			Measure:     getTime,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        setTime.Name(),
			Description: setTime.Description(),
			Measure:     setTime,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        setRetries.Name(),
			Description: setRetries.Description(),
			Measure:     setRetries,
			Aggregation: view.Distribution(0, 1, 2, 3, 4, 5, 10),
		},
		&view.View{
			Name:        updateSyncPulseTime.Name(),
			Description: updateSyncPulseTime.Description(),
			Measure:     updateSyncPulseTime,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        updateSyncPulseRetries.Name(),
			Description: updateSyncPulseRetries.Description(),
			Measure:     updateSyncPulseRetries,
			Aggregation: view.Distribution(0, 1, 2, 3, 4, 5, 10),
		},
		&view.View{
			Name:        TruncateHeadTime.Name(),
			Description: TruncateHeadTime.Description(),
			Measure:     TruncateHeadTime,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        TruncateHeadRetries.Name(),
			Description: TruncateHeadRetries.Description(),
			Measure:     TruncateHeadRetries,
			Aggregation: view.Distribution(0, 1, 2, 3, 4, 5, 10),
		},
		&view.View{
			Name:        statRecordInDrop.Name(),
			Description: statRecordInDrop.Description(),
			Measure:     statRecordInDrop,
			TagKeys:     []tag.Key{TagJetID},
			Aggregation: view.Count(),
		},
	)
	if err != nil {
		panic(err)
	}
}
