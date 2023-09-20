package metrics

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	LogicRunnerOnPulseTiming = stats.Float64(
		"vm_logicrunner_onpulse_timing",
		"time spent on pulse change of LogicRunner",
		stats.UnitMilliseconds,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        LogicRunnerOnPulseTiming.Name(),
			Description: LogicRunnerOnPulseTiming.Description(),
			Measure:     LogicRunnerOnPulseTiming,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000, 10000, 20000),
		},
	)
	if err != nil {
		panic(err)
	}
}
