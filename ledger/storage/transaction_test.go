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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
)

type tTxStat struct {
	attempts int
	err      error
}

type tConflictResult struct {
	tx1   tTxStat
	tx2   tTxStat
	value string
}

/*
 emulates transaction conflict in two parallel transactions tx1 and tx2
 which reads and writes the same key simultaneously

 [tx1]         [tx2]
   |             |
<start>          |
 get(k)          |
 set(k)          |
   |          <start> (if conflicts<attempts)
   |            get(k)
   |            set(k)
   |            commit()
   |          <end>
 commit()
<end>

tx1 returns storage.ErrConflict error after commit in this scenario.
*/
func testconflict(t *testing.T, db *storage.DB, key []byte, conflicts int) *tConflictResult {
	tlog := getlog()
	tlog.Println("use key", key)
	newvalue := newvalueGenerator()
	seterr := db.Set(key, newvalue())
	assert.NoError(t, seterr)

	var tx1stat tTxStat

	inflightTx1 := make(chan bool)
	endTx1 := make(chan bool)
	doneTx1 := make(chan bool)
	go func() {
		v1 := newvalue()
		tx1stat.err = db.Update(func(tx *storage.TransactionManager) error {
			tx1stat.attempts++
			tlog.Printf("tx1 [%v]: start", tx1stat.attempts)
			vgot, geterr := tx.Get(key)
			if geterr != nil {
				return geterr
			}
			tlog.Printf("tx1 [%v]: got '%v'\n", tx1stat.attempts, string(vgot))

			seterr := tx.Set(key, v1)
			if seterr != nil {
				return seterr
			}
			tlog.Printf("tx1 [%v]: set '%v'\n", tx1stat.attempts, string(v1))
			inflightTx1 <- true
			<-endTx1
			return seterr
		})
		tlog.Printf("tx1 [%v]: done", tx1stat.attempts)
		close(doneTx1)
	}()

	var tx2stat tTxStat
	tx2fn := func() {
		v2 := newvalue()
		tx2stat.err = db.Update(func(tx *storage.TransactionManager) error {
			tx2stat.attempts++
			tlog.Printf("tx2 [%v]: start", tx2stat.attempts)

			vgot, geterr := tx.Get(key)
			if geterr != nil {
				return geterr
			}
			tlog.Printf("tx2 [%v]: got '%v'", tx2stat.attempts, string(vgot))

			seterr := tx.Set(key, v2)
			if seterr == nil {
				tlog.Printf("tx2 [%v]: set '%v'\n", tx2stat.attempts, string(v2))
			}
			return seterr
		})
		tlog.Printf("tx2 [%v]: done (error=%v)\n", tx2stat.attempts, tx2stat.err)
	}

TRY_LOOP:
	for {
		select {
		case <-inflightTx1:
			// tx2 makes conflict for tx1 here util specified conflicts counter is reached
			if tx2stat.attempts < conflicts {
				tx2fn()
			}
			endTx1 <- true
		case <-doneTx1:
			tlog.Println("goroutine with cycle done")
			break TRY_LOOP
		}
	}
	<-doneTx1

	vGot, err := db.Get(key)
	assert.NoError(t, err)
	return &tConflictResult{
		tx1:   tx1stat,
		tx2:   tx2stat,
		value: string(vGot),
	}
}

func TestStore_TransactionConflict(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	t.Run("no_retry", func(t *testing.T) {
		db.SetTxRetiries(0)
		res := testconflict(t, db, genuniqkey(), 1)

		assert.Error(t, res.tx1.err)
		assert.NoError(t, res.tx2.err)
		assert.Equal(t, res.tx1.err, storage.ErrConflict)
		assert.Equal(t, 1, res.tx1.attempts)
		assert.Equal(t, "v2", res.value)
	})
	t.Run("with_retry", func(t *testing.T) {
		db.SetTxRetiries(2)
		res := testconflict(t, db, genuniqkey(), 1)

		assert.NoError(t, res.tx1.err)
		assert.NoError(t, res.tx2.err)
		assert.Equal(t, 2, res.tx1.attempts)
		assert.Equal(t, "v1", res.value)
	})
}

func TestStore_TransactionRetryOver(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	tx1attemptsExpect := 3
	db.SetTxRetiries(tx1attemptsExpect - 1)
	res := testconflict(t, db, genuniqkey(), tx1attemptsExpect*2)

	assert.Error(t, res.tx1.err)
	assert.Equal(t, res.tx1.err, storage.ErrConflictRetriesOver)
	assert.Equal(t, tx1attemptsExpect, res.tx1.attempts)
	assert.Equal(t, fmt.Sprintf("v%v", tx1attemptsExpect+1), res.value)
}

var keycounter int32

func genuniqkey() []byte {
	return []byte(fmt.Sprintf("k%v", atomic.AddInt32(&keycounter, 1)))
}

type valueGen func() []byte

func newvalueGenerator() valueGen {
	var valcounter int32
	return func() []byte {
		return []byte(fmt.Sprintf("v%v", atomic.AddInt32(&valcounter, 1)-1))
	}
}

func getlog() *log.Logger {
	if os.Getenv("INSOLAR_TESTS_CONFLICTS_DEBUG") != "" {
		return log.New(os.Stderr, "", log.LstdFlags)
	}
	return log.New(ioutil.Discard, "", log.LstdFlags)
}
