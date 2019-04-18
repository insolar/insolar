/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package replication

import (
	"context"
	"fmt"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/pulse"
	"go.opencensus.io/stats"
)

// ToHeavySyncer is a base interface for a sync component
type ToHeavySyncer interface {
	// NotifyAboutPulse is method for notifying a sync component about new pulse
	NotifyAboutPulse(ctx context.Context, pn insolar.PulseNumber)
}

type toHeavySyncer struct {
	once sync.Once

	syncWaitingPulses chan insolar.PulseNumber

	jetCalculator   jet.Calculator
	dataGatherer    DataGatherer
	cleaner         Cleaner
	msgBus          insolar.MessageBus
	pulseCalculator pulse.Calculator
}

// NewToHeavySyncer creates new instance of ToHeavySyncer
func NewToHeavySyncer(
	jetCalculator jet.Calculator,
	dataGatherer DataGatherer,
	cleaner Cleaner,
	msgBus insolar.MessageBus,
	calculator pulse.Calculator,
) ToHeavySyncer {
	return &toHeavySyncer{
		jetCalculator:     jetCalculator,
		dataGatherer:      dataGatherer,
		cleaner:           cleaner,
		msgBus:            msgBus,
		pulseCalculator:   calculator,
		syncWaitingPulses: make(chan insolar.PulseNumber),
	}
}

// NotifyAboutPulse is method for notifying a sync component about new pulse
// When it's called, a provided pulse is added to a channel.
// There is a special gorutine that is reading that channel. When a new pulse is being received,
// the routine starts to gather data (with using of LightDataGatherer). After gathering all the data,
// it attempts to send it to the heavy. After sending a heavy payload to a heavy, data is deleted
// with help of Cleaner
func (t *toHeavySyncer) NotifyAboutPulse(ctx context.Context, pn insolar.PulseNumber) {
	t.once.Do(func() {
		go t.sync(ctx)
	})

	logger := inslogger.FromContext(ctx)
	logger.Debugf("[NotifyAboutPulse] pn - %v", pn)

	prevPN, err := t.pulseCalculator.Backwards(ctx, pn, 1)
	if err != nil {
		logger.Error("[NotifyAboutPulse]", err)
		return
	}

	logger.Debugf("[NotifyAboutPulse] prevPn - %v", prevPN.PulseNumber)
	t.syncWaitingPulses <- prevPN.PulseNumber
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
				panic(
					fmt.Sprintf(
						"[sync] Problems with gather data for a pulse - %v and jet - %v. err - %v",
						pn,
						jID.DebugString(),
						err,
					),
				)
			}
			err = t.sendToHeavy(ctx, msg)
			if err != nil {
				logger.Errorf("[sync] Problems with sending msg to a heavy node", err)
			} else {
				logger.Debugf("[sync] data has been sent to a heavy. pn - %v, jetID - %v", msg.PulseNum, msg.JetID.DebugString())
			}
		}

		t.cleaner.NotifyAboutPulse(ctx, pn)
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
