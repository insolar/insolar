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
	index := NewRecentObjectsIndex(123)
	require.NotNil(t, index)
	require.NotNil(t, index.recentObjects)
	require.Equal(t, 123, index.DefaultTTL)
}

func Test_AddId(t *testing.T) {
	index := NewRecentObjectsIndex(123)

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		index.AddID(core.NewRecordID(123, []byte{1}))
		wg.Done()
	}()
	go func() {
		index.AddID(core.NewRecordID(123, []byte{2}))
		wg.Done()
	}()
	go func() {
		index.AddID(core.NewRecordID(123, []byte{3}))
		wg.Done()
	}()

	wg.Wait()
	require.Equal(t, 3, len(index.GetObjects()))
}

func Test_Clear(t *testing.T) {
	index := NewRecentObjectsIndex(123)
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		index.AddID(core.NewRecordID(123, []byte{1}))
		wg.Done()
	}()
	go func() {
		index.AddID(core.NewRecordID(123, []byte{2}))
		wg.Done()
	}()
	go func() {
		index.AddID(core.NewRecordID(123, []byte{3}))
		wg.Done()
	}()
	wg.Wait()

	index.ClearObjects()

	require.Equal(t, 0, len(index.GetObjects()))
}
