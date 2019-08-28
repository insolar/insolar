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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type HandleStillExecuting struct {
	dep *Dependencies

	Message payload.Meta
	Parcel  insolar.Parcel
}

func (h *HandleStillExecuting) Present(ctx context.Context, f flow.Flow) error {
	message := payload.StillExecuting{}
	err := message.Unmarshal(h.Message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message")
	}

	ctx, logger := inslogger.WithField(ctx, "object", message.ObjectRef.String())
	logger.Debug("handling still executing message")

	done, err := h.dep.WriteAccessor.Begin(ctx, flow.Pulse(ctx))
	if err != nil {
		logger.Warn("late still executing message, ignoring: ", err.Error())
		return nil
	}
	defer done()

	h.dep.ResultsMatcher.AddStillExecution(ctx, &message)

	broker := h.dep.StateStorage.UpsertExecutionState(message.ObjectRef)
	broker.PrevExecutorStillExecuting(ctx)

	return nil
}
