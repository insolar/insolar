/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */


package controller

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	tagMessageType = insmetrics.MustTagKey("messageType")
)

var (
	statParcelsSentSizeBytes = stats.Int64(
		"network/parcels/sent/size",
		"size of sent parcels",
		stats.UnitBytes,
	)
	statParcelsReplySizeBytes = stats.Int64(
		"network/parcels/reply/size",
		"size of replies to parcels",
		stats.UnitBytes,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Measure:     statParcelsSentSizeBytes,
			Aggregation: view.Distribution(16, 32, 64, 128, 256, 512, 1024, 16*1<<10, 512*1<<10, 1<<20 ),
			TagKeys:     []tag.Key{tagMessageType},
		},
	)
	if err != nil {
		panic(err)
	}
}

