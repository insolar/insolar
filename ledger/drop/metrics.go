// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package drop

import (
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	inmemoryStorage = insmetrics.MustTagKey("inmemorystorage")
)

var (
	statDropInMemoryAddedCount = stats.Int64(
		"dropstorage_added_count",
		"How many drop-records have been saved in a drop storage",
		stats.UnitDimensionless,
	)
	statDropInMemoryRemovedCount = stats.Int64(
		"dropstorage_removed_count",
		"How many drop-records have been removed from a drop storage",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statDropInMemoryAddedCount.Name(),
			Description: statDropInMemoryAddedCount.Description(),
			Measure:     statDropInMemoryAddedCount,
			Aggregation: view.Count(),
			TagKeys:     []tag.Key{inmemoryStorage},
		},
		&view.View{
			Name:        statDropInMemoryRemovedCount.Name(),
			Description: statDropInMemoryRemovedCount.Description(),
			Measure:     statDropInMemoryRemovedCount,
			Aggregation: view.Count(),
			TagKeys:     []tag.Key{inmemoryStorage},
		},
	)
	if err != nil {
		panic(err)
	}
}
