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

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/pkg/errors"
)

type FetchJet struct {
	target  insolar.ID
	pulse   insolar.PulseNumber
	message *watermillMsg.Message

	Result struct {
		Jet insolar.JetID
	}

	Dep struct {
		JetAccessor jet.Accessor
		Coordinator jet.Coordinator
		JetUpdater  jet.Fetcher
		JetFetcher  jet.Fetcher
		Sender      bus.Sender
	}
}

func NewFetchJet(target insolar.ID, pn insolar.PulseNumber, message *watermillMsg.Message) *FetchJet {
	return &FetchJet{
		target:  target,
		pulse:   pn,
		message: message,
	}
}

func (p *FetchJet) Proceed(ctx context.Context) error {
	// Special case for genesis pulse. No one was executor at that time, so anyone can fetch data from it.
	if p.pulse <= insolar.FirstPulseNumber {
		p.Result.Jet = *insolar.NewJetID(0, nil)
		return nil
	}

	jetID, err := p.Dep.JetFetcher.Fetch(ctx, p.target, p.pulse)
	if err != nil {
		err := errors.Wrap(err, "failed to fetch jet")
		if err != nil {
			msg := bus.ErrorAsMessage(ctx, err)
			p.Dep.Sender.Reply(ctx, p.message, msg)
		}
		return err
	}
	executor, err := p.Dep.Coordinator.LightExecutorForJet(ctx, *jetID, p.pulse)
	if err != nil {
		err := errors.Wrap(err, "failed to calculate executor for jet")
		if err != nil {
			msg := bus.ErrorAsMessage(ctx, err)
			p.Dep.Sender.Reply(ctx, p.message, msg)
		}
		return err
	}
	if *executor != p.Dep.Coordinator.Me() {
		msg := bus.ReplyAsMessage(ctx, &reply.JetMiss{JetID: *jetID, Pulse: p.pulse})
		p.Dep.Sender.Reply(ctx, p.message, msg)
		return errors.New("jet miss")
	}

	p.Result.Jet = insolar.JetID(*jetID)
	return nil
}

type WaitHot struct {
	jetID   insolar.JetID
	pulse   insolar.PulseNumber
	message *watermillMsg.Message

	Dep struct {
		Waiter hot.JetWaiter
		Sender bus.Sender
	}
}

func NewWaitHot(j insolar.JetID, pn insolar.PulseNumber, message *watermillMsg.Message) *WaitHot {
	return &WaitHot{
		jetID:   j,
		pulse:   pn,
		message: message,
	}
}

func (p *WaitHot) Proceed(ctx context.Context) error {
	err := p.Dep.Waiter.Wait(ctx, insolar.ID(p.jetID), p.pulse)
	if err != nil {
		msg := bus.ReplyAsMessage(ctx, &reply.Error{ErrType: reply.ErrHotDataTimeout})
		p.Dep.Sender.Reply(ctx, p.message, msg)
		return err
	}

	return nil
}
