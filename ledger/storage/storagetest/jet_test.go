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
	"context"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func addDropSizeToDB(ctx context.Context, t *testing.T, db *DB, jetID core.RecordID, dropSize uint64) {
	dropSizeData := jet.DropSizeData{
		JetID:    jetID,
		PulseNo:  core.FirstPulseNumber, // TODO: is pulse correct ?
		DropSize: dropSize,
	}

	cryptoServiceMock := testutils.NewCryptographyServiceMock(t)
	cryptoServiceMock.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes(nil)
		return &signature, nil
	}
	signature, err := cryptoServiceMock.Sign(dropSizeData.Bytes(ctx))
	jetDropSize := &jet.JetDropSize{
		SizeData:  dropSizeData,
		Signature: signature.Bytes(),
	}

	require.NoError(t, err)

	err = db.AddDropSize(ctx, jetDropSize)
	require.NoError(t, err)
}

func findElement(testSize uint64, dropSizes []jet.JetDropSize) bool {
	for _, ds := range dropSizes {
		if ds.SizeData.DropSize == testSize {
			return true
		}
	}

	return false
}

func TestAddDropSize(t *testing.T) {
	ctx := inslogger.TestContext(t)
	jetID := core.TODOJetID

	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	dropSizes := []uint64{100, 200, 300}

	for _, s := range dropSizes {
		addDropSizeToDB(ctx, t, db, jetID, s)
	}

	dropSizeList, err := db.GetDropSizeList(ctx)
	require.NoError(t, err)

	dropSizeArray := []jet.JetDropSize(dropSizeList)

	require.Equal(t, len(dropSizes), len(dropSizeArray))

	for _, s := range dropSizes {
		require.True(t, findElement(s, dropSizeArray))
	}

}
