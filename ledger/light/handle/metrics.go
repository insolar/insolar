package handle

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	KeyMsgType   = insmetrics.MustTagKey("msg_type")
	KeyErrorCode = insmetrics.MustTagKey("error_code")
)

var (
	statHandlerError = stats.Int64(
		"handler_errors",
		"How many procedures return errors",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statHandlerError.Name(),
			Description: statHandlerError.Description(),
			Measure:     statHandlerError,
			TagKeys:     []tag.Key{KeyMsgType, KeyErrorCode},
			Aggregation: view.Count(),
		},
	)
	if err != nil {
		panic(err)
	}
}
