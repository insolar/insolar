package pulsenetwork

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statSendPulseErrorsCount = stats.Int64("pulsar_sending_pulse_errors", "count of errors while sending pulse", stats.UnitDimensionless)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statSendPulseErrorsCount.Name(),
			Description: statSendPulseErrorsCount.Description(),
			Measure:     statSendPulseErrorsCount,
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		panic(err)
	}
}
