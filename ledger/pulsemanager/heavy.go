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

	"github.com/pkg/errors"

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
	retry bool,
) error {
	inslog := inslogger.FromContext(ctx)

	if retry {
		inslog.Infof("send reset message for pulse %v (retry sync)", pn)
		resetMsg := &message.HeavyReset{PulseNum: pn}
		if _, reseterr := m.Bus.Send(ctx, resetMsg, nil); reseterr != nil {
			return reseterr
		}
	}

	signalMsg := &message.HeavyStartStop{PulseNum: pn}
	_, starterr := m.Bus.Send(ctx, signalMsg, nil)
	// TODO: check if locked
	if starterr != nil {
		return starterr
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
			return senderr
		}
	}

	signalMsg.Finished = true
	_, stoperr := m.Bus.Send(ctx, signalMsg, nil)
	if stoperr != nil {
		return stoperr
	}
	inslog.Debugf("synchronize, sucessfully send start message for range [%v:%v]", pn, pn+1)

	lastmeetpulse := replicator.LastPulse()
	inslog.Debugf("synchronize on [%v:%v] finised (maximum record pulse is %v)",
		pn, pn+1, lastmeetpulse)
	return nil
}

// NextSyncPulses returns next pulse number for syncing to heavy node.
// If nothing to sync it returns 0, nil.
func (m *PulseManager) NextSyncPulses(ctx context.Context) (core.PulseNumber, error) {
	var (
		replicated core.PulseNumber
		err        error
	)
	if replicated, err = m.db.GetReplicatedPulse(ctx); err != nil {
		return 0, err
	}

	if replicated == 0 {
		return core.FirstPulseNumber, nil
	}
	return m.findnext(ctx, replicated)
}

func (m *PulseManager) findnext(ctx context.Context, from core.PulseNumber) (core.PulseNumber, error) {
	// start should be after replicated pulse or at least from FirstPulseNumber (for zero case)
	pulse, err := m.db.GetPulse(ctx, from)
	if err != nil {
		return 0, errors.Wrapf(err, "GetPulse with pulse num %v failed", from)
	}
	if pulse.Next == nil {
		return 0, nil
	}

	iwasalight, err := m.JetCoordinator.IsAuthorized(
		ctx,
		core.RoleLightExecutor,
		nil,
		*pulse.Next,
		m.NodeNet.GetOrigin().ID(),
	)
	if err != nil {
		return 0, errors.Wrapf(err, "Light checking failed")
	}
	if iwasalight {
		return *pulse.Next, nil
	}
	return m.findnext(ctx, *pulse.Next)
}
