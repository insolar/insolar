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
	buswm "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
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
		BusWM            buswm.Sender
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
	logger := inslogger.FromContext(ctx)
	lfl, err := p.Dep.LifelineAccessor.ForID(ctx, p.pn, p.objID)
	if err != nil {
		return errors.Wrap(err, "[RefreshPendingFilament] can't fetch a lifeline state")
	}

	logger.Debugf("RefreshPendingFilament objID - %V     lfl.PendingPointer == %v || lfl.EarliestOpenRequest == %v", p.objID.DebugString(), lfl.PendingPointer, *lfl.EarliestOpenRequest)

	// No pendings
	if lfl.PendingPointer == nil {
		return nil
	}
	// No open pendings
	if lfl.EarliestOpenRequest == nil {
		return nil
	}
	// If an earliest pending created during current pulse
	if lfl.EarliestOpenRequest != nil && *lfl.EarliestOpenRequest == p.pn {
		return nil
	}

	fp, err := p.Dep.PendingAccessor.FirstPending(ctx, p.pn, p.objID)
	if err != nil {
		panic(err)
		return err
	}

	logger.Debugf("RefreshPendingFilament fp == %v, obj - %v", fp, p.objID.DebugString())

	if fp == nil || fp.PreviousRecord == nil {
		err = p.fillPendingFilament(ctx, p.pn, p.objID, lfl.PendingPointer.Pulse(), *lfl.EarliestOpenRequest)
		if err != nil {
			panic(err)
			return err
		}
	} else {
		err = p.fillPendingFilament(ctx, p.pn, p.objID, fp.PreviousRecord.Pulse(), *lfl.EarliestOpenRequest)
		if err != nil {
			panic(err)
			return err
		}
	}

	err = p.Dep.PendingModifier.RefreshState(ctx, p.pn, p.objID)
	if err != nil {
		panic(err)
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
			panic(err)
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
			panic(err)
			return err
		}

		msg, err := payload.NewMessage(&payload.GetPendingFilament{
			ObjectID: objID,
			Pulse:    destPN,
		})
		if err != nil {
			panic(err)
			return errors.Wrap(err, "failed to create a GetPendingFilament message")
		}

		rep, done := p.Dep.BusWM.SendTarget(ctx, msg, *node)
		defer done()

		res, ok := <-rep
		if !ok {
			panic(err)
			return errors.New("no reply")
		}

		pl, err := payload.UnmarshalFromMeta(res.Payload)
		if err != nil {
			panic(err)
			return errors.Wrap(err, "failed to unmarshal reply")
		}

		switch r := pl.(type) {
		case *payload.PendingFilament:
			err := p.Dep.PendingModifier.SetFilament(ctx, p.pn, objID, destPN, r.Records)
			if err != nil {
				panic(err)
				return err
			}

			if len(r.Records) == 0 {
				panic("unexpected behaviour")
			}

			if r.Records[0].Meta.PreviousRecord == nil {
				panic("we get here")
				continueFilling = false
			}

			// If know border read to the start of the chain
			// In other words, we read until limit
			if earlistOpenRequest == 0 || r.Records[0].Meta.PreviousRecord.Pulse() > earlistOpenRequest {
				destPN = r.Records[0].Meta.PreviousRecord.Pulse()
			} else {
				continueFilling = false
			}
		case *payload.Error:
			panic(err)
			return errors.New(r.Text)
		default:
			panic(err)
			return fmt.Errorf("RefreshPendingFilament: unexpected reply: %#v", p)
		}
	}

	return nil
}
