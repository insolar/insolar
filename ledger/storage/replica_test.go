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

	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

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

func Test_SyncClientJetPulses(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	jetID := testutils.RandomJet()

	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	var expectEmpty []core.PulseNumber
	gotEmpty, err := db.GetSyncClientJetPulses(ctx, jetID)
	require.NoError(t, err)
	assert.Equal(t, expectEmpty, gotEmpty)

	expect := []core.PulseNumber{100, 500, 100500}
	err = db.SetSyncClientJetPulses(ctx, jetID, expect)
	require.NoError(t, err)

	got, err := db.GetSyncClientJetPulses(ctx, jetID)
	require.NoError(t, err)
	assert.Equal(t, expect, got)
}

func Test_GetAllSyncClientJets(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	tt := []struct {
		jetID  core.RecordID
		pulses []core.PulseNumber
	}{
		{
			jetID:  testutils.RandomJet(),
			pulses: []core.PulseNumber{100, 500, 100500},
		},
		{
			jetID:  testutils.RandomJet(),
			pulses: []core.PulseNumber{100, 500},
		},
		{
			jetID: testutils.RandomJet(),
		},
		{
			jetID:  testutils.RandomJet(),
			pulses: []core.PulseNumber{100500},
		},
	}

	for _, tCase := range tt {
		err := db.SetSyncClientJetPulses(ctx, tCase.jetID, tCase.pulses)
		require.NoError(t, err)
	}

	gotJets, err := db.GetAllNonEmptySyncClientJets(ctx)
	require.NoError(t, err)
	// fmt.Printf("%#v\n", gotJets)

	for i, tCase := range tt {
		gotPulses, ok := gotJets[tCase.jetID]
		if tCase.pulses == nil {
			assert.Falsef(t, ok, "jet should not present jetID=%v", tCase.jetID)
		} else {
			require.Truef(t, ok, "jet should  present jetID=%v", tCase.jetID)
			assert.Equalf(t, tCase.pulses, gotPulses, "pulses not found for jet number %v: %v", i, tCase.jetID)
		}
	}

	gotJets, err = db.GetAllSyncClientJets(ctx)
	require.NoError(t, err)
	// fmt.Printf("%#v\n", gotJets)

	for i, tCase := range tt {
		gotPulses, ok := gotJets[tCase.jetID]
		require.Truef(t, ok, "jet should  present jetID=%v", tCase.jetID)
		assert.Equalf(t, tCase.pulses, gotPulses, "pulses not found for jet number %v: %v", i, tCase.jetID)
	}
}
