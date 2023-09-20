package insolar

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	postgresConnectionLatency = stats.Float64(
		"postgres_connection_latency",
		"time spent on acquiring connection",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        "postgres_conenction_latency_milliseconds",
			Description: "acquiring connection latency",
			Measure:     postgresConnectionLatency,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
	)
	if err != nil {
		panic(err)
	}
}
