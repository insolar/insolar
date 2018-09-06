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

package storagetest

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/insolar/insolar/ledger/storage"
)

// TmpStore returns BadgerDB's store implementation and cleanup function.
//
// Creates BadgerDB in temporary directory and uses t for errors reporting.
func TmpStore(t *testing.T, dir string) (*storage.DB, func()) {
	tmpdir, err := ioutil.TempDir(dir, "bdb-test-")
	if err != nil {
		t.Fatal(err)
	}
	store, err := storage.NewStore(tmpdir, nil)
	if err != nil {
		t.Fatal(err)
	}
	return store, func() {
		closeErr := store.Close()
		rmErr := os.RemoveAll(tmpdir)
		if closeErr != nil {
			t.Error("temporary db close failed", closeErr)
		}
		if rmErr != nil {
			t.Fatal("temporary db dir cleanup failed", rmErr)
		}
	}
}
