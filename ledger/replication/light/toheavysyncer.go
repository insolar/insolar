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
	"github.com/insolar/insolar/ledger/storage/pulse"
	"github.com/insolar/insolar/utils/backoff"
	"go.opencensus.io/stats"
)

type ToHeavySyncer interface {
	NotifyAboutPulse(ctx context.Context, pn insolar.PulseNumber)
}

type toHeavySyncer struct {
	once sync.Once

	checker *time.Ticker

	syncWaitingPulses chan insolar.PulseNumber

	notSentPayloadsMux sync.Mutex
	notSentPayloads    []notSentPayload

	jetCalculator   jet.Calculator
	dataGatherer    DataGatherer
	cleaner         Cleaner
	msgBus          insolar.MessageBus
	pulseCalculator pulse.Calculator

	conf configuration.LightToHeavySync
}

func NewToHeavySyncer(
	jetCalculator jet.Calculator,
	dataGatherer DataGatherer,
	cleaner Cleaner,
	msgBus insolar.MessageBus,
	conf configuration.LightToHeavySync,
	calculator pulse.Calculator,
) ToHeavySyncer {
	return &toHeavySyncer{
		jetCalculator:     jetCalculator,
		dataGatherer:      dataGatherer,
		cleaner:           cleaner,
		msgBus:            msgBus,
		pulseCalculator:   calculator,
		conf:              conf,
		syncWaitingPulses: make(chan insolar.PulseNumber),
	}
}

type notSentPayload struct {
	msg *message.HeavyPayload

	lastAttempt time.Time
	backoff     *backoff.Backoff
}

func (t *toHeavySyncer) addToNotSentPayloads(payload *message.HeavyPayload) {
	t.notSentPayloadsMux.Lock()
	defer t.notSentPayloadsMux.Unlock()

	t.notSentPayloads = append(t.notSentPayloads,
		notSentPayload{
			msg:         payload,
			backoff:     backoffFromConfig(t.conf.Backoff),
			lastAttempt: time.Now(),
		},
	)
}

func (t *toHeavySyncer) reAddToNotSentPayloads(ctx context.Context, payload notSentPayload) {
	t.notSentPayloadsMux.Lock()
	defer t.notSentPayloadsMux.Unlock()

	if payload.backoff.Attempt() > t.conf.Backoff.MaxAttempts {
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

func (t *toHeavySyncer) NotifyAboutPulse(ctx context.Context, pn insolar.PulseNumber) {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("[NotifyAboutPulse] pn - %v", pn)

	t.lazyInit(ctx)
	prevPN, err := t.pulseCalculator.Backwards(ctx, pn, 1)
	if err != nil {
		logger.Error("[NotifyAboutPulse]", err)
		return
	}

	logger.Debugf("[NotifyAboutPulse] prevPn - %v", prevPN.PulseNumber)
	t.syncWaitingPulses <- prevPN.PulseNumber
}

func (t *toHeavySyncer) lazyInit(ctx context.Context) {
	t.once.Do(func() {
		go t.sync(ctx)
		inslogger.FromContext(ctx).Debugf("[lazyInit] start rechecker with duration - %v", t.conf.RetryLoopDuration)
		t.checker = time.NewTicker(t.conf.RetryLoopDuration)
		go func() {
			for range t.checker.C {
				go t.retrySync(ctx)
			}
		}()
	})
}

func (t *toHeavySyncer) sync(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	for pn := range t.syncWaitingPulses {
		logger.Debugf("[sync] pn received - %v", pn)

		jets := t.jetCalculator.MineForPulse(ctx, pn)
		logger.Debugf("[sync] founds %v jets", len(jets))
		for _, jID := range jets {
			msg, err := t.dataGatherer.ForPulseAndJet(ctx, pn, jID)
			if err != nil {
				panic(fmt.Sprintf("[sync] Problems with gather data for a pulse - %v and jet - %v. err - %v", pn, jID.DebugString(), err))
			}
			err = t.sendToHeavy(ctx, msg)
			if err != nil {
				logger.Errorf("[sync] Problems with sending msg to a heavy node", err)
				t.addToNotSentPayloads(msg)
				continue
			}
			logger.Debugf("[sync] data has been sent to a heavy. pn - %v, jetID - %v", msg.PulseNum, msg.JetID.DebugString())
		}

		t.cleaner.NotifyAboutPulse(ctx, pn)
	}
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
		stats.Record(ctx,
			statRetryHeavyPayloadCount.M(1),
		)
	}
}

func (t *toHeavySyncer) sendToHeavy(ctx context.Context, data insolar.Message) error {
	rep, err := t.msgBus.Send(ctx, data, nil)
	if err != nil {
		stats.Record(ctx,
			statErrHeavyPayloadCount.M(1),
		)
		return err
	}
	if rep != nil {
		err, ok := rep.(*reply.HeavyError)
		if ok {
			stats.Record(ctx,
				statErrHeavyPayloadCount.M(1),
			)
			return err
		}
	}
	stats.Record(ctx,
		statHeavyPayloadCount.M(1),
	)
	return nil
}
