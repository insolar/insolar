package storage

import (
	"sync"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/require"
)

func TestNewRecentObjectsIndex(t *testing.T) {
	index := NewRecentObjectsIndex(123)
	require.NotNil(t, index)
	require.NotNil(t, index.RecentObjects)
	require.Equal(t, 123, index.DefaultTTL)
}

func Test_addToFetched(t *testing.T) {
	index := NewRecentObjectsIndex(123)

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		index.AddId(core.NewRecordID(123, []byte{1}))
		wg.Done()
	}()
	go func() {
		index.AddId(core.NewRecordID(123, []byte{2}))
		wg.Done()
	}()
	go func() {
		index.AddId(core.NewRecordID(123, []byte{3}))
		wg.Done()
	}()

	wg.Wait()
	require.Equal(t, 3, len(index.RecentObjects))
}

//
// func Test_clear(t *testing.T) {
// 	index := NewRecentObjectsIndex()
// 	wg := sync.WaitGroup{}
// 	wg.Add(3)
// 	go func() {
// 		index.addId(core.NewRecordID(123, []byte{1}))
// 		wg.Done()
// 	}()
// 	go func() {
// 		index.addId(core.NewRecordID(123, []byte{2}))
// 		wg.Done()
// 	}()
// 	go func() {
// 		index.addId(core.NewRecordID(123, []byte{3}))
// 		wg.Done()
// 	}()
// 	wg.Wait()
//
// 	index.clear()
//
// 	require.Equal(t, 0, len(index.recentObjects))
// }
