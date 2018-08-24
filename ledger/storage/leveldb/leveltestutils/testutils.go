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

// Package leveltestutils provides sharable utils for testing LevelDB ledger implementation.
package leveltestutils

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/insolar/insolar/ledger/storage/leveldb"
)

// TmpDB returns LevelDB's ledger implementation and cleanup function.
//
// Creates LevelDB in temporary directory and uses t for errors reporting.
func TmpDB(t *testing.T, dir string) (*leveldb.LevelLedger, func()) {
	tmpdir, err := ioutil.TempDir(dir, "ldb-test-")
	if err != nil {
		t.Fatal(err)
	}
	ledger, err := leveldb.InitDB(tmpdir, nil)
	if err != nil {
		t.Fatal(err)
	}
	return ledger, func() {
		closeErr := ledger.Close()
		rmErr := os.RemoveAll(tmpdir)
		if closeErr != nil {
			t.Error("temporary db close failed", closeErr)
		}
		if rmErr != nil {
			t.Fatal("temporary db dir cleanup failed", rmErr)
		}
	}
}
