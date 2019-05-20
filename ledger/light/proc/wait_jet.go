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

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/pkg/errors"
)

type FetchJet struct {
	message *message.Message
	target  insolar.ID
	pulse   insolar.PulseNumber

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

func NewFetchJet(target insolar.ID, pn insolar.PulseNumber, msg *message.Message) *FetchJet {
	return &FetchJet{
		target:  target,
		pulse:   pn,
		message: msg,
	}
}

func (p *FetchJet) Proceed(ctx context.Context) error {
	err := p.proceed(ctx)
	if err != nil {
		go func() {
			pl := payload.Error{
				Text: err.Error(),
			}
			buf, err := pl.Marshal()
			if err != nil {
				inslogger.FromContext(ctx).Error("failed to encode payload")
			}
			p.Dep.Sender.Reply(ctx, p.message, message.NewMessage(watermill.NewUUID(), buf))
		}()
	}
	return err
}

func (p *FetchJet) proceed(ctx context.Context) error {
	// Special case for genesis pulse. No one was executor at that time, so anyone can fetch data from it.
	if p.pulse <= insolar.FirstPulseNumber {
		p.Result.Jet = *insolar.NewJetID(0, nil)
		return nil
	}

	jetID, err := p.Dep.JetFetcher.Fetch(ctx, p.target, p.pulse)
	if err != nil {
		err := errors.Wrap(err, "failed to fetch jet")
		return err
	}
	executor, err := p.Dep.Coordinator.LightExecutorForJet(ctx, *jetID, p.pulse)
	if err != nil {
		err := errors.Wrap(err, "failed to calculate executor for jet")
		return err
	}
	if *executor != p.Dep.Coordinator.Me() {
		return errors.New("jet miss")
	}

	p.Result.Jet = insolar.JetID(*jetID)
	return nil
}

type WaitHot struct {
	jetID   insolar.JetID
	pulse   insolar.PulseNumber
	message *message.Message

	Dep struct {
		Waiter hot.JetWaiter
		Sender bus.Sender
	}
}

func NewWaitHot(j insolar.JetID, pn insolar.PulseNumber, msg *message.Message) *WaitHot {
	return &WaitHot{
		jetID:   j,
		pulse:   pn,
		message: msg,
	}
}

func (p *WaitHot) Proceed(ctx context.Context) error {
	err := p.Dep.Waiter.Wait(ctx, insolar.ID(p.jetID))
	if err != nil {
		go func() {
			pl := payload.Error{
				Text: err.Error(),
			}
			buf, err := pl.Marshal()
			if err != nil {
				inslogger.FromContext(ctx).Error("failed to encode payload")
			}
			p.Dep.Sender.Reply(ctx, p.message, message.NewMessage(watermill.NewUUID(), buf))
		}()
	}

	return err
}
