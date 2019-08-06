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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/common"
)

type HandlePendingFinished struct {
	dep *Dependencies

	Message payload.Meta
	Parcel  insolar.Parcel
}

func (h *HandlePendingFinished) Present(ctx context.Context, _ flow.Flow) error {
	ctx = common.LoggerWithTargetID(ctx, h.Parcel)
	inslogger.FromContext(ctx).Debug("HandlePendingFinished.Present starts ...")

	msg := h.Parcel.Message().(*message.PendingFinished)

	broker := h.dep.StateStorage.UpsertExecutionState(msg.Reference)

	err := broker.PrevExecutorFinishedPending(ctx)
	if err != nil {
		err = errors.Wrap(err, "can not finish pending")
		inslogger.FromContext(ctx).Error(err.Error())
		return err
	}

	replyOk := bus.ReplyAsMessage(ctx, &reply.OK{})
	h.dep.Sender.Reply(ctx, h.Message, replyOk)
	return nil
}
