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
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/light/recentstorage"
)

type GetPendingRequestID struct {
	replyTo  chan<- bus.Reply
	msg      *message.GetPendingRequestID
	jet      insolar.JetID
	reqPulse insolar.PulseNumber

	Dep struct {
		RecentStorageProvider recentstorage.Provider
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

func (p *GetPendingRequestID) Proceed(ctx context.Context) error {
	msg := p.msg
	jetID := insolar.ID(p.jet)

	requests := p.Dep.RecentStorageProvider.GetPendingStorage(ctx, jetID).GetRequestsForObject(msg.ObjectID)
	if len(requests) == 0 {
		p.replyTo <- bus.Reply{Reply: &reply.Error{ErrType: reply.ErrNoPendingRequests}}
		return nil
	}
	p.replyTo <- bus.Reply{Reply: &reply.ID{ID: requests[0]}}

	return nil
}
