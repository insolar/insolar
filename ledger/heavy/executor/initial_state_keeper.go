// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/heavy/executor.InitialStateAccessor -o ./ -s _mock.go -g

// InitialStateAccessor interface can receive initial state for lights
type InitialStateAccessor interface {
	// Get method returns InitialState filled only with data for passed light node
	// If node isn't lightExecutor for any jets - arrays in InitialState will be empty
	// Passed pulse is current pulse for checking light executor for jets (not a topSyncPulse)
	Get(ctx context.Context, lightExecutor insolar.Reference, pulse insolar.PulseNumber) *InitialState
}

type InitialState struct {
	// JetIds for passed executor (not all ids). If JetDrop for this jet has Split flag - both jets will be in slice
	JetIDs []insolar.JetID
	// Drops for JetIDs above
	Drops []drop.Drop
	// Indexes only for Lifelines that has pending requests
	Indexes []record.Index
}

// InitialStateKeeper prepares state for LMEs
type InitialStateKeeper struct {
	jetAccessor    jet.Accessor
	jetCoordinator jet.Coordinator
	indexStorage   object.MemoryIndexAccessor
	dropStorage    drop.Accessor

	syncPulse insolar.PulseNumber

	lock                  sync.RWMutex
	jetSiblings           map[insolar.JetID]insolar.JetID
	jetDrops              map[insolar.JetID]drop.Drop
	abandonRequestIndexes map[insolar.JetID][]record.Index
}

func NewInitialStateKeeper(
	jetKeeper JetKeeper,
	jetAccessor jet.Accessor,
	jetCoordinator jet.Coordinator,
	indexStorage object.MemoryIndexAccessor,
	dropStorage drop.Accessor,
) *InitialStateKeeper {
	return &InitialStateKeeper{
		jetAccessor:           jetAccessor,
		jetCoordinator:        jetCoordinator,
		indexStorage:          indexStorage,
		dropStorage:           dropStorage,
		syncPulse:             jetKeeper.TopSyncPulse(),
		jetSiblings:           make(map[insolar.JetID]insolar.JetID),
		jetDrops:              make(map[insolar.JetID]drop.Drop),
		abandonRequestIndexes: make(map[insolar.JetID][]record.Index),
	}
}

// Start method prepares state before network starts
func (isk *InitialStateKeeper) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	isk.lock.Lock()
	defer isk.lock.Unlock()

	logger.Info("[ InitialStateKeeper ] Prepare drops for JetIds")
	isk.prepareDrops(ctx)
	logger.Info("[ InitialStateKeeper ] Prepare abandon request indexes")
	isk.prepareAbandonRequests(ctx)
	logger.Info("[ InitialStateKeeper ] Initial state prepared")

	return nil
}

func (isk *InitialStateKeeper) prepareDrops(ctx context.Context) {
	for _, jetID := range isk.jetAccessor.All(ctx, isk.syncPulse) {
		dr, err := isk.dropStorage.ForPulse(ctx, jetID, isk.syncPulse)
		if err != nil {
			inslogger.FromContext(ctx).Fatal("No drop found for pulse: ", isk.syncPulse.String())
		}

		if dr.Split {
			left, right := jet.Siblings(jetID)

			isk.jetSiblings[left] = right
			isk.jetSiblings[right] = left

			isk.jetDrops[left] = dr
			isk.jetDrops[right] = dr
		} else {
			isk.jetDrops[jetID] = dr
		}
	}
}

func (isk *InitialStateKeeper) prepareAbandonRequests(ctx context.Context) {
	logger := inslogger.FromContext(ctx)

	tree := jet.NewTree(true)
	for jetID := range isk.jetDrops {
		tree.Update(jetID, true)
		isk.abandonRequestIndexes[jetID] = []record.Index{}
	}

	indexes, err := isk.indexStorage.ForPulse(ctx, isk.syncPulse)
	if err != nil {
		if err == object.ErrIndexNotFound {
			logger.Warnf("[ InitialStateKeeper ] No object indexes found in lastSyncPulseNumber: %s", isk.syncPulse.String())
			return
		}
		logger.Fatal("Cant receive initial state indexes: ", err.Error())
	}

	// Fill the map of indexes with abandon requests
	for _, index := range indexes {

		if index.Lifeline.EarliestOpenRequest != nil {
			isk.addIndexToState(ctx, index, tree)
		}
	}
}

func (isk *InitialStateKeeper) addIndexToState(ctx context.Context, index record.Index, tree *jet.Tree) {
	logger := inslogger.FromContext(ctx)
	indexJet, _ := tree.Find(index.ObjID)
	indexes, ok := isk.abandonRequestIndexes[indexJet]
	if !ok {
		// Someone changed jetTree in sync pulse while starting heavy material node
		// If this ever happens - we need to stop network
		logger.Fatal("Jet tree changed on preparing state. New jet: ", indexJet)
	}
	logger.Info("Prepare index with abandon request: %s in jet %s", index.ObjID.DebugString(), indexJet.DebugString())
	isk.abandonRequestIndexes[indexJet] = append(indexes, index)
}

func (isk *InitialStateKeeper) Get(ctx context.Context, lightExecutor insolar.Reference, pulse insolar.PulseNumber) *InitialState {
	logger := inslogger.FromContext(ctx)

	isk.lock.RLock()
	defer isk.lock.RUnlock()

	jetIDs := make([]insolar.JetID, 0)
	drops := make([]drop.Drop, 0)
	indexes := make([]record.Index, 0)

	logger.Infof("[ InitialStateKeeper ] Getting drops for: %s in pulse: %s", lightExecutor.String(), pulse.String())

	// Must not send two equal drops to single LME after split
	existingDrops := make(map[insolar.JetID]struct{})

	for id, jetDrop := range isk.jetDrops {
		light, err := isk.jetCoordinator.LightExecutorForJet(ctx, insolar.ID(id), pulse)
		if err != nil {
			logger.Fatal("No drop found for pulse ", isk.syncPulse.String())
		}

		if light.Equal(lightExecutor) {
			jetIDs = append(jetIDs, id)

			if _, ok := existingDrops[id]; ok {
				continue
			}

			drops = append(drops, jetDrop)
			if siblingID, ok := isk.jetSiblings[id]; ok {
				existingDrops[siblingID] = struct{}{}
			}
		}
	}

	logger.Infof("[ InitialStateKeeper ] Getting indexes for: %s in pulse: %s", lightExecutor.String(), pulse.String())
	for _, id := range jetIDs {
		indexes = append(indexes, isk.abandonRequestIndexes[id]...)
	}

	return &InitialState{
		JetIDs:  jetIDs,
		Drops:   drops,
		Indexes: indexes,
	}
}
