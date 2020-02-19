// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package drop

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/dgraph-io/badger"
	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/pulse"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func BadgerDefaultOptions(dir string) badger.Options {
	ops := badger.DefaultOptions(dir)
	ops.CompactL0OnClose = false
	ops.SyncWrites = false

	return ops
}

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

func TestBadgerDropStorageDB(t *testing.T) {
	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	ds := NewBadgerDB(db)

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

func TestBadgerDropStorageCompare(t *testing.T) {
	ctx := inslogger.TestContext(t)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	ds := NewBadgerDB(db)
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
