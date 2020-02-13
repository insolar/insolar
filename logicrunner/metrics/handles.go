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
	TagHandlePayloadType = insmetrics.MustTagKey("vm_handle_payload_type")

	TagFinishedWithError = insmetrics.MustTagKey("vm_handle_finished_with_error")
)

// common handling stats
var (
	HandleStarted = stats.Int64(
		"vm_handle_started",
		"started message handling",
		stats.UnitDimensionless,
	)
	HandleFinished = stats.Int64(
		"vm_handle_finished",
		"finished message handling",
		stats.UnitDimensionless,
	)
	HandleFuture = stats.Int64(
		"vm_handle_future",
		"handling messages from future pulse",
		stats.UnitDimensionless,
	)
	HandlePast = stats.Int64(
		"vm_handle_past",
		"handling messages from previous pulse",
		stats.UnitDimensionless,
	)
	HandlePastFlowCancelled = stats.Int64(
		"vm_handle_past_flow_cancelled",
		"handling messages from previous pulse that flow cancelled",
		stats.UnitDimensionless,
	)
	HandleTiming = stats.Float64(
		"vm_handle_latency",
		"time spent on handling messages",
		stats.UnitMilliseconds,
	)

	HandlingParsingError = stats.Int64(
		"vm_handle_parsing_error",
		"handling messages that could not be parsed",
		stats.UnitDimensionless,
	)
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
			TagKeys:     []tag.Key{TagHandlePayloadType},
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000, 10000, 20000),
		},
		&view.View{
			Name:        HandleStarted.Name(),
			Description: HandleStarted.Description(),
			Measure:     HandleStarted,
			TagKeys:     []tag.Key{TagHandlePayloadType},
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        HandlePast.Name(),
			Description: HandlePast.Description(),
			Measure:     HandlePast,
			TagKeys:     []tag.Key{TagHandlePayloadType},
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        HandlePastFlowCancelled.Name(),
			Description: HandlePastFlowCancelled.Description(),
			Measure:     HandlePastFlowCancelled,
			TagKeys:     []tag.Key{TagHandlePayloadType},
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        HandleFuture.Name(),
			Description: HandleFuture.Description(),
			Measure:     HandleFuture,
			TagKeys:     []tag.Key{TagHandlePayloadType},
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        HandleFinished.Name(),
			Description: HandleFinished.Description(),
			Measure:     HandleFinished,
			TagKeys:     []tag.Key{TagHandlePayloadType, TagFinishedWithError},
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        HandleUnknownMessageType.Name(),
			Description: HandleUnknownMessageType.Description(),
			Measure:     HandleUnknownMessageType,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        CallMethodLogicalError.Name(),
			Description: CallMethodLogicalError.Description(),
			Measure:     CallMethodLogicalError,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        CallMethodAdditionalCall.Name(),
			Description: CallMethodAdditionalCall.Description(),
			Measure:     CallMethodAdditionalCall,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        CallMethodLoopDetected.Name(),
			Description: CallMethodLoopDetected.Description(),
			Measure:     CallMethodLoopDetected,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        ExecutorResultsRequestsFromPrevExecutor.Name(),
			Description: ExecutorResultsRequestsFromPrevExecutor.Description(),
			Measure:     ExecutorResultsRequestsFromPrevExecutor,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        PendingFinishedAlreadyExecuting.Name(),
			Description: PendingFinishedAlreadyExecuting.Description(),
			Measure:     PendingFinishedAlreadyExecuting,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        StillExecutingAlreadyExecuting.Name(),
			Description: StillExecutingAlreadyExecuting.Description(),
			Measure:     StillExecutingAlreadyExecuting,
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		panic(err)
	}
}

// unknown message type error
var (
	HandleUnknownMessageType = stats.Int64(
		"vm_handle_unknown_message",
		"handling unknown message type",
		stats.UnitDimensionless,
	)
)

// CallMethod specific stats
var (
	CallMethodLogicalError = stats.Int64(
		"vm_call_method_logical_error",
		"call method with logical error",
		stats.UnitDimensionless,
	)
	CallMethodAdditionalCall = stats.Int64(
		"vm_call_method_additional_call",
		"call method with additional call",
		stats.UnitDimensionless,
	)
	CallMethodLoopDetected = stats.Int64(
		"vm_call_method_loop_detected",
		"call method with loop detected",
		stats.UnitDimensionless,
	)
)

// ExecutorResults specific stats
var (
	ExecutorResultsRequestsFromPrevExecutor = stats.Int64(
		"vm_executor_results_prev_executor",
		"ExecutorResults with AddRequestsFromPrevExecutor calls",
		stats.UnitDimensionless,
	)
)

// PendingFinished specific stats
var (
	PendingFinishedAlreadyExecuting = stats.Int64(
		"vm_pending_finished_already_executing",
		"PendingFinished with AlreadyExecuting error",
		stats.UnitDimensionless,
	)
)

// StillExecuting specific stats
var (
	StillExecutingAlreadyExecuting = stats.Int64(
		"vm_still_executing_already_executing",
		"StillExecuting with AlreadyExecuting error",
		stats.UnitDimensionless,
	)
)
