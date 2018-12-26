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

package heavyserver

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestHeavy_SyncBasic(t *testing.T) {
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	var err error
	var pnum core.PulseNumber
	kvalues := []core.KV{
		{K: []byte("100"), V: []byte("500")},
	}

	// TODO: call every case in subtest
	jetID := testutils.RandomID()

	sync := NewSync(db)
	err = sync.Start(ctx, jetID, pnum)
	require.Error(t, err, "start with zero pulse")

	err = sync.Store(ctx, jetID, pnum, kvalues)
	require.Error(t, err, "store values on non started sync")

	err = sync.Stop(ctx, jetID, pnum)
	require.Error(t, err, "stop on non started sync")

	pnum = 5
	err = sync.Start(ctx, jetID, pnum)
	require.Error(t, err, "last synced pulse is less when 'first pulse number'")

	pnum = core.FirstPulseNumber
	err = sync.Start(ctx, jetID, pnum)
	require.Error(t, err, "start from first pulse on empty storage")

	pnum = core.FirstPulseNumber + 1
	err = sync.Start(ctx, jetID, pnum)
	require.NoError(t, err, "start sync on empty heavy jet with non first pulse number")

	err = sync.Start(ctx, jetID, pnum)
	require.Error(t, err, "double start")

	pnumNext := pnum + 1
	err = sync.Start(ctx, jetID, pnumNext)
	require.Error(t, err, "start next pulse sync when previous not end")

	// stop previous
	err = sync.Stop(ctx, jetID, pnum)
	require.NoError(t, err)

	// start next
	pnumNextPlus := pnumNext + 1
	err = sync.Start(ctx, jetID, pnumNextPlus)
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

	err = sync.Start(ctx, jetID, pnumNext)
	require.NoError(t, err, "start next pulse")

	err = sync.Store(ctx, jetID, pnumNextPlus, kvalues)
	require.Error(t, err, "store from other pulse at the same jet")

	err = sync.Stop(ctx, jetID, pnumNextPlus)
	require.Error(t, err, "stop from other pulse at the same jet")

	err = sync.Store(ctx, jetID, pnumNext, kvalues)
	require.NoError(t, err, "store on current range")
	err = sync.Store(ctx, jetID, pnumNext, kvalues)
	require.NoError(t, err, "store the same on current range")
	err = sync.Stop(ctx, jetID, pnumNext)
	require.NoError(t, err, "stop current range")

	preparepulse(pnumNextPlus) // should set corret next for previous pulse
	sync = NewSync(db)
	err = sync.Start(ctx, jetID, pnumNextPlus)
	require.NoError(t, err, "start next+1 range on new sync instance (checkpoint check)")
	err = sync.Store(ctx, jetID, pnumNextPlus, kvalues)
	require.NoError(t, err, "store next+1 pulse")
	err = sync.Stop(ctx, jetID, pnumNextPlus)
	require.NoError(t, err, "stop next+1 range on new sync instance")
}

func TestHeavy_SyncByJet(t *testing.T) {
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	var err error
	var pnum core.PulseNumber
	kvalues1 := []core.KV{
		{K: []byte("1_11"), V: []byte("1_12")},
	}
	kvalues2 := []core.KV{
		{K: []byte("2_21"), V: []byte("2_22")},
	}

	// TODO: call every case in subtest
	jetID1 := testutils.RandomID()
	jetID2 := jetID1
	// flip first bit of jetID2 for different prefix
	jetID2[0] ^= jetID2[0]

	// prepare pulse helper
	preparepulse := func(pn core.PulseNumber) {
		pulse := core.Pulse{PulseNumber: pn}
		// fmt.Printf("Store pulse: %v\n", pulse.PulseNumber)
		err = db.AddPulse(ctx, pulse)
		require.NoError(t, err)
	}

	sync := NewSync(db)

	pnum = core.FirstPulseNumber + 1
	pnumNext := pnum + 1
	preparepulse(pnum)
	preparepulse(pnumNext) // should set correct next for previous pulse

	err = sync.Start(ctx, jetID1, core.FirstPulseNumber)
	require.Error(t, err)

	err = sync.Start(ctx, jetID1, pnum)
	require.NoError(t, err, "start from first+1 pulse on empty storage, jet1")

	err = sync.Start(ctx, jetID2, pnum)
	require.NoError(t, err, "start from first+1 pulse on empty storage, jet2")

	err = sync.Store(ctx, jetID2, pnum, kvalues2)
	require.NoError(t, err, "store jet2 pulse")

	err = sync.Store(ctx, jetID1, pnum, kvalues1)
	require.NoError(t, err, "store jet1 pulse")

	// stop previous
	err = sync.Stop(ctx, jetID1, pnum)
	err = sync.Stop(ctx, jetID2, pnum)
	require.NoError(t, err)
}
