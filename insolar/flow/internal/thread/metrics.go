// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package thread

import (
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	tagProcedureName = insmetrics.MustTagKey("proc_type")
	tagResult        = insmetrics.MustTagKey("result")
)

var (
	procCallTime = stats.Float64(
		"flow_procedure_latency",
		"time spent on procedure run",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        procCallTime.Name(),
			Description: procCallTime.Description(),
			Measure:     procCallTime,
			Aggregation: view.Distribution(1, 10, 100, 1000, 5000, 10000),
			TagKeys:     []tag.Key{tagProcedureName, tagResult},
		},
	)
	if err != nil {
		panic(err)
	}
}
