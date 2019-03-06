/*
 *    Copyright 2019 Insolar Technologies
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

package blob_test

import (
	"math/rand"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryBlob(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	blobStorage := blob.NewStorage()

	type tempBlob struct {
		id core.RecordID
		b  blob.Blob
	}

	var blobs []tempBlob

	f := fuzz.New().Funcs(func(t *tempBlob, c fuzz.Continue) {
		t.id = gen.ID()
		t.b = blob.Blob{
			Value: slice(),
			JetID: gen.JetID(),
		}
	})
	f.NumElements(5, 10).NilChance(0).Fuzz(&blobs)

	t.Run("saves correct blob-value", func(t *testing.T) {
		for _, bl := range blobs {
			err := blobStorage.Set(ctx, bl.id, bl.b)
			require.NoError(t, err)
		}

		for _, bl := range blobs {
			resBlob, err := blobStorage.Get(ctx, bl.id)
			require.NoError(t, err)

			assert.Equal(t, bl.b, resBlob)
			assert.Equal(t, bl.b.Value, resBlob.Value)
			assert.Equal(t, bl.b.JetID, resBlob.JetID)
		}
	})

	t.Run("returns override error when saving with the same id", func(t *testing.T) {
		for _, bl := range blobs {
			err := blobStorage.Set(ctx, bl.id, bl.b)
			require.Error(t, err)
			assert.Equal(t, blob.ErrOverride, err)
		}
	})

	t.Run("returns error when no blob-value for id", func(t *testing.T) {
		t.Parallel()

		for i := int32(0); i < rand.Int31n(10); i++ {
			_, err := blobStorage.Get(ctx, gen.ID())
			require.Error(t, err)
			assert.Equal(t, blob.ErrNotFound, err)
		}
	})
}

// sizedSlice generates random byte slice fixed size
func sizedSlice(size int32) (blob []byte) {
	blob = make([]byte, size)
	rand.Read(blob)
	return
}

// slice generates random byte slice with random size between 0 and 1024
func slice() []byte {
	size := rand.Int31n(1024)
	return sizedSlice(size)
}
