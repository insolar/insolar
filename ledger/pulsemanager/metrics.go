/*
 *    Copyright 2019 Insolar Technologies
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

package pulsemanager

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statCleanLatencyTotal = stats.Int64("lightcleanup/latency/total", "Light storage cleanup time in milliseconds", stats.UnitMilliseconds)
)

func init() {
	err := view.Register(

		&view.View{
			Name:        statCleanLatencyTotal.Name(),
			Description: statCleanLatencyTotal.Description(),
			Measure:     statCleanLatencyTotal,
			Aggregation: view.Distribution(100, 500, 1000, 5000, 10000),
		},
	)
	if err != nil {
		panic(err)
	}
}
