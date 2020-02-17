// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package heavy

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statBadgerStartTime = stats.Float64(
		"badger_start_time",
		"Time of last badger starting",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statBadgerStartTime.Name(),
			Description: statBadgerStartTime.Description(),
			Measure:     statBadgerStartTime,
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}
