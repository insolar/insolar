/*
 *    Copyright 2019 Insolar Technologies
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
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/require"
)

func testFuture() future {
	return future{
		result: make(chan insolar.Reply, 1),
	}
}

func TestNewFuture(t *testing.T) {
	conveyorFuture := NewFuture()
	require.Implements(t, (*insolar.ConveyorFuture)(nil), conveyorFuture)
	f := conveyorFuture.(*future)
	require.NotNil(t, f.result)
	require.EqualValues(t, 0, f.finished)
}

func TestFuture_Result(t *testing.T) {
	f := testFuture()
	require.Empty(t, f.Result())

	f.result <- testReply
	require.NotEmpty(t, f.Result())
	res := <-f.Result()
	require.Equal(t, testReply, res)
}

func TestFuture_SetResult(t *testing.T) {
	rep := testReply
	f := testFuture()
	require.Empty(t, f.result)

	f.SetResult(rep)
	require.NotEmpty(t, f.result)
	res := <-f.result
	require.Equal(t, rep, res)
	require.EqualValues(t, 1, f.finished)
}

func TestFuture_SetResult_Multiple(t *testing.T) {
	f := testFuture()
	require.Empty(t, f.result)

	for i := 0; i < 100; i++ {
		f.SetResult(replyMock(i))
	}

	require.Len(t, f.result, 1)
	res := <-f.result
	require.Equal(t, replyMock(0), res)
	require.EqualValues(t, 1, f.finished)
}

func TestFuture_GetResult(t *testing.T) {
	f := testFuture()

	_, err := f.GetResult(10 * time.Millisecond)
	require.EqualError(t, err, ErrFutureTimeout.Error())
	require.EqualValues(t, 1, f.finished)
}

func TestFuture_GetResult_AfterSet(t *testing.T) {
	rep := replyMock(111)
	f := testFuture()
	go func() {
		time.Sleep(time.Millisecond)
		f.SetResult(rep)
	}()

	res, err := f.GetResult(10 * time.Millisecond)
	require.NoError(t, err)
	require.EqualValues(t, rep, res)
	require.EqualValues(t, 1, f.finished)
}

func TestFuture_GetResult_AfterCancel(t *testing.T) {
	f := testFuture()
	go func() {
		time.Sleep(time.Millisecond)
		f.Cancel()
	}()

	_, err := f.GetResult(10 * time.Millisecond)
	require.EqualError(t, err, ErrFutureChannelClosed.Error())
	require.EqualValues(t, 1, f.finished)
}

func TestFuture_Cancel_Multiple(t *testing.T) {
	f := testFuture()

	for i := 0; i < 100; i++ {
		f.Cancel()
	}

	require.EqualValues(t, 1, f.finished)
}

func TestFuture_SetResult_Cancel_Concurrency(t *testing.T) {
	rep := replyMock(111)

	f := testFuture()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		f.Cancel()
		wg.Done()
	}()
	go func() {
		f.SetResult(rep)
		wg.Done()
	}()

	wg.Wait()
	res, ok := <-f.Result()

	cancelDone := res == nil && !ok
	resultDone := res != nil && ok

	require.True(t, cancelDone || resultDone)
	require.EqualValues(t, 1, f.finished)
}
