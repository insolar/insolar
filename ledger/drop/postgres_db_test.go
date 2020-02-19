// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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

func cleanupDatabase() {
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

func TestPostgresDropDBKey(t *testing.T) {
	t.Parallel()

	testPulseNumber := insolar.GenesisPulse.PulseNumber
	expectedKey := dropDbKey{jetPrefix: []byte("HelloWorld"), pn: testPulseNumber}

	rawID := expectedKey.ID()

	actualKey := newDropDbKey(rawID)
	require.Equal(t, expectedKey, actualKey)
}

func TestPostgresDropStorageDB_TruncateHead_NoSuchPulse(t *testing.T) {
	defer cleanupDatabase()

	ctx := inslogger.TestContext(t)
	db := NewPostgresDB(getPool())

	err := db.TruncateHead(ctx, insolar.GenesisPulse.PulseNumber)
	require.NoError(t, err)
}

func TestPostgresDropStorageDB_TruncateHead(t *testing.T) {
	defer cleanupDatabase()

	ctx := inslogger.TestContext(t)
	db := NewPostgresDB(getPool())

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

func TestPostgresDropStorageDB_Set(t *testing.T) {
	defer cleanupDatabase()

	ctx := inslogger.TestContext(t)
	db := NewPostgresDB(getPool())

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

func TestPostgresDropStorageDB_Set_ErrOverride(t *testing.T) {
	defer cleanupDatabase()

	ctx := inslogger.TestContext(t)
	db := NewPostgresDB(getPool())

	dr := Drop{
		Pulse: gen.PulseNumber(),
		JetID: gen.JetID(),
	}

	err := db.Set(ctx, dr)
	err = db.Set(ctx, dr)

	require.Error(t, err)
	require.Equal(t, ErrOverride, err)
}

func TestPostgresDropStorageDB_ForPulse(t *testing.T) {
	defer cleanupDatabase()

	ctx := inslogger.TestContext(t)
	db := NewPostgresDB(getPool())

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

func TestPostgresDropStorageDB_ForPulse_NotExist(t *testing.T) {
	defer cleanupDatabase()

	ctx := inslogger.TestContext(t)
	db := NewPostgresDB(getPool())

	jetID := gen.JetID()
	pn := gen.PulseNumber()

	_, err := db.ForPulse(ctx, jetID, pn)

	require.Error(t, err)
	require.Equal(t, ErrNotFound, err)
}
