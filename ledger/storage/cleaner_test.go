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
	"sort"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type cleanerSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()

	objectStorage  storage.ObjectStorage
	storageCleaner storage.Cleaner

	jetID core.RecordID
}

func NewCleanerSuite() *cleanerSuite {
	return &cleanerSuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestCleaner(t *testing.T) {
	suite.Run(t, NewCleanerSuite())
}

func (s *cleanerSuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())
	// TODO: just use two cases: zero and non zero jetID
	s.jetID = core.TODOJetID

	db, cleaner := storagetest.TmpDB(s.ctx, s.T())
	s.cleaner = cleaner

	s.objectStorage = storage.NewObjectStorage()
	s.storageCleaner = storage.NewCleaner()

	s.cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		db,
		s.objectStorage,
		s.storageCleaner,
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

func (s *cleanerSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

func (s *cleanerSuite) Test_RemoveJetIndexesUntil_Basic() {
	removeJetIndexesUntil(s, false)
}

func (s *cleanerSuite) Test_RemoveJetIndexesUntil_WithSkips() {
	removeJetIndexesUntil(s, true)
}

func removeJetIndexesUntil(s *cleanerSuite, skip bool) {
	// if we operate on zero jetID
	var expectLeftIDs []core.RecordID
	err := s.objectStorage.IterateIndexIDs(s.ctx, s.jetID, func(id core.RecordID) error {
		if id.Pulse() == core.FirstPulseNumber {
			expectLeftIDs = append(expectLeftIDs, id)
		}
		return nil
	})
	require.NoError(s.T(), err)

	pulsesCount := 10
	untilIdx := pulsesCount / 2
	var until core.PulseNumber

	var pulses []core.PulseNumber
	expectedRmCount := 0
	for i := 0; i < pulsesCount; i++ {
		pn := core.FirstPulseNumber + core.PulseNumber(i)
		if i == untilIdx {
			until = pn
			if skip {
				// skip index saving with 'until' pulse (corner case)
				continue
			}
		}
		pulses = append(pulses, pn)
		objID := testutils.RandomID()
		copy(objID[:core.PulseNumberSize], pn.Bytes())
		err := s.objectStorage.SetObjectIndex(s.ctx, s.jetID, &objID, &index.ObjectLifeline{
			State:       record.StateActivation,
			LatestState: &objID,
		})
		require.NoError(s.T(), err)
		if (pn == core.FirstPulseNumber) || (i >= untilIdx) {
			expectLeftIDs = append(expectLeftIDs, objID)
		} else {
			expectedRmCount += 1
		}
	}

	rmcount, err := s.storageCleaner.RemoveJetIndexesUntil(s.ctx, s.jetID, until, nil)
	require.NoError(s.T(), err)

	var foundIDs []core.RecordID
	err = s.objectStorage.IterateIndexIDs(s.ctx, s.jetID, func(id core.RecordID) error {
		foundIDs = append(foundIDs, id)
		return nil
	})
	require.NoError(s.T(), err)

	assert.Equal(s.T(), int64(expectedRmCount), rmcount.Removed)
	assert.Equalf(s.T(), sortIDS(expectLeftIDs), sortIDS(foundIDs),
		"expected keys and found indexes, doesn't match, jetID=%v", s.jetID.DebugString())
}

func sortIDS(ids []core.RecordID) []core.RecordID {
	sort.Slice(ids, func(i, j int) bool {
		return bytes.Compare(ids[i][:], ids[j][:]) < 0
	})
	return ids
}
