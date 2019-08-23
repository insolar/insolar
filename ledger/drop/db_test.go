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
	"context"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func TestDropDBKey(t *testing.T) {
	t.Parallel()

	testPulseNumber := insolar.GenesisPulse.PulseNumber
	expectedKey := dropDbKey{jetPrefix: []byte("HelloWorld"), pn: testPulseNumber}

	rawID := expectedKey.ID()

	actualKey := newDropDbKey(rawID)
	require.Equal(t, expectedKey, actualKey)
}

func TestNewStorageDB(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	dbStore := NewDB(db)
	require.NotNil(t, dbStore)
}

type setInput struct {
	jetID insolar.JetID
	dr    Drop
}

func TestDropStorageDB_TruncateHead_NoSuchPulse(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	dbMock, err := store.NewBadgerDB(ops)
	defer dbMock.Stop(ctx)
	require.NoError(t, err)

	dropStore := NewDB(dbMock)

	err = dropStore.TruncateHead(ctx, insolar.GenesisPulse.PulseNumber)
	require.NoError(t, err)
}

func TestDropStorageDB_TruncateHead(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	dbMock, err := store.NewBadgerDB(ops)
	defer dbMock.Stop(ctx)
	require.NoError(t, err)

	dropStore := NewDB(dbMock)

	numElements := 10

	// it's used for writing pulses in random order to db
	indexes := make([]int, numElements)
	for i := 0; i < numElements; i++ {
		indexes[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(indexes), func(i, j int) { indexes[i], indexes[j] = indexes[j], indexes[i] })

	startPulseNumber := insolar.GenesisPulse.PulseNumber
	jets := make([]insolar.JetID, numElements)
	for _, idx := range indexes {
		drop := Drop{}
		drop.Pulse = startPulseNumber + insolar.PulseNumber(idx)
		jets[idx] = *insolar.NewJetID(uint8(idx), gen.ID().Bytes())

		drop.JetID = jets[idx]
		err := dropStore.Set(ctx, drop)
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			drop.JetID = *insolar.NewJetID(uint8(idx+i+50), gen.ID().Bytes())
			err = dropStore.Set(ctx, drop)
			require.NoError(t, err)
		}
	}

	for i := 0; i < numElements; i++ {
		_, err := dropStore.ForPulse(ctx, jets[i], startPulseNumber+insolar.PulseNumber(i))
		require.NoError(t, err)
	}

	numLeftElements := numElements / 2
	err = dropStore.TruncateHead(ctx, startPulseNumber+insolar.PulseNumber(numLeftElements))
	require.NoError(t, err)

	for i := 0; i < numLeftElements; i++ {
		p := startPulseNumber + insolar.PulseNumber(i)
		_, err := dropStore.ForPulse(ctx, jets[i], p)
		require.NoError(t, err, "Pulse: ", p.String())
	}

	for i := numElements - 1; i >= numLeftElements; i-- {
		p := startPulseNumber + insolar.PulseNumber(i)
		_, err := dropStore.ForPulse(ctx, jets[i], p)
		require.EqualError(t, err, ErrNotFound.Error(), "Pulse: ", p.String())
	}
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

		encoded := MustEncode(&inp.dr)
		encodedDrops[string(encoded)] = struct{}{}
	}).NumElements(5, 5000).NilChance(0)
	f.Fuzz(&inputs)

	dbMock := store.NewDBMock(t)
	dbMock.SetMock.Set(func(p store.Key, p1 []byte) (r error) {
		_, ok := encodedDrops[string(p1)]
		require.Equal(t, true, ok)
		return nil
	})
	dbMock.GetMock.Return(nil, ErrNotFound)

	dropStore := NewDB(dbMock)

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

	dbMock := store.NewDBMock(t)
	dbMock.GetMock.Return(nil, nil)

	dropStore := NewDB(dbMock)

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
	buf := MustEncode(&dr)

	dbMock := store.NewDBMock(t)
	dbMock.GetMock.Return(buf, nil)

	dropStore := NewDB(dbMock)

	resDr, err := dropStore.ForPulse(ctx, jetID, pn)

	require.NoError(t, err)
	require.Equal(t, dr, resDr)
}

func TestDropStorageDB_ForPulse_NotExist(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jetID := gen.JetID()
	pn := gen.PulseNumber()

	dbMock := store.NewDBMock(t)
	dbMock.GetMock.Return(nil, ErrNotFound)

	dropStore := NewDB(dbMock)

	_, err := dropStore.ForPulse(ctx, jetID, pn)

	require.Error(t, err, ErrNotFound)
}

func TestDropStorageDB_ForPulse_ProblemsWithDecoding(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jetID := gen.JetID()
	pn := gen.PulseNumber()

	dbMock := store.NewDBMock(t)
	dbMock.GetMock.Return([]byte{1, 2, 3}, nil)

	dropStore := NewDB(dbMock)

	_, err := dropStore.ForPulse(ctx, jetID, pn)

	require.Error(t, err)
}
