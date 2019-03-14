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

package index_test

import (
	"math/rand"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryIndex(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	indexStorage := index.NewStorageMemory()

	type tempIndex struct {
		id  core.RecordID
		idx index.ObjectLifeline
	}

	var indices []tempIndex

	f := fuzz.New().Funcs(func(t *tempIndex, c fuzz.Continue) {
		t.id = gen.ID()
		ls := gen.ID()
		pn := gen.PulseNumber()
		t.idx = index.ObjectLifeline{
			LatestState:  &ls,
			LatestUpdate: pn,
			JetID:        gen.JetID(),
		}
	})
	f.NumElements(5, 10).NilChance(0).Fuzz(&indices)

	t.Run("saves correct index-value", func(t *testing.T) {
		for _, i := range indices {
			err := indexStorage.Set(ctx, i.id, i.idx)
			require.NoError(t, err)
		}

		for _, i := range indices {
			resIndex, err := indexStorage.ForID(ctx, i.id)
			require.NoError(t, err)

			assert.Equal(t, i.idx, resIndex)
			assert.Equal(t, i.idx.JetID, resIndex.JetID)
			assert.Equal(t, i.idx.LatestState, resIndex.LatestState)
			assert.Equal(t, i.idx.LatestUpdate, resIndex.LatestUpdate)
		}
	})

	t.Run("returns error when no index-value for id", func(t *testing.T) {
		t.Parallel()

		for i := int32(0); i < rand.Int31n(10); i++ {
			_, err := indexStorage.ForID(ctx, gen.ID())
			require.Error(t, err)
			assert.Equal(t, index.ErrNotFound, err)
		}
	})

	t.Run("override indices is ok", func(t *testing.T) {
		t.Parallel()

		indexStorage := index.NewStorageMemory()
		for _, i := range indices {
			err := indexStorage.Set(ctx, i.id, i.idx)
			require.NoError(t, err)
		}

		for _, i := range indices {
			err := indexStorage.Set(ctx, i.id, i.idx)
			assert.NoError(t, err)
		}
	})
}
