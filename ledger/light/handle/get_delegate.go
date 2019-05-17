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
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetDelegate struct {
	dep     *proc.Dependencies
	replyTo chan<- bus.Reply
	parcel  insolar.Parcel
}

func NewGetDelegate(dep *proc.Dependencies, rep chan<- bus.Reply, parcel insolar.Parcel) *GetDelegate {
	return &GetDelegate{
		dep:     dep,
		parcel:  parcel,
		replyTo: rep,
	}
}

func (s *GetDelegate) Present(ctx context.Context, f flow.Flow) error {
	msg := s.parcel.Message().(*message.GetDelegate)

	jet := proc.NewFetchJet(*msg.DefaultTarget().Record(), flow.Pulse(ctx), s.replyTo)
	s.dep.FetchJet(jet)
	if err := f.Procedure(ctx, jet, false); err != nil {
		return err
	}

	idx := proc.NewGetIndex(msg.Head, jet.Result.Jet, s.replyTo, s.parcel.Pulse())
	s.dep.GetIndex(idx)
	if err := f.Procedure(ctx, idx, false); err != nil {
		return err
	}

	getDelegate := proc.NewGetDelegate(msg, &idx.Result.Index, s.replyTo)
	if err := f.Procedure(ctx, getDelegate, false); err != nil {
		return err
	}
	return nil
}
