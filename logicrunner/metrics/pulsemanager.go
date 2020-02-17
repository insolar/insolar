// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package metrics

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	PulseManagerOnPulseTiming = stats.Float64(
		"vm_pulse_manager_onpulse_timing",
		"time spent on pulse set in pulsemanager",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        PulseManagerOnPulseTiming.Name(),
			Description: PulseManagerOnPulseTiming.Description(),
			Measure:     PulseManagerOnPulseTiming,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000, 10000, 20000),
		},
	)
	if err != nil {
		panic(err)
	}
}
