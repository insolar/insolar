// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package exporter

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	TagHeavyExporterMethodName = insmetrics.MustTagKey("heavy_exporter_method_name")
)

var (
	HeavyExporterMethodTiming = stats.Float64(
		"heavy_exporter_method_timing",
		"time spent in exporter method",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        HeavyExporterMethodTiming.Name(),
			Description: HeavyExporterMethodTiming.Description(),
			Measure:     HeavyExporterMethodTiming,
			TagKeys:     []tag.Key{TagHeavyExporterMethodName},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000, 10000, 20000),
		},
	)
	if err != nil {
		panic(err)
	}
}
