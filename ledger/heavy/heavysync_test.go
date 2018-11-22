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

package heavy

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/storagetest"
)

func TestHeavy_Sync(t *testing.T) {
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	var err error
	prange := core.PulseRange{}
	kvalues := []core.KV{
		{K: []byte("100"), V: []byte("500")},
	}

	// TODO: call every case in subtest

	sync := NewSync(db)
	err = sync.Start(ctx, prange)
	require.Error(t, err, "start with zero range")

	err = sync.Store(ctx, prange, kvalues)
	require.Error(t, err, "store values on non started sync")

	err = sync.Stop(ctx, prange)
	require.Error(t, err, "stop on non started sync")

	prange.Begin = 5
	prange.End = 6
	err = sync.Start(ctx, prange)
	require.Error(t, err, "last synced pulse is less when first pulse number")

	prange.Begin = core.FirstPulseNumber + 1
	prange.End = core.FirstPulseNumber + 2
	err = sync.Start(ctx, prange)
	require.Error(t, err, "start sync on empty store with non first pulse number")

	prange.Begin = core.FirstPulseNumber
	prange.End = core.FirstPulseNumber
	err = sync.Start(ctx, prange)
	require.Error(t, err, "case Begin<=End range")

	prange.End = prange.Begin + 1
	err = sync.Start(ctx, prange)
	require.NoError(t, err, "start from first pulse on empty storage")

	err = sync.Start(ctx, prange)
	require.Error(t, err, "double start")

	prangeNext := prange
	prangeNext.Begin++
	prangeNext.End++
	err = sync.Start(ctx, prangeNext)
	require.Error(t, err, "start next pulse when other sync already run")

	// stop previous
	err = sync.Stop(ctx, prange)

	// start next
	prangeNextPlus := prangeNext
	prangeNextPlus.Begin++
	prangeNextPlus.End++
	err = sync.Start(ctx, prangeNextPlus)
	require.Error(t, err, "start when previous pulses not synced")

	err = sync.Start(ctx, prangeNext)
	require.NoError(t, err, "start next pulse")

	err = sync.Store(ctx, prangeNextPlus, kvalues)
	require.Error(t, err, "store from other pulse at the same jet")

	err = sync.Stop(ctx, prangeNextPlus)
	require.Error(t, err, "stop from other pulse at the same jet")

	err = sync.Store(ctx, prangeNext, kvalues)
	require.NoError(t, err, "store on current range")
	err = sync.Store(ctx, prangeNext, kvalues)
	require.NoError(t, err, "store the same on current range")
	err = sync.Stop(ctx, prangeNext)
	require.NoError(t, err, "stop current range")

	sync = NewSync(db)
	err = sync.Start(ctx, prangeNextPlus)
	require.NoError(t, err, "start next+1 range on new sync instance (checkpoint check)")
	err = sync.Store(ctx, prangeNextPlus, kvalues)
	require.NoError(t, err, "store next+1 pulse")
	err = sync.Stop(ctx, prangeNextPlus)
	require.NoError(t, err, "stop next+1 range on new sync instance")
}
