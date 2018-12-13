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

package storagetest

import (
	"context"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func addDropSizeToDB(ctx context.Context, t *testing.T, db *storage.DB, jetID core.RecordID, dropSize uint64) {
	dropSizeData := &jet.DropSize{
		JetID:    jetID,
		PulseNo:  core.FirstPulseNumber,
		DropSize: dropSize,
	}

	cryptoServiceMock := testutils.NewCryptographyServiceMock(t)
	cryptoServiceMock.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes(nil)
		return &signature, nil
	}

	hasher := testutils.NewPlatformCryptographyScheme().IntegrityHasher()
	_, err := dropSizeData.WriteHashData(hasher)
	require.NoError(t, err)

	signature, err := cryptoServiceMock.Sign(hasher.Sum(nil))
	require.NoError(t, err)

	dropSizeData.Signature = signature.Bytes()

	err = db.AddDropSize(ctx, dropSizeData)
	require.NoError(t, err)
}

func findSize(testSize uint64, dropSizes []jet.DropSize) bool {
	for _, ds := range dropSizes {
		if ds.DropSize == testSize {
			return true
		}
	}

	return false
}

func TestAddAndGetDropSize(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jetID := core.TODOJetID

	db, cleaner := TmpDB(ctx, t)
	defer cleaner()

	dropSizes := []uint64{100, 200, 300, 400}

	for _, s := range dropSizes {
		addDropSizeToDB(ctx, t, db, jetID, s)
	}

	dropSizeHistory, err := db.GetDropSizeHistory(ctx)
	require.NoError(t, err)

	dropSizeArray := []jet.DropSize(dropSizeHistory)

	require.Equal(t, len(dropSizes), len(dropSizeArray))

	for _, s := range dropSizes {
		require.True(t, findSize(s, dropSizeArray))
	}
}

func TestAddDropSizeAndIncreaseLimit(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jetID := core.TODOJetID

	db, cleaner := TmpDB(ctx, t)
	defer cleaner()

	numElements := db.GetJetSizesHistoryDepth() * 2

	for i := 0; i <= numElements; i++ {
		addDropSizeToDB(ctx, t, db, jetID, uint64(i))
	}

	dropSizeHistory, err := db.GetDropSizeHistory(ctx)
	require.NoError(t, err)

	dropSizeArray := []jet.DropSize(dropSizeHistory)
	require.Equal(t, db.GetJetSizesHistoryDepth(), len(dropSizeArray))

	for i := numElements; i > (numElements - db.GetJetSizesHistoryDepth()); i-- {
		require.True(t, findSize(uint64(i), dropSizeArray), "Couldn't find %d", i)
	}
}
