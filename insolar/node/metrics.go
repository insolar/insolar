// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package node

import (
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	NodesPostgresDB = insmetrics.MustTagKey("node_postgres_db")
)

var (
	SetRetries = stats.Int64(
		"node_set_retries",
		"retries while truncating head",
		stats.UnitDimensionless,
	)
	SetTime = stats.Float64(
		"node_set_time",
		"time spent on truncating head",
		stats.UnitMilliseconds,
	)
	AllTime = stats.Float64(
		"node_all_time",
		"time spent on truncating head",
		stats.UnitMilliseconds,
	)
	InRoleTime = stats.Float64(
		"node_inrole_time",
		"time spent on truncating head",
		stats.UnitMilliseconds,
	)
	TruncateRetries = stats.Int64(
		"node_truncate_retries",
		"retries while truncating head",
		stats.UnitDimensionless,
	)
	TruncateTime = stats.Float64(
		"node_truncate_time",
		"time spent on truncating head",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        SetTime.Name(),
			Description: SetTime.Description(),
			Measure:     SetTime,
			TagKeys:     []tag.Key{NodesPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        SetRetries.Name(),
			Description: SetRetries.Description(),
			Measure:     SetRetries,
			TagKeys:     []tag.Key{NodesPostgresDB},
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        AllTime.Name(),
			Description: AllTime.Description(),
			Measure:     AllTime,
			TagKeys:     []tag.Key{NodesPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        InRoleTime.Name(),
			Description: InRoleTime.Description(),
			Measure:     InRoleTime,
			TagKeys:     []tag.Key{NodesPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        TruncateTime.Name(),
			Description: TruncateTime.Description(),
			Measure:     TruncateTime,
			TagKeys:     []tag.Key{NodesPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        TruncateRetries.Name(),
			Description: TruncateRetries.Description(),
			Measure:     TruncateRetries,
			TagKeys:     []tag.Key{NodesPostgresDB},
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		panic(err)
	}
}
