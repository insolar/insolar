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

package handle

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetObject struct {
	dep *proc.Dependencies

	Message bus.Message
}

func (s *GetObject) Present(ctx context.Context, f flow.Flow) error {
	msg := s.Message.Parcel.Message().(*message.GetObject)
	ctx, _ = inslogger.WithField(ctx, "object", msg.Head.Record().DebugString())

	// Workaround to fetch object states.
	var pn insolar.PulseNumber
	if s.Message.Parcel.DelegationToken() == nil {
		pn = flow.Pulse(ctx)
	} else {
		pn = msg.State.Pulse()
	}
	jet := proc.NewFetchJet(*msg.Head.Record(), pn, s.Message.ReplyTo)
	s.dep.FetchJet(jet)
	if err := jet.Proceed(ctx); err != nil {
		return err
	}

	hot := proc.NewWaitHot(jet.Result.Jet, flow.Pulse(ctx), s.Message.ReplyTo)
	s.dep.WaitHot(hot)
	if err := jet.Proceed(ctx); err != nil {
		return err
	}

	idx := proc.NewGetIndex(msg.Head, jet.Result.Jet, s.Message.ReplyTo)
	s.dep.GetIndex(idx)
	if err := idx.Proceed(ctx); err != nil {
		return err
	}

	send := proc.NewSendObject(s.Message, jet.Result.Jet, idx.Result.Index)
	s.dep.SendObject(send)
	return send.Proceed(ctx)
}
