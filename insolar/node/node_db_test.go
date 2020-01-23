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
	"testing"

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

	db = NewStorageDB(pool)

	// Run all tests
	code := m.Run()

	log.Info("[TestMain] Cleaning up...")
	stopPostgreSQL()
	os.Exit(code)
}

//func BadgerDefaultOptions(dir string) badger.Options {
//	ops := badger.DefaultOptions(dir)
//	ops.CompactL0OnClose = false
//	ops.SyncWrites = false
//
//	return ops
//}
//
//func TestNodeStorageDB_All(t *testing.T) {
//	t.Parallel()
//
//	var all []insolar.Node
//	f := fuzz.New().Funcs(func(e *insolar.Node, c fuzz.Continue) {
//		e.ID = gen.Reference()
//	})
//	f.NumElements(5, 10).NilChance(0).Fuzz(&all)
//	pulse := gen.PulseNumber()
//
//	t.Run("returns correct nodes", func(t *testing.T) {
//		tmpdir, err := ioutil.TempDir("", "bdb-test-")
//		defer os.RemoveAll(tmpdir)
//		require.NoError(t, err)
//
//		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
//		require.NoError(t, err)
//		defer db.Stop(context.Background())
//
//		nodeStorage := NewStorageDB(db, pool)
//		err = nodeStorage.Set(pulse, all)
//		require.NoError(t, err)
//
//		result, err := nodeStorage.All(pulse)
//
//		require.NoError(t, err)
//		require.Equal(t, all, result)
//	})
//
//	t.Run("returns nil when empty nodes", func(t *testing.T) {
//		tmpdir, err := ioutil.TempDir("", "bdb-test-")
//		defer os.RemoveAll(tmpdir)
//		require.NoError(t, err)
//
//		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
//		require.NoError(t, err)
//		defer db.Stop(context.Background())
//
//		nodeStorage := NewStorageDB(db)
//		err = nodeStorage.Set(pulse, nil)
//		require.NoError(t, err)
//
//		result, err := nodeStorage.All(pulse)
//
//		require.NoError(t, err)
//		require.Nil(t, result)
//	})
//
//	t.Run("returns error when no nodes", func(t *testing.T) {
//		tmpdir, err := ioutil.TempDir("", "bdb-test-")
//		defer os.RemoveAll(tmpdir)
//		require.NoError(t, err)
//
//		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
//		require.NoError(t, err)
//		defer db.Stop(context.Background())
//
//		nodeStorage := NewStorageDB(db)
//
//		result, err := nodeStorage.All(pulse)
//
//		require.Equal(t, ErrNoNodes, err)
//		require.Nil(t, result)
//	})
//}
//
//func TestNodeStorageDB_InRole(t *testing.T) {
//	t.Parallel()
//
//	var (
//		virtuals  []insolar.Node
//		materials []insolar.Node
//		all       []insolar.Node
//	)
//	{
//		f := fuzz.New().Funcs(func(e *insolar.Node, c fuzz.Continue) {
//			e.ID = gen.Reference()
//			e.Role = insolar.StaticRoleVirtual
//		})
//		f.NumElements(5, 10).NilChance(0).Fuzz(&virtuals)
//	}
//	{
//		f := fuzz.New().Funcs(func(e *insolar.Node, c fuzz.Continue) {
//			e.ID = gen.Reference()
//			e.Role = insolar.StaticRoleLightMaterial
//		})
//		f.NumElements(5, 10).NilChance(0).Fuzz(&materials)
//	}
//	all = append(virtuals, materials...)
//	pulse := gen.PulseNumber()
//
//	t.Run("returns correct nodes", func(t *testing.T) {
//		tmpdir, err := ioutil.TempDir("", "bdb-test-")
//		defer os.RemoveAll(tmpdir)
//		require.NoError(t, err)
//
//		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
//		require.NoError(t, err)
//		defer db.Stop(context.Background())
//
//		nodeStorage := NewStorageDB(db)
//		err = nodeStorage.Set(pulse, all)
//		require.NoError(t, err)
//		{
//			result, err := nodeStorage.InRole(pulse, insolar.StaticRoleVirtual)
//			require.NoError(t, err)
//			require.Equal(t, virtuals, result)
//		}
//		{
//			result, err := nodeStorage.InRole(pulse, insolar.StaticRoleLightMaterial)
//			require.NoError(t, err)
//			require.Equal(t, materials, result)
//		}
//	})
//
//	t.Run("returns nil when empty nodes", func(t *testing.T) {
//		tmpdir, err := ioutil.TempDir("", "bdb-test-")
//		defer os.RemoveAll(tmpdir)
//		require.NoError(t, err)
//
//		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
//		require.NoError(t, err)
//		defer db.Stop(context.Background())
//
//		nodeStorage := NewStorageDB(db)
//		err = nodeStorage.Set(pulse, nil)
//		result, err := nodeStorage.InRole(pulse, insolar.StaticRoleVirtual)
//		assert.NoError(t, err)
//		assert.Nil(t, result)
//	})
//
//	t.Run("returns error when no nodes", func(t *testing.T) {
//		tmpdir, err := ioutil.TempDir("", "bdb-test-")
//		defer os.RemoveAll(tmpdir)
//		require.NoError(t, err)
//
//		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
//		require.NoError(t, err)
//		defer db.Stop(context.Background())
//
//		nodeStorage := NewStorageDB(db)
//		result, err := nodeStorage.InRole(pulse, insolar.StaticRoleVirtual)
//		assert.Equal(t, ErrNoNodes, err)
//		assert.Nil(t, result)
//	})
//}
//
//func TestNodeStorageDB_Set(t *testing.T) {
//	t.Parallel()
//
//	var nodes []insolar.Node
//	f := fuzz.New().Funcs(func(e *insolar.Node, c fuzz.Continue) {
//		e.ID = gen.Reference()
//	})
//	f.NumElements(5, 10).NilChance(0).Fuzz(&nodes)
//	pulse := gen.PulseNumber()
//
//	t.Run("saves correct nodes", func(t *testing.T) {
//		tmpdir, err := ioutil.TempDir("", "bdb-test-")
//		defer os.RemoveAll(tmpdir)
//		require.NoError(t, err)
//
//		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
//		require.NoError(t, err)
//		defer db.Stop(context.Background())
//
//		nodeStorage := NewStorageDB(db)
//
//		err = nodeStorage.Set(pulse, nodes)
//		require.NoError(t, err)
//
//		res, err := nodeStorage.All(pulse)
//
//		require.NoError(t, err)
//		require.Equal(t, nodes, res)
//	})
//
//	t.Run("saves nil if empty nodes", func(t *testing.T) {
//		tmpdir, err := ioutil.TempDir("", "bdb-test-")
//		defer os.RemoveAll(tmpdir)
//		require.NoError(t, err)
//
//		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
//		require.NoError(t, err)
//		defer db.Stop(context.Background())
//
//		nodeStorage := NewStorageDB(db)
//
//		err = nodeStorage.Set(pulse, []insolar.Node{})
//		require.NoError(t, err)
//
//		res, err := nodeStorage.All(pulse)
//
//		require.NoError(t, err)
//		require.Nil(t, res)
//	})
//
//	t.Run("returns error when saving with the same pulse", func(t *testing.T) {
//		tmpdir, err := ioutil.TempDir("", "bdb-test-")
//		defer os.RemoveAll(tmpdir)
//		require.NoError(t, err)
//
//		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
//		require.NoError(t, err)
//		defer db.Stop(context.Background())
//
//		nodeStorage := NewStorageDB(db)
//
//		_ = nodeStorage.Set(pulse, nodes)
//		err = nodeStorage.Set(pulse, nodes)
//		require.Equal(t, ErrOverride, err)
//	})
//}
//
//func TestNodeStorageDB_TruncateHead_NoSuchPulse(t *testing.T) {
//	t.Parallel()
//
//	ctx := inslogger.TestContext(t)
//	tmpdir, err := ioutil.TempDir("", "bdb-test-")
//	defer os.RemoveAll(tmpdir)
//	assert.NoError(t, err)
//
//	ops := BadgerDefaultOptions(tmpdir)
//	dbMock, err := store.NewBadgerDB(ops)
//	defer dbMock.Stop(ctx)
//	require.NoError(t, err)
//
//	dropStore := NewStorageDB(dbMock)
//
//	err = dropStore.TruncateHead(ctx, insolar.GenesisPulse.PulseNumber)
//	require.NoError(t, err)
//}
//
//func TestDropStorageDB_TruncateHead(t *testing.T) {
//	t.Parallel()
//
//	ctx := inslogger.TestContext(t)
//	tmpdir, err := ioutil.TempDir("", "bdb-test-")
//	defer os.RemoveAll(tmpdir)
//	assert.NoError(t, err)
//
//	ops := BadgerDefaultOptions(tmpdir)
//	dbMock, err := store.NewBadgerDB(ops)
//	defer dbMock.Stop(ctx)
//	require.NoError(t, err)
//
//	nodeStor := NewStorageDB(dbMock)
//
//	nodeSets := make([]struct {
//		nodes []insolar.Node
//		pn    insolar.PulseNumber
//	}, 10)
//
//	for i := range nodeSets {
//		nodeSets[i].pn = pulse.Number(pulse.MinTimePulse + (i * 10))
//		nodeSets[i].nodes = []insolar.Node{
//			{
//				Role: insolar.StaticRoleHeavyMaterial,
//			},
//			{
//				Role: insolar.StaticRoleLightMaterial,
//			},
//			{
//				Role: insolar.StaticRoleVirtual,
//			},
//		}
//	}
//
//	rand.Seed(time.Now().UnixNano())
//	rand.Shuffle(len(nodeSets), func(i, j int) { nodeSets[i], nodeSets[j] = nodeSets[j], nodeSets[i] })
//
//	for _, nodeSet := range nodeSets {
//		err := nodeStor.Set(nodeSet.pn, nodeSet.nodes)
//		require.NoError(t, err)
//	}
//
//	for i := 0; i < 10; i++ {
//		_, err := nodeStor.All(pulse.Number(pulse.MinTimePulse + (i * 10)))
//		require.NoError(t, err)
//	}
//
//	numLeftElements := 10 / 2
//	err = nodeStor.TruncateHead(ctx, pulse.MinTimePulse+insolar.PulseNumber(numLeftElements*10))
//	require.NoError(t, err)
//
//	for i := 0; i < numLeftElements; i++ {
//		_, err := nodeStor.All(pulse.Number(pulse.MinTimePulse + (i * 10)))
//		require.NoError(t, err)
//	}
//
//	for i := numLeftElements - 1; i >= numLeftElements; i-- {
//		p := pulse.MinTimePulse + insolar.PulseNumber(numLeftElements*10)
//		_, err := nodeStor.All(p)
//		require.EqualError(t, err, ErrNoNodes.Error(), "Pulse: ", p.String())
//	}
//}
