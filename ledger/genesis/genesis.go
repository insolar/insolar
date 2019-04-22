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
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
)

// BaseRecord provides methods for genesis base record manipulation.
type BaseRecord struct {
	DB             store.DB
	DropModifier   drop.Modifier
	PulseAppender  pulse.Appender
	PulseAccessor  pulse.Accessor
	RecordModifier object.RecordModifier
	IndexModifier  object.IndexModifier
}

// Key is genesis key.
type Key struct{}

func (Key) ID() []byte {
	return []byte{0x01}
}

func (Key) Scope() store.Scope {
	return store.ScopeGenesis
}

// CreateIfNeeded creates new base genesis record if needed.
// Returns reference of genesis record and flag if base record have been created.
func (gi *BaseRecord) CreateIfNeeded(ctx context.Context) (bool, error) {
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

	createGenesisRecord := func() error {
		err := gi.PulseAppender.Append(
			ctx,
			insolar.Pulse{
				PulseNumber: insolar.GenesisPulse.PulseNumber,
				Entropy:     insolar.GenesisPulse.Entropy,
			},
		)
		if err != nil {
			return errors.Wrap(err, "fail to set genesis pulse")
		}
		// Add initial drop
		err = gi.DropModifier.Set(ctx, drop.Drop{JetID: insolar.ZeroJetID})
		if err != nil {
			return errors.Wrap(err, "fail to set initial drop")
		}

		lastPulse, err := gi.PulseAccessor.Latest(ctx)
		if err != nil {
			return errors.Wrap(err, "fail to get last pulse")
		}
		if lastPulse.PulseNumber != insolar.GenesisPulse.PulseNumber {
			return fmt.Errorf(
				"last pulse number %v is not equal to genesis special value %v",
				lastPulse.PulseNumber,
				insolar.GenesisPulse.PulseNumber,
			)
		}

		genesisID := insolar.GenesisRecord.ID()
		rec := record.MaterialRecord{
			Record: &object.GenesisRecord{
				VirtualRecord: insolar.GenesisRecord,
			},
			JetID: insolar.ZeroJetID,
		}
		err = gi.RecordModifier.Set(ctx, genesisID, rec)
		if err != nil {
			return errors.Wrap(err, "can't save genesis record into storage")
		}

		err = gi.IndexModifier.Set(
			ctx,
			genesisID,
			object.Lifeline{
				LatestState:         &genesisID,
				LatestStateApproved: &genesisID,
				JetID:               insolar.ZeroJetID,
			},
		)
		if err != nil {
			return errors.Wrap(err, "fail to set genesis index")
		}

		return gi.DB.Set(Key{}, insolar.GenesisRecord.Ref().Bytes())
	}

	_, err := getGenesisRef()
	if err == nil {
		return false, nil
	}
	if err != store.ErrNotFound {
		return false, errors.Wrap(err, "genesis bootstrap failed")
	}

	err = createGenesisRecord()
	if err != nil {
		return true, err
	}

	return true, nil
}
