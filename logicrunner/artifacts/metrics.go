// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package artifacts

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
	statCalls   = stats.Int64("artifactmanager_calls", "The number of AM method calls", stats.UnitDimensionless)
	statLatency = stats.Int64("artifactmanager_latency", "The latency in milliseconds per AM call", stats.UnitMilliseconds)

	statRedirects = stats.Int64("artifactmanager_redirects", "The number redirects happens on AM", stats.UnitDimensionless)
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
			Name:        "artifactmanager_latency",
			Description: statLatency.Description(),
			Measure:     statLatency,
			Aggregation: view.Distribution(25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000),
			TagKeys:     commontags,
		},

		&view.View{
			Name:        statRedirects.Name(),
			Description: statRedirects.Description(),
			Measure:     statRedirects,
			Aggregation: view.Count(),
		},
	)
	if err != nil {
		panic(err)
	}
}
