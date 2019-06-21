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

package recentstorage

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	tagJet = insmetrics.MustTagKey("jet")
)

var (
	statRecentStorageObjectsAdded   = stats.Int64("storage/recent/objects/added/count", "recent storage objects added", stats.UnitDimensionless)
	statRecentStorageObjectsRemoved = stats.Int64("storage/recent/objects/removed/count", "recent storage objects removed", stats.UnitDimensionless)

	statRecentStoragePendingsAdded   = stats.Int64("storage/recent/pending/added/count", "recent storage pending requests added", stats.UnitDimensionless)
	statRecentStoragePendingsRemoved = stats.Int64("storage/recent/pending/removed/count", "recent storage pending requests removed", stats.UnitDimensionless)
)

func init() {
	commontags := []tag.Key{tagJet}
	err := view.Register(
		&view.View{
			Name:        statRecentStorageObjectsAdded.Name(),
			Description: statRecentStorageObjectsAdded.Description(),
			Measure:     statRecentStorageObjectsAdded,
			Aggregation: view.Sum(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        statRecentStorageObjectsRemoved.Name(),
			Description: statRecentStorageObjectsRemoved.Description(),
			Measure:     statRecentStorageObjectsRemoved,
			Aggregation: view.Sum(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        statRecentStoragePendingsAdded.Name(),
			Description: statRecentStoragePendingsAdded.Description(),
			Measure:     statRecentStoragePendingsAdded,
			Aggregation: view.Sum(),
			TagKeys:     commontags,
		},
		&view.View{
			Name:        statRecentStoragePendingsRemoved.Name(),
			Description: statRecentStoragePendingsRemoved.Description(),
			Measure:     statRecentStoragePendingsRemoved,
			Aggregation: view.Sum(),
			TagKeys:     commontags,
		},
	)
	if err != nil {
		panic(err)
	}
}
