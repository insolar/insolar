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
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetObject struct {
	dep *proc.Dependencies

	payload payload.GetObject
	message bus.Message
}

func NewGetObject(dep *proc.Dependencies, msg bus.Message, pl payload.GetObject) *GetObject {
	return &GetObject{
		dep:     dep,
		payload: pl,
		message: msg,
	}
}

func (s *GetObject) Present(ctx context.Context, f flow.Flow) error {
	ctx, _ = inslogger.WithField(ctx, "object", s.payload.ObjectID.DebugString())

	var (
		objJetID, stateJetID insolar.JetID
		stateID              insolar.ID
	)

	jet := proc.NewFetchJetWM(s.payload.ObjectID, flow.Pulse(ctx), s.message.WatermillMsg)
	s.dep.FetchJetWM(jet)
	if err := f.Procedure(ctx, jet, false); err != nil {
		return err
	}
	objJetID = jet.Result.Jet

	hot := proc.NewWaitHotWM(objJetID, flow.Pulse(ctx), s.message.WatermillMsg)
	s.dep.WaitHotWM(hot)
	if err := f.Procedure(ctx, hot, false); err != nil {
		return err
	}

	idx := proc.NewGetIndexWM(s.payload.ObjectID, objJetID, s.message.WatermillMsg)
	s.dep.GetIndexWM(idx)
	if err := f.Procedure(ctx, idx, false); err != nil {
		return err
	}
	stateID = *idx.Result.Index.LatestState

	jet = proc.NewFetchJetWM(stateID, stateID.Pulse(), s.message.WatermillMsg)
	s.dep.FetchJetWM(jet)
	if err := f.Procedure(ctx, jet, false); err != nil {
		return err
	}
	stateJetID = jet.Result.Jet

	send := proc.NewSendObject(s.message.WatermillMsg, s.payload, objJetID, stateJetID, idx.Result.Index)
	s.dep.SendObject(send)
	return f.Procedure(ctx, send, false)
}
