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
	"bytes"
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
	// Arrange
	ctx := inslogger.TestContext(t)
	jetID := core.TODOJetID
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
	recentMock.ClearZeroTTLObjectsMock.Return()
	recentMock.ClearObjectsMock.Return()
	recentMock.GetObjectsMock.Return(map[core.RecordID]int{
		*firstID: 1,
	})
	recentMock.GetRequestsMock.Return(map[core.RecordID]map[core.RecordID]struct{}{objID: {*secondID: struct{}{}}})
	recentMock.IsMineFunc = func(inputID core.RecordID) (r bool) {
		return bytes.Equal(firstID.Bytes(), inputID.Bytes())
	}

	providerMock := recentstorage.NewProviderMock(t)
	providerMock.GetStorageMock.Return(recentMock)

	mbMock := testutils.NewMessageBusMock(t)
	mbMock.SendFunc = func(p context.Context, p1 core.Message, _ core.Pulse, p2 *core.MessageSendOptions) (r core.Reply, r1 error) {
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

	nodeNetworkMock := network.NewNodeNetworkMock(t)
	nodeNetworkMock.GetActiveNodesMock.Return([]core.Node{nodeMock})
	nodeNetworkMock.GetOriginMock.Return(nodeMock)

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

	pm.LR = lr

	pm.RecentStorageProvider = providerMock
	pm.Bus = mbMock
	pm.NodeNet = nodeNetworkMock
	pm.GIL = gil
	pm.ActiveListSwapper = alsMock
	pm.CryptographyService = cryptoServiceMock
	pm.PlatformCryptographyScheme = testutils.NewPlatformCryptographyScheme()
	pm.PulseStorage = pulseStorageMock

	// Act
	err := pm.Set(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1}, true)
	require.NoError(t, err)
	assert.Equal(t, uint64(2), mbMock.SendMinimockCounter()) // 1 validator drop + 1 executor (no split)
	savedIndex, err := db.GetObjectIndex(ctx, jetID, firstID, false)
	require.NoError(t, err)

	// Assert
	require.NotNil(t, savedIndex)
	require.NotNil(t, firstIndex, savedIndex)
	recentMock.MinimockFinish()
}

func TestPulseManager_Set_PerformsSplit(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jetID := *jet.NewID(0, nil)

	lr := testutils.NewLogicRunnerMock(t)
	lr.OnPulseMock.Return(nil)

	db, dbcancel := storagetest.TmpDB(ctx, t)
	defer dbcancel()

	err := db.AddDropSize(ctx, &jet.DropSize{
		JetID:    jetID,
		DropSize: 100,
	})
	require.NoError(t, err)

	recentMock := recentstorage.NewRecentStorageMock(t)
	recentMock.ClearZeroTTLObjectsMock.Return()
	recentMock.ClearObjectsMock.Return()
	recentMock.GetObjectsMock.Return(nil)
	recentMock.GetRequestsMock.Return(nil)
	recentMock.IsMineMock.Return(true)

	providerMock := recentstorage.NewProviderMock(t)
	providerMock.GetStorageMock.Return(recentMock)

	mbMock := testutils.NewMessageBusMock(t)
	mbMock.SendMock.Return(nil, nil)

	nodeMock := network.NewNodeMock(t)
	nodeMock.RoleMock.Return(core.StaticRoleLightMaterial)

	nodeNetworkMock := network.NewNodeNetworkMock(t)
	nodeNetworkMock.GetActiveNodesMock.Return([]core.Node{nodeMock})
	nodeNetworkMock.GetOriginMock.Return(nodeMock)

	pulseStorage := pulsemanager.NewpulseStoragePmMock(t)
	pulseStorage.LockMock.Return()
	pulseStorage.UnlockMock.Return()

	pm := pulsemanager.NewPulseManager(db, configuration.Ledger{
		JetSizesHistoryDepth: 2,
		PulseManager:         configuration.PulseManager{SplitThreshold: 0},
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

	pm.LR = lr
	pm.RecentStorageProvider = providerMock
	pm.Bus = mbMock
	pm.NodeNet = nodeNetworkMock
	pm.GIL = gil
	pm.ActiveListSwapper = alsMock
	pm.CryptographyService = cryptoServiceMock
	pm.PlatformCryptographyScheme = testutils.NewPlatformCryptographyScheme()
	pm.PulseStorage = pulseStorage

	err = pm.Set(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1}, true)
	require.NoError(t, err)
	assert.Equal(t, uint64(3), mbMock.SendMinimockCounter()) // 1 validator drop + 2 executors (split)
}
