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

package thread

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNewThread(t *testing.T) {
	t.Parallel()
	msg := bus.Message{}
	ch := make(chan struct{})
	controller := &Controller{
		cancel: ch,
	}
	thread := NewThread(msg, controller)
	require.NotNil(t, thread)
	require.Equal(t, controller, thread.controller)
	require.Equal(t, ch, thread.controller.cancel)
	require.NotNil(t, thread.procedures)
	require.Equal(t, msg, thread.message)
}

func TestThread_Handle_CancelledBefore(t *testing.T) {
	t.Parallel()
	cancel := make(chan struct{})
	thread := Thread{
		cancel: cancel,
	}
	close(cancel)
	handle := func(ctx context.Context, f flow.Flow) error {
		return nil
	}
	err := thread.Handle(context.Background(), handle)
	require.Error(t, err)
	require.Equal(t, err, flow.ErrCancelled)
}

func TestThread_Handle_Error(t *testing.T) {
	t.Parallel()
	thread := Thread{}

	handleError := errors.New("test error")
	handle := func(ctx context.Context, f flow.Flow) error {
		return handleError
	}
	err := thread.Handle(context.Background(), handle)
	require.Error(t, err)
	require.Equal(t, err, handleError)
}

func TestThread_Handle_CanceledAfter(t *testing.T) {
	t.Parallel()
	cancel := make(chan struct{})
	thread := Thread{
		cancel: cancel,
	}
	handle := func(ctx context.Context, f flow.Flow) error {
		close(cancel)
		return nil
	}
	err := thread.Handle(context.Background(), handle)
	require.Error(t, err)
	require.Equal(t, err, flow.ErrCancelled)
}

func TestThread_Handle(t *testing.T) {
	t.Parallel()
	thread := Thread{}
	handle := func(ctx context.Context, f flow.Flow) error {
		return nil
	}
	err := thread.Handle(context.Background(), handle)
	require.NoError(t, err)
}

func TestThread_Procedure_CancelledBefore(t *testing.T) {
	t.Parallel()
	cancel := make(chan struct{})
	thread := Thread{
		cancel: cancel,
	}
	close(cancel)
	err := thread.Procedure(context.Background(), nil)
	require.Error(t, err)
	require.Equal(t, err, flow.ErrCancelled)
}

func TestThread_Procedure_NilProcedureError(t *testing.T) {
	t.Parallel()
	thread := Thread{}
	err := thread.Procedure(context.Background(), nil)
	require.EqualError(t, err, "procedure called with nil procedure")
}

func TestThread_Procedure_CancelledWhenProcedureWorks(t *testing.T) {
	t.Parallel()
	cancel := make(chan struct{})
	thread := Thread{
		cancel:     cancel,
		procedures: map[flow.Procedure]chan error{},
	}
	pm := flow.NewProcedureMock(t)
	pm.ProceedFunc = func(ctx context.Context) error {
		close(cancel)
		runtime.Gosched()
		<-cancel
		return nil
	}
	err := thread.Procedure(context.Background(), pm)
	require.Error(t, err)
	require.Equal(t, flow.ErrCancelled, err)
}

func TestThread_Procedure_ProceedReturnsError(t *testing.T) {
	t.Parallel()
	thread := Thread{
		procedures: map[flow.Procedure]chan error{},
	}
	pm := flow.NewProcedureMock(t)
	pm.ProceedFunc = func(ctx context.Context) error {
		return errors.New("proceed test error")
	}
	err := thread.Procedure(context.Background(), pm)
	require.Error(t, err)
	require.EqualError(t, err, "proceed test error")
}

func TestThread_Procedure(t *testing.T) {
	t.Parallel()
	thread := Thread{
		procedures: map[flow.Procedure]chan error{},
	}
	pm := flow.NewProcedureMock(t)
	pm.ProceedFunc = func(ctx context.Context) error {
		return nil
	}
	err := thread.Procedure(context.Background(), pm)
	require.NoError(t, err)
}

func TestThread_Migrate_MigratedError(t *testing.T) {
	t.Parallel()
	thread := Thread{
		migrated: true,
	}
	err := thread.Migrate(context.Background(), nil)
	require.EqualError(t, err, "migrate called on migrated flow")
}

func TestThread_Migrate_HandleReturnsError(t *testing.T) {
	t.Parallel()
	controller := &Controller{
		cancel: make(chan struct{}),
	}
	thread := Thread{
		controller: controller,
		cancel:     controller.cancel,
	}
	close(controller.cancel)

	handle := func(ctx context.Context, f flow.Flow) error {
		require.NotEqual(t, &thread, f)
		return errors.New("test error")
	}
	err := thread.Migrate(context.Background(), handle)
	require.EqualError(t, err, "test error")
}

func TestThread_Migrate(t *testing.T) {
	t.Parallel()
	controller := &Controller{
		cancel: make(chan struct{}),
	}
	thread := Thread{
		controller: controller,
		cancel:     controller.cancel,
	}
	close(controller.cancel)

	handle := func(ctx context.Context, f flow.Flow) error {
		require.NotEqual(t, &thread, f)
		return nil
	}
	err := thread.Migrate(context.Background(), handle)
	require.NoError(t, err)
	require.True(t, thread.migrated)
}

func TestThread_Continue(t *testing.T) {
	controllerCancel := make(chan struct{})
	threadCancel := make(chan struct{})
	thread := Thread{
		controller: &Controller{
			cancel: controllerCancel,
		},
		cancel: threadCancel,
	}
	close(threadCancel)
	thread.Continue(context.Background())
	var expected <-chan struct{} = controllerCancel
	require.Equal(t, expected, thread.cancel)
}

func TestThread_Run_Error(t *testing.T) {
	t.Parallel()
	thread := Thread{}
	handle := func(ctx context.Context, f flow.Flow) error {
		require.Equal(t, &thread, f)
		return errors.New("test error")
	}
	err := thread.Run(context.Background(), handle)
	require.EqualError(t, err, "test error")
}

func TestThread_Run(t *testing.T) {
	t.Parallel()
	thread := Thread{}
	handle := func(ctx context.Context, f flow.Flow) error {
		require.Equal(t, &thread, f)
		return nil
	}
	err := thread.Run(context.Background(), handle)
	require.NoError(t, err)
}

func TestThread_procedure_AlreadyExists(t *testing.T) {
	t.Parallel()
	thread := Thread{
		procedures: map[flow.Procedure]chan error{},
	}
	procedure := flow.NewProcedureMock(t)
	done := make(chan error, 1)
	thread.procedures[procedure] = done

	ch := thread.procedure(context.Background(), procedure)
	var expected <-chan error = done
	require.Equal(t, expected, ch)
}

func TestThread_procedure(t *testing.T) {
	t.Parallel()
	thread := Thread{
		procedures: map[flow.Procedure]chan error{},
	}
	procedure := flow.NewProcedureMock(t)
	procedure.ProceedFunc = func(ctx context.Context) error {
		return errors.New("test error")
	}

	ch := thread.procedure(context.Background(), procedure)
	require.NotNil(t, ch)
	timer := time.NewTimer(10 * time.Millisecond)
	select {
	case result := <-ch:
		require.EqualError(t, result, "test error")
	case <-timer.C:
		t.Fatal("timeout")
	}
}

func TestThread_canceled_Canceled(t *testing.T) {
	t.Parallel()
	cancel := make(chan struct{})
	thread := Thread{
		cancel: cancel,
	}
	close(cancel)
	result := thread.cancelled()
	require.True(t, result)
}

func TestThread_canceled(t *testing.T) {
	t.Parallel()
	cancel := make(chan struct{})
	thread := Thread{
		cancel: cancel,
	}
	result := thread.cancelled()
	require.False(t, result)
}
