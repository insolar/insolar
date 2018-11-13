package storage

import (
	"sync"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/require"
)

func TestNewRecentObjectsIndex(t *testing.T) {
	index := NewRecentObjectsIndex()
	require.NotNil(t, index)
	require.NotNil(t, index.updatedObjects)
	require.NotNil(t, index.fetchedObjects)
}

func Test_addToFetched(t *testing.T) {
	index := NewRecentObjectsIndex()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		index.addToFetched(core.NewRecordID(123, []byte{1}))
		wg.Done()
	}()
	go func() {
		index.addToFetched(core.NewRecordID(123, []byte{2}))
		wg.Done()
	}()
	go func() {
		index.addToFetched(core.NewRecordID(123, []byte{3}))
		wg.Done()
	}()

	wg.Wait()
	require.Equal(t, 3, len(index.fetchedObjects))
}

func Test_addToUpdated(t *testing.T) {
	index := NewRecentObjectsIndex()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		index.addToUpdated(core.NewRecordID(123, []byte{1}))
		wg.Done()
	}()
	go func() {
		index.addToUpdated(core.NewRecordID(123, []byte{2}))
		wg.Done()
	}()
	go func() {
		index.addToUpdated(core.NewRecordID(123, []byte{3}))
		wg.Done()
	}()

	wg.Wait()
	require.Equal(t, 3, len(index.updatedObjects))
}

func Test_clear(t *testing.T) {
	index := NewRecentObjectsIndex()
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		index.addToUpdated(core.NewRecordID(123, []byte{1}))
		wg.Done()
	}()
	go func() {
		index.addToUpdated(core.NewRecordID(123, []byte{2}))
		wg.Done()
	}()
	go func() {
		index.addToFetched(core.NewRecordID(123, []byte{3}))
		wg.Done()
	}()
	wg.Wait()

	index.clear()

	require.Equal(t, 0, len(index.updatedObjects))
	require.Equal(t, 0, len(index.fetchedObjects))
}
