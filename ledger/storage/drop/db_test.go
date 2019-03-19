/*
*    Copyright 2019 Insolar Technologies
*
*    Licensed under the Apache License, Version 2.0 (the "License");
*    you may not use this file except in compliance with the License.
*    You may obtain a copy of the License at
*
*        http://www.apache.org/licenses/LICENSE-2.0
*
*    Unless required by applicable law or agreed to in writing, software
*    distributed under the License is distributed on an "AS IS" BASIS,
*    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*    See the License for the specific language governing permissions and
*    limitations under the License.
 */

package drop

import (
	"math/rand"
	"testing"

	"github.com/google/gofuzz"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/db"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/stretchr/testify/require"
)

func TestNewStorageDB(t *testing.T) {
	dbStor := NewStorageDB()

	require.NotNil(t, dbStor)
}

type setInput struct {
	jetID core.JetID
	dr    jet.Drop
}

func TestDropStorageDB_Set(t *testing.T) {
	ctx := inslogger.TestContext(t)
	var inputs []setInput
	encodedDrops := map[string]struct{}{}
	f := fuzz.New().Funcs(func(inp *setInput, c fuzz.Continue) {
		inp.dr = jet.Drop{
			Size:  rand.Uint64(),
			Pulse: gen.PulseNumber(),
			JetID: gen.JetID(),
		}

		encoded, _ := jet.Encode(&inp.dr)
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

	dropStor := NewStorageDB()
	dropStor.DB = dbMock

	for _, inp := range inputs {
		err := dropStor.Set(ctx, inp.dr)
		require.NoError(t, err)
	}
}

func TestDropStorageDB_Set_ErrOverride(t *testing.T) {
	ctx := inslogger.TestContext(t)
	dr := jet.Drop{
		Size:  rand.Uint64(),
		Pulse: gen.PulseNumber(),
		JetID: gen.JetID(),
	}

	dbMock := db.NewDBMock(t)
	dbMock.GetMock.Return(nil, nil)

	dropStor := NewStorageDB()
	dropStor.DB = dbMock

	err := dropStor.Set(ctx, dr)

	require.Error(t, err, ErrNotFound)
}

func TestDropStorageDB_ForPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jetID := gen.JetID()
	pn := gen.PulseNumber()
	dr := jet.Drop{
		Size:  rand.Uint64(),
		Pulse: gen.PulseNumber(),
	}
	buf, _ := jet.Encode(&dr)

	dbMock := db.NewDBMock(t)
	dbMock.GetMock.Return(buf, nil)

	dropStor := NewStorageDB()
	dropStor.DB = dbMock

	resDr, err := dropStor.ForPulse(ctx, jetID, pn)

	require.NoError(t, err)
	require.Equal(t, dr, resDr)
}

func TestDropStorageDB_ForPulse_NotExist(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jetID := gen.JetID()
	pn := gen.PulseNumber()

	dbMock := db.NewDBMock(t)
	dbMock.GetMock.Return(nil, ErrNotFound)

	dropStor := NewStorageDB()
	dropStor.DB = dbMock

	_, err := dropStor.ForPulse(ctx, jetID, pn)

	require.Error(t, err, ErrNotFound)
}

func TestDropStorageDB_ForPulse_ProblemsWithDecoding(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jetID := gen.JetID()
	pn := gen.PulseNumber()

	dbMock := db.NewDBMock(t)
	dbMock.GetMock.Return([]byte{1, 2, 3}, nil)

	dropStor := NewStorageDB()
	dropStor.DB = dbMock

	_, err := dropStor.ForPulse(ctx, jetID, pn)

	require.Error(t, err)
}
