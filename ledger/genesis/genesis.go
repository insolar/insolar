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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/ledger/storage/pulse"
)

// BaseRecord provides methods for genesis base record manipulation.
type BaseRecord struct {
	PCS insolar.PlatformCryptographyScheme

	DB             store.DB
	DropModifier   drop.Modifier
	PulseAppender  pulse.Appender
	PulseAccessor  pulse.Accessor
	RecordModifier object.RecordModifier
	// TODO: @imarkin 28.03.2019 - remove ObjectStorage after all new storages integration (INS-2013, etc)
	ObjectStorage storage.ObjectStorage
}

// Key is genesis key.
type Key struct{}

func (Key) ID() []byte {
	return []byte{0x01}
}

func (Key) Scope() store.Scope {
	return store.ScopeSystem
}

// CreateIfNeeded creates new base genesis record if needed.
// Returns reference of genesis record and flag if base record have been created.
func (gi *BaseRecord) CreateIfNeeded(ctx context.Context) (*insolar.Reference, bool, error) {
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
		err := gi.PulseAppender.Append(
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

		lastPulse, err := gi.PulseAccessor.Latest(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "fail to get last pulse")
		}
		if lastPulse.PulseNumber != insolar.GenesisPulse.PulseNumber {
			return nil, fmt.Errorf(
				"last pulse number %v is not equal to genesis special value %v",
				lastPulse.PulseNumber,
				insolar.GenesisPulse.PulseNumber,
			)
		}

		virtRec := &object.GenesisRecord{}
		genesisID := object.NewRecordIDFromRecord(gi.PCS, lastPulse.PulseNumber, virtRec)
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
	genesisRef, err := getGenesisRef()
	if err == nil {
		return genesisRef, false, nil
	}
	if err != store.ErrNotFound {
		return nil, false, errors.Wrap(err, "genesis bootstrap failed")
	}

	genesisRef, err = createGenesisRecord()
	if err != nil {
		return nil, true, err
	}

	return genesisRef, true, nil
}
