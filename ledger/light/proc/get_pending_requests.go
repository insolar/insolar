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
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/object"
)

type GetPendingRequests struct {
	replyTo  chan<- bus.Reply
	msg      *message.GetPendingRequests
	jet      insolar.JetID
	reqPulse insolar.PulseNumber

	Dep struct {
		PendingAccessor object.PendingAccessor
	}
}

func NewGetPendingRequests(jetID insolar.JetID, replyTo chan<- bus.Reply, msg *message.GetPendingRequests, reqPulse insolar.PulseNumber) *GetPendingRequests {
	return &GetPendingRequests{
		msg:      msg,
		replyTo:  replyTo,
		jet:      jetID,
		reqPulse: reqPulse,
	}
}

func (p *GetPendingRequests) Proceed(ctx context.Context) error {
	msg := p.msg

	pends, err := p.Dep.PendingAccessor.OpenRequestsForObjID(ctx, flow.Pulse(ctx), *msg.Object.Record(), 1)
	if err != nil || p == nil || len(pends) == 0 {
		p.replyTo <- bus.Reply{Reply: &reply.HasPendingRequests{Has: false}}
		return nil
	}
	p.replyTo <- bus.Reply{Reply: &reply.HasPendingRequests{Has: true}}
	return nil
}
