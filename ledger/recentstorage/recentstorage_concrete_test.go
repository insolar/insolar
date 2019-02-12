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

package recentstorage

import (
	"sync"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestNewRecentIndexStorage(t *testing.T) {
	t.Parallel()
	jetID := testutils.RandomID()
	index := NewRecentIndexStorage(jetID, 123)
	require.NotNil(t, index)
	require.NotNil(t, index.indexes)
	require.Equal(t, 123, index.DefaultTTL)
}

func TestNewRecentIndexStorage_AddId(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	jetID := testutils.RandomID()
	s := NewRecentIndexStorage(jetID, 123)

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		s.AddObject(ctx, *core.NewRecordID(123, []byte{1}))
		wg.Done()
	}()
	go func() {
		s.AddObject(ctx, *core.NewRecordID(123, []byte{2}))
		wg.Done()
	}()
	go func() {
		s.AddObject(ctx, *core.NewRecordID(123, []byte{3}))
		wg.Done()
	}()

	wg.Wait()
	require.Equal(t, 3, len(s.GetObjects()))
}

func TestPendingStorage_AddPendingRequest(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	jetID := testutils.RandomID()

	s := NewPendingStorage(jetID)

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

	require.Equal(
		t,
		map[core.RecordID]PendingObjectContext{
			obj1: {Active: true, Requests: []core.RecordID{expectedIDs[0], expectedIDs[1]}},
			obj2: {Active: true, Requests: []core.RecordID{expectedIDs[2]}},
		},
		s.GetRequests(),
	)
}

func TestPendingStorage_RemovePendingRequest(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	jetID := testutils.RandomID()

	s := NewPendingStorage(jetID)

	obj := *core.NewRecordID(0, nil)

	expectedIDs := []core.RecordID{
		*core.NewRecordID(123, []byte{1}),
	}
	extraIDs := []core.RecordID{
		*core.NewRecordID(123, []byte{2}),
		*core.NewRecordID(123, []byte{3}),
		*core.NewRecordID(123, []byte{4}),
	}
	s.requests = map[core.RecordID]*lockedPendingObjectContext{
		obj: {
			Context: &PendingObjectContext{
				Requests: []core.RecordID{expectedIDs[0], extraIDs[0], extraIDs[1], extraIDs[2]},
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
		map[core.RecordID]PendingObjectContext{
			obj: {
				Active:   true,
				Requests: []core.RecordID{expectedIDs[0]},
			},
		},
		s.GetRequests(),
	)
}

func TestPendingStorage_RemovePendingRequest_RemoveNothingIfThereIsNothing(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	jetID := testutils.RandomID()
	objID := testutils.RandomID()
	anotherObj := *core.NewRecordID(123, nil)
	s := NewPendingStorage(jetID)

	s.requests = map[core.RecordID]*lockedPendingObjectContext{
		objID: {
			Context: &PendingObjectContext{
				Requests: []core.RecordID{},
			},
		},
	}

	s.RemovePendingRequest(ctx, anotherObj, testutils.RandomID())

	require.Equal(
		t,
		map[core.RecordID]PendingObjectContext{
			objID: {
				Requests: []core.RecordID{},
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
		go func() {
			id := testutils.RandomJet()
			storage := provider.GetIndexStorage(inslogger.TestContext(t), id)
			require.NotNil(t, storage)
			pendingStorage := provider.GetPendingStorage(inslogger.TestContext(t), id)
			require.NotNil(t, pendingStorage)
			wg.Done()
		}()
	}

	wg.Wait()

	// Assert
	require.Equal(t, 8, len(provider.indexStorages))
	require.Equal(t, 8, len(provider.pendingStorages))
}

func TestRecentStorage_markForDelete(t *testing.T) {
	t.Parallel()
	candidates := make([]core.RecordID, 0, 100)
	expect := make([]core.RecordID, 0, 50)
	recentStorageMap := make(map[core.RecordID]recentObjectMeta)

	for i := 0; i < 100; i++ {
		rID := testutils.RandomID()
		candidates = append(candidates, rID)
		if i%2 == 0 {
			expect = append(expect, rID)
		} else {
			recentStorageMap[rID] = recentObjectMeta{i}
		}
	}

	recentStorage := NewRecentIndexStorage(testutils.RandomJet(), 888)
	recentStorage.indexes = recentStorageMap

	markedCandidates := recentStorage.markForDelete(candidates)

	require.Equal(t, expect, markedCandidates)
}

func TestRecentStorageProvider_DecreaseIndexesTTL(t *testing.T) {
	t.Parallel()
	// Arrange
	ctx := inslogger.TestContext(t)

	firstJet := testutils.RandomJet()
	secondJet := testutils.RandomJet()

	provider := NewRecentStorageProvider(8)
	provider.GetIndexStorage(ctx, firstJet).AddObject(ctx, testutils.RandomID())
	provider.GetIndexStorage(ctx, firstJet).AddObject(ctx, testutils.RandomID())

	removedFirst := testutils.RandomID()
	removedSecond := testutils.RandomID()
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
	require.Equal(t, map[core.RecordID][]core.RecordID{}, result)
}

func TestPendingStorageConcrete_GetRequestsForObject(t *testing.T) {
	t.Parallel()

	objID := testutils.RandomID()
	requestID := testutils.RandomID()

	unexpectedID := testutils.RandomID()
	unexpectedReqID := testutils.RandomID()

	pendingStorage := &PendingStorageConcrete{
		requests: map[core.RecordID]*lockedPendingObjectContext{
			objID: {
				Context: &PendingObjectContext{
					Requests: []core.RecordID{requestID},
				},
			},
			unexpectedID: {
				Context: &PendingObjectContext{
					Requests: []core.RecordID{unexpectedReqID},
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

	unexpectedReqID := testutils.RandomID()

	pendingStorage := &PendingStorageConcrete{
		requests: map[core.RecordID]*lockedPendingObjectContext{},
	}

	requests := pendingStorage.GetRequestsForObject(unexpectedReqID)

	require.Nil(t, requests)
}

func TestPendingStorageConcrete_SetContextToObject(t *testing.T) {
	t.Parallel()
	pendingStorage := &PendingStorageConcrete{
		requests: map[core.RecordID]*lockedPendingObjectContext{},
	}
	expectedObj := testutils.RandomID()
	expectedContext := PendingObjectContext{
		Active:   true,
		Requests: []core.RecordID{testutils.RandomID(), testutils.RandomID()},
	}

	pendingStorage.SetContextToObject(inslogger.TestContext(t), expectedObj, expectedContext)

	require.Equal(t, 1, len(pendingStorage.requests))
	require.Equal(t, expectedContext, *pendingStorage.requests[expectedObj].Context)
}
