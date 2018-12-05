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
	"bytes"
	"sort"
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
		s.AddObject(*core.NewRecordID(123, []byte{1}), false)
		wg.Done()
	}()
	go func() {
		s.AddObject(*core.NewRecordID(123, []byte{2}), false)
		wg.Done()
	}()
	go func() {
		s.AddObject(*core.NewRecordID(123, []byte{3}), false)
		wg.Done()
	}()

	wg.Wait()
	require.Equal(t, 3, len(s.GetObjects()))
}

func TestRecentObjectsIndex_AddPendingRequest(t *testing.T) {
	s := NewRecentStorage(123)

	wg := sync.WaitGroup{}
	wg.Add(3)

	expectedIDs := []core.RecordID{
		*core.NewRecordID(123, []byte{1}),
		*core.NewRecordID(123, []byte{2}),
		*core.NewRecordID(123, []byte{3}),
	}
	go func() {
		s.AddPendingRequest(expectedIDs[0])
		wg.Done()
	}()
	go func() {
		s.AddPendingRequest(expectedIDs[1])
		wg.Done()
	}()
	go func() {
		s.AddPendingRequest(expectedIDs[2])
		wg.Done()
	}()
	wg.Wait()

	actualRequests := s.GetRequests()
	sort.Slice(actualRequests, func(i, j int) bool {
		return bytes.Compare(actualRequests[i][:], actualRequests[j][:]) == -1
	})
	require.Equal(t, expectedIDs, actualRequests)
}

func TestRecentObjectsIndex_RemovePendingRequest(t *testing.T) {
	s := NewRecentStorage(123)

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
	s.pendingRequests = map[core.RecordID]struct{}{
		expectedIDs[0]: {},
		extraIDs[0]:    {},
		extraIDs[1]:    {},
		extraIDs[2]:    {},
	}

	go func() {
		s.RemovePendingRequest(extraIDs[0])
		wg.Done()
	}()
	go func() {
		s.RemovePendingRequest(extraIDs[1])
		wg.Done()
	}()
	go func() {
		s.RemovePendingRequest(extraIDs[2])
		wg.Done()
	}()
	wg.Wait()

	actualRequests := s.GetRequests()
	sort.Slice(actualRequests, func(i, j int) bool {
		return bytes.Compare(actualRequests[i][:], actualRequests[j][:]) == -1
	})
	require.Equal(t, expectedIDs, actualRequests)
}

func TestRecentObjectsIndex_ClearObjects(t *testing.T) {
	index := NewRecentStorage(123)
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		index.AddObject(*core.NewRecordID(123, []byte{1}), false)
		wg.Done()
	}()
	go func() {
		index.AddObject(*core.NewRecordID(123, []byte{2}), false)
		wg.Done()
	}()
	go func() {
		index.AddObject(*core.NewRecordID(123, []byte{3}), false)
		wg.Done()
	}()
	wg.Wait()

	index.ClearObjects()

	require.Equal(t, 0, len(index.GetObjects()))
}
