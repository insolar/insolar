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

package storagetest

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/assert"
)

// TmpDB returns BadgerDB's storage implementation and cleanup function.
//
// Creates BadgerDB in temporary directory and uses t for errors reporting.
func TmpDB(ctx context.Context, t testing.TB, dir string) (*storage.DB, func()) {
	tmpdir, err := ioutil.TempDir(dir, "bdb-test-")
	if err != nil {
		t.Fatal(err)
	}
	db, err := storage.NewDB(configuration.Ledger{
		Storage: configuration.Storage{
			DataDirectory: tmpdir,
		},
	}, nil)
	db.PlatformCryptographyScheme = platformpolicy.NewPlatformCryptographyScheme()
	if err != nil {
		t.Fatal(err)
	}
	// Bootstrap
	err = db.Bootstrap(ctx)
	assert.NoError(t, err)

	return db, func() {
		closeErr := db.Close()
		rmErr := os.RemoveAll(tmpdir)
		if closeErr != nil {
			t.Error("temporary db close failed", closeErr)
		}
		if rmErr != nil {
			t.Fatal("temporary db dir cleanup failed", rmErr)
		}
	}
}
