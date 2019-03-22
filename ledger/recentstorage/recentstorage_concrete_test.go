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

package recentstorage

import (
	"bytes"
	"sync"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/stretchr/testify/require"
)

func TestNewRecentIndexStorage(t *testing.T) {
	t.Parallel()
	jetID := *insolar.NewID(123, []byte{1})
	index := NewRecentIndexStorage(jetID, 123)
	require.NotNil(t, index)
	require.NotNil(t, index.indexes)
	require.Equal(t, 123, index.DefaultTTL)
}

func TestNewRecentIndexStorage_AddId(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	jetID := *insolar.NewID(123, []byte{1})
	s := NewRecentIndexStorage(jetID, 123)

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		s.AddObject(ctx, *insolar.NewID(123, []byte{1}))
		wg.Done()
	}()
	go func() {
		s.AddObject(ctx, *insolar.NewID(123, []byte{2}))
		wg.Done()
	}()
	go func() {
		s.AddObject(ctx, *insolar.NewID(123, []byte{3}))
		wg.Done()
	}()

	wg.Wait()
	require.Equal(t, 3, len(s.GetObjects()))
}

func TestPendingStorage_AddPendingRequest(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	jetID := *insolar.NewID(123, []byte{99})

	s := NewPendingStorage(jetID)

	obj1 := *insolar.NewID(0, nil)
	obj2 := *insolar.NewID(1, nil)

	wg := sync.WaitGroup{}
	wg.Add(3)

	expectedIDs := []insolar.ID{
		*insolar.NewID(123, []byte{1}),
		*insolar.NewID(123, []byte{2}),
		*insolar.NewID(123, []byte{3}),
	}
	go func() {
		s.AddPendingRequest(ctx, obj1, expectedIDs[0])
		wg.Done()
	}()
	go func() {
		s.AddPendingRequest(ctx, obj1, expectedIDs[1])
		wg.Done()
	}()
	go func() {
		s.AddPendingRequest(ctx, obj2, expectedIDs[2])
		wg.Done()
	}()
	wg.Wait()

	contains := func(slice []insolar.ID, x insolar.ID) bool {
		for _, n := range slice {
			if x == n {
				return true
			}
		}
		return false
	}
	requests := s.GetRequests()
	require.Equal(t, 2, len(requests))
	for key, objContext := range requests {
		if bytes.Equal(key.Bytes(), obj1.Bytes()) {
			require.Equal(t, 2, len(objContext.Requests))
			require.Equal(t, true, contains(objContext.Requests, expectedIDs[0]))
			require.Equal(t, true, contains(objContext.Requests, expectedIDs[1]))
		} else {
			require.Equal(t, 1, len(objContext.Requests))
			require.Equal(t, expectedIDs[2], objContext.Requests[0])
		}
	}
}

func TestPendingStorage_RemovePendingRequest(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	jetID := *insolar.NewID(123, []byte{99})

	s := NewPendingStorage(jetID)

	obj := *insolar.NewID(0, nil)

	expectedIDs := []insolar.ID{
		*insolar.NewID(123, []byte{1}),
	}
	extraIDs := []insolar.ID{
		*insolar.NewID(123, []byte{2}),
		*insolar.NewID(123, []byte{3}),
		*insolar.NewID(123, []byte{4}),
	}
	s.requests = map[insolar.ID]*lockedPendingObjectContext{
		obj: {
			Context: &PendingObjectContext{
				Requests: []insolar.ID{expectedIDs[0], extraIDs[0], extraIDs[1], extraIDs[2]},
			},
		},
	}

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		s.RemovePendingRequest(ctx, obj, extraIDs[0])
		wg.Done()
	}()
	go func() {
		s.RemovePendingRequest(ctx, obj, extraIDs[1])
		wg.Done()
	}()
	go func() {
		s.RemovePendingRequest(ctx, obj, extraIDs[2])
		wg.Done()
	}()
	wg.Wait()

	require.Equal(
		t,
		map[insolar.ID]PendingObjectContext{
			obj: {
				Active:   true,
				Requests: []insolar.ID{expectedIDs[0]},
			},
		},
		s.GetRequests(),
	)
}

func TestPendingStorage_RemovePendingRequest_RemoveNothingIfThereIsNothing(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	jetID := *insolar.NewID(123, []byte{99})
	objID := *insolar.NewID(123, []byte{1})
	anotherObj := *insolar.NewID(123, nil)
	s := NewPendingStorage(jetID)

	s.requests = map[insolar.ID]*lockedPendingObjectContext{
		objID: {
			Context: &PendingObjectContext{
				Requests: []insolar.ID{},
			},
		},
	}

	s.RemovePendingRequest(ctx, anotherObj, *insolar.NewID(123, []byte{2}))

	require.Equal(
		t,
		map[insolar.ID]PendingObjectContext{
			objID: {
				Requests: []insolar.ID{},
			},
		},
		s.GetRequests(),
	)
}

func TestNewRecentStorageProvider(t *testing.T) {
	t.Parallel()
	// Act
	provider := NewRecentStorageProvider(888)

	// Assert
	require.Equal(t, 888, provider.DefaultTTL)
	require.NotNil(t, provider.pendingStorages)
	require.NotNil(t, provider.indexStorages)
}

func TestRecentStorageProvider_GetStorage(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := NewRecentStorageProvider(8)

	// Act
	wg := sync.WaitGroup{}
	wg.Add(8)

	for i := 0; i < 8; i++ {
		go func(jetIndex int) {
			jetID := *insolar.NewID(123, []byte{byte(jetIndex)})
			indStorage := provider.GetIndexStorage(inslogger.TestContext(t), jetID)
			require.NotNil(t, indStorage)
			pendingStorage := provider.GetPendingStorage(inslogger.TestContext(t), jetID)
			require.NotNil(t, pendingStorage)
			wg.Done()
		}(i + 1)
	}

	wg.Wait()

	// Assert
	require.Equal(t, 8, len(provider.indexStorages))
	require.Equal(t, 8, len(provider.pendingStorages))
}

func TestRecentStorage_markForDelete(t *testing.T) {
	t.Parallel()
	candidates := make([]insolar.ID, 0, 100)
	expect := make([]insolar.ID, 0, 50)
	recentStorageMap := make(map[insolar.ID]recentObjectMeta)
	jetID := *insolar.NewID(123, []byte{99})

	for i := 0; i < 100; i++ {
		rID := *insolar.NewID(123, []byte{byte(i)})
		candidates = append(candidates, rID)
		if i%2 == 0 {
			expect = append(expect, rID)
		} else {
			recentStorageMap[rID] = recentObjectMeta{i}
		}
	}

	recentStorage := NewRecentIndexStorage(jetID, 888)
	recentStorage.indexes = recentStorageMap

	markedCandidates := recentStorage.markForDelete(candidates)

	require.Equal(t, expect, markedCandidates)
}

func TestRecentStorageProvider_DecreaseIndexesTTL(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)

	firstJet := *insolar.NewID(123, []byte{1})
	secondJet := *insolar.NewID(123, []byte{2})

	provider := NewRecentStorageProvider(8)
	provider.GetIndexStorage(ctx, firstJet).AddObject(ctx, *insolar.NewID(123, []byte{22}))
	provider.GetIndexStorage(ctx, firstJet).AddObject(ctx, *insolar.NewID(123, []byte{33}))

	removedFirst := *insolar.NewID(123, []byte{21})
	removedSecond := *insolar.NewID(123, []byte{22})
	provider.GetIndexStorage(ctx, secondJet).AddObjectWithTLL(ctx, removedFirst, 1)
	provider.GetIndexStorage(ctx, secondJet).AddObjectWithTLL(ctx, removedSecond, 1)

	// Act
	result := provider.DecreaseIndexesTTL(ctx)

	// Assert
	provider.indexLock.Lock()
	defer provider.indexLock.Unlock()
	require.NotNil(t, result)
	require.Equal(t, 1, len(provider.indexStorages))
	require.Equal(t, 1, len(result))
	require.Equal(t, 2, len(result[secondJet]))
	if removedFirst != result[secondJet][0] && removedFirst != result[secondJet][1] {
		require.Fail(t, "return result is broken")
	}
	if removedSecond != result[secondJet][1] && removedSecond != result[secondJet][0] {
		require.Fail(t, "return result is broken")
	}
	for _, index := range provider.indexStorages[firstJet].indexes {
		require.Equal(t, 7, index.ttl)
	}
}

func TestRecentStorageProvider_DecreaseIndexesTTL_WorksOnEmptyStorage(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)
	provider := NewRecentStorageProvider(8)

	// Act
	result := provider.DecreaseIndexesTTL(ctx)

	// Assert
	require.Equal(t, map[insolar.ID][]insolar.ID{}, result)
}

func TestPendingStorageConcrete_GetRequestsForObject(t *testing.T) {
	t.Parallel()

	objID := *insolar.NewID(123, []byte{1})
	requestID := *insolar.NewID(123, []byte{2})

	unexpectedID := *insolar.NewID(123, []byte{3})
	unexpectedReqID := *insolar.NewID(123, []byte{4})

	pendingStorage := &PendingStorageConcrete{
		requests: map[insolar.ID]*lockedPendingObjectContext{
			objID: {
				Context: &PendingObjectContext{
					Requests: []insolar.ID{requestID},
				},
			},
			unexpectedID: {
				Context: &PendingObjectContext{
					Requests: []insolar.ID{unexpectedReqID},
				},
			},
		},
	}

	requests := pendingStorage.GetRequestsForObject(objID)

	require.Equal(t, 1, len(requests))
	require.Equal(t, requestID, requests[0])
}

func TestPendingStorageConcrete_GetRequestsForObject_NoObject(t *testing.T) {
	t.Parallel()

	unexpectedReqID := *insolar.NewID(123, []byte{2})

	pendingStorage := &PendingStorageConcrete{
		requests: map[insolar.ID]*lockedPendingObjectContext{},
	}

	requests := pendingStorage.GetRequestsForObject(unexpectedReqID)

	require.Nil(t, requests)
}

func TestPendingStorageConcrete_SetContextToObject(t *testing.T) {
	t.Parallel()

	pendingStorage := &PendingStorageConcrete{
		requests: map[insolar.ID]*lockedPendingObjectContext{},
	}
	expectedObj := *insolar.NewID(123, []byte{1})
	expectedContext := PendingObjectContext{
		Active:   true,
		Requests: []insolar.ID{*insolar.NewID(123, []byte{2}), *insolar.NewID(123, []byte{3})},
	}

	pendingStorage.SetContextToObject(inslogger.TestContext(t), expectedObj, expectedContext)

	require.Equal(t, 1, len(pendingStorage.requests))
	require.Equal(t, expectedContext, *pendingStorage.requests[expectedObj].Context)
}

func TestPendingStorageConcrete_RemovePendingRequest_RemoveFromStart(t *testing.T) {
	t.Parallel()

	objID := *insolar.NewID(123, []byte{100})
	first := *insolar.NewID(123, []byte{1})
	second := *insolar.NewID(123, []byte{2})
	pendingStorage := &PendingStorageConcrete{
		requests: map[insolar.ID]*lockedPendingObjectContext{

			objID: {
				Context: &PendingObjectContext{
					Requests: []insolar.ID{first, second},
				},
			},
		},
	}

	pendingStorage.RemovePendingRequest(inslogger.TestContext(t), objID, first)

	require.Equal(t, 1, len(pendingStorage.requests[objID].Context.Requests))
	require.Equal(t, second, pendingStorage.requests[objID].Context.Requests[0])
}

func TestPendingStorageConcrete_RemovePendingRequest_RemoveFromEnd(t *testing.T) {
	t.Parallel()

	objID := *insolar.NewID(123, []byte{100})
	first := *insolar.NewID(123, []byte{1})
	second := *insolar.NewID(123, []byte{2})
	pendingStorage := &PendingStorageConcrete{
		requests: map[insolar.ID]*lockedPendingObjectContext{

			objID: {
				Context: &PendingObjectContext{
					Requests: []insolar.ID{first, second},
				},
			},
		},
	}

	pendingStorage.RemovePendingRequest(inslogger.TestContext(t), objID, second)

	require.Equal(t, 1, len(pendingStorage.requests[objID].Context.Requests))
	require.Equal(t, first, pendingStorage.requests[objID].Context.Requests[0])
}

func TestPendingStorageConcrete_RemovePendingRequest_RemoveFromMiddle(t *testing.T) {
	t.Parallel()

	objID := *insolar.NewID(123, []byte{100})
	first := *insolar.NewID(123, []byte{1})
	second := *insolar.NewID(123, []byte{2})
	third := *insolar.NewID(123, []byte{3})
	pendingStorage := &PendingStorageConcrete{
		requests: map[insolar.ID]*lockedPendingObjectContext{

			objID: {
				Context: &PendingObjectContext{
					Requests: []insolar.ID{first, second, third},
				},
			},
		},
	}

	pendingStorage.RemovePendingRequest(inslogger.TestContext(t), objID, second)

	require.Equal(t, 2, len(pendingStorage.requests[objID].Context.Requests))
	require.Equal(t, first, pendingStorage.requests[objID].Context.Requests[0])
	require.Equal(t, third, pendingStorage.requests[objID].Context.Requests[1])
}

func TestPendingStorageConcrete_RemovePendingRequest_NothingHappensIfNoRequests(t *testing.T) {
	t.Parallel()

	objID := *insolar.NewID(123, []byte{100})
	pendingStorage := &PendingStorageConcrete{
		requests: map[insolar.ID]*lockedPendingObjectContext{

			objID: {
				Context: &PendingObjectContext{},
			},
		},
	}

	pendingStorage.RemovePendingRequest(inslogger.TestContext(t), objID, *insolar.NewID(123, []byte{1}))

	require.Equal(t, 1, len(pendingStorage.requests))
	_, ok := pendingStorage.requests[objID]
	require.Equal(t, true, ok)
}

func TestPendingStorageConcrete_RemovePendingRequest_RemoveOnlyOne(t *testing.T) {
	t.Parallel()

	objID := *insolar.NewID(123, []byte{100})
	first := *insolar.NewID(123, []byte{1})

	pendingStorage := &PendingStorageConcrete{
		requests: map[insolar.ID]*lockedPendingObjectContext{
			objID: {
				Context: &PendingObjectContext{
					Requests: []insolar.ID{first},
				},
			},
		},
	}

	pendingStorage.RemovePendingRequest(inslogger.TestContext(t), objID, first)

	require.Equal(t, 0, len(pendingStorage.requests[objID].Context.Requests))
}

func TestPendingStorageConcrete_RemovePendingRequest_RemoveNotExisting(t *testing.T) {
	t.Parallel()

	objID := *insolar.NewID(123, []byte{100})
	first := *insolar.NewID(123, []byte{1})
	second := *insolar.NewID(123, []byte{2})
	third := *insolar.NewID(123, []byte{3})
	pendingStorage := &PendingStorageConcrete{
		requests: map[insolar.ID]*lockedPendingObjectContext{

			objID: {
				Context: &PendingObjectContext{
					Requests: []insolar.ID{first, second},
				},
			},
		},
	}

	pendingStorage.RemovePendingRequest(inslogger.TestContext(t), objID, third)

	require.Equal(t, 2, len(pendingStorage.requests[objID].Context.Requests))
	require.Equal(t, first, pendingStorage.requests[objID].Context.Requests[0])
	require.Equal(t, second, pendingStorage.requests[objID].Context.Requests[1])
}
