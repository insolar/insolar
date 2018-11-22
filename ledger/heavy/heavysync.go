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
	"errors"
	"fmt"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
)

// processable errors by client (i.e. it could retry)
var (
	ErrSyncInProgress = errors.New("Heavy node already syncing")
)

// in testnet we start with only one jet
type syncstate struct {
	sync.Mutex
	lastok core.PulseNumber
	// insyncend core.PulseNumber
	syncrange *core.PulseRange
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

// Start try to start heavy sync in provided range of pulses.
func (s *Sync) Start(ctx context.Context, prange core.PulseRange) error {
	s.Lock()
	defer s.Unlock()
	if s.lastok == 0 {
		pnum, err := s.db.GetHeavySyncedPulse(ctx)
		if err != nil {
			return err
		}
		s.lastok = pnum
	}

	if prange.Begin >= prange.End {
		return errors.New("Wrong pulse range")
	}

	if s.syncrange != nil {
		return ErrSyncInProgress
	}

	if prange.Begin <= s.lastok {
		return errors.New("Range already synced")
	}

	if s.lastok == 0 {
		if prange.Begin != core.FirstPulseNumber {
			return errors.New("Range should start with first pulse on empty heavy")
		}
		// for passing next check
		// TODO: move to storage method logic?
		s.lastok = prange.Begin - 1
	}

	// TODO: increase lastok by one if lastok more than N pulses in past
	if prange.Begin-1 > s.lastok {
		return fmt.Errorf("last synced pulse is %v, but range is %v", s.lastok, prange)
	}
	s.syncrange = &prange
	return nil
}

// Store stores recieved key/value pairs at heavy storage.
//
// TODO: check actual pulse in keys
func (s *Sync) Store(ctx context.Context, prange core.PulseRange, kvs []core.KV) error {
	err := func() error {
		s.Lock()
		defer s.Unlock()
		if s.syncrange == nil {
			return errors.New("Jet not in sync mode")
		}
		if *s.syncrange != prange {
			return errors.New("Passed range doesn't match range in sync")
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
	return s.db.StoreKeyValues(ctx, kvs)
}

// Stop stops replication with specified pulses range.
//
// TODO: call Stop if range sync too long
func (s *Sync) Stop(ctx context.Context, prange core.PulseRange) error {
	s.Lock()
	defer s.Unlock()
	if s.syncrange == nil {
		return errors.New("Jet not in sync mode")
	}
	if *s.syncrange != prange {
		return errors.New("Passed range doesn't match range in sync")
	}
	if s.insync {
		return errors.New("Can't stop heavy repliction that still in store mode")
	}
	s.syncrange = nil

	// TODO: store lastok
	lastok := prange.End - 1
	err := s.db.SetHeavySyncedPulse(ctx, lastok)
	if err != nil {
		return err
	}
	s.lastok = lastok
	return nil
}
