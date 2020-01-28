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

package executor

import (
	"context"
	"sync"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
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
//func TestMain(m *testing.M) {
//	ctx := context.Background()
//	log.Info("[TestMain] About to start PostgreSQL...")
//	pgURL, stopPostgreSQL := common.StartPostgreSQL()
//	log.Info("[TestMain] PostgreSQL started!")
//
//	pool, err := pgxpool.Connect(ctx, pgURL)
//	if err != nil {
//		stopPostgreSQL()
//		log.Panicf("[TestMain] pgxpool.Connect() failed: %v", err)
//	}
//
//	migrationPath := "../../../migration"
//	cwd, err := os.Getwd()
//	if err != nil {
//		stopPostgreSQL()
//		panic(errors.Wrap(err, "[TestMain] os.Getwd failed"))
//	}
//	log.Infof("[TestMain] About to run PostgreSQL migration, cwd = %s, migration migrationPath = %s", cwd, migrationPath)
//	ver, err := migration.MigrateDatabase(ctx, pool, migrationPath)
//	if err != nil {
//		stopPostgreSQL()
//		panic(errors.Wrap(err, "Unable to migrate database"))
//	}
//	log.Infof("[TestMain] PostgreSQL database migration done, current schema version: %d", ver)
//
//	setPool(pool)
//	// Run all tests
//	code := m.Run()
//
//	log.Info("[TestMain] Cleaning up...")
//	stopPostgreSQL()
//	os.Exit(code)
//}

func initDB(t *testing.T, testPulse insolar.PulseNumber) (*DBJetKeeper, *jet.DBStore, *pulse.DB) {
	ctx := context.Background()
	jets := jet.NewDBStore(getPool())
	pulses := pulse.NewDB(getPool())
	err := pulses.Append(ctx, insolar.Pulse{PulseNumber: insolar.GenesisPulse.PulseNumber})
	require.NoError(t, err)

	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: testPulse})
	require.NoError(t, err)

	jetKeeper := NewJetKeeper(jets, getPool(), pulses)

	return jetKeeper, jets, pulses
}

func Test_TruncateHead(t *testing.T) {
	ctx := inslogger.TestContext(t)
	testPulse := insolar.GenesisPulse.PulseNumber + 10
	ji, jets, _ := initDB(t, testPulse)

	testJet := insolar.ZeroJetID

	err := jets.Update(ctx, testPulse, true, testJet)
	require.NoError(t, err)
	err = ji.AddHotConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)
	err = ji.AddDropConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)
	err = ji.AddBackupConfirmation(ctx, testPulse)
	require.NoError(t, err)

	require.Equal(t, testPulse, ji.TopSyncPulse())

	_, err = ji.get(testPulse)
	require.NoError(t, err)

	nextPulse := testPulse + 10

	err = ji.AddDropConfirmation(ctx, nextPulse, gen.JetID(), false)
	require.NoError(t, err)
	err = ji.AddHotConfirmation(ctx, nextPulse, gen.JetID(), false)
	require.NoError(t, err)

	_, err = ji.get(nextPulse)
	require.NoError(t, err)

	err = ji.TruncateHead(ctx, nextPulse)
	require.NoError(t, err)

	_, err = ji.get(testPulse)
	require.NoError(t, err)
	_, err = ji.get(nextPulse)
	require.EqualError(t, err, "value not found")
}
