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

package executor

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/gofuzz"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
)

func TestJetInfoIsConfirmed_OneDropOneHot(t *testing.T) {
	ji := jetInfo{}
	require.False(t, ji.isConfirmed())

	jet := gen.JetID()

	ji.addDrop(jet, false)
	require.False(t, ji.isConfirmed())

	ji.addHot(jet, false)
	require.True(t, ji.isConfirmed())
}

func TestJetInfoIsConfirmed_Split(t *testing.T) {
	ji := jetInfo{}
	require.False(t, ji.isConfirmed())
	testJet := gen.JetID()

	left, right := jet.Siblings(testJet)

	err := ji.addHot(left, true)
	require.NoError(t, err)
	require.False(t, ji.isConfirmed())

	err = ji.addHot(right, true)
	require.NoError(t, err)
	require.False(t, ji.isConfirmed())

	err = ji.addDrop(testJet, true)
	require.NoError(t, err)
	require.True(t, ji.isConfirmed())
}

func TestJetInfoIsConfirmed_Split_And_DifferentOrderOfComingConfirmations(t *testing.T) {
	ji := jetInfo{}
	require.False(t, ji.isConfirmed())
	testJet := gen.JetID()

	left, right := jet.Siblings(testJet)

	err := ji.addHot(left, true)
	require.NoError(t, err)
	require.False(t, ji.isConfirmed())

	err = ji.addDrop(testJet, true)
	require.NoError(t, err)
	require.False(t, ji.isConfirmed())

	err = ji.addHot(right, true)
	require.NoError(t, err)
	require.True(t, ji.isConfirmed())
}

func TestJetInfo_RewriteWithDifferentParent(t *testing.T) {
	ji := jetInfo{}
	require.False(t, ji.isConfirmed())
	testJet := gen.JetID()

	left, _ := jet.Siblings(testJet)

	err := ji.addHot(left, true)
	require.NoError(t, err)
	require.False(t, ji.isConfirmed())

	err = ji.addHot(testJet, true)
	require.Contains(t, err.Error(), "try to rewrite jet with different parent")
	require.False(t, ji.isConfirmed())
}

func TestJetInfo_ExistingHot(t *testing.T) {
	ji := jetInfo{}
	require.False(t, ji.isConfirmed())
	testJet := gen.JetID()

	left, _ := jet.Siblings(testJet)

	err := ji.addHot(left, false)
	require.NoError(t, err)

	err = ji.addHot(left, false)
	require.Contains(t, err.Error(), "try add already existing hot confirmation")
}

func TestJetInfo_ExistingDrop(t *testing.T) {
	ji := jetInfo{}
	require.False(t, ji.isConfirmed())
	testJet := gen.JetID()
	err := ji.addDrop(testJet, false)
	require.NoError(t, err)

	err = ji.addDrop(testJet, false)
	require.Contains(t, err.Error(), "try to rewrite drop confirmation")
}

func TestJetInfo_ExceedNumHotConfirmations(t *testing.T) {
	ji := jetInfo{}
	require.False(t, ji.isConfirmed())
	testJet := gen.JetID()

	left, right := jet.Siblings(testJet)

	err := ji.addHot(left, false)
	require.NoError(t, err)

	err = ji.addHot(right, false)
	require.NoError(t, err)

	err = ji.addHot(testJet, false)
	require.Contains(t, err.Error(), "num hot confirmations exceeds")
}

func TestNewJetKeeper(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	db, err := store.NewBadgerDB(tmpdir)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	jets := jet.NewDBStore(db)
	pulses := pulse.NewCalculatorMock(t)
	jetKeeper := NewJetKeeper(jets, db, pulses)
	require.NotNil(t, jetKeeper)
}

func TestDbJetKeeper_AddDropConfirmation(t *testing.T) {
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
	jetKeeper := NewJetKeeper(jets, db, pulses)

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

	jetKeeper := NewJetKeeper(jets, db, pulses)

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

	jetKeeper := NewJetKeeper(jets, db, pulses)

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

	inslogger.FromContext(ctx).Debug("INIT: JET: ", jet.DebugString())

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

// func TestDbJetKeeper_Add_CantGetPulse(t *testing.T) {
// 	ctx := inslogger.TestContext(t)
// 	dbMock := store.NewDBMock(t)
//
// 	pn := insolar.GenesisPulse.PulseNumber
//
// 	dbMock.GetMock.Expect(jetKeeperKey(pn)).Return([]byte{}, nil)
//
// 	jets := jet.NewStorageMock(t)
// 	pulses := pulse.NewCalculatorMock(t)
//
// 	jetKeeper := NewJetKeeper(jets, dbMock, pulses)
// 	err := jetKeeper.AddHotConfirmation(ctx, pn, insolar.ZeroJetID)
// 	require.Error(t, err)
//
// 	err = jetKeeper.AddDropConfirmation(ctx, pn, insolar.ZeroJetID)
// 	require.Error(t, err)
//
// }
