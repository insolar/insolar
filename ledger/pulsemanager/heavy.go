/*
 *    Copyright 2018 Insolar
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

package pulsemanager

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
)

// HeavySync syncs records from light to heavy node, returns last synced pulse and error.
//
// It syncs records from start to end of provided pulse numbers.
func (m *PulseManager) HeavySync(
	ctx context.Context,
	pn core.PulseNumber,
) (core.PulseNumber, error) {
	inslog := inslogger.FromContext(ctx)

	signalMsg := &message.HeavyStartStop{PulseNum: pn}
	_, starterr := m.Bus.Send(ctx, signalMsg, nil)
	// TODO: check if locked
	if starterr != nil {
		return 0, starterr
	}
	inslog.Debugf("synchronize, sucessfully send start message for range [%v:%v]", pn, pn+1)

	replicator := storage.NewReplicaIter(
		ctx, m.db, pn, pn+1, m.options.syncmessagelimit)
	for {
		recs, err := replicator.NextRecords()
		if err == storage.ErrReplicatorDone {
			break
		}
		if err != nil {
			panic(err)
		}
		msg := &message.HeavyPayload{Records: recs}
		_, senderr := m.Bus.Send(ctx, msg, nil)
		if senderr != nil {
			return 0, senderr
		}
	}

	signalMsg.Finished = true
	_, stoperr := m.Bus.Send(ctx, signalMsg, nil)
	if stoperr != nil {
		return 0, stoperr
	}
	inslog.Debugf("synchronize, sucessfully send start message for range [%v:%v]", pn, pn+1)

	lastmeetpulse := replicator.LastPulse()
	inslog.Debugf("synchronize on [%v:%v] finised (maximum record pulse is %v)",
		pn, pn+1, lastmeetpulse)
	return lastmeetpulse, nil
}

// NextSyncPulses returns pulse numbers range for syncing to heavy node.
// If nothing to sync it returns 0, 0, nil.
func (m *PulseManager) NextSyncPulses(ctx context.Context) (start, end core.PulseNumber, err error) {
	var (
		replicated core.PulseNumber
		last       core.PulseNumber
	)
	if replicated, err = m.db.GetReplicatedPulse(ctx); err != nil {
		return
	}
	if last, err = m.db.GetLastPulseAsLightMaterial(ctx); err != nil {
		return
	}
	// if replicated pulse is not less than "last light material pulse", nothing to sync
	if replicated >= last {
		return
	}

	// start should be after replicated pulse or at least from FirstPulseNumber (for zero case)
	start = replicated + 1
	if replicated == 0 {
		start = core.FirstPulseNumber
	}
	// end should be after "last light material pulse" + 1 or at least next pulse after start
	end = last + 1
	if last == 0 {
		end = start + 1
	}
	return
}
