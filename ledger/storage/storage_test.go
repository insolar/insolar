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

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/ledger/storage/pulse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/storage"
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

	indexAccessor object.IndexAccessor
	indexModifier object.IndexModifier

	dropModifier drop.Modifier
	dropAccessor drop.Accessor

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

	tmpDB, _, cleaner := storagetest.TmpDB(s.ctx, s.T())
	s.db = tmpDB
	s.cleaner = cleaner

	idxStor := object.NewIndexMemory()
	s.indexModifier = idxStor
	s.indexAccessor = idxStor

	storageDB := store.NewMemoryMockDB()
	dropStorage := drop.NewDB(storageDB)
	s.dropAccessor = dropStorage
	s.dropModifier = dropStorage
	s.jetID = testutils.RandomJet()

	s.cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		s.db,
		store.NewMemoryMockDB(),
		idxStor,
		s.dropModifier,
		s.dropAccessor,
		pulse.NewStorageMem(),
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

func (s *storageSuite) TestDB_SetObjectIndex_ReturnsNotFoundIfNoIndex() {
	idx, err := s.indexAccessor.ForID(s.ctx, *insolar.NewID(0, hexhash("5000")))
	assert.Equal(s.T(), object.ErrIndexNotFound, err)
	assert.Equal(s.T(), idx, object.Lifeline{})
}

func (s *storageSuite) TestDB_SetObjectIndex_StoresCorrectDataInStorage() {
	idx := object.Lifeline{
		LatestState: insolar.NewID(0, hexhash("20")),
		JetID:       insolar.JetID(s.jetID),
		Delegates:   map[insolar.Reference]insolar.Reference{},
	}
	zeroid := insolar.NewID(0, hexhash(""))
	err := s.indexModifier.Set(s.ctx, *zeroid, idx)
	assert.Nil(s.T(), err)

	storedIndex, err := s.indexAccessor.ForID(s.ctx, *zeroid)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), storedIndex, idx)
}

func (s *storageSuite) TestDB_SetObjectIndex_SaveLastUpdate() {
	// Arrange
	jetID := testutils.RandomJet()

	idx := object.Lifeline{
		LatestState:  insolar.NewID(0, hexhash("20")),
		LatestUpdate: 1239,
		JetID:        insolar.JetID(jetID),
		Delegates:    map[insolar.Reference]insolar.Reference{},
	}
	zeroid := insolar.NewID(0, hexhash(""))

	// Act
	err := s.indexModifier.Set(s.ctx, *zeroid, idx)
	assert.Nil(s.T(), err)

	// Assert
	storedIndex, err := s.indexAccessor.ForID(s.ctx, *zeroid)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), storedIndex, idx)
	assert.Equal(s.T(), 1239, int(idx.LatestUpdate))
}

func (s *storageSuite) TestDB_GetDrop_ReturnsNotFoundIfNoDrop() {
	d, err := s.dropAccessor.ForPulse(s.ctx, insolar.JetID(testutils.RandomJet()), 1)
	assert.Equal(s.T(), err, store.ErrNotFound)
	assert.Equal(s.T(), drop.Drop{}, d)
}

func (s *storageSuite) TestDB_SetDrop() {
	jetID := gen.JetID()
	drop42 := drop.Drop{
		Pulse: 42,
		Hash:  []byte{0xFF},
		JetID: jetID,
	}
	err := s.dropModifier.Set(s.ctx, drop42)
	assert.NoError(s.T(), err)

	got, err := s.dropAccessor.ForPulse(s.ctx, jetID, 42)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), got, drop42)
}
