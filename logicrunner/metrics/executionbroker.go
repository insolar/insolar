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
	TagExecutionQueueName = insmetrics.MustTagKey("vm_execution_broker_queue_name")
)

var (
	ExecutionBrokerOnPulseTiming = stats.Float64(
		"vm_execution_broker_onpulse_timing",
		"time spent on pulse in execution broker",
		stats.UnitMilliseconds,
	)
	ExecutionBrokerOnPulseNotConfirmed = stats.Int64(
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
		"vm_execution_broker_execution_finished",
		"execution broker started execution",
		stats.UnitDimensionless,
	)
)

func init() {
	err := view.Register(
		&view.View{
			Name:        ExecutionBrokerOnPulseTiming.Name(),
			Description: ExecutionBrokerOnPulseTiming.Description(),
			Measure:     ExecutionBrokerOnPulseTiming,
			Aggregation: view.Distribution(0.001, 0.01, 0.1, 1, 10, 100, 1000, 5000, 10000, 20000),
		},
		&view.View{
			Name:        ExecutionBrokerOnPulseNotConfirmed.Name(),
			Description: ExecutionBrokerOnPulseNotConfirmed.Description(),
			Measure:     ExecutionBrokerOnPulseNotConfirmed,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        ExecutionBrokerTruncatedRequests.Name(),
			Description: ExecutionBrokerTruncatedRequests.Description(),
			Measure:     ExecutionBrokerTruncatedRequests,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        ExecutionBrokerTranscriptRegistered.Name(),
			Description: ExecutionBrokerTranscriptRegistered.Description(),
			Measure:     ExecutionBrokerTranscriptRegistered,
			TagKeys:     []tag.Key{TagExecutionQueueName},
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        ExecutionBrokerTranscriptDuplicate.Name(),
			Description: ExecutionBrokerTranscriptDuplicate.Description(),
			Measure:     ExecutionBrokerTranscriptDuplicate,
			TagKeys:     []tag.Key{TagExecutionQueueName},
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        ExecutionBrokerTranscriptExecuting.Name(),
			Description: ExecutionBrokerTranscriptExecuting.Description(),
			Measure:     ExecutionBrokerTranscriptExecuting,
			TagKeys:     []tag.Key{TagExecutionQueueName},
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        ExecutionBrokerTranscriptAlreadyRegistered.Name(),
			Description: ExecutionBrokerTranscriptAlreadyRegistered.Description(),
			Measure:     ExecutionBrokerTranscriptAlreadyRegistered,
			TagKeys:     []tag.Key{TagExecutionQueueName},
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        ExecutionBrokerExecutionStarted.Name(),
			Description: ExecutionBrokerExecutionStarted.Description(),
			Measure:     ExecutionBrokerExecutionStarted,
			Aggregation: view.Sum(),
		},
		&view.View{
			Name:        ExecutionBrokerExecutionFinished.Name(),
			Description: ExecutionBrokerExecutionFinished.Description(),
			Measure:     ExecutionBrokerExecutionFinished,
			Aggregation: view.Sum(),
		},
	)
	if err != nil {
		panic(err)
	}
}
