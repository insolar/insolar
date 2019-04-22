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

	"github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
)

func TestBlobStorages(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	memStorage := NewStorageMemory()
	dbStorage := NewDB(store.NewMemoryMockDB())
	type storage interface {
		Accessor
		Modifier
	}
	storages := map[string]storage{
		"blobsStor": memStorage,
		"badger":    dbStorage,
	}

	seen := map[insolar.ID]bool{}
	newUnseenID := func() (id insolar.ID) {
		for id = gen.ID(); seen[id]; id = gen.ID() {
			seen[id] = true
		}
		return
	}

	type blobToID struct {
		id insolar.ID
		b  Blob
	}
	var blobs []blobToID
	f := fuzz.New().Funcs(func(elem *blobToID, c fuzz.Continue) {
		// IN REAL client code object.CalculateIDForBlob(os.PCS, gen.PulseNumber(), t.b)
		elem.b = Blob{
			Value: slice(),
			JetID: gen.JetID(),
		}
		elem.id = newUnseenID()
	})
	f.NumElements(5, 10).NilChance(0).Fuzz(&blobs)

	for name, s := range storages {
		t.Run(name+" saves correct blob-value", func(t *testing.T) {
			for _, bl := range blobs {
				err := s.Set(ctx, bl.id, bl.b)
				require.NoError(t, err, "set failed")
			}

			for _, bl := range blobs {
				resBlob, err := s.ForID(ctx, bl.id)
				require.NoError(t, err)

				assert.Equal(t, bl.b, resBlob)
				assert.Equal(t, bl.b.Value, resBlob.Value)
				assert.Equal(t, bl.b.JetID, resBlob.JetID)
			}
		})

		t.Run(name+" returns error when no blob-value for id", func(t *testing.T) {
			for i := 0; i < 10; i++ {
				_, err := s.ForID(ctx, newUnseenID())
				require.Error(t, err)
				assert.Equal(t, ErrNotFound, err)
			}
		})

		t.Run(name+" returns override error when saving with the same id", func(t *testing.T) {
			for _, bl := range blobs {
				err := s.Set(ctx, bl.id, bl.b)
				require.Error(t, err)
				assert.Equal(t, ErrOverride, err)
			}
		})
	}

	t.Run("compare blobsStor and storage implementations", func(t *testing.T) {
		for _, bl := range blobs {
			resMem, err := memStorage.ForID(ctx, bl.id)
			require.NoError(t, err)
			dbBlob, err := dbStorage.ForID(ctx, bl.id)
			require.NoError(t, err)
			assert.Equal(t, resMem, dbBlob, "blobsStor and persistent result should be match")
		}
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
