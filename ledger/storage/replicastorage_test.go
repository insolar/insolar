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

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
)

type replicaSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()

	replicaStorage storage.ReplicaStorage

	jetID core.RecordID
}

func NewReplicaSuite() *replicaSuite {
	return &replicaSuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestReplicaStorage(t *testing.T) {
	suite.Run(t, NewReplicaSuite())
}

func (s *replicaSuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())
	s.jetID = core.TODOJetID

	db, cleaner := storagetest.TmpDB(s.ctx, nil, s.T())
	s.cleaner = cleaner
	s.replicaStorage = storage.NewReplicaStorage()

	s.cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		db,
		s.replicaStorage,
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

func (s *replicaSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

func (s *replicaSuite) Test_ReplicatedPulse() {
	// test {Set/Get}HeavySyncedPulse methods pair
	heavyGot0, err := s.replicaStorage.GetHeavySyncedPulse(s.ctx, s.jetID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), core.PulseNumber(0), heavyGot0)

	expectHeavy := core.PulseNumber(100500)
	err = s.replicaStorage.SetHeavySyncedPulse(s.ctx, s.jetID, expectHeavy)
	require.NoError(s.T(), err)

	gotHeavy, err := s.replicaStorage.GetHeavySyncedPulse(s.ctx, s.jetID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), expectHeavy, gotHeavy)
}

func (s *replicaSuite) Test_SyncClientJetPulses() {
	var expectEmpty []core.PulseNumber
	gotEmpty, err := s.replicaStorage.GetSyncClientJetPulses(s.ctx, s.jetID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), expectEmpty, gotEmpty)

	expect := []core.PulseNumber{100, 500, 100500}
	err = s.replicaStorage.SetSyncClientJetPulses(s.ctx, s.jetID, expect)
	require.NoError(s.T(), err)

	got, err := s.replicaStorage.GetSyncClientJetPulses(s.ctx, s.jetID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), expect, got)
}

func (s *replicaSuite) Test_GetAllSyncClientJets() {
	tt := []struct {
		jetID  core.RecordID
		pulses []core.PulseNumber
	}{
		{
			jetID:  testutils.RandomJet(),
			pulses: []core.PulseNumber{100, 500, 100500},
		},
		{
			jetID:  testutils.RandomJet(),
			pulses: []core.PulseNumber{100, 500},
		},
		{
			jetID: testutils.RandomJet(),
		},
		{
			jetID:  testutils.RandomJet(),
			pulses: []core.PulseNumber{100500},
		},
	}

	for _, tCase := range tt {
		err := s.replicaStorage.SetSyncClientJetPulses(s.ctx, tCase.jetID, tCase.pulses)
		require.NoError(s.T(), err)
	}

	gotJets, err := s.replicaStorage.GetAllNonEmptySyncClientJets(s.ctx)
	require.NoError(s.T(), err)

	for i, tCase := range tt {
		gotPulses, ok := gotJets[tCase.jetID]
		if tCase.pulses == nil {
			assert.Falsef(s.T(), ok, "jet should not present jetID=%v", tCase.jetID)
		} else {
			require.Truef(s.T(), ok, "jet should  present jetID=%v", tCase.jetID)
			assert.Equalf(s.T(), tCase.pulses, gotPulses, "pulses not found for jet number %v: %v", i, tCase.jetID)
		}
	}

	gotJets, err = s.replicaStorage.GetAllSyncClientJets(s.ctx)
	require.NoError(s.T(), err)

	for i, tCase := range tt {
		gotPulses, ok := gotJets[tCase.jetID]
		require.Truef(s.T(), ok, "jet should  present jetID=%v", tCase.jetID)
		assert.Equalf(s.T(), tCase.pulses, gotPulses, "pulses not found for jet number %v: %v", i, tCase.jetID)
	}
}
