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
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type RefreshPendingFilament struct {
	replyTo chan<- bus.Reply
	objID   insolar.ID
	pn      insolar.PulseNumber

	Dep struct {
		PendingAccessor  object.PendingAccessor
		PendingModifier  object.PendingModifier
		LifelineAccessor object.LifelineAccessor
		Coordinator      jet.Coordinator
		Bus              insolar.MessageBus
	}
}

func NewRefreshPendingFilament(replyTo chan<- bus.Reply, pn insolar.PulseNumber, objID insolar.ID) *RefreshPendingFilament {
	return &RefreshPendingFilament{replyTo: replyTo, objID: objID, pn: pn}
}

func (p *RefreshPendingFilament) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, fmt.Sprintf("RefreshPendingFilament"))
	defer span.End()

	err := p.process(ctx)
	if err != nil {
		p.replyTo <- bus.Reply{Err: err}
	}
	return err
}

func (p *RefreshPendingFilament) process(ctx context.Context) error {
	lfl, err := p.Dep.LifelineAccessor.ForID(ctx, p.pn, p.objID)
	if err != nil {
		return errors.Wrap(err, "[RefreshPendingFilament] can't fetch a lifeline state")
	}

	if lfl.PendingPointer == nil || lfl.EarliestOpenRequest == nil {
		return nil
	}

	fp, err := p.Dep.PendingAccessor.FirstPending(ctx, p.pn, p.objID)
	if err != nil {
		return err
	}

	if fp == nil || fp.PreviousRecord == nil {
		err = p.fillPendingFilament(ctx, p.pn, p.objID, lfl.PendingPointer.Pulse(), *lfl.EarliestOpenRequest)
		if err != nil {
			return err
		}
	} else {
		err = p.fillPendingFilament(ctx, p.pn, p.objID, fp.PreviousRecord.Pulse(), *lfl.EarliestOpenRequest)
		if err != nil {
			return err
		}
	}

	err = p.Dep.PendingModifier.RefreshState(ctx, p.pn, p.objID)
	if err != nil {
		return err
	}

	return nil
}

func (p *RefreshPendingFilament) fillPendingFilament(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID, destPN insolar.PulseNumber, earlistOpenRequest insolar.PulseNumber) error {
	ctx, span := instracer.StartSpan(ctx, fmt.Sprintf("RefreshPendingFilament.fillPendingFilament"))
	defer span.End()

	continueFilling := true

	for continueFilling {
		isBeyond, err := p.Dep.Coordinator.IsBeyondLimit(ctx, currentPN, destPN)
		if err != nil {
			return err
		}
		if isBeyond {
			panic("we don't want to be here")
			// We need to update our chain
			// If oldest at the heavy
			return nil
		}

		node, err := p.Dep.Coordinator.NodeForObject(ctx, objID, currentPN, destPN)
		if err != nil {
			return err
		}

		rep, err := p.Dep.Bus.Send(
			ctx,
			&message.GetPendingFilament{ObjectID: objID},
			&insolar.MessageSendOptions{
				Receiver: node,
			},
		)
		if err != nil {
			return err
		}

		switch r := rep.(type) {
		case *reply.PendingFilament:
			err := p.Dep.PendingModifier.SetFilament(ctx, p.pn, objID, destPN, r.Records)
			if err != nil {
				return err
			}

			if len(r.Records) == 0 {
				panic("unexpected behaviour")
			}

			if r.Records[0].Meta.PreviousRecord.Pulse() == 0 {
				continueFilling = false
			}

			// If know border read to the start of the chain
			// In other words, we read until limit
			if earlistOpenRequest == 0 || r.Records[0].Meta.PreviousRecord.Pulse() > earlistOpenRequest {
				destPN = r.Records[0].Meta.PreviousRecord.Pulse()
			} else {
				continueFilling = false
			}
		case *reply.Error:
			return r.Error()
		default:
			return fmt.Errorf("fillPendingFilament: unexpected reply: %#v", rep)
		}
	}

	return nil
}
