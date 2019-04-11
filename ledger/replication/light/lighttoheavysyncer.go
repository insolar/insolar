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

package light

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/utils/backoff"
)

type ToHeavySyncer interface {
	SyncPulse(pn insolar.PulseNumber) error
}

type toHeavySyncer struct {
	once sync.Once

	checker *time.Ticker

	waitingMux       sync.Mutex
	syncWaitingSlots []syncSlot

	sendingProblemMux   sync.Mutex
	sendingProblemSlots []syncSlot

	jetAccessor  jet.Accessor
	dataGatherer DataGatherer
	cleaner      Cleaner
	msgBus       insolar.MessageBus

	bconf configuration.Backoff
}

type syncSlot struct {
	pn    insolar.PulseNumber
	jetID insolar.JetID

	lastAttempt time.Time
	backoff     *backoff.Backoff
}

func (t *toHeavySyncer) addToWaitingSlots(pn insolar.PulseNumber) {
	t.waitingMux.Lock()
	defer t.waitingMux.Unlock()

	t.syncWaitingSlots = append(t.syncWaitingSlots, syncSlot{pn: pn})
}

func (t *toHeavySyncer) extractWaitingSlot() (syncSlot, bool) {
	t.waitingMux.Lock()
	defer t.waitingMux.Unlock()

	if len(t.syncWaitingSlots) == 0 {
		return syncSlot{}, false
	}

	slot := t.syncWaitingSlots[0]
	// it's a copy
	t.syncWaitingSlots = append([]syncSlot(nil), t.syncWaitingSlots[1:]...)
	return slot, true
}

func (t *toHeavySyncer) addToProblemsSlots(pn insolar.PulseNumber, jetID insolar.JetID) {
	t.sendingProblemMux.Lock()
	defer t.sendingProblemMux.Unlock()

	t.sendingProblemSlots = append(t.sendingProblemSlots,
		syncSlot{
			pn:          pn,
			jetID:       jetID,
			lastAttempt: time.Now(),
			backoff:     backoffFromConfig(t.bconf)},
	)
}

func (t *toHeavySyncer) reAddToProblemsSlots(ctx context.Context, slot syncSlot) {
	t.sendingProblemMux.Lock()
	defer t.sendingProblemMux.Unlock()

	if slot.backoff.Attempt() > t.bconf.MaxAttempts {
		inslogger.FromContext(ctx).Errorf("Failed to sync pulse - %v with jetID - %v. Attempts - %v", slot.pn, slot.jetID, slot.backoff.Attempt()-1)
		return
	}

	t.sendingProblemSlots = append(t.sendingProblemSlots, slot)
}

func (t *toHeavySyncer) extractProblemSlot() (syncSlot, bool) {
	t.sendingProblemMux.Lock()
	defer t.sendingProblemMux.Unlock()

	for i, slot := range t.sendingProblemSlots {
		pause := slot.backoff.ForAttempt(slot.backoff.Attempt())
		timeDiff := time.Now().Truncate(pause)
		if slot.lastAttempt.Before(timeDiff) {
			temp := append(t.sendingProblemSlots[:i], t.sendingProblemSlots[i+1:]...)
			// it's a copy
			t.sendingProblemSlots = append([]syncSlot(nil), temp...)
			return slot, true
		}
	}

	return syncSlot{}, false
}

func backoffFromConfig(bconf configuration.Backoff) *backoff.Backoff {
	return &backoff.Backoff{
		Jitter: bconf.Jitter,
		Min:    bconf.Min,
		Max:    bconf.Max,
		Factor: bconf.Factor,
	}
}

func (t *toHeavySyncer) SyncPulse(ctx context.Context, pn insolar.PulseNumber) {
	t.addToWaitingSlots(pn)
	t.lazyInit(ctx)
}

func (t *toHeavySyncer) lazyInit(ctx context.Context) {
	t.once.Do(func() {
		t.checker = time.NewTicker(500 * time.Millisecond)
		go func() {
			for range t.checker.C {
				go t.sync(ctx)
				go t.retrySync(ctx)
			}
		}()
	})
}

func (t *toHeavySyncer) sync(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	slot, ok := t.extractWaitingSlot()
	if !ok {
		logger.Infof("Sync queue is empty")
		return
	}

	jets := t.jetAccessor.All(ctx, slot.pn)
	for _, jID := range jets {
		msg, err := t.gatherForPnAndJet(ctx, slot.pn, jID)
		if err != nil {
			logger.Error(fmt.Sprintf("Problems with gather data for a pulse - %v and jet - %v", slot.pn, jID.DebugString()))
			continue
		}
		err = t.sendToHeavy(ctx, msg)
		if err != nil {
			logger.Errorf("Problems with sending msg to a heavy node", err)
			t.addToProblemsSlots(slot.pn, jID)
		}
	}
}

func (t *toHeavySyncer) retrySync(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	slot, ok := t.extractProblemSlot()
	if !ok {
		logger.Infof("Retry queue is empty")
		return
	}

	msg, err := t.dataGatherer.ForPulseAndJet(ctx, slot.pn, slot.jetID)
	if err != nil {
		logger.Error(fmt.Sprintf("Problems with gather data for a pulse - %v and jet - %v", slot.pn, slot.jetID.DebugString()))
		slot.backoff.Duration()
		t.reAddToProblemsSlots(ctx, slot)
		return
	}
	err = t.sendToHeavy(ctx, msg)
	if err != nil {
		logger.Errorf("Problems with sending msg to a heavy node", err)
		slot.backoff.Duration()
		t.reAddToProblemsSlots(ctx, slot)
		return
	}
}

func (t *toHeavySyncer) sendToHeavy(ctx context.Context, data *message.HeavyPayload) error {
	rep, err := t.msgBus.Send(ctx, data, nil)
	if err != nil {
		return err
	}
	if rep != nil {
		err, ok := rep.(*reply.HeavyError)
		if ok {
			return err
		}
	}
	return nil
}
