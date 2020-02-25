// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

import (
	"context"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/ledger/light/executor"
)

type FetchJet struct {
	target  insolar.ID
	pulse   insolar.PulseNumber
	message payload.Meta
	pass    bool

	Result struct {
		Jet insolar.JetID
	}

	dep struct {
		jetAccessor jet.Accessor
		jetFetcher  executor.JetFetcher
		coordinator jet.Coordinator
		sender      bus.Sender
	}
}

func NewFetchJet(target insolar.ID, pn insolar.PulseNumber, msg payload.Meta, pass bool) *FetchJet {
	return &FetchJet{
		target:  target,
		pulse:   pn,
		message: msg,
		pass:    pass,
	}
}

func (p *FetchJet) Dep(
	jets jet.Accessor,
	fetcher executor.JetFetcher,
	c jet.Coordinator,
	s bus.Sender,
) {
	p.dep.jetAccessor = jets
	p.dep.jetFetcher = fetcher
	p.dep.coordinator = c
	p.dep.sender = s
}

func (p *FetchJet) Proceed(ctx context.Context) error {
	jetID, err := p.dep.jetFetcher.Fetch(ctx, p.target, p.pulse)
	if err != nil {
		return errors.Wrap(err, "failed to fetch jet")
	}

	worker, err := p.dep.coordinator.LightExecutorForJet(ctx, *jetID, p.pulse)
	if err != nil {
		return errors.Wrap(err, "failed to calculate executor for jet")
	}
	if *worker != p.dep.coordinator.Me() {
		inslogger.FromContext(ctx).Warn("virtual node missed jet")
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
		_, done := p.dep.sender.SendTarget(ctx, msg, *worker)
		done()

		// Send calculated jet to virtual node.
		msg, err = payload.NewMessage(&payload.UpdateJet{
			Pulse: p.pulse,
			JetID: insolar.JetID(*jetID),
		})
		if err != nil {
			return errors.Wrap(err, "failed to create jet message")
		}
		_, done = p.dep.sender.SendTarget(ctx, msg, p.message.Sender)
		done()
		return ErrNotExecutor
	}

	p.Result.Jet = insolar.JetID(*jetID)
	return nil
}

type WaitHot struct {
	jetID   insolar.JetID
	pulse   insolar.PulseNumber
	message payload.Meta

	dep struct {
		waiter executor.JetWaiter
		sender bus.Sender
	}
}

func NewWaitHot(j insolar.JetID, pn insolar.PulseNumber, msg payload.Meta) *WaitHot {
	return &WaitHot{
		jetID:   j,
		pulse:   pn,
		message: msg,
	}
}

func (p *WaitHot) Dep(
	w executor.JetWaiter,
) {
	p.dep.waiter = w
}

func (p *WaitHot) Proceed(ctx context.Context) error {
	return p.dep.waiter.Wait(ctx, p.jetID, p.pulse)
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
