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

	"github.com/google/gofuzz"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
)

func initDB(t *testing.T, testPulse insolar.PulseNumber) (executor.JetKeeper, string, *store.BadgerDB, *jet.DBStore) {
	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")

	require.NoError(t, err)

	db, err := store.NewBadgerDB(tmpdir)
	require.NoError(t, err)

	jets := jet.NewDBStore(db)
	pulses := pulse.NewDB(db)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: insolar.GenesisPulse.PulseNumber})
	require.NoError(t, err)

	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: testPulse})
	require.NoError(t, err)

	jetKeeper := executor.NewJetKeeper(jets, db, pulses)

	return jetKeeper, tmpdir, db, jets
}

func TestJetInfoIsConfirmed_OneDropOneHot(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	testPulse := insolar.GenesisPulse.PulseNumber + 10
	ji, tmpDir, db, jets := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(ctx)

	testJet := gen.JetID()

	err := jets.Update(ctx, testPulse, true, testJet)
	require.NoError(t, err)

	err = ji.AddHotConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)
	require.Equal(t, insolar.GenesisPulse.PulseNumber, ji.TopSyncPulse())

	err = ji.AddDropConfirmation(ctx, testPulse, testJet, false)
	require.NoError(t, err)
	require.Equal(t, testPulse, ji.TopSyncPulse())
}

func TestJetInfoIsConfirmed_Split(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	testPulse := insolar.GenesisPulse.PulseNumber + 10
	ji, tmpDir, db, jets := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(ctx)

	testJet := gen.JetID()

	err := jets.Update(ctx, testPulse, true, testJet)
	require.NoError(t, err)
	left, right, err := jets.Split(ctx, testPulse, testJet)
	require.NoError(t, err)

	err = ji.AddHotConfirmation(ctx, testPulse, left, true)
	require.NoError(t, err)
	require.Equal(t, insolar.GenesisPulse.PulseNumber, ji.TopSyncPulse())

	err = ji.AddHotConfirmation(ctx, testPulse, right, true)
	require.NoError(t, err)
	require.Equal(t, insolar.GenesisPulse.PulseNumber, ji.TopSyncPulse())

	err = ji.AddDropConfirmation(ctx, testPulse, testJet, true)
	require.NoError(t, err)
	require.Equal(t, testPulse, ji.TopSyncPulse())
}

func TestJetInfoIsConfirmed_Split_And_DifferentOrderOfComingConfirmations(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	testPulse := insolar.GenesisPulse.PulseNumber + 10
	ji, tmpDir, db, jets := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(ctx)

	testJet := gen.JetID()

	err := jets.Update(ctx, testPulse, true, testJet)
	require.NoError(t, err)
	left, right, err := jets.Split(ctx, testPulse, testJet)
	require.NoError(t, err)

	err = ji.AddHotConfirmation(ctx, testPulse, left, true)
	require.NoError(t, err)
	require.Equal(t, insolar.GenesisPulse.PulseNumber, ji.TopSyncPulse())

	err = ji.AddDropConfirmation(ctx, testPulse, testJet, true)
	require.NoError(t, err)
	require.Equal(t, insolar.GenesisPulse.PulseNumber, ji.TopSyncPulse())

	err = ji.AddHotConfirmation(ctx, testPulse, right, true)
	require.NoError(t, err)
	require.Equal(t, testPulse, ji.TopSyncPulse())
}

func TestJetInfo_ExistingDrop(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	testPulse := insolar.GenesisPulse.PulseNumber + 10
	jetKeeper, tmpDir, db, _ := initDB(t, testPulse)
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
	jetKeeper, tmpDir, db, _ := initDB(t, testPulse)
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
	jetKeeper, tmpDir, db, _ := initDB(t, testPulse)
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

	db, err := store.NewBadgerDB(tmpdir)
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
	jetKeeper, tmpDir, db, jets := initDB(t, testPulse)
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
	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())
}

func TestDbJetKeeper_DifferentBymberOfActualAndExpectedJets(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	testPulse := insolar.GenesisPulse.PulseNumber + 10
	jetKeeper, tmpDir, db, jets := initDB(t, testPulse)
	defer os.RemoveAll(tmpDir)
	defer db.Stop(context.Background())

	testJet := gen.JetID()
	left, right := jet.Siblings(testJet)

	err := jets.Update(ctx, testPulse, true, testJet)
	require.NoError(t, err)

	err = jetKeeper.AddHotConfirmation(ctx, testPulse, left, true)
	require.NoError(t, err)

	err = jetKeeper.AddHotConfirmation(ctx, testPulse, right, true)
	require.NoError(t, err)

	err = jetKeeper.AddDropConfirmation(ctx, testPulse, testJet, true)
	require.NoError(t, err)
	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())
}

func TestDbJetKeeper_AddDropConfirmation(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	db, err := store.NewBadgerDB(tmpdir)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	jets := jet.NewDBStore(db)
	pulses := pulse.NewCalculatorMock(t)
	pulses.BackwardsFunc = func(p context.Context, p1 insolar.PulseNumber, p2 int) (r insolar.Pulse, r1 error) {
		return insolar.Pulse{PulseNumber: p1 - insolar.PulseNumber(p2)}, nil
	}
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

func TestDbJetKeeper_TopSyncPulse(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	db, err := store.NewBadgerDB(tmpdir)
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
		jet          insolar.JetID
	)
	currentPulse = insolar.GenesisPulse.PulseNumber + 10
	nextPulse = insolar.GenesisPulse.PulseNumber + 20
	jet = insolar.ZeroJetID

	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: currentPulse})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: nextPulse})
	require.NoError(t, err)

	err = jets.Update(ctx, currentPulse, true, jet)
	require.NoError(t, err)
	err = jetKeeper.AddDropConfirmation(ctx, currentPulse, jet, false)
	require.NoError(t, err)
	// it's still top confirmed
	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())

	err = jetKeeper.AddHotConfirmation(ctx, currentPulse, jet, false)
	require.NoError(t, err)
	require.Equal(t, currentPulse, jetKeeper.TopSyncPulse())

	err = jets.Clone(ctx, currentPulse, nextPulse, true)
	require.NoError(t, err)
	left, right, err := jets.Split(ctx, nextPulse, jet)
	require.NoError(t, err)

	err = jetKeeper.AddDropConfirmation(ctx, nextPulse, jet, true)
	require.NoError(t, err)

	err = jetKeeper.AddHotConfirmation(ctx, nextPulse, right, true)
	require.NoError(t, err)
	require.Equal(t, currentPulse, jetKeeper.TopSyncPulse())
	err = jetKeeper.AddHotConfirmation(ctx, nextPulse, left, true)
	require.NoError(t, err)
	require.Equal(t, nextPulse, jetKeeper.TopSyncPulse())
}

func TestDbJetKeeper_TopSyncPulse_FinalizeMultiple(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	db, err := store.NewBadgerDB(tmpdir)
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
		jet          insolar.JetID
	)
	currentPulse = insolar.GenesisPulse.PulseNumber + 10
	nextPulse = insolar.GenesisPulse.PulseNumber + 20
	futurePulse = insolar.GenesisPulse.PulseNumber + 30
	jet = insolar.ZeroJetID

	err = jets.Update(ctx, currentPulse, true, jet)
	require.NoError(t, err)

	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: currentPulse})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: nextPulse})
	require.NoError(t, err)
	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())

	// Complete currentPulse pulse
	{
		err = jetKeeper.AddHotConfirmation(ctx, currentPulse, jet, false)
		require.NoError(t, err)
		require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())
		err = jetKeeper.AddDropConfirmation(ctx, currentPulse, jet, false)
		require.NoError(t, err)
		require.Equal(t, currentPulse, jetKeeper.TopSyncPulse())
	}

	err = jets.Clone(ctx, currentPulse, nextPulse, true)
	require.NoError(t, err)
	left, right, err := jets.Split(ctx, nextPulse, jet)
	require.NoError(t, err)

	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: futurePulse})
	require.NoError(t, err)

	// Complete future pulse
	{
		err = jets.Clone(ctx, nextPulse, futurePulse, true)
		require.NoError(t, err)
		leftFuture, rightFuture, err := jets.Split(ctx, futurePulse, left)
		require.NoError(t, err)
		err = jetKeeper.AddDropConfirmation(ctx, futurePulse, left, true)
		require.NoError(t, err)

		require.Equal(t, currentPulse, jetKeeper.TopSyncPulse())
		err = jetKeeper.AddDropConfirmation(ctx, futurePulse, right, false)
		require.Equal(t, currentPulse, jetKeeper.TopSyncPulse())

		err = jetKeeper.AddHotConfirmation(ctx, futurePulse, rightFuture, true)
		require.NoError(t, err)
		require.Equal(t, currentPulse, jetKeeper.TopSyncPulse())
		err = jetKeeper.AddHotConfirmation(ctx, futurePulse, leftFuture, true)
		require.NoError(t, err)
		require.Equal(t, currentPulse, jetKeeper.TopSyncPulse())
		err = jetKeeper.AddHotConfirmation(ctx, futurePulse, right, false)
		require.NoError(t, err)
		require.Equal(t, currentPulse, jetKeeper.TopSyncPulse())
	}

	// complete next pulse
	{
		err = jetKeeper.AddDropConfirmation(ctx, nextPulse, jet, true)
		require.NoError(t, err)
		require.Equal(t, currentPulse, jetKeeper.TopSyncPulse())

		err = jetKeeper.AddHotConfirmation(ctx, nextPulse, left, true)
		require.NoError(t, err)
		require.Equal(t, currentPulse, jetKeeper.TopSyncPulse())
		err = jetKeeper.AddHotConfirmation(ctx, nextPulse, right, true)
		require.NoError(t, err)
	}

	require.Equal(t, futurePulse, jetKeeper.TopSyncPulse())

}

// func TestDbJetKeeper_SubscribeUpdate(t *testing.T) {
// 	var (
// 		targetPulse = insolar.GenesisPulse.PulseNumber
// 		handler     = func(present insolar.PulseNumber) {
// 			require.Equal(t, targetPulse, present)
// 		}
// 		db        = store.NewMemoryMockDB()
// 		jets      = jet.NewDBStore(db)
// 		jetKeeper = NewJetKeeper(jets, db)
// 	)
//
// 	jetKeeper.Subscribe(targetPulse, handler)
// 	err := jetKeeper.Update(targetPulse)
// 	require.NoError(t, err)
// }
