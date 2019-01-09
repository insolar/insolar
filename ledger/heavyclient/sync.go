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
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
)

func messageToHeavy(ctx context.Context, bus core.MessageBus, msg core.Message) error {
	busreply, buserr := bus.Send(ctx, msg, nil)
	if buserr != nil {
		return buserr
	}
	if busreply != nil {
		herr, ok := busreply.(*reply.HeavyError)
		if ok {
			return herr
		}
	}
	return nil
}

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
	if retry {
		inslog.Info("synchronize: send reset message (retry sync)")
		resetMsg := &message.HeavyReset{
			JetID:    jetID,
			PulseNum: pn,
		}
		if err := messageToHeavy(ctx, c.Bus, resetMsg); err != nil {
			return err
		}
	}

	signalMsg := &message.HeavyStartStop{
		JetID:    jetID,
		PulseNum: pn,
	}
	if err := messageToHeavy(ctx, c.Bus, signalMsg); err != nil {
		return err
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
			JetID:    jetID,
			PulseNum: pn,
			Records:  recs,
		}
		if err := messageToHeavy(ctx, c.Bus, msg); err != nil {
			return err
		}
		inslog.Debug("synchronize: sucessfully send save message")
	}

	signalMsg.Finished = true
	if err := messageToHeavy(ctx, c.Bus, signalMsg); err != nil {
		inslog.Error("synchronize: finish send error", err.Error())
		return err
	}
	inslog.Debug("synchronize: sucessfully send finish message")

	lastMeetPulse := replicator.LastSeenPulse()
	inslog.Debugf("synchronize: finished (maximum pulse of saved messages is %v)", lastMeetPulse)
	return nil
}
