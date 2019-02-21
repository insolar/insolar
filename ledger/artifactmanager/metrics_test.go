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

package artifactmanager

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/nodes"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/testutils/testmetrics"
)

type metricSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()
	db      storage.DBContext

	scheme        core.PlatformCryptographyScheme
	pulseTracker  storage.PulseTracker
	nodeStorage   nodes.Accessor
	objectStorage storage.ObjectStorage
	jetStorage    storage.JetStorage
	dropStorage   storage.DropStorage
	genesisState  storage.GenesisState
}

func NewMetricSuite() *metricSuite {
	return &metricSuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestMetricSuite(t *testing.T) {
	suite.Run(t, NewMetricSuite())
}

func (s *metricSuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())

	db, cleaner := storagetest.TmpDB(s.ctx, s.T())
	s.cleaner = cleaner
	s.db = db
	s.scheme = testutils.NewPlatformCryptographyScheme()
	s.jetStorage = storage.NewJetStorage()
	s.nodeStorage = nodes.NewStorage()
	s.pulseTracker = storage.NewPulseTracker()
	s.objectStorage = storage.NewObjectStorage()
	s.dropStorage = storage.NewDropStorage(10)
	s.genesisState = storage.NewGenesisInitializer()

	s.cm.Inject(
		s.scheme,
		s.db,
		s.jetStorage,
		s.nodeStorage,
		s.pulseTracker,
		s.objectStorage,
		s.dropStorage,
		s.genesisState,
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

func (s *metricSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

func (s *metricSuite) TestLedgerArtifactManager_Metrics() {
	// BEWARE: this test should not be run in parallel!
	mc := minimock.NewController(s.T())
	defer mc.Finish()

	amPulseStorageMock := testutils.NewPulseStorageMock(s.T())
	amPulseStorageMock.CurrentFunc = func(p context.Context) (r *core.Pulse, r1 error) {
		pulse, err := s.pulseTracker.GetLatestPulse(p)
		require.NoError(s.T(), err)
		return &pulse.Pulse, err
	}

	mb := testutils.NewMessageBusMock(mc)
	mb.SendMock.Return(&reply.ID{}, nil)
	cs := platformpolicy.NewPlatformCryptographyScheme()
	am := NewArtifactManger()
	am.DB = s.db
	am.PlatformCryptographyScheme = cs
	am.DefaultBus = mb
	am.PulseStorage = amPulseStorageMock
	am.GenesisState = s.genesisState

	tmetrics := testmetrics.Start(s.ctx)
	defer tmetrics.Stop()

	msg := message.GenesisRequest{Name: "4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.4FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa"}
	_, err := am.RegisterRequest(s.ctx, *am.GenesisRef(), &message.Parcel{Msg: &msg})
	require.NoError(s.T(), err)

	time.Sleep(1500 * time.Millisecond)

	_ = am
	content, err := tmetrics.FetchContent()
	require.NoError(s.T(), err)

	assert.Contains(s.T(), content, `insolar_artifactmanager_latency_count{method="RegisterRequest",result="2xx"}`)
	assert.Contains(s.T(), content, `insolar_artifactmanager_calls{method="RegisterRequest",result="2xx"}`)
}
