///
// Copyright 2020 Insolar Technologies GmbH
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
///

package main

import (
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statWholePulseWrite = stats.Int64(
		"time_to_store_whole_pulse_data",
		"Time it takes to write whole",
		stats.UnitMilliseconds,
	)
	statIter = stats.Int64(
		"iter",
		"Iteration number",
		stats.UnitDimensionless,
	)

	statVlogSize = stats.Int64(
		"vlog_size",
		"Size of vlog",
		stats.UnitBytes,
	)

	statLSMSize = stats.Int64(
		"lsm_size",
		"Size of lsm",
		stats.UnitBytes,
	)

	statNumRecords = stats.Int64(
		"num_record",
		"Number of records",
		stats.UnitDimensionless)

	statNumIndexes = stats.Int64(
		"num_indexes",
		"Number of indexes",
		stats.UnitDimensionless)

	statNumDrops = stats.Int64(
		"num_drops",
		"Number of drops",
		stats.UnitDimensionless)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statWholePulseWrite.Name(),
			Description: statWholePulseWrite.Description(),
			Measure:     statWholePulseWrite,
			Aggregation: view.Distribution(
				float64(time.Millisecond*5),
				float64(time.Millisecond*10),
				float64(time.Millisecond*20),
				float64(time.Millisecond*50),
				float64(time.Millisecond*100),
				float64(time.Millisecond*200),
				float64(time.Millisecond*800),
				float64(time.Second*2),
			),
		},
		&view.View{
			Name:        statIter.Name(),
			Description: statIter.Description(),
			Measure:     statIter,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statVlogSize.Name(),
			Description: statVlogSize.Description(),
			Measure:     statVlogSize,
			Aggregation: view.LastValue(),
		},
		&view.View{
			Name:        statLSMSize.Name(),
			Description: statLSMSize.Description(),
			Measure:     statLSMSize,
			Aggregation: view.LastValue(),
		},
		&view.View{
			Name:        statNumDrops.Name(),
			Description: statNumDrops.Description(),
			Measure:     statNumDrops,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statNumRecords.Name(),
			Description: statNumRecords.Description(),
			Measure:     statNumRecords,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statNumIndexes.Name(),
			Description: statNumIndexes.Description(),
			Measure:     statNumIndexes,
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		panic(err)
	}
}
