package blockexplorer

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	tagMethod = insmetrics.MustTagKey("method")
	tagResult = insmetrics.MustTagKey("result")
)

var (
	statCalls   = stats.Int64("blockexplorermanager/calls", "The number of BE method calls", stats.UnitDimensionless)
	statLatency = stats.Int64("blockexplorermanager/latency", "The latency in milliseconds per BE call", stats.UnitMilliseconds)
)

func init() {
	commontags := []tag.Key{tagMethod, tagResult}
	err := view.Register(
		&view.View{
			Name:        statCalls.Name(),
			Description: statCalls.Description(),
			Measure:     statCalls,
			Aggregation: view.Count(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        "blockexplorermanager_latency",
			Description: statLatency.Description(),
			Measure:     statLatency,
			Aggregation: view.Distribution(0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000),
			TagKeys:     commontags,
		},
	)
	if err != nil {
		panic(err)
	}
}
