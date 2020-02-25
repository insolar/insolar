// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package writecontroller

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
)

var (
	ErrWriteClosed = errors.New("requested pulse is closed for writing")
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner/writecontroller.Accessor -o ./ -s _mock.go -g
type Accessor interface {
	// Begin requests writing access for pulse number. If requested pulse is closed, ErrWriteClosed will be returned.
	// The caller must call returned "done" function when finished writing.
	Begin(context.Context, insolar.PulseNumber) (done func(), err error)

	// Wait for Open to be called after CloseAndWait
	WaitOpened(ctx context.Context)
}

//go:generate minimock -i github.com/insolar/insolar/logicrunner/writecontroller.Manager -o ./ -s _mock.go -g
type Manager interface {
	// Open marks pulse number as opened for writing. It can be used later by Begin from accessor.
	Open(context.Context, insolar.PulseNumber) error
	// CloseAndWait immediately marks pulse number as closed for writing and blocks until all writes are done.
	CloseAndWait(context.Context, insolar.PulseNumber) error
}

//go:generate minimock -i github.com/insolar/insolar/logicrunner/writecontroller.WriteController -o ./ -s _mock.go -g
type WriteController interface {
	Begin(ctx context.Context, pulse insolar.PulseNumber) (func(), error)
	Open(ctx context.Context, pulse insolar.PulseNumber) error
	WaitOpened(ctx context.Context)
	CloseAndWait(ctx context.Context, pulse insolar.PulseNumber) error
}

type writeController struct {
	lock    sync.RWMutex
	current insolar.PulseNumber
	closed  bool

	wg *sync.WaitGroup
}

func NewWriteController() WriteController {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	return &writeController{
		current: 0,
		closed:  true,
		wg:      wg,
	}
}

func (m *writeController) Begin(ctx context.Context, pulse insolar.PulseNumber) (func(), error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if pulse != m.current {
		return func() {}, ErrWriteClosed
	}
	if m.closed {
		return func() {}, ErrWriteClosed
	}
	m.wg.Add(1)

	return func() { m.wg.Done() }, nil
}

func (m *writeController) Open(ctx context.Context, pulse insolar.PulseNumber) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if pulse < m.current {
		return fmt.Errorf("can't open past pulse for writing: %v", pulse)
	}
	if pulse == m.current {
		return fmt.Errorf("requested pulse already opened for writing: %v", pulse)
	}

	if m.wg != nil {
		// we signaling that we're opened, only if it's not our first iteration
		m.wg.Done()
	}

	m.wg = &sync.WaitGroup{}
	m.current = pulse
	m.closed = false

	return nil
}

func (m *writeController) WaitOpened(ctx context.Context) {
	m.lock.RLock()

	if !m.closed {
		m.lock.RUnlock()
		return
	}

	wg := m.wg

	m.lock.RUnlock()

	// we won't have race condition here, since every new Open we'll have new WaitGroup
	// we're assured that we have old WaitGroup, when we copy pointer to 'wg' var
	wg.Wait()
}

func (m *writeController) CloseAndWait(ctx context.Context, pulse insolar.PulseNumber) error {
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

	// we signaling that we're closed
	m.wg.Add(1)

	return nil
}
