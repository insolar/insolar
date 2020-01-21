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

package object

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

func TestInMemoryIndex_SetIndex(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	objID := gen.ID()
	lflID := gen.ID()
	buck := record.Index{
		ObjID: objID,
		Lifeline: record.Lifeline{
			LatestState: &lflID,
		},
	}

	t.Run("saves correct bucket", func(t *testing.T) {
		pn := gen.PulseNumber()
		index := NewIndexStorageMemory()

		index.Set(ctx, pn, buck)

		savedBuck := index.buckets[pn][objID]
		require.NotNil(t, savedBuck)

		buckBuf, _ := buck.Marshal()
		savedBuckBuf, _ := savedBuck.Marshal()

		require.Equal(t, buckBuf, savedBuckBuf)
	})

	t.Run("re-save works fine", func(t *testing.T) {
		pn := gen.PulseNumber()
		index := NewIndexStorageMemory()

		index.Set(ctx, pn, buck)

		sLlflID := insolar.NewID(lflID.Pulse()+1, lflID.Hash())
		sBuck := record.Index{
			ObjID: objID,
			Lifeline: record.Lifeline{
				LatestState: sLlflID,
			},
		}

		index.Set(ctx, pn, sBuck)

		savedBuck := index.buckets[pn][objID]
		require.NotNil(t, savedBuck)

		sBuckBuf, _ := sBuck.Marshal()
		savedBuckBuf, _ := savedBuck.Marshal()

		require.Equal(t, sBuckBuf, savedBuckBuf)
	})
}

func TestNewInMemoryIndex_DeleteForPN(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	fPn := gen.PulseNumber()
	sPn := fPn + 1
	tPn := sPn + 1

	index := NewIndexStorageMemory()

	index.buckets[fPn] = map[insolar.ID]*record.Index{}
	index.buckets[sPn] = map[insolar.ID]*record.Index{}
	index.buckets[tPn] = map[insolar.ID]*record.Index{}

	index.DeleteForPN(ctx, sPn)

	_, ok := index.buckets[fPn]
	require.Equal(t, true, ok)
	_, ok = index.buckets[sPn]
	require.Equal(t, false, ok)
	_, ok = index.buckets[tPn]
	require.Equal(t, true, ok)
}
