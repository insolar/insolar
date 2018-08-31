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

package storagetestutils

import (
	"log"
	"os"
	"testing"

	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/badgerdb/badgertestutils"
	"github.com/insolar/insolar/ledger/storage/leveldb/leveltestutils"
)

// TmpStore returns current store implementation and cleanup function.
func TmpStore(t *testing.T) (storage.Store, func()) {
	if os.Getenv("INSOLAR_STORAGE_ENGINE") == "level" {
		log.Println("Use LevelDB implemenatation (Depricated)")
		return tmpstoreLevel(t)
	}
	return tmpstoreBadger(t)
}

func tmpstoreLevel(t *testing.T) (storage.Store, func()) {
	return leveltestutils.TmpDB(t, "")
}

func tmpstoreBadger(t *testing.T) (storage.Store, func()) {
	return badgertestutils.TmpDB(t, "")
}
