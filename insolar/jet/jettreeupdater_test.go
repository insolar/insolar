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

package jet

import (
	"context"
	"sync"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/ledger/storage/node"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

func TestJetTreeUpdater_otherNodesForPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	jc := testutils.NewJetCoordinatorMock(mc)
	ans := node.NewAccessorMock(mc)
	js := NewStorageMock(mc)
	jtu := &treeUpdater{
		Nodes:          ans,
		JetStorage:     js,
		JetCoordinator: jc,
	}

	t.Run("active nodes storage returns error", func(t *testing.T) {
		ans.InRoleMock.Expect(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			nil, errors.New("some"),
		)

		nodes, err := jtu.otherNodesForPulse(ctx, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Empty(t, nodes)
	})

	meRef := testutils.RandomRef()
	jc.MeMock.Return(meRef)

	t.Run("no active nodes at all", func(t *testing.T) {

		ans.InRoleMock.Expect(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			[]insolar.Node{}, nil,
		)

		nodes, err := jtu.otherNodesForPulse(ctx, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Empty(t, nodes)
	})

	t.Run("one active node, it's me", func(t *testing.T) {
		ans.InRoleMock.Expect(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			[]insolar.Node{{ID: meRef}}, nil,
		)

		nodes, err := jtu.otherNodesForPulse(ctx, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Empty(t, nodes)
	})

	t.Run("active node", func(t *testing.T) {
		someNode := insolar.Node{ID: gen.Reference()}
		ans.InRoleMock.Expect(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			[]insolar.Node{someNode}, nil,
		)

		nodes, err := jtu.otherNodesForPulse(ctx, insolar.PulseNumber(100))
		require.NoError(t, err)
		require.Contains(t, nodes, someNode)
	})

	t.Run("active node and me", func(t *testing.T) {
		meNode := insolar.Node{ID: meRef}
		someNode := insolar.Node{ID: gen.Reference()}

		ans.InRoleMock.Expect(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			[]insolar.Node{someNode, meNode}, nil,
		)

		nodes, err := jtu.otherNodesForPulse(ctx, insolar.PulseNumber(100))
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
	ans := node.NewAccessorMock(mc)
	js := NewStorageMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	jtu := &treeUpdater{
		Nodes:          ans,
		JetStorage:     js,
		JetCoordinator: jc,
		MessageBus:     mb,
	}

	meRef := testutils.RandomRef()
	jc.MeMock.Return(meRef)

	ans.InRoleMock.Expect(
		100, insolar.StaticRoleLightMaterial,
	).Return(
		[]insolar.Node{{ID: gen.Reference()}}, nil,
	)

	t.Run("MB error on fetching actual jet", func(t *testing.T) {
		target := testutils.RandomID()

		mb.SendMock.Return(nil, errors.New("some"))

		jetID, err := jtu.fetchActualJetFromOtherNodes(ctx, target, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Nil(t, jetID)
	})

	t.Run("MB got one not actual jet", func(t *testing.T) {
		target := testutils.RandomID()

		mb.SendMock.Return(
			&reply.Jet{ID: insolar.ID(*insolar.NewJetID(0, nil)), Actual: false},
			nil,
		)

		jetID, err := jtu.fetchActualJetFromOtherNodes(ctx, target, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Nil(t, jetID)
	})
	t.Run("MB got one actual jet", func(t *testing.T) {
		target := testutils.RandomID()

		mb.SendMock.Return(
			&reply.Jet{ID: insolar.ID(*insolar.NewJetID(0, nil)), Actual: true},
			nil,
		)

		jetID, err := jtu.fetchActualJetFromOtherNodes(ctx, target, insolar.PulseNumber(100))
		require.NoError(t, err)
		require.Equal(t, insolar.ID(*insolar.NewJetID(0, nil)), *jetID)
	})

	// TODO: multiple nodes returned different results
	// TODO: multiple nodes returned the same result
}

func TestJetTreeUpdater_fetchJet(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	jc := testutils.NewJetCoordinatorMock(mc)
	ans := node.NewAccessorMock(mc)
	js := NewStorageMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	jtu := &treeUpdater{
		Nodes:          ans,
		JetStorage:     js,
		JetCoordinator: jc,
		MessageBus:     mb,
		sequencer:      map[seqKey]*seqEntry{},
	}

	target := testutils.RandomID()

	t.Run("quick reply, data is up to date", func(t *testing.T) {
		fjmr := *insolar.NewJetID(0, nil)
		js.ForIDMock.Return(fjmr, true)
		jetID, err := jtu.FetchJet(ctx, target, insolar.PulseNumber(100))
		require.NoError(t, err)
		require.Equal(t, fjmr, insolar.JetID(*jetID))
	})

	t.Run("fetch jet from friends", func(t *testing.T) {
		meRef := testutils.RandomRef()
		jc.MeMock.Return(meRef)

		ans.InRoleMock.Expect(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			[]insolar.Node{{ID: gen.Reference()}}, nil,
		)
		mb.SendMock.Return(
			&reply.Jet{ID: insolar.ID(*insolar.NewJetID(0, nil)), Actual: true},
			nil,
		)

		fjm := *insolar.NewJetID(0, nil)
		js.ForIDMock.Return(fjm, false)
		js.UpdateFunc = func(ctx context.Context, pn insolar.PulseNumber, actual bool, jets ...insolar.JetID) {
			require.Equal(t, insolar.PulseNumber(100), pn)
			require.True(t, actual)
			require.Equal(t, []insolar.JetID{*insolar.NewJetID(0, nil)}, jets)
		}

		jetID, err := jtu.FetchJet(ctx, target, insolar.PulseNumber(100))
		require.NoError(t, err)
		require.Equal(t, insolar.ID(*insolar.NewJetID(0, nil)), *jetID)
	})
}

func TestJetTreeUpdater_Concurrency(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	jc := testutils.NewJetCoordinatorMock(mc)
	ans := node.NewAccessorMock(mc)
	js := NewStorageMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	jtu := &treeUpdater{
		Nodes:          ans,
		JetStorage:     js,
		JetCoordinator: jc,
		MessageBus:     mb,
		sequencer:      map[seqKey]*seqEntry{},
	}

	meRef := testutils.RandomRef()
	jc.MeMock.Return(meRef)

	nodes := []insolar.Node{{ID: gen.Reference()}, {ID: gen.Reference()}, {ID: gen.Reference()}}

	ans.InRoleMock.Return(nodes, nil)

	dataMu := sync.Mutex{}

	first := insolar.ID(*insolar.NewJetID(2, []byte{0}))
	second := insolar.ID(*insolar.NewJetID(2, []byte{0}))
	third := insolar.ID(*insolar.NewJetID(2, []byte{0}))
	fourth := insolar.ID(*insolar.NewJetID(2, []byte{0}))

	data := map[byte]*insolar.ID{
		0:   &first,  // 00
		128: &second, // 10
		64:  &third,  // 01
		192: &fourth, // 11
	}

	mb.SendFunc = func(ctx context.Context, msg insolar.Message, opt *insolar.MessageSendOptions) (insolar.Reply, error) {
		dataMu.Lock()
		defer dataMu.Unlock()

		b := msg.(*message.GetJet).Object.Bytes()[0]
		return &reply.Jet{ID: *data[b], Actual: true}, nil
	}

	i := 100
	for i > 0 {
		i--

		treeMu := sync.Mutex{}
		tree := NewTree(false)

		js.UpdateFunc = func(ctx context.Context, pn insolar.PulseNumber, actual bool, jets ...insolar.JetID) {
			treeMu.Lock()
			defer treeMu.Unlock()

			for _, id := range jets {
				tree.Update(id, actual)
			}
		}
		js.ForIDFunc = func(ctx context.Context, pulse insolar.PulseNumber, id insolar.ID) (insolar.JetID, bool) {
			treeMu.Lock()
			defer treeMu.Unlock()

			return tree.Find(id)
		}

		wg := sync.WaitGroup{}
		wg.Add(4)

		for _, b := range []byte{0, 128, 192} {
			go func(b byte) {
				target := insolar.NewID(0, []byte{b})

				jetID, err := jtu.FetchJet(ctx, *target, insolar.PulseNumber(100))
				require.NoError(t, err)

				dataMu.Lock()
				require.Equal(t, data[b], jetID)
				dataMu.Unlock()

				wg.Done()
			}(b)
		}
		go func() {
			dataMu.Lock()
			jtu.ReleaseJet(ctx, *data[128], insolar.PulseNumber(100))
			dataMu.Unlock()

			wg.Done()
		}()
		wg.Wait()
	}
}
