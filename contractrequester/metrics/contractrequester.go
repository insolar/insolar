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
		"contract_requester_handle_parsing_error",
		"handling messages that could not be parsed",
		stats.UnitDimensionless,
	)
)

var (
	CallMethodName = insmetrics.MustTagKey("contract_requester_method_name")
	CallReturnMode = insmetrics.MustTagKey("contract_requester_method_return_mode")
)
