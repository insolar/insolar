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
	"sync"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
)

func TestJetTreeUpdater_otherNodesForPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	jc := testutils.NewJetCoordinatorMock(mc)
	ans := storage.NewNodeStorageMock(mc)
	js := storage.NewJetStorageMock(mc)
	jtu := &jetTreeUpdater{
		ActiveNodesStorage: ans,
		JetStorage:         js,
		JetCoordinator:     jc,
	}

	t.Run("active nodes storage returns error", func(t *testing.T) {
		ans.GetActiveNodesByRoleMock.Expect(
			100, core.StaticRoleLightMaterial,
		).Return(
			nil, errors.New("some"),
		)

		nodes, err := jtu.otherNodesForPulse(ctx, core.PulseNumber(100))
		require.Error(t, err)
		require.Empty(t, nodes)
	})

	meRef := testutils.RandomRef()
	jc.MeMock.Return(meRef)

	t.Run("no active nodes at all", func(t *testing.T) {

		ans.GetActiveNodesByRoleMock.Expect(
			100, core.StaticRoleLightMaterial,
		).Return(
			[]core.Node{}, nil,
		)

		nodes, err := jtu.otherNodesForPulse(ctx, core.PulseNumber(100))
		require.Error(t, err)
		require.Empty(t, nodes)
	})

	t.Run("one active node, it's me", func(t *testing.T) {

		someNode := network.NewNodeMock(mc)
		someNode.IDMock.Return(meRef)

		ans.GetActiveNodesByRoleMock.Expect(
			100, core.StaticRoleLightMaterial,
		).Return(
			[]core.Node{someNode}, nil,
		)

		nodes, err := jtu.otherNodesForPulse(ctx, core.PulseNumber(100))
		require.Error(t, err)
		require.Empty(t, nodes)
	})

	t.Run("active node", func(t *testing.T) {
		someNode := network.NewNodeMock(mc)
		someNode.IDMock.Return(testutils.RandomRef())

		ans.GetActiveNodesByRoleMock.Expect(
			100, core.StaticRoleLightMaterial,
		).Return(
			[]core.Node{someNode}, nil,
		)

		nodes, err := jtu.otherNodesForPulse(ctx, core.PulseNumber(100))
		require.NoError(t, err)
		require.Contains(t, nodes, someNode)
	})

	t.Run("active node and me", func(t *testing.T) {
		meNode := network.NewNodeMock(mc)
		meNode.IDMock.Return(meRef)

		someNode := network.NewNodeMock(mc)
		someNode.IDMock.Return(testutils.RandomRef())

		ans.GetActiveNodesByRoleMock.Expect(
			100, core.StaticRoleLightMaterial,
		).Return(
			[]core.Node{someNode, meNode}, nil,
		)

		nodes, err := jtu.otherNodesForPulse(ctx, core.PulseNumber(100))
		require.NoError(t, err)
		require.Contains(t, nodes, someNode)
		require.NotContains(t, nodes, meNode)
	})
}

func TestJetTreeUpdater_fetchActualJetFromOtherNodes(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	jc := testutils.NewJetCoordinatorMock(mc)
	ans := storage.NewNodeStorageMock(mc)
	js := storage.NewJetStorageMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	jtu := &jetTreeUpdater{
		ActiveNodesStorage: ans,
		JetStorage:         js,
		JetCoordinator:     jc,
		MessageBus:         mb,
	}

	meRef := testutils.RandomRef()
	jc.MeMock.Return(meRef)

	otherNode := network.NewNodeMock(mc)
	otherNode.IDMock.Return(testutils.RandomRef())

	ans.GetActiveNodesByRoleMock.Expect(
		100, core.StaticRoleLightMaterial,
	).Return(
		[]core.Node{otherNode}, nil,
	)

	t.Run("MB error on fetching actual jet", func(t *testing.T) {
		target := testutils.RandomID()

		mb.SendMock.Return(nil, errors.New("some"))

		jetID, err := jtu.fetchActualJetFromOtherNodes(ctx, target, core.PulseNumber(100))
		require.Error(t, err)
		require.Nil(t, jetID)
	})

	t.Run("MB got one not actual jet", func(t *testing.T) {
		target := testutils.RandomID()

		mb.SendMock.Return(
			&reply.Jet{ID: *jet.NewID(0, nil), Actual: false},
			nil,
		)

		jetID, err := jtu.fetchActualJetFromOtherNodes(ctx, target, core.PulseNumber(100))
		require.Error(t, err)
		require.Nil(t, jetID)
	})
	t.Run("MB got one actual jet", func(t *testing.T) {
		target := testutils.RandomID()

		mb.SendMock.Return(
			&reply.Jet{ID: *jet.NewID(0, nil), Actual: true},
			nil,
		)

		jetID, err := jtu.fetchActualJetFromOtherNodes(ctx, target, core.PulseNumber(100))
		require.NoError(t, err)
		require.Equal(t, jet.NewID(0, nil), jetID)
	})

	// TODO: multiple nodes returned different results
	// TODO: multiple nodes returned the same result
}

func TestJetTreeUpdater_fetchJet(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	jc := testutils.NewJetCoordinatorMock(mc)
	ans := storage.NewNodeStorageMock(mc)
	js := storage.NewJetStorageMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	jtu := &jetTreeUpdater{
		ActiveNodesStorage: ans,
		JetStorage:         js,
		JetCoordinator:     jc,
		MessageBus:         mb,
		sequencer: map[string]*struct {
			sync.Mutex
			done bool
		}{},
	}

	target := testutils.RandomID()

	t.Run("wrong tree", func(t *testing.T) {
		js.GetJetTreeMock.Return(nil, errors.New("some"))
		jetID, err := jtu.fetchJet(ctx, target, core.PulseNumber(100))
		require.Error(t, err)
		require.Nil(t, jetID)
	})

	t.Run("quick reply, data is up to date", func(t *testing.T) {
		js.GetJetTreeMock.Return(
			jet.NewTree(true), nil,
		)
		jetID, err := jtu.fetchJet(ctx, target, core.PulseNumber(100))
		require.NoError(t, err)
		require.Equal(t, jet.NewID(0, nil), jetID)
	})

	t.Run("fetch jet from friends", func(t *testing.T) {
		meRef := testutils.RandomRef()
		jc.MeMock.Return(meRef)

		otherNode := network.NewNodeMock(mc)
		otherNode.IDMock.Return(testutils.RandomRef())

		ans.GetActiveNodesByRoleMock.Expect(
			100, core.StaticRoleLightMaterial,
		).Return(
			[]core.Node{otherNode}, nil,
		)
		mb.SendMock.Return(
			&reply.Jet{ID: *jet.NewID(0, nil), Actual: true},
			nil,
		)

		js.GetJetTreeMock.Return(
			jet.NewTree(false), nil,
		)
		js.UpdateJetTreeFunc = func(ctx context.Context, pn core.PulseNumber, actual bool, jets ...core.RecordID) (r error) {
			require.Equal(t, core.PulseNumber(100), pn)
			require.True(t, actual)
			require.Equal(t, []core.RecordID{*jet.NewID(0, nil)}, jets)
			return nil
		}

		jetID, err := jtu.fetchJet(ctx, target, core.PulseNumber(100))
		require.NoError(t, err)
		require.Equal(t, jet.NewID(0, nil), jetID)
	})
}
