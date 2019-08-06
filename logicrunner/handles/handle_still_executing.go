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

package handles

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/logicrunner/common"
)

type HandleStillExecuting struct {
	dep *Dependencies

	Message payload.Meta
	Parcel  insolar.Parcel
}

func (h *HandleStillExecuting) Present(ctx context.Context, f flow.Flow) error {
	ctx = common.LoggerWithTargetID(ctx, h.Parcel)
	msg := h.Parcel.Message().(*message.StillExecuting)
	h.dep.ResultsMatcher.AddStillExecution(ctx, msg)

	broker := h.dep.StateStorage.UpsertExecutionState(msg.Reference)
	broker.PrevExecutorStillExecuting(ctx)

	// logger := inslogger.FromContext(ctx)
	// for _, reqRef := range msg.RequestRefs {
	// 	if broker.IsKnownRequest(ctx, reqRef) {
	// 		logger.Debug("skipping known request ", reqRef.String())
	// 		continue
	// 	}
	//
	// 	request, err := broker.Req.am.GetIncomingRequest(ctx, rf.object, reqRef)
	// 	if err != nil {
	// 		logger.Error("couldn't get request: ", err.Error())
	// 		continue
	// 	}
	//
	// 	select {
	// 	case <-ctx.Done():
	// 		logger.Debug("quiting fetching requests, was stopped")
	// 		return nil
	// 	default:
	// 	}
	//
	// 	logger.Errorf("fetch req from ledger: %s, %s", reqRef.String(), request.String())
	// 	requestCtx := freshContextFromContextAndRequest(ctx, *request)
	// 	tr := NewTranscript(requestCtx, reqRef, *request)
	// 	rf.broker.AddRequestsFromLedger(ctx, tr)
	// 	tr := NewTranscript(ctx, msg.Reference, reqRef)
	// 	broker.AddAdditionalRequestFromPrevExecutor(ctx, reqRef)
	// }
	replyOk := bus.ReplyAsMessage(ctx, &reply.OK{})
	h.dep.Sender.Reply(ctx, h.Message, replyOk)
	return nil
}
