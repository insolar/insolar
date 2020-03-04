// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulse

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	ForPulseNumberTime = stats.Float64(
		"pulse_for_pulse_number_time",
		"time spent on ForPulseNumber",
		stats.UnitMilliseconds,
	)
	LatestTime = stats.Float64(
		"pulse_latest_time",
		"time spent on Latest",
		stats.UnitMilliseconds,
	)
	TruncateHeadRetries = stats.Int64(
		"pulse_truncate_head_retries",
		"retries while TruncateHead",
		stats.UnitDimensionless,
	)
	TruncateHeadTime = stats.Float64(
		"pulse_truncate_head_time",
		"time spent on TruncateHead",
		stats.UnitMilliseconds,
	)
	AppendRetries = stats.Int64(
		"pulse_append_retries",
		"retries while Append",
		stats.UnitDimensionless,
	)
	AppendTime = stats.Float64(
		"pulse_append_time",
		"time spent on Append",
		stats.UnitMilliseconds,
	)
	ForwardsTime = stats.Float64(
		"pulse_forwards_time",
		"time spent on Forwards",
		stats.UnitMilliseconds,
	)
	BackwardsTime = stats.Float64(
		"pulse_backwards_time",
		"time spent on Backwards",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        ForPulseNumberTime.Name(),
			Description: ForPulseNumberTime.Description(),
			Measure:     ForPulseNumberTime,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        LatestTime.Name(),
			Description: LatestTime.Description(),
			Measure:     LatestTime,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
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
			Name:        AppendTime.Name(),
			Description: AppendTime.Description(),
			Measure:     AppendTime,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        AppendRetries.Name(),
			Description: AppendRetries.Description(),
			Measure:     AppendRetries,
			Aggregation: view.Distribution(0, 1, 2, 3, 4, 5, 10),
		},
		&view.View{
			Name:        ForwardsTime.Name(),
			Description: ForwardsTime.Description(),
			Measure:     ForwardsTime,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        BackwardsTime.Name(),
			Description: BackwardsTime.Description(),
			Measure:     BackwardsTime,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
	)
	if err != nil {
		panic(err)
	}
}
