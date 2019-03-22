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

package blob

import (
	"math/rand"
	"testing"

	"github.com/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlobStorage_NewStorageMemory(t *testing.T) {
	t.Parallel()

	blobStorage := NewStorageMemory()
	assert.Equal(t, 0, len(blobStorage.memory))
}

func TestBlobStorage_Set(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	rawBlob := slice()
	blob := Blob{
		JetID: jetID,
		Value: rawBlob,
	}

	jetIndex := db.NewJetIndexModifierMock(t)
	jetIndex.AddMock.Expect(id, jetID)

	t.Run("saves correct blob-value", func(t *testing.T) {
		t.Parallel()

		blobStorage := &StorageMemory{
			memory:   map[insolar.ID]Blob{},
			jetIndex: jetIndex,
		}
		err := blobStorage.Set(ctx, id, blob)
		require.NoError(t, err)
		assert.Equal(t, 1, len(blobStorage.memory))
		assert.Equal(t, blob, blobStorage.memory[id])
		assert.Equal(t, rawBlob, blobStorage.memory[id].Value)
		assert.Equal(t, jetID, blobStorage.memory[id].JetID)
	})

	t.Run("returns override error when saving with the same id", func(t *testing.T) {
		t.Parallel()

		blobStorage := &StorageMemory{
			memory:   map[insolar.ID]Blob{},
			jetIndex: jetIndex,
		}
		err := blobStorage.Set(ctx, id, blob)
		require.NoError(t, err)

		err = blobStorage.Set(ctx, id, blob)
		require.Error(t, err)
		assert.Equal(t, ErrOverride, err)
	})
}

func TestBlobStorage_ForID(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	rawBlob := slice()
	blob := Blob{
		JetID: jetID,
		Value: rawBlob,
	}

	t.Run("returns correct blob-value", func(t *testing.T) {
		t.Parallel()

		blobStorage := &StorageMemory{
			memory: map[insolar.ID]Blob{},
		}
		blobStorage.memory[id] = blob

		resultBlob, err := blobStorage.ForID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, blob, resultBlob)
		assert.Equal(t, rawBlob, resultBlob.Value)
		assert.Equal(t, jetID, resultBlob.JetID)
	})

	t.Run("returns error when no blob-value for id", func(t *testing.T) {
		t.Parallel()

		blobStorage := &StorageMemory{
			memory: map[insolar.ID]Blob{},
		}
		blobStorage.memory[id] = blob

		_, err := blobStorage.ForID(ctx, gen.ID())
		require.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
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
