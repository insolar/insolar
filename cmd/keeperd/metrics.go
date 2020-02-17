// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

func GetRegistry() *prometheus.Registry {
	registry := prometheus.NewRegistry()
	registerer := prometheus.WrapRegistererWith(prometheus.Labels{"role": "keeper"}, registry)

	registerer.MustRegister(prometheus.NewProcessCollector(
		prometheus.ProcessCollectorOpts{Namespace: "insolar"},
	))
	registerer.MustRegister(prometheus.NewGoCollector())

	return registry
}

var (
	IsAvailable = stats.Int64(
		"is_available",
		"1 if all metrics are OK and platform is available for requests, 0 otherwise",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        IsAvailable.Name(),
			Description: IsAvailable.Description(),
			Measure:     IsAvailable,
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}
