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

package blob

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
	statBlobInMemorySize = stats.Int64(
		"blobstorage/inmemory/size",
		"Size of the blob-records in in-memory blob storage",
		stats.UnitBytes,
	)
	statBlobInMemoryCount = stats.Int64(
		"blobstorage/inmemory/count",
		"How many blob-records have been saved in in-memory blob storage",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statBlobInMemorySize.Name(),
			Description: statBlobInMemorySize.Description(),
			Measure:     statBlobInMemorySize,
			Aggregation: view.Sum(),
			TagKeys:     []tag.Key{inmemoryStorage},
		},
		&view.View{
			Name:        statBlobInMemoryCount.Name(),
			Description: statBlobInMemoryCount.Description(),
			Measure:     statBlobInMemoryCount,
			Aggregation: view.Count(),
			TagKeys:     []tag.Key{inmemoryStorage},
		},
	)
	if err != nil {
		panic(err)
	}
}
