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

package contractrequester

import (
	"context"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/pkg/errors"
)

type HandleReturnResults struct {
	dep *Dependencies

	Message bus.Message
}

func (h *HandleReturnResults) Present(ctx context.Context, f flow.Flow) error {
	parcel := h.Message.Parcel
	cr := h.dep.cr
	inslogger.FromContext(ctx).Debug("HandleReturnResults.Present starts ...")
	replyOk := bus.Reply{Reply: &reply.OK{}, Err: nil}

	msg, ok := parcel.Message().(*message.ReturnResults)
	if !ok {
		return errors.New("HandleReturnResults accepts only message.ReturnResults")
	}

	ctx, span := instracer.StartSpan(ctx, "ContractRequester.ReceiveResult")
	defer span.End()

	cr.ResultMutex.Lock()
	defer cr.ResultMutex.Unlock()

	logger := inslogger.FromContext(ctx)
	c, ok := cr.ResultMap[msg.Sequence]
	if !ok {
		logger.Info("oops unwaited results seq=", msg.Sequence)
		h.Message.ReplyTo <- replyOk
		return nil
	}
	logger.Debug("Got wanted results seq=", msg.Sequence)

	c <- msg
	delete(cr.ResultMap, msg.Sequence)

	h.Message.ReplyTo <- replyOk
	return nil

}
