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

package executor

import (
	"context"
	"sync"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/pulse"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func TestJetTreeUpdater_otherNodesForPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	jc := jet.NewCoordinatorMock(mc)
	ans := node.NewAccessorMock(mc)
	js := jet.NewStorageMock(mc)
	jtu := &fetcher{
		Nodes:       ans,
		JetStorage:  js,
		coordinator: jc,
	}

	t.Run("active light nodes storage returns error", func(t *testing.T) {
		ans.InRoleMock.Expect(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			nil, errors.New("some"),
		)

		nodes, err := jtu.nodesForPulse(ctx, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Empty(t, nodes)
	})

	meRef := gen.Reference()
	jc.MeMock.Return(meRef)

	t.Run("no active nodes at all", func(t *testing.T) {
		ans.InRoleMock.Expect(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			[]insolar.Node{}, nil,
		)

		nodes, err := jtu.nodesForPulse(ctx, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Empty(t, nodes)
	})

	getNodes := func(refs ...insolar.Reference) []insolar.Node {
		res := []insolar.Node{}
		for _, ref := range refs {
			res = append(res, insolar.Node{ID: ref})
		}
		return res
	}

	t.Run("one active node, it's me", func(t *testing.T) {
		ans.InRoleMock.Expect(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			getNodes(meRef), nil,
		)

		nodes, err := jtu.nodesForPulse(ctx, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Empty(t, nodes)
	})

	t.Run("active node", func(t *testing.T) {
		someNode := insolar.Node{ID: gen.Reference()}
		ans.InRoleMock.Expect(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			getNodes(someNode.ID), nil,
		)

		nodes, err := jtu.nodesForPulse(ctx, insolar.PulseNumber(100))
		require.NoError(t, err)
		require.Contains(t, nodes, someNode)
	})

	t.Run("active node and me", func(t *testing.T) {
		meNode := insolar.Node{ID: meRef}
		someNode := insolar.Node{ID: gen.Reference()}

		ans.InRoleMock.Expect(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			getNodes(meNode.ID, someNode.ID), nil,
		)

		nodes, err := jtu.nodesForPulse(ctx, insolar.PulseNumber(100))
		require.NoError(t, err)
		require.Contains(t, nodes, someNode)
		require.NotContains(t, nodes, meNode)
	})
}

func TestJetTreeUpdater_fetchActualJetFromOtherNodes(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	meRef := gen.Reference()
	expectJetID := insolar.ID(*insolar.NewJetID(0, nil))

	js := jet.NewStorageMock(mc)
	jc := jet.NewCoordinatorMock(mc)
	jc.MeMock.Return(meRef)
	bm := bus.NewSenderMock(mc)

	initNodes := func(mc *minimock.Controller) node.Accessor {
		ans := node.NewAccessorMock(mc)
		ans.InRoleMock.Expect(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			[]insolar.Node{
				{ID: gen.Reference()},
				{ID: meRef},
			}, nil,
		)

		return ans
	}

	t.Run("MB error on fetching actual jet", func(t *testing.T) {
		target := gen.ID()
		jtu := &fetcher{
			Nodes:       initNodes(mc),
			JetStorage:  js,
			coordinator: jc,
			sender:      bm,
		}

		bm.SendTargetMock.Set(func(_ context.Context, msg *message.Message, target insolar.Reference) (<-chan *message.Message, func()) {
			res := make(chan *message.Message)
			close(res)
			return res, func() {}
		})

		jetID, err := jtu.fetch(ctx, target, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Nil(t, jetID)
	})

	t.Run("MB got one not actual jet", func(t *testing.T) {
		objectID := gen.ID()
		jtu := &fetcher{
			Nodes:       initNodes(mc),
			JetStorage:  js,
			coordinator: jc,
			sender:      bm,
		}

		bm.SendTargetMock.Set(func(_ context.Context, msg *message.Message, node insolar.Reference) (<-chan *message.Message, func()) {
			getJet := payload.GetJet{}
			err := getJet.Unmarshal(msg.Payload)
			require.NoError(t, err)

			require.Equal(t, objectID, getJet.ObjectID)

			reqMsg, err := payload.NewMessage(&payload.Jet{
				JetID:  gen.JetID(),
				Actual: false,
			})
			require.NoError(t, err)

			meta := payload.Meta{Payload: reqMsg.Payload}
			buf, err := meta.Marshal()
			require.NoError(t, err)
			reqMsg.Payload = buf
			ch := make(chan *message.Message, 1)
			ch <- reqMsg
			return ch, func() {}
		})

		jetID, err := jtu.fetch(ctx, objectID, insolar.PulseNumber(100))
		require.Error(t, err)
		require.Nil(t, jetID)
	})

	t.Run("MB got one actual jet ( from light )", func(t *testing.T) {
		objectID := gen.ID()
		jtu := &fetcher{
			Nodes:       initNodes(mc),
			JetStorage:  js,
			coordinator: jc,
			sender:      bm,
		}

		expectedJetID := insolar.NewJetID(0, nil)

		bm.SendTargetMock.Set(func(_ context.Context, msg *message.Message, node insolar.Reference) (<-chan *message.Message, func()) {
			getJet := payload.GetJet{}
			err := getJet.Unmarshal(msg.Payload)
			require.NoError(t, err)

			require.Equal(t, objectID, getJet.ObjectID)

			reqMsg, err := payload.NewMessage(&payload.Jet{
				JetID:  *expectedJetID,
				Actual: true,
			})
			require.NoError(t, err)

			meta := payload.Meta{Payload: reqMsg.Payload}
			buf, err := meta.Marshal()
			require.NoError(t, err)
			reqMsg.Payload = buf
			ch := make(chan *message.Message, 1)
			ch <- reqMsg
			return ch, func() {}
		})

		jetID, err := jtu.fetch(ctx, objectID, insolar.PulseNumber(100))
		require.NoError(t, err)
		require.Equal(t, insolar.ID(*expectedJetID), *jetID)
	})

	t.Run("MB got one actual jet ( from other light )", func(t *testing.T) {
		ans := node.NewAccessorMock(mc)
		objectID := gen.ID()
		target := insolar.NewReference(objectID)
		ans.InRoleMock.Expect(
			100, insolar.StaticRoleLightMaterial,
		).Return(
			[]insolar.Node{
				{ID: *target},
			}, nil,
		)

		jtu := &fetcher{
			Nodes:       ans,
			JetStorage:  js,
			coordinator: jc,
			sender:      bm,
		}

		expectedJetID := insolar.NewJetID(0, nil)

		bm.SendTargetMock.Set(func(_ context.Context, msg *message.Message, node insolar.Reference) (<-chan *message.Message, func()) {
			getJet := payload.GetJet{}
			err := getJet.Unmarshal(msg.Payload)
			require.NoError(t, err)

			require.Equal(t, *target, node, "send to other target")
			require.Equal(t, objectID, getJet.ObjectID)

			reqMsg, err := payload.NewMessage(&payload.Jet{
				JetID:  *expectedJetID,
				Actual: true,
			})
			require.NoError(t, err)

			meta := payload.Meta{Payload: reqMsg.Payload}
			buf, err := meta.Marshal()
			require.NoError(t, err)
			reqMsg.Payload = buf
			ch := make(chan *message.Message, 1)
			ch <- reqMsg
			return ch, func() {}
		})

		jetID, err := jtu.fetch(ctx, objectID, insolar.PulseNumber(100))
		require.NoError(t, err)
		require.Equal(t, expectJetID, *jetID)
	})

	// TODO: multiple nodes returned different results
	// TODO: multiple nodes returned the same result
}

func TestJetTreeUpdater_fetchJet(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	jc := jet.NewCoordinatorMock(mc)
	ans := node.NewAccessorMock(mc)
	js := jet.NewStorageMock(mc)
	bm := bus.NewSenderMock(mc)
	jtu := &fetcher{
		Nodes:       ans,
		JetStorage:  js,
		coordinator: jc,
		sender:      bm,
		sequencer:   map[seqKey]*seqEntry{},
	}

	target := gen.ID()

	t.Run("quick reply, data is up to date", func(t *testing.T) {
		fjmr := *insolar.NewJetID(0, nil)
		js.ForIDMock.Return(fjmr, true)
		jetID, err := jtu.Fetch(ctx, target, pulse.MinTimePulse+insolar.PulseNumber(100))
		require.NoError(t, err)
		require.Equal(t, fjmr, insolar.JetID(*jetID))
	})

	t.Run("fetch jet from friends", func(t *testing.T) {
		meRef := gen.Reference()
		jc.MeMock.Return(meRef)

		getNodes := func() []insolar.Node {
			return []insolar.Node{{ID: gen.Reference()}}
		}

		ans.InRoleMock.Expect(
			pulse.MinTimePulse+100, insolar.StaticRoleLightMaterial,
		).Return(
			getNodes(), nil,
		)

		expectedJetID := insolar.NewJetID(0, nil)

		bm.SendTargetMock.Set(func(_ context.Context, msg *message.Message, node insolar.Reference) (<-chan *message.Message, func()) {
			getJet := payload.GetJet{}
			err := getJet.Unmarshal(msg.Payload)
			require.NoError(t, err)

			reqMsg, err := payload.NewMessage(&payload.Jet{
				JetID:  *expectedJetID,
				Actual: true,
			})
			require.NoError(t, err)

			meta := payload.Meta{Payload: reqMsg.Payload}
			buf, err := meta.Marshal()
			require.NoError(t, err)
			reqMsg.Payload = buf
			ch := make(chan *message.Message, 1)
			ch <- reqMsg
			return ch, func() {}
		})

		fjm := *insolar.NewJetID(0, nil)
		js.ForIDMock.Return(fjm, false)
		js.UpdateMock.Set(func(ctx context.Context, pn insolar.PulseNumber, actual bool, jets ...insolar.JetID) error {
			require.Equal(t, pulse.MinTimePulse+insolar.PulseNumber(100), pn)
			require.True(t, actual)
			require.Equal(t, []insolar.JetID{*insolar.NewJetID(0, nil)}, jets)
			return nil
		})

		jetID, err := jtu.Fetch(ctx, target, pulse.MinTimePulse+insolar.PulseNumber(100))
		require.NoError(t, err)
		require.Equal(t, insolar.ID(*insolar.NewJetID(0, nil)), *jetID)
	})
}

func TestJetTreeUpdater_Concurrency(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	jc := jet.NewCoordinatorMock(mc)
	ans := node.NewAccessorMock(mc)
	js := jet.NewStorageMock(mc)
	// mb := testutils.NewMessageBusMock(mc)
	bm := bus.NewSenderMock(mc)
	jtu := &fetcher{
		Nodes:       ans,
		JetStorage:  js,
		coordinator: jc,
		sender:      bm,
		sequencer:   map[seqKey]*seqEntry{},
	}

	meRef := gen.Reference()
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

	bm.SendTargetMock.Set(func(_ context.Context, msg *message.Message, target insolar.Reference) (<-chan *message.Message, func()) {
		dataMu.Lock()
		defer dataMu.Unlock()

		getJet := payload.GetJet{}
		err := getJet.Unmarshal(msg.Payload)
		require.NoError(t, err)

		b := getJet.ObjectID.Bytes()[0]
		id := insolar.JetID(*data[b])

		reqMsg, err := payload.NewMessage(&payload.Jet{
			JetID:  id,
			Actual: true,
		})
		require.NoError(t, err)

		meta := payload.Meta{Payload: reqMsg.Payload}
		buf, err := meta.Marshal()
		require.NoError(t, err)
		reqMsg.Payload = buf
		ch := make(chan *message.Message, 1)
		ch <- reqMsg
		return ch, func() {}
	})

	i := 100
	for i > 0 {
		i--

		treeMu := sync.Mutex{}
		tree := jet.NewTree(false)

		js.UpdateMock.Set(func(ctx context.Context, pn insolar.PulseNumber, actual bool, jets ...insolar.JetID) error {
			treeMu.Lock()
			defer treeMu.Unlock()

			for _, id := range jets {
				tree.Update(id, actual)
			}
			return nil
		})
		js.ForIDMock.Set(func(ctx context.Context, pulse insolar.PulseNumber, id insolar.ID) (insolar.JetID, bool) {
			treeMu.Lock()
			defer treeMu.Unlock()

			return tree.Find(id)
		})

		wg := sync.WaitGroup{}
		wg.Add(4)

		for _, b := range []byte{0, 128, 192} {
			go func(b byte) {
				target := insolar.NewID(pulse.MinTimePulse+50, []byte{b})

				jetID, err := jtu.Fetch(ctx, *target, pulse.MinTimePulse+insolar.PulseNumber(100))
				require.NoError(t, err)

				dataMu.Lock()
				require.Equal(t, data[b], jetID)
				dataMu.Unlock()

				wg.Done()
			}(b)
		}
		go func() {
			dataMu.Lock()
			jtu.Fetch(ctx, *data[128], pulse.MinTimePulse+insolar.PulseNumber(100))
			dataMu.Unlock()

			wg.Done()
		}()
		wg.Wait()
	}
}
