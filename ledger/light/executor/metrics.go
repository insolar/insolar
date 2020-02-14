// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	// StatRequestsOpened specifies a metric about opened requests for the current light
	StatRequestsOpened = stats.Int64(
		"requests_opened",
		"How many requests are opened",
		stats.UnitDimensionless,
	)
	// StatRequestsClosed specifies a metric about closed requests for the current light
	StatRequestsClosed = stats.Int64(
		"requests_closed",
		"How many requests are closed",
		stats.UnitDimensionless,
	)

	statHotObjectsTotal = stats.Int64("hotdata_objects_total", "Amount of hot records prepared to send for next executors", stats.UnitDimensionless)
	statHotObjectsSend  = stats.Int64("hotdata_objects_send", "Amount of hot records actually sent to next executors", stats.UnitDimensionless)

	statJets      = stats.Int64("jets", "jets counter", stats.UnitDimensionless)
	statJetSplits = stats.Int64("jet_splits", "jet splits counter", stats.UnitDimensionless)

	statDrop        = stats.Int64("drops", "How many drop records have created", stats.UnitDimensionless)
	statDropRecords = stats.Int64("drop_records", "Amount of records in drop", stats.UnitDimensionless)

	statLastReplicatedPulse = stats.Int64(
		"light_last_sent_pulse",
		"last pulse sent to heavy",
		stats.UnitDimensionless,
	)

	statFilamentLength = stats.Int64(
		"filament_length",
		"How many records are in filaments during iterations",
		stats.UnitDimensionless,
	)
	statFilamentFetchedCount = stats.Int64(
		"filament_fetched_count",
		"How many records are in fetched from network filament segment",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        StatRequestsOpened.Name(),
			Description: StatRequestsOpened.Description(),
			Measure:     StatRequestsOpened,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        StatRequestsClosed.Name(),
			Description: StatRequestsClosed.Description(),
			Measure:     StatRequestsClosed,
			Aggregation: view.Sum(),
		},

		&view.View{
			Name:        statHotObjectsTotal.Name(),
			Description: statHotObjectsTotal.Description(),
			Measure:     statHotObjectsTotal,
			Aggregation: view.Sum(),
		},

		&view.View{
			Name:        statHotObjectsSend.Name(),
			Description: statHotObjectsSend.Description(),
			Measure:     statHotObjectsSend,
			Aggregation: view.Sum(),
		},

		&view.View{
			Name:        "jets_counter",
			Description: "how many jets on start of pulse",
			Measure:     statJets,
			Aggregation: view.LastValue(),
		},
		&view.View{
			Name:        "jet_splits_total",
			Description: "how many jet splits performed",
			Measure:     statJetSplits,
			Aggregation: view.Sum(),
		},

		&view.View{
			Name:        "drops_total",
			Description: statDrop.Description(),
			Measure:     statDrop,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        "drop_records_total",
			Description: statDropRecords.Description(),
			Measure:     statDropRecords,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statLastReplicatedPulse.Name(),
			Description: statLastReplicatedPulse.Description(),
			Measure:     statLastReplicatedPulse,
			Aggregation: view.LastValue(),
		},

		&view.View{
			Name:        statFilamentLength.Name(),
			Description: statFilamentLength.Description(),
			Measure:     statFilamentLength,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statFilamentFetchedCount.Name(),
			Description: statFilamentFetchedCount.Description(),
			Measure:     statFilamentFetchedCount,
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		panic(err)
	}
}
