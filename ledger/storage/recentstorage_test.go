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

package storage

import (
	"sync"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/require"
)

func TestNewRecentObjectsIndex(t *testing.T) {
	index := NewRecentStorage(123)
	require.NotNil(t, index)
	require.NotNil(t, index.recentObjects)
	require.Equal(t, 123, index.DefaultTTL)
}

func TestRecentObjectsIndex_AddId(t *testing.T) {
	s := NewRecentStorage(123)

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		s.AddObject(*core.NewRecordID(123, []byte{1}))
		wg.Done()
	}()
	go func() {
		s.AddObject(*core.NewRecordID(123, []byte{2}))
		wg.Done()
	}()
	go func() {
		s.AddObject(*core.NewRecordID(123, []byte{3}))
		wg.Done()
	}()

	wg.Wait()
	require.Equal(t, 3, len(s.GetObjects()))
}

func TestRecentObjectsIndex_AddPendingRequest(t *testing.T) {
	s := NewRecentStorage(123)

	obj1 := *core.NewRecordID(0, nil)
	obj2 := *core.NewRecordID(1, nil)

	wg := sync.WaitGroup{}
	wg.Add(3)

	expectedIDs := []core.RecordID{
		*core.NewRecordID(123, []byte{1}),
		*core.NewRecordID(123, []byte{2}),
		*core.NewRecordID(123, []byte{3}),
	}
	go func() {
		s.AddPendingRequest(obj1, expectedIDs[0])
		wg.Done()
	}()
	go func() {
		s.AddPendingRequest(obj1, expectedIDs[1])
		wg.Done()
	}()
	go func() {
		s.AddPendingRequest(obj2, expectedIDs[2])
		wg.Done()
	}()
	wg.Wait()

	require.Equal(t, map[core.RecordID]map[core.RecordID]struct{}{
		obj1: {
			expectedIDs[0]: struct{}{},
			expectedIDs[1]: struct{}{},
		},
		obj2: {
			expectedIDs[2]: struct{}{},
		},
	}, s.GetRequests())
}

func TestRecentObjectsIndex_RemovePendingRequest(t *testing.T) {
	s := NewRecentStorage(123)

	obj := *core.NewRecordID(0, nil)

	wg := sync.WaitGroup{}
	wg.Add(3)

	expectedIDs := []core.RecordID{
		*core.NewRecordID(123, []byte{1}),
	}
	extraIDs := []core.RecordID{
		*core.NewRecordID(123, []byte{2}),
		*core.NewRecordID(123, []byte{3}),
		*core.NewRecordID(123, []byte{4}),
	}
	s.pendingRequests = map[core.RecordID]map[core.RecordID]struct{}{
		obj: {
			expectedIDs[0]: {},
			extraIDs[0]:    {},
			extraIDs[1]:    {},
			extraIDs[2]:    {},
		},
	}

	go func() {
		s.RemovePendingRequest(obj, extraIDs[0])
		wg.Done()
	}()
	go func() {
		s.RemovePendingRequest(obj, extraIDs[1])
		wg.Done()
	}()
	go func() {
		s.RemovePendingRequest(obj, extraIDs[2])
		wg.Done()
	}()
	wg.Wait()

	require.Equal(t, map[core.RecordID]map[core.RecordID]struct{}{
		obj: {
			expectedIDs[0]: struct{}{},
		},
	}, s.GetRequests())
}

func TestRecentObjectsIndex_ClearObjects(t *testing.T) {
	index := NewRecentStorage(123)
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		index.AddObject(*core.NewRecordID(123, []byte{1}))
		wg.Done()
	}()
	go func() {
		index.AddObject(*core.NewRecordID(123, []byte{2}))
		wg.Done()
	}()
	go func() {
		index.AddObject(*core.NewRecordID(123, []byte{3}))
		wg.Done()
	}()
	wg.Wait()

	index.ClearObjects()

	require.Equal(t, 0, len(index.GetObjects()))
}

func TestNewRecentStorageProvider(t *testing.T) {
	// Act
	provider := NewRecentStorageProvider(888)

	// Assert
	require.Equal(t, 888, provider.DefaultTTL)
	require.NotNil(t, provider.storage)
}

func TestRecentStorageProvider_GetStorage(t *testing.T) {
	// Arrange
	provider := NewRecentStorageProvider(8)

	// Act
	wg := sync.WaitGroup{}
	wg.Add(8)

	for i := 0; i < 8; i++ {
		i := i
		go func() {
			id := core.NewRecordID(core.FirstPulseNumber, []byte{byte(i)})
			storage := provider.GetStorage(*id)
			require.NotNil(t, storage)
			wg.Done()
		}()
	}

	wg.Wait()

	// Assert
	require.Equal(t, 8, len(provider.storage))
}
