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
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/pkg/errors"
)

type GetPendingRequestID struct {
	replyTo  chan<- bus.Reply
	msg      *message.GetPendingRequestID
	jet      insolar.JetID
	reqPulse insolar.PulseNumber

	dep struct {
		filaments executor.FilamentCalculator
	}
}

func NewGetPendingRequestID(jetID insolar.JetID, replyTo chan<- bus.Reply, msg *message.GetPendingRequestID, reqPulse insolar.PulseNumber) *GetPendingRequestID {
	return &GetPendingRequestID{
		msg:      msg,
		replyTo:  replyTo,
		jet:      jetID,
		reqPulse: reqPulse,
	}
}

func (p *GetPendingRequestID) Dep(filaments executor.FilamentCalculator) {
	p.dep.filaments = filaments
}

func (p *GetPendingRequestID) Proceed(ctx context.Context) error {
	ids, err := p.dep.filaments.PendingRequests(ctx, flow.Pulse(ctx), p.msg.ObjectID)
	if err != nil {
		return errors.Wrap(err, "failed to calculate pending")
	}
	if len(ids) == 0 {
		p.replyTo <- bus.Reply{Reply: &reply.Error{ErrType: reply.ErrNoPendingRequests}}
		return nil
	}

	p.replyTo <- bus.Reply{Reply: &reply.ID{ID: ids[0]}}
	return nil
}
