// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package record_test

import (
	"crypto/sha256"
	"math/rand"
	"testing"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/stretchr/testify/assert"
)

func TestHashVirtual(t *testing.T) {
	t.Parallel()

	t.Run("check consistent hash for virtual record", func(t *testing.T) {
		t.Parallel()

		rec := getVirtualRecord()
		h := sha256.New()
		hash1 := record.HashVirtual(h, rec)

		h = sha256.New()
		hash2 := record.HashVirtual(h, rec)
		assert.Equal(t, hash1, hash2)
	})

	t.Run("same hash for virtual record with and without signature", func(t *testing.T) {
		t.Parallel()

		rec := getVirtualRecord()
		h := sha256.New()
		hash1 := record.HashVirtual(h, rec)

		h = sha256.New()
		rec.Signature = slice()
		hash2 := record.HashVirtual(h, rec)
		assert.Equal(t, hash1, hash2)
	})

	t.Run("different hash for changed virtual record", func(t *testing.T) {
		t.Parallel()

		rec := getVirtualRecord()
		h := sha256.New()
		hashBefore := record.HashVirtual(h, rec)

		rec.Union = &record.Virtual_IncomingRequest{
			IncomingRequest: &record.IncomingRequest{},
		}
		h = sha256.New()
		hashAfter := record.HashVirtual(h, rec)
		assert.NotEqual(t, hashBefore, hashAfter)
	})

	t.Run("different hashes for different virtual records", func(t *testing.T) {
		t.Parallel()

		recFoo := getVirtualRecord()
		h := sha256.New()
		hashFoo := record.HashVirtual(h, recFoo)

		recBar := getVirtualRecord()
		h = sha256.New()
		hashBar := record.HashVirtual(h, recBar)

		assert.NotEqual(t, hashFoo, hashBar)
	})
}

// getVirtualRecord generates random Virtual record
func getVirtualRecord() record.Virtual {
	var requestRecord record.IncomingRequest

	obj := gen.Reference()
	requestRecord.Object = &obj

	virtualRecord := record.Virtual{
		Union: &record.Virtual_IncomingRequest{
			IncomingRequest: &requestRecord,
		},
	}

	return virtualRecord
}

// sizedSlice generates random byte slice fixed size.
func sizedSlice(size int32) (blob []byte) {
	blob = make([]byte, size)
	rand.Read(blob)
	return
}

// slice generates random byte slice with random size between 0 and 1024.
func slice() []byte {
	size := rand.Int31n(1024)
	return sizedSlice(size)
}
