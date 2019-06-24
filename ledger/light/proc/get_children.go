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
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetChildren struct {
	msg     *message.GetChildren
	parcel  insolar.Parcel
	replyTo chan<- bus.Reply

	Dep struct {
		IndexLocker            object.IndexLocker
		IndexAccessor          object.IndexAccessor
		Coordinator            jet.Coordinator
		RecordAccessor         object.RecordAccessor
		JetStorage             jet.Storage
		JetTreeUpdater         jet.Fetcher
		DelegationTokenFactory insolar.DelegationTokenFactory
	}
}

func NewGetChildren(msg *message.GetChildren, parcel insolar.Parcel, replyTo chan<- bus.Reply) *GetChildren {
	return &GetChildren{
		msg:     msg,
		parcel:  parcel,
		replyTo: replyTo,
	}
}

func (p *GetChildren) Proceed(ctx context.Context) error {
	p.replyTo <- p.reply(ctx)
	return nil
}

func (p *GetChildren) reply(ctx context.Context) bus.Reply {
	p.Dep.IndexLocker.Lock(p.msg.Parent.Record())
	defer p.Dep.IndexLocker.Unlock(p.msg.Parent.Record())

	idx, err := p.Dep.IndexAccessor.ForID(ctx, flow.Pulse(ctx), *p.msg.Parent.Record())
	if err != nil {
		return bus.Reply{Err: err}
	}

	// The object has no children.
	if idx.Lifeline.ChildPointer == nil {
		return bus.Reply{
			Reply: &reply.Children{Refs: nil, NextFrom: nil},
		}
	}

	var currentChild *insolar.ID

	// Counting from specified child or the latest.
	if p.msg.FromChild != nil {
		currentChild = p.msg.FromChild
	} else {
		currentChild = idx.Lifeline.ChildPointer
	}

	// The object has no children.
	if currentChild == nil {
		return bus.Reply{
			Reply: &reply.Children{Refs: nil, NextFrom: nil},
		}
	}

	var childJet *insolar.ID
	onHeavy, err := p.Dep.Coordinator.IsBeyondLimit(ctx, p.parcel.Pulse(), currentChild.Pulse())
	if err != nil && err != pulse.ErrNotFound {
		return bus.Reply{Err: err}
	}
	if onHeavy {
		node, err := p.Dep.Coordinator.Heavy(ctx, p.parcel.Pulse())
		if err != nil {
			return bus.Reply{Err: err}
		}
		repl, err := reply.NewGetChildrenRedirect(p.Dep.DelegationTokenFactory, p.parcel, node, *currentChild)
		if err != nil {
			return bus.Reply{Err: err}
		}
		return bus.Reply{Reply: repl}

	}

	childJetID, actual := p.Dep.JetStorage.ForID(ctx, currentChild.Pulse(), *p.msg.Parent.Record())
	childJet = (*insolar.ID)(&childJetID)

	if !actual {
		actualJet, err := p.Dep.JetTreeUpdater.Fetch(ctx, *p.msg.Parent.Record(), currentChild.Pulse())
		if err != nil {
			return bus.Reply{Err: err}
		}
		childJet = actualJet
	}

	// Try to fetch the first child.
	_, err = p.Dep.RecordAccessor.ForID(ctx, *currentChild)

	if err == object.ErrNotFound {
		node, err := p.Dep.Coordinator.NodeForJet(ctx, *childJet, p.parcel.Pulse(), currentChild.Pulse())
		if err != nil {
			return bus.Reply{Err: err}
		}
		repl, err := reply.NewGetChildrenRedirect(p.Dep.DelegationTokenFactory, p.parcel, node, *currentChild)
		if err != nil {
			return bus.Reply{Err: err}
		}
		return bus.Reply{Reply: repl}
	}

	if err != nil {
		return bus.Reply{Err: errors.Wrap(err, "failed to fetch child")}
	}

	var refs []insolar.Reference
	counter := 0
	for currentChild != nil {
		// We have enough results.
		if counter >= p.msg.Amount {
			return bus.Reply{Reply: &reply.Children{Refs: refs, NextFrom: currentChild}}
		}
		counter++

		rec, err := p.Dep.RecordAccessor.ForID(ctx, *currentChild)

		// We don't have this child reference. Return what was collected.
		if err == object.ErrNotFound {
			return bus.Reply{Reply: &reply.Children{Refs: refs, NextFrom: currentChild}}
		}
		if err != nil {
			return bus.Reply{Err: errors.New("failed to retrieve children")}
		}

		virtRec := rec.Virtual
		concrete := record.Unwrap(virtRec)
		childRec, ok := concrete.(*record.Child)
		if !ok {
			return bus.Reply{Err: errors.New("failed to retrieve children")}
		}
		currentChild = &childRec.PrevChild

		// Skip records later than specified pulse.
		recPulse := childRec.Ref.Record().Pulse()
		if p.msg.FromPulse != nil && recPulse > *p.msg.FromPulse {
			continue
		}
		refs = append(refs, childRec.Ref)
	}

	return bus.Reply{Reply: &reply.Children{Refs: refs, NextFrom: nil}}
}
