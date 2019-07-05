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

package replica

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
)

func TestNewJetKeeper(t *testing.T) {
	db := store.NewMemoryMockDB()
	jets := jet.NewDBStore(db)
	pulses := pulse.NewDB(db)
	jetKeeper := NewJetKeeper(jets, db, pulses)
	require.NotNil(t, jetKeeper)
}

func TestDbJetKeeper_Add(t *testing.T) {
	ctx := inslogger.TestContext(t)

	db := store.NewMemoryMockDB()
	jets := jet.NewDBStore(db)
	pulses := pulse.NewDB(db)
	jetKeeper := NewJetKeeper(jets, db, pulses)

	var (
		err   error
		pulse insolar.PulseNumber
		jet   insolar.JetID
	)
	f := fuzz.New()
	f.Fuzz(&pulse)
	f.Fuzz(&jet)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: insolar.GenesisPulse.PulseNumber})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: pulse - 10})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: pulse})
	require.NoError(t, err)

	err = jetKeeper.Add(ctx, pulse, jet)
	require.NoError(t, err)
}

func TestDbJetKeeper_TopSyncPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)

	db := store.NewMemoryMockDB()
	jets := jet.NewDBStore(db)
	pulses := pulse.NewDB(db)
	jetKeeper := NewJetKeeper(jets, db, pulses)

	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())

	var (
		err          error
		genesisPulse = insolar.GenesisPulse.PulseNumber
		emptyPulse   = genesisPulse + insolar.PulseNumber(10)
		pulse        = genesisPulse + insolar.PulseNumber(20)
		futurePulse  = genesisPulse + insolar.PulseNumber(30)
		jet          = insolar.ZeroJetID
	)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: genesisPulse})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: emptyPulse})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: pulse})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: futurePulse})
	require.NoError(t, err)

	err = jets.Update(ctx, pulse, true, jet)
	require.NoError(t, err)
	err = jetKeeper.Add(ctx, pulse, jet)
	require.NoError(t, err)

	require.Equal(t, pulse, jetKeeper.TopSyncPulse())

	err = jets.Clone(ctx, pulse, futurePulse)
	require.NoError(t, err)
	left, right, err := jets.Split(ctx, futurePulse, jet)
	require.NoError(t, err)
	err = jetKeeper.Add(ctx, futurePulse, left)
	require.NoError(t, err)
	err = jetKeeper.Add(ctx, futurePulse, right)
	require.NoError(t, err)

	require.Equal(t, futurePulse, jetKeeper.TopSyncPulse())
}

func TestDbJetKeeper_OvertakePulse(t *testing.T) {
	ctx := inslogger.TestContext(t)

	db := store.NewMemoryMockDB()
	jets := jet.NewDBStore(db)
	pulses := pulse.NewDB(db)
	jetKeeper := NewJetKeeper(jets, db, pulses)

	var (
		err          error
		genesisPulse = insolar.GenesisPulse.PulseNumber
		emptyPulse   = genesisPulse + insolar.PulseNumber(10)
		P1           = genesisPulse + insolar.PulseNumber(20)
		P2           = genesisPulse + insolar.PulseNumber(30)
		P3           = genesisPulse + insolar.PulseNumber(40)
		jet          = insolar.ZeroJetID
	)

	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: genesisPulse})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: emptyPulse})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: P1})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: P2})
	require.NoError(t, err)
	err = pulses.Append(ctx, insolar.Pulse{PulseNumber: P3})
	require.NoError(t, err)

	// genesis
	require.Equal(t, genesisPulse, jetKeeper.TopSyncPulse())

	// P1 (normal flow)
	err = jets.Update(ctx, P1, true, jet)
	require.NoError(t, err)
	err = jetKeeper.Add(ctx, P1, jet)
	require.NoError(t, err)

	require.Equal(t, P1, jetKeeper.TopSyncPulse())

	// P1 try to overtake P2
	err = jets.Clone(ctx, P1, P2)
	require.NoError(t, err)
	err = jets.Clone(ctx, P2, P3)
	require.NoError(t, err)
	err = jetKeeper.Add(ctx, P3, jet)
	require.NoError(t, err)

	require.Equal(t, P1, jetKeeper.TopSyncPulse())

	// P3 catch up
	err = jetKeeper.Add(ctx, P2, jet)
	require.NoError(t, err)

	require.Equal(t, P3, jetKeeper.TopSyncPulse())
}
