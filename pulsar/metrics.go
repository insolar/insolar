// Copyright 2020 Insolar Network Ltd.
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

package pulsar

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statPulseGenerated = stats.Int64("pulsar_pulse_generated", "count of generated pulses", stats.UnitDimensionless)
	statCurrentPulse   = stats.Int64("pulsar_current_pulse", "last generated pulse", stats.UnitDimensionless)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statPulseGenerated.Name(),
			Description: statPulseGenerated.Description(),
			Measure:     statPulseGenerated,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statCurrentPulse.Name(),
			Description: statCurrentPulse.Description(),
			Measure:     statCurrentPulse,
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}
