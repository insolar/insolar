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
		ctx, c.db, jetID, pn, pn+1, c.opts.SyncMessageLimit)
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
