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

package heavyserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/ledger/storage"
)

const defaultTimeout = time.Second * 10

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
	timer     *time.Timer
}

func (s *syncstate) resetTimeout(ctx context.Context, timeout time.Duration) {
	if s.timer != nil {
		s.timer.Reset(timeout)
	} else {
		s.timer = time.NewTimer(timeout)
	}
	timer := s.timer
	go func() {
		<-timer.C

		s.Lock()
		if s.timer == timer {
			stats.Record(ctx, statSyncedTimeout.M(1))
			s.syncpulse = nil
			s.timer = nil
		}
		s.Unlock()
	}()
}

type jetprefix [core.JetPrefixSize]byte

// Sync provides methods for syncing records to heavy storage.
type Sync struct {
	DropModifier   drop.Modifier          `inject:""`
	ReplicaStorage storage.ReplicaStorage `inject:""`
	DBContext      storage.DBContext

	sync.Mutex
	jetSyncStates map[jetprefix]*syncstate
}

// NewSync creates new Sync instance.
func NewSync(db storage.DBContext) *Sync {
	return &Sync{
		DBContext:     db,
		jetSyncStates: map[jetprefix]*syncstate{},
	}
}

func (s *Sync) checkIsNextPulse(ctx context.Context, jetID core.RecordID, jetstate *syncstate, pn core.PulseNumber) error {
	var (
		checkpoint core.PulseNumber
		err        error
	)

	checkpoint = jetstate.lastok
	if checkpoint == 0 {
		checkpoint, err = s.ReplicaStorage.GetHeavySyncedPulse(ctx, jetID)
		if err != nil {
			return errors.Wrap(err, "heavyserver: GetHeavySyncedPulse failed")
		}
	}

	// just start sync on first sync
	if checkpoint == 0 {
		return nil
	}

	if pn <= jetstate.lastok {
		return fmt.Errorf("heavyserver: pulse %v is not greater than last synced pulse %v (jet=%v)",
			pn, jetstate.lastok, jetID)
	}

	return nil
}

func (s *Sync) getJetSyncState(ctx context.Context, jetID core.RecordID) *syncstate {
	var jp jetprefix
	jpBuf := core.JetID(jetID).Prefix()
	copy(jp[:], jpBuf)
	s.Lock()
	jetState, ok := s.jetSyncStates[jp]
	if !ok {
		jetState = &syncstate{}
		s.jetSyncStates[jp] = jetState
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
		if *jetState.syncpulse >= pn {
			return fmt.Errorf("heavyserver: pulse %v is not greater than current in-sync pulse %v (jet=%v)",
				pn, *jetState.syncpulse, jetID)
		}
		return errSyncInProgress(jetID, pn)
	}

	if pn <= core.FirstPulseNumber {
		return fmt.Errorf("heavyserver: sync pulse should be greater than first pulse %v (got %v)", core.FirstPulseNumber, pn)
	}

	if err := s.checkIsNextPulse(ctx, jetID, jetState, pn); err != nil {
		return err
	}

	jetState.syncpulse = &pn
	jetState.resetTimeout(ctx, defaultTimeout)
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
		jetState.resetTimeout(ctx, defaultTimeout)
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
	err = s.DBContext.StoreKeyValues(ctx, kvs)
	if err != nil {
		return errors.Wrapf(err, "heavyserver: store failed")
	}

	// heavy stats
	recordsCount := int64(len(kvs))
	recordsSize := core.KVSize(kvs)
	inslog.Debugf("heavy store stat: JetID=%v, recordsCount+=%v, recordsSize+=%v\n", jetID.DebugString(), recordsCount, recordsSize)

	ctx = insmetrics.InsertTag(ctx, tagJet, jetID.DebugString())
	stats.Record(ctx,
		statSyncedCount.M(1),
		statSyncedRecords.M(recordsCount),
		statSyncedPulse.M(int64(pn)),
		statSyncedBytes.M(recordsSize),
	)
	return nil
}

// StoreDrop saves a jet.Drop to a heavy db
func (s *Sync) StoreDrop(ctx context.Context, jetID core.JetID, rawDrop []byte) error {
	err := s.DropModifier.Set(ctx, drop.Deserialize(rawDrop))
	if err != nil {
		return errors.Wrapf(err, "heavyserver: drop storing failed")
	}

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

	err := s.ReplicaStorage.SetHeavySyncedPulse(ctx, jetID, pn)
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

	if jetState.lastok == pn {
		// Sync is finished. No need to reset.
		return nil
	}

	inslogger.FromContext(ctx).Debugf("heavyserver: Reset sync: jetID=%v, pulse=%v", jetID, pn)
	jetState.syncpulse = nil
	return nil
}
