/*
 *    Copyright 2020 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package insolar

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	postgresConnectionLatency = stats.Float64(
		"postgres_connection_latency",
		"time spent on acquiring connection",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        "postgres_conenction_latency_milliseconds",
			Description: "acquiring connection latency",
			Measure:     postgresConnectionLatency,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000),
		},
	)
	if err != nil {
		panic(err)
	}
}
