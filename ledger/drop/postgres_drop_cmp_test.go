// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build slowtest

package drop

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/pulse"
)

func TestPostgresDropStorageDB(t *testing.T) {
	defer cleanupDatabase()

	ctx := inslogger.TestContext(t)
	db := NewPostgresDB(getPool())

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
		err := db.Set(ctx, dr)
		require.NoError(t, err)
	}

	// Fetch
	for inp := range genInputs {
		_, err := db.ForPulse(ctx, inp.jetID, inp.pn)
		require.NoError(t, err)
	}
}

func TestDropStorageCompare(t *testing.T) {
	defer cleanupDatabase()

	ctx := inslogger.TestContext(t)

	db := NewPostgresDB(getPool())
	ms := NewStorageMemory()

	var drops []Drop
	pn := insolar.PulseNumber(pulse.MinTimePulse)
	genInputs := map[jetPulse]struct{}{}
	f := fuzz.New().Funcs(func(jd *Drop, c fuzz.Continue) {
		jd.Pulse = pn
		jetID := gen.JetID()
		jd.JetID = jetID

		genInputs[jetPulse{jetID: jetID, pn: pn}] = struct{}{}

		pn++
	}).NumElements(5, 1000)
	f.Fuzz(&drops)

	// Add
	for _, dr := range drops {
		err := db.Set(ctx, dr)
		require.NoError(t, err)
		err = ms.Set(ctx, dr)
		require.NoError(t, err)
	}

	// Fetch
	for inp := range genInputs {
		dbDrop, err := db.ForPulse(ctx, inp.jetID, inp.pn)
		require.NoError(t, err)

		memDrop, err := ms.ForPulse(ctx, inp.jetID, inp.pn)
		require.NoError(t, err)

		require.Equal(t, dbDrop, memDrop)
	}

}
