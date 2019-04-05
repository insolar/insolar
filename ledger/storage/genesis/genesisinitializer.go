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

package genesis

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/object"
)

type State interface {
	component.Initer
	GenesisRef() *insolar.Reference
	Happened() bool
}

type genesisInitializer struct {
	PCS insolar.PlatformCryptographyScheme `inject:""`

	DB             store.DB              `inject:""`
	DropModifier   drop.Modifier         `inject:""`
	RecordModifier object.RecordModifier `inject:""`

	// TODO: @imarkin 28.03.2019 - remove ObjectStorage and  PulseTracker after all new storages integration (INS-2013, etc)
	ObjectStorage storage.ObjectStorage `inject:""`
	PulseTracker  storage.PulseTracker  `inject:""`

	genesisRef *insolar.Reference
	happened   bool
}

func NewGenesisInitializer() State {
	return new(genesisInitializer)
}

// GenesisRef returns the genesis record reference.
//
// Genesis record is the parent for all top-level records.
func (gi *genesisInitializer) GenesisRef() *insolar.Reference {
	return gi.genesisRef
}

type Key struct{}

func (gk Key) ID() []byte {
	return []byte{0x01}
}

func (gk Key) Scope() store.Scope {
	return store.ScopeSystem
}

func (gi *genesisInitializer) Happened() bool {
	return gi.happened
}

func (gi *genesisInitializer) Init(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("start storage bootstrap")

	getGenesisRef := func() (*insolar.Reference, error) {
		buff, err := gi.DB.Get(Key{})
		if err != nil {
			return nil, err
		}
		var genesisRef insolar.Reference
		copy(genesisRef[:], buff)
		return &genesisRef, nil
	}

	createGenesisRecord := func() (*insolar.Reference, error) {
		err := gi.PulseTracker.AddPulse(
			ctx,
			insolar.Pulse{
				PulseNumber: insolar.GenesisPulse.PulseNumber,
				Entropy:     insolar.GenesisPulse.Entropy,
			},
		)
		if err != nil {
			return nil, errors.Wrap(err, "fail to set genesis pulse")
		}
		// Add initial drop
		err = gi.DropModifier.Set(ctx, drop.Drop{JetID: insolar.ZeroJetID})
		if err != nil {
			return nil, errors.Wrap(err, "fail to set initial drop")
		}

		lastPulse, err := gi.PulseTracker.GetLatestPulse(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "fail to get last pulse")
		}
		if lastPulse.Pulse.PulseNumber != insolar.GenesisPulse.PulseNumber {
			return nil, fmt.Errorf(
				"last pulse number %v is not equal to genesis special value %v",
				lastPulse.Pulse.PulseNumber,
				insolar.GenesisPulse.PulseNumber,
			)
		}

		virtRec := &object.GenesisRecord{}
		genesisID := object.NewRecordIDFromRecord(gi.PCS, lastPulse.Pulse.PulseNumber, virtRec)
		rec := record.MaterialRecord{
			Record: virtRec,
			JetID:  insolar.ZeroJetID,
		}
		err = gi.RecordModifier.Set(ctx, *genesisID, rec)
		if err != nil {
			return nil, errors.Wrap(err, "can't save genesis record into storage")
		}

		err = gi.ObjectStorage.SetObjectIndex(
			ctx,
			insolar.ID(insolar.ZeroJetID),
			genesisID,
			&object.Lifeline{
				LatestState:         genesisID,
				LatestStateApproved: genesisID,
			},
		)
		if err != nil {
			return nil, errors.Wrap(err, "fail to set genesis index")
		}

		genesisRef := insolar.NewReference(*genesisID, *genesisID)
		return genesisRef, gi.DB.Set(Key{}, genesisRef[:])
	}

	var err error
	gi.genesisRef, err = getGenesisRef()
	if err == nil {
		return nil
	}
	if err != store.ErrNotFound {
		return errors.Wrap(err, "genesis bootstrap failed")
	}

	gi.genesisRef, err = createGenesisRecord()
	if err != nil {
		return err
	}

	gi.happened = true
	return nil
}
