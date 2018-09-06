/*
 *    Copyright 2018 INS Ecosystem
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

	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
)

func TestStore_DropWaitWrites(t *testing.T) {
	t.Parallel()
	store, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	var (
		wgStart sync.WaitGroup
		wgEnd   sync.WaitGroup
	)
	wgStart.Add(2)
	wgEnd.Add(2)
	go func() {
		store.Update(func(tx *storage.TransactionManager) error {
			wgStart.Done()
			// log.Println("start tx1")
			return nil
		})
		// log.Println("end tx1")
		wgEnd.Done()
	}()
	tx2finish := make(chan bool)
	go func() {
		store.Update(func(tx *storage.TransactionManager) error {
			wgStart.Done()
			// log.Println("start tx2")
			<-tx2finish
			return nil
		})
		// log.Println("end tx2")
		wgEnd.Done()
	}()

	// all transactions started
	dropdone := make(chan bool)
	wgStart.Wait()
	go func() {
		prevdrop := &jetdrop.JetDrop{}
		// log.Println("start SetDrop")
		_, _ = store.SetDrop(0, prevdrop)
		close(dropdone)
	}()

	txFinCh := make(chan time.Time)
	dropFinCh := make(chan time.Time)
	go func() {
		// wait all transactions are finished
		wgEnd.Wait()
		txFinCh <- time.Now()
	}()
	go func() {
		// wait drop is ready
		<-dropdone
		dropFinCh <- time.Now()
	}()

	tx2waittime := time.Millisecond * 200
	go func() {
		time.Sleep(tx2waittime)
		close(tx2finish)
	}()

	txFin := <-txFinCh
	dropFin := <-dropFinCh
	// log.Println("R: tx end t:", txFin)
	// log.Println("R: drop   t:", dropFin)

	assert.Conditionf(t, func() bool {
		return dropFin.After(txFin)
	}, "drop should happens after transaction ending")
}
