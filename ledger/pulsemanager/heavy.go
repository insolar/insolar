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
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
)

// HeavyErr holds core.Reply and implements core.Retryable and error interfaces.
type HeavyErr struct {
	reply core.Reply
	err   error
}

// Error implements error interface.
func (he HeavyErr) Error() string {
	return he.err.Error()
}

// IsRetryable checks retryability of message.
func (he HeavyErr) IsRetryable() bool {
	herr, ok := he.reply.(*reply.HeavyError)
	if !ok {
		return false
	}
	return herr.ConcreteType() == reply.ErrHeavySyncInProgress
}

// HeavySync syncs records from light to heavy node, returns last synced pulse and error.
//
// It syncs records from start to end of provided pulse numbers.
func (m *PulseManager) HeavySync(
	ctx context.Context,
	pn core.PulseNumber,
	retry bool,
) error {
	inslog := inslogger.FromContext(ctx)
	var (
		busreply core.Reply
		buserr   error
	)

	if retry {
		inslog.Infof("send reset message for pulse %v (retry sync)", pn)
		resetMsg := &message.HeavyReset{PulseNum: pn}
		if busreply, buserr := m.Bus.Send(ctx, resetMsg, nil); buserr != nil {
			return HeavyErr{reply: busreply, err: buserr}
		}
	}

	signalMsg := &message.HeavyStartStop{PulseNum: pn}
	busreply, buserr = m.Bus.Send(ctx, signalMsg, nil)
	// TODO: check if locked
	if buserr != nil {
		return HeavyErr{reply: busreply, err: buserr}
	}
	inslog.Infof("synchronize, sucessfully send start message for pulse %v", pn)

	replicator := storage.NewReplicaIter(
		ctx, m.db, pn, pn+1, m.options.syncMessageLimit)
	for {
		recs, err := replicator.NextRecords()
		if err == storage.ErrReplicatorDone {
			break
		}
		if err != nil {
			panic(err)
		}
		msg := &message.HeavyPayload{Records: recs}
		busreply, buserr = m.Bus.Send(ctx, msg, nil)
		if buserr != nil {
			return HeavyErr{reply: busreply, err: buserr}
		}
	}

	signalMsg.Finished = true
	busreply, buserr = m.Bus.Send(ctx, signalMsg, nil)
	if buserr != nil {
		return HeavyErr{reply: busreply, err: buserr}
	}
	inslog.Infof("synchronize, sucessfully send finish message for pulse %v", pn)

	lastmeetpulse := replicator.LastPulse()
	inslog.Infof("synchronize on %v finised (maximum record pulse is %v)",
		pn, lastmeetpulse)
	return nil
}

// NextSyncPulses returns next pulse numbers for syncing to heavy node.
// If nothing to sync it returns nil, nil.
func (m *PulseManager) NextSyncPulses(ctx context.Context) ([]core.PulseNumber, error) {
	var (
		replicated core.PulseNumber
		err        error
	)
	if replicated, err = m.db.GetReplicatedPulse(ctx); err != nil {
		return nil, err
	}

	if replicated == 0 {
		return m.findAllCompleted(ctx, core.FirstPulseNumber)
	}
	next, nexterr := m.findnext(ctx, replicated)
	if nexterr != nil {
		return nil, nexterr
	}
	if next == nil {
		return nil, nil
	}
	return m.findAllCompleted(ctx, *next)
}

func (m *PulseManager) findAllCompleted(ctx context.Context, from core.PulseNumber) ([]core.PulseNumber, error) {
	wasalight, err := m.JetCoordinator.IsAuthorized(
		ctx,
		core.DynamicRoleLightExecutor,
		// TODO: pass JetID RecordRef here, when it would be ready
		nil,
		from,
		m.NodeNet.GetOrigin().ID(),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "Light checking failed")
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
	extra, err := m.findAllCompleted(ctx, *next)
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
