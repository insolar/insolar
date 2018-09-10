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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
)

func TestStore_TransactionConflict(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	conflictK := genuniqkey()
	v0, v1, v2 := []byte("v0"), []byte("v1"), []byte("v2")
	seterr := db.Set(conflictK, v0)
	assert.NoError(t, seterr)

	testconflict := func(
		t *testing.T,
		checkfn func(tx1err error, tx2err error, value string, retries int),
	) {
		var tries = 0
		startTx1 := make(chan bool)
		endTx1 := make(chan bool)
		doneTx1 := make(chan bool)
		// create two transaction and make it conflictable
		var tx1err error
		go func() {
			iter := 0
			tx1err = db.Update(func(tx *storage.TransactionManager) error {
				tries++
				iter++
				if iter == 1 {
					<-startTx1
					// log.Println("tx1: start")
				}
				vgot, geterr := tx.Get(conflictK)
				if geterr != nil {
					return geterr
				}
				_ = vgot
				// log.Println("tx1: got", string(vgot))

				seterr := tx.Set(conflictK, v1)
				if seterr != nil {
					return seterr
				}
				// log.Println("tx1: set", string(v1))

				if iter == 1 {
					<-endTx1
				}
				return seterr
			})
			// log.Println("tx1: done")
			close(doneTx1)
		}()

		startTx2 := make(chan bool)
		endTx2 := make(chan bool)
		doneTx2 := make(chan bool)
		var tx2err error
		go func() {
			iter := 0
			tx2err = db.Update(func(tx *storage.TransactionManager) error {
				iter++
				if iter == 1 {
					<-startTx2
					// log.Println("tx2: start")
				}

				vgot, geterr := tx.Get(conflictK)
				if geterr != nil {
					return geterr
				}
				_ = vgot
				// log.Println("tx2: got", string(vgot))

				seterr := tx.Set(conflictK, v2)
				if iter == 1 {
					<-endTx2
				}
				if seterr == nil {
					// log.Println("tx2: set", string(v2))
				}
				return seterr
			})
			// log.Println("tx2: done")
			close(doneTx2)
		}()

		close(startTx1)
		time.Sleep(50 * time.Millisecond)
		close(startTx2)
		time.Sleep(50 * time.Millisecond)
		close(endTx2)
		time.Sleep(50 * time.Millisecond)
		close(endTx1)

		<-doneTx1
		<-doneTx2

		// log.Printf("tx1err: %v", tx1err)
		// log.Printf("tx2err: %v", tx2err)
		vGot, err := db.Get(conflictK)
		assert.NoError(t, err)
		// log.Println(t.Name(), "vGot:", string(vGot))

		checkfn(tx1err, tx2err, string(vGot), tries)
	}

	t.Run("no retry", func(t *testing.T) {
		testconflict(t, func(tx1err error, tx2err error, value string, retries int) {
			assert.Error(t, tx1err)
			assert.NoError(t, tx2err)
			assert.Equal(t, tx1err, storage.ErrConflict)
			assert.Equal(t, 1, retries)
			assert.Equal(t, "v2", value)
		})
	})
	t.Run("with retry", func(t *testing.T) {
		db.SetTxRetiries(2)
		testconflict(t, func(tx1err error, tx2err error, value string, retries int) {
			assert.NoError(t, tx1err)
			assert.NoError(t, tx2err)
			assert.Equal(t, 2, retries)
			assert.Equal(t, "v1", value)
		})
	})
}

func TestStore_TransactionRetryOver(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()
	tlog := log.New(ioutil.Discard, "", log.LstdFlags)
	if os.Getenv("INSOLAR_TESTS_CONFLICTS_DEBUG") != "" {
		tlog = log.New(os.Stderr, "", log.LstdFlags)
	}

	conflictK := genuniqkey()
	var valcounter int32
	newvalue := func() []byte {
		return []byte(fmt.Sprintf("%v", atomic.AddInt32(&valcounter, 1)-1))
	}
	seterr := db.Set(conflictK, newvalue())
	assert.NoError(t, seterr)

	tx1tries := 0
	tx1triesExpect := 3
	db.SetTxRetiries(tx1triesExpect - 1)
	inflightTx1 := make(chan bool)
	endTx1 := make(chan bool)
	doneTx1 := make(chan bool)
	var tx1err error
	go func() {
		v1 := newvalue()
		tx1err = db.Update(func(tx *storage.TransactionManager) error {
			tx1tries++
			tlog.Printf("tx1 [%v]: start", tx1tries)
			vgot, geterr := tx.Get(conflictK)
			if geterr != nil {
				return geterr
			}
			tlog.Printf("tx1 [%v]: got '%v'\n", tx1tries, string(vgot))

			seterr := tx.Set(conflictK, v1)
			if seterr != nil {
				return seterr
			}
			tlog.Printf("tx1 [%v]: set '%v'\n", tx1tries, string(v1))
			inflightTx1 <- true
			<-endTx1
			return seterr
		})
		tlog.Printf("tx1 [%v]: done", tx1tries)
		close(doneTx1)
	}()

	tx2fn := func() {
		iter := 0
		v2 := newvalue()
		tx2err := db.Update(func(tx *storage.TransactionManager) error {
			iter++
			tlog.Printf("tx2 [%v]: start", iter)

			vgot, geterr := tx.Get(conflictK)
			if geterr != nil {
				return geterr
			}
			tlog.Printf("tx2 [%v]: got '%v'", iter, string(vgot))

			seterr := tx.Set(conflictK, v2)
			if seterr == nil {
				tlog.Printf("tx2 [%v]: set '%v'\n", iter, string(v2))
			}
			return seterr
		})
		tlog.Printf("tx2 [%v]: done (error=%v)\n", iter, tx2err)
	}

	go func() {
		for {
			select {
			case <-inflightTx1:
				// tx1 conflicts on every iteration with tx2 results here
				tx2fn()
				endTx1 <- true
			case <-doneTx1:
				tlog.Println("goroutine with cycle done")
				return
			}
		}
	}()
	// tx1 give up
	<-doneTx1

	tlog.Printf("tx1err: %v", tx1err)
	vGot, err := db.Get(conflictK)
	assert.NoError(t, err)
	tlog.Println("result key:", string(vGot))

	assert.Error(t, tx1err)
	assert.Equal(t, tx1err, storage.ErrConflictRetriesOver)
	assert.Equal(t, tx1triesExpect, tx1tries)
}

var keycounter int32

func genuniqkey() []byte {
	keycounter++
	return []byte(fmt.Sprintf("k%v", atomic.AddInt32(&keycounter, 1)))
}
