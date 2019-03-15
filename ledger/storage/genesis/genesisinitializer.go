/*
 *    Copyright 2019 Insolar
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

package genesis

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/object"
)

type GenesisState interface {
	component.Initer
	GenesisRef() *core.RecordRef
}

type genesisInitializer struct {
	DB            storage.DBContext     `inject:""`
	ObjectStorage storage.ObjectStorage `inject:""`
	PulseTracker  storage.PulseTracker  `inject:""`
	DropModifier  drop.Modifier         `inject:""`

	genesisRef *core.RecordRef
}

func NewGenesisInitializer() GenesisState {
	return new(genesisInitializer)
}

// GenesisRef returns the genesis record reference.
//
// Genesis record is the parent for all top-level records.
func (gi *genesisInitializer) GenesisRef() *core.RecordRef {
	return gi.genesisRef
}

func (gi *genesisInitializer) Init(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("start storage bootstrap")
	jetID := *core.NewJetID(0, nil)

	getGenesisRef := func() (*core.RecordRef, error) {
		buff, err := gi.DB.Get(ctx, storage.GenesisPrefixKey())
		if err != nil {
			return nil, err
		}
		var genesisRef core.RecordRef
		copy(genesisRef[:], buff)
		return &genesisRef, nil
	}

	createGenesisRecord := func() (*core.RecordRef, error) {
		err := gi.PulseTracker.AddPulse(
			ctx,
			core.Pulse{
				PulseNumber: core.GenesisPulse.PulseNumber,
				Entropy:     core.GenesisPulse.Entropy,
			},
		)
		if err != nil {
			return nil, err
		}
		// It should be 0. Because pulse after 65537 will try to use a hash of drop between 0 - 65537
		err = gi.DropModifier.Set(ctx, jetID, drop.Drop{})
		if err != nil {
			return nil, err
		}

		lastPulse, err := gi.PulseTracker.GetLatestPulse(ctx)
		if err != nil {
			return nil, err
		}
		genesisID, err := gi.ObjectStorage.SetRecord(ctx, core.RecordID(jetID), lastPulse.Pulse.PulseNumber, &object.GenesisRecord{})
		if err != nil {
			return nil, err
		}
		err = gi.ObjectStorage.SetObjectIndex(
			ctx,
			core.RecordID(jetID),
			genesisID,
			&object.Lifeline{LatestState: genesisID, LatestStateApproved: genesisID},
		)
		if err != nil {
			return nil, err
		}

		genesisRef := core.NewRecordRef(*genesisID, *genesisID)
		return genesisRef, gi.DB.Set(ctx, storage.GenesisPrefixKey(), genesisRef[:])
	}

	var err error
	gi.genesisRef, err = getGenesisRef()
	if err == core.ErrNotFound {
		gi.genesisRef, err = createGenesisRecord()
	}
	if err != nil {
		return errors.Wrap(err, "bootstrap failed")
	}

	return nil
}
