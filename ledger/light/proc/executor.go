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
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/hot"
)

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

	worker, err := p.Dep.Coordinator.LightExecutorForJet(ctx, *jetID, p.pulse)
	if err != nil {
		return errors.Wrap(err, "failed to calculate executor for jet")
	}
	if *worker != p.Dep.Coordinator.Me() {
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
			_, done := p.Dep.Sender.SendTarget(ctx, msg, *worker)
			done()
		}()
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
		waiter hot.JetWaiter
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
	w hot.JetWaiter,
	s bus.Sender,
) {
	p.dep.waiter = w
	p.dep.sender = s
}

func (p *WaitHot) Proceed(ctx context.Context) error {
	return p.dep.waiter.Wait(ctx, insolar.ID(p.jetID), p.pulse)
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
