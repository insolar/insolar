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
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/heavy/executor.InitialStateAccessor -o ./ -s _mock.go -g

// InitialStateAccessor
type InitialStateAccessor interface {
	Get(ctx context.Context, lightExecutor insolar.Reference, pulse insolar.PulseNumber) *InitialState
}

type InitialState struct {
	JetIDs  []insolar.JetID
	Drops   [][]byte
	Indexes []record.Index
}

type InitialStateKeeper struct {
	jetAccessor    jet.Accessor
	jetCoordinator jet.Coordinator
	indexStorage   object.IndexAccessor
	dropStorage    drop.Accessor

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
	lastSyncPulseNumber := jetKeeper.TopSyncPulse()

	return &InitialStateKeeper{
		jetAccessor:           jetAccessor,
		jetCoordinator:        jetCoordinator,
		indexStorage:          indexStorage,
		dropStorage:           dropStorage,
		syncPulse:             lastSyncPulseNumber,
		jetDrops:              make(map[insolar.JetID][]byte),
		abandonRequestIndexes: make(map[insolar.JetID][]record.Index),
	}
}

func (isk *InitialStateKeeper) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	// logger.Debugf("[ InitialStateKeeper ] Last finalized pulse number: %s", isk.syncPulse.String())
	// if isk.syncPulse == insolar.GenesisPulse.PulseNumber {
	// 	logger.Debug("[ InitialStateKeeper ] First start. No need to prepare state")
	// 	return nil
	// }

	logger.Debug("[ InitialStateKeeper ] prepareDrops")
	isk.prepareDrops(ctx)
	logger.Debug("[ InitialStateKeeper ] prepareAbandonRequests")
	isk.prepareAbandonRequests(ctx)
	logger.Debug("[ InitialStateKeeper ] initial state prepared")

	return nil
}

func (isk *InitialStateKeeper) prepareDrops(ctx context.Context) {
	for _, jetID := range isk.jetAccessor.All(ctx, isk.syncPulse) {
		dr, err := isk.dropStorage.ForPulse(ctx, jetID, isk.syncPulse)
		if err != nil {
			panic(fmt.Sprintf("No drop found for pulse: %s", isk.syncPulse.String()))
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
		panic(fmt.Sprintf("[ InitialStateKeeper ] Cant recieve initial state indexes: %s", err.Error()))
	}
	if len(indexes) == 0 {
		logger.Warnf("[ InitialStateKeeper ] No object indexes found in lastSyncPulseNumber: %s", isk.syncPulse.String())
	}

	// Fill the map of indexes with abandon requests
	for _, index := range indexes {

		if index.Lifeline.EarliestOpenRequest != nil {
			isk.addIndexToState(ctx, index)
		}
	}
}

func (isk *InitialStateKeeper) addIndexToState(ctx context.Context, index record.Index) {
	indexJet, _ := isk.jetAccessor.ForID(ctx, isk.syncPulse, index.ObjID)
	indexes, ok := isk.abandonRequestIndexes[indexJet]
	if !ok {
		// Someone changed jetTree in sync pulse while starting heavy material node
		// If this ever happens - we need to stop network
		panic(fmt.Sprintf("[ InitialStateKeeper ] Jet tree changed on preparing state. New jet: %s", indexJet))
	}
	isk.abandonRequestIndexes[indexJet] = append(indexes, index)
}

func (isk *InitialStateKeeper) Get(ctx context.Context, lightExecutor insolar.Reference, pulse insolar.PulseNumber) *InitialState {
	logger := inslogger.FromContext(ctx)

	var jetIDs []insolar.JetID
	var drops [][]byte
	var indexes []record.Index

	logger.Debugf("[ InitialStateKeeper ] Getting drops for: %s in pulse: %s", lightExecutor.String(), pulse.String())
	for id, jetDrop := range isk.jetDrops {
		light, err := isk.jetCoordinator.LightExecutorForJet(ctx, insolar.ID(id), pulse)
		if err != nil {
			panic(fmt.Sprintf("No drop found for pulse: %s", isk.syncPulse.String()))
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
