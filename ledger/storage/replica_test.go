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

package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
)

func Test_ReplicatedPulse(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	jetID := testutils.RandomJet()

	db, cleaner := storagetest.TmpDB(ctx, t, storagetest.DisableBootstrap())
	defer cleaner()

	// test {Set/Get}ReplicatedPulse methods pair
	got0, err := db.GetReplicatedPulse(ctx, jetID)
	require.NoError(t, err)
	assert.Equal(t, core.PulseNumber(0), got0)

	expect := core.PulseNumber(100500)
	err = db.SetReplicatedPulse(ctx, jetID, expect)
	require.NoError(t, err)

	got, err := db.GetReplicatedPulse(ctx, jetID)
	require.NoError(t, err)
	assert.Equal(t, expect, got)

	// test {Set/Get}HeavySyncedPulse methods pair
	heavyGot0, err := db.GetHeavySyncedPulse(ctx, jetID)
	require.NoError(t, err)
	assert.Equal(t, core.PulseNumber(0), heavyGot0)

	expectHeavy := core.PulseNumber(100500)
	err = db.SetHeavySyncedPulse(ctx, jetID, expectHeavy)
	require.NoError(t, err)

	gotHeavy, err := db.GetHeavySyncedPulse(ctx, jetID)
	require.NoError(t, err)
	assert.Equal(t, expectHeavy, gotHeavy)
}
