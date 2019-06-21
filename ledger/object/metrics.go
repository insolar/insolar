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
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statBucketAddedCount = stats.Int64(
		"object/bucket/added/count",
		"How many bucket have been created on a node",
		stats.UnitDimensionless,
	)
	statBucketRemovedCount = stats.Int64(
		"object/bucket/removed/count",
		"How many bucket have been removed from a node",
		stats.UnitDimensionless,
	)
	statObjectPendingRequestsInMemoryAddedCount = stats.Int64(
		"оbject/pendings/added/count",
		"How many index-records have been saved in in-indexStorage index storage",
		stats.UnitDimensionless,
	)
	statObjectPendingRequestsInMemoryRemovedCount = stats.Int64(
		"оbject/pendings/removed/count",
		"How many index-records have been saved in in-indexStorage index storage",
		stats.UnitDimensionless,
	)
	statRecordInMemoryAddedCount = stats.Int64(
		"recordstorage/added/count",
		"How many records have been saved to a in-memory storage",
		stats.UnitDimensionless,
	)
	statRecordInMemoryRemovedCount = stats.Int64(
		"recordstorage/removed/count",
		"How many records have been removed from a in-memory storage",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statBucketAddedCount.Name(),
			Description: statBucketAddedCount.Description(),
			Measure:     statBucketAddedCount,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statBucketRemovedCount.Name(),
			Description: statBucketRemovedCount.Description(),
			Measure:     statBucketRemovedCount,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statRecordInMemoryAddedCount.Name(),
			Description: statRecordInMemoryAddedCount.Description(),
			Measure:     statRecordInMemoryAddedCount,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statRecordInMemoryRemovedCount.Name(),
			Description: statRecordInMemoryRemovedCount.Description(),
			Measure:     statRecordInMemoryRemovedCount,
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		panic(err)
	}
}
