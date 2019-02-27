/*
 *    Copyright 2019 Insolar Technologies
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

package storage_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/insolar/insolar/component"
	base58 "github.com/jbenet/go-base58"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type storageSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()
	db      storage.DBContext

	objectStorage storage.ObjectStorage
	dropStorage   storage.DropStorage
	pulseTracker  storage.PulseTracker

	jetID core.RecordID
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

	db, cleaner := storagetest.TmpDB(s.ctx, s.T())
	s.db = db
	s.cleaner = cleaner

	s.objectStorage = storage.NewObjectStorage()
	s.dropStorage = storage.NewDropStorage(10)
	s.pulseTracker = storage.NewPulseTracker()
	s.jetID = testutils.RandomJet()

	s.cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		s.db,
		s.objectStorage,
		s.dropStorage,
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
	rec, err := s.objectStorage.GetRecord(s.ctx, s.jetID, &core.RecordID{})
	assert.Equal(s.T(), err, core.ErrNotFound)
	assert.Nil(s.T(), rec)
}

func (s *storageSuite) TestDB_SetRecord() {
	rec := &record.RequestRecord{}
	gotRef, err := s.objectStorage.SetRecord(s.ctx, s.jetID, core.GenesisPulse.PulseNumber, rec)
	assert.Nil(s.T(), err)

	gotRec, err := s.objectStorage.GetRecord(s.ctx, s.jetID, gotRef)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), rec, gotRec)

	_, err = s.objectStorage.SetRecord(s.ctx, s.jetID, core.GenesisPulse.PulseNumber, rec)
	assert.Equalf(s.T(), err, storage.ErrOverride, "records override should be forbidden")
}

func (s *storageSuite) TestDB_SetObjectIndex_ReturnsNotFoundIfNoIndex() {
	idx, err := s.objectStorage.GetObjectIndex(s.ctx, s.jetID, core.NewRecordID(0, hexhash("5000")), false)
	assert.Equal(s.T(), core.ErrNotFound, err)
	assert.Nil(s.T(), idx)
}

func (s *storageSuite) TestDB_SetObjectIndex_StoresCorrectDataInStorage() {
	idx := index.ObjectLifeline{
		LatestState: core.NewRecordID(0, hexhash("20")),
	}
	zeroid := core.NewRecordID(0, hexhash(""))
	err := s.objectStorage.SetObjectIndex(s.ctx, s.jetID, zeroid, &idx)
	assert.Nil(s.T(), err)

	storedIndex, err := s.objectStorage.GetObjectIndex(s.ctx, s.jetID, zeroid, false)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), *storedIndex, idx)
}

func (s *storageSuite) TestDB_SetObjectIndex_SaveLastUpdate() {
	// Arrange
	jetID := testutils.RandomJet()

	idx := index.ObjectLifeline{
		LatestState:  core.NewRecordID(0, hexhash("20")),
		LatestUpdate: 1239,
	}
	zeroid := core.NewRecordID(0, hexhash(""))

	// Act
	err := s.objectStorage.SetObjectIndex(s.ctx, jetID, zeroid, &idx)
	assert.Nil(s.T(), err)

	// Assert
	storedIndex, err := s.objectStorage.GetObjectIndex(s.ctx, jetID, zeroid, false)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), *storedIndex, idx)
	assert.Equal(s.T(), 1239, int(idx.LatestUpdate))
}

func (s *storageSuite) TestDB_GetDrop_ReturnsNotFoundIfNoDrop() {
	drop, err := s.dropStorage.GetDrop(s.ctx, testutils.RandomJet(), 1)
	assert.Equal(s.T(), err, core.ErrNotFound)
	assert.Nil(s.T(), drop)
}

func (s *storageSuite) TestDB_CreateDrop() {
	// FIXME: should work with random jet
	// jetID := testutils.RandomJet()
	jetID := *jet.NewID(0, nil)

	pulse := core.PulseNumber(core.FirstPulseNumber + 10)
	err := s.pulseTracker.AddPulse(
		s.ctx,
		core.Pulse{
			PulseNumber: pulse,
			Entropy:     core.Entropy{1, 2, 3},
		},
	)
	cs := platformpolicy.NewPlatformCryptographyScheme()

	msgCount := 3
	for i := 1; i < 1+msgCount; i++ {
		setRecordMessage := message.SetRecord{
			Record: record.SerializeRecord(&record.CodeRecord{
				Code: record.CalculateIDForBlob(cs, pulse, []byte{byte(i)}),
			}),
		}
		err = s.objectStorage.SetMessage(s.ctx, jetID, pulse, &setRecordMessage)
		require.NoError(s.T(), err)
		_, err = s.objectStorage.SetBlob(s.ctx, jetID, pulse, []byte{byte(i)})
		require.NoError(s.T(), err)
	}

	drop, messages, dropSize, err := s.dropStorage.CreateDrop(s.ctx, jetID, pulse, []byte{4, 5, 6})
	require.NoError(s.T(), err)
	require.NotEqual(s.T(), 0, dropSize)
	// TODO: messages collection was disabled in ab46d01, validation is not active ATM
	require.Equal(s.T(), 0, len(messages))
	require.Equal(s.T(), pulse, drop.Pulse)
	require.Equal(s.T(), "2aCdao6DhZSWQNTrtrxJW7QQZRb6UJ1ssRi9cg", base58.Encode(drop.Hash))

	for _, rawMessage := range messages {
		formatedMessage, err := message.Deserialize(bytes.NewBuffer(rawMessage))
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), core.TypeSetRecord, formatedMessage.Type())
	}
}

func (s *storageSuite) TestDB_SetDrop() {
	drop42 := jet.JetDrop{
		Pulse: 42,
		Hash:  []byte{0xFF},
	}
	// FIXME: should work with random jet
	// jetID := testutils.RandomJet()
	jetID := *jet.NewID(0, nil)
	err := s.dropStorage.SetDrop(s.ctx, jetID, &drop42)
	assert.NoError(s.T(), err)

	got, err := s.dropStorage.GetDrop(s.ctx, jetID, 42)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), *got, drop42)
}

func (s *storageSuite) TestDB_AddPulse() {
	pulse42 := core.Pulse{PulseNumber: 42, Entropy: core.Entropy{1, 2, 3}}
	err := s.pulseTracker.AddPulse(s.ctx, pulse42)
	require.NoError(s.T(), err)

	latestPulse, err := s.pulseTracker.GetLatestPulse(s.ctx)
	assert.Equal(s.T(), core.PulseNumber(42), latestPulse.Pulse.PulseNumber)

	pulse, err := s.pulseTracker.GetPulse(s.ctx, latestPulse.Pulse.PulseNumber)
	require.NoError(s.T(), err)

	prevPulse, err := s.pulseTracker.GetPulse(s.ctx, *latestPulse.Prev)
	require.NoError(s.T(), err)

	prevPN := core.PulseNumber(core.FirstPulseNumber)
	expectPulse := storage.Pulse{
		Prev:         &prevPN,
		Pulse:        pulse42,
		SerialNumber: prevPulse.SerialNumber + 1,
	}
	assert.Equal(s.T(), expectPulse, *pulse)
}

func TestDB_Close(t *testing.T) {
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)

	jetID := testutils.RandomJet()

	os := storage.NewObjectStorage()
	ds := storage.NewDropStorage(10)

	cm := &component.Manager{}
	cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		db,
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

	rec, err := os.GetRecord(ctx, jetID, &core.RecordID{})
	assert.Nil(t, rec)
	assert.Equal(t, err, storage.ErrClosed)

	rec = &record.RequestRecord{}
	gotRef, err := os.SetRecord(ctx, jetID, core.GenesisPulse.PulseNumber, rec)
	assert.Nil(t, gotRef)
	assert.Equal(t, err, storage.ErrClosed)
}
