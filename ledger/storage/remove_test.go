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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func Test_RemoveJetIndexesUntil_Basic(t *testing.T) {
	t.Parallel()
	removeJetIndexesUntil(t, false)
}

func Test_RemoveJetIndexesUntil_WithSkips(t *testing.T) {
	t.Parallel()
	removeJetIndexesUntil(t, true)
}

func removeJetIndexesUntil(t *testing.T, skip bool) {
	ctx := inslogger.TestContext(t)
	jetID := testutils.RandomJet()

	db, cleaner := storagetest.TmpDB(ctx, t, storagetest.DisableBootstrap())
	defer cleaner()

	pulsesCount := 10
	untilIdx := pulsesCount / 2
	var until core.PulseNumber

	pulses := []core.PulseNumber{}
	var expectLeftIDs []core.RecordID
	expectedRmCount := 0
	for i := 0; i < pulsesCount; i++ {
		pn := core.FirstPulseNumber + core.PulseNumber(i)
		if i == untilIdx {
			until = pn
			if skip {
				// skip index saving with 'until' pulse (corner case)
				continue
			}
		}
		pulses = append(pulses, pn)
		objID := testutils.RandomID()
		copy(objID[:core.PulseNumberSize], pn.Bytes())
		err := db.SetObjectIndex(ctx, jetID, &objID, &index.ObjectLifeline{
			State:       record.StateActivation,
			LatestState: &objID,
		})
		require.NoError(t, err)
		// fmt.Println("..save", objID, "on pulse", pn)
		if (pn == core.FirstPulseNumber) || (i >= untilIdx) {
			expectLeftIDs = append(expectLeftIDs, objID)
		} else {
			expectedRmCount += 1
		}
	}
	rmcount, err := db.RemoveJetIndexesUntil(ctx, jetID, until)
	require.NoError(t, err)

	// fmt.Println("expectLeftIDs:", expectLeftIDs)
	var foundIDs []core.RecordID
	db.IterateIndexIDs(ctx, jetID, func(id core.RecordID) error {
		// fmt.Println("found:", id)
		foundIDs = append(foundIDs, id)
		return nil
	})

	require.Equal(t, expectedRmCount, rmcount)
	require.Equal(t, expectLeftIDs, foundIDs)
}
