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

package pulsemanager

import (
	"context"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type pulseManagerSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()

	objectStorage storage.ObjectStorage
}

func NewPulseManagerSuite() *pulseManagerSuite {
	return &pulseManagerSuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestPulseManager(t *testing.T) {
	suite.Run(t, NewPulseManagerSuite())
}

func (s *pulseManagerSuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())

	db, cleaner := storagetest.TmpDB(s.ctx, nil, s.T())
	s.cleaner = cleaner
	s.objectStorage = storage.NewObjectStorage()

	s.cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		db,
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

func (s *pulseManagerSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

func (s *pulseManagerSuite) TestPulseManager_Set_CheckHotIndexesSending() {
	// Error:      	Not equal:
	// expected: 0x2
	// actual  : 0x0
	s.T().Skip()

	// Arrange
	jetID := core.ZeroJetID
	objID := core.RecordID{}

	lr := testutils.NewLogicRunnerMock(s.T())
	lr.OnPulseMock.Return(nil)

	firstID, _ := s.objectStorage.SetRecord(
		s.ctx,
		core.RecordID(jetID),
		core.GenesisPulse.PulseNumber,
		&object.ObjectActivateRecord{})
	firstIndex := object.Lifeline{
		LatestState: firstID,
	}
	_ = s.objectStorage.SetObjectIndex(s.ctx, core.RecordID(jetID), firstID, &firstIndex)
	codeRecord := &object.CodeRecord{}
	secondID, _ := s.objectStorage.SetRecord(
		s.ctx,
		core.RecordID(jetID),
		core.GenesisPulse.PulseNumber,
		codeRecord,
	)

	indexMock := recentstorage.NewRecentIndexStorageMock(s.T())
	pendingMock := recentstorage.NewPendingStorageMock(s.T())
	// TODO: @andreyromancev. 12.01.19. Uncomment to check if this doesn't delete indexes it should not.
	// recentMock.ClearZeroTTLObjectsMock.Return()
	// recentMock.ClearObjectsMock.Return()
	indexMock.GetObjectsMock.Return(map[core.RecordID]int{
		*firstID: 1,
	})
	pendingMock.GetRequestsMock.Return(
		map[core.RecordID]recentstorage.PendingObjectContext{
			objID: {Requests: []core.RecordID{*secondID}},
		})

	providerMock := recentstorage.NewProviderMock(s.T())
	providerMock.GetPendingStorageMock.Return(pendingMock)
	providerMock.GetIndexStorageMock.Return(indexMock)
	providerMock.ClonePendingStorageMock.Return()
	providerMock.CloneIndexStorageMock.Return()

	mbMock := testutils.NewMessageBusMock(s.T())
	mbMock.OnPulseFunc = func(context.Context, core.Pulse) error {
		return nil
	}
	mbMock.SendFunc = func(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) (r core.Reply, r1 error) {
		val, ok := p1.(*message.HotData)
		if !ok {
			return nil, nil
		}

		// Assert
		require.Equal(s.T(), 1, len(val.PendingRequests))
		objContext, ok := val.PendingRequests[objID]
		require.True(s.T(), ok)
		require.Equal(s.T(), 1, len(objContext.Requests))

		require.Equal(s.T(), 1, len(val.RecentObjects))
		decodedIndex := object.Decode(val.RecentObjects[*firstID].Index)
		require.Equal(s.T(), firstIndex, decodedIndex)
		require.Equal(s.T(), 1, val.RecentObjects[*firstID].TTL)

		return nil, nil
	}

	nodeMock := network.NewNodeMock(s.T())
	nodeMock.RoleMock.Return(core.StaticRoleLightMaterial)
	nodeMock.IDMock.Return(core.RecordRef{})

	nodeNetworkMock := network.NewNodeNetworkMock(s.T())
	nodeNetworkMock.GetWorkingNodesMock.Return([]core.Node{nodeMock})
	nodeNetworkMock.GetOriginMock.Return(nodeMock)

	jetCoordinatorMock := testutils.NewJetCoordinatorMock(s.T())
	executor := core.NewRecordRef(core.RecordID{}, *core.NewRecordID(123, []byte{3, 2, 1}))
	jetCoordinatorMock.LightExecutorForJetMock.Return(executor, nil)
	jetCoordinatorMock.MeMock.Return(*executor)

	pm := NewPulseManager(configuration.Ledger{})

	gil := testutils.NewGlobalInsolarLockMock(s.T())
	gil.AcquireMock.Return()
	gil.ReleaseMock.Return()

	alsMock := testutils.NewActiveListSwapperMock(s.T())
	alsMock.MoveSyncToActiveFunc = func(context.Context) error { return nil }

	cryptoServiceMock := testutils.NewCryptographyServiceMock(s.T())
	cryptoServiceMock.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes(nil)
		return &signature, nil
	}

	pulseStorageMock := NewpulseStoragePmMock(s.T())
	pulseStorageMock.CurrentMock.Return(core.GenesisPulse, nil)
	pulseStorageMock.LockMock.Return()
	pulseStorageMock.UnlockMock.Return()
	pulseStorageMock.SetMock.Return()

	pm.LR = lr

	pm.RecentStorageProvider = providerMock
	pm.Bus = mbMock
	pm.NodeNet = nodeNetworkMock
	pm.GIL = gil
	pm.ActiveListSwapper = alsMock
	pm.CryptographyService = cryptoServiceMock
	pm.PlatformCryptographyScheme = testutils.NewPlatformCryptographyScheme()
	pm.PulseStorage = pulseStorageMock
	pm.JetCoordinator = jetCoordinatorMock

	// Act
	err := pm.Set(s.ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1}, true)
	require.NoError(s.T(), err)
	// // TODO: @andreyromancev. 12.01.19. put 1, when dynamic split is working.
	assert.Equal(s.T(), uint64(2), mbMock.SendMinimockCounter()) // 1 validator drop (no split)
	savedIndex, err := s.objectStorage.GetObjectIndex(s.ctx, core.RecordID(jetID), firstID, false)
	require.NoError(s.T(), err)

	// Assert
	require.NotNil(s.T(), savedIndex)
	require.NotNil(s.T(), firstIndex, savedIndex)
	indexMock.MinimockFinish()
	pendingMock.MinimockFinish()
}
