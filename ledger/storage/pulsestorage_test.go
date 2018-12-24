/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package storage_test

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestNewPulseStorage(t *testing.T) {
	t.Parallel()

	// Arrange
	testDb := storage.DB{}

	// Act
	pStorage := storage.NewPulseStorage(&testDb)

	// Assert
	require.NotNil(t, pStorage)
}

func TestLockUnlock(t *testing.T) {
	t.Parallel()

	// Arrange
	testDb := storage.DB{}
	pStorage := storage.NewPulseStorage(&testDb)

	// Act
	pStorage.Lock()
	pStorage.Unlock()
}

func TestCurrent_OneThread(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := inslogger.TestContext(t)
	testDb, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	err := testDb.AddPulse(ctx, *core.GenesisPulse)
	require.NoError(t, err)
	pStorage := storage.NewPulseStorage(testDb)

	// Act
	pulse, err := pStorage.Current(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, core.GenesisPulse, pulse)
}

func TestCurrent_ThreeThreads(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx := inslogger.TestContext(t)
	testDb, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	pStorage := storage.NewPulseStorage(testDb)

	// Act
	var g errgroup.Group
	g.Go(func() error {
		pStorage.Lock()
		defer pStorage.Unlock()
		err := testDb.AddPulse(ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 123})
		return err
	})
	g.Go(func() error {
		_, err := pStorage.Current(ctx)
		return err
	})
	g.Go(func() error {
		_, err := pStorage.Current(ctx)
		return err
	})
	err := g.Wait()
	require.NoError(t, err)
	pulse, err := pStorage.Current(ctx)

	// Assert
	require.NoError(t, err)
	require.Equal(t, core.GenesisPulse.PulseNumber+123, pulse.PulseNumber)
}
