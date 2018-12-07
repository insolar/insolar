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
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/ledgertestutils"
	"github.com/insolar/insolar/ledger/pulsemanager"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPulseManager_Current(t *testing.T) {
	ctx := inslogger.TestContext(t)
	lr, err := logicrunner.NewLogicRunner(&configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	})
	assert.NoError(t, err)
	c := core.Components{LogicRunner: lr}
	// FIXME: TmpLedger is deprecated. Use mocks instead.
	ledger, cleaner := ledgertestutils.TmpLedger(t, "", c)
	defer cleaner()

	pm := ledger.GetPulseManager()

	pulse, err := pm.Current(ctx)
	assert.NoError(t, err)
	assert.Equal(t, core.GenesisPulse.PulseNumber, pulse.PulseNumber)
}

func TestPulseManager_Set_CheckHotIndexesSending(t *testing.T) {
	// Arrange
	ctx := inslogger.TestContext(t)

	lr := testutils.NewLogicRunnerMock(t)
	lr.OnPulseMock.Return(nil)

	db, dbcancel := storagetest.TmpDB(ctx, t)
	defer dbcancel()
	firstID, _ := db.SetRecord(
		ctx,
		core.GenesisPulse.PulseNumber,
		&record.ObjectActivateRecord{})
	firstIndex := index.ObjectLifeline{
		LatestState: firstID,
	}
	_ = db.SetObjectIndex(ctx, firstID, &firstIndex)
	codeRecord := &record.CodeRecord{}
	secondID, _ := db.SetRecord(
		ctx,
		core.GenesisPulse.PulseNumber,
		codeRecord,
	)

	recentMock := recentstorage.NewRecentStorageMock(t)
	recentMock.ClearZeroTTLObjectsMock.Return()
	recentMock.ClearObjectsMock.Return()
	recentMock.GetObjectsMock.Return(map[core.RecordID]int{
		*firstID: 1,
	})
	recentMock.GetRequestsMock.Return([]core.RecordID{*secondID})
	recentMock.IsMineFunc = func(inputID core.RecordID) (r bool) {
		return bytes.Equal(firstID.Bytes(), inputID.Bytes())
	}

	mbMock := testutils.NewMessageBusMock(t)
	mbMock.SendFunc = func(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) (r core.Reply, r1 error) {
		val, ok := p1.(*message.HotData)
		if !ok {
			return nil, nil
		}

		// Assert
		require.Equal(t, 1, len(val.PendingRequests))
		require.Equal(t, codeRecord, record.DeserializeRecord(val.PendingRequests[*secondID]))

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

	pm := pulsemanager.NewPulseManager(db, configuration.Ledger{})

	gil := testutils.NewGlobalInsolarLockMock(t)
	gil.AcquireMock.Return()
	gil.ReleaseMock.Return()

	alsMock := testutils.NewActiveListSwapperMock(t)
	alsMock.MoveSyncToActiveFunc = func() {}

	pm.LR = lr
	pm.Recent = recentMock
	pm.Bus = mbMock
	pm.NodeNet = nodeNetworkMock
	pm.GIL = gil
	pm.ActiveListSwapper = alsMock

	// Act
	err := pm.Set(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1}, false)
	require.NoError(t, err)
	savedIndex, err := db.GetObjectIndex(ctx, firstID, false)
	require.NoError(t, err)

	// Assert
	require.NotNil(t, savedIndex)
	require.NotNil(t, firstIndex, savedIndex)
	recentMock.MinimockFinish()
}
