// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
)

func TestDBRollback_TruncateReturnError(t *testing.T) {
	jetKeeper := NewJetKeeperMock(t)
	testPulse := insolar.GenesisPulse.PulseNumber + 1
	jetKeeper.TopSyncPulseMock.Set(func() (r insolar.PulseNumber) {
		return testPulse
	})

	calculator := pulse.NewCalculatorMock(t)
	calculator.ForwardsMock.Set(func(p context.Context, p1 insolar.PulseNumber, p2 int) (r insolar.Pulse, r1 error) {
		return insolar.Pulse{PulseNumber: p1 + 1}, nil
	})

	testError := errors.New("Hello")
	drops := NewHeadTruncaterMock(t)
	drops.TruncateHeadMock.Return(testError)
	rollback := NewDBRollback(jetKeeper, drops)
	err := rollback.Start(context.Background())
	require.Contains(t, err.Error(), testError.Error(), err)
}

func TestDBRollback_HappyPath_Badger(t *testing.T) {
	jetKeeper := NewJetKeeperMock(t)
	testPulse := insolar.GenesisPulse.PulseNumber + 1
	jetKeeper.TopSyncPulseMock.Set(func() (r insolar.PulseNumber) {
		return testPulse
	})

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	badger, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer badger.Stop(context.Background())

	db := store.NewDBMock(t)
	hits := make(map[store.Scope]int)
	db.SetMock.Return(nil)

	db.GetMock.Return([]byte{}, nil)

	db.DeleteMock.Return(nil)
	iterNum := 0
	db.NewIteratorMock.Set(func(p store.Key, p1 bool) (r store.Iterator) {
		num, _ := hits[p.Scope()]
		hits[p.Scope()] = num + 1

		iterMock := store.NewIteratorMock(t)
		iterMock.NextMock.Set(func() (r bool) {
			iterNum++
			return iterNum%2 != 0
		})

		// IndexDB key
		var key []byte
		key = append(key, testPulse.Bytes()...)
		key = append(key, p.ID()...)

		iterMock.KeyMock.Return(key)
		iterMock.CloseMock.Return()
		iterMock.ValueMock.Return([]byte{}, nil)
		return iterMock
	})

	drops := drop.NewBadgerDB(db)

	records := NewHeadTruncaterMock(t)
	records.TruncateHeadMock.Set(func(ctx context.Context, from insolar.PulseNumber) (err error) {
		hits[store.ScopeRecord] = 1
		return nil
	})

	indexes := object.NewBadgerIndexDB(db, nil)

	jets := jet.NewBadgerDBStore(db)
	txManager, err := object.NewBadgerTxManager(badger.Backend())
	require.NoError(t, err)
	pulses := pulse.NewBadgerDB(badger, txManager)

	nodes := NewHeadTruncaterMock(t)
	nodes.TruncateHeadMock.Set(func(ctx context.Context, from insolar.PulseNumber) (err error) {
		hits[store.ScopeNodeHistory] = 1
		return nil
	})

	calculator := pulse.NewCalculatorMock(t)
	calculator.ForwardsMock.Set(func(p context.Context, p1 insolar.PulseNumber, p2 int) (r insolar.Pulse, r1 error) {
		return insolar.Pulse{PulseNumber: p1 + 1}, nil
	})

	rollback := NewDBRollback(jetKeeper, drops, records, indexes, jets, pulses, nodes)
	err = rollback.Start(context.Background())
	expectedScopes := []struct {
		scope   store.Scope
		numHits int
	}{
		{store.ScopeJetDrop, 1},
		{store.ScopeRecord, 1},
		{store.ScopeIndex, 2},
		{store.ScopeJetTree, 1},
		{store.ScopeNodeHistory, 1}}
	for _, s := range expectedScopes {
		actualNum, ok := hits[s.scope]
		require.True(t, ok, "Scope: ", s.scope)
		require.Equal(t, s.numHits, actualNum, "Scope: ", s.scope)
	}
	require.Len(t, hits, 5) // drops, record, jets, indexes, nodes
	require.NoError(t, err)
}
