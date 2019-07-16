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
	"github.com/insolar/insolar/insolar/node"
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

	jc := NewCoordinatorMock(mc)
	ans := node.NewAccessorMock(mc)
	js := NewStorageMock(mc)
	jtu := &fetcher{
		Nodes:       ans,
		JetStorage:  js,
		coordinator: jc,
	}

	t.Run("active light nodes storage returns error", func(t *testing.T) {
		ans.InRoleMock.ExpectOnce(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			nil, errors.New("some"),
		)

		nodes, err := jtu.nodesForPulse(ctx, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Empty(t, nodes)
	})

	t.Run("active heavy nodes storage returns error", func(t *testing.T) {
		ans.InRoleMock.ExpectOnce(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			[]insolar.Node{insolar.Node{}}, nil,
		)

		ans.InRoleMock.ExpectOnce(
			100, insolar.StaticRoleHeavyMaterial,
		).Return(
			nil, errors.New("some"),
		)

		nodes, err := jtu.nodesForPulse(ctx, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Empty(t, nodes)
	})

	meRef := testutils.RandomRef()
	jc.MeMock.Return(meRef)

	t.Run("no active nodes at all", func(t *testing.T) {

		ans.InRoleMock.ExpectOnce(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			[]insolar.Node{}, nil,
		)

		ans.InRoleMock.ExpectOnce(
			100, insolar.StaticRoleHeavyMaterial,
		).Return(
			[]insolar.Node{}, nil,
		)

		nodes, err := jtu.nodesForPulse(ctx, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Empty(t, nodes)
	})

	ttable := []struct{ name string }{{name: "Light"}, {name: "Heavy"}}
	getNodes := func(tname string, isFromLight bool, refs ...insolar.Reference) []insolar.Node {
		res := []insolar.Node{}
		for _, ref := range refs {
			res = append(res, insolar.Node{ID: ref})
		}
		if tname == "Light" && isFromLight {
			return res
		}
		if tname == "Heavy" && !isFromLight {
			return res
		}
		return []insolar.Node{}
	}

	t.Run("one active node, it's me", func(t *testing.T) {
		for _, tt := range ttable {
			t.Run(tt.name, func(t *testing.T) {
				ans.InRoleMock.ExpectOnce(
					100, insolar.StaticRoleLightMaterial,
				).Return(
					getNodes(tt.name, true, meRef), nil,
				)

				ans.InRoleMock.ExpectOnce(
					100, insolar.StaticRoleHeavyMaterial,
				).Return(
					getNodes(tt.name, false, meRef), nil,
				)

				nodes, err := jtu.nodesForPulse(ctx, insolar.PulseNumber(100))
				require.Error(t, err)
				require.Empty(t, nodes)
			})
		}

	})

	t.Run("active node", func(t *testing.T) {
		someNode := insolar.Node{ID: gen.Reference()}
		for _, tt := range ttable {
			t.Run(tt.name, func(t *testing.T) {
				ans.InRoleMock.ExpectOnce(
					100, insolar.StaticRoleLightMaterial,
				).Return(
					getNodes(tt.name, true, someNode.ID), nil,
				)

				ans.InRoleMock.ExpectOnce(
					100, insolar.StaticRoleHeavyMaterial,
				).Return(
					getNodes(tt.name, false, someNode.ID), nil,
				)

				nodes, err := jtu.nodesForPulse(ctx, insolar.PulseNumber(100))
				require.NoError(t, err)
				require.Contains(t, nodes, someNode)
			})
		}
	})

	t.Run("active node and me", func(t *testing.T) {
		meNode := insolar.Node{ID: meRef}
		someNode := insolar.Node{ID: gen.Reference()}

		for _, tt := range ttable {
			t.Run(tt.name, func(t *testing.T) {
				ans.InRoleMock.ExpectOnce(
					100, insolar.StaticRoleLightMaterial,
				).Return(
					getNodes(tt.name, true, meNode.ID, someNode.ID), nil,
				)

				ans.InRoleMock.ExpectOnce(
					100, insolar.StaticRoleHeavyMaterial,
				).Return(
					getNodes(tt.name, false, meNode.ID, someNode.ID), nil,
				)

				nodes, err := jtu.nodesForPulse(ctx, insolar.PulseNumber(100))
				require.NoError(t, err)
				require.Contains(t, nodes, someNode)
				require.NotContains(t, nodes, meNode)
			})
		}
	})

	t.Run("active nodes (heavy and light)", func(t *testing.T) {
		someNode1 := insolar.Node{ID: gen.Reference()}
		someNode2 := insolar.Node{ID: gen.Reference()}

		ans.InRoleMock.ExpectOnce(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			[]insolar.Node{someNode1}, nil,
		)
		ans.InRoleMock.ExpectOnce(
			100, insolar.StaticRoleHeavyMaterial,
		).Return(
			[]insolar.Node{someNode2}, nil,
		)

		nodes, err := jtu.nodesForPulse(ctx, insolar.PulseNumber(100))
		require.NoError(t, err)
		require.Len(t, nodes, 2)
		require.Contains(t, nodes, someNode1)
		require.Contains(t, nodes, someNode2)
	})
}

func TestJetTreeUpdater_fetchActualJetFromOtherNodes(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	jc := NewCoordinatorMock(mc)

	js := NewStorageMock(mc)
	mb := testutils.NewMessageBusMock(mc)

	meRef := testutils.RandomRef()
	jc.MeMock.Return(meRef)

	initNodes := func(mc *minimock.Controller) node.Accessor {
		ans := node.NewAccessorMock(mc)
		ans.InRoleMock.ExpectOnce(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			[]insolar.Node{{ID: gen.Reference()}}, nil,
		)

		ans.InRoleMock.ExpectOnce(
			100, insolar.StaticRoleHeavyMaterial,
		).Return(
			[]insolar.Node{}, nil,
		)

		return ans
	}

	t.Run("MB error on fetching actual jet", func(t *testing.T) {
		target := testutils.RandomID()
		jtu := &fetcher{
			Nodes:       initNodes(mc),
			JetStorage:  js,
			coordinator: jc,
			MessageBus:  mb,
		}

		mb.SendMock.Return(nil, errors.New("some"))

		jetID, err := jtu.fetch(ctx, target, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Nil(t, jetID)
	})

	t.Run("MB got one not actual jet", func(t *testing.T) {
		target := testutils.RandomID()
		jtu := &fetcher{
			Nodes:       initNodes(mc),
			JetStorage:  js,
			coordinator: jc,
			MessageBus:  mb,
		}

		mb.SendMock.Return(
			&reply.Jet{ID: insolar.ID(*insolar.NewJetID(0, nil)), Actual: false},
			nil,
		)

		jetID, err := jtu.fetch(ctx, target, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Nil(t, jetID)
	})

	t.Run("MB got one actual jet ( from light )", func(t *testing.T) {
		target := testutils.RandomID()
		jtu := &fetcher{
			Nodes:       initNodes(mc),
			JetStorage:  js,
			coordinator: jc,
			MessageBus:  mb,
		}

		mb.SendMock.Return(
			&reply.Jet{ID: insolar.ID(*insolar.NewJetID(0, nil)), Actual: true},
			nil,
		)

		jetID, err := jtu.fetch(ctx, target, insolar.PulseNumber(100))
		require.NoError(t, err)
		require.Equal(t, insolar.ID(*insolar.NewJetID(0, nil)), *jetID)
	})

	t.Run("MB got one actual jet ( from heavy )", func(t *testing.T) {
		ans := node.NewAccessorMock(mc)
		target := testutils.RandomID()
		ans.InRoleMock.ExpectOnce(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			[]insolar.Node{}, nil,
		)

		ans.InRoleMock.ExpectOnce(
			100, insolar.StaticRoleHeavyMaterial,
		).Return(
			[]insolar.Node{{ID: gen.Reference()}}, nil,
		)
		jtu := &fetcher{
			Nodes:       ans,
			JetStorage:  js,
			coordinator: jc,
			MessageBus:  mb,
		}

		mb.SendMock.Return(
			&reply.Jet{ID: insolar.ID(*insolar.NewJetID(0, nil)), Actual: true},
			nil,
		)

		jetID, err := jtu.fetch(ctx, target, insolar.PulseNumber(100))
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

	jc := NewCoordinatorMock(mc)
	ans := node.NewAccessorMock(mc)
	js := NewStorageMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	jtu := &fetcher{
		Nodes:       ans,
		JetStorage:  js,
		coordinator: jc,
		MessageBus:  mb,
		sequencer:   map[seqKey]*seqEntry{},
	}

	target := testutils.RandomID()

	t.Run("quick reply, data is up to date", func(t *testing.T) {
		fjmr := *insolar.NewJetID(0, nil)
		js.ForIDMock.Return(fjmr, true)
		jetID, err := jtu.Fetch(ctx, target, insolar.FirstPulseNumber+insolar.PulseNumber(100))
		require.NoError(t, err)
		require.Equal(t, fjmr, insolar.JetID(*jetID))
	})

	t.Run("fetch jet from friends", func(t *testing.T) {
		meRef := testutils.RandomRef()
		jc.MeMock.Return(meRef)

		getNodes := func(tname string, isFromLight bool) []insolar.Node {
			if tname == "Light" && isFromLight {
				return []insolar.Node{{ID: gen.Reference()}}
			}
			if tname == "Heavy" && !isFromLight {
				return []insolar.Node{{ID: gen.Reference()}}
			}
			return []insolar.Node{}
		}

		ttable := []struct{ name string }{{name: "Light"}, {name: "Heavy"}}

		for _, tt := range ttable {
			t.Run(tt.name, func(t *testing.T) {
				ans.InRoleMock.ExpectOnce(
					insolar.FirstPulseNumber+100, insolar.StaticRoleLightMaterial,
				).Return(
					getNodes(tt.name, true), nil,
				)
				ans.InRoleMock.ExpectOnce(
					insolar.FirstPulseNumber+100, insolar.StaticRoleHeavyMaterial,
				).Return(
					getNodes(tt.name, false), nil,
				)
				mb.SendMock.Return(
					&reply.Jet{ID: insolar.ID(*insolar.NewJetID(0, nil)), Actual: true},
					nil,
				)

				fjm := *insolar.NewJetID(0, nil)
				js.ForIDMock.Return(fjm, false)
				js.UpdateFunc = func(ctx context.Context, pn insolar.PulseNumber, actual bool, jets ...insolar.JetID) error {
					require.Equal(t, insolar.FirstPulseNumber+insolar.PulseNumber(100), pn)
					require.True(t, actual)
					require.Equal(t, []insolar.JetID{*insolar.NewJetID(0, nil)}, jets)
					return nil
				}

				jetID, err := jtu.Fetch(ctx, target, insolar.FirstPulseNumber+insolar.PulseNumber(100))
				require.NoError(t, err)
				require.Equal(t, insolar.ID(*insolar.NewJetID(0, nil)), *jetID)
			})
		}

	})
}

func TestJetTreeUpdater_Concurrency(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	jc := NewCoordinatorMock(mc)
	ans := node.NewAccessorMock(mc)
	js := NewStorageMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	jtu := &fetcher{
		Nodes:       ans,
		JetStorage:  js,
		coordinator: jc,
		MessageBus:  mb,
		sequencer:   map[seqKey]*seqEntry{},
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

		js.UpdateFunc = func(ctx context.Context, pn insolar.PulseNumber, actual bool, jets ...insolar.JetID) error {
			treeMu.Lock()
			defer treeMu.Unlock()

			for _, id := range jets {
				tree.Update(id, actual)
			}
			return nil
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

				jetID, err := jtu.Fetch(ctx, *target, insolar.FirstPulseNumber+insolar.PulseNumber(100))
				require.NoError(t, err)

				dataMu.Lock()
				require.Equal(t, data[b], jetID)
				dataMu.Unlock()

				wg.Done()
			}(b)
		}
		go func() {
			dataMu.Lock()
			jtu.Fetch(ctx, *data[128], insolar.FirstPulseNumber+insolar.PulseNumber(100))
			dataMu.Unlock()

			wg.Done()
		}()
		wg.Wait()
	}
}
