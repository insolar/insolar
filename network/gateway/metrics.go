// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package gateway

import (
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	tagNodeRef = insmetrics.MustTagKey("nodeRef")
)

var (
	statPulse = stats.Int64(
		"current_pulse",
		"current node pulse",
		stats.UnitDimensionless,
	)
	networkState = stats.Int64(
		"network_state",
		"current network state",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statPulse.Name(),
			Description: statPulse.Description(),
			Measure:     statPulse,
			Aggregation: view.LastValue(),
		},
		&view.View{
			Name:        networkState.Name(),
			Description: networkState.Description(),
			Measure:     networkState,
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}
