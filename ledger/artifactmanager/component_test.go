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

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/nodes"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/testmessagebus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type componentSuite struct {
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
}

func NewComponentSuite() *componentSuite {
	return &componentSuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestComponentSuite(t *testing.T) {
	suite.Run(t, NewComponentSuite())
}

func (s *componentSuite) BeforeTest(suiteName, testName string) {
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

	s.cm.Inject(
		s.scheme,
		s.db,
		s.jetStorage,
		s.nodeStorage,
		s.pulseTracker,
		s.objectStorage,
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

func (s *componentSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

func (s *componentSuite) TestLedgerArtifactManager_PendingRequest() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()

	jetID := *jet.NewID(0, nil)

	amPulseStorageMock := testutils.NewPulseStorageMock(s.T())
	amPulseStorageMock.CurrentFunc = func(p context.Context) (r *core.Pulse, r1 error) {
		pulse, err := s.pulseTracker.GetLatestPulse(p)
		require.NoError(s.T(), err)
		return &pulse.Pulse, err
	}

	jcMock := testutils.NewJetCoordinatorMock(s.T())
	jcMock.LightExecutorForJetMock.Return(&core.RecordRef{}, nil)
	jcMock.MeMock.Return(core.RecordRef{})

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	cs := testutils.NewPlatformCryptographyScheme()
	mb := testmessagebus.NewTestMessageBus(s.T())
	mb.PulseStorage = amPulseStorageMock

	am := NewArtifactManger()
	am.PulseStorage = amPulseStorageMock
	am.PlatformCryptographyScheme = cs
	am.DefaultBus = mb
	am.PlatformCryptographyScheme = platformpolicy.NewPlatformCryptographyScheme()

	provider := recentstorage.NewRecentStorageProvider(0)

	cryptoScheme := testutils.NewPlatformCryptographyScheme()

	handler := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 10,
	},
		certificate)

	handler.JetStorage = s.jetStorage
	handler.Nodes = s.nodeStorage
	handler.DBContext = s.db
	handler.PulseTracker = s.pulseTracker
	handler.ObjectStorage = s.objectStorage

	handler.PlatformCryptographyScheme = cryptoScheme
	handler.Bus = mb
	handler.RecentStorageProvider = provider
	handler.JetCoordinator = jcMock

	handler.HotDataWaiter = NewHotDataWaiterConcrete()
	err := handler.HotDataWaiter.Unlock(s.ctx, jetID)
	require.NoError(s.T(), err)

	err = handler.Init(s.ctx)
	require.NoError(s.T(), err)
	objRef := *genRandomRef(0)

	s.jetStorage.UpdateJetTree(s.ctx, core.FirstPulseNumber, true, jetID)
	s.jetStorage.UpdateJetTree(s.ctx, core.FirstPulseNumber+1, true, jetID)

	// Register request
	reqID, err := am.RegisterRequest(s.ctx, objRef, &message.Parcel{Msg: &message.CallMethod{}, PulseNumber: core.FirstPulseNumber})
	require.NoError(s.T(), err)

	// Change pulse.
	err = s.pulseTracker.AddPulse(s.ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1})
	require.NoError(s.T(), err)

	// Should have pending request.
	has, err := am.HasPendingRequests(s.ctx, objRef)
	require.NoError(s.T(), err)
	assert.True(s.T(), has)

	// Register result.
	reqRef := *core.NewRecordRef(core.DomainID, *reqID)
	_, err = am.RegisterResult(s.ctx, objRef, reqRef, nil)
	require.NoError(s.T(), err)

	// Should not have pending request.
	has, err = am.HasPendingRequests(s.ctx, objRef)
	require.NoError(s.T(), err)
	assert.False(s.T(), has)
}
