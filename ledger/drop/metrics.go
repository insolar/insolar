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

package drop

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
	statDropInMemoryAddedCount = stats.Int64(
		"dropstorage/added/count",
		"How many drop-records have been saved in a drop storage",
		stats.UnitDimensionless,
	)
	statDropInMemoryRemovedCount = stats.Int64(
		"dropstorage/removed/count",
		"How many drop-records have been removed from a drop storage",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statDropInMemoryAddedCount.Name(),
			Description: statDropInMemoryAddedCount.Description(),
			Measure:     statDropInMemoryAddedCount,
			Aggregation: view.Count(),
			TagKeys:     []tag.Key{inmemoryStorage},
		},
		&view.View{
			Name:        statDropInMemoryRemovedCount.Name(),
			Description: statDropInMemoryRemovedCount.Description(),
			Measure:     statDropInMemoryRemovedCount,
			Aggregation: view.Count(),
			TagKeys:     []tag.Key{inmemoryStorage},
		},
	)
	if err != nil {
		panic(err)
	}
}
