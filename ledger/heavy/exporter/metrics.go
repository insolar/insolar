// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package exporter

import (
	"context"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"google.golang.org/grpc/metadata"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

const ObsID = "observer_id"

var (
	TagHeavyExporterMethodName = insmetrics.MustTagKey("heavy_exporter_method_name")
	TagHeavyIdObserver         = insmetrics.MustTagKey("heavy_exporter_observer_id")
)

var (
	HeavyExporterMethodTiming = stats.Float64(
		"heavy_exporter_method_timing",
		"time spent in exporter method",
		stats.UnitMilliseconds,
	)

	HeavyExporterLastExportedPulse = stats.Int64(
		"heavy_exporter_last_exported_pulse",
		"the last pulse fetched by an observer",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        HeavyExporterMethodTiming.Name(),
			Description: HeavyExporterMethodTiming.Description(),
			Measure:     HeavyExporterMethodTiming,
			TagKeys:     []tag.Key{TagHeavyExporterMethodName, TagHeavyIdObserver},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000, 10000, 20000),
		},
		&view.View{
			Name:        HeavyExporterLastExportedPulse.Name(),
			Description: HeavyExporterLastExportedPulse.Description(),
			Measure:     HeavyExporterLastExportedPulse,
			TagKeys:     []tag.Key{TagHeavyExporterMethodName, TagHeavyIdObserver},
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}

func addTagsForExporterMethodTiming(authRequired bool, ctx context.Context, methodName string) context.Context {
	observer := "unknown"
	if !authRequired {
		ctx = insmetrics.InsertTag(ctx, TagHeavyIdObserver, observer)
		ctx = insmetrics.InsertTag(ctx, TagHeavyExporterMethodName, methodName)
		return ctx
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if _, isContain := md[ObsID]; isContain && ok {
		observer = md.Get(ObsID)[0]
	}
	ctx = insmetrics.InsertTag(ctx, TagHeavyIdObserver, observer)
	ctx = insmetrics.InsertTag(ctx, TagHeavyExporterMethodName, methodName)
	return ctx
}
