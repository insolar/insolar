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
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statBlobInMemorySize = stats.Int64(
		"blobstorage/inmemory/size",
		"Size of the blob-records saved in blobsStor",
		stats.UnitBytes,
	)
	statBlobInMemoryCount = stats.Int64(
		"blobstorage/inmemory/count",
		"How many blob-records saved in blobsStor",
		stats.UnitDimensionless,
	)
	statBlobInMemoryRemovedCount = stats.Int64(
		"blobstorage/inmemory/removed/count",
		"How many blob-records removed from blobsStor",
		stats.UnitDimensionless,
	)
	statBlobInStorageSize = stats.Int64(
		"blobstorage/persistent/size",
		"Size of the blob-records persisted in blob storage",
		stats.UnitBytes,
	)
	statBlobInStorageCount = stats.Int64(
		"blobstorage/persistent/count",
		"How many blob-records persisted in blob storage",
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
		},
		&view.View{
			Name:        statBlobInMemoryCount.Name(),
			Description: statBlobInMemoryCount.Description(),
			Measure:     statBlobInMemoryCount,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        statBlobInMemoryRemovedCount.Name(),
			Description: statBlobInMemoryRemovedCount.Description(),
			Measure:     statBlobInMemoryRemovedCount,
			Aggregation: view.Count(),
		},
		&view.View{
			Name:        statBlobInStorageSize.Name(),
			Description: statBlobInStorageSize.Description(),
			Measure:     statBlobInStorageSize,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statBlobInStorageCount.Name(),
			Description: statBlobInStorageCount.Description(),
			Measure:     statBlobInStorageCount,
			Aggregation: view.Count(),
		},
	)
	if err != nil {
		panic(err)
	}
}
