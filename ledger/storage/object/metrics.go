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

package object

import (
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	inmemoryStorage = insmetrics.MustTagKey("inmemorystorage")
)

var (
	statIndexInMemoryCount = stats.Int64(
		"indexstorage/inmemory/count",
		"How many index-records have been saved in in-memory index storage",
		stats.UnitDimensionless,
	)
	statRecordInMemoryCount = stats.Int64(
		"recordstorage/inmemory/count",
		"How many records have been saved in in-memory record storage",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statIndexInMemoryCount.Name(),
			Description: statIndexInMemoryCount.Description(),
			Measure:     statIndexInMemoryCount,
			Aggregation: view.Count(),
			TagKeys:     []tag.Key{inmemoryStorage},
		},
		&view.View{
			Name:        statRecordInMemoryCount.Name(),
			Description: statRecordInMemoryCount.Description(),
			Measure:     statRecordInMemoryCount,
			Aggregation: view.Count(),
			TagKeys:     []tag.Key{inmemoryStorage},
		},
	)
	if err != nil {
		panic(err)
	}
}
