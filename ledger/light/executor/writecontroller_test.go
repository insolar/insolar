// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/stretchr/testify/require"
)

func TestWriteController_Open(t *testing.T) {
	t.Parallel()

	t.Run("open for correct pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := executor.NewWriteController()
		err := m.Open(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("multiple open for same pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := executor.NewWriteController()
		err := m.Open(ctx, 1)
		require.NoError(t, err)

		err = m.Open(ctx, 1)
		require.Error(t, err)
	})

	t.Run("try to open previous pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := executor.NewWriteController()
		err := m.Open(ctx, 2)
		require.NoError(t, err)

		err = m.Open(ctx, 1)
		require.Error(t, err)
	})
}

func TestWriteController_CloseAndWait(t *testing.T) {
	t.Parallel()

	t.Run("close correct pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := executor.NewWriteController()
		_ = m.Open(ctx, 1)
		err := m.CloseAndWait(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("multiple close for same pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := executor.NewWriteController()
		_ = m.Open(ctx, 1)
		err := m.CloseAndWait(ctx, 1)
		require.NoError(t, err)

		err = m.CloseAndWait(ctx, 1)
		require.Error(t, err)
	})

	t.Run("try to close incorrect pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := executor.NewWriteController()
		err := m.Open(ctx, 2)
		require.NoError(t, err)

		err = m.CloseAndWait(ctx, 1)
		require.Error(t, err)

		err = m.CloseAndWait(ctx, 3)
		require.Error(t, err)
	})
}

func TestWriteController_Begin(t *testing.T) {
	t.Parallel()

	t.Run("begin for not-opened pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := executor.NewWriteController()
		_, err := m.Begin(ctx, 1)
		require.Error(t, err)
	})

	t.Run("begin for closed pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := executor.NewWriteController()
		err := m.Open(ctx, 1)
		require.NoError(t, err)
		err = m.CloseAndWait(ctx, 1)
		require.NoError(t, err)

		_, err = m.Begin(ctx, 1)
		require.Error(t, err)
	})

	t.Run("begin for correct pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := executor.NewWriteController()
		err := m.Open(ctx, 1)
		require.NoError(t, err)

		for i := 0; i < 1000; i++ {
			done, _ := m.Begin(ctx, 1)
			go func() {
				time.Sleep((time.Duration)(rand.Int31n(100)) * time.Millisecond)
				done()
			}()
		}
		err = m.CloseAndWait(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("begin while waiting pulse closing", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := executor.NewWriteController()
		err := m.Open(ctx, 1)
		require.NoError(t, err)

		done, _ := m.Begin(ctx, 1)
		started := make(chan struct{})

		go func() {
			close(started)
			err = m.CloseAndWait(ctx, 1)
			require.NoError(t, err)
		}()
		<-started
		time.Sleep(time.Millisecond * 100)

		_, err = m.Begin(ctx, 1)
		require.Error(t, err)

		done()
	})
}
