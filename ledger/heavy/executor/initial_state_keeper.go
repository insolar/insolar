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
	Drops [][]byte
	// Indexes only for Lifelines that has pending requests
	Indexes []record.Index
}

// InitialStateKeeper prepares state for LMEs
type InitialStateKeeper struct {
	jetAccessor    jet.Accessor
	jetCoordinator jet.Coordinator
	indexStorage   object.IndexAccessor
	dropStorage    drop.Accessor

	lock                  sync.RWMutex
	syncPulse             insolar.PulseNumber
	jetDrops              map[insolar.JetID][]byte
	abandonRequestIndexes map[insolar.JetID][]record.Index
}

func NewInitialStateKeeper(
	jetKeeper JetKeeper,
	jetAccessor jet.Accessor,
	jetCoordinator jet.Coordinator,
	indexStorage object.IndexAccessor,
	dropStorage drop.Accessor,
) *InitialStateKeeper {
	return &InitialStateKeeper{
		jetAccessor:           jetAccessor,
		jetCoordinator:        jetCoordinator,
		indexStorage:          indexStorage,
		dropStorage:           dropStorage,
		syncPulse:             jetKeeper.TopSyncPulse(),
		jetDrops:              make(map[insolar.JetID][]byte),
		abandonRequestIndexes: make(map[insolar.JetID][]record.Index),
	}
}

// Start method prepares state before network starts
func (isk *InitialStateKeeper) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	isk.lock.Lock()
	defer isk.lock.Unlock()

	logger.Debug("[ InitialStateKeeper ] Prepare drops for JetIds")
	isk.prepareDrops(ctx)
	logger.Debug("[ InitialStateKeeper ] Prepare abandon request indexes")
	isk.prepareAbandonRequests(ctx)
	logger.Debug("[ InitialStateKeeper ] Initial state prepared")

	return nil
}

func (isk *InitialStateKeeper) prepareDrops(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	for _, jetID := range isk.jetAccessor.All(ctx, isk.syncPulse) {
		dr, err := isk.dropStorage.ForPulse(ctx, jetID, isk.syncPulse)
		if err != nil {
			logger.Fatal("No drop found for pulse: ", isk.syncPulse.String())
		}

		jetDrop := drop.MustEncode(&dr)

		if dr.Split {
			left, right := jet.Siblings(jetID)
			isk.jetDrops[left] = jetDrop
			isk.jetDrops[right] = jetDrop
		} else {
			isk.jetDrops[jetID] = jetDrop
		}
	}
}

func (isk *InitialStateKeeper) prepareAbandonRequests(ctx context.Context) {
	logger := inslogger.FromContext(ctx)

	for jetID := range isk.jetDrops {
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
			isk.addIndexToState(ctx, index)
		}
	}
}

func (isk *InitialStateKeeper) addIndexToState(ctx context.Context, index record.Index) {
	logger := inslogger.FromContext(ctx)
	indexJet, _ := isk.jetAccessor.ForID(ctx, isk.syncPulse, index.ObjID)
	indexes, ok := isk.abandonRequestIndexes[indexJet]
	if !ok {
		// Someone changed jetTree in sync pulse while starting heavy material node
		// If this ever happens - we need to stop network
		logger.Fatal("Jet tree changed on preparing state. New jet: ", indexJet)
	}
	isk.abandonRequestIndexes[indexJet] = append(indexes, index)
}

func (isk *InitialStateKeeper) Get(ctx context.Context, lightExecutor insolar.Reference, pulse insolar.PulseNumber) *InitialState {
	logger := inslogger.FromContext(ctx)

	isk.lock.RLock()
	defer isk.lock.RUnlock()

	jetIDs := make([]insolar.JetID, 0)
	drops := make([][]byte, 0)
	indexes := make([]record.Index, 0)

	logger.Debugf("[ InitialStateKeeper ] Getting drops for: %s in pulse: %s", lightExecutor.String(), pulse.String())
	for id, jetDrop := range isk.jetDrops {
		light, err := isk.jetCoordinator.LightExecutorForJet(ctx, insolar.ID(id), pulse)
		if err != nil {
			logger.Fatal("No drop found for pulse ", isk.syncPulse.String())
		}

		if light.Equal(lightExecutor) {
			jetIDs = append(jetIDs, id)
			drops = append(drops, jetDrop)
		}
	}

	logger.Debugf("[ InitialStateKeeper ] Getting indexes for: %s in pulse: %s", lightExecutor.String(), pulse.String())
	for _, id := range jetIDs {
		indexes = append(indexes, isk.abandonRequestIndexes[id]...)
	}

	return &InitialState{
		JetIDs:  jetIDs,
		Drops:   drops,
		Indexes: indexes,
	}
}
