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
func (c *jetSyncClient) HeavySync(
	ctx context.Context,
	pn core.PulseNumber,
	retry bool,
) error {
	inslog := inslogger.FromContext(ctx)
	jetID := c.jetID

	current, err := c.PulseStorage.Current(ctx)
	if err != nil {
		return err
	}

	var (
		busreply core.Reply
		buserr   error
	)

	if retry {
		inslog.Infof("send reset message for pulse %v (retry sync)", pn)
		resetMsg := &message.HeavyReset{PulseNum: pn}
		if busreply, buserr := c.Bus.Send(ctx, resetMsg, *current, nil); buserr != nil {
			return HeavyErr{reply: busreply, err: buserr}
		}
	}

	signalMsg := &message.HeavyStartStop{PulseNum: pn}
	busreply, buserr = c.Bus.Send(ctx, signalMsg, *current, nil)
	// TODO: check if locked
	if buserr != nil {
		return HeavyErr{reply: busreply, err: buserr}
	}
	inslog.Infof("synchronize, sucessfully send start message for pulse %v", pn)

	replicator := storage.NewReplicaIter(
		ctx, c.db, jetID, pn, pn+1, c.syncMessageLimit)
	for {
		recs, err := replicator.NextRecords()
		if err == storage.ErrReplicatorDone {
			break
		}
		if err != nil {
			panic(err)
		}
		msg := &message.HeavyPayload{Records: recs}
		busreply, buserr = c.Bus.Send(ctx, msg, *current, nil)
		if buserr != nil {
			return HeavyErr{reply: busreply, err: buserr}
		}
	}

	signalMsg.Finished = true
	busreply, buserr = c.Bus.Send(ctx, signalMsg, *current, nil)
	if buserr != nil {
		return HeavyErr{reply: busreply, err: buserr}
	}
	inslog.Infof("synchronize, sucessfully send finish message for pulse %v", pn)

	lastMeetPulse := replicator.LastPulse()
	inslog.Infof("synchronize on %v finised (maximum record pulse is %v)",
		pn, lastMeetPulse)
	return nil
}

func (m *PulseManager) initJetSyncState(ctx context.Context) error {
	allJets, err := m.db.GetJets(ctx)
	if err != nil {
		return err
	}
	// not so effective, because we rescan pulses
	// but for now it is easier to do this in this way
	for jetID := range allJets {
		pulseNums, err := m.NextSyncPulses(ctx, jetID)
		if err != nil {
			return err
		}
		m.syncClientsPool.AddPulsesToSyncClient(ctx, jetID, false, pulseNums...)
	}
	return nil
}

// NextSyncPulses returns next pulse numbers for syncing to heavy node.
// If nothing to sync it returns nil, nil.
func (m *PulseManager) NextSyncPulses(ctx context.Context, jetID core.RecordID) ([]core.PulseNumber, error) {
	var (
		replicated core.PulseNumber
		err        error
	)
	if replicated, err = m.db.GetReplicatedPulse(ctx, jetID); err != nil {
		return nil, err
	}

	if replicated == 0 {
		return m.findAllCompleted(ctx, jetID, core.FirstPulseNumber)
	}
	next, nexterr := m.findnext(ctx, replicated)
	if nexterr != nil {
		return nil, nexterr
	}
	if next == nil {
		return nil, nil
	}
	return m.findAllCompleted(ctx, jetID, *next)
}

func (m *PulseManager) findAllCompleted(ctx context.Context, jetID core.RecordID, from core.PulseNumber) ([]core.PulseNumber, error) {
	wasalight, err := m.JetCoordinator.IsAuthorized(
		ctx,
		core.DynamicRoleLightExecutor,
		&jetID,
		from,
		m.NodeNet.GetOrigin().ID(),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "check 'am I light' for pulse num %v failed", from)
	}
	next, err := m.findnext(ctx, from)
	if err != nil {
		return nil, err
	}
	if next == nil {
		// if next is not found, we haven't got next pulse
		// in such case we don't want to replicate unfinished pulse
		return nil, nil
	}

	var found []core.PulseNumber
	if wasalight {
		found = append(found, from)
	}
	extra, err := m.findAllCompleted(ctx, jetID, *next)
	if err != nil {
		return nil, err
	}
	return append(found, extra...), nil
}

func (m *PulseManager) findnext(ctx context.Context, from core.PulseNumber) (*core.PulseNumber, error) {
	pulse, err := m.db.GetPulse(ctx, from)
	if err != nil {
		return nil, errors.Wrapf(err, "GetPulse with pulse num %v failed", from)
	}
	return pulse.Next, nil
}
