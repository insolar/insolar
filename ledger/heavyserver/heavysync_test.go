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

package heavyserver

import (
	"context"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type heavysyncSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()
	db      storage.DBContext

	pulseTracker   storage.PulseTracker
	replicaStorage storage.ReplicaStorage

	sync *Sync
}

func NewHeavysyncSuite() *heavysyncSuite {
	return &heavysyncSuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestHeavysync(t *testing.T) {
	suite.Run(t, NewHeavysyncSuite())
}

func (s *heavysyncSuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())

	db, cleaner := storagetest.TmpDB(s.ctx, s.T())

	s.db = db
	s.cleaner = cleaner
	s.pulseTracker = storage.NewPulseTracker()
	s.replicaStorage = storage.NewReplicaStorage()

	s.cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		s.db,
		s.pulseTracker,
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

func (s *heavysyncSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

func (s *heavysyncSuite) TestHeavy_SyncBasic() {
	var err error
	var pnum core.PulseNumber
	kvalues := []core.KV{
		{K: []byte("100"), V: []byte("500")},
	}

	// TODO: call every case in subtest
	jetID := testutils.RandomJet()

	sync := NewSync(s.db)
	sync.ReplicaStorage = s.replicaStorage
	err = sync.Start(s.ctx, jetID, pnum)
	require.Error(s.T(), err, "start with zero pulse")

	err = sync.Store(s.ctx, jetID, pnum, kvalues)
	require.Error(s.T(), err, "store values on non started sync")

	err = sync.Stop(s.ctx, jetID, pnum)
	require.Error(s.T(), err, "stop on non started sync")

	pnum = 5
	err = sync.Start(s.ctx, jetID, pnum)
	require.Error(s.T(), err, "last synced pulse is less when 'first pulse number'")

	pnum = core.FirstPulseNumber
	err = sync.Start(s.ctx, jetID, pnum)
	require.Error(s.T(), err, "start from first pulse on empty storage")

	pnum = core.FirstPulseNumber + 1
	err = sync.Start(s.ctx, jetID, pnum)
	require.NoError(s.T(), err, "start sync on empty heavy jet with non first pulse number")

	err = sync.Start(s.ctx, jetID, pnum)
	require.Error(s.T(), err, "double start")

	pnumNext := pnum + 1
	err = sync.Start(s.ctx, jetID, pnumNext)
	require.Error(s.T(), err, "start next pulse sync when previous not end")

	// stop previous
	err = sync.Stop(s.ctx, jetID, pnum)
	require.NoError(s.T(), err)

	// start sparse next
	pnumNextPlus := pnumNext + 1
	err = sync.Start(s.ctx, jetID, pnumNextPlus)
	require.NoError(s.T(), err, "sparse sync is ok")
	err = sync.Stop(s.ctx, jetID, pnumNextPlus)
	require.NoError(s.T(), err)

	// prepare pulse helper
	preparepulse := func(pn core.PulseNumber) {
		pulse := core.Pulse{PulseNumber: pn}
		err = s.pulseTracker.AddPulse(s.ctx, pulse)
		require.NoError(s.T(), err)
	}
	pnum = pnumNextPlus + 1
	pnumNext = pnum + 1
	preparepulse(pnum)
	preparepulse(pnumNext) // should set correct next for previous pulse

	err = sync.Start(s.ctx, jetID, pnumNext)
	require.NoError(s.T(), err, "start next pulse")

	err = sync.Store(s.ctx, jetID, pnumNextPlus, kvalues)
	require.Error(s.T(), err, "store from other pulse at the same jet")

	err = sync.Stop(s.ctx, jetID, pnumNextPlus)
	require.Error(s.T(), err, "stop from other pulse at the same jet")

	err = sync.Store(s.ctx, jetID, pnumNext, kvalues)
	require.NoError(s.T(), err, "store on current range")
	err = sync.Store(s.ctx, jetID, pnumNext, kvalues)
	require.NoError(s.T(), err, "store the same on current range")
	err = sync.Stop(s.ctx, jetID, pnumNext)
	require.NoError(s.T(), err, "stop current range")

	preparepulse(pnumNextPlus) // should set corret next for previous pulse
	sync = NewSync(s.db)
	sync.ReplicaStorage = s.replicaStorage
	err = sync.Start(s.ctx, jetID, pnumNextPlus)
	require.NoError(s.T(), err, "start next+1 range on new sync instance (checkpoint check)")
	err = sync.Store(s.ctx, jetID, pnumNextPlus, kvalues)
	require.NoError(s.T(), err, "store next+1 pulse")
	err = sync.Stop(s.ctx, jetID, pnumNextPlus)
	require.NoError(s.T(), err, "stop next+1 range on new sync instance")
}

func (s *heavysyncSuite) TestHeavy_SyncByJet() {
	var err error
	var pnum core.PulseNumber
	kvalues1 := []core.KV{
		{K: []byte("1_11"), V: []byte("1_12")},
	}
	kvalues2 := []core.KV{
		{K: []byte("2_21"), V: []byte("2_22")},
	}

	// TODO: call every case in subtest
	jetID1 := testutils.RandomJet()
	jetID2 := jetID1
	// flip first bit of jetID2 for different prefix
	lastidx := len(jetID1) - 1
	jetID2[lastidx] ^= 0xFF

	// prepare pulse helper
	preparepulse := func(pn core.PulseNumber) {
		pulse := core.Pulse{PulseNumber: pn}
		err = s.pulseTracker.AddPulse(s.ctx, pulse)
		require.NoError(s.T(), err)
	}

	sync := NewSync(s.db)
	sync.ReplicaStorage = s.replicaStorage

	pnum = core.FirstPulseNumber + 1
	pnumNext := pnum + 1
	preparepulse(pnum)
	preparepulse(pnumNext) // should set correct next for previous pulse

	err = sync.Start(s.ctx, jetID1, core.FirstPulseNumber)
	require.Error(s.T(), err)

	err = sync.Start(s.ctx, jetID1, pnum)
	require.NoError(s.T(), err, "start from first+1 pulse on empty storage, jet1")

	err = sync.Start(s.ctx, jetID2, pnum)
	require.NoError(s.T(), err, "start from first+1 pulse on empty storage, jet2")

	err = sync.Store(s.ctx, jetID2, pnum, kvalues2)
	require.NoError(s.T(), err, "store jet2 pulse")

	err = sync.Store(s.ctx, jetID1, pnum, kvalues1)
	require.NoError(s.T(), err, "store jet1 pulse")

	// stop previous
	err = sync.Stop(s.ctx, jetID1, pnum)
	err = sync.Stop(s.ctx, jetID2, pnum)
	require.NoError(s.T(), err)
}
