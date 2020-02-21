// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package instrumenter

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	tagError     = insmetrics.MustTagKey("error")
	tagMethod    = insmetrics.MustTagKey("method")
	tagSubMethod = insmetrics.MustTagKey("subMethod")
)

var (
	incomingRequests = stats.Int64("api_incoming", "Count of incoming requests", stats.UnitDimensionless)
	statLatency      = stats.Int64("api_time", "The latency in milliseconds per API call", stats.UnitMilliseconds)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statLatency.Name(),
			Description: statLatency.Description(),
			Measure:     statLatency,
			Aggregation: view.Distribution(25, 500, 1000, 5000, 10000, 15000, 24800),
			TagKeys:     []tag.Key{tagMethod, tagSubMethod, tagError},
		},
		&view.View{
			Name:        incomingRequests.Name(),
			Description: incomingRequests.Description(),
			TagKeys:     []tag.Key{tagMethod},
			Measure:     incomingRequests,
			Aggregation: view.Count(),
		},
	)

	if err != nil {
		panic(err)
	}
}
