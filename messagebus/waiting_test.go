/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package messagebus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewWaitingPool(t *testing.T) {
	t.Parallel()
	wp := newWaitingPool(100)
	require.NotNil(t, wp.waiting)
	require.Equal(t, waitingPool{
		waiting: wp.waiting,
		limit:   100,
	}, wp)
}

func TestWaitingPool_WaitingChan(t *testing.T) {
	t.Parallel()
	wp := newWaitingPool(5)
	stop := true

	go func() {
		time.Sleep(time.Second)
		stop = false
		wp.clear()
	}()

	<-wp.waiting
	require.False(t, stop)
}

func TestWaitingPool_LockInAcquire(t *testing.T) {
	t.Parallel()
	wp := newWaitingPool(5)
	stop := true

	go func() {
		time.Sleep(time.Second)
		stop = false
		wp.locker.Unlock()
	}()

	wp.locker.Lock()
	wp.acquire()
	require.False(t, stop)
}

func TestWaitingPool_LockInClear(t *testing.T) {
	t.Parallel()
	wp := newWaitingPool(5)
	stop := true

	go func() {
		time.Sleep(time.Second)
		stop = false
		wp.locker.Unlock()
	}()

	wp.locker.Lock()
	wp.clear()
	require.False(t, stop)
}

func TestWaitingPool_AcquireLimit(t *testing.T) {
	t.Parallel()
	wp := newWaitingPool(5)

	require.True(t, wp.acquire())
	require.True(t, wp.acquire())
	require.True(t, wp.acquire())
	require.True(t, wp.acquire())
	require.True(t, wp.acquire())
	require.False(t, wp.acquire())

	wp.clear()
	require.True(t, wp.acquire())
}
