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

package executor_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/dgraph-io/badger"
	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func BadgerDefaultOptions(dir string) badger.Options {
	ops := badger.DefaultOptions(dir)
	ops.CompactL0OnClose = false
	ops.SyncWrites = false

	return ops
}

func initDB(t *testing.T, testPulse insolar.PulseNumber) (executor.JetKeeper, string, *store.BadgerDB, *jet.DBStore, *pulse.DB) {
	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")

	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)

	jets := jet.NewDBStore(db)
	pulses := pulse.NewDB(db)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: insolar.GenesisPulse.PulseNumber})
	require.NoError(t, err)

	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: testPulse})
	require.NoError(t, err)

	jetKeeper := executor.NewJetKeeper(jets, db, pulses)

	return jetKeeper, tmpdir, db, jets, pulses
}

func Test_TruncateHead_TryToTruncateTopSync(t *testing.T) {
	ctx := inslogger.TestContext(t)
	testPulse := insolar.GenesisPulse.PulseNumber + 10
	ji, tmpDir, db, _, _ := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(ctx)
	err := ji.(*executor.DBJetKeeper).TruncateHead(ctx, 1)
	require.EqualError(t, err, "try to truncate top sync pulse")
}

func TestJetInfoIsConfirmed_OneDropOneHot(t *testing.T) {
	t.Parallel()
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
	require.Equal(t, insolar.GenesisPulse.PulseNumber, ji.TopSyncPulse())

	err = ji.AddDropConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)

	require.Equal(t, insolar.GenesisPulse.PulseNumber, ji.TopSyncPulse())

	err = ji.AddBackupConfirmation(ctx, testPulse)
	require.NoError(t, err)
	require.Equal(t, testPulse, ji.TopSyncPulse())
}

func Test_DifferentSplitFlagsInDropsAndHots(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	testPulse := insolar.GenesisPulse.PulseNumber + 10
	ji, tmpDir, db, _, _ := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(ctx)

	testJet := insolar.ZeroJetID

	// AddHotConfirmation: 'true' come first
	err := ji.AddDropConfirmation(ctx, testPulse, testJet, true)
	require.NoError(t, err)
	err = ji.AddHotConfirmation(ctx, testPulse, testJet, false)
	require.Contains(t, err.Error(), "try to change split from true to false")

	// AddHotConfirmation: 'false' comes first
	left, _ := jet.Siblings(testJet)
	leftLeft, rightLeft := jet.Siblings(left)
	err = ji.AddHotConfirmation(ctx, testPulse, left, false)
	require.NoError(t, err)
	err = ji.AddHotConfirmation(ctx, testPulse, leftLeft, true)
	require.Contains(t, err.Error(), "try to change split from false to true")

	// AddDropConfirmation
	err = ji.AddHotConfirmation(ctx, testPulse, rightLeft, false)
	require.NoError(t, err)
	err = ji.AddDropConfirmation(ctx, testPulse, rightLeft, true)
	require.Contains(t, err.Error(), "try to change split from false to true")
}

func TestJetInfoIsConfirmed_Split(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	testPulse := insolar.GenesisPulse.PulseNumber + 10
	ji, tmpDir, db, jets, pulses := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(ctx)

	testJet := insolar.ZeroJetID

	nextPulse := insolar.GenesisPulse.PulseNumber + 20

	err := jets.Update(ctx, testPulse, true, testJet)
	require.NoError(t, err)
	err = ji.AddHotConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)
	require.Equal(t, insolar.GenesisPulse.PulseNumber, ji.TopSyncPulse())
	err = ji.AddDropConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)
	require.Equal(t, insolar.GenesisPulse.PulseNumber, ji.TopSyncPulse())
	err = ji.AddBackupConfirmation(ctx, testPulse)
	require.NoError(t, err)
	require.Equal(t, testPulse, ji.TopSyncPulse())

	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: nextPulse})
	require.NoError(t, err)

	left, right := jet.Siblings(testJet)
	err = jets.Update(ctx, nextPulse, true, testJet)
	require.NoError(t, err)
	err = ji.AddDropConfirmation(ctx, nextPulse, testJet, true)
	require.NoError(t, err)
	err = ji.AddHotConfirmation(ctx, nextPulse, left, true)
	require.NoError(t, err)
	err = ji.AddHotConfirmation(ctx, nextPulse, right, true)
	require.NoError(t, err)
	err = ji.AddBackupConfirmation(ctx, nextPulse)
	require.NoError(t, err)
	require.Equal(t, nextPulse, ji.TopSyncPulse())
}

func TestJetInfo_BackupConfirmComesFirst(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	testPulse := insolar.GenesisPulse.PulseNumber + 10
	jetKeeper, tmpDir, db, _, _ := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(context.Background())

	err := jetKeeper.AddBackupConfirmation(ctx, testPulse)
	require.Contains(t, err.Error(), "Received backup confirmation before replication data")
}

func TestJetInfo_ExistingDrop(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	testPulse := insolar.GenesisPulse.PulseNumber + 10
	jetKeeper, tmpDir, db, _, _ := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(ctx)

	testJet := gen.JetID()
	err := jetKeeper.AddDropConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)

	err = jetKeeper.AddDropConfirmation(ctx, testPulse, testJet, false)
	require.Contains(t, err.Error(), "try to rewrite drop confirmation")
	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())
}

func TestJetInfo_ExistingHot(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	testPulse := insolar.GenesisPulse.PulseNumber + 10
	jetKeeper, tmpDir, db, _, _ := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(ctx)

	testJet := gen.JetID()
	err := jetKeeper.AddHotConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)

	err = jetKeeper.AddHotConfirmation(ctx, testPulse, testJet, false)
	require.Contains(t, err.Error(), "try add already existing hot confirmation")
}

func TestJetInfo_ExceedNumHotConfirmations(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	testPulse := insolar.GenesisPulse.PulseNumber + 10
	jetKeeper, tmpDir, db, _, _ := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(context.Background())

	testJet := gen.JetID()
	left, right := jet.Siblings(testJet)

	err := jetKeeper.AddHotConfirmation(ctx, testPulse, left, true)
	require.NoError(t, err)

	err = jetKeeper.AddHotConfirmation(ctx, testPulse, right, true)
	require.NoError(t, err)

	err = jetKeeper.AddHotConfirmation(ctx, testPulse, left, true)
	require.Contains(t, err.Error(), "num hot confirmations exceeds")
	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())
}

func TestNewJetKeeper(t *testing.T) {
	t.Parallel()
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	jets := jet.NewDBStore(db)
	pulses := pulse.NewCalculatorMock(t)
	jetKeeper := executor.NewJetKeeper(jets, db, pulses)
	require.NotNil(t, jetKeeper)
}

func TestDbJetKeeper_DifferentActualAndExpectedJets(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	testPulse := insolar.GenesisPulse.PulseNumber + 10
	jetKeeper, tmpDir, db, jets, _ := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(context.Background())

	testJet := gen.JetID()
	left, _ := jet.Siblings(testJet)

	err := jets.Update(ctx, testPulse, true, left)
	require.NoError(t, err)

	err = jetKeeper.AddHotConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)

	err = jetKeeper.AddDropConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)

	require.False(t, jetKeeper.HasAllJetConfirms(ctx, testPulse))

	err = jetKeeper.AddBackupConfirmation(ctx, testPulse)
	require.NoError(t, err)
	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())
	require.False(t, jetKeeper.HasAllJetConfirms(ctx, testPulse))
}

func TestDbJetKeeper_DifferentNumberOfActualAndExpectedJets(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	testPulse := insolar.GenesisPulse.PulseNumber + 10
	jetKeeper, tmpDir, db, jets, _ := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(context.Background())

	testJet := gen.JetID()
	left, right := jet.Siblings(testJet)

	err := jets.Update(ctx, testPulse, true, testJet)
	require.NoError(t, err)

	err = jetKeeper.AddHotConfirmation(ctx, testPulse, left, false)
	require.NoError(t, err)

	err = jetKeeper.AddHotConfirmation(ctx, testPulse, right, false)
	require.NoError(t, err)

	err = jetKeeper.AddDropConfirmation(ctx, testPulse, right, false)
	require.NoError(t, err)

	err = jetKeeper.AddDropConfirmation(ctx, testPulse, left, false)
	require.NoError(t, err)

	err = jetKeeper.AddBackupConfirmation(ctx, testPulse)
	require.NoError(t, err)

	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())
}

func TestDbJetKeeper_AddDropConfirmation(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	jets := jet.NewDBStore(db)
	pulses := pulse.NewCalculatorMock(t)
	pulses.BackwardsMock.Set(func(p context.Context, p1 insolar.PulseNumber, p2 int) (r insolar.Pulse, r1 error) {
		return insolar.Pulse{PulseNumber: p1 - insolar.PulseNumber(p2)}, nil
	})
	jetKeeper := executor.NewJetKeeper(jets, db, pulses)

	var (
		pulse insolar.PulseNumber
		jet   insolar.JetID
	)
	f := fuzz.New()
	f.Fuzz(&pulse)
	f.Fuzz(&jet)
	err = jetKeeper.AddDropConfirmation(ctx, pulse, jet, false)
	require.NoError(t, err)
}

func TestDbJetKeeper_CheckJetTreeFail(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	testPulse := insolar.GenesisPulse.PulseNumber + 10
	ji, tmpDir, db, _, _ := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(ctx)

	testJet := insolar.ZeroJetID

	err := ji.AddHotConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)
	require.Equal(t, insolar.GenesisPulse.PulseNumber, ji.TopSyncPulse())
	err = ji.AddDropConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)
	require.Equal(t, insolar.GenesisPulse.PulseNumber, ji.TopSyncPulse())
	err = ji.AddBackupConfirmation(ctx, testPulse)
	require.NoError(t, err)
	require.Equal(t, insolar.GenesisPulse.PulseNumber, ji.TopSyncPulse())
	require.False(t, false, ji.HasAllJetConfirms(ctx, testPulse))
}

func TestDbJetKeeper_TopSyncPulse(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	jets := jet.NewDBStore(db)
	pulses := pulse.NewDB(db)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: insolar.GenesisPulse.PulseNumber})
	require.NoError(t, err)

	jetKeeper := executor.NewJetKeeper(jets, db, pulses)

	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())

	var (
		currentPulse insolar.PulseNumber
		nextPulse    insolar.PulseNumber
		testJet      insolar.JetID
	)
	currentPulse = insolar.GenesisPulse.PulseNumber + 10
	nextPulse = insolar.GenesisPulse.PulseNumber + 20
	testJet = insolar.ZeroJetID

	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: currentPulse})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: nextPulse})
	require.NoError(t, err)

	err = jets.Update(ctx, currentPulse, true, testJet)
	require.NoError(t, err)
	err = jetKeeper.AddDropConfirmation(ctx, currentPulse, testJet, false)
	require.NoError(t, err)
	// it's still top confirmed
	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())

	err = jetKeeper.AddHotConfirmation(ctx, currentPulse, testJet, false)
	require.NoError(t, err)
	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())

	err = jetKeeper.AddBackupConfirmation(ctx, currentPulse)
	require.NoError(t, err)
	require.Equal(t, currentPulse, jetKeeper.TopSyncPulse())

	err = jets.Clone(ctx, currentPulse, nextPulse, true)
	require.NoError(t, err)
	left, right := jet.Siblings(testJet)

	err = jetKeeper.AddDropConfirmation(ctx, nextPulse, testJet, true)
	require.NoError(t, err)

	err = jetKeeper.AddHotConfirmation(ctx, nextPulse, right, true)
	require.NoError(t, err)
	require.Equal(t, currentPulse, jetKeeper.TopSyncPulse())
	err = jetKeeper.AddHotConfirmation(ctx, nextPulse, left, true)
	require.NoError(t, err)
	require.Equal(t, currentPulse, jetKeeper.TopSyncPulse())

	err = jetKeeper.AddBackupConfirmation(ctx, nextPulse)
	require.NoError(t, err)
	require.Equal(t, nextPulse, jetKeeper.TopSyncPulse())
}

func TestDbJetKeeper_LostDataOnNextPulseAfterSplit(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	jets := jet.NewDBStore(db)
	pulses := pulse.NewDB(db)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: insolar.GenesisPulse.PulseNumber})
	require.NoError(t, err)

	jetKeeper := executor.NewJetKeeper(jets, db, pulses)

	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())

	var (
		currentPulse insolar.PulseNumber
		nextPulse    insolar.PulseNumber
		futurePulse  insolar.PulseNumber
		testJet      insolar.JetID
	)
	currentPulse = insolar.GenesisPulse.PulseNumber + 10
	nextPulse = insolar.GenesisPulse.PulseNumber + 20
	futurePulse = insolar.GenesisPulse.PulseNumber + 30
	testJet = insolar.ZeroJetID

	err = jets.Update(ctx, currentPulse, true, testJet)
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: currentPulse})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: nextPulse})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: futurePulse})
	require.NoError(t, err)

	// finalize currentPulse
	{
		err = jetKeeper.AddHotConfirmation(ctx, currentPulse, testJet, false)
		require.NoError(t, err)
		err = jetKeeper.AddDropConfirmation(ctx, currentPulse, testJet, false)
		require.NoError(t, err)
		require.True(t, jetKeeper.HasAllJetConfirms(ctx, currentPulse))
		err = jetKeeper.AddBackupConfirmation(ctx, currentPulse)
		require.Equal(t, currentPulse, jetKeeper.TopSyncPulse())
	}

	left, right := jet.Siblings(testJet)
	// finalize nextPulse
	{
		err = jets.Update(ctx, nextPulse, true, testJet)
		require.NoError(t, err)
		err = jetKeeper.AddDropConfirmation(ctx, nextPulse, testJet, true)
		require.NoError(t, err)
		err = jetKeeper.AddHotConfirmation(ctx, nextPulse, left, true)
		require.NoError(t, err)
		require.False(t, jetKeeper.HasAllJetConfirms(ctx, nextPulse))
		err = jetKeeper.AddHotConfirmation(ctx, nextPulse, right, true)
		require.NoError(t, err)

		require.True(t, jetKeeper.HasAllJetConfirms(ctx, currentPulse))
		require.True(t, jetKeeper.HasAllJetConfirms(ctx, nextPulse))
		err = jetKeeper.AddBackupConfirmation(ctx, nextPulse)
		require.NoError(t, err)
		require.Equal(t, nextPulse, jetKeeper.TopSyncPulse())
	}

	err = jets.Update(ctx, futurePulse, true, left)
	require.NoError(t, err)
	err = jetKeeper.AddDropConfirmation(ctx, futurePulse, left, false)
	require.NoError(t, err)
	err = jetKeeper.AddHotConfirmation(ctx, futurePulse, left, false)
	require.NoError(t, err)
	require.True(t, jetKeeper.HasAllJetConfirms(ctx, currentPulse))
	require.False(t, jetKeeper.HasAllJetConfirms(ctx, futurePulse))

	err = jets.Update(ctx, futurePulse, true, right)
	err = jetKeeper.AddDropConfirmation(ctx, futurePulse, right, false)
	require.NoError(t, err)
	err = jetKeeper.AddHotConfirmation(ctx, futurePulse, right, false)
	require.NoError(t, err)

	require.True(t, jetKeeper.HasAllJetConfirms(ctx, currentPulse))
	require.True(t, jetKeeper.HasAllJetConfirms(ctx, nextPulse))
	require.True(t, jetKeeper.HasAllJetConfirms(ctx, futurePulse))

	err = jetKeeper.AddBackupConfirmation(ctx, futurePulse)
	require.NoError(t, err)
	require.Equal(t, futurePulse, jetKeeper.TopSyncPulse())
}
