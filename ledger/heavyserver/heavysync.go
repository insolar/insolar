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

package heavyserver

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/ledger/storage"
)

func errSyncInProgress(jetID core.RecordID, pn core.PulseNumber) *reply.HeavyError {
	return &reply.HeavyError{
		Message:  "Heavy node sync in progress",
		SubType:  reply.ErrHeavySyncInProgress,
		JetID:    jetID,
		PulseNum: pn,
	}
}

// in testnet we start with only one jet
type syncstate struct {
	sync.Mutex
	lastok core.PulseNumber
	// insyncend core.PulseNumber
	syncpulse *core.PulseNumber
	insync    bool
}

// Sync provides methods for syncing records to heavy storage.
type Sync struct {
	db *storage.DB

	sync.Mutex
	jetSyncStates map[core.RecordID]*syncstate
}

// NewSync creates new Sync instance.
func NewSync(db *storage.DB) *Sync {
	return &Sync{
		db:            db,
		jetSyncStates: map[core.RecordID]*syncstate{},
	}
}

func (s *Sync) checkIsNextPulse(ctx context.Context, jetID core.RecordID, jetstate *syncstate, pn core.PulseNumber) error {
	var (
		checkpoint core.PulseNumber
		err        error
	)

	checkpoint = jetstate.lastok
	if checkpoint == 0 {
		checkpoint, err = s.db.GetHeavySyncedPulse(ctx, jetID)
		if err != nil {
			return errors.Wrap(err, "heavyserver: GetHeavySyncedPulse failed")
		}
	}

	// just start sync on first sync
	if checkpoint == 0 {
		return nil
	}

	if pn <= jetstate.lastok {
		return fmt.Errorf("heavyserver: pulse %v is not greater than last synced pulse %v", pn, jetstate.lastok)
	}

	return nil
}

func (s *Sync) getJetSyncState(ctx context.Context, jetID core.RecordID) *syncstate {
	s.Lock()
	jetState, ok := s.jetSyncStates[jetID]
	if !ok {
		jetState = &syncstate{}
		s.jetSyncStates[jetID] = jetState
	}
	s.Unlock()
	return jetState
}

// Start try to start heavy sync for provided pulse.
func (s *Sync) Start(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) error {
	jetState := s.getJetSyncState(ctx, jetID)
	jetState.Lock()
	defer jetState.Unlock()

	if jetState.syncpulse != nil {
		return errSyncInProgress(jetID, pn)
	}

	if pn <= core.FirstPulseNumber {
		return fmt.Errorf("heavyserver: sync pulse should be greater than first pulse %v (got %v)", core.FirstPulseNumber, pn)
	}

	if err := s.checkIsNextPulse(ctx, jetID, jetState, pn); err != nil {
		return err
	}

	jetState.syncpulse = &pn
	return nil
}

// Store stores recieved key/value pairs at heavy storage.
//
// TODO: check actual jet and pulse in keys
func (s *Sync) Store(ctx context.Context, jetID core.RecordID, pn core.PulseNumber, kvs []core.KV) error {
	inslog := inslogger.FromContext(ctx)
	jetState := s.getJetSyncState(ctx, jetID)

	err := func() error {
		jetState.Lock()
		defer jetState.Unlock()
		if jetState.syncpulse == nil {
			return fmt.Errorf("heavyserver: jet %v not in sync mode", jetID)
		}
		if *jetState.syncpulse != pn {
			return fmt.Errorf("heavyserver: passed pulse %v doesn't match in-sync pulse %v", pn, *jetState.syncpulse)
		}
		if jetState.insync {
			return errSyncInProgress(jetID, pn)
		}
		jetState.insync = true
		return nil
	}()
	if err != nil {
		return err
	}

	defer func() {
		jetState.Lock()
		jetState.insync = false
		jetState.Unlock()
	}()
	// TODO: check jet in keys?
	err = s.db.StoreKeyValues(ctx, kvs)
	if err != nil {
		return errors.Wrapf(err, "heavyserver: store failed")
	}

	// heavy stats
	recordsCount := int64(len(kvs))
	recordsSize := core.KVSize(kvs)
	inslog.Debugf("heavy store stat: JetID=%v, recordsCount+=%v, recordsSize+=%v\n", jetID.String(), recordsCount, recordsSize)

	ctx = insmetrics.InsertTag(ctx, tagJet, jetID.String())
	stats.Record(ctx,
		statSyncedCount.M(1),
		statSyncedRecords.M(recordsCount),
		statSyncedPulse.M(int64(pn)),
		statSyncedBytes.M(recordsSize),
	)
	return nil
}

// Stop successfully stops replication for specified pulse.
//
// TODO: call Stop if range sync too long
func (s *Sync) Stop(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) error {
	jetState := s.getJetSyncState(ctx, jetID)
	jetState.Lock()
	defer jetState.Unlock()

	if jetState.syncpulse == nil {
		return errors.Errorf("Jet %v not in sync mode", jetID)
	}
	if *jetState.syncpulse != pn {
		return fmt.Errorf(
			"heavyserver: Passed pulse %v doesn't match pulse %v current in sync for jet %v",
			pn, *jetState.syncpulse, jetID)
	}
	if jetState.insync {
		return errSyncInProgress(jetID, pn)
	}
	jetState.syncpulse = nil

	err := s.db.SetHeavySyncedPulse(ctx, jetID, pn)
	if err != nil {
		return err
	}
	inslogger.FromContext(ctx).Debugf("heavyserver: Fin sync: jetID=%v, pulse=%v", jetID, pn)
	jetState.lastok = pn
	return nil
}

// Reset resets sync for provided pulse.
func (s *Sync) Reset(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) error {
	jetState := s.getJetSyncState(ctx, jetID)
	jetState.Lock()
	defer jetState.Unlock()

	if jetState.insync {
		return errSyncInProgress(jetID, pn)
	}

	inslogger.FromContext(ctx).Debugf("heavyserver: Reset sync: jetID=%v, pulse=%v", jetID, pn)
	jetState.syncpulse = nil
	return nil
}
