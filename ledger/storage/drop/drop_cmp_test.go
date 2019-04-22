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

package drop

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type jetPulse struct {
	jetID insolar.JetID
	pn    insolar.PulseNumber
}

func TestDropStorageMemory(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ms := NewStorageMemory()

	var drops []Drop
	genInputs := map[jetPulse]struct{}{}
	f := fuzz.New().Funcs(func(jd *Drop, c fuzz.Continue) {
		pn := gen.PulseNumber()
		jd.Pulse = pn

		jetID := gen.JetID()
		jd.JetID = jetID

		genInputs[jetPulse{jetID: jetID, pn: pn}] = struct{}{}
	}).NumElements(5, 1000)
	f.Fuzz(&drops)

	// Add
	for _, dr := range drops {
		err := ms.Set(ctx, dr)
		require.NoError(t, err)
	}

	// Fetch
	for inp := range genInputs {
		_, err := ms.ForPulse(ctx, inp.jetID, inp.pn)
		require.NoError(t, err)
	}

	// Delete
	for inp := range genInputs {
		ms.DeleteForPN(ctx, inp.pn)
	}

	// Check for deleting
	for inp := range genInputs {
		_, err := ms.ForPulse(ctx, inp.jetID, inp.pn)
		require.Error(t, err, ErrNotFound)
	}
}

func TestDropStorageDB(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ds := NewDB(store.NewMemoryMockDB())

	var drops []Drop
	genInputs := map[jetPulse]struct{}{}
	f := fuzz.New().Funcs(func(jd *Drop, c fuzz.Continue) {
		pn := gen.PulseNumber()
		jd.Pulse = pn

		jetID := gen.JetID()
		jd.JetID = jetID

		genInputs[jetPulse{jetID: jetID, pn: pn}] = struct{}{}
	}).NumElements(5, 1000)
	f.Fuzz(&drops)

	// Add
	for _, dr := range drops {
		err := ds.Set(ctx, dr)
		require.NoError(t, err)
	}

	// Fetch
	for inp := range genInputs {
		_, err := ds.ForPulse(ctx, inp.jetID, inp.pn)
		require.NoError(t, err)
	}
}

func TestDropStorageCompare(t *testing.T) {
	ctx := inslogger.TestContext(t)

	ds := NewDB(store.NewMemoryMockDB())
	ms := NewStorageMemory()

	var drops []Drop
	genInputs := map[jetPulse]struct{}{}
	f := fuzz.New().Funcs(func(jd *Drop, c fuzz.Continue) {
		pn := gen.PulseNumber()
		jd.Pulse = pn

		jetID := gen.JetID()
		jd.JetID = jetID

		genInputs[jetPulse{jetID: jetID, pn: pn}] = struct{}{}
	}).NumElements(5, 1000)
	f.Fuzz(&drops)

	// Add
	for _, dr := range drops {
		err := ds.Set(ctx, dr)
		require.NoError(t, err)
		err = ms.Set(ctx, dr)
		require.NoError(t, err)
	}

	// Fetch
	for inp := range genInputs {
		dbDrop, err := ds.ForPulse(ctx, inp.jetID, inp.pn)
		require.NoError(t, err)

		memDrop, err := ms.ForPulse(ctx, inp.jetID, inp.pn)
		require.NoError(t, err)

		require.Equal(t, dbDrop, memDrop)
	}

}
