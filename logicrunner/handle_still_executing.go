// Copyright 2020 Insolar Network Ltd.
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

package logicrunner

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/metrics"
)

func checkPayloadStillExecuting(msg payload.StillExecuting) error {
	if msg.ObjectRef.IsEmpty() {
		return errors.New("Got StillExecuting message, but field ObjectRef is empty")
	}
	if !msg.ObjectRef.IsObjectReference() {
		return errors.Errorf("StillExecuting.ObjectRef should be ObjectReference; ref=%s", msg.ObjectRef.String())
	}

	if !msg.Executor.IsObjectReference() {
		return errors.Errorf("StillExecuting.Executor should be ObjectReference; ref=%s", msg.Executor.String())
	}

	if len(msg.RequestRefs) == 0 {
		return errors.New("StillExecuting.RequestRefs should have list of elements, got empty list")
	}

	for _, requestRef := range msg.RequestRefs {
		if !requestRef.IsRecordScope() {
			return errors.Errorf("StillExecuting.RequestRefs should have only RecordReferences; ref=%s", requestRef.String())
		}
	}

	return nil
}

type HandleStillExecuting struct {
	dep *Dependencies

	Message payload.Meta
}

func (h *HandleStillExecuting) Present(ctx context.Context, f flow.Flow) error {
	message := payload.StillExecuting{}
	err := message.Unmarshal(h.Message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message")
	}

	ctx, logger := inslogger.WithFields(ctx, map[string]interface{}{
		"object": message.ObjectRef.String(),
		"sender": h.Message.Sender.String(),
	})
	logger.Debug("handle StillExecuting message")

	if err := checkPayloadStillExecuting(message); err != nil {
		return nil
	}

	done, err := h.dep.WriteAccessor.Begin(ctx, flow.Pulse(ctx))
	if err != nil {
		logger.Warn("late StillExecuting message, ignoring: ", err.Error())
		return nil
	}
	defer done()

	h.dep.ResultsMatcher.AddStillExecution(ctx, message)

	broker := h.dep.StateStorage.UpsertExecutionState(message.ObjectRef)
	err = broker.PrevExecutorStillExecuting(ctx)
	if err != nil {
		logger.Warn(err)
		if err == ErrNotInPending {
			stats.Record(ctx, metrics.StillExecutingAlreadyExecuting.M(1))
		}
	}

	return nil
}
