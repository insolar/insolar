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

package sequence

import (
	"bytes"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"sort"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/genesis"
	"github.com/insolar/insolar/ledger/object"
)

func TestSequencer_Records(t *testing.T) {
	ctx := inslogger.TestContext(t)

	db := store.NewMemoryMockDB()
	storage := object.NewRecordDB(db)
	records := NewSequencer(db, store.ScopeRecord)

	type tempRecord struct {
		id  insolar.ID
		rec record.Material
	}

	const (
		pulse = 10
		skip  = 5
		limit = 10
	)

	var (
		expected   []tempRecord
		unexpected []tempRecord
	)

	f := fuzz.New().Funcs(func(t *tempRecord, c fuzz.Continue) {
		var hash [insolar.RecordIDSize - insolar.RecordHashOffset]byte
		c.Fuzz(&hash)
		t.id = *insolar.NewID(pulse, hash[:])
		t.rec = record.Material{Polymorph: c.Int31()}
	})
	f.NilChance(0)
	f.NumElements(limit, limit)
	f.Fuzz(&expected)

	f = fuzz.New().Funcs(func(t *tempRecord, c fuzz.Continue) {
		var hash [insolar.RecordIDSize - insolar.RecordHashOffset]byte
		c.Fuzz(&hash)
		oneOrZero := rand.Intn(2)
		t.id = *insolar.NewID(insolar.PulseNumber(pulse+oneOrZero*(-10)+(1-oneOrZero)*(10)), hash[:])
		t.rec = record.Material{Polymorph: c.Int31()}
	})
	f.NilChance(0)
	f.NumElements(limit, limit)
	f.Fuzz(&unexpected)

	for i := 0; i < limit; i++ {
		storage.Set(ctx, expected[i].id, expected[i].rec)
		storage.Set(ctx, unexpected[i].id, unexpected[i].rec)
	}

	t.Run("slice length", func(t *testing.T) {
		recs := records.Slice(pulse, skip, pulse+10, limit)
		require.Equal(t, limit-skip, len(recs))
	})

	t.Run("first last", func(t *testing.T) {
		all := append(expected, unexpected...)
		sort.Slice(all, func(i, j int) bool {
			return bytes.Compare(all[i].id.Bytes(), all[j].id.Bytes()) == -1
		})

		require.Equal(t, all[0].id.Bytes(), records.First().Key)
		require.Equal(t, all[len(all)-1].id.Bytes(), records.Last().Key)
	})
}

func TestSequencer_Blobs(t *testing.T) {
	ctx := inslogger.TestContext(t)

	db := store.NewMemoryMockDB()
	storage := blob.NewDB(db)
	blobs := NewSequencer(db, store.ScopeBlob)

	type tempBlob struct {
		id   insolar.ID
		blob blob.Blob
	}

	const (
		pulse = 10
		skip  = 5
		limit = 10
	)

	var (
		expected   []tempBlob
		unexpected []tempBlob
	)

	f := fuzz.New().Funcs(func(t *tempBlob, c fuzz.Continue) {
		var memory [10]byte
		c.Fuzz(&memory)
		var hash [insolar.RecordIDSize - insolar.RecordHashOffset]byte
		c.Fuzz(&hash)
		oneOrZero := rand.Intn(2)
		t.id = *insolar.NewID(insolar.PulseNumber(pulse+oneOrZero*(-10)+(1-oneOrZero)*(10)), hash[:])
		t.blob = blob.Blob{
			JetID: insolar.ZeroJetID,
			Value: memory[:],
		}
	})
	f.NilChance(0)
	f.NumElements(limit, limit)
	f.Fuzz(&expected)

	f = fuzz.New().Funcs(func(t *tempBlob, c fuzz.Continue) {
		var memory [10]byte
		c.Fuzz(&memory)
		var hash [insolar.RecordIDSize - insolar.RecordHashOffset]byte
		c.Fuzz(&hash)
		t.id = *insolar.NewID(pulse, hash[:])
		t.blob = blob.Blob{
			JetID: insolar.ZeroJetID,
			Value: memory[:],
		}
	})
	f.NilChance(0)
	f.NumElements(limit, limit)
	f.Fuzz(&unexpected)

	for i := 0; i < limit; i++ {
		storage.Set(ctx, expected[i].id, expected[i].blob)
		storage.Set(ctx, unexpected[i].id, unexpected[i].blob)
	}

	t.Run("slice length", func(t *testing.T) {
		blbs := blobs.Slice(pulse, skip, pulse+10, limit)
		require.Equal(t, limit-skip, len(blbs))
	})

	t.Run("first last", func(t *testing.T) {
		all := append(expected, unexpected...)
		sort.Slice(all, func(i, j int) bool {
			return bytes.Compare(all[i].id.Bytes(), all[j].id.Bytes()) == -1
		})

		require.Equal(t, all[0].id.Bytes(), blobs.First().Key)
		require.Equal(t, all[len(all)-1].id.Bytes(), blobs.Last().Key)
	})
}

func TestSequencer_Drops(t *testing.T) {
	ctx := inslogger.TestContext(t)

	db := store.NewMemoryMockDB()
	storage := drop.NewDB(db)
	drops := NewSequencer(db, store.ScopeJetDrop)

	type tempDrop struct {
		drop drop.Drop
	}

	const (
		pulse = 10
		skip  = 5
		limit = 10
	)

	var (
		expected   []tempDrop
		unexpected []tempDrop
	)

	f := fuzz.New().Funcs(func(t *tempDrop, c fuzz.Continue) {
		var hash [insolar.RecordIDSize - insolar.RecordHashOffset]byte
		c.Fuzz(&hash)
		t.drop = drop.Drop{
			Pulse: pulse,
			JetID: gen.JetID(),
		}
	})
	f.NilChance(0)
	f.NumElements(limit, limit)
	f.Fuzz(&expected)

	f = fuzz.New().Funcs(func(t *tempDrop, c fuzz.Continue) {
		var hash [insolar.RecordIDSize - insolar.RecordHashOffset]byte
		c.Fuzz(&hash)
		oneOrZero := rand.Intn(2)
		t.drop = drop.Drop{
			Pulse: insolar.PulseNumber(pulse + oneOrZero*(-10) + (1-oneOrZero)*(10)),
			JetID: gen.JetID(),
		}
	})
	f.NilChance(0)
	f.NumElements(limit, limit)
	f.Fuzz(&unexpected)

	for i := 0; i < limit; i++ {
		storage.Set(ctx, expected[i].drop)
		storage.Set(ctx, unexpected[i].drop)
	}

	t.Run("slice length", func(t *testing.T) {
		drps := drops.Slice(pulse, skip, pulse+10, limit)
		require.Equal(t, limit-skip, len(drps))
	})

	t.Run("first last", func(t *testing.T) {
		all := append(expected, unexpected...)
		sort.Slice(all, func(i, j int) bool {
			return bytes.Compare(
				append(all[i].drop.Pulse.Bytes(), all[i].drop.JetID.Prefix()...),
				append(all[j].drop.Pulse.Bytes(), all[j].drop.JetID.Prefix()...),
			) == -1
		})

		require.Equal(t, append(all[0].drop.Pulse.Bytes(), all[0].drop.JetID.Prefix()...), drops.First().Key)
		require.Equal(t, append(all[len(all)-1].drop.Pulse.Bytes(), all[len(all)-1].drop.JetID.Prefix()...), drops.Last().Key)
	})
}

func TestSequencer_Genesis(t *testing.T) {
	ctx := inslogger.TestContext(t)

	db := store.NewMemoryMockDB()
	pulseDB := pulse.NewDB(db)
	dropDB := drop.NewDB(db)
	recordDB := object.NewRecordDB(db)
	lifelineDB := object.NewIndexDB(db)
	genesisCreator := genesis.BaseRecord{
		DB:                    db,
		PulseAccessor:         pulseDB,
		PulseAppender:         pulseDB,
		DropModifier:          dropDB,
		RecordModifier:        recordDB,
		IndexLifelineModifier: lifelineDB,
	}

	pulses := NewSequencer(db, store.ScopePulse)
	drops := NewSequencer(db, store.ScopeJetDrop)
	records := NewSequencer(db, store.ScopeRecord)
	lifelines := NewSequencer(db, store.ScopeIndex)
	lastKnownIndex := NewSequencer(db, store.ScopeLastKnownIndexPN)
	baseRecord := NewSequencer(db, store.ScopeGenesis)

	genesisCreator.Create(ctx)
	genesisCreator.Done(ctx)

	var (
		pulse     = insolar.GenesisPulse.PulseNumber
		nextPulse = pulse + 10
	)

	t.Run("check pulses", func(t *testing.T) {
		actualPulses := pulses.Slice(pulse, 0, nextPulse, 1)
		require.Equal(t, 1, pulses.Len(pulse))
		require.Equal(t, 1, len(actualPulses))
		require.Equal(t, pulse, insolar.NewPulseNumber(actualPulses[0].Key))
	})

	t.Run("check drops", func(t *testing.T) {
		actualDrops := drops.Slice(pulse, 0, nextPulse, 1)
		require.Equal(t, 1, drops.Len(pulse))
		require.Equal(t, 1, len(actualDrops))
		drop, err := drop.Decode(actualDrops[0].Value)
		require.NoError(t, err)
		require.Equal(t, pulse, drop.Pulse)
		require.Equal(t, insolar.ZeroJetID, drop.JetID)
	})

	t.Run("check records and filaments", func(t *testing.T) {
		actualRecords := records.Slice(pulse, 0, nextPulse, 1)
		require.Equal(t, 1, records.Len(pulse))
		require.Equal(t, 1, len(actualRecords))
		rec := record.Material{}
		err := rec.Unmarshal(actualRecords[0].Value)
		require.NoError(t, err)
		require.Equal(t, insolar.ZeroJetID, rec.JetID)
	})

	t.Run("check filaments", func(t *testing.T) {
		actualFilaments := lifelines.Slice(pulse, 0, nextPulse, 1)
		require.Equal(t, 1, lifelines.Len(pulse))
		require.Equal(t, 1, len(actualFilaments))
		filament := object.FilamentIndex{}
		err := filament.Unmarshal(actualFilaments[0].Value)
		require.NoError(t, err)
		require.Equal(t, insolar.GenesisRecord.ID(), filament.ObjID)

		actualLastKnownIndex := lastKnownIndex.Slice(pulse, 0, nextPulse, 1)
		require.Equal(t, 1, lastKnownIndex.Len(pulse))
		require.Equal(t, 1, len(actualLastKnownIndex))
		actualLastPulse := insolar.NewPulseNumber(actualLastKnownIndex[0].Value)
		require.Equal(t, pulse, actualLastPulse)
	})

	t.Run("check baserecord", func(t *testing.T) {
		actualBaseRecords := baseRecord.Slice(pulse, 0, nextPulse, 1)
		require.Equal(t, 1, baseRecord.Len(pulse))
		require.Equal(t, 1, len(actualBaseRecords))
	})

}

func TestSequencer_Pulses(t *testing.T) {
	ctx := inslogger.TestContext(t)

	db := store.NewMemoryMockDB()
	storage := pulse.NewDB(db)
	records := NewSequencer(db, store.ScopePulse)

	const (
		pulse = 10
		limit = 10
		total = 30
	)

	var all []insolar.Pulse

	f := fuzz.New().Funcs(func(t *insolar.Pulse, c fuzz.Continue) {
		t.PulseNumber = pulse
	})
	f.NilChance(0)
	f.NumElements(30, 30)
	f.Fuzz(&all)

	for i := 0; i < total; i++ {
		storage.Append(ctx, insolar.Pulse{PulseNumber: insolar.PulseNumber(pulse + 10*i)})
	}

	recs := records.Slice(pulse, 0, math.MaxUint32, limit)
	require.Equal(t, limit, len(recs))
}

func TestSequencer_Upsert(t *testing.T) {
	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	db, err := store.NewBadgerDB(tmpdir)
	defer db.Stop(ctx)
	require.NoError(t, err)
	records := NewSequencer(db, store.ScopeRecord)

	type tempRecord struct {
		id  insolar.ID
		rec record.Material
	}
	const (
		pulse = 10
		total = 3
	)

	var all []tempRecord

	f := fuzz.New().Funcs(func(t *tempRecord, c fuzz.Continue) {
		var hash [insolar.RecordIDSize - insolar.RecordHashOffset]byte
		c.Fuzz(&hash)
		t.id = *insolar.NewID(pulse, hash[:])
		t.rec = record.Material{Polymorph: c.Int31()}
	})
	f.NilChance(0)
	f.NumElements(total, total)
	f.Fuzz(&all)

	require.Equal(t, 0, records.Len(pulse))

	items := []Item{}
	for _, rec := range all {
		val, _ := rec.rec.Marshal()
		items = append(items, Item{Key: rec.id.Bytes(), Value: val})
	}

	records.Upsert(items)
	require.Equal(t, len(items), records.Len(pulse))

	seq := records.Slice(pulse, 0, math.MaxUint32, total)
	for _, item := range seq {
		id := insolar.ID{}
		copy(id[:], item.Key)
		_, err := db.Get(&dbKey{id})
		require.NoError(t, err)
	}
}

type dbKey struct {
	id insolar.ID
}

func (k *dbKey) Scope() store.Scope {
	return store.ScopeRecord
}

func (k *dbKey) ID() []byte {
	return k.id[:]
}
