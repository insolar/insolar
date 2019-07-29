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

	"github.com/insolar/insolar/insolar/bus"
	wbus "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/object"
)

const (
	pendingNotifyThreshold = 2
)

type HotObjects struct {
	meta    payload.Meta
	jetID   insolar.JetID
	drop    drop.Drop
	indexes []record.Index
	pulse   insolar.PulseNumber

	dep struct {
		drops       drop.Modifier
		indices     object.IndexModifier
		jetStorage  jet.Storage
		jetFetcher  executor.JetFetcher
		jetReleaser hot.JetReleaser
		coordinator jet.Coordinator
		calculator  pulse.Calculator
		sender      wbus.Sender
	}
}

func NewHotObjects(
	meta payload.Meta,
	pn insolar.PulseNumber,
	jetID insolar.JetID,
	drop drop.Drop,
	indexes []record.Index,
) *HotObjects {
	return &HotObjects{
		meta:    meta,
		jetID:   jetID,
		drop:    drop,
		indexes: indexes,
		pulse:   pn,
	}
}

func (p *HotObjects) Dep(
	drops drop.Modifier,
	indices object.IndexModifier,
	jStore jet.Storage,
	jFetcher executor.JetFetcher,
	jReleaser hot.JetReleaser,
	coordinator jet.Coordinator,
	pCalc pulse.Calculator,
	sender bus.Sender,
) {
	p.dep.drops = drops
	p.dep.indices = indices
	p.dep.jetStorage = jStore
	p.dep.jetFetcher = jFetcher
	p.dep.jetReleaser = jReleaser
	p.dep.coordinator = coordinator
	p.dep.calculator = pCalc
	p.dep.sender = sender
}

func (p *HotObjects) Proceed(ctx context.Context) error {
	inslogger.FromContext(ctx).Debug("Get hot. pulse: ", p.drop.Pulse, " jet: ", p.drop.JetID.DebugString())

	err := p.dep.drops.Set(ctx, p.drop)
	if err == drop.ErrOverride {
		err = nil
	}
	if err != nil {
		return errors.Wrapf(err, "[HotObjects.process]: drop error (pulse: %v)", p.drop.Pulse)
	}

	err = p.dep.jetStorage.Update(
		ctx, p.pulse, true, p.jetID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update jet tree")
	}

	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"jet": p.jetID.DebugString(),
	})

	pendingNotifyPulse, err := p.dep.calculator.Backwards(ctx, flow.Pulse(ctx), pendingNotifyThreshold)
	if err != nil {
		if err == pulse.ErrNotFound {
			pendingNotifyPulse = *insolar.GenesisPulse
		} else {
			return errors.Wrap(err, "failed to calculate pending notify pulse")
		}
	}

	logger.Debugf("received %v hot indexes for jet %s and pulse %s", len(p.indexes), p.jetID.DebugString(), p.pulse)
	for _, idx := range p.indexes {
		objJetID, _ := p.dep.jetStorage.ForID(ctx, p.pulse, idx.ObjID)
		if objJetID != p.jetID {
			logger.Warn("received wrong id")
			continue
		}

		err = p.dep.indices.SetIndex(
			ctx,
			p.pulse,
			record.Index{
				ObjID:            idx.ObjID,
				Lifeline:         idx.Lifeline,
				LifelineLastUsed: idx.LifelineLastUsed,
				PendingRecords:   []insolar.ID{},
			},
		)
		if err != nil {
			logger.Error(errors.Wrapf(err, "[handleHotRecords] failed to save index - %v", idx.ObjID.DebugString()))
			continue
		}
		logger.Debugf("[handleHotRecords] lifeline with id - %v saved", idx.ObjID.DebugString())

		go p.notifyPending(ctx, idx.ObjID, idx.Lifeline, pendingNotifyPulse.PulseNumber)
	}

	p.dep.jetFetcher.Release(ctx, p.jetID, p.pulse)
	err = p.dep.jetReleaser.Unlock(ctx, insolar.ID(p.jetID))
	if err != nil {
		return errors.Wrap(err, "failed to release jets")
	}

	p.sendConfirmationToHeavy(ctx, p.jetID, p.drop.Pulse, p.drop.Split)
	return nil
}

func (p *HotObjects) sendConfirmationToHeavy(ctx context.Context, jetID insolar.JetID, pn insolar.PulseNumber, split bool) {
	msg, err := payload.NewMessage(&payload.GotHotConfirmation{
		JetID: jetID,
		Pulse: pn,
		Split: split,
	})

	if err != nil {
		inslogger.FromContext(ctx).Error("Can't create GotHotConfirmation message: ", err)
		return
	}

	inslogger.FromContext(ctx).Debug("Send hot confirmation to heavy. pulse: ", pn, " jet: ", p.drop.JetID.DebugString())

	_, done := p.dep.sender.SendRole(ctx, msg, insolar.DynamicRoleHeavyExecutor, *insolar.NewReference(insolar.ID(jetID)))
	done()
}

func (p *HotObjects) notifyPending(
	ctx context.Context,
	objectID insolar.ID,
	lifeline record.Lifeline,
	notifyLimit insolar.PulseNumber,
) {
	// No pending requests.
	if lifeline.EarliestOpenRequest == nil {
		return
	}

	// Too early to notify.
	if *lifeline.EarliestOpenRequest >= notifyLimit {
		return
	}

	msg, err := payload.NewMessage(&payload.AbandonedRequestsNotification{
		ObjectID: objectID,
	})
	if err != nil {
		inslogger.FromContext(ctx).Error("failed to notify about pending requests: ", err.Error())
		return
	}

	// Hot data was sent to us by another light.
	// Notification should be send to virtual for this object.
	_, done := p.dep.sender.SendRole(ctx, msg, insolar.DynamicRoleVirtualExecutor, *insolar.NewReference(objectID))
	done()
}
