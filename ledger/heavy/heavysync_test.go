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
	var pnum core.PulseNumber
	kvalues := []core.KV{
		{K: []byte("100"), V: []byte("500")},
	}

	// TODO: call every case in subtest

	sync := NewSync(db)
	err = sync.Start(ctx, pnum)
	require.Error(t, err, "start with zero pulse")

	err = sync.Store(ctx, pnum, kvalues)
	require.Error(t, err, "store values on non started sync")

	err = sync.Stop(ctx, pnum)
	require.Error(t, err, "stop on non started sync")

	pnum = 5
	err = sync.Start(ctx, pnum)
	require.Error(t, err, "last synced pulse is less when 'first pulse number'")

	pnum = core.FirstPulseNumber + 1
	err = sync.Start(ctx, pnum)
	require.Error(t, err, "start sync on empty store with non first pulse number")

	pnum = core.FirstPulseNumber
	err = sync.Start(ctx, pnum)
	require.NoError(t, err, "start from first pulse on empty storage")

	err = sync.Start(ctx, pnum)
	require.Error(t, err, "double start")

	pnumNext := pnum + 1
	err = sync.Start(ctx, pnumNext)
	require.Error(t, err, "start next pulse sync when previous not end")

	// stop previous
	err = sync.Stop(ctx, pnum)
	require.NoError(t, err)

	// start next
	pnumNextPlus := pnumNext + 1
	err = sync.Start(ctx, pnumNextPlus)
	require.Error(t, err, "start when previous pulses not synced")

	// prepare pulse helper
	preparepulse := func(pn core.PulseNumber) {
		pulse := core.Pulse{PulseNumber: pn}
		// fmt.Printf("Store pulse: %v\n", pulse.PulseNumber)
		err = db.AddPulse(ctx, pulse)
		require.NoError(t, err)
	}
	preparepulse(pnum)
	preparepulse(pnumNext) // should set correct next for previous pulse

	err = sync.Start(ctx, pnumNext)
	require.NoError(t, err, "start next pulse")

	err = sync.Store(ctx, pnumNextPlus, kvalues)
	require.Error(t, err, "store from other pulse at the same jet")

	err = sync.Stop(ctx, pnumNextPlus)
	require.Error(t, err, "stop from other pulse at the same jet")

	err = sync.Store(ctx, pnumNext, kvalues)
	require.NoError(t, err, "store on current range")
	err = sync.Store(ctx, pnumNext, kvalues)
	require.NoError(t, err, "store the same on current range")
	err = sync.Stop(ctx, pnumNext)
	require.NoError(t, err, "stop current range")

	preparepulse(pnumNextPlus) // should set corret next for previous pulse
	sync = NewSync(db)
	err = sync.Start(ctx, pnumNextPlus)
	require.NoError(t, err, "start next+1 range on new sync instance (checkpoint check)")
	err = sync.Store(ctx, pnumNextPlus, kvalues)
	require.NoError(t, err, "store next+1 pulse")
	err = sync.Stop(ctx, pnumNextPlus)
	require.NoError(t, err, "stop next+1 range on new sync instance")
}
