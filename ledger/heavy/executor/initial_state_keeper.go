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

type InitialStateAccessor interface {
	ReadForJetID(context.Context, insolar.JetID) []record.Index
}

type InitialState struct {
	JetIDs  []insolar.JetID
	Drops   [][]byte
	Indexes []record.Index
}

type InitialStateKeeper struct {
	jetAccessor  jet.Accessor
	indexStorage object.IndexAccessor
	dropStorage  drop.Accessor

	syncPulse       insolar.PulseNumber
	jetDrops        map[insolar.JetID][]byte
	abandonRequests map[insolar.JetID][]record.Index
}

func NewInitialStateKeeper(jetKeeper JetKeeper, jetAccessor jet.Accessor, indexStorage object.IndexAccessor) *InitialStateKeeper {
	lastSyncPulseNumber := jetKeeper.TopSyncPulse()

	return &InitialStateKeeper{
		jetAccessor:     jetAccessor,
		indexStorage:    indexStorage,
		syncPulse:       lastSyncPulseNumber,
		abandonRequests: make(map[insolar.JetID][]record.Index),
	}
}

func (isk *InitialStateKeeper) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	logger.Debugf("[ InitialStateKeeper ] Last finalized pulse number: %s", isk.syncPulse.String())
	if isk.syncPulse == insolar.GenesisPulse.PulseNumber {
		logger.Debug("[ InitialStateKeeper ] First start. No need to prepare state")
		return nil
	}

	logger.Debug("[ InitialStateKeeper ] prepareDrops")
	isk.prepareDrops(ctx)
	logger.Debug("[ InitialStateKeeper ] prepareAbandonRequests")
	isk.prepareAbandonRequests(ctx)
	logger.Debug("[ InitialStateKeeper ] Initial state prepared")

	return nil
}

func (isk *InitialStateKeeper) prepareDrops(ctx context.Context) {
	for _, jetID := range isk.jetAccessor.All(ctx, isk.syncPulse) {
		dr, err := isk.dropStorage.ForPulse(ctx, jetID, isk.syncPulse)
		if err != nil {
			// TODO: panic
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
		isk.abandonRequests[jetID] = []record.Index{}
	}

	indexes, err := isk.indexStorage.ForPulse(ctx, isk.syncPulse)
	if err != nil {
		panic(fmt.Sprintf("[ InitialStateKeeper ] Cant recieve initial state indexes: %s", err.Error()))
	}
	if len(indexes) == 0 {
		logger.Warnf("[ InitialStateKeeper ] No object indexes found in lastSyncPulseNumber: %s", isk.syncPulse.String())
	}

	// Fill map
	for _, index := range indexes {

		if index.Lifeline.EarliestOpenRequest != nil {
			isk.addIndexToState(ctx, index)
		}
	}
}

func (isk *InitialStateKeeper) addIndexToState(ctx context.Context, index record.Index) {
	indexJet, _ := isk.jetAccessor.ForID(ctx, isk.syncPulse, index.ObjID)
	indexes, ok := isk.abandonRequests[indexJet]
	if !ok {
		// Someone changed jetTree in sync pulse while starting heavy material node
		// If this ever happens - we need to stop network
		panic(fmt.Sprintf("[ InitialStateKeeper ] Jet tree changed on preparing state. New jet: %s", indexJet))
	}
	isk.abandonRequests[indexJet] = append(indexes, index)
}

func (isk *InitialStateKeeper) Get(ctx context.Context, lightExecutor insolar.Reference) *InitialState {

	return &InitialState{}
}

func (isk *InitialStateKeeper) ReadForJetID(ctx context.Context, jetID insolar.JetID) []record.Index {
	indexes, ok := isk.abandonRequests[jetID]
	if !ok {
		// Someone changed jetTree in sync pulse while starting network
		// If this ever happens - we need to stop network
		panic(fmt.Sprintf("[ InitialStateKeeper ] Jet is not known: %s", jetID))
	}
	return indexes
}
