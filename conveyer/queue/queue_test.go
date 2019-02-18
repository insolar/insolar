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

package queue

import (
	"fmt"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getFilledQueue(t *testing.T, numElements int, expectedResult *[]OutputElement) IQueue {
	queue := makeTestQueue()
	for i := 0; i < numElements; i++ {
		require.NoError(t, queue.SinkPush(i))
		*expectedResult = append(*expectedResult, OutputElement{data: i})
	}

	return queue
}

func TestSequentialAccess(t *testing.T) {
	numElements := 300
	expectedResult := make([]OutputElement, 0, numElements*2)
	queue := getFilledQueue(t, numElements, &expectedResult)

	buf := make([]interface{}, 0, numElements)
	for i := 0; i < numElements; i++ {
		el := numElements + i
		buf = append(buf, el)
		expectedResult = append(expectedResult, OutputElement{data: el})
	}
	queue.SinkPushAll(buf)

	total := queue.RemoveAll()
	require.Equal(t, numElements*2, len(total))

	require.EqualValues(t, expectedResult, total)
}

func TestSinkPushAllToEmptyQueue(t *testing.T) {
	queue := makeTestQueue()
	expected := []interface{}{3, 5, 55}
	queue.SinkPushAll(expected)
	require.EqualValues(t, []OutputElement{
		OutputElement{data: expected[0]},
		OutputElement{data: expected[1]},
		OutputElement{data: expected[2]},
	}, queue.RemoveAll())
}

func TestGetFromEmptyQueue(t *testing.T) {
	queue := makeTestQueue()
	for i := 0; i < 100; i++ {
		require.Empty(t, queue.RemoveAll())
	}
}

func TestBlockAndRemoveAll(t *testing.T) {
	numElements := 300
	expectedResult := make([]OutputElement, 0, numElements*2)
	queue := getFilledQueue(t, numElements, &expectedResult)

	total := queue.BlockAndRemoveAll()
	require.EqualValues(t, expectedResult, total)

	require.Contains(t, queue.SinkPush(77).Error(), "Queue is blocked")
	require.Contains(t, queue.SinkPushAll([]interface{}{33}).Error(), "Queue is blocked")
	require.Empty(t, queue.BlockAndRemoveAll())
	require.Empty(t, queue.BlockAndRemoveAll())

}

func TestBlockUnblock(t *testing.T) {
	numElements := 300
	expectedResult := make([]OutputElement, 0, numElements*2)
	queue := getFilledQueue(t, numElements, &expectedResult)

	total := queue.BlockAndRemoveAll()
	require.EqualValues(t, expectedResult, total)
	require.True(t, queue.Unblock())

	testValue := 88
	require.NoError(t, queue.SinkPush(testValue))
	require.EqualValues(t, []OutputElement{OutputElement{data: testValue}}, queue.BlockAndRemoveAll())
}

func TestBlockUnblockEmptyQueue(t *testing.T) {
	queue := makeTestQueue()
	require.False(t, queue.Unblock())
	require.False(t, queue.Unblock())
	require.EqualValues(t, []OutputElement{}, queue.BlockAndRemoveAll())
	require.EqualValues(t, []OutputElement{}, queue.BlockAndRemoveAll())
	require.True(t, queue.Unblock())

	testValue := 88
	require.NoError(t, queue.SinkPush(testValue))
	require.EqualValues(t, []OutputElement{OutputElement{data: testValue}}, queue.RemoveAll())
}

func TestMultipleBlockUnblock(t *testing.T) {
	numElements := 300
	expectedResult := make([]OutputElement, 0, numElements*2)
	queue := getFilledQueue(t, numElements, &expectedResult)

	require.False(t, queue.Unblock())

	total := queue.BlockAndRemoveAll()
	require.EqualValues(t, expectedResult, total)

	require.True(t, queue.Unblock())
}

func TestRemoveAllAfterBlock(t *testing.T) {
	numElements := 200
	expectedResult := make([]OutputElement, 0, numElements*2)
	queue := getFilledQueue(t, numElements, &expectedResult)
	require.EqualValues(t, expectedResult, queue.BlockAndRemoveAll())
	require.Empty(t, queue.RemoveAll())
}

type mockSyncDone struct{}

func (m mockSyncDone) done() {}

// with signals
func TestSimplePushSignals(t *testing.T) {
	queue := makeTestQueue()

	sig := mockSyncDone{}
	numElements := 100

	expected := make([]OutputElement, numElements-1)
	for i := 1; i < numElements; i++ {
		queue.PushSignal(uint32(i), sig)
		expected[numElements-i-1] = OutputElement{data: sig, itemType: uint32(i)}
	}

	require.NotEmpty(t, expected)
	require.EqualValues(t, expected, queue.RemoveAll())
}

func TestPushSignalsAndMessages(t *testing.T) {
	queue := makeTestQueue()

	callback := mockSyncDone{}
	numElements := 201

	expectedSignals := make([]OutputElement, numElements)
	expectedMessages := make([]OutputElement, numElements)
	for i := 0; i < numElements; i++ {
		signal := i + 1
		element := signal * numElements
		queue.SinkPush(element)
		queue.PushSignal(uint32(signal), callback)
		// add signals as LIFO
		expectedSignals[numElements-i-1] = OutputElement{data: callback, itemType: uint32(signal)}

		// add messages as FIFO
		expectedMessages[i] = OutputElement{data: element}
	}

	require.NotEmpty(t, expectedSignals)
	require.NotEmpty(t, expectedMessages)
	require.EqualValues(t, append(expectedSignals, expectedMessages...), queue.RemoveAll())
}

func TestHasSignalAfterPushAll(t *testing.T) {
	queue := makeTestQueue()

	callback := mockSyncDone{}

	require.False(t, queue.HasSignal())
	queue.PushSignal(11, callback)
	queue.SinkPushAll([]interface{}{77, 88, 99})

	require.True(t, queue.HasSignal())
}

func TestPushInvalidSignal(t *testing.T) {
	queue := makeTestQueue()
	assert.Contains(t, queue.PushSignal(0, mockSyncDone{}).Error(), "Unsupported signalType")
}

func TestHasSignal(t *testing.T) {
	queue := makeTestQueue()
	require.False(t, queue.HasSignal())

	queue.PushSignal(3, mockSyncDone{})
	require.True(t, queue.HasSignal())

	queue.RemoveAll()
	require.False(t, queue.HasSignal())

	queue.PushSignal(3, mockSyncDone{})
	queue.SinkPush(3)
	require.True(t, queue.HasSignal())

	queue.RemoveAll()
	require.False(t, queue.HasSignal())

	queue.SinkPush(3)
	queue.PushSignal(3, mockSyncDone{})
	require.True(t, queue.HasSignal())

	queue.RemoveAll()
	require.False(t, queue.HasSignal())
}

func TestBlockUnblockAndHasSignal(t *testing.T) {
	queue := makeTestQueue()
	require.False(t, queue.Unblock())
	require.False(t, queue.HasSignal())

	queue.BlockAndRemoveAll()
	require.False(t, queue.HasSignal())

	require.Error(t, queue.PushSignal(3, mockSyncDone{}))
	require.False(t, queue.HasSignal())
}

func chanToSortedArray(in chan OutputElement, additional []OutputElement) []OutputElement {
	result := make([]OutputElement, 0, len(in)+len(additional))
	for len(in) != 0 {
		result = append(result, <-in)
	}

	result = append(result, additional...)

	sort.Slice(result, func(i, j int) bool {
		isLeftMsg := (result[i].itemType == 0)
		isRightMsg := (result[j].itemType == 0)

		// if both are messages -> compare data fields
		if isLeftMsg && isRightMsg {
			return result[i].data.(int) > result[j].data.(int)
		}
		// if both are signals -> compare itemType fields
		if !isLeftMsg && !isRightMsg {
			return result[i].itemType > result[j].itemType
		}

		return isLeftMsg

	})

	return result
}

func TestParallelAccess(t *testing.T) {
	queue := NewMutexQueue()

	parallelHasSignal := 23
	parallelGet := 25
	parallelPut := 31
	parallelPushSignal := 20
	wg := sync.WaitGroup{}
	wg.Add(parallelPut*2 + parallelGet*2 + parallelHasSignal + parallelPushSignal)

	numIterations := 2

	totalNumOperations := (parallelGet+parallelPut)*numIterations*2 + parallelPushSignal
	addedElements := make(chan OutputElement, totalNumOperations)
	gotElements := make(chan OutputElement, totalNumOperations)
	blockedRequests := make(chan OutputElement, totalNumOperations)

	// SinkPush
	for i := 0; i < parallelPut; i++ {
		go func(wg *sync.WaitGroup, q IQueue, added chan OutputElement, blocked chan OutputElement) {
			for i := 0; i < numIterations; i++ {
				if q.SinkPush(i) == nil {
					added <- OutputElement{data: i}
				} else {
					blockedRequests <- OutputElement{data: i}
				}
			}
			wg.Done()
		}(&wg, queue, addedElements, blockedRequests)
	}

	// RemoveAll
	for i := 0; i < parallelGet; i++ {
		go func(wg *sync.WaitGroup, q IQueue, got chan OutputElement) {
			for i := 0; i < numIterations; i++ {
				results := q.RemoveAll()
				for _, el := range results {
					got <- el
				}
			}
			wg.Done()
		}(&wg, queue, gotElements)
	}

	// SinkPushAll
	for i := 0; i < parallelPut; i++ {
		go func(wg *sync.WaitGroup, q IQueue, added chan OutputElement, blocked chan OutputElement) {
			input := []interface{}{}
			for i := 0; i < numIterations; i++ {
				input = append(input, i)
			}
			if queue.SinkPushAll(input) == nil {
				for _, el := range input {
					added <- OutputElement{data: el}
				}
			} else {
				for _, el := range input {
					blocked <- OutputElement{data: el}
				}
			}
			wg.Done()
		}(&wg, queue, addedElements, blockedRequests)
	}

	// BlockAndRemoveAll - Unblock
	for i := 0; i < parallelGet; i++ {
		go func(wg *sync.WaitGroup, q IQueue, got chan OutputElement) {
			for i := 0; i < numIterations; i++ {
				results := q.BlockAndRemoveAll()
				q.Unblock()
				for _, el := range results {
					got <- el
				}
			}
			wg.Done()
		}(&wg, queue, gotElements)
	}

	// HasSignal
	for i := 0; i < parallelHasSignal; i++ {
		go func(wg *sync.WaitGroup, q IQueue) {
			for i := 0; i < numIterations; i++ {
				q.HasSignal()
			}
			wg.Done()
		}(&wg, queue)
	}

	// PushSignal
	for i := 0; i < parallelPushSignal; i++ {
		go func(wg *sync.WaitGroup, q IQueue, added chan OutputElement, blocked chan OutputElement) {
			for i := 0; i < numIterations; i++ {
				element := OutputElement{data: mockSyncDone{}, itemType: uint32(i)}
				if q.PushSignal(uint32(i), mockSyncDone{}) == nil {
					added <- element
				} else {
					blockedRequests <- element
				}
			}
			wg.Done()
		}(&wg, queue, addedElements, blockedRequests)
	}

	wg.Wait()

	restResult := queue.RemoveAll()

	fmt.Println("Got:                 ", len(gotElements))
	fmt.Println("Added:               ", len(addedElements))
	fmt.Println("Rest:                ", len(restResult))
	fmt.Println("Num blocked requests:", len(blockedRequests))

	require.NotEqual(t, 0, len(gotElements))
	require.NotEqual(t, 0, len(addedElements))

	leftResults := make([]OutputElement, 0, len(restResult))
	for _, el := range restResult {
		leftResults = append(leftResults, el)
	}

	allAddedElements := chanToSortedArray(addedElements, []OutputElement{})
	allGotElements := chanToSortedArray(gotElements, leftResults)

	require.EqualValues(t, allAddedElements, allGotElements)
}
