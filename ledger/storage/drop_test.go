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

	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/log"
)

// testing storage SetDrop lock logic
//
// 1) wait update transaction start, when start SetDrop (should lock)
// 2) transaction waits start of SetDrop call and waits 'waittime' (200ms)
// (could be unstable in really slow environments)
// 3) wait SetDrop and transaction finished
// 4) compare finish time of SetDrop and transaction
// SetDrop should happen after transaction (after 'waittime' timeout happens)
func TestStore_DropWaitWrites(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	var txFin time.Time
	var dropFin time.Time
	waittime := time.Millisecond * 200

	var wg sync.WaitGroup
	wg.Add(2)
	txstarted := make(chan bool)
	dropwaits := make(chan bool)
	go func() {
		db.Update(func(tx *storage.TransactionManager) error {
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
		log.Debugln("start SetDrop")
		close(dropwaits)
		_, _, _ = db.CreateDrop(0, []byte{})
		dropFin = time.Now()
		log.Debugln("end SetDrop")
		wg.Done()
	}()
	wg.Wait()

	log.Debugln("R: tx end t:", txFin)
	log.Debugln("R: drop   t:", dropFin)

	assert.Conditionf(t, func() bool {
		return dropFin.After(txFin)
	}, "drop should happens after transaction ending")
}
