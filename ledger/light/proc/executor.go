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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/hot"
)

type FetchJet struct {
	target  insolar.ID
	pulse   insolar.PulseNumber
	message payload.Meta

	Result struct {
		Jet insolar.JetID
	}

	Dep struct {
		JetAccessor jet.Accessor
		Coordinator jet.Coordinator
		JetUpdater  executor.JetFetcher
		JetFetcher  executor.JetFetcher
		Sender      bus.Sender
	}
}

func NewFetchJet(target insolar.ID, pn insolar.PulseNumber, message payload.Meta) *FetchJet {
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
			msg, err := payload.NewMessage(&payload.Error{Text: err.Error()})
			if err != nil {
				return err
			}
			go p.Dep.Sender.Reply(ctx, p.message, msg)
		}
		return err
	}
	executor, err := p.Dep.Coordinator.LightExecutorForJet(ctx, *jetID, p.pulse)
	if err != nil {
		err := errors.Wrap(err, "failed to calculate executor for jet")
		if err != nil {
			msg, err := payload.NewMessage(&payload.Error{Text: err.Error()})
			if err != nil {
				return err
			}
			go p.Dep.Sender.Reply(ctx, p.message, msg)
		}
		return err
	}
	if *executor != p.Dep.Coordinator.Me() {
		msg := bus.ReplyAsMessage(ctx, &reply.JetMiss{JetID: *jetID, Pulse: p.pulse})
		go p.Dep.Sender.Reply(ctx, p.message, msg)
		return errors.New("jet miss")
	}

	p.Result.Jet = insolar.JetID(*jetID)
	return nil
}

type WaitHot struct {
	jetID   insolar.JetID
	pulse   insolar.PulseNumber
	message payload.Meta

	Dep struct {
		Waiter hot.JetWaiter
		Sender bus.Sender
	}
}

func NewWaitHot(j insolar.JetID, pn insolar.PulseNumber, message payload.Meta) *WaitHot {
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
		go p.Dep.Sender.Reply(ctx, p.message, msg)
		return err
	}

	return nil
}

type CheckJet struct {
	target  insolar.ID
	pulse   insolar.PulseNumber
	message payload.Meta
	pass    bool

	Result struct {
		Jet insolar.JetID
	}

	Dep struct {
		JetAccessor jet.Accessor
		Coordinator jet.Coordinator
		JetFetcher  executor.JetFetcher
		Sender      bus.Sender
	}
}

func NewCheckJet(target insolar.ID, pn insolar.PulseNumber, msg payload.Meta, pass bool) *CheckJet {
	return &CheckJet{
		target:  target,
		pulse:   pn,
		message: msg,
		pass:    pass,
	}
}

func (p *CheckJet) Proceed(ctx context.Context) error {
	jetID, err := p.Dep.JetFetcher.Fetch(ctx, p.target, p.pulse)
	if err != nil {
		return errors.Wrap(err, "failed to fetch jet")
	}

	executor, err := p.Dep.Coordinator.LightExecutorForJet(ctx, *jetID, p.pulse)
	if err != nil {
		return errors.Wrap(err, "failed to calculate executor for jet")
	}
	if *executor != p.Dep.Coordinator.Me() {
		if !p.pass {
			return ErrNotExecutor
		}

		buf, err := p.message.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal origin meta message")
		}
		msg, err := payload.NewMessage(&payload.Pass{
			Origin: buf,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}
		go func() {
			_, done := p.Dep.Sender.SendTarget(ctx, msg, *executor)
			done()
		}()
		return ErrNotExecutor
	}

	p.Result.Jet = insolar.JetID(*jetID)
	return nil
}

type WaitHotWM struct {
	jetID   insolar.JetID
	pulse   insolar.PulseNumber
	message payload.Meta

	Dep struct {
		Waiter hot.JetWaiter
		Sender bus.Sender
	}
}

func NewWaitHotWM(j insolar.JetID, pn insolar.PulseNumber, msg payload.Meta) *WaitHotWM {
	return &WaitHotWM{
		jetID:   j,
		pulse:   pn,
		message: msg,
	}
}

func (p *WaitHotWM) Proceed(ctx context.Context) error {
	return p.Dep.Waiter.Wait(ctx, insolar.ID(p.jetID), p.pulse)
}

type CalculateID struct {
	payload []byte
	pulse   insolar.PulseNumber

	Result struct {
		ID insolar.ID
	}

	dep struct {
		pcs insolar.PlatformCryptographyScheme
	}
}

func NewCalculateID(payload []byte, pulse insolar.PulseNumber) *CalculateID {
	return &CalculateID{
		payload: payload,
		pulse:   pulse,
	}
}

func (p *CalculateID) Dep(pcs insolar.PlatformCryptographyScheme) {
	p.dep.pcs = pcs
}

func (p *CalculateID) Proceed(ctx context.Context) error {
	h := p.dep.pcs.ReferenceHasher()
	_, err := h.Write(p.payload)
	if err != nil {
		return errors.Wrap(err, "failed to calculate id")
	}

	p.Result.ID = *insolar.NewID(p.pulse, h.Sum(nil))
	return nil
}
