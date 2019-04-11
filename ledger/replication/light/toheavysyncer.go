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
	SyncPulse(ctx context.Context, pn insolar.PulseNumber)
}

type toHeavySyncer struct {
	once sync.Once

	checker *time.Ticker

	waitingMux        sync.Mutex
	syncWaitingPulses []insolar.PulseNumber

	notSentPayloadsMux sync.Mutex
	notSentPayloads    []notSentPayload

	jetAccessor  jet.Accessor
	dataGatherer DataGatherer
	cleaner      Cleaner
	msgBus       insolar.MessageBus

	bconf configuration.Backoff
}

func NewToHeavySyncer(
	jetAccessor jet.Accessor,
	dataGatherer DataGatherer,
	cleaner Cleaner,
	msgBus insolar.MessageBus,
	bconf configuration.Backoff,
) ToHeavySyncer {
	return &toHeavySyncer{
		jetAccessor:  jetAccessor,
		dataGatherer: dataGatherer,
		cleaner:      cleaner,
		msgBus:       msgBus,
		bconf:        bconf,
	}
}

type notSentPayload struct {
	msg *message.HeavyPayload

	lastAttempt time.Time
	backoff     *backoff.Backoff
}

func (t *toHeavySyncer) addToWaitingPulses(pn insolar.PulseNumber) {
	t.waitingMux.Lock()
	defer t.waitingMux.Unlock()

	t.syncWaitingPulses = append(t.syncWaitingPulses, pn)
}

func (t *toHeavySyncer) extractWaitingPulse() (insolar.PulseNumber, bool) {
	t.waitingMux.Lock()
	defer t.waitingMux.Unlock()

	if len(t.syncWaitingPulses) == 0 {
		return insolar.FirstPulseNumber, false
	}

	slot := t.syncWaitingPulses[0]
	// it's a copy
	t.syncWaitingPulses = append([]insolar.PulseNumber(nil), t.syncWaitingPulses[1:]...)
	return slot, true
}

func (t *toHeavySyncer) addToNotSentPayloads(payload *message.HeavyPayload) {
	t.notSentPayloadsMux.Lock()
	defer t.notSentPayloadsMux.Unlock()

	t.notSentPayloads = append(t.notSentPayloads,
		notSentPayload{
			msg:         payload,
			backoff:     backoffFromConfig(t.bconf),
			lastAttempt: time.Now(),
		},
	)
}

func (t *toHeavySyncer) reAddToNotSentPayloads(ctx context.Context, payload notSentPayload) {
	t.notSentPayloadsMux.Lock()
	defer t.notSentPayloadsMux.Unlock()

	if payload.backoff.Attempt() > t.bconf.MaxAttempts {
		inslogger.FromContext(ctx).Errorf("Failed to sync pulse - %v with jetID - %v. Attempts - %v", payload.msg.PulseNum, payload.msg.JetID, payload.backoff.Attempt()-1)
		return
	}

	t.notSentPayloads = append(t.notSentPayloads, payload)
}

func (t *toHeavySyncer) extractNotSentPayload() (notSentPayload, bool) {
	t.notSentPayloadsMux.Lock()
	defer t.notSentPayloadsMux.Unlock()

	for i, notSent := range t.notSentPayloads {
		pause := notSent.backoff.ForAttempt(notSent.backoff.Attempt())
		timeDiff := time.Now().Truncate(pause)
		if notSent.lastAttempt.Before(timeDiff) {
			temp := append(t.notSentPayloads[:i], t.notSentPayloads[i+1:]...)
			// it's a copy
			t.notSentPayloads = append([]notSentPayload(nil), temp...)
			return notSent, true
		}
	}

	return notSentPayload{}, false
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
	t.addToWaitingPulses(pn)
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
	pn, ok := t.extractWaitingPulse()
	if !ok {
		logger.Infof("Sync queue is empty")
		return
	}

	jets := t.jetAccessor.All(ctx, pn)
	for _, jID := range jets {
		msg, err := t.dataGatherer.ForPulseAndJet(ctx, pn, jID)
		if err != nil {
			panic(fmt.Sprintf("Problems with gather data for a pulse - %v and jet - %v", pn, jID.DebugString()))
		}
		err = t.sendToHeavy(ctx, msg)
		if err != nil {
			logger.Errorf("Problems with sending msg to a heavy node", err)
			t.addToNotSentPayloads(msg)
			continue
		}
	}

	t.cleaner.Clean(ctx, pn)
}

func (t *toHeavySyncer) retrySync(ctx context.Context) {
	logger := inslogger.FromContext(ctx)

	for payload, ok := t.extractNotSentPayload(); ok; payload, ok = t.extractNotSentPayload() {
		err := t.sendToHeavy(ctx, payload.msg)
		if err != nil {
			logger.Errorf("Problems with sending msg to a heavy node", err)
			payload.backoff.Duration()
			t.reAddToNotSentPayloads(ctx, payload)
		}
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
