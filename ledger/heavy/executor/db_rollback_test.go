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
	"errors"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/stretchr/testify/require"
)

func TestDBRollback_HasOnlyGenesisPulse(t *testing.T) {
	jetKeeper := NewJetKeeperMock(t)
	jetKeeper.TopSyncPulseFunc = func() (r insolar.PulseNumber) {
		return insolar.GenesisPulse.PulseNumber
	}
	db := store.NewDBMock(t)
	drops := drop.NewDB(db)

	rollback := NewDBRollback(drops, jetKeeper)
	err := rollback.Start(context.Background())
	require.NoError(t, err)
}

func TestDBRollback_TruncateHeadError(t *testing.T) {
	jetKeeper := NewJetKeeperMock(t)
	testPulseNumber := insolar.GenesisPulse.PulseNumber + 1
	jetKeeper.TopSyncPulseFunc = func() (r insolar.PulseNumber) {
		return testPulseNumber
	}
	db := store.NewDBMock(t)
	db.NewIteratorFunc = func(p store.Key, p1 bool) (r store.Iterator) {
		iterMock := store.NewIteratorMock(t)
		iterMock.NextMock.Expect().Return(true)
		iterMock.CloseMock.Expect().Return()
		iterMock.KeyMock.Return((testPulseNumber + 1).Bytes())
		return iterMock
	}
	db.DeleteMock.Return(errors.New("Test"))

	drops := drop.NewDB(db)

	rollback := NewDBRollback(drops, jetKeeper)
	err := rollback.Start(context.Background())
	require.Error(t, err)
}

func TestDBRollback(t *testing.T) {
	jetKeeper := NewJetKeeperMock(t)
	jetKeeper.TopSyncPulseFunc = func() (r insolar.PulseNumber) {
		return insolar.GenesisPulse.PulseNumber + 1
	}
	db := store.NewDBMock(t)
	db.NewIteratorFunc = func(p store.Key, p1 bool) (r store.Iterator) {
		iterMock := store.NewIteratorMock(t)
		iterMock.NextMock.Expect().Return(false)
		iterMock.CloseMock.Expect().Return()
		return iterMock
	}

	drops := drop.NewDB(db)

	rollback := NewDBRollback(drops, jetKeeper)
	err := rollback.Start(context.Background())
	require.NoError(t, err)
}
