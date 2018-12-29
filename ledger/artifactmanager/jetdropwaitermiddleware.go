package artifactmanager

import (
	"context"
	"sync"
	"time"

	"github.com/insolar/insolar/core"
)

type jetDropTimeoutProvider struct {
	waiters          map[core.RecordID]*jetDropTimeout
	waitersInitLocks map[core.RecordID]*sync.RWMutex

	waitersLock          sync.RWMutex
	waitersInitLocksLock sync.Mutex
}

func (p *jetDropTimeoutProvider) getLock(jetID core.RecordID) *sync.RWMutex {
	p.waitersInitLocksLock.Lock()
	defer p.waitersInitLocksLock.Unlock()

	if _, ok := p.waitersInitLocks[jetID]; !ok {
		p.waitersInitLocks[jetID] = &sync.RWMutex{}
	}

	return p.waitersInitLocks[jetID]
}

func (p *jetDropTimeoutProvider) getWaiter(jetID core.RecordID) *jetDropTimeout {
	p.waitersLock.RLock()
	defer p.waitersLock.RUnlock()

	return p.waiters[jetID]
}

type jetDropTimeout struct {
	lastJdPulseLock sync.RWMutex
	lastJdPulse     core.PulseNumber

	jetDropLocker chan struct{}
	timeoutLocker chan struct{}

	isTimeoutRunLock sync.Mutex
	isTimeoutRun     bool
}

func (jdw *jetDropTimeout) getLastJdPulse() core.PulseNumber {
	jdw.lastJdPulseLock.RLock()
	defer jdw.lastJdPulseLock.RUnlock()

	return jdw.lastJdPulse
}

func (jdw *jetDropTimeout) setLastJdPulse(pn core.PulseNumber) {
	jdw.lastJdPulseLock.Lock()
	defer jdw.lastJdPulseLock.Unlock()

	jdw.lastJdPulse = pn
}

func (m *middleware) waitForDrop(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		jetID := jetFromContext(ctx)
		lock := m.jetDropTimeoutProvider.getLock(jetID)
		waiter := m.jetDropTimeoutProvider.getWaiter(jetID)

		lock.RLock()
		if waiter == nil {
			lock.RUnlock()
			return handler(ctx, parcel)
		}
		lock.RUnlock()

		if waiter.getLastJdPulse() != parcel.Pulse() {
			waiter.runDropWaitingTimeout()

			select {
			case <-waiter.jetDropLocker:
			case <-waiter.timeoutLocker:
			}

			waiter.isTimeoutRunLock.Lock()
			waiter.isTimeoutRun = false
			waiter.isTimeoutRunLock.Unlock()
		}

		return handler(ctx, parcel)
	}
}

func (w *jetDropTimeout) runDropWaitingTimeout() {
	w.isTimeoutRunLock.Lock()
	defer w.isTimeoutRunLock.Unlock()

	if w.isTimeoutRun {
		return
	}

	w.isTimeoutRun = true
	go func() {
		time.Sleep(2 * time.Second)

		close(w.timeoutLocker)
		w.timeoutLocker = make(chan struct{})

		w.isTimeoutRunLock.Lock()
		w.isTimeoutRun = false
		w.isTimeoutRunLock.Unlock()
	}()
}

func (m *middleware) unlockDropWaiters(handler core.MessageHandler) core.MessageHandler {
	return func(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
		jetID := jetFromContext(ctx)
		lock := m.jetDropTimeoutProvider.getLock(jetID)
		waiter := m.jetDropTimeoutProvider.getWaiter(jetID)

		lock.Lock()
		defer lock.Unlock()

		if waiter == nil {
			waiter = &jetDropTimeout{
				jetDropLocker: make(chan struct{}),
				timeoutLocker: make(chan struct{}),
			}
			m.jetDropTimeoutProvider.waiters[jetID] = waiter
		}
		resp, err := handler(ctx, parcel)

		waiter.setLastJdPulse(parcel.Pulse())
		close(waiter.jetDropLocker)
		waiter.jetDropLocker = make(chan struct{})

		return resp, err
	}
}
