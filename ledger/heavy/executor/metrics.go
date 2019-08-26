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

package executor

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statFinalizedPulse = stats.Int64(
		"heavy_finalized_pulse",
		"last pulse with fully finalized data",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statFinalizedPulse.Name(),
			Description: statFinalizedPulse.Description(),
			Measure:     statFinalizedPulse,
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}
