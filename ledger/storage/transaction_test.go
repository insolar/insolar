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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
)

/*
check lock on select for update in 2 parallel transactions tx1 and tx2
which try reads and writes the same key simultaneously

  tx1                    tx2
   |                      |
<start>                 <start>
 get(k), for_update=T      |
 set(k)
   |----- proceed -------->|
 ..sleep..               get(k), for_update=T/F
 commit()                set(k)
  <end>                  commit()
                        <end>
*/

func TestStore_Transaction_LockOnUpdate(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	jetID := testutils.RandomID()

	objid := core.NewRecordID(100500, nil)
	idxid := core.NewRecordID(0, nil)
	objvalue0 := &index.ObjectLifeline{
		LatestState: objid,
	}
	db.SetObjectIndex(ctx, jetID, idxid, objvalue0)

	lockfn := func(t *testing.T, withlock bool) *index.ObjectLifeline {
		started2 := make(chan bool)
		proceed2 := make(chan bool)
		var wg sync.WaitGroup
		var tx1err error
		var tx2err error
		wg.Add(1)
		go func() {
			tx1err = db.Update(ctx, func(tx *storage.TransactionManager) error {
				// log.Debugf("tx1: start")
				<-started2
				// log.Debug("tx1: GetObjectIndex before")
				idxlife, geterr := tx.GetObjectIndex(ctx, jetID, idxid, true)
				// log.Debug("tx1: GetObjectIndex after")
				if geterr != nil {
					return geterr
				}

				seterr := tx.SetObjectIndex(ctx, jetID, idxid, idxlife)
				if seterr != nil {
					return seterr
				}
				// log.Debugf("tx1: set %+v\n", idxlife)
				close(proceed2)
				time.Sleep(100 * time.Millisecond)
				return seterr
			})
			// log.Debugf("tx1: finished")
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			tx2err = db.Update(ctx, func(tx *storage.TransactionManager) error {
				close(started2)
				// log.Debug("tx2: start")
				<-proceed2
				// log.Debug("tx2: GetObjectIndex before")
				idxlife, geterr := tx.GetObjectIndex(ctx, jetID, idxid, withlock)
				// log.Debug("tx2: GetObjectIndex after")
				if geterr != nil {
					return geterr
				}

				seterr := tx.SetObjectIndex(ctx, jetID, idxid, idxlife)
				if seterr != nil {
					return seterr
				}
				// log.Debugf("tx2: set %+v\n", idxlife)
				return seterr
			})
			// log.Debugf("tx2: finished")
			wg.Done()
		}()
		wg.Wait()

		assert.NoError(t, tx1err)
		assert.NoError(t, tx2err)
		idxlife, geterr := db.GetObjectIndex(ctx, jetID, idxid, false)
		assert.NoError(t, geterr)
		// log.Debugf("withlock=%v) result: got %+v", withlock, idxlife)

		// cleanup AmendRefs
		assert.NoError(t, db.SetObjectIndex(ctx, jetID, idxid, objvalue0))
		return idxlife
	}
	t.Run("with lock", func(t *testing.T) {
		idxlife := lockfn(t, true)
		assert.Equal(t, objid, idxlife.LatestState)
	})
	t.Run("no lock", func(t *testing.T) {
		idxlife := lockfn(t, false)
		assert.Equal(t, objid, idxlife.LatestState)
	})
}
