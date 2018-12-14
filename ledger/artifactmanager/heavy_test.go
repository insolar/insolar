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

package artifactmanager

import (
	"context"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/testutils"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLedgerArtifactManager_handleHeavy(t *testing.T) {
	t.Parallel()
	ctx, db, _, cleaner := getTestData(t)
	defer cleaner()
	jetID := testutils.RandomID()

	// prepare mock
	heavysync := testutils.NewHeavySyncMock(t)
	heavysync.StartMock.Return(nil)
	heavysync.StoreMock.Set(func(ctx context.Context, jetID core.RecordID, pn core.PulseNumber, kvs []core.KV) error {
		return db.StoreKeyValues(ctx, kvs)
	})
	heavysync.StopMock.Return(nil)

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	// message hanler with mok
	mh := NewMessageHandler(db, nil)

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetStorageFunc = func(p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	mh.RecentStorageProvider = provideMock

	mh.HeavySync = heavysync

	payload := []core.KV{
		{K: []byte("ABC"), V: []byte("CDE")},
		{K: []byte("ABC"), V: []byte("CDE")},
		{K: []byte("CDE"), V: []byte("ABC")},
	}

	parcel := &message.Parcel{
		Msg: &message.HeavyPayload{
			JetID:   jetID,
			Records: payload,
		},
	}

	var err error
	_, err = mh.handleHeavyPayload(ctx, parcel)
	require.NoError(t, err)

	badgerdb := db.GetBadgerDB()
	err = badgerdb.View(func(tx *badger.Txn) error {
		for _, kv := range payload {
			item, err := tx.Get(kv.K)
			if !assert.NoError(t, err) {
				continue
			}
			value, err := item.Value()
			if !assert.NoError(t, err) {
				continue
			}
			// fmt.Println("Got key:", string(item.Key()))
			assert.Equal(t, kv.V, value)
		}
		return nil
	})
	require.NoError(t, err)
}
