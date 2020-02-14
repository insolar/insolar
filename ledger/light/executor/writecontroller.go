// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"context"
	"fmt"
	"sync"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.WriteAccessor -o ./ -s _mock.go -g

type WriteAccessor interface {
	// Begin requests writing access for pulse number. If requested pulse is closed, ErrWriteClosed will be returned.
	// The caller must call returned "done" function when finished writing.
	Begin(context.Context, insolar.PulseNumber) (done func(), err error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.WriteManager -o ./ -s _mock.go -g

type WriteManager interface {
	// Open marks pulse number as opened for writing. It can be used later by Begin from accessor.
	Open(context.Context, insolar.PulseNumber) error
	// CloseAndWait immediately marks pulse number as closed for writing and blocks until all writes are done.
	CloseAndWait(context.Context, insolar.PulseNumber) error
}

type WriteController struct {
	lock    sync.RWMutex
	current insolar.PulseNumber
	closed  bool

	wg sync.WaitGroup
}

func NewWriteController() *WriteController {
	return &WriteController{
		current: 0,
		closed:  true,
	}
}

func (m *WriteController) Begin(ctx context.Context, pulse insolar.PulseNumber) (func(), error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if pulse != m.current {
		return nil, ErrWriteClosed
	}
	if m.closed {
		return nil, ErrWriteClosed
	}
	m.wg.Add(1)

	return func() { m.wg.Done() }, nil
}

func (m *WriteController) Open(ctx context.Context, pulse insolar.PulseNumber) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if pulse < m.current {
		return fmt.Errorf("can't open past pulse for writing: %v", pulse)
	}
	if pulse == m.current {
		return fmt.Errorf("requested pulse already opened for writing: %v", pulse)
	}

	m.wg = sync.WaitGroup{}
	m.current = pulse
	m.closed = false

	return nil
}

func (m *WriteController) CloseAndWait(ctx context.Context, pulse insolar.PulseNumber) error {
	m.lock.Lock()

	if m.current == 0 {
		m.lock.Unlock()
		return nil
	}

	if pulse != m.current {
		m.lock.Unlock()
		return fmt.Errorf("wrong pulse for closing: opened - %v, requested = %v", m.current, pulse)
	}

	if m.closed {
		m.lock.Unlock()
		return fmt.Errorf("requested pulse already closed for writing: %v", pulse)
	}

	m.closed = true
	m.lock.Unlock()

	m.wg.Wait()

	return nil
}
