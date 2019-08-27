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

package executor_test

import (
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_HotDataWaiterConcrete_WaitUnlock(t *testing.T) {
	t.Parallel()
	// Arrange
	waitingStarted := make(chan struct{}, 1)
	waitingFinished := make(chan struct{})

	hdw := executor.NewChannelWaiter()
	hdwLock := sync.Mutex{}
	hdwGetter := func() *executor.ChannelWaiter {
		hdwLock.Lock()
		defer hdwLock.Unlock()

		return hdw
	}
	jetID := gen.JetID()
	pulse := gen.PulseNumber()

	// Act
	go func() {
		waitingStarted <- struct{}{}
		err := hdwGetter().Wait(inslogger.TestContext(t), jetID, pulse)
		require.Nil(t, err)
		close(waitingFinished)
	}()

	<-waitingStarted
	time.Sleep(1 * time.Second)

	// Closing waiter the first time, no error.
	err := hdwGetter().Unlock(inslogger.TestContext(t), pulse, jetID)
	require.NoError(t, err)

	<-waitingFinished

	// Closing waiter the second time, error.
	err = hdwGetter().Unlock(inslogger.TestContext(t), pulse, jetID)
	assert.Error(t, err)
}

func Test_HotDataWaiterConcrete_Close(t *testing.T) {
	t.Parallel()
	// Arrange
	waitingStarted := make(chan struct{}, 1)
	waitingFinished := make(chan struct{})

	hdw := executor.NewChannelWaiter()
	hdwLock := sync.Mutex{}
	hdwGetter := func() *executor.ChannelWaiter {
		hdwLock.Lock()
		defer hdwLock.Unlock()

		return hdw
	}

	jetID := gen.JetID()
	pulse := gen.PulseNumber()

	// Act
	go func() {
		waitingStarted <- struct{}{}
		err := hdwGetter().Wait(inslogger.TestContext(t), jetID, pulse)
		require.NotNil(t, err)
		require.Equal(t, insolar.ErrHotDataTimeout, err)
		close(waitingFinished)
	}()

	<-waitingStarted
	time.Sleep(1 * time.Second)

	hdwGetter().CloseAllUntil(inslogger.TestContext(t), pulse)

	<-waitingFinished

	err := hdwGetter().Wait(inslogger.TestContext(t), jetID, pulse)
	require.Nil(t, err)
}

func Test_HotDataWaiterConcrete_WaitClose_MultipleMembers(t *testing.T) {
	t.Parallel()
	// Arrange
	waitingStarted := make(chan struct{}, 2)
	waitingRes := make(chan error, 2)

	hdw := executor.NewChannelWaiter()
	hdwLock := sync.Mutex{}
	hdwGetter := func() *executor.ChannelWaiter {
		hdwLock.Lock()
		defer hdwLock.Unlock()

		return hdw
	}

	pulse := gen.PulseNumber()
	jetID := gen.JetID()
	secondJetID := gen.JetID()

	// Act
	go func() {
		waitingStarted <- struct{}{}
		waitingRes <- hdwGetter().Wait(inslogger.TestContext(t), jetID, pulse)
	}()
	go func() {
		waitingStarted <- struct{}{}
		waitingRes <- hdwGetter().Wait(inslogger.TestContext(t), secondJetID, pulse)
	}()

	<-waitingStarted
	<-waitingStarted
	time.Sleep(1 * time.Second)

	hdwGetter().CloseAllUntil(inslogger.TestContext(t), pulse)

	err := <-waitingRes
	require.Equal(t, err, insolar.ErrHotDataTimeout)
	err = <-waitingRes
	require.Equal(t, err, insolar.ErrHotDataTimeout)
}
