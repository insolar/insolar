package artifactmanager

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// HotDataWaiter provides waiting system for a specific jet
// We tend to think, that it will be used for waiting hot-data in handler
// Also, because of the some jet pitfalls, we need to have an instrument
// to handler edge-cases from pulse manager
//go:generate minimock -i github.com/insolar/insolar/ledger/storage.HotDataWaiter -o ./ -s _mock.go
type HotDataWaiter interface {
	Wait(ctx context.Context, jetID core.RecordID) error
	Unlock(ctx context.Context, jetID core.RecordID)
	ThrowTimeout(ctx context.Context)
}

// HotDataWaiterConcrete is an implementation of HotDataWaiter
type HotDataWaiterConcrete struct {
	waitersMapLock sync.Mutex
	waiters        map[core.RecordID]*waiter
}

// NewHotDataWaiterConcrete is a constructor
func NewHotDataWaiterConcrete() *HotDataWaiterConcrete {
	return &HotDataWaiterConcrete{}
}

type waiter struct {
	hotDataChannel chan struct{}
	timeoutChannel chan struct{}
}

func (hdw *HotDataWaiterConcrete) getWaiter(ctx context.Context, jetID core.RecordID) *waiter {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("[getWaiter] jetID - %v", jetID.DebugString())

	hdw.waitersMapLock.Lock()
	defer hdw.waitersMapLock.Unlock()

	if _, ok := hdw.waiters[jetID]; !ok {
		logger.Debugf("[getWaiter] create new  - %v", jetID.DebugString())
		hdw.waiters[jetID] = &waiter{
			hotDataChannel: make(chan struct{}),
			timeoutChannel: make(chan struct{}),
		}
	}

	return hdw.waiters[jetID]
}

// Wait waits for the raising one of two channels.
// If hotDataChannel or timeoutChannel was raised, the method returns error
// Either nil or ErrHotDataTimeout
func (hdw *HotDataWaiterConcrete) Wait(ctx context.Context, jetID core.RecordID) error {
	logger := inslogger.FromContext(ctx)
	waiter := hdw.getWaiter(ctx, jetID)

	logger.Debugf("[Wait] before pause of request with jet - %v", jetID.DebugString())
	select {
	case <-waiter.hotDataChannel:
		logger.Debugf("[Wait] hotDataChannel's events was raised")
		return nil
	case <-waiter.timeoutChannel:
		logger.Errorf("[Wait] timeout was raised for jet - %v", jetID.DebugString())
		return core.ErrHotDataTimeout
	}
}

// Unlock raises hotDataChannel
func (hdw *HotDataWaiterConcrete) Unlock(ctx context.Context, jetID core.RecordID) {
	logger := inslogger.FromContext(ctx)
	waiter := hdw.getWaiter(ctx, jetID)

	hdw.waitersMapLock.Lock()
	defer hdw.waitersMapLock.Unlock()

	logger.Debugf("[Unlock] release all requests for jet - %v from waiting", jetID.DebugString())
	close(waiter.hotDataChannel)
}

// ThrowTimeout raises all timeoutChannel
func (hdw *HotDataWaiterConcrete) ThrowTimeout(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("[ThrowTimeout] start method. waiters will be notified about timeout")

	hdw.waitersMapLock.Lock()
	defer hdw.waitersMapLock.Unlock()

	for jetID, waiter := range hdw.waiters {
		logger.Debugf("[ThrowTimeout] raising timeout for requests with jetID - %v", jetID.DebugString())
		close(waiter.timeoutChannel)
	}

	hdw.waiters = map[core.RecordID]*waiter{}
}
