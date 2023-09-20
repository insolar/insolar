package dispatcher

import (
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	tagMessageType = insmetrics.MustTagKey("message_type")
	tagResult      = insmetrics.MustTagKey("result")
)

var (
	statProcessTime = stats.Float64(
		"flow_dispatcher_process_latency",
		"process handlers call latency (handlers latency+overhead)",
		stats.UnitMilliseconds,
	)
	statHandlerTime = stats.Float64(
		"flow_dispatcher_handler_latency",
		"handlers latency",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statProcessTime.Name(),
			Description: statProcessTime.Description(),
			Measure:     statProcessTime,
			Aggregation: view.Distribution(1, 10, 100, 1000, 5000, 10000),
			TagKeys:     []tag.Key{tagMessageType, tagResult},
		},
		&view.View{
			Name:        statHandlerTime.Name(),
			Description: statHandlerTime.Description(),
			Measure:     statHandlerTime,
			Aggregation: view.Distribution(1, 10, 100, 1000, 5000, 10000),
			TagKeys:     []tag.Key{tagMessageType, tagResult},
		},
	)
	if err != nil {
		panic(err)
	}
}
