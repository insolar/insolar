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
	"github.com/insolar/insolar/ledger/storage/nodes"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
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
	ans := nodes.NewAccessorMock(mc)
	js := storage.NewJetStorageMock(mc)
	jtu := &jetTreeUpdater{
		Nodes:          ans,
		JetStorage:     js,
		JetCoordinator: jc,
	}

	t.Run("active nodes storage returns error", func(t *testing.T) {
		ans.InRoleMock.Expect(
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

		ans.InRoleMock.Expect(
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

		ans.InRoleMock.Expect(
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

		ans.InRoleMock.Expect(
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

		ans.InRoleMock.Expect(
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
	ans := nodes.NewAccessorMock(mc)
	js := storage.NewJetStorageMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	jtu := &jetTreeUpdater{
		Nodes:          ans,
		JetStorage:     js,
		JetCoordinator: jc,
		MessageBus:     mb,
	}

	meRef := testutils.RandomRef()
	jc.MeMock.Return(meRef)

	otherNode := network.NewNodeMock(mc)
	otherNode.IDMock.Return(testutils.RandomRef())

	ans.InRoleMock.Expect(
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
	ans := nodes.NewAccessorMock(mc)
	js := storage.NewJetStorageMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	jtu := &jetTreeUpdater{
		Nodes:          ans,
		JetStorage:     js,
		JetCoordinator: jc,
		MessageBus:     mb,
		sequencer:      map[seqKey]*seqEntry{},
	}

	target := testutils.RandomID()

	t.Run("quick reply, data is up to date", func(t *testing.T) {
		js.FindJetMock.Return(jet.NewID(0, nil), true)
		jetID, err := jtu.fetchJet(ctx, target, core.PulseNumber(100))
		require.NoError(t, err)
		require.Equal(t, jet.NewID(0, nil), jetID)
	})

	t.Run("fetch jet from friends", func(t *testing.T) {
		meRef := testutils.RandomRef()
		jc.MeMock.Return(meRef)

		otherNode := network.NewNodeMock(mc)
		otherNode.IDMock.Return(testutils.RandomRef())

		ans.InRoleMock.Expect(
			100, core.StaticRoleLightMaterial,
		).Return(
			[]core.Node{otherNode}, nil,
		)
		mb.SendMock.Return(
			&reply.Jet{ID: *jet.NewID(0, nil), Actual: true},
			nil,
		)

		js.FindJetMock.Return(jet.NewID(0, nil), false)
		js.UpdateJetTreeFunc = func(ctx context.Context, pn core.PulseNumber, actual bool, jets ...core.RecordID) {
			require.Equal(t, core.PulseNumber(100), pn)
			require.True(t, actual)
			require.Equal(t, []core.RecordID{*jet.NewID(0, nil)}, jets)
		}

		jetID, err := jtu.fetchJet(ctx, target, core.PulseNumber(100))
		require.NoError(t, err)
		require.Equal(t, jet.NewID(0, nil), jetID)
	})
}

func TestJetTreeUpdater_Concurrency(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	jc := testutils.NewJetCoordinatorMock(mc)
	ans := nodes.NewAccessorMock(mc)
	js := storage.NewJetStorageMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	jtu := &jetTreeUpdater{
		Nodes:          ans,
		JetStorage:     js,
		JetCoordinator: jc,
		MessageBus:     mb,
		sequencer:      map[seqKey]*seqEntry{},
	}

	meRef := testutils.RandomRef()
	jc.MeMock.Return(meRef)

	node := network.NewNodeMock(mc)
	node.IDMock.Return(testutils.RandomRef())
	nodes := []core.Node{node, node, node}

	ans.InRoleMock.Return(nodes, nil)

	dataMu := sync.Mutex{}
	data := map[byte]*core.RecordID{
		0:   jet.NewID(2, []byte{0}), // 00
		128: jet.NewID(2, []byte{0}), // 10
		64:  jet.NewID(2, []byte{0}), // 01
		192: jet.NewID(2, []byte{0}), // 11
	}

	mb.SendFunc = func(ctx context.Context, msg core.Message, opt *core.MessageSendOptions) (core.Reply, error) {
		dataMu.Lock()
		defer dataMu.Unlock()

		b := msg.(*message.GetJet).Object.Bytes()[0]
		return &reply.Jet{ID: *data[b], Actual: true}, nil
	}

	i := 100
	for i > 0 {
		i--

		treeMu := sync.Mutex{}
		tree := jet.NewTree(false)

		js.UpdateJetTreeFunc = func(ctx context.Context, pn core.PulseNumber, actual bool, jets ...core.RecordID) {
			treeMu.Lock()
			defer treeMu.Unlock()

			for _, id := range jets {
				tree.Update(id, actual)
			}
		}
		js.FindJetFunc = func(ctx context.Context, pulse core.PulseNumber, id core.RecordID) (*core.RecordID, bool) {
			treeMu.Lock()
			defer treeMu.Unlock()

			return tree.Find(id)
		}

		wg := sync.WaitGroup{}
		wg.Add(4)

		for _, b := range []byte{0, 128, 192} {
			go func(b byte) {
				target := core.NewRecordID(0, []byte{b})

				jetID, err := jtu.fetchJet(ctx, *target, core.PulseNumber(100))
				require.NoError(t, err)

				dataMu.Lock()
				require.Equal(t, data[b], jetID)
				dataMu.Unlock()

				wg.Done()
			}(b)
		}
		go func() {
			dataMu.Lock()
			jtu.releaseJet(ctx, *data[128], core.PulseNumber(100))
			dataMu.Unlock()

			wg.Done()
		}()
		wg.Wait()
	}
}
