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

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

type HandlePendingFinished struct {
	dep *Dependencies

	Message bus.Message
}

func (h *HandlePendingFinished) Present(ctx context.Context, f flow.Flow) error {
	parcel := h.Message.Parcel
	ctx = loggerWithTargetID(ctx, parcel)
	lr := h.dep.lr
	inslogger.FromContext(ctx).Debug("HandlePendingFinished.Present starts ...")
	replyOk := bus.Reply{Reply: &reply.OK{}, Err: nil}

	msg := parcel.Message().(*message.PendingFinished)
	ref := msg.DefaultTarget()
	os := lr.StateStorage.UpsertObjectState(*ref)

	os.Lock()
	if os.ExecutionState == nil {
		// we are first, strange, soon ExecuteResults message should come
		os.ExecutionState = NewExecutionState(*ref)
		os.ExecutionState.pending = message.NotPending
		os.ExecutionState.RegisterLogicRunner(lr)
		os.Unlock()

		h.Message.ReplyTo <- replyOk
		return nil
	}
	es := os.ExecutionState
	os.Unlock()

	es.Lock()
	es.pending = message.NotPending
	if !es.Broker.currentList.Empty() {
		es.Unlock()
		return errors.New("[ HandlePendingFinished ] received PendingFinished when we are already executing")
	}
	es.Unlock()

	es.Broker.StartProcessorIfNeeded(ctx)

	h.Message.ReplyTo <- replyOk
	return nil

}
