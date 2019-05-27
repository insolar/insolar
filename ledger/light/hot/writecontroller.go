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

package hot

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/hot.WriteAccessor -o ./ -s _mock.go

type WriteAccessor interface {
	// Begin requests writing access for pulse number. If requested pulse is closed, ErrClosed will be returned.
	// The caller must call returned "done" function when finished writing.
	Begin(context.Context, insolar.PulseNumber) (done func(), err error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/light/hot.WriteManager -o ./ -s _mock.go

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

	logger := inslogger.FromContext(ctx)

	if pulse != m.current {
		logger.Errorf("can't begin writing for wrong pulse: opened - %v, requested - %v", m.current, pulse)
		return func() {}, ErrWriteClosed
	}
	if m.closed {
		logger.Error("requested pulse is closed for writing")
		return func() {}, ErrWriteClosed
	}
	m.wg.Add(1)

	return func() { m.wg.Done() }, nil
}

func (m *WriteController) Open(ctx context.Context, pulse insolar.PulseNumber) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	logger := inslogger.FromContext(ctx)

	if pulse < m.current {
		logger.Error("can't open past pulse for writing: ", pulse)
		return ErrWriteClosed
	}
	if pulse == m.current {
		logger.Warn("requested pulse already opened for writing: ", pulse)
		return nil
	}

	var wg sync.WaitGroup
	m.wg = wg
	m.current = pulse
	m.closed = false

	return nil
}

func (m *WriteController) CloseAndWait(ctx context.Context, pulse insolar.PulseNumber) error {
	m.lock.Lock()

	logger := inslogger.FromContext(ctx)

	if pulse != m.current {
		logger.Errorf("wrong pulse for closing: opened - %v, requested = %v", m.current, pulse)
		m.lock.Unlock()
		return ErrWriteClosed
	}

	if pulse == m.current && m.closed {
		logger.Error("requested pulse already closed for writing: ", pulse)
		m.lock.Unlock()
		return ErrWriteClosed
	}

	m.closed = true
	m.lock.Unlock()

	m.wg.Wait()

	return nil
}
