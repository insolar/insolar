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
	ForPulseTime = stats.Float64(
		"drop_for_pulse_time",
		"time spent on ForPulse",
		stats.UnitMilliseconds,
	)
	SetTime = stats.Float64(
		"drop_set_time",
		"time spent on Set",
		stats.UnitMilliseconds,
	)
	SetRetries = stats.Int64(
		"drop_set_retries",
		"retries while Set",
		stats.UnitDimensionless,
	)
	TruncateHeadTime = stats.Float64(
		"drop_truncate_head_time",
		"time spent on TruncateHead",
		stats.UnitMilliseconds,
	)
	TruncateHeadRetries = stats.Int64(
		"drop_truncate_head_retries",
		"retries while TruncateHead",
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
		&view.View{
			Name:        ForPulseTime.Name(),
			Description: ForPulseTime.Description(),
			Measure:     ForPulseTime,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        SetTime.Name(),
			Description: SetTime.Description(),
			Measure:     SetTime,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        SetRetries.Name(),
			Description: SetRetries.Description(),
			Measure:     SetRetries,
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
	)
	if err != nil {
		panic(err)
	}
}
