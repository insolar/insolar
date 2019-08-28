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
	statHotObjectsTotal = stats.Int64("hotdata/objects/total", "Amount of hot records for next executors", stats.UnitDimensionless)
	statHotObjectsSend  = stats.Int64("hotdata/objects/send", "Amount of hot records actually sent to next executors", stats.UnitDimensionless)

	statJets      = stats.Int64("jets", "jets counter", stats.UnitDimensionless)
	statJetSplits = stats.Int64("jet/splits", "jet splits counter", stats.UnitDimensionless)

	statDrop        = stats.Int64("drops", "How many drop records have created", stats.UnitDimensionless)
	statDropRecords = stats.Int64("drop/records", "Amount of records in drop", stats.UnitDimensionless)

	statLastReplicatedPulse = stats.Int64(
		"light_last_sent_pulse",
		"last pulse sent to heavy",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statHotObjectsTotal.Name(),
			Description: statHotObjectsTotal.Description(),
			Measure:     statHotObjectsTotal,
			Aggregation: view.Sum(),
		},

		&view.View{
			Name:        statHotObjectsSend.Name(),
			Description: statHotObjectsSend.Description(),
			Measure:     statHotObjectsSend,
			Aggregation: view.Sum(),
		},

		&view.View{
			Name:        "jets_counter",
			Description: "how many jets on start of pulse",
			Measure:     statJets,
			Aggregation: view.LastValue(),
		},
		&view.View{
			Name:        "jet_splits_total",
			Description: "how many jet splits performed",
			Measure:     statJetSplits,
			Aggregation: view.Sum(),
		},

		&view.View{
			Name:        "drops_total",
			Description: statDrop.Description(),
			Measure:     statDrop,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        "drop_records_total",
			Description: statDropRecords.Description(),
			Measure:     statDropRecords,
			Aggregation: view.Sum(),
		},

		&view.View{
			Name:        statLastReplicatedPulse.Name(),
			Description: statLastReplicatedPulse.Description(),
			Measure:     statLastReplicatedPulse,
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}
