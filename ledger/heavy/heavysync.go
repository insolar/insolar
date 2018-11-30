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

package heavy

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
	syncstate
}

// NewSync creates new Sync instance.
func NewSync(db *storage.DB) *Sync {
	return &Sync{
		db: db,
	}
}

func (s *Sync) checkIsNextPulse(ctx context.Context, pn core.PulseNumber) error {
	var (
		checkpoint core.PulseNumber
		err        error
	)

	checkpoint = s.lastok
	if checkpoint == 0 {
		checkpoint, err = s.db.GetHeavySyncedPulse(ctx)
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

	if pn <= s.lastok {
		return fmt.Errorf("Pulse %v is not greater than last synced pulse %v", pn, s.lastok)
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

// Start try to start heavy sync for provided pulse.
func (s *Sync) Start(ctx context.Context, pn core.PulseNumber) error {
	s.Lock()
	defer s.Unlock()

	if s.syncpulse != nil {
		return ErrSyncInProgress
	}

	if err := s.checkIsNextPulse(ctx, pn); err != nil {
		return err
	}

	s.syncpulse = &pn
	return nil
}

// Store stores recieved key/value pairs at heavy storage.
//
// TODO: check actual pulse in keys
func (s *Sync) Store(ctx context.Context, pn core.PulseNumber, kvs []core.KV) error {
	err := func() error {
		s.Lock()
		defer s.Unlock()
		if s.syncpulse == nil {
			return errors.New("Jet not in sync mode")
		}
		if *s.syncpulse != pn {
			return fmt.Errorf("Passed pulse %v doesn't math in-sync pulse %v", pn, *s.syncpulse)
		}
		s.insync = true
		return nil
	}()
	if err != nil {
		return err
	}

	defer func() {
		s.Lock()
		s.insync = false
		s.Unlock()
	}()
	// could be retryable error only on low level storage issues
	// TODO: check error value and wrap to repeatable if it is not integrity errors
	return s.db.StoreKeyValues(ctx, kvs)
}

// Stop successfully stops replication for specified pulse.
//
// TODO: call Stop if range sync too long
func (s *Sync) Stop(ctx context.Context, pn core.PulseNumber) error {
	s.Lock()
	defer s.Unlock()
	if s.syncpulse == nil {
		return errors.New("Jet not in sync mode")
	}
	if *s.syncpulse != pn {
		return fmt.Errorf("Passed pulse %v doesn't match pulse %v current in sync", pn, *s.syncpulse)
	}
	if s.insync {
		return ErrSyncInProgress
	}
	s.syncpulse = nil

	err := s.db.SetHeavySyncedPulse(ctx, pn)
	if err != nil {
		return err
	}
	s.lastok = pn
	return nil
}

// Reset resets sync for provided pulse.
func (s *Sync) Reset(ctx context.Context, pn core.PulseNumber) error {
	s.Lock()
	defer s.Unlock()

	if s.insync {
		return ErrSyncInProgress
	}

	s.syncpulse = nil
	return nil
}
