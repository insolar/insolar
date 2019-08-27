/*
 * Copyright 2019 Insolar Technologies GmbH
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package executor_test

import (
	"bytes"
	"sort"
	"strings"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/stretchr/testify/require"
)

var (
	pending insolar.PulseNumber = 780
	topSync insolar.PulseNumber = 800
	current insolar.PulseNumber = 850
)

func singleJetFixture() []insolar.JetID {
	tree := jet.NewTree(true)
	return tree.LeafIDs()
}

func jetsFixture() []insolar.JetID {
	tree := jet.NewTree(true)
	left, _, _ := tree.Split(tree.LeafIDs()[0])
	_, _, _ = tree.Split(left)
	return tree.LeafIDs()
}

func indexesFixture() []record.Index {
	return []record.Index{
		{
			ObjID: *insolar.NewID(pending, []byte{0x01}), // JetIDs[0]
			Lifeline: record.Lifeline{
				EarliestOpenRequest: &pending,
			},
		},
		{
			ObjID: *insolar.NewID(pending, []byte{0x02}), // JetIDs[0]
			Lifeline: record.Lifeline{
				EarliestOpenRequest: &pending,
			},
		},
		{
			ObjID: *insolar.NewID(pending, []byte{0xD5}), // 11010101 -> JetIDs[2]
			Lifeline: record.Lifeline{
				EarliestOpenRequest: &pending,
			},
		},
		{
			ObjID: *insolar.NewID(pending, []byte{0xD6}), // 11010110 -> JetIDs[2]
			Lifeline: record.Lifeline{
				EarliestOpenRequest: nil,
			},
		},
		{
			ObjID: *insolar.NewID(pending, []byte{0xAC}), // 10101100 -> JetIDs[1]
			Lifeline: record.Lifeline{
				EarliestOpenRequest: nil,
			},
		},
	}
}

func dropsFixture() []drop.Drop {
	ids := gen.UniqueIDs(3)
	return []drop.Drop{
		{Split: false, Pulse: ids[0].Pulse(), Hash: ids[0].Bytes()},
		{Split: true, Pulse: ids[1].Pulse(), Hash: ids[1].Bytes()},
		{Split: false, Pulse: ids[2].Pulse(), Hash: ids[2].Bytes()},
	}
}

func sortDrops(drops []drop.Drop) {
	cmp := func(i, j int) bool {
		return drops[i].Pulse < drops[j].Pulse
	}
	sort.Slice(drops, cmp)
}

func sortJets(jets []insolar.JetID) {
	cmp := func(i, j int) bool {
		cmp := strings.Compare(jets[i].DebugString(), jets[j].DebugString())
		return cmp < 0
	}
	sort.Slice(jets, cmp)
}

func sortIndexes(indexes []record.Index) {
	cmp := func(i, j int) bool {
		cmp := bytes.Compare(indexes[i].ObjID.Bytes(), indexes[j].ObjID.Bytes())
		return cmp < 0
	}
	sort.Slice(indexes, cmp)
}

func TestInitialStateKeeper_Get_AfterRestart(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	ctx := inslogger.TestContext(t)

	jetKeeper := executor.NewJetKeeperMock(mc)
	jetKeeper.TopSyncPulseMock.Return(topSync)
	jetIDs := jetsFixture()
	// split jet depends on fixture data
	left, right := jet.Siblings(jetIDs[1])

	jetAccessor := jet.NewAccessorMock(mc)
	jetAccessor.AllMock.Expect(ctx, topSync).Return(jetIDs)

	drops := dropsFixture()
	dropAccessor := drop.NewAccessorMock(mc)
	dropAccessor.ForPulseMock.When(ctx, jetIDs[0], topSync).Then(drops[0], nil)
	dropAccessor.ForPulseMock.When(ctx, jetIDs[1], topSync).Then(drops[1], nil)
	dropAccessor.ForPulseMock.When(ctx, jetIDs[2], topSync).Then(drops[2], nil)

	indexes := indexesFixture()
	indexAccessor := object.NewIndexAccessorMock(mc)
	indexAccessor.ForPulseMock.Expect(ctx, topSync).Return(indexes, nil)

	jetCoordinator := jet.NewCoordinatorMock(mc)

	stateKeeper := executor.NewInitialStateKeeper(jetKeeper, jetAccessor, jetCoordinator, indexAccessor, dropAccessor)
	err := stateKeeper.Start(ctx)
	require.NoError(t, err)

	currentLight := gen.Reference()
	anotherLight := gen.Reference()

	jetCoordinator.LightExecutorForJetMock.When(ctx, insolar.ID(jetIDs[0]), current).Then(&currentLight, nil)
	jetCoordinator.LightExecutorForJetMock.When(ctx, insolar.ID(left), current).Then(&currentLight, nil)
	jetCoordinator.LightExecutorForJetMock.When(ctx, insolar.ID(right), current).Then(&anotherLight, nil)
	jetCoordinator.LightExecutorForJetMock.When(ctx, insolar.ID(jetIDs[2]), current).Then(&anotherLight, nil)

	// Get for currentLight
	state := stateKeeper.Get(ctx, currentLight, current)

	expectedIndexes := []record.Index{indexes[0], indexes[1]}
	sortIndexes(expectedIndexes)
	sortIndexes(state.Indexes)
	require.Equal(t, expectedIndexes, state.Indexes)

	expectedDrops := []drop.Drop{
		drops[0],
		drops[1],
	}
	sortDrops(expectedDrops)
	sortDrops(state.Drops)
	require.Equal(t, expectedDrops, state.Drops)

	expectedJets := []insolar.JetID{jetIDs[0], left}
	sortJets(expectedJets)
	sortJets(state.JetIDs)
	require.Equal(t, expectedJets, state.JetIDs)

	// Get for anotherLight
	state = stateKeeper.Get(ctx, anotherLight, current)

	expectedIndexes = []record.Index{indexes[2]}
	sortIndexes(expectedIndexes)
	sortIndexes(state.Indexes)
	require.Equal(t, []record.Index{indexes[2]}, state.Indexes)

	expectedDrops = []drop.Drop{
		drops[1],
		drops[2],
	}
	sortDrops(expectedDrops)
	sortDrops(state.Drops)
	require.Equal(t, expectedDrops, state.Drops)

	expectedJets = []insolar.JetID{right, jetIDs[2]}
	sortJets(expectedJets)
	sortJets(state.JetIDs)
	require.Equal(t, expectedJets, state.JetIDs)
}

func TestInitialStateKeeper_Get_EmptyAfterRestart(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	ctx := inslogger.TestContext(t)

	jetKeeper := executor.NewJetKeeperMock(mc)
	jetKeeper.TopSyncPulseMock.Return(topSync)

	jetIDs := singleJetFixture()
	jetAccessor := jet.NewAccessorMock(mc)
	jetAccessor.AllMock.Expect(ctx, topSync).Return(jetIDs)

	jetDrop := drop.Drop{Split: false, Hash: gen.ID().Bytes()}
	dropAccessor := drop.NewAccessorMock(mc)
	dropAccessor.ForPulseMock.When(ctx, jetIDs[0], topSync).Then(jetDrop, nil)

	indexAccessor := object.NewIndexAccessorMock(mc)
	indexAccessor.ForPulseMock.Expect(ctx, topSync).Return(nil, object.ErrIndexNotFound)

	jetCoordinator := jet.NewCoordinatorMock(mc)

	stateKeeper := executor.NewInitialStateKeeper(jetKeeper, jetAccessor, jetCoordinator, indexAccessor, dropAccessor)
	err := stateKeeper.Start(ctx)
	require.NoError(t, err)

	currentLight := gen.Reference()
	anotherLight := gen.Reference()

	jetCoordinator.LightExecutorForJetMock.When(ctx, insolar.ID(jetIDs[0]), current).Then(&currentLight, nil)

	// Get for currentLight
	state := stateKeeper.Get(ctx, currentLight, current)

	require.Equal(t, []record.Index{}, state.Indexes)
	require.Equal(t, []drop.Drop{jetDrop}, state.Drops)
	require.Equal(t, []insolar.JetID{jetIDs[0]}, state.JetIDs)

	// Get for anotherLight
	state = stateKeeper.Get(ctx, anotherLight, current)

	require.Equal(t, []record.Index{}, state.Indexes)
	require.Equal(t, []drop.Drop{}, state.Drops)
	require.Equal(t, []insolar.JetID{}, state.JetIDs)

}

func TestInitialStateKeeper_Get_WithDuplicatedDrops(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	ctx := inslogger.TestContext(t)

	jetKeeper := executor.NewJetKeeperMock(mc)
	jetKeeper.TopSyncPulseMock.Return(topSync)

	jetIDs := jetsFixture()
	// split jet depends on fixture data
	left, right := jet.Siblings(jetIDs[1])

	jetAccessor := jet.NewAccessorMock(mc)
	jetAccessor.AllMock.Expect(ctx, topSync).Return(jetIDs)

	drops := dropsFixture()
	dropAccessor := drop.NewAccessorMock(mc)
	dropAccessor.ForPulseMock.When(ctx, jetIDs[0], topSync).Then(drops[0], nil)
	dropAccessor.ForPulseMock.When(ctx, jetIDs[1], topSync).Then(drops[1], nil)
	dropAccessor.ForPulseMock.When(ctx, jetIDs[2], topSync).Then(drops[2], nil)

	indexes := indexesFixture()
	indexAccessor := object.NewIndexAccessorMock(mc)
	indexAccessor.ForPulseMock.Expect(ctx, topSync).Return(indexes, nil)

	jetCoordinator := jet.NewCoordinatorMock(mc)

	stateKeeper := executor.NewInitialStateKeeper(jetKeeper, jetAccessor, jetCoordinator, indexAccessor, dropAccessor)

	err := stateKeeper.Start(ctx)
	require.NoError(t, err)

	currentLight := gen.Reference()
	anotherLight := gen.Reference()

	// Both jets after split executed by currentLight
	jetCoordinator.LightExecutorForJetMock.When(ctx, insolar.ID(jetIDs[0]), current).Then(&currentLight, nil)
	jetCoordinator.LightExecutorForJetMock.When(ctx, insolar.ID(left), current).Then(&currentLight, nil)
	jetCoordinator.LightExecutorForJetMock.When(ctx, insolar.ID(right), current).Then(&currentLight, nil)
	jetCoordinator.LightExecutorForJetMock.When(ctx, insolar.ID(jetIDs[2]), current).Then(&anotherLight, nil)

	// Get for currentLight
	state := stateKeeper.Get(ctx, currentLight, current)

	expectedIndexes := []record.Index{indexes[0], indexes[1]}
	sortIndexes(expectedIndexes)
	sortIndexes(state.Indexes)
	require.Equal(t, expectedIndexes, state.Indexes)

	expectedDrops := []drop.Drop{
		drops[0],
		drops[1],
	}
	sortDrops(expectedDrops)
	sortDrops(state.Drops)
	require.Equal(t, expectedDrops, state.Drops)

	expectedJets := []insolar.JetID{jetIDs[0], left, right}
	sortJets(expectedJets)
	sortJets(state.JetIDs)
	require.Equal(t, expectedJets, state.JetIDs)

	// Get for anotherLight
	state = stateKeeper.Get(ctx, anotherLight, current)

	expectedIndexes = []record.Index{indexes[2]}
	sortIndexes(expectedIndexes)
	sortIndexes(state.Indexes)
	require.Equal(t, []record.Index{indexes[2]}, state.Indexes)

	expectedDrops = []drop.Drop{
		drops[2],
	}
	sortDrops(expectedDrops)
	sortDrops(state.Drops)
	require.Equal(t, expectedDrops, state.Drops)

	expectedJets = []insolar.JetID{jetIDs[2]}
	sortJets(expectedJets)
	sortJets(state.JetIDs)
	require.Equal(t, expectedJets, state.JetIDs)
}
