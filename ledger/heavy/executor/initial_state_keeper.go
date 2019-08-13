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
	"github.com/insolar/insolar/ledger/object"
)

type InitialStateReader interface {
	ReadForJetID(context.Context, insolar.JetID) ([]record.Index, error)
}

type InitialStateKeeper struct {
	jetAccessor  jet.Accessor
	indexStorage object.IndexAccessor

	syncPulse       insolar.PulseNumber
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

	logger.Debug("[ InitialStateKeeper ] Preparing initial state")
	indexes, err := isk.indexStorage.ForPulse(ctx, isk.syncPulse)

	if err != nil {
		panic(fmt.Sprintf("[ InitialStateKeeper ] Cant recieve initial state indexes: %s", err.Error()))
	}

	if len(indexes) == 0 {
		logger.Warnf("[ InitialStateKeeper ] No object indexes found in lastSyncPulseNumber: %s", isk.syncPulse.String())
	}

	// Build empty map for with all known JetIDs
	knownJets := isk.jetAccessor.All(ctx, isk.syncPulse)
	for _, jetID := range knownJets {
		isk.abandonRequests[jetID] = []record.Index{}
	}

	// Fill abandonRequests with pending indexes
	for i := 0; i < len(indexes); i++ {
		index := indexes[i]

		// Check if lifeline has open pending requests
		if index.Lifeline.EarliestOpenRequest != nil {
			isk.addIndexToState(ctx, index)
		}
	}

	logger.Debug("[ InitialStateKeeper ] Initial state prepared")
	return nil
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

func (isk *InitialStateKeeper) ReadForJetID(ctx context.Context, jetID insolar.JetID) ([]record.Index, error) {
	indexes, ok := isk.abandonRequests[jetID]
	if !ok {
		// Someone changed jetTree in sync pulse while starting network
		// If this ever happens - we need to stop network
		panic(fmt.Sprintf("[ InitialStateKeeper ] Jet is not known: %s", jetID))
	}
	return indexes, nil
}
