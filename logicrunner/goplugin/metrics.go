package goplugin

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	tagMethodName = insmetrics.MustTagKey("methodName")
)

var (
	statGopluginContractMethodTime = stats.Float64(
		"goplugin/contract/method/time",
		"time spent on execution contract, measured in goplugin",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Measure:     statGopluginContractMethodTime,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000, 10000, 20000),
			TagKeys:     []tag.Key{tagMethodName},
		},
	)
	if err != nil {
		panic(err)
	}
}

