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
	IndexModifier              object.IndexModifier               `inject:""`
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

// StoreIndexes stores recieved key/value pairs for indices at heavy storage.
func (s *Sync) StoreIndexes(ctx context.Context, jet insolar.ID, pn insolar.PulseNumber, rawIndexes map[insolar.ID][]byte) error {
	for id, rwi := range rawIndexes {
		idx, err := object.DecodeIndex(rwi)
		if err != nil {
			inslogger.FromContext(ctx).Error(err)
			continue
		}

		err = s.IndexModifier.Set(ctx, id, idx)
		if err != nil {
			return errors.Wrapf(err, "heavyserver: index storing failed")
		}
	}

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
	inslog := inslogger.FromContext(ctx)

	for _, rwb := range rawBlobs {
		b, err := blob.Decode(rwb)
		if err != nil {
			inslog.Error(err, "heavyserver: deserialize blob failed")
			continue
		}

		blobID := object.CalculateIDForBlob(s.PlatformCryptographyScheme, pn, b.Value)

		err = s.BlobModifier.Set(ctx, *blobID, *b)
		if err != nil {
			inslog.Error(err, "heavyserver: blob storing failed")
			continue
		}
	}
	return nil
}

// StoreRecords stores recieved records at heavy storage.
func (s *Sync) StoreRecords(ctx context.Context, jetID insolar.ID, pn insolar.PulseNumber, rawRecords [][]byte) {
	inslog := inslogger.FromContext(ctx)

	for _, rawRec := range rawRecords {
		rec, err := object.DecodeMaterial(rawRec)
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
