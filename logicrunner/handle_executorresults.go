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

package logicrunner

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/metrics"
	"github.com/insolar/insolar/logicrunner/writecontroller"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

type HandleExecutorResults struct {
	dep *Dependencies

	meta payload.Meta
}

func checkPayloadExecutorResults(ctx context.Context, results payload.ExecutorResults) error {
	if !results.Caller.IsEmpty() && !results.Caller.IsObjectReference() {
		return errors.Errorf("results.Caller should be ObjectReference; ref=%s", results.Caller.String())
	}
	if !results.RecordRef.IsObjectReference() {
		return errors.Errorf("results.RecordRef should be ObjectReference; ref=%s", results.RecordRef.String())
	}

	for _, elem := range results.Queue {
		if !elem.RequestRef.IsRecordScope() {
			return errors.Errorf("results.RecordRef should be RecordReference; ref=%s", results.RecordRef.String())
		}
		if err := checkIncomingRequest(ctx, elem.Incoming); err != nil {
			return errors.Wrap(err, "failed to check ExecutionQueue of ExecutorResults")
		}
	}

	return nil
}

func (h *HandleExecutorResults) Present(ctx context.Context, f flow.Flow) error {
	message := payload.ExecutorResults{}
	err := message.Unmarshal(h.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message")
	}

	ctx, logger := inslogger.WithField(ctx, "object", message.RecordRef.String())
	logger.Debug("handling ExecutorResults")

	ctx, span := instracer.StartSpan(ctx, "HandleExecutorResults.Present")
	defer span.End()

	if err := checkPayloadExecutorResults(ctx, message); err != nil {
		return err
	}

	return h.handleMessage(ctx, message)
}

func (h *HandleExecutorResults) handleMessage(ctx context.Context, msg payload.ExecutorResults) error {
	done, err := h.dep.WriteAccessor.Begin(ctx, flow.Pulse(ctx))
	if err != nil {
		if err == writecontroller.ErrWriteClosed {
			return flow.ErrCancelled
		}
		return errors.Wrap(err, "failed to acquire write access")
	}
	defer done()

	broker := h.dep.StateStorage.UpsertExecutionState(msg.RecordRef)
	broker.PrevExecutorPendingResult(ctx, msg.Pending)

	if msg.LedgerHasMoreRequests {
		stats.Record(ctx, metrics.ExecutorResultsRequestsFromPrevExecutor.M(1))
		broker.MoreRequestsOnLedger(ctx)

		return nil
	}

	if len(msg.Queue) > 0 {
		transcripts := make([]*common.Transcript, len(msg.Queue))
		for i, qe := range msg.Queue {
			transcripts[i] = common.NewTranscriptCloneContext(qe.ServiceData, qe.RequestRef, *qe.Incoming)
		}
		broker.AddRequestsFromPrevExecutor(ctx, transcripts...)
	}

	return nil
}
