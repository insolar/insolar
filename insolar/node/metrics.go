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
		"retries on Set",
		stats.UnitDimensionless,
	)
	SetTime = stats.Float64(
		"node_set_time",
		"time spent on Set",
		stats.UnitMilliseconds,
	)
	AllTime = stats.Float64(
		"node_all_time",
		"time spent on All",
		stats.UnitMilliseconds,
	)
	InRoleTime = stats.Float64(
		"node_inrole_time",
		"time spent on InRole",
		stats.UnitMilliseconds,
	)
	TruncateHeadRetries = stats.Int64(
		"node_truncate_retries",
		"retries on TruncateHead",
		stats.UnitDimensionless,
	)
	TruncateHeadTime = stats.Float64(
		"node_truncate_time",
		"time spent on TruncateHead",
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
			Aggregation: view.Distribution(0, 1, 2, 3, 4, 5, 10),
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
			Name:        TruncateHeadTime.Name(),
			Description: TruncateHeadTime.Description(),
			Measure:     TruncateHeadTime,
			TagKeys:     []tag.Key{NodesPostgresDB},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
		&view.View{
			Name:        TruncateHeadRetries.Name(),
			Description: TruncateHeadRetries.Description(),
			Measure:     TruncateHeadRetries,
			TagKeys:     []tag.Key{NodesPostgresDB},
			Aggregation: view.Distribution(0, 1, 2, 3, 4, 5, 10),
		},
	)
	if err != nil {
		panic(err)
	}
}
