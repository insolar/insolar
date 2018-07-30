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

package ledger

import (
	"os"
	"testing"

	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/leveldb"
)

var storer storage.LedgerStorer

func TestMain(m *testing.M) {
	storer = levelDBInit()
	retCode := m.Run()
	os.Exit(retCode)
}

func TestLedger_LevelDB_Init(t *testing.T) {
	ledger := Ledger{
		Store: storer,
	}
	_, _ = ledger.Store.Get(record.Key{
		Hash:     []byte("test"),
		TimeSlot: 1,
	})
}

func levelDBInit() storage.LedgerStorer {
	// ledger, err := newLedger()
	store, err := leveldb.InitDB()
	if err != nil {
		panic(err)
	}
	return store
}
