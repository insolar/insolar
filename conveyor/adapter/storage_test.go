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

package adapter

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitializer_NoRegistered(t *testing.T) {
	storage := NewStorage()

	require.Empty(t, storage.GetRegisteredAdapters())
}

func TestInitializer_GetAdapterWhileNoRegistered(t *testing.T) {
	storage := NewStorage()

	require.Nil(t, storage.GetAdapterByID(444))
}

func TestInitializer_RegisterAndGet(t *testing.T) {
	storage := NewStorage()

	for i := 0; i < 20; i++ {
		testAdapterID := uint32(i * i)
		sinkMock := NewTaskSinkMock(t)
		sinkMock.GetAdapterIDFunc = func() (r uint32) {
			return testAdapterID
		}

		storage.Register(sinkMock)
		require.Equal(t, sinkMock, storage.GetAdapterByID(testAdapterID))
	}
}

func TestInitializer_RegisterDuplicatingID(t *testing.T) {
	storage := NewStorage()

	testAdapterID := uint32(142)
	sinkMock := NewTaskSinkMock(t)
	sinkMock.GetAdapterIDFunc = func() (r uint32) {
		return testAdapterID
	}

	storage.Register(sinkMock)
	require.PanicsWithValue(t, "[ StorageManager.Register ] adapter ID 'ID(142)' already exists",
		func() {
			storage.Register(sinkMock)
		})
}

func TestInitializer_GetRegisteredAdapters(t *testing.T) {
	storage := NewStorage()

	numRegistered := 100

	for i := 0; i < numRegistered; i++ {
		testAdapterID := uint32(i)
		sinkMock := NewTaskSinkMock(t)
		sinkMock.GetAdapterIDFunc = func() (r uint32) {
			return testAdapterID
		}

		storage.Register(sinkMock)
	}

	registered := storage.GetRegisteredAdapters()
	require.Len(t, registered, numRegistered)

	// we need sort here, since adapters are stored in storage in map
	sort.Slice(registered, func(i, j int) bool {
		left := registered[i].(TaskSink).GetAdapterID()
		right := registered[j].(TaskSink).GetAdapterID()
		return left < right
	})

	for i := 0; i < numRegistered; i++ {
		require.Equal(t, uint32(i), registered[i].(TaskSink).GetAdapterID())
	}
}
