// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package metrics

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

// IncomingRequests related stats
var (
	IncomingRequestsNew = stats.Int64(
		"vm_incoming_request_new",
		"New IncomingRequests created",
		stats.UnitDimensionless,
	)
	IncomingRequestsDuplicate = stats.Int64(
		"vm_incoming_request_duplicate",
		"Duplicated IncomingRequests registered",
		stats.UnitDimensionless,
	)
	IncomingRequestsClosed = stats.Int64(
		"vm_incoming_request_closed",
		"Duplicated IncomingRequests with results",
		stats.UnitDimensionless,
	)
)

// OutgoingRequests related stats
var (
	OutgoingRequestsNew = stats.Int64(
		"vm_outgoing_request_new",
		"New OutgoingRequests created",
		stats.UnitDimensionless,
	)
	OutgoingRequestsDuplicate = stats.Int64(
		"vm_outgoing_request_duplicate",
		"Duplicated OutgoingRequests registered",
		stats.UnitDimensionless,
	)
	OutgoingRequestsClosed = stats.Int64(
		"vm_outgoing_request_closed",
		"Duplicated OutgoingRequests with results",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        IncomingRequestsNew.Name(),
			Description: IncomingRequestsNew.Description(),
			Measure:     IncomingRequestsNew,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        IncomingRequestsDuplicate.Name(),
			Description: IncomingRequestsDuplicate.Description(),
			Measure:     IncomingRequestsDuplicate,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        IncomingRequestsClosed.Name(),
			Description: IncomingRequestsClosed.Description(),
			Measure:     IncomingRequestsClosed,
			Aggregation: view.Sum(),
		},

		&view.View{
			Name:        OutgoingRequestsNew.Name(),
			Description: OutgoingRequestsNew.Description(),
			Measure:     OutgoingRequestsNew,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        OutgoingRequestsDuplicate.Name(),
			Description: OutgoingRequestsDuplicate.Description(),
			Measure:     OutgoingRequestsDuplicate,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        OutgoingRequestsClosed.Name(),
			Description: OutgoingRequestsClosed.Description(),
			Measure:     OutgoingRequestsClosed,
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		panic(err)
	}
}
