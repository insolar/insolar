//
// Copyright 2019 Insolar Technologies GmbH
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
//

package record_test

import (
	"crypto/sha256"
	"math/rand"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestHashMaterial(t *testing.T) {
	t.Parallel()

	t.Run("check error for nil virtual record in material record", func(t *testing.T) {
		t.Parallel()

		rec := getMaterialRecord()
		h := sha256.New()
		rec.Virtual = nil
		_, err := record.HashMaterial(h, rec)
		require.Error(t, err)
	})

	t.Run("correct consistent hash for material record", func(t *testing.T) {
		t.Parallel()

		rec := getMaterialRecord()
		h := sha256.New()
		hash1, err := record.HashMaterial(h, rec)
		require.NoError(t, err)

		h = sha256.New()
		hash2, err := record.HashMaterial(h, rec)
		require.NoError(t, err)
		assert.Equal(t, hash1, hash2)
	})

	t.Run("same hash for material record with and without signature", func(t *testing.T) {
		t.Parallel()

		rec := getMaterialRecord()
		h := sha256.New()
		hash1, err := record.HashMaterial(h, rec)
		require.NoError(t, err)

		h = sha256.New()
		rec.Signature = slice()
		hash2, err := record.HashMaterial(h, rec)
		require.NoError(t, err)
		assert.Equal(t, hash1, hash2)
	})

	t.Run("different hash for changed material record", func(t *testing.T) {
		t.Parallel()

		rec := getMaterialRecord()
		h := sha256.New()
		hashBefore, err := record.HashMaterial(h, rec)
		require.NoError(t, err)

		rec.JetID = *insolar.NewJetID(uint8(rand.Int()), gen.ID().Bytes())
		h = sha256.New()
		hashAfter, err := record.HashMaterial(h, rec)
		require.NoError(t, err)
		assert.NotEqual(t, hashBefore, hashAfter)
	})

	t.Run("different hashes for different material records", func(t *testing.T) {
		t.Parallel()

		recFoo := getMaterialRecord()
		h := sha256.New()
		hashFoo, err := record.HashMaterial(h, recFoo)
		require.NoError(t, err)

		recBar := getMaterialRecord()
		h = sha256.New()
		hashBar, err := record.HashMaterial(h, recBar)
		require.NoError(t, err)

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

// getMaterialRecord generates random Material record
func getMaterialRecord() record.Material {
	virtRec := getVirtualRecord()

	materialRecord := record.Material{
		Virtual: &virtRec,
		JetID:   *insolar.NewJetID(uint8(rand.Int()), gen.ID().Bytes()),
	}

	return materialRecord
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
