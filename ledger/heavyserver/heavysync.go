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

	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/ledger/storage"
)

const defaultTimeout = time.Second * 10

func errSyncInProgress(jetID insolar.ID, pn insolar.PulseNumber) *reply.HeavyError {
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
	lastok insolar.PulseNumber
	// insyncend insolar.PulseNumber
	syncpulse *insolar.PulseNumber
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

type jetprefix [insolar.JetPrefixSize]byte

// Sync provides methods for syncing records to heavy storage.
type Sync struct {
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
	DropModifier               drop.Modifier                      `inject:""`
	BlobModifier               blob.Modifier                      `inject:""`
	ReplicaStorage             storage.ReplicaStorage             `inject:""`
	DBContext                  storage.DBContext

	RecordModifier object.RecordModifier

	sync.Mutex
	jetSyncStates map[jetprefix]*syncstate
}

// NewSync creates new Sync instance.
func NewSync(db storage.DBContext, records object.RecordModifier) *Sync {
	return &Sync{
		DBContext:      db,
		RecordModifier: records,
		jetSyncStates:  map[jetprefix]*syncstate{},
	}
}

func (s *Sync) checkIsNextPulse(ctx context.Context, jetID insolar.ID, jetstate *syncstate, pn insolar.PulseNumber) error {
	var (
		checkpoint insolar.PulseNumber
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

func (s *Sync) getJetSyncState(ctx context.Context, jetID insolar.ID) *syncstate {
	var jp jetprefix
	jpBuf := insolar.JetID(jetID).Prefix()
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
func (s *Sync) Start(ctx context.Context, jetID insolar.ID, pn insolar.PulseNumber) error {
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

	if pn <= insolar.FirstPulseNumber {
		return fmt.Errorf("heavyserver: sync pulse should be greater than first pulse %v (got %v)", insolar.FirstPulseNumber, pn)
	}

	if err := s.checkIsNextPulse(ctx, jetID, jetState, pn); err != nil {
		return err
	}

	jetState.syncpulse = &pn
	jetState.resetTimeout(ctx, defaultTimeout)
	return nil
}

// StoreIndices stores recieved key/value pairs for indices at heavy storage.
func (s *Sync) StoreIndices(ctx context.Context, jet insolar.ID, pn insolar.PulseNumber, kvs []insolar.KV) error {
	inslog := inslogger.FromContext(ctx)
	jetState := s.getJetSyncState(ctx, jet)

	err := func() error {
		jetState.Lock()
		defer jetState.Unlock()

		if jetState.syncpulse == nil {
			return fmt.Errorf("heavyserver: jet %v not in sync mode", jet)
		}
		if *jetState.syncpulse != pn {
			return fmt.Errorf("heavyserver: passed pulse %v doesn't match in-sync pulse %v", pn, *jetState.syncpulse)
		}
		if jetState.insync {
			return errSyncInProgress(jet, pn)
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
	err = s.DBContext.StoreKeyValues(ctx, kvs)
	if err != nil {
		return errors.Wrapf(err, "heavyserver: store failed")
	}

	// heavy stats
	recordsCount := int64(len(kvs))
	recordsSize := insolar.KVSize(kvs)
	inslog.Debugf("heavy store stat: JetID=%v, recordsCount+=%v, recordsSize+=%v\n", jet.DebugString(), recordsCount, recordsSize)

	ctx = insmetrics.InsertTag(ctx, tagJet, jet.DebugString())
	stats.Record(ctx,
		statSyncedCount.M(1),
		statSyncedRecords.M(recordsCount),
		statSyncedPulse.M(int64(pn)),
		statSyncedBytes.M(recordsSize),
	)
	return nil
}

// StoreDrop saves a jet.Drop to a heavy db
func (s *Sync) StoreDrop(ctx context.Context, jetID insolar.JetID, rawDrop []byte) error {
	d, err := drop.Decode(rawDrop)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
		return err
	}
	err = s.DropModifier.Set(ctx, *d)
	if err != nil {
		return errors.Wrapf(err, "heavyserver: drop storing failed")
	}

	return nil
}

// StoreBlobs saves a collection of blobs to a heavy's storage
func (s *Sync) StoreBlobs(ctx context.Context, pn insolar.PulseNumber, rawBlobs [][]byte) error {
	for _, rwb := range rawBlobs {
		b, err := blob.Decode(rwb)
		if err != nil {
			inslogger.FromContext(ctx).Error(err)
			continue
		}

		blobID := object.CalculateIDForBlob(s.PlatformCryptographyScheme, pn, rwb)
		err = s.BlobModifier.Set(ctx, *blobID, *b)
		if err != nil {
			return errors.Wrapf(err, "heavyserver: blob storing failed")
		}
	}
	return nil
}

// StoreRecords stores recieved records at heavy storage.
func (s *Sync) StoreRecords(ctx context.Context, jetID insolar.ID, pn insolar.PulseNumber, rawRecords [][]byte) {
	inslog := inslogger.FromContext(ctx)

	for _, rawRec := range rawRecords {
		rec, err := object.DecodeRecord(rawRec)
		if err != nil {
			inslog.Error(err, "heavyserver: deserialize record failed")
			continue
		}

		virtRec := rec.Record

		id := object.NewRecordIDFromRecord(s.PlatformCryptographyScheme, pn, virtRec)
		err = s.RecordModifier.Set(ctx, *id, rec)
		if err != nil {
			inslog.Error(err, "heavyserver: store record failed")
			continue
		}
	}
}

// Stop successfully stops replication for specified pulse.
func (s *Sync) Stop(ctx context.Context, jetID insolar.ID, pn insolar.PulseNumber) error {
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
func (s *Sync) Reset(ctx context.Context, jetID insolar.ID, pn insolar.PulseNumber) error {
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
