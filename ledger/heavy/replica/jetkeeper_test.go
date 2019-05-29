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
	"github.com/insolar/insolar/internal/ledger/store"
)

func TestNewJetKeeper(t *testing.T) {
	db := store.NewMemoryMockDB()
	jetKeeper := NewJetKeeper(db)
	require.NotNil(t, jetKeeper)
}

func TestDbJetKeeper_Add(t *testing.T) {
	db := store.NewMemoryMockDB()
	jetKeeper := NewJetKeeper(db)

	var (
		pulse insolar.PulseNumber
		jet   insolar.JetID
	)
	f := fuzz.New()
	f.Fuzz(&pulse)
	f.Fuzz(&jet)
	err := jetKeeper.Add(pulse, jet)
	require.NoError(t, err)
}

func TestDbJetKeeper_All(t *testing.T) {
	db := store.NewMemoryMockDB()
	jetKeeper := NewJetKeeper(db)
	const (
		pulse    = 10
		jetCount = 100
	)
	jets, err := jetKeeper.All(pulse)
	require.NoError(t, err)
	require.Empty(t, jets)

	var (
		jet insolar.JetID
	)
	f := fuzz.New()
	for i := 0; i < jetCount; i++ {
		f.Fuzz(&jet)
		err = jetKeeper.Add(pulse, jet)
		require.NoError(t, err)
	}

	jets, err = jetKeeper.All(pulse)
	require.NoError(t, err)
	require.Len(t, jets, jetCount)
}
