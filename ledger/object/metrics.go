// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package object

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statIndexesAddedCount = stats.Int64(
		"object_indexes_added_count",
		"How many bucket have been created on a node",
		stats.UnitDimensionless,
	)
	statIndexesRemovedCount = stats.Int64(
		"object_indexes_removed_count",
		"How many bucket have been removed from a node",
		stats.UnitDimensionless,
	)
	statRecordInMemoryAddedCount = stats.Int64(
		"record_storage_added_count",
		"How many records have been saved to a in-memory storage",
		stats.UnitDimensionless,
	)
	statRecordInMemoryRemovedCount = stats.Int64(
		"record_storage_removed_count",
		"How many records have been removed from a in-memory storage",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statIndexesAddedCount.Name(),
			Description: statIndexesAddedCount.Description(),
			Measure:     statIndexesAddedCount,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statIndexesRemovedCount.Name(),
			Description: statIndexesRemovedCount.Description(),
			Measure:     statIndexesRemovedCount,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statRecordInMemoryAddedCount.Name(),
			Description: statRecordInMemoryAddedCount.Description(),
			Measure:     statRecordInMemoryAddedCount,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statRecordInMemoryRemovedCount.Name(),
			Description: statRecordInMemoryRemovedCount.Description(),
			Measure:     statRecordInMemoryRemovedCount,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statRecordInMemoryRemovedCount.Name(),
			Description: statRecordInMemoryRemovedCount.Description(),
			Measure:     statRecordInMemoryRemovedCount,
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		panic(err)
	}
}
