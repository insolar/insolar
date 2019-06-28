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
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
)

func TestNewJetKeeper(t *testing.T) {
	db := store.NewMemoryMockDB()
	jets := jet.NewDBStore(db)
	jetKeeper := NewJetKeeper(jets, db)
	require.NotNil(t, jetKeeper)
}

func TestDbJetKeeper_Add(t *testing.T) {
	ctx := inslogger.TestContext(t)

	db := store.NewMemoryMockDB()
	jets := jet.NewDBStore(db)
	jetKeeper := NewJetKeeper(jets, db)

	var (
		pulse insolar.PulseNumber
		jet   insolar.JetID
	)
	f := fuzz.New()
	f.Fuzz(&pulse)
	f.Fuzz(&jet)
	err := jetKeeper.Add(ctx, pulse, jet)
	require.NoError(t, err)
}

func TestDbJetKeeper_TopSyncPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)

	db := store.NewMemoryMockDB()
	jets := jet.NewDBStore(db)
	jetKeeper := NewJetKeeper(jets, db)

	require.Equal(t, insolar.GenesisPulse.PulseNumber, jetKeeper.TopSyncPulse())

	var (
		pulse       insolar.PulseNumber
		futurePulse insolar.PulseNumber
		jet         insolar.JetID
	)
	pulse = 10
	futurePulse = 20
	jet = insolar.ZeroJetID

	err := jetKeeper.Add(ctx, pulse, jet)
	require.NoError(t, err)

	err = jets.Update(ctx, pulse, false, jet)
	require.NoError(t, err)
	left, right, err := jets.Split(ctx, pulse, jet)
	require.NoError(t, err)

	require.Equal(t, pulse, jetKeeper.TopSyncPulse())

	err = jets.Clone(ctx, pulse, futurePulse)
	require.NoError(t, err)
	err = jetKeeper.Add(ctx, futurePulse, left)
	require.NoError(t, err)
	err = jetKeeper.Add(ctx, futurePulse, right)
	require.NoError(t, err)

	require.Equal(t, futurePulse, jetKeeper.TopSyncPulse())
}
