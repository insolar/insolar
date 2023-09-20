package pulsar

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	statPulseGenerated = stats.Int64("pulsar_pulse_generated", "count of generated pulses", stats.UnitDimensionless)
	statCurrentPulse   = stats.Int64("pulsar_current_pulse", "last generated pulse", stats.UnitDimensionless)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statPulseGenerated.Name(),
			Description: statPulseGenerated.Description(),
			Measure:     statPulseGenerated,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        statCurrentPulse.Name(),
			Description: statCurrentPulse.Description(),
			Measure:     statCurrentPulse,
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}
