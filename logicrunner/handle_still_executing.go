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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type HandleStillExecuting struct {
	dep *Dependencies

	Message payload.Meta
	Parcel  insolar.Parcel
}

func (h *HandleStillExecuting) Present(ctx context.Context, f flow.Flow) error {
	logger := inslogger.FromContext(ctx)
	lr := h.dep.lr
	replyOk := bus.ReplyAsMessage(ctx, &reply.OK{})

	inslogger.FromContext(ctx).Debug("HandleStillExecuting.Present starts ...")

	msg := h.Parcel.Message().(*message.StillExecuting)
	ref := msg.DefaultTarget()
	broker := lr.StateStorage.UpsertExecutionState(*ref)
	es := &broker.executionState

	logger.Debugf("Got information that %s is still executing", ref.String())

	es.Lock()
	switch es.pending {
	case insolar.NotPending:
		// It might be when StillExecuting comes after PendingFinished
		logger.Error("got StillExecuting message, but our state says that it's not in pending")
	case insolar.InPending:
		es.PendingConfirmed = true
	case insolar.PendingUnknown:
		// we are first, strange, soon ExecuteResults message should come
		es.pending = insolar.InPending
		es.PendingConfirmed = true
	}

	es.Unlock()
	h.dep.Sender.Reply(ctx, h.Message, replyOk)

	return nil

}
