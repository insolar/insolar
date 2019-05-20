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
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetChildren struct {
	dep     *proc.Dependencies
	replyTo chan<- bus.Reply

	Message bus.Message
}

func (s *GetChildren) Present(ctx context.Context, f flow.Flow) error {
	msg := s.Message.Parcel.Message().(*message.GetChildren)

	var jetID insolar.JetID
	if s.Message.Parcel.DelegationToken() == nil {
		jet := proc.NewFetchJet(*msg.DefaultTarget().Record(), flow.Pulse(ctx), s.Message.ReplyTo)
		s.dep.FetchJet(jet)
		if err := f.Procedure(ctx, jet, false); err != nil {
			return err
		}
		hot := proc.NewWaitHot(jet.Result.Jet, flow.Pulse(ctx), s.Message.ReplyTo)
		s.dep.WaitHot(hot)
		if err := f.Procedure(ctx, hot, false); err != nil {
			return err
		}

		jetID = jet.Result.Jet
	} else {
		// Workaround to fetch object states.
		jet := proc.NewFetchJet(*msg.DefaultTarget().Record(), msg.FromChild.Pulse(), s.Message.ReplyTo)
		s.dep.FetchJet(jet)
		if err := f.Procedure(ctx, jet, false); err != nil {
			return err
		}
		jetID = jet.Result.Jet
	}

	getIndex := proc.NewGetIndex(msg.Parent, jetID, s.replyTo)
	s.dep.GetIndex(getIndex)
	if err := f.Procedure(ctx, getIndex, true); err != nil {
		return err
	}
	// The object has no children.
	if getIndex.Result.Index.ChildPointer == nil {
		s.replyTo <- bus.Reply{
			Reply: &reply.Children{Refs: nil, NextFrom: nil},
		}
		return nil
	}

	var currentChild *insolar.ID

	// Counting from specified child or the latest.
	if msg.FromChild != nil {
		currentChild = msg.FromChild
	} else {
		currentChild = getIndex.Result.Index.ChildPointer
	}

	// The object has no children.
	if currentChild == nil {
		s.replyTo <- bus.Reply{
			Reply: &reply.Children{Refs: nil, NextFrom: nil},
		}
		return nil
	}

	getChildren := proc.NewGetChildren(currentChild, msg, s.Message.Parcel, s.replyTo)
	s.dep.GetChildren(getChildren)
	if err := f.Procedure(ctx, getChildren, true); err != nil {
		return err
	}

	return nil
}
