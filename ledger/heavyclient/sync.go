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

package heavyclient

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/drop"
)

func messageToHeavy(ctx context.Context, bus insolar.MessageBus, msg insolar.Message) error {
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
	pn insolar.PulseNumber,
) error {
	jetID := c.jetID
	inslog := inslogger.FromContext(ctx)
	inslog = inslog.WithField("jetID", jetID.DebugString())
	inslog = inslog.WithField("pulseNum", pn)

	signalMsg := &message.HeavyStartStop{
		JetID:    jetID,
		PulseNum: pn,
	}
	if err := messageToHeavy(ctx, c.bus, signalMsg); err != nil {
		inslog.Error("synchronize: start failed")
		return err
	}

	dr, err := c.dropAccessor.ForPulse(ctx, jetID, pn)
	if err != nil {
		inslog.Error("synchronize: can't fetch a drop")
		return err
	}

	replicator := storage.NewReplicaIter(
		ctx, c.db, insolar.ID(jetID), pn, pn+1, c.opts.SyncMessageLimit)
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
			Drop:     drop.Serialize(dr),
		}
		if err := messageToHeavy(ctx, c.bus, msg); err != nil {
			inslog.Error("synchronize: payload failed")
			return err
		}
	}

	signalMsg.Finished = true
	if err := messageToHeavy(ctx, c.bus, signalMsg); err != nil {
		inslog.Error("synchronize: finish failed")
		return err
	}

	return nil
}
