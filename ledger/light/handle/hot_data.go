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

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/ledger/light/proc"
)

type HotData struct {
	dep     *proc.Dependencies
	replyTo chan<- bus.Reply
	message *message.HotData
	//pulse   insolar.PulseNumber
}

func NewHotData(dep *proc.Dependencies, rep chan<- bus.Reply, msg *message.HotData /*, pulse insolar.PulseNumber*/) *HotData {
	return &HotData{
		dep:     dep,
		replyTo: rep,
		message: msg,
		// pulse:   pulse,
	}
}

func (s *HotData) Present(ctx context.Context, f flow.Flow) error {
	// TODO implement
	return nil
	/***
	jet := proc.NewFetchJet(*s.message.DefaultTarget().Record(), flow.Pulse(ctx), s.replyTo)
	s.dep.FetchJet(jet)
	if err := f.Procedure(ctx, jet, false); err != nil {
		return err
	}

	hot := proc.NewWaitHot(jet.Result.Jet, flow.Pulse(ctx), s.replyTo)
	s.dep.WaitHot(hot)
	if err := f.Procedure(ctx, hot, false); err != nil {

		return err
	}

	getIndex := proc.NewGetIndex(s.message.Parent, jet.Result.Jet, s.replyTo)
	s.dep.GetIndex(getIndex)
	err := f.Procedure(ctx, getIndex, false)
	if err != nil {
		return err
	}

	registerChild := proc.NewRegisterChild(jet.Result.Jet, s.message, s.pulse, getIndex.Result.Index, s.replyTo)
	s.dep.RegisterChild(registerChild)
	return f.Procedure(ctx, registerChild, false)
	 ***/
}
