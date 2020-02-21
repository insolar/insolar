// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/insolar/insolar/ledger/object"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

func BadgerDefaultOptions(dir string) badger.Options {
	ops := badger.DefaultOptions(dir)
	ops.CompactL0OnClose = false
	ops.SyncWrites = false

	return ops
}

func initDB(t *testing.T, testPulse insolar.PulseNumber) (JetKeeper, string, *store.BadgerDB, *jet.BadgerDBStore, *pulse.BadgerDB) {
	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")

	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)

	jets := jet.NewBadgerDBStore(db)
	txManager, err := object.NewBadgerTxManager(db.Backend())
	require.NoError(t, err)
	pulses := pulse.NewBadgerDB(db, txManager)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: insolar.GenesisPulse.PulseNumber})
	require.NoError(t, err)

	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: testPulse})
	require.NoError(t, err)

	jetKeeper := NewBadgerJetKeeper(jets, db, pulses)

	return jetKeeper, tmpdir, db, jets, pulses
}

func Test_JetKeeperKey(t *testing.T) {
	k := jetKeeperKey(insolar.GenesisPulse.PulseNumber)
	d := k.ID()
	require.Equal(t, k, newJetKeeperKey(d))
}

func Test_TruncateHead(t *testing.T) {
	ctx := inslogger.TestContext(t)
	testPulse := insolar.GenesisPulse.PulseNumber + 10
	ji, tmpDir, db, jets, _ := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(ctx)

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

	_, err = db.Get(jetKeeperKey(testPulse))
	require.NoError(t, err)

	nextPulse := testPulse + 10

	err = ji.AddDropConfirmation(ctx, nextPulse, gen.JetID(), false)
	require.NoError(t, err)
	err = ji.AddHotConfirmation(ctx, nextPulse, gen.JetID(), false)
	require.NoError(t, err)

	_, err = db.Get(jetKeeperKey(nextPulse))
	require.NoError(t, err)

	err = ji.(*BadgerDBJetKeeper).TruncateHead(ctx, nextPulse)
	require.NoError(t, err)

	_, err = db.Get(jetKeeperKey(testPulse))
	require.NoError(t, err)
	_, err = db.Get(jetKeeperKey(nextPulse))
	require.EqualError(t, err, "value not found")
}
