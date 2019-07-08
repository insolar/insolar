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
)

type HandleStillExecuting struct {
	dep *Dependencies

	Message bus.Message
}

func (h *HandleStillExecuting) Present(ctx context.Context, f flow.Flow) error {
	logger := inslogger.FromContext(ctx)
	parcel := h.Message.Parcel
	ctx = loggerWithTargetID(ctx, parcel)
	lr := h.dep.lr
	replyOk := bus.Reply{Reply: &reply.OK{}, Err: nil}

	logger.Debug("HandleStillExecuting.Present starts ...")

	msg := parcel.Message().(*message.StillExecuting)
	ref := msg.DefaultTarget()
	os := lr.StateStorage.UpsertObjectState(*ref)

	logger.Debugf("Got informatlogicrunner/executionstate_test.goion that %s is still executing", ref.String())

	es, _ := os.InitAndGetExecution(lr, ref)

	es.Lock()
	switch os.ExecutionState.pending {
	case message.NotPending:
		// It might be when StillExecuting comes after PendingFinished
		logger.Error("got StillExecuting message, but our state says that it's not in pending")
	case message.InPending:
		os.ExecutionState.PendingConfirmed = true
	case message.PendingUnknown:
		// we are first, strange, soon ExecuteResults message should come
		os.ExecutionState.pending = message.InPending
		os.ExecutionState.PendingConfirmed = true
	}
	es.Unlock()

	h.Message.ReplyTo <- replyOk
	return nil

}
