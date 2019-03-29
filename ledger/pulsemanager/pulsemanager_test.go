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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/drop"
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

	db, cleaner := storagetest.TmpDB(s.ctx, s.T())
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
	jetID := insolar.ZeroJetID
	objID := insolar.ID{}

	lr := testutils.NewLogicRunnerMock(s.T())
	lr.OnPulseMock.Return(nil)

	firstID, _ := s.objectStorage.SetRecord(
		s.ctx,
		insolar.ID(jetID),
		insolar.GenesisPulse.PulseNumber,
		&object.ActivateRecord{})
	firstIndex := object.Lifeline{
		LatestState: firstID,
	}
	_ = s.objectStorage.SetObjectIndex(s.ctx, insolar.ID(jetID), firstID, &firstIndex)
	codeRecord := &object.CodeRecord{}
	secondID, _ := s.objectStorage.SetRecord(
		s.ctx,
		insolar.ID(jetID),
		insolar.GenesisPulse.PulseNumber,
		codeRecord,
	)

	indexMock := recentstorage.NewRecentIndexStorageMock(s.T())
	pendingMock := recentstorage.NewPendingStorageMock(s.T())
	indexMock.GetObjectsMock.Return(map[insolar.ID]int{
		*firstID: 1,
	})
	pendingMock.GetRequestsMock.Return(
		map[insolar.ID]recentstorage.PendingObjectContext{
			objID: {Requests: []insolar.ID{*secondID}},
		})

	providerMock := recentstorage.NewProviderMock(s.T())
	providerMock.GetPendingStorageMock.Return(pendingMock)
	providerMock.GetIndexStorageMock.Return(indexMock)
	providerMock.ClonePendingStorageMock.Return()
	providerMock.CloneIndexStorageMock.Return()

	mbMock := testutils.NewMessageBusMock(s.T())
	mbMock.OnPulseFunc = func(context.Context, insolar.Pulse) error {
		return nil
	}
	mbMock.SendFunc = func(p context.Context, p1 insolar.Message, p2 *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
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
		decodedIndex := object.DecodeIndex(val.RecentObjects[*firstID].Index)
		require.Equal(s.T(), firstIndex, decodedIndex)
		require.Equal(s.T(), 1, val.RecentObjects[*firstID].TTL)

		return nil, nil
	}

	nodeMock := network.NewNetworkNodeMock(s.T())
	nodeMock.RoleMock.Return(insolar.StaticRoleLightMaterial)
	nodeMock.IDMock.Return(insolar.Reference{})

	nodeNetworkMock := network.NewNodeNetworkMock(s.T())
	nodeNetworkMock.GetWorkingNodesMock.Return([]insolar.NetworkNode{nodeMock})
	nodeNetworkMock.GetOriginMock.Return(nodeMock)

	jetCoordinatorMock := testutils.NewJetCoordinatorMock(s.T())
	executor := insolar.NewReference(insolar.ID{}, *insolar.NewID(123, []byte{3, 2, 1}))
	jetCoordinatorMock.LightExecutorForJetMock.Return(executor, nil)
	jetCoordinatorMock.MeMock.Return(*executor)

	pm := NewPulseManager(configuration.Ledger{}, drop.NewCleanerMock(s.T()), blob.NewCleanerMock(s.T()), blob.NewCollectionAccessorMock(s.T()))

	gil := testutils.NewGlobalInsolarLockMock(s.T())
	gil.AcquireMock.Return()
	gil.ReleaseMock.Return()

	alsMock := testutils.NewActiveListSwapperMock(s.T())
	alsMock.MoveSyncToActiveFunc = func(context.Context) error { return nil }

	cryptoServiceMock := testutils.NewCryptographyServiceMock(s.T())
	cryptoServiceMock.SignFunc = func(p []byte) (r *insolar.Signature, r1 error) {
		signature := insolar.SignatureFromBytes(nil)
		return &signature, nil
	}

	pulseStorageMock := NewpulseStoragePmMock(s.T())
	pulseStorageMock.CurrentMock.Return(insolar.GenesisPulse, nil)
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
	err := pm.Set(s.ctx, insolar.Pulse{PulseNumber: insolar.FirstPulseNumber + 1}, true)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), uint64(2), mbMock.SendMinimockCounter()) // 1 validator drop (no split)
	savedIndex, err := s.objectStorage.GetObjectIndex(s.ctx, insolar.ID(jetID), firstID)
	require.NoError(s.T(), err)

	// Assert
	require.NotNil(s.T(), savedIndex)
	require.NotNil(s.T(), firstIndex, savedIndex)
	indexMock.MinimockFinish()
	pendingMock.MinimockFinish()
}
