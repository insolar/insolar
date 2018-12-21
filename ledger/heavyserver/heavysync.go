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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/ledger/storage"
)

// ErrSyncInProgress indicates that sync for provided jet is in sync
var ErrSyncInProgress = &reply.HeavyError{
	Message: "Heavy node sync in progress",
	SubType: reply.ErrHeavySyncInProgress,
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
			return errors.Wrap(err, "GetHeavySyncedPulse failed")
		}
	}

	// TODO: not sure how to handle this case properly
	if checkpoint == 0 {
		if pn != core.FirstPulseNumber {
			return errors.New("Pulse should be equal first pulse number if sync checkpoint on heavy not found")
		}
		return nil
	}

	if pn <= jetstate.lastok {
		return fmt.Errorf("Pulse %v is not greater than last synced pulse %v", pn, jetstate.lastok)
	}

	pulse, err := s.db.GetPulse(ctx, checkpoint)
	if err != nil {
		return errors.Wrapf(err, "GetPulse with pulse num %v failed", checkpoint)
	}
	if pulse.Next == nil {
		return fmt.Errorf("next pulse after %v not found", checkpoint)
	}

	if pn != *pulse.Next {
		return fmt.Errorf("pulse %v is not next after %v", pn, *pulse.Next)
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
		return ErrSyncInProgress
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
	jetState := s.getJetSyncState(ctx, jetID)

	err := func() error {
		jetState.Lock()
		defer jetState.Unlock()
		if jetState.syncpulse == nil {
			return fmt.Errorf("Jet %v not in sync mode", jetID)
		}
		if *jetState.syncpulse != pn {
			return fmt.Errorf("Passed pulse %v doesn't math in-sync pulse %v", pn, *jetState.syncpulse)
		}
		if jetState.insync {
			return ErrSyncInProgress
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
	return s.db.StoreKeyValues(ctx, kvs)
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
			"Passed pulse %v doesn't match pulse %v current in sync for jet %v",
			pn, *jetState.syncpulse, jetID)
	}
	if jetState.insync {
		return ErrSyncInProgress
	}
	jetState.syncpulse = nil

	err := s.db.SetHeavySyncedPulse(ctx, jetID, pn)
	if err != nil {
		return err
	}
	jetState.lastok = pn
	return nil
}

// Reset resets sync for provided pulse.
func (s *Sync) Reset(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) error {
	jetState := s.getJetSyncState(ctx, jetID)
	jetState.Lock()
	defer jetState.Unlock()

	if jetState.insync {
		return ErrSyncInProgress
	}

	jetState.syncpulse = nil
	return nil
}
