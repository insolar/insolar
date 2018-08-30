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

package badgerdb_test

import (
	"bytes"
	"sort"
	"testing"

	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage/badgerdb"
	"github.com/insolar/insolar/ledger/storage/badgerdb/badgertestutils"
	"github.com/stretchr/testify/assert"
)

func sorthashes(hashes [][]byte) {
	sort.Slice(hashes, func(i, j int) bool {
		return bytes.Compare(hashes[i], hashes[j]) == -1
	})
}

func TestStore_SlotIterate(t *testing.T) {
	store, cleaner := badgertestutils.TmpDB(t, "")
	defer cleaner()

	var recset = []record.Record{
		&record.ClassActivateRecord{},
		&record.ObjectActivateRecord{},
	}
	pulse1 := record.PulseNum(1)
	pulse2 := record.PulseNum(2)
	pulses := []record.PulseNum{pulse1, pulse2}

	// save records set in different pulses
	for _, pulse := range pulses {
		store.SetCurrentPulse(pulse)

		for _, rec := range recset {
			ref, err := store.SetRecord(rec)
			assert.NoError(t, err)
			assert.NotNil(t, ref)
		}
	}

	// iterate over pulse1
	var iterErr error
	var allhashes1expect [][]byte
	iterErr = store.ProcessSlotHashes(pulse1, func(it badgerdb.HashIterator) error {
		for i := 1; it.Next(); i++ {
			h := it.Hash()
			allhashes1expect = append(allhashes1expect, h)
		}
		return nil
	})
	assert.NoError(t, iterErr)

	allhashes1got, err := store.GetSlotHashes(pulse1)
	assert.NoError(t, err)
	assert.Equalf(t, len(recset), len(allhashes1got), "hashes count the same as records count")
	assert.Equalf(t, allhashes1expect, allhashes1got, "all hashes the same")

	sorthashes(allhashes1expect)
	assert.Equalf(t, allhashes1expect, allhashes1got, "GetSlotHashes returns sorted hashes")

	// iterate over pulse2
	var allhashes2expect [][]byte
	iterErr = store.ProcessSlotHashes(pulse2, func(it badgerdb.HashIterator) error {
		for i := 1; it.Next(); i++ {
			h := it.Hash()
			// log.Printf("%v: got hash: %x\n", i, h)
			allhashes2expect = append(allhashes2expect, h)
		}
		return nil
	})
	assert.NoError(t, iterErr)

	allhashes2got, err := store.GetSlotHashes(pulse2)
	assert.NoError(t, err)
	assert.Equalf(t, len(recset), len(allhashes2got), "hashes count the same as records count")
	assert.Equalf(t, allhashes2expect, allhashes2got, "all hashes the same")

	sorthashes(allhashes2expect)
	assert.Equalf(t, allhashes2expect, allhashes2got, "GetSlotHashes returns sorted hashes")

	assert.NotEqualf(t, allhashes1got, allhashes2got,
		"hash sets for different pulses should not be equal")
}
