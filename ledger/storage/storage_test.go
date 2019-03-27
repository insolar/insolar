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

package storage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/db"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

type storageSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()
	db      storage.DBContext

	objectStorage storage.ObjectStorage
	dropModifier  drop.Modifier
	dropAccessor  drop.Accessor
	pulseTracker  storage.PulseTracker

	jetID insolar.ID
}

func NewStorageSuite() *storageSuite {
	return &storageSuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestStorage(t *testing.T) {
	suite.Run(t, NewStorageSuite())
}

func (s *storageSuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())

	tmpDB, cleaner := storagetest.TmpDB(s.ctx, s.T())
	s.db = tmpDB
	s.cleaner = cleaner

	s.objectStorage = storage.NewObjectStorage()

	storageDB := db.NewDBWithBadger(tmpDB.GetBadgerDB())
	dropStorage := drop.NewStorageDB(storageDB)
	s.dropAccessor = dropStorage
	s.dropModifier = dropStorage

	s.pulseTracker = storage.NewPulseTracker()
	s.jetID = testutils.RandomJet()

	s.cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		s.db,
		db.NewMemoryMockDB(),
		s.objectStorage,
		s.dropModifier,
		s.dropAccessor,
		s.pulseTracker,
	)

	err := s.cm.Init(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager init failed", err)
	}
	err = s.cm.Start(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager start failed", err)
	}
}

func (s *storageSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

func (s *storageSuite) TestDB_GetRecordNotFound() {
	rec, err := s.objectStorage.GetRecord(s.ctx, s.jetID, &insolar.ID{})
	assert.Equal(s.T(), err, insolar.ErrNotFound)
	assert.Nil(s.T(), rec)
}

func (s *storageSuite) TestDB_SetRecord() {
	rec := &object.RequestRecord{}
	gotRef, err := s.objectStorage.SetRecord(s.ctx, s.jetID, insolar.GenesisPulse.PulseNumber, rec)
	assert.Nil(s.T(), err)

	gotRec, err := s.objectStorage.GetRecord(s.ctx, s.jetID, gotRef)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), rec, gotRec)

	_, err = s.objectStorage.SetRecord(s.ctx, s.jetID, insolar.GenesisPulse.PulseNumber, rec)
	assert.Equalf(s.T(), err, storage.ErrOverride, "records override should be forbidden")
}

func (s *storageSuite) TestDB_SetObjectIndex_ReturnsNotFoundIfNoIndex() {
	idx, err := s.objectStorage.GetObjectIndex(s.ctx, s.jetID, insolar.NewID(0, hexhash("5000")))
	assert.Equal(s.T(), insolar.ErrNotFound, err)
	assert.Nil(s.T(), idx)
}

func (s *storageSuite) TestDB_SetObjectIndex_StoresCorrectDataInStorage() {
	idx := object.Lifeline{
		LatestState: insolar.NewID(0, hexhash("20")),
	}
	zeroid := insolar.NewID(0, hexhash(""))
	err := s.objectStorage.SetObjectIndex(s.ctx, s.jetID, zeroid, &idx)
	assert.Nil(s.T(), err)

	storedIndex, err := s.objectStorage.GetObjectIndex(s.ctx, s.jetID, zeroid)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), *storedIndex, idx)
}

func (s *storageSuite) TestDB_SetObjectIndex_SaveLastUpdate() {
	// Arrange
	jetID := testutils.RandomJet()

	idx := object.Lifeline{
		LatestState:  insolar.NewID(0, hexhash("20")),
		LatestUpdate: 1239,
	}
	zeroid := insolar.NewID(0, hexhash(""))

	// Act
	err := s.objectStorage.SetObjectIndex(s.ctx, jetID, zeroid, &idx)
	assert.Nil(s.T(), err)

	// Assert
	storedIndex, err := s.objectStorage.GetObjectIndex(s.ctx, jetID, zeroid)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), *storedIndex, idx)
	assert.Equal(s.T(), 1239, int(idx.LatestUpdate))
}

func (s *storageSuite) TestDB_GetDrop_ReturnsNotFoundIfNoDrop() {
	d, err := s.dropAccessor.ForPulse(s.ctx, insolar.JetID(testutils.RandomJet()), 1)
	assert.Equal(s.T(), err, db.ErrNotFound)
	assert.Equal(s.T(), drop.Drop{}, d)
}

func (s *storageSuite) TestDB_SetDrop() {
	jetID := *insolar.NewJetID(0, nil)
	drop42 := drop.Drop{
		Pulse: 42,
		Hash:  []byte{0xFF},
		JetID: jetID,
	}
	// FIXME: should work with random jet
	// jetID := testutils.RandomJet()
	err := s.dropModifier.Set(s.ctx, drop42)
	assert.NoError(s.T(), err)

	got, err := s.dropAccessor.ForPulse(s.ctx, jetID, 42)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), got, drop42)
}

func (s *storageSuite) TestDB_AddPulse() {
	pulse42 := insolar.Pulse{PulseNumber: 42, Entropy: insolar.Entropy{1, 2, 3}}
	err := s.pulseTracker.AddPulse(s.ctx, pulse42)
	require.NoError(s.T(), err)

	latestPulse, err := s.pulseTracker.GetLatestPulse(s.ctx)
	assert.Equal(s.T(), insolar.PulseNumber(42), latestPulse.Pulse.PulseNumber)

	pulse, err := s.pulseTracker.GetPulse(s.ctx, latestPulse.Pulse.PulseNumber)
	require.NoError(s.T(), err)

	prevPulse, err := s.pulseTracker.GetPulse(s.ctx, *latestPulse.Prev)
	require.NoError(s.T(), err)

	prevPN := insolar.PulseNumber(insolar.FirstPulseNumber)
	expectPulse := storage.Pulse{
		Prev:         &prevPN,
		Pulse:        pulse42,
		SerialNumber: prevPulse.SerialNumber + 1,
	}
	assert.Equal(s.T(), expectPulse, *pulse)
}

func TestDB_Close(t *testing.T) {
	ctx := inslogger.TestContext(t)
	tmpDB, cleaner := storagetest.TmpDB(ctx, t)

	jetID := testutils.RandomJet()

	os := storage.NewObjectStorage()
	storageDB := db.NewDBWithBadger(tmpDB.GetBadgerDB())
	ds := drop.NewStorageDB(storageDB)

	cm := &component.Manager{}
	cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		tmpDB,
		db.NewMemoryMockDB(),
		os,
		ds,
	)
	err := cm.Init(ctx)
	if err != nil {
		t.Error("ComponentManager init failed", err)
	}
	err = cm.Start(ctx)
	if err != nil {
		t.Error("ComponentManager start failed", err)
	}

	err = cm.Stop(ctx)
	if err != nil {
		t.Error("ComponentManager stop failed", err)
	}

	cleaner()

	rec, err := os.GetRecord(ctx, jetID, &insolar.ID{})
	assert.Nil(t, rec)
	assert.Equal(t, err, storage.ErrClosed)

	rec = &object.RequestRecord{}
	gotRef, err := os.SetRecord(ctx, jetID, insolar.GenesisPulse.PulseNumber, rec)
	assert.Nil(t, gotRef)
	assert.Equal(t, err, storage.ErrClosed)
}
