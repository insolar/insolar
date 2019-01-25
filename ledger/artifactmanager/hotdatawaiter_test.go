package artifactmanager

import (
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestNewHotDataWaiterConcrete(t *testing.T) {
	// Act
	hdw := NewHotDataWaiterConcrete()

	// Assert
	require.NotNil(t, hdw)
	require.NotNil(t, hdw.waiters)
}

func TestHotDataWaiterConcrete_Get_CreateIfNil(t *testing.T) {
	// Arrange
	hdw := NewHotDataWaiterConcrete()
	jetID := testutils.RandomID()

	// Act
	waiter := hdw.getWaiter(inslogger.TestContext(t), jetID)

	// Assert
	require.NotNil(t, waiter)
	require.NotNil(t, waiter.hotDataChannel)
	require.NotNil(t, waiter.timeoutChannel)
	require.Equal(t, waiter, hdw.waiters[jetID])
	require.Equal(t, 1, len(hdw.waiters))
}

func TestHotDataWaiterConcrete_Wait_UnlockHotData(t *testing.T) {
	// Arrange
	syncChannel := make(chan struct{})
	hdw := NewHotDataWaiterConcrete()
	hdwLock := sync.Mutex{}
	hdwGetter := func() HotDataWaiter {
		hdwLock.Lock()
		defer hdwLock.Unlock()

		return hdw
	}
	jetID := testutils.RandomID()
	_ = hdw.getWaiter(inslogger.TestContext(t), jetID)

	// Act
	go func() {
		err := hdwGetter().Wait(inslogger.TestContext(t), jetID)
		require.Nil(t, err)
		close(syncChannel)
	}()

	time.Sleep(10 * time.Millisecond)

	hdwGetter().Unlock(inslogger.TestContext(t), jetID)

	<-syncChannel
}

func TestHotDataWaiterConcrete_Wait_ThrowTimeout(t *testing.T) {
	// Arrange
	syncChannel := make(chan struct{})
	hdw := NewHotDataWaiterConcrete()
	hdwLock := sync.Mutex{}
	hdwGetter := func() HotDataWaiter {
		hdwLock.Lock()
		defer hdwLock.Unlock()

		return hdw
	}
	hdwLengthGetter := func() int {
		hdwLock.Lock()
		defer hdwLock.Unlock()

		return len(hdw.waiters)
	}
	jetID := testutils.RandomID()
	_ = hdw.getWaiter(inslogger.TestContext(t), jetID)

	// Act
	go func() {
		err := hdwGetter().Wait(inslogger.TestContext(t), jetID)
		require.NotNil(t, err)
		require.Equal(t, core.ErrHotDataTimeout, err)
		close(syncChannel)
	}()

	time.Sleep(10 * time.Millisecond)

	hdwGetter().ThrowTimeout(inslogger.TestContext(t))

	<-syncChannel
	require.Equal(t, 0, hdwLengthGetter())
}

func TestHotDataWaiterConcrete_Wait_ThrowTimeout_MultipleMembers(t *testing.T) {
	// Arrange
	syncChannel := make(chan struct{})
	hdw := NewHotDataWaiterConcrete()
	hdwLock := sync.Mutex{}
	hdwGetter := func() HotDataWaiter {
		hdwLock.Lock()
		defer hdwLock.Unlock()

		return hdw
	}
	hdwLengthGetter := func() int {
		hdwLock.Lock()
		defer hdwLock.Unlock()

		return len(hdw.waiters)
	}
	jetID := testutils.RandomID()
	secondJetID := testutils.RandomID()
	_ = hdw.getWaiter(inslogger.TestContext(t), jetID)

	// Act
	go func() {
		err := hdwGetter().Wait(inslogger.TestContext(t), jetID)
		require.NotNil(t, err)
		require.Equal(t, core.ErrHotDataTimeout, err)
		syncChannel <- struct{}{}
	}()
	go func() {
		err := hdwGetter().Wait(inslogger.TestContext(t), secondJetID)
		require.NotNil(t, err)
		require.Equal(t, core.ErrHotDataTimeout, err)
		syncChannel <- struct{}{}
	}()

	time.Sleep(10 * time.Millisecond)

	hdwGetter().ThrowTimeout(inslogger.TestContext(t))

	<-syncChannel
	<-syncChannel

	require.Equal(t, 0, hdwLengthGetter())
}
