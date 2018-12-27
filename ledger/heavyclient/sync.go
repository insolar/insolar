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

package heavyclient

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
func (c *JetClient) HeavySync(
	ctx context.Context,
	pn core.PulseNumber,
	retry bool,
) error {
	inslog := inslogger.FromContext(ctx)
	jetID := c.jetID
	inslog = inslog.WithField("jetID", jetID).WithField("pulseNum", pn)

	inslog.Debug("JetClient.HeavySync")
	var (
		busreply core.Reply
		buserr   error
	)

	if retry {
		inslog.Info("synchronize: send reset message (retry sync)")
		resetMsg := &message.HeavyReset{PulseNum: pn}
		if busreply, buserr = c.Bus.Send(ctx, resetMsg, nil); buserr != nil {
			return HeavyErr{reply: busreply, err: buserr}
		}
	}

	signalMsg := &message.HeavyStartStop{PulseNum: pn}
	busreply, buserr = c.Bus.Send(ctx, signalMsg, nil)
	// TODO: check if locked
	if buserr != nil {
		inslog.Error("synchronize: start send error", buserr.Error())
		return HeavyErr{reply: busreply, err: buserr}
	}
	inslog.Debug("synchronize: sucessfully send start message")

	replicator := storage.NewReplicaIter(
		ctx, c.db, jetID, pn, pn+1, c.opts.SyncMessageLimit)
	for {
		recs, err := replicator.NextRecords()
		if err == storage.ErrReplicatorDone {
			break
		}
		if err != nil {
			panic(err)
		}
		msg := &message.HeavyPayload{
			PulseNum: pn,
			Records:  recs,
		}
		busreply, buserr = c.Bus.Send(ctx, msg, nil)
		if buserr != nil {
			inslog.Error("synchronize: payload send error", buserr.Error())
			return HeavyErr{reply: busreply, err: buserr}
		}
		inslog.Debug("synchronize: sucessfully send save message")
	}

	signalMsg.Finished = true
	busreply, buserr = c.Bus.Send(ctx, signalMsg, nil)
	if buserr != nil {
		inslog.Error("synchronize: finish send error", buserr.Error())
		return HeavyErr{reply: busreply, err: buserr}
	}
	inslog.Debug("synchronize: sucessfully send finish message")

	lastMeetPulse := replicator.LastSeenPulse()
	inslog.Debugf("synchronize: finished (maximum pulse of saved messages is %v)", lastMeetPulse)
	return nil
}
