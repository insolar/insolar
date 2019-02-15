/*
 *    Copyright 2019 Insolar
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

package queue

import (
	"fmt"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func getFilledQueue(t *testing.T, numElements int, expectedResult *[]interface{}) IQueue {
	queue := makeTestQueue()
	for i := 0; i < numElements; i++ {
		require.True(t, queue.SinkPush(i))
		*expectedResult = append(*expectedResult, i)
	}

	return queue
}

func TestSequentialAccess(t *testing.T) {
	numElements := 300
	expectedResult := make([]interface{}, 0, numElements*2)
	queue := getFilledQueue(t, numElements, &expectedResult)

	buf := make([]interface{}, 0, numElements)
	for i := 0; i < numElements; i++ {
		el := numElements + i
		buf = append(buf, el)
		expectedResult = append(expectedResult, el)
	}
	queue.SinkPushAll(buf)

	fmt.Println("SECOND PART: ", buf)

	total := queue.RemoveAll()
	require.Equal(t, numElements*2, len(total))

	require.EqualValues(t, expectedResult, total)
}

func TestGetFromEmptyQueue(t *testing.T) {
	queue := makeTestQueue()
	for i := 0; i < 100; i++ {
		require.Empty(t, queue.RemoveAll())
	}
}

func TestBlockAndRemoveAll(t *testing.T) {
	numElements := 300
	expectedResult := make([]interface{}, 0, numElements*2)
	queue := getFilledQueue(t, numElements, &expectedResult)

	total := queue.BlockAndRemoveAll()
	require.EqualValues(t, expectedResult, total)

	require.False(t, queue.SinkPush(77))
	require.False(t, queue.SinkPushAll([]interface{}{33}))
	require.Empty(t, queue.BlockAndRemoveAll())
	require.Empty(t, queue.BlockAndRemoveAll())

}

func TestBlockUnblock(t *testing.T) {
	numElements := 300
	expectedResult := make([]interface{}, 0, numElements*2)
	queue := getFilledQueue(t, numElements, &expectedResult)

	total := queue.BlockAndRemoveAll()
	require.EqualValues(t, expectedResult, total)
	require.True(t, queue.Unblock())

	testValue := 88
	require.True(t, queue.SinkPush(testValue))
	require.EqualValues(t, []interface{}{testValue}, queue.BlockAndRemoveAll())
}

func TestBlockUnblockEmptyQueue(t *testing.T) {
	queue := makeTestQueue()
	require.False(t, queue.Unblock())
	require.False(t, queue.Unblock())
	require.EqualValues(t, []interface{}{}, queue.BlockAndRemoveAll())
	require.EqualValues(t, []interface{}{}, queue.BlockAndRemoveAll())
	require.False(t, queue.Unblock())

	testValue := 88
	require.True(t, queue.SinkPush(testValue))
	require.EqualValues(t, []interface{}{testValue}, queue.RemoveAll())
}

func TestMultipleBlockUnblock(t *testing.T) {
	numElements := 300
	expectedResult := make([]interface{}, 0, numElements*2)
	queue := getFilledQueue(t, numElements, &expectedResult)

	require.False(t, queue.Unblock())

	total := queue.BlockAndRemoveAll()
	require.EqualValues(t, expectedResult, total)

	require.True(t, queue.Unblock())
}

func TestRemoveAllAfterBlock(t *testing.T) {
	numElements := 300
	expectedResult := make([]interface{}, 0, numElements*2)
	queue := getFilledQueue(t, numElements, &expectedResult)
	require.EqualValues(t, expectedResult, queue.BlockAndRemoveAll())
	require.Empty(t, queue.RemoveAll())
}

func TestParallelAccess(t *testing.T) {
	queue := NewMutexQueue()

	parallelGet := 30

	parallelPut := 30
	wg := sync.WaitGroup{}
	wg.Add(parallelPut*2 + parallelGet)

	numIterations := 100

	for i := 0; i < parallelPut; i++ {
		go func(wg *sync.WaitGroup, q IQueue) {
			for i := 0; i < numIterations; i++ {
				q.SinkPush(i)
			}
			wg.Done()
		}(&wg, queue)
	}

	for i := 0; i < parallelGet; i++ {
		go func(wg *sync.WaitGroup, q IQueue) {
			for i := 0; i < numIterations; i++ {
				q.RemoveAll()
			}
			wg.Done()
		}(&wg, queue)
	}

	for i := 0; i < parallelPut; i++ {
		go func(wg *sync.WaitGroup, q IQueue) {
			input := []interface{}{}
			for i := 0; i < numIterations; i++ {
				input = append(input, i)
			}
			queue.SinkPushAll(input)
			wg.Done()
		}(&wg, queue)
	}

	wg.Wait()

	result := queue.RemoveAll()

	sort.Slice(result, func(i, j int) bool {
		return result[i].(int) > result[j].(int)
	})

	fmt.Println(result)

	fmt.Printf("Num Elements: %d . Must be: %d\n", len(result), parallelPut*numIterations)

}
