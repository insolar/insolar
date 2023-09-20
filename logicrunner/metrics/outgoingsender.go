package metrics

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	OutgoingSenderActorGoroutines = stats.Int64(
		"vm_outgoing_sender_goroutines",
		"OutgoingSender goroutines",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        OutgoingSenderActorGoroutines.Name(),
			Description: OutgoingSenderActorGoroutines.Description(),
			Measure:     OutgoingSenderActorGoroutines,
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		panic(err)
	}
}
