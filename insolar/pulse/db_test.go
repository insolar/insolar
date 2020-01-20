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

package pulse

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/ledger/heavy/migration"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/tests/common"
	"github.com/stretchr/testify/require"
)

var db *DB

// TestMain does the before and after setup
func TestMain(m *testing.M) {
	ctx := context.Background()
	log.Info("[TestMain] About to start PostgreSQL...")
	pgURL, stopPostgreSQL := common.StartPostgreSQL()
	log.Info("[TestMain] PostgreSQL started!")
	defer stopPostgreSQL()

	pool, err := pgxpool.Connect(ctx, pgURL)
	if err != nil {
		log.Panicf("[TestMain] pgxpool.Connect() failed: %v", err)
	}

	migrationPath := "../../migration"
	cwd, err := os.Getwd()
	if err != nil {
		panic(errors.Wrap(err, "[TestMain] os.Getwd failed"))
	}
	log.Infof("[TestMain] About to run PostgreSQL migration, cwd = %s, migration migrationPath = %s", cwd, migrationPath)
	ver, err := migration.MigrateDatabase(ctx, pool, migrationPath)
	if err != nil {
		panic(errors.Wrap(err, "Unable to migrate database"))
	}
	log.Infof("[TestMain] PostgreSQL database migration done, current schema version: %d", ver)

	db = NewDB(pool)

	// Run all tests
	code := m.Run()

	log.Info("[TestMain] Cleaning up...")
	os.Exit(code)
}

func TestAppend(t *testing.T) {
	ctx := context.Background()
	pn := gen.PulseNumber()

	// Make sure there is no such pulse in DB yet
	_, err := db.ForPulseNumber(ctx, pn)
	require.Error(t, err)

	conf := insolar.PulseSenderConfirmation{
		PulseNumber:     pn,
		ChosenPublicKey: "lol",
		Entropy:         [insolar.EntropySize]byte{3, 3, 2, 2, 1, 1},
		Signature:       []byte{1, 1, 2, 2, 3, 3},
	}
	signs := make(map[string]insolar.PulseSenderConfirmation, 1)
	signs[conf.ChosenPublicKey] = conf
	writePulse := insolar.Pulse{
		PulseNumber:      pn,
		PrevPulseNumber:  gen.PulseNumber(),
		NextPulseNumber:  gen.PulseNumber(),
		PulseTimestamp:   123456789,
		EpochPulseNumber: pulse.Epoch(1234),
		OriginID:         [insolar.OriginIDSize]byte{3, 2, 1},
		Entropy:          [insolar.EntropySize]byte{1, 2, 3},
		Signs:            signs,
	}

	err = db.Append(ctx, writePulse)
	require.NoError(t, err)

	readPulse, err := db.ForPulseNumber(ctx, pn)
	require.NoError(t, err)
	require.Equal(t, writePulse, readPulse)
}
