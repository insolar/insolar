// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulse

import (
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	PulsePostgresDB = insmetrics.MustTagKey("pulse_postgres_db")
)

var (
	ForPulseNumberTime = stats.Float64(
		"pulse_for_pulse_time",
		"time spent on ForPulseNumber",
		stats.UnitMilliseconds,
	)
	LatestTime = stats.Float64(
		"pulse_latest_time",
		"time spent on Latest",
		stats.UnitMilliseconds,
	)
	TruncateHeadRetries = stats.Int64(
		"pulse_truncate_retries",
		"retries while truncating head",
		stats.UnitDimensionless,
	)
	TruncateHeadTime = stats.Float64(
		"pulse_truncate_time",
		"time spent on truncating head",
		stats.UnitMilliseconds,
	)
	AppendRetries = stats.Int64(
		"pulse_append_retries",
		"retries while appending pulse",
		stats.UnitDimensionless,
	)
	AppendTime = stats.Float64(
		"pulse_append_time",
		"time spent on appending pulse",
		stats.UnitMilliseconds,
	)
	ForwardsTime = stats.Float64(
		"pulse_forwards_time",
		"time spent on forwards",
		stats.UnitMilliseconds,
	)
	BackwardsTime = stats.Float64(
		"pulse_backwards_time",
		"time spent on backwards",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        ForPulseNumberTime.Name(),
			Description: ForPulseNumberTime.Description(),
			Measure:     ForPulseNumberTime,
			TagKeys:     []tag.Key{PulsePostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        LatestTime.Name(),
			Description: LatestTime.Description(),
			Measure:     LatestTime,
			TagKeys:     []tag.Key{PulsePostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        TruncateHeadTime.Name(),
			Description: TruncateHeadTime.Description(),
			Measure:     TruncateHeadTime,
			TagKeys:     []tag.Key{PulsePostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        TruncateHeadRetries.Name(),
			Description: TruncateHeadRetries.Description(),
			Measure:     TruncateHeadRetries,
			TagKeys:     []tag.Key{PulsePostgresDB},
			Aggregation: view.Distribution(0, 1, 2, 3, 4, 5, 10),
		},
		&view.View{
			Name:        AppendTime.Name(),
			Description: AppendTime.Description(),
			Measure:     AppendTime,
			TagKeys:     []tag.Key{PulsePostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        AppendRetries.Name(),
			Description: AppendRetries.Description(),
			Measure:     AppendRetries,
			TagKeys:     []tag.Key{PulsePostgresDB},
			Aggregation: view.Distribution(0, 1, 2, 3, 4, 5, 10),
		},
		&view.View{
			Name:        ForwardsTime.Name(),
			Description: ForwardsTime.Description(),
			Measure:     ForwardsTime,
			TagKeys:     []tag.Key{PulsePostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        BackwardsTime.Name(),
			Description: BackwardsTime.Description(),
			Measure:     BackwardsTime,
			TagKeys:     []tag.Key{PulsePostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
	)
	if err != nil {
		panic(err)
	}
}
