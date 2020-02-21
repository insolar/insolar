// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulse

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/insolar/insolar/ledger/object"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BadgerDefaultOptions(dir string) badger.Options {
	ops := badger.DefaultOptions(dir)
	ops.CompactL0OnClose = false
	ops.SyncWrites = false

	return ops
}

func TestBadgerPulseKey(t *testing.T) {
	t.Parallel()

	expectedKey := pulseKey(insolar.GenesisPulse.PulseNumber)

	rawID := expectedKey.ID()

	actualKey := newPulseKey(rawID)
	require.Equal(t, expectedKey, actualKey)
}

func TestBadgerDropStorageDB_TruncateHead_NoSuchPulse(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	dbMock, err := store.NewBadgerDB(ops)
	defer dbMock.Stop(ctx)
	require.NoError(t, err)

	txManager, err := object.NewBadgerTxManager(dbMock.Backend())
	require.NoError(t, err)
	pulseStore := NewBadgerDB(dbMock, txManager)

	err = pulseStore.TruncateHead(ctx, 77)
	require.NoError(t, err)
}

func TestBadgerDBStore_TruncateHead(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	dbMock, err := store.NewBadgerDB(ops)
	defer dbMock.Stop(ctx)
	require.NoError(t, err)

	txManager, err := object.NewBadgerTxManager(dbMock.Backend())
	require.NoError(t, err)
	dbStore := NewBadgerDB(dbMock, txManager)

	numElements := 10

	startPulseNumber := insolar.GenesisPulse.PulseNumber
	for i := 0; i < numElements; i++ {
		pn := startPulseNumber + insolar.PulseNumber(i)
		pulse := *pulsar.NewPulse(0, pn, &entropygenerator.StandardEntropyGenerator{})
		err := dbStore.Append(ctx, pulse)
		require.NoError(t, err)
	}

	for i := 0; i < numElements; i++ {
		_, err := dbStore.ForPulseNumber(ctx, startPulseNumber+insolar.PulseNumber(i))
		require.NoError(t, err)
	}

	numLeftElements := numElements / 2
	err = dbStore.TruncateHead(ctx, startPulseNumber+insolar.PulseNumber(numLeftElements))
	require.NoError(t, err)

	for i := 0; i < numLeftElements; i++ {
		_, err := dbStore.ForPulseNumber(ctx, startPulseNumber+insolar.PulseNumber(i))
		require.NoError(t, err)
	}

	for i := numElements - 1; i >= numLeftElements; i-- {
		_, err := dbStore.ForPulseNumber(ctx, startPulseNumber+insolar.PulseNumber(i))
		require.EqualError(t, err, ErrNotFound.Error())
	}
}
