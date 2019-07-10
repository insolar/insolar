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

	rollback := NewDBRollback(nil, nil, nil, nil, jetKeeper)
	err := rollback.Start(context.Background())
	require.NoError(t, err)
}

func TestDBRollback_HappyPath(t *testing.T) {
	jetKeeper := NewJetKeeperMock(t)
	jetKeeper.TopSyncPulseFunc = func() (r insolar.PulseNumber) {
		return insolar.GenesisPulse.PulseNumber + 1
	}
	db := store.NewDBMock(t)
	hits := make(map[store.Scope]struct{})

	db.NewIteratorFunc = func(p store.Key, p1 bool) (r store.Iterator) {
		// check that every 'scope' called once
		_, exists := hits[p.Scope()]
		require.False(t, exists)
		hits[p.Scope()] = struct{}{}

		iterMock := store.NewIteratorMock(t)
		iterMock.NextMock.Expect().Return(false)
		iterMock.CloseMock.Expect().Return()
		return iterMock
	}

	drops := drop.NewDB(db)
	records := object.NewRecordDB(db)
	indexes := object.NewIndexDB(db)
	jets := jet.NewDBStore(db)

	rollback := NewDBRollback(drops, records, indexes, jets, jetKeeper)
	err := rollback.Start(context.Background())
	require.Len(t, hits, 4) // drops, record, jets, indexes
	expectedScopes := []store.Scope{store.ScopeJetDrop, store.ScopeRecord, store.ScopeIndex, store.ScopeJetTree}
	for _, s := range expectedScopes {
		_, ok := hits[s]
		require.True(t, ok)
	}

	require.NoError(t, err)
}
