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

const IDobs = "idobserver"

var (
	TagHeavyExporterMethodName = insmetrics.MustTagKey("heavy_exporter_method_name")
	// public - data from observer on public side(crypto exchange). internal - from internal network
	TagHeavyIdObserver = insmetrics.MustTagKey("heavy_exporter_type_observer")
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
			TagKeys:     []tag.Key{TagHeavyExporterMethodName, TagHeavyIdObserver},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000, 10000, 20000),
		},
	)
	if err != nil {
		panic(err)
	}
}

func addTagsForExporterMethodTiming(ctx context.Context, methodName string) context.Context {
	typeObserver := "internal"
	md, ok := metadata.FromIncomingContext(ctx)
	if _, isContain := md[IDobs]; isContain && ok {
		typeObserver = md.Get(IDobs)[0]
	}
	ctx = insmetrics.InsertTag(ctx, TagHeavyIdObserver, typeObserver)
	ctx = insmetrics.InsertTag(ctx, TagHeavyExporterMethodName, methodName)
	return ctx
}
