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
	TagExecutionBrokerName = insmetrics.MustTagKey("vm_execution_broker_name")

	TagExecutionQueueName = insmetrics.MustTagKey("vm_execution_broker_queue_name")
)

var (
	ExecutionBrokerOnPulseTiming = stats.Float64(
		"vm_execution_broker_onpulse_timing",
		"time spent on pulse in execution broker",
		stats.UnitMilliseconds,
	)
	ExecutionBrokerOnPulseNotConfirmed = stats.Float64(
		"vm_execution_broker_onpulse_notconfirmed",
		"not confirmed execution brokers",
		stats.UnitDimensionless,
	)
)

var (
	ExecutionBrokerTruncatedRequests = stats.Int64(
		"vm_execution_broker_onpulse_truncated",
		"execution broker truncated requests onpulse",
		stats.UnitDimensionless,
	)
	ExecutionBrokerTranscriptRegistered = stats.Int64(
		"vm_execution_broker_transcript_new",
		"execution broker new transcript registered",
		stats.UnitDimensionless,
	)
	ExecutionBrokerTranscriptDuplicate = stats.Int64(
		"vm_execution_broker_transcript_duplicate",
		"execution broker duplicate transcript registered",
		stats.UnitDimensionless,
	)
	ExecutionBrokerTranscriptExecuting = stats.Int64(
		"vm_execution_broker_transcript_executing",
		"execution broker already executing transcript",
		stats.UnitDimensionless,
	)
	ExecutionBrokerTranscriptAlreadyRegistered = stats.Int64(
		"vm_execution_broker_transcript_already_registered",
		"execution broker already registered transcript",
		stats.UnitDimensionless,
	)
)

var (
	ExecutionBrokerExecutionStarted = stats.Int64(
		"vm_execution_broker_execution_started",
		"execution broker started execution",
		stats.UnitDimensionless,
	)
	ExecutionBrokerExecutionFinished = stats.Int64(
		"vm_execution_broker_execution_started",
		"execution broker started execution",
		stats.UnitDimensionless,
	)
)
