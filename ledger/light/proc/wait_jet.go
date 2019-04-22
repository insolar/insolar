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
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/pkg/errors"
)

type FetchJet struct {
	Parcel insolar.Parcel

	Result struct {
		Jet   insolar.JetID
		Miss  bool
		Pulse insolar.PulseNumber
	}

	Dep struct {
		JetAccessor jet.Accessor
		Coordinator jet.Coordinator
		JetUpdater  jet.Fetcher
	}
}

func (p *FetchJet) Proceed(ctx context.Context) error {
	msg := p.Parcel.Message()
	if msg.DefaultTarget() == nil {
		return errors.New("unexpected message")
	}

	// Hack to temporary allow any genesis request.
	if p.Parcel.Pulse() <= insolar.FirstPulseNumber {
		p.Result.Jet = *insolar.NewJetID(0, nil)
		return nil
	}

	// Check token jet.
	token := p.Parcel.DelegationToken()
	if token != nil {
		// Calculate jet for target pulse.
		target := *msg.DefaultTarget().Record()
		pulse := target.Pulse()
		switch tm := msg.(type) {
		case *message.GetObject:
			pulse = tm.State.Pulse()
		case *message.GetChildren:
			if tm.FromChild == nil {
				return errors.New("fetching children without child pointer is forbidden")
			}
			pulse = tm.FromChild.Pulse()
		case *message.GetRequest:
			pulse = tm.Request.Pulse()
		}
		jetID, actual := p.Dep.JetAccessor.ForID(ctx, pulse, target)
		if !actual {
			inslogger.FromContext(ctx).WithFields(map[string]interface{}{
				"msg":   msg.Type().String(),
				"jet":   jetID.DebugString(),
				"pulse": pulse,
			}).Error("jet is not actual")
		}

		p.Result.Jet = jetID
		return nil
	}

	// Calculate jet for current pulse.
	// Calculate jet and pulse.
	var jetID insolar.ID
	var pulse insolar.PulseNumber
	if msg.DefaultTarget().Record().Pulse() == insolar.PulseNumberJet {
		jetID = *msg.DefaultTarget().Record()
	} else {
		if gr, ok := msg.(*message.GetRequest); ok {
			pulse = gr.Request.Pulse()
		} else {
			pulse = p.Parcel.Pulse()
		}

		j, err := p.Dep.JetUpdater.Fetch(ctx, *msg.DefaultTarget().Record(), p.Parcel.Pulse())
		if err != nil {
			return errors.Wrap(err, "failed to fetch jet tree")
		}

		jetID = *j
	}

	// Check if jet is ours.
	node, err := p.Dep.Coordinator.LightExecutorForJet(ctx, jetID, p.Parcel.Pulse())
	if err != nil {
		return errors.Wrap(err, "failed to calculate executor for jet")
	}

	if *node != p.Dep.Coordinator.Me() {
		p.Result.Miss = true
		p.Result.Pulse = pulse
		p.Result.Jet = insolar.JetID(jetID)
		return nil
	}

	p.Result.Jet = insolar.JetID(jetID)
	return nil
}

type WaitHot struct {
	Parcel insolar.Parcel
	JetID  insolar.JetID

	Res struct {
		Timeout bool
	}

	Dep struct {
		Waiter hot.JetWaiter
	}
}

func (p *WaitHot) Proceed(ctx context.Context) error {
	parcel := p.Parcel
	// Hack is needed for genesis:
	// because we don't have hot data on first pulse and without this we would stale.
	if parcel.Pulse() <= insolar.FirstPulseNumber {
		return nil
	}

	// If the call is a call in redirect-chain
	// skip waiting for the hot records
	if parcel.DelegationToken() != nil {
		return nil
	}

	err := p.Dep.Waiter.Wait(ctx, insolar.ID(p.JetID))
	if err != nil {
		p.Res.Timeout = true
	}

	return nil
}
