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
	"io"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/flow"
)

func TestNewThread(t *testing.T) {
	t.Parallel()
	msg := &message.Message{}
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

func TestThread_Procedure_CancelledBefore(t *testing.T) {
	t.Parallel()
	cancel := make(chan struct{})
	thread := Thread{
		cancel: cancel,
	}
	close(cancel)
	proc := flow.NewProcedureMock(t)
	err := thread.Procedure(context.Background(), proc, true)
	require.Error(t, err)
	require.Equal(t, err, flow.ErrCancelled)
}

type loggerMock struct {
	fatalChecker func()
}

func (l *loggerMock) WithLevel(string) (insolar.Logger, error) {
	panic("implement me")
}

func (l *loggerMock) WithLevelNumber(level insolar.LogLevel) (insolar.Logger, error) {
	panic("implement me")
}

func (l *loggerMock) WithFormat(format insolar.LogFormat) (insolar.Logger, error) {
	panic("implement me")
}

func (l *loggerMock) WithCaller(flag bool) insolar.Logger {
	panic("implement me")
}

func (l *loggerMock) WithSkipFrameCount(delta int) insolar.Logger {
	panic("implement me")
}

func (l *loggerMock) WithFuncName(flag bool) insolar.Logger {
	panic("implement me")
}

func (l *loggerMock) Debug(...interface{}) {
	panic("implement me")
}

func (l *loggerMock) Debugf(string, ...interface{}) {
	panic("implement me")
}

func (l *loggerMock) Info(...interface{}) {
	panic("implement me")
}

func (l *loggerMock) Infof(string, ...interface{}) {
	panic("implement me")
}

func (l *loggerMock) Warn(...interface{}) {
	panic("implement me")
}

func (l *loggerMock) Warnf(string, ...interface{}) {
	panic("implement me")
}

func (l *loggerMock) Error(...interface{}) {
	panic("implement me")
}

func (l *loggerMock) Errorf(string, ...interface{}) {
	panic("implement me")
}

func (l *loggerMock) Fatal(...interface{}) {
	l.fatalChecker()
}

func (l *loggerMock) Fatalf(string, ...interface{}) {
	panic("implement me")
}

func (l *loggerMock) Panic(...interface{}) {
	panic("implement me")
}

func (l *loggerMock) Panicf(string, ...interface{}) {
	panic("implement me")
}

func (l *loggerMock) WithOutput(w io.Writer) insolar.Logger {
	panic("implement me")
}

func (l *loggerMock) WithFields(map[string]interface{}) insolar.Logger {
	panic("implement me")
}

func (l *loggerMock) WithField(string, interface{}) insolar.Logger {
	panic("implement me")
}

func (l *loggerMock) Is(level insolar.LogLevel) bool {
	panic("implement me")
}

func TestThread_Procedure_LoggerFailProcedureError(t *testing.T) {
	t.Parallel()
	thread := Thread{}
	isCalled := false

	logger := &loggerMock{
		fatalChecker: func() {
			isCalled = true
		},
	}
	ctx := context.TODO()
	ctx = inslogger.SetLogger(ctx, logger)
	_ = thread.Procedure(ctx, nil, true)

	require.True(t, isCalled)
}

func TestThread_Procedure_CancelledWhenProcedureWorks(t *testing.T) {
	t.Parallel()
	cancel := make(chan struct{})
	finish := make(chan struct{})
	thread := Thread{
		cancel:     cancel,
		procedures: map[flow.Procedure]*result{},
	}
	pm := flow.NewProcedureMock(t)
	pm.ProceedMock.Set(func(ctx context.Context) error {
		close(cancel)
		<-finish
		return nil
	})
	err := thread.Procedure(context.Background(), pm, true)
	require.Error(t, err)
	require.Equal(t, flow.ErrCancelled, err)
	close(finish)
}

func TestThread_Procedure_NotCancelled(t *testing.T) {
	t.Parallel()
	cancel := make(chan struct{})
	thread := Thread{
		cancel:     cancel,
		procedures: map[flow.Procedure]*result{},
	}
	close(cancel)
	proc := flow.NewProcedureMock(t)
	proc.ProceedMock.Return(nil)
	err := thread.Procedure(context.Background(), proc, false)
	require.NoError(t, err)
}

func TestThread_Procedure_ProceedReturnsError(t *testing.T) {
	t.Parallel()
	thread := Thread{
		procedures: map[flow.Procedure]*result{},
	}
	pm := flow.NewProcedureMock(t)
	pm.ProceedMock.Set(func(ctx context.Context) error {
		return errors.New("proceed test error")
	})
	err := thread.Procedure(context.Background(), pm, true)
	require.Error(t, err)
	require.EqualError(t, err, "proceed test error")
}

func TestThread_Procedure(t *testing.T) {
	t.Parallel()
	thread := Thread{
		procedures: map[flow.Procedure]*result{},
	}
	pm := flow.NewProcedureMock(t)
	pm.ProceedMock.Set(func(ctx context.Context) error {
		return nil
	})
	err := thread.Procedure(context.Background(), pm, true)
	require.NoError(t, err)
}

func TestThread_Procedure_Reattach(t *testing.T) {
	t.Parallel()
	thread := Thread{
		procedures: map[flow.Procedure]*result{},
	}
	procErr := errors.New("test error")
	pm := flow.NewProcedureMock(t)
	pm.ProceedMock.Set(func(ctx context.Context) error {
		return procErr
	})
	err := thread.Procedure(context.Background(), pm, true)
	require.Equal(t, procErr, err)
	done := make(chan struct{})
	go func() {
		err = thread.Procedure(context.Background(), pm, true)
		require.Equal(t, procErr, err)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		assert.Fail(t, "reattach deadlock")
	}
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
		canBegin: make(chan struct{}),
	}
	thread := Thread{
		controller: controller,
		canBegin:   controller.canBegin,
	}
	close(controller.canBegin)

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
		canBegin: make(chan struct{}),
	}
	thread := Thread{
		controller: controller,
		canBegin:   controller.canBegin,
	}
	close(controller.canBegin)

	handle := func(ctx context.Context, f flow.Flow) error {
		require.NotEqual(t, &thread, f)
		return nil
	}
	err := thread.Migrate(context.Background(), handle)
	require.NoError(t, err)
	require.True(t, thread.migrated)
}

func TestThread_Continue(t *testing.T) {
	controllerBegin := make(chan struct{})
	threadBegin := make(chan struct{})
	thread := Thread{
		controller: &Controller{
			canBegin: controllerBegin,
		},
		canBegin: threadBegin,
	}
	close(threadBegin)
	thread.Continue(context.Background())
	var expected <-chan struct{} = controllerBegin
	require.Equal(t, expected, thread.canBegin)
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
