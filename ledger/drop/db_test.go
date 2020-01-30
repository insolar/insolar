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

// +build slowtest

package drop

import (
	"context"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	fuzz "github.com/google/gofuzz"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/migration"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/tests/common"
)

var (
	poolLock     sync.Mutex
	globalPgPool *pgxpool.Pool
)

func setPool(pool *pgxpool.Pool) {
	poolLock.Lock()
	defer poolLock.Unlock()
	globalPgPool = pool
}

func getPool() *pgxpool.Pool {
	poolLock.Lock()
	defer poolLock.Unlock()
	return globalPgPool
}

// TestMain does the before and after setup
func TestMain(m *testing.M) {
	ctx := context.Background()
	log.Info("[TestMain] About to start PostgreSQL...")
	pgURL, stopPostgreSQL := common.StartPostgreSQL()
	log.Info("[TestMain] PostgreSQL started!")

	pool, err := pgxpool.Connect(ctx, pgURL)
	if err != nil {
		stopPostgreSQL()
		log.Panicf("[TestMain] pgxpool.Connect() failed: %v", err)
	}

	migrationPath := "../../migration"
	cwd, err := os.Getwd()
	if err != nil {
		stopPostgreSQL()
		panic(errors.Wrap(err, "[TestMain] os.Getwd failed"))
	}
	log.Infof("[TestMain] About to run PostgreSQL migration, cwd = %s, migration migrationPath = %s", cwd, migrationPath)
	ver, err := migration.MigrateDatabase(ctx, pool, migrationPath)
	if err != nil {
		stopPostgreSQL()
		panic(errors.Wrap(err, "Unable to migrate database"))
	}
	log.Infof("[TestMain] PostgreSQL database migration done, current schema version: %d", ver)

	setPool(pool)

	// Run all tests
	code := m.Run()

	log.Info("[TestMain] Cleaning up...")
	stopPostgreSQL()
	os.Exit(code)
}

func cleanDropsTable() {
	ctx := context.Background()
	conn, err := getPool().Acquire(ctx)
	if err != nil {
		panic("Unable to acquire a database connection")
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "DELETE FROM drops CASCADE")
	if err != nil {
		panic(err)
	}
}

func TestDropDBKey(t *testing.T) {
	t.Parallel()

	testPulseNumber := insolar.GenesisPulse.PulseNumber
	expectedKey := dropDbKey{jetPrefix: []byte("HelloWorld"), pn: testPulseNumber}

	rawID := expectedKey.ID()

	actualKey := newDropDbKey(rawID)
	require.Equal(t, expectedKey, actualKey)
}

type setInput struct {
	jetID insolar.JetID
	dr    Drop
}

func TestDropStorageDB_TruncateHead_NoSuchPulse(t *testing.T) {
	defer cleanDropsTable()

	ctx := inslogger.TestContext(t)
	db := NewDB(getPool())

	err := db.TruncateHead(ctx, insolar.GenesisPulse.PulseNumber)
	require.NoError(t, err)
}

func TestDropStorageDB_TruncateHead(t *testing.T) {
	defer cleanDropsTable()

	ctx := inslogger.TestContext(t)
	db := NewDB(getPool())

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
		err := db.Set(ctx, drop)
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			drop.JetID = *insolar.NewJetID(uint8(idx+i+50), gen.ID().Bytes())
			err = db.Set(ctx, drop)
			require.NoError(t, err)
		}
	}

	for i := 0; i < numElements; i++ {
		_, err := db.ForPulse(ctx, jets[i], startPulseNumber+insolar.PulseNumber(i))
		require.NoError(t, err)
	}

	numLeftElements := numElements / 2
	err := db.TruncateHead(ctx, startPulseNumber+insolar.PulseNumber(numLeftElements))
	require.NoError(t, err)

	for i := 0; i < numLeftElements; i++ {
		p := startPulseNumber + insolar.PulseNumber(i)
		_, err := db.ForPulse(ctx, jets[i], p)
		require.NoError(t, err, "Pulse: ", p.String())
	}

	for i := numElements - 1; i >= numLeftElements; i-- {
		p := startPulseNumber + insolar.PulseNumber(i)
		_, err := db.ForPulse(ctx, jets[i], p)
		require.EqualError(t, err, ErrNotFound.Error(), "Pulse: ", p.String())
	}
}

func TestDropStorageDB_Set(t *testing.T) {
	defer cleanDropsTable()

	ctx := inslogger.TestContext(t)
	db := NewDB(getPool())

	var inputs []setInput
	encodedDrops := map[string]struct{}{}
	f := fuzz.New().Funcs(func(inp *setInput, c fuzz.Continue) {
		inp.dr = Drop{
			Pulse: gen.PulseNumber(),
			JetID: gen.JetID(),
		}

		encoded, err := inp.dr.Marshal()
		require.NoError(t, err)
		encodedDrops[string(encoded)] = struct{}{}
	}).NumElements(5, 5000).NilChance(0)
	f.Fuzz(&inputs)

	for _, inp := range inputs {
		err := db.Set(ctx, inp.dr)
		require.NoError(t, err)
	}
}

func TestDropStorageDB_Set_ErrOverride(t *testing.T) {
	defer cleanDropsTable()

	ctx := inslogger.TestContext(t)
	db := NewDB(getPool())

	dr := Drop{
		Pulse: gen.PulseNumber(),
		JetID: gen.JetID(),
	}

	err := db.Set(ctx, dr)
	err = db.Set(ctx, dr)

	require.Error(t, err)
	require.Equal(t, ErrOverride, err)
}

func TestDropStorageDB_ForPulse(t *testing.T) {
	defer cleanDropsTable()

	ctx := inslogger.TestContext(t)
	db := NewDB(getPool())

	jetID := gen.JetID()
	pn := gen.PulseNumber()
	dr := Drop{
		Pulse: pn,
		JetID: jetID,
	}

	err := db.Set(ctx, dr)
	require.NoError(t, err)

	resDr, err := db.ForPulse(ctx, jetID, pn)

	require.NoError(t, err)
	require.Equal(t, dr, resDr)

}

func TestDropStorageDB_ForPulse_NotExist(t *testing.T) {
	defer cleanDropsTable()

	ctx := inslogger.TestContext(t)
	db := NewDB(getPool())

	jetID := gen.JetID()
	pn := gen.PulseNumber()

	_, err := db.ForPulse(ctx, jetID, pn)

	require.Error(t, err)
	require.Equal(t, ErrNotFound, err)
}
