// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statReceivedHeavyPayloadCount = stats.Int64(
		"heavysyncer_heavypayload_count",
		"How many heavy-payload messages were received from a light-node",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statReceivedHeavyPayloadCount.Name(),
			Description: statReceivedHeavyPayloadCount.Description(),
			Measure:     statReceivedHeavyPayloadCount,
			Aggregation: view.Count(),
		},
	)
	if err != nil {
		panic(err)
	}
}
