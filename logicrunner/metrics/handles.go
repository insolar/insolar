//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package metrics

import (
	"go.opencensus.io/stats"

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
