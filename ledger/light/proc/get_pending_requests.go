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

type GetPendingRequests struct {
	replyTo  chan<- bus.Reply
	msg      *message.GetPendingRequests
	jet      insolar.JetID
	reqPulse insolar.PulseNumber

	Dep struct {
		RecentStorageProvider recentstorage.Provider
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
	p.replyTo <- p.reply(ctx)
	return nil
}

func (p *GetPendingRequests) reply(ctx context.Context) bus.Reply {
	msg := p.msg
	jetID := insolar.ID(p.jet)

	hasPendingRequests := false
	pendingStorage := p.Dep.RecentStorageProvider.GetPendingStorage(ctx, jetID)
	for _, reqID := range pendingStorage.GetRequestsForObject(*msg.Object.Record()) {
		if reqID.Pulse() < p.reqPulse {
			hasPendingRequests = true
		}
	}
	return bus.Reply{Reply: &reply.HasPendingRequests{Has: hasPendingRequests}}
}
