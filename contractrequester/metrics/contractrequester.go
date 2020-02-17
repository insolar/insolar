// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package metrics

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/insolar/insolar/instrumentation/insmetrics"
)

var (
	TagFinishedWithError = insmetrics.MustTagKey("vm_handle_finished_with_error")
)

// common handling stats
var (
	HandleStarted = stats.Int64(
		"contract_requester_handle_started",
		"started message handling",
		stats.UnitDimensionless,
	)
	HandleFinished = stats.Int64(
		"contract_requester_handle_finished",
		"finished message handling",
		stats.UnitDimensionless,
	)
	HandleFuture = stats.Int64(
		"contract_requester_handle_future",
		"handling messages from future pulse",
		stats.UnitDimensionless,
	)
	HandlePast = stats.Int64(
		"contract_requester_handle_past",
		"handling messages from previous pulse",
		stats.UnitDimensionless,
	)
	HandlePastFlowCancelled = stats.Int64(
		"contract_requester_handle_past_flow_cancelled",
		"handling messages from previous pulse that flow cancelled",
		stats.UnitDimensionless,
	)
	HandleTiming = stats.Float64(
		"contract_requester_handle_latency",
		"time spent on handling messages",
		stats.UnitMilliseconds,
	)

	HandlingParsingError = stats.Int64(
		"contract_requester_handle_parsing_error",
		"handling messages that could not be parsed",
		stats.UnitDimensionless,
	)
)

var (
	SendMessageTiming = stats.Float64(
		"contract_requester_send_message_timing",
		"SendMessage timing",
		stats.UnitDimensionless,
	)
)

var (
	CallMethodName = insmetrics.MustTagKey("contract_requester_method_name")
	CallReturnMode = insmetrics.MustTagKey("contract_requester_method_return_mode")
)

func init() {
	err := view.Register(
		&view.View{
			Name:        HandlingParsingError.Name(),
			Description: HandlingParsingError.Description(),
			Measure:     HandlingParsingError,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        HandleTiming.Name(),
			Description: HandleTiming.Description(),
			Measure:     HandleTiming,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000, 10000, 20000),
		},
		&view.View{
			Name:        HandleStarted.Name(),
			Description: HandleStarted.Description(),
			Measure:     HandleStarted,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        HandlePast.Name(),
			Description: HandlePast.Description(),
			Measure:     HandlePast,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        HandlePastFlowCancelled.Name(),
			Description: HandlePastFlowCancelled.Description(),
			Measure:     HandlePastFlowCancelled,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        HandleFuture.Name(),
			Description: HandleFuture.Description(),
			Measure:     HandleFuture,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        HandleFinished.Name(),
			Description: HandleFinished.Description(),
			Measure:     HandleFinished,
			TagKeys:     []tag.Key{TagFinishedWithError},
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        SendMessageTiming.Name(),
			Description: SendMessageTiming.Description(),
			Measure:     SendMessageTiming,
			TagKeys:     []tag.Key{CallMethodName, CallReturnMode},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000, 10000, 20000),
		},
	)
	if err != nil {
		panic(err)
	}
}
