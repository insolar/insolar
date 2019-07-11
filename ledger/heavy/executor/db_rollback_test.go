///
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
///

package executor

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/stretchr/testify/require"
)

func TestDBRollback_HasOnlyGenesisPulse(t *testing.T) {
	jetKeeper := NewJetKeeperMock(t)
	jetKeeper.TopSyncPulseFunc = func() (r insolar.PulseNumber) {
		return insolar.GenesisPulse.PulseNumber
	}

	rollback := NewDBRollback(jetKeeper, nil, nil, nil, nil, nil)
	err := rollback.Start(context.Background())
	require.NoError(t, err)
}

func TestDBRollback_HappyPath(t *testing.T) {
	jetKeeper := NewJetKeeperMock(t)
	jetKeeper.TopSyncPulseFunc = func() (r insolar.PulseNumber) {
		return insolar.GenesisPulse.PulseNumber + 1
	}
	db := store.NewDBMock(t)
	hits := make(map[store.Scope]int)

	db.DeleteMock.Return(nil)
	iterNum := 0
	db.NewIteratorFunc = func(p store.Key, p1 bool) (r store.Iterator) {
		num, _ := hits[p.Scope()]
		hits[p.Scope()] = num + 1

		iterMock := store.NewIteratorMock(t)
		iterMock.NextFunc = func() (r bool) {
			iterNum++
			return iterNum%2 != 0
		}

		iterMock.KeyMock.Return(p.ID())
		iterMock.CloseMock.Return()
		return iterMock
	}

	drops := drop.NewDB(db)
	records := object.NewRecordDB(db)
	indexes := object.NewIndexDB(db)
	jets := jet.NewDBStore(db)
	pulses := pulse.NewDB(db)

	rollback := NewDBRollback(jetKeeper, drops, records, indexes, jets, pulses)
	err := rollback.Start(context.Background())
	require.Len(t, hits, 5) // drops, record, jets, indexes, pulses
	expectedScopes := []struct {
		scope   store.Scope
		numHits int
	}{
		{store.ScopeJetDrop, 1},
		{store.ScopeRecord, 1},
		{store.ScopeIndex, 1},
		{store.ScopeJetTree, 1},
		{store.ScopePulse, 1}}
	for _, s := range expectedScopes {
		actualNum, ok := hits[s.scope]
		require.True(t, ok, "Scope: ", s.scope)
		require.Equal(t, s.numHits, actualNum, "Scope: ", s.scope)
	}

	require.NoError(t, err)
}
