//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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
