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

package proc

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

type Replication struct {
	message payload.Meta
	cfg     configuration.Ledger

	dep struct {
		records          object.RecordModifier
		recordsPositions object.RecordPositionModifier
		indexes          object.IndexModifier
		pcs              insolar.PlatformCryptographyScheme
		pulses           pulse.Accessor
		drops            drop.Modifier
		jets             jet.Modifier
		keeper           executor.JetKeeper
	}
}

func NewReplication(msg payload.Meta, cfg configuration.Ledger) *Replication {
	return &Replication{
		message: msg,
		cfg:     cfg,
	}
}

func (p *Replication) Dep(
	records object.RecordModifier,
	indexes object.IndexModifier,
	recordsPositions object.RecordPositionModifier,
	pcs insolar.PlatformCryptographyScheme,
	pulses pulse.Accessor,
	drops drop.Modifier,
	jets jet.Modifier,
	keeper executor.JetKeeper,
) {
	p.dep.records = records
	p.dep.indexes = indexes
	p.dep.recordsPositions = recordsPositions
	p.dep.pcs = pcs
	p.dep.pulses = pulses
	p.dep.drops = drops
	p.dep.jets = jets
	p.dep.keeper = keeper
}

func (p *Replication) Proceed(ctx context.Context) error {
	pl, err := payload.Unmarshal(p.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload")
	}
	msg, ok := pl.(*payload.Replication)
	if !ok {
		return fmt.Errorf("unexpected payload %T", pl)
	}

	storeRecords(ctx, p.dep.records, p.dep.recordsPositions, p.dep.pcs, msg.Pulse, msg.Records)
	if err := storeIndexes(ctx, p.dep.indexes, msg.Indexes, msg.Pulse); err != nil {
		return errors.Wrap(err, "failed to store indexes")
	}

	dr, err := storeDrop(ctx, p.dep.drops, msg.Drop)
	if err != nil {
		return errors.Wrap(err, "failed to store drop")
	}

	if err := p.dep.keeper.AddDropConfirmation(ctx, dr.Pulse, dr.JetID, dr.Split); err != nil {
		return errors.Wrapf(err, "failed to add jet to JetKeeper jet=%v", dr.JetID.DebugString())
	}

	stats.Record(ctx,
		statReceivedHeavyPayloadCount.M(1),
	)

	return nil
}

func storeIndexes(
	ctx context.Context,
	mod object.IndexModifier,
	indexes []record.Index,
	pn insolar.PulseNumber,
) error {
	for _, idx := range indexes {
		err := mod.SetIndex(ctx, pn, idx)
		if err != nil {
			return errors.Wrapf(err, "heavyserver: index storing failed")
		}
	}

	return nil
}

func storeDrop(
	ctx context.Context,
	drops drop.Modifier,
	rawDrop []byte,
) (*drop.Drop, error) {
	d, err := drop.Decode(rawDrop)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
		return nil, err
	}
	err = drops.Set(ctx, *d)
	if err != nil {
		return nil, errors.Wrapf(err, "heavyserver: drop storing failed")
	}

	return d, nil
}

func storeRecords(
	ctx context.Context,
	recordStorage object.RecordModifier,
	recordIndex object.RecordPositionModifier,
	pcs insolar.PlatformCryptographyScheme,
	pn insolar.PulseNumber,
	records []record.Material,
) {
	inslog := inslogger.FromContext(ctx)

	for _, rec := range records {
		hash := record.HashVirtual(pcs.ReferenceHasher(), rec.Virtual)
		id := *insolar.NewID(pn, hash)
		// FIXME: skipping errors will lead to inconsistent state.
		if rec.ID != id {
			inslog.Error(fmt.Errorf(
				"record id does not match (calculated: %s, received: %s)",
				id.DebugString(),
				rec.ID.DebugString(),
			))
			continue
		}
		err := recordStorage.Set(ctx, rec)
		if err != nil {
			inslog.Error(err, "heavyserver: store record failed")
			continue
		}
		err = recordIndex.IncrementPosition(id)
		if err != nil {
			inslog.Error(err, "heavyserver: fail to store record position")
			continue
		}
	}
}
