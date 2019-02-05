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

package jetcoordinator

import (
	"context"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJetCoordinator_QueryRole(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	db, cleaner := storagetest.TmpDB(ctx, t)

	defer cleaner()
	jc := JetCoordinator{
		JetStorage:                 db,
		PulseTracker:               db,
		ActiveNodesStorage:         db,
		PlatformCryptographyScheme: testutils.NewPlatformCryptographyScheme(),
	}
	ps := storage.NewPulseStorage()
	ps.PulseTracker = db
	jc.PulseStorage = ps

	err := db.AddPulse(ctx, core.Pulse{PulseNumber: 0, Entropy: core.Entropy{1, 2, 3}})
	require.NoError(t, err)
	var nodes []core.Node
	var nodeRefs []core.RecordRef
	for i := 0; i < 100; i++ {
		ref := *core.NewRecordRef(core.DomainID, *core.NewRecordID(0, []byte{byte(i)}))
		nodes = append(nodes, storage.Node{FID: ref, FRole: core.StaticRoleLightMaterial})
		nodeRefs = append(nodeRefs, ref)
	}
	err = db.SetActiveNodes(0, nodes)
	require.NoError(t, err)

	objID := core.NewRecordID(0, []byte{1, 42, 123})
	err = db.UpdateJetTree(ctx, 0, true, *jet.NewID(50, []byte{1, 42, 123}))
	require.NoError(t, err)

	selected, err := jc.QueryRole(ctx, core.DynamicRoleLightValidator, *objID, 0)
	require.NoError(t, err)
	assert.Equal(t, 3, len(selected))

	// Indexes are hard-coded from previously calculated values.
	assert.Equal(t, []core.RecordRef{nodeRefs[16], nodeRefs[21], nodeRefs[78]}, selected)
}

func TestNewJetCoordinator(t *testing.T) {
	t.Parallel()
	// Act
	calc := NewJetCoordinator(12)

	// Assert
	require.NotNil(t, calc)
	require.Equal(t, 12, calc.lightChainLimit)
}

func TestJetCoordinator_IsBeyondLimit_ProblemsWithTracker(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTrackerMock := storage.NewPulseTrackerMock(t)
	pulseTrackerMock.GetPulseMock.Return(nil, errors.New("it's expected"))
	calc := NewJetCoordinator(12)
	calc.PulseTracker = pulseTrackerMock

	// Act
	res, err := calc.IsBeyondLimit(ctx, core.FirstPulseNumber, 0)

	// Assert
	require.NotNil(t, err)
	require.Equal(t, false, res)
}

func TestJetCoordinator_IsBeyondLimit_ProblemsWithTracker_SecondCall(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTrackerMock := storage.NewPulseTrackerMock(t)
	pulseTrackerMock.GetPulseFunc = func(p context.Context, p1 core.PulseNumber) (r *storage.Pulse, r1 error) {
		if p1 == core.FirstPulseNumber {
			return &storage.Pulse{}, nil
		}

		return nil, errors.New("it's expected")
	}
	calc := NewJetCoordinator(12)
	calc.PulseTracker = pulseTrackerMock

	// Act
	res, err := calc.IsBeyondLimit(ctx, core.FirstPulseNumber, 0)

	// Assert
	require.NotNil(t, err)
	require.Equal(t, false, res)
}

func TestJetCoordinator_IsBeyondLimit_OutsideOfLightChainLimit(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTrackerMock := storage.NewPulseTrackerMock(t)
	pulseTrackerMock.GetPulseFunc = func(p context.Context, p1 core.PulseNumber) (r *storage.Pulse, r1 error) {
		if p1 == core.FirstPulseNumber {
			return &storage.Pulse{SerialNumber: 50}, nil
		}

		return &storage.Pulse{SerialNumber: 24}, nil
	}
	calc := NewJetCoordinator(25)
	calc.PulseTracker = pulseTrackerMock

	// Act
	res, err := calc.IsBeyondLimit(ctx, core.FirstPulseNumber, 0)

	// Assert
	require.Nil(t, err)
	require.Equal(t, true, res)
}

func TestJetCoordinator_IsBeyondLimit_InsideOfLightChainLimit(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTrackerMock := storage.NewPulseTrackerMock(t)
	pulseTrackerMock.GetPulseFunc = func(p context.Context, p1 core.PulseNumber) (r *storage.Pulse, r1 error) {
		if p1 == core.FirstPulseNumber {
			return &storage.Pulse{SerialNumber: 50}, nil
		}

		return &storage.Pulse{SerialNumber: 34}, nil
	}
	calc := NewJetCoordinator(25)
	calc.PulseTracker = pulseTrackerMock

	// Act
	res, err := calc.IsBeyondLimit(ctx, core.FirstPulseNumber, 0)

	// Assert
	require.Nil(t, err)
	require.Equal(t, false, res)
}

func TestJetCoordinator_NodeForJet_CheckLimitFailed(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTrackerMock := storage.NewPulseTrackerMock(t)
	pulseTrackerMock.GetPulseMock.Return(nil, errors.New("it's expected"))
	calc := NewJetCoordinator(12)
	calc.PulseTracker = pulseTrackerMock

	// Act
	res, err := calc.NodeForJet(ctx, testutils.RandomJet(), core.FirstPulseNumber, 0)

	// Assert
	require.NotNil(t, err)
	require.Nil(t, res)
}

func TestJetCoordinator_NodeForJet_GoToHeavy(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTrackerMock := storage.NewPulseTrackerMock(t)
	pulseTrackerMock.GetPulseFunc = func(p context.Context, p1 core.PulseNumber) (r *storage.Pulse, r1 error) {
		if p1 == core.FirstPulseNumber {
			return &storage.Pulse{SerialNumber: 50}, nil
		}

		return &storage.Pulse{SerialNumber: 24}, nil
	}
	expectedID := core.NewRecordRef(testutils.RandomID(), testutils.RandomID())
	nodeMock := network.NewNodeMock(t)
	nodeMock.IDMock.Return(*expectedID)
	activeNodesStorageMock := storage.NewActiveNodesStorageMock(t)
	activeNodesStorageMock.GetActiveNodesByRoleFunc = func(p core.PulseNumber, p1 core.StaticRole) (r []core.Node, r1 error) {
		require.Equal(t, core.FirstPulseNumber, int(p))
		require.Equal(t, core.StaticRoleHeavyMaterial, p1)

		return []core.Node{nodeMock}, nil
	}

	pulseStorageMock := testutils.NewPulseStorageMock(t)
	pulseStorageMock.CurrentFunc = func(p context.Context) (r *core.Pulse, r1 error) {
		generator := entropygenerator.StandardEntropyGenerator{}
		return &core.Pulse{PulseNumber: core.FirstPulseNumber, Entropy: generator.GenerateEntropy()}, nil
	}

	calc := NewJetCoordinator(25)
	calc.PulseTracker = pulseTrackerMock
	calc.ActiveNodesStorage = activeNodesStorageMock
	calc.PulseStorage = pulseStorageMock
	calc.PlatformCryptographyScheme = platformpolicy.NewPlatformCryptographyScheme()

	// Act
	resNode, err := calc.NodeForJet(ctx, testutils.RandomJet(), core.FirstPulseNumber, 0)

	// Assert
	require.Nil(t, err)
	require.Equal(t, expectedID, resNode)
}

func TestJetCoordinator_NodeForJet_GoToLight(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	pulseTrackerMock := storage.NewPulseTrackerMock(t)
	pulseTrackerMock.GetPulseFunc = func(p context.Context, p1 core.PulseNumber) (r *storage.Pulse, r1 error) {
		if p1 == core.FirstPulseNumber {
			return &storage.Pulse{SerialNumber: 50}, nil
		}

		return &storage.Pulse{SerialNumber: 49}, nil
	}
	expectedID := core.NewRecordRef(testutils.RandomID(), testutils.RandomID())
	nodeMock := network.NewNodeMock(t)
	nodeMock.IDMock.Return(*expectedID)
	activeNodesStorageMock := storage.NewActiveNodesStorageMock(t)
	activeNodesStorageMock.GetActiveNodesByRoleFunc = func(p core.PulseNumber, p1 core.StaticRole) (r []core.Node, r1 error) {
		require.Equal(t, 0, int(p))
		require.Equal(t, core.StaticRoleLightMaterial, p1)

		return []core.Node{nodeMock}, nil
	}

	pulseStorageMock := testutils.NewPulseStorageMock(t)
	pulseStorageMock.CurrentFunc = func(p context.Context) (r *core.Pulse, r1 error) {
		generator := entropygenerator.StandardEntropyGenerator{}
		return &core.Pulse{PulseNumber: core.FirstPulseNumber, Entropy: generator.GenerateEntropy()}, nil
	}

	calc := NewJetCoordinator(25)
	calc.PulseTracker = pulseTrackerMock
	calc.ActiveNodesStorage = activeNodesStorageMock
	calc.PulseStorage = pulseStorageMock
	calc.PlatformCryptographyScheme = platformpolicy.NewPlatformCryptographyScheme()

	// Act
	resNode, err := calc.NodeForJet(ctx, testutils.RandomJet(), core.FirstPulseNumber, 0)

	// Assert
	require.Nil(t, err)
	require.Equal(t, expectedID, resNode)
}
