// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	wbus "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
)

const (
	pendingNotifyThreshold = 2
)

type HotObjects struct {
	meta                   payload.Meta
	jetID                  insolar.JetID
	drop                   drop.Drop
	indexes                []record.Index
	pulse                  insolar.PulseNumber
	availableNotifications uint

	dep struct {
		drops       drop.Modifier
		indexes     object.MemoryIndexModifier
		jetStorage  jet.Storage
		jetFetcher  executor.JetFetcher
		jetReleaser executor.JetReleaser
		coordinator jet.Coordinator
		calculator  pulse.Calculator
		sender      wbus.Sender
		registry    executor.MetricsRegistry
	}
}

func NewHotObjects(
	meta payload.Meta,
	pn insolar.PulseNumber,
	jetID insolar.JetID,
	drop drop.Drop,
	indexes []record.Index,
	availableNotifications uint,
) *HotObjects {
	return &HotObjects{
		meta:                   meta,
		jetID:                  jetID,
		drop:                   drop,
		indexes:                indexes,
		pulse:                  pn,
		availableNotifications: availableNotifications,
	}
}

func (p *HotObjects) Dep(
	drops drop.Modifier,
	indexes object.MemoryIndexModifier,
	jStore jet.Storage,
	jFetcher executor.JetFetcher,
	jReleaser executor.JetReleaser,
	coordinator jet.Coordinator,
	pCalc pulse.Calculator,
	sender bus.Sender,
	registry executor.MetricsRegistry,
) {
	p.dep.drops = drops
	p.dep.indexes = indexes
	p.dep.jetStorage = jStore
	p.dep.jetFetcher = jFetcher
	p.dep.jetReleaser = jReleaser
	p.dep.coordinator = coordinator
	p.dep.calculator = pCalc
	p.dep.sender = sender
	p.dep.registry = registry
}

func (p *HotObjects) Proceed(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("start hotObjects processing")

	logger.Debug("get hots. pulse: ", p.drop.Pulse, " jet: ", p.drop.JetID.DebugString())

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

		p.dep.indexes.Set(
			ctx,
			p.pulse,
			record.Index{
				ObjID:            idx.ObjID,
				Lifeline:         idx.Lifeline,
				LifelineLastUsed: idx.LifelineLastUsed,
				PendingRecords:   []insolar.ID{},
			},
		)
		logger.Debugf("[handleHotRecords] lifeline with id - %v saved", idx.ObjID.DebugString())

		p.notifyPending(ctx, idx.ObjID, idx.Lifeline, pendingNotifyPulse.PulseNumber)
	}

	logger.Infof("before releasing jetFetcher for jet %s and pulse %s", p.jetID.DebugString(), p.pulse)
	p.dep.jetFetcher.Release(ctx, p.jetID, p.pulse)

	logger.Infof("before unlocking jetReleaser for jet %s and pulse %s", p.jetID.DebugString(), p.pulse)
	err = p.dep.jetReleaser.Unlock(ctx, p.pulse, p.jetID)
	if err != nil {
		return errors.Wrap(err, "failed to release jets")
	}

	logger.Infof("jetFetcher and jetReleaser for jet %s and pulse %s are released", p.jetID.DebugString(), p.pulse)

	p.sendConfirmationToHeavy(ctx, p.jetID, p.drop.Pulse, p.drop.Split)

	logger.Infof("finish hotObjects processing")
	return nil
}

func (p *HotObjects) sendConfirmationToHeavy(ctx context.Context, jetID insolar.JetID, pn insolar.PulseNumber, split bool) {
	logger := inslogger.FromContext(ctx)
	msg, err := payload.NewMessage(&payload.GotHotConfirmation{
		JetID: jetID,
		Pulse: pn,
		Split: split,
	})

	if err != nil {
		logger.Error("can't create GotHotConfirmation message: ", err)
		return
	}

	logger.Info("send hot confirmation to heavy. pulse: ", pn, " jet: ", p.drop.JetID.DebugString())

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

	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"object":        objectID.String(),
		"earliestPulse": *lifeline.EarliestOpenRequest,
	})

	if p.availableNotifications <= 0 {
		logger.Warn("out of AbandonedRequestsNotification limit")
		return
	}
	p.availableNotifications--

	// This is needed for metrics: we collect the earliest requests pulse we got
	p.dep.registry.SetOldestAbandonedRequestAge(int(p.pulse - *lifeline.EarliestOpenRequest))
	logger.Debug("sending AbandonedRequestNotification")

	msg, err := payload.NewMessage(&payload.AbandonedRequestsNotification{
		ObjectID: objectID,
	})
	if err != nil {
		logger.Error("failed to notify about pending requests: ", err.Error())
		return
	}

	// Hot data was sent to us by another light.
	// Notification should be send to virtual for this object.
	_, done := p.dep.sender.SendRole(ctx, msg, insolar.DynamicRoleVirtualExecutor, *insolar.NewReference(objectID))
	done()
	stats.Record(ctx, statHotsAbandoned.M(1))
}
