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

package drop

import (
	"math/rand"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/db"
)

func TestNewStorageDB(t *testing.T) {
	dbStore := NewStorageDB(db.NewMemoryMockDB())
	require.NotNil(t, dbStore)
}

type setInput struct {
	jetID insolar.JetID
	dr    Drop
}

func TestDropStorageDB_Set(t *testing.T) {
	ctx := inslogger.TestContext(t)
	var inputs []setInput
	encodedDrops := map[string]struct{}{}
	f := fuzz.New().Funcs(func(inp *setInput, c fuzz.Continue) {
		inp.dr = Drop{
			Size:  rand.Uint64(),
			Pulse: gen.PulseNumber(),
			JetID: gen.JetID(),
		}

		encoded, _ := Encode(&inp.dr)
		encodedDrops[string(encoded)] = struct{}{}
	}).NumElements(5, 5000).NilChance(0)
	f.Fuzz(&inputs)

	dbMock := db.NewDBMock(t)
	dbMock.SetFunc = func(p db.Key, p1 []byte) (r error) {
		_, ok := encodedDrops[string(p1)]
		require.Equal(t, true, ok)
		return nil
	}
	dbMock.GetMock.Return(nil, ErrNotFound)

	dropStore := NewStorageDB(dbMock)

	for _, inp := range inputs {
		err := dropStore.Set(ctx, inp.dr)
		require.NoError(t, err)
	}
}

func TestDropStorageDB_Set_ErrOverride(t *testing.T) {
	ctx := inslogger.TestContext(t)
	dr := Drop{
		Size:  rand.Uint64(),
		Pulse: gen.PulseNumber(),
		JetID: gen.JetID(),
	}

	dbMock := db.NewDBMock(t)
	dbMock.GetMock.Return(nil, nil)

	dropStore := NewStorageDB(dbMock)

	err := dropStore.Set(ctx, dr)

	require.Error(t, err, ErrNotFound)
}

func TestDropStorageDB_ForPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jetID := gen.JetID()
	pn := gen.PulseNumber()
	dr := Drop{
		Size:  rand.Uint64(),
		Pulse: gen.PulseNumber(),
	}
	buf, _ := Encode(&dr)

	dbMock := db.NewDBMock(t)
	dbMock.GetMock.Return(buf, nil)

	dropStore := NewStorageDB(dbMock)

	resDr, err := dropStore.ForPulse(ctx, jetID, pn)

	require.NoError(t, err)
	require.Equal(t, dr, resDr)
}

func TestDropStorageDB_ForPulse_NotExist(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jetID := gen.JetID()
	pn := gen.PulseNumber()

	dbMock := db.NewDBMock(t)
	dbMock.GetMock.Return(nil, ErrNotFound)

	dropStore := NewStorageDB(dbMock)

	_, err := dropStore.ForPulse(ctx, jetID, pn)

	require.Error(t, err, ErrNotFound)
}

func TestDropStorageDB_ForPulse_ProblemsWithDecoding(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jetID := gen.JetID()
	pn := gen.PulseNumber()

	dbMock := db.NewDBMock(t)
	dbMock.GetMock.Return([]byte{1, 2, 3}, nil)

	dropStore := NewStorageDB(dbMock)

	_, err := dropStore.ForPulse(ctx, jetID, pn)

	require.Error(t, err)
}
