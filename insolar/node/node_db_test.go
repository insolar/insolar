//
// Copyright 2020 Insolar Technologies GmbH
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

package node

import (
	"context"
	"os"
	"sort"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/ledger/heavy/migration"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/tests/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

var db *StorageDB

// TestMain does the before and after setup
func TestMain(m *testing.M) {
	ctx := context.Background()
	pgURL, stopDBMS := common.StartDBMS()

	pool, err := pgxpool.Connect(ctx, pgURL)
	if err != nil {
		stopDBMS()
		log.Panicf("[TestMain] pgxpool.Connect() failed: %v", err)
	}

	migrationPath := "../../migration"
	cwd, err := os.Getwd()
	if err != nil {
		stopDBMS()
		panic(errors.Wrap(err, "[TestMain] os.Getwd failed"))
	}
	log.Infof("[TestMain] About to run PostgreSQL migration, cwd = %s, migration migrationPath = %s", cwd, migrationPath)
	ver, err := migration.MigrateDatabase(ctx, pool, migrationPath)
	if err != nil {
		stopDBMS()
		panic(errors.Wrap(err, "Unable to migrate database"))
	}
	log.Infof("[TestMain] PostgreSQL database migration done, current schema version: %d", ver)

	db = NewStorageDB(pool)

	// Run all tests
	code := m.Run()

	log.Info("[TestMain] Cleaning up...")
	stopDBMS()
	os.Exit(code)
}

func TestInsertSelect(t *testing.T) {
	pn := gen.PulseNumber()
	// Make sure there are no nodes for a given pulse yet
	{
		readNodes, err := db.All(pn)
		require.NoError(t, err)
		require.Empty(t, readNodes)
	}
	{
		readNodes, err := db.InRole(pn, insolar.StaticRoleVirtual)
		require.NoError(t, err)
		require.Empty(t, readNodes)
	}

	// Insert nodes for a given pulse
	nodes := []insolar.Node{
		{
			Polymorph: 123,
			ID:        gen.Reference(),
			Role:      insolar.StaticRoleVirtual,
		},
		{
			Polymorph: 123,
			ID:        gen.Reference(),
			Role:      insolar.StaticRoleHeavyMaterial,
		},
		{
			Polymorph: 123,
			ID:        gen.Reference(),
			Role:      insolar.StaticRoleLightMaterial,
		},
		{
			Polymorph: 123,
			ID:        gen.Reference(),
			Role:      insolar.StaticRoleLightMaterial,
		},
	}
	err := db.Set(pn, nodes)
	require.NoError(t, err)

	// Make sure .All returns all nodes in the same order as saved
	{
		readNodes, err := db.All(pn)
		require.NoError(t, err)
		require.NotEmpty(t, readNodes)
		require.Equal(t, nodes, readNodes)
	}

	// Make sure .InRole returns only nodes that have a given role
	{
		readNodes, err := db.InRole(pn, insolar.StaticRoleVirtual)
		require.NoError(t, err)
		require.Equal(t, 1, len(readNodes))
		require.Equal(t, nodes[0], readNodes[0])
	}
	{
		readNodes, err := db.InRole(pn, insolar.StaticRoleHeavyMaterial)
		require.NoError(t, err)
		require.Equal(t, 1, len(readNodes))
		require.Equal(t, nodes[1], readNodes[0])
	}
	{
		readNodes, err := db.InRole(pn, insolar.StaticRoleLightMaterial)
		require.NoError(t, err)
		require.Equal(t, 2, len(readNodes))
		require.Equal(t, nodes[2], readNodes[0])
		require.Equal(t, nodes[3], readNodes[1])
	}
}

func TestTruncateHead(t *testing.T) {
	// Generate sorted list of pulses
	pulses := make([]insolar.PulseNumber, 5)
	for i := 0; i < len(pulses); i++ {
		pulses[i] = gen.PulseNumber()
	}
	sort.Slice(pulses, func(i, j int) bool {
		return pulses[i] < pulses[j]
	})

	// Insert some nodes for each pulse
	nodes := make([][]insolar.Node, len(pulses))
	for i := 0; i < len(pulses); i++ {
		nodes[i] = []insolar.Node{
			{
				Polymorph: 123,
				ID:        gen.Reference(),
				Role:      insolar.StaticRoleVirtual,
			},
		}
		err := db.Set(pulses[i], nodes[i])
		require.NoError(t, err)
	}

	// Make sure data for all pulses is present
	for i := 0; i < len(pulses); i++ {
		readNodes, err := db.All(pulses[i])
		require.NoError(t, err)
		require.Equal(t, nodes[i], readNodes)
	}

	// Truncate head
	err := db.TruncateHead(context.Background(), pulses[2])
	require.NoError(t, err)

	// Make sure nodes for pulses [0,1] is still here
	for i := 0; i < 2; i++ {
		readNodes, err := db.All(pulses[i])
		require.NoError(t, err)
		require.Equal(t, nodes[i], readNodes)
	}

	// Make sure nodes for pulses [2,3,4] were deleted
	for i := 2; i <= 4; i++ {
		readNodes, err := db.All(pulses[i])
		require.NoError(t, err)
		require.Empty(t, readNodes)
	}
}
