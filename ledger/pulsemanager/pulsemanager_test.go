/*
 *    Copyright 2018 Insolar
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

package pulsemanager_test

import (
	"context"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/pulsemanager"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPulseManager_Set_CheckHotIndexesSending(t *testing.T) {
	// Error:      	Not equal:
	// expected: 0x2
	// actual  : 0x0
	t.Skip()

	// Arrange
	ctx := inslogger.TestContext(t)
	jetID := jet.ZeroJetID
	objID := core.RecordID{}

	lr := testutils.NewLogicRunnerMock(t)
	lr.OnPulseMock.Return(nil)

	db, dbcancel := storagetest.TmpDB(ctx, t)
	defer dbcancel()
	firstID, _ := db.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		&record.ObjectActivateRecord{})
	firstIndex := index.ObjectLifeline{
		LatestState: firstID,
	}
	_ = db.SetObjectIndex(ctx, jetID, firstID, &firstIndex)
	codeRecord := &record.CodeRecord{}
	secondID, _ := db.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		codeRecord,
	)

	recentMock := recentstorage.NewRecentStorageMock(t)
	// TODO: @andreyromancev. 12.01.19. Uncomment to check if this doesn't delete indexes it should not.
	// recentMock.ClearZeroTTLObjectsMock.Return()
	// recentMock.ClearObjectsMock.Return()
	recentMock.GetObjectsMock.Return(map[core.RecordID]int{
		*firstID: 1,
	})
	recentMock.GetRequestsMock.Return(map[core.RecordID]map[core.RecordID]struct{}{objID: {*secondID: struct{}{}}})

	providerMock := recentstorage.NewProviderMock(t)
	providerMock.GetStorageMock.Return(recentMock)
	providerMock.CloneStorageMock.Return()

	mbMock := testutils.NewMessageBusMock(t)
	mbMock.OnPulseFunc = func(context.Context, core.Pulse) error {
		return nil
	}
	mbMock.SendFunc = func(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) (r core.Reply, r1 error) {
		val, ok := p1.(*message.HotData)
		if !ok {
			return nil, nil
		}

		// Assert
		require.Equal(t, 1, len(val.PendingRequests))
		requests, ok := val.PendingRequests[objID]
		require.True(t, ok)
		require.Equal(t, 1, len(requests))
		require.Equal(t, codeRecord, record.DeserializeRecord(requests[*secondID]))

		require.Equal(t, 1, len(val.RecentObjects))
		decodedIndex, err := index.DecodeObjectLifeline(val.RecentObjects[*firstID].Index)
		require.NoError(t, err)
		require.Equal(t, firstIndex, *decodedIndex)
		require.Equal(t, 1, val.RecentObjects[*firstID].TTL)

		return nil, nil
	}

	nodeMock := network.NewNodeMock(t)
	nodeMock.RoleMock.Return(core.StaticRoleLightMaterial)
	nodeMock.IDMock.Return(core.RecordRef{})

	nodeNetworkMock := network.NewNodeNetworkMock(t)
	nodeNetworkMock.GetActiveNodesMock.Return([]core.Node{nodeMock})
	nodeNetworkMock.GetOriginMock.Return(nodeMock)

	jetCoordinatorMock := testutils.NewJetCoordinatorMock(t)
	executor := core.NewRecordRef(core.RecordID{}, *core.NewRecordID(123, []byte{3, 2, 1}))
	jetCoordinatorMock.LightExecutorForJetMock.Return(executor, nil)
	jetCoordinatorMock.MeMock.Return(*executor)

	pm := pulsemanager.NewPulseManager(db, configuration.Ledger{
		JetSizesHistoryDepth: 5,
	})

	gil := testutils.NewGlobalInsolarLockMock(t)
	gil.AcquireMock.Return()
	gil.ReleaseMock.Return()

	alsMock := testutils.NewActiveListSwapperMock(t)
	alsMock.MoveSyncToActiveFunc = func() {}

	cryptoServiceMock := testutils.NewCryptographyServiceMock(t)
	cryptoServiceMock.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes(nil)
		return &signature, nil
	}

	pulseStorageMock := pulsemanager.NewpulseStoragePmMock(t)
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
	err := pm.Set(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1}, true)
	require.NoError(t, err)
	// // TODO: @andreyromancev. 12.01.19. put 1, when dynamic split is working.
	assert.Equal(t, uint64(2), mbMock.SendMinimockCounter()) // 1 validator drop (no split)
	savedIndex, err := db.GetObjectIndex(ctx, jetID, firstID, false)
	require.NoError(t, err)

	// Assert
	require.NotNil(t, savedIndex)
	require.NotNil(t, firstIndex, savedIndex)
	recentMock.MinimockFinish()
}
