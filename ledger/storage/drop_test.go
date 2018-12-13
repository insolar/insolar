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
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/log"
)

// testing storage CreateDrop lock logic
//
// 1) wait update transaction start, when start CreateDrop (should lock)
// 2) transaction waits start of CreateDrop call and waits 'waittime' (200ms)
// (could be unstable in really slow environments)
// 3) wait CreateDrop and transaction finished
// 4) compare finish time of CreateDrop and transaction
// CreateDrop should happen after transaction (after 'waittime' timeout happens)
func TestStore_DropWaitWrites(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	var txFin time.Time
	var dropFin time.Time
	waittime := time.Millisecond * 200

	var wg sync.WaitGroup
	wg.Add(2)
	txstarted := make(chan bool)
	dropwaits := make(chan bool)
	go func() {
		db.Update(ctx, func(tx *storage.TransactionManager) error {
			log.Debugln("start tx")
			close(txstarted)
			<-dropwaits
			time.Sleep(waittime)
			return nil
		})
		txFin = time.Now()
		log.Debugln("end tx")
		wg.Done()
	}()

	go func() {
		<-txstarted
		log.Debugln("start CreateDrop")
		close(dropwaits)
		_, _, dropSize, droperr := db.CreateDrop(ctx, core.TODOJetID, 0, []byte{})
		if droperr != nil {
			panic(droperr)
		}
		require.NotEqual(t, 0, dropSize)
		dropFin = time.Now()
		log.Debugln("end CreateDrop")
		wg.Done()
	}()
	wg.Wait()

	log.Debugln("R: tx end t:", txFin)
	log.Debugln("R: drop   t:", dropFin)

	assert.Conditionf(t, func() bool {
		return dropFin.After(txFin)
	}, "drop should happens after transaction ending")
}
