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

package proc

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/object"
)

type GetPendingRequests struct {
	message  payload.Meta
	msg      *message.GetPendingRequests
	jet      insolar.JetID
	reqPulse insolar.PulseNumber

	dep struct {
		index  object.IndexAccessor
		sender bus.Sender
	}
}

func NewGetPendingRequests(jetID insolar.JetID, message payload.Meta, msg *message.GetPendingRequests, reqPulse insolar.PulseNumber) *GetPendingRequests {
	return &GetPendingRequests{
		msg:      msg,
		message:  message,
		jet:      jetID,
		reqPulse: reqPulse,
	}
}

func (p *GetPendingRequests) Dep(index object.IndexAccessor, sender bus.Sender) {
	p.dep.index = index
	p.dep.sender = sender
}

func (p *GetPendingRequests) Proceed(ctx context.Context) error {
	idx, err := p.dep.index.ForID(ctx, flow.Pulse(ctx), *p.msg.Object.Record())
	if err != nil {
		return err
	}
	rep := bus.ReplyAsMessage(ctx, &reply.HasPendingRequests{
		Has: idx.Lifeline.EarliestOpenRequest != nil && *idx.Lifeline.EarliestOpenRequest < flow.Pulse(ctx),
	})
	go p.dep.sender.Reply(ctx, p.message, rep)
	return nil
}
