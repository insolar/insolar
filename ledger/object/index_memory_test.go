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

package object

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

func TestInMemoryIndex_SetIndex(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	objID := gen.ID()
	lflID := gen.ID()
	jetID := gen.JetID()
	buck := FilamentIndex{
		ObjID: objID,
		Lifeline: Lifeline{
			LatestState: &lflID,
			JetID:       jetID,
			Delegates:   []LifelineDelegate{},
		},
	}

	t.Run("saves correct bucket", func(t *testing.T) {
		pn := gen.PulseNumber()
		index := NewIndexStorageMemory()

		err := index.SetIndex(ctx, pn, buck)
		require.NoError(t, err)

		savedBuck := index.buckets[pn][objID]
		require.NotNil(t, savedBuck)

		buckBuf, _ := buck.Marshal()
		savedBuckBuf, _ := savedBuck.Marshal()

		require.Equal(t, buckBuf, savedBuckBuf)
	})

	t.Run("re-save works fine", func(t *testing.T) {
		pn := gen.PulseNumber()
		index := NewIndexStorageMemory()

		err := index.SetIndex(ctx, pn, buck)
		require.NoError(t, err)

		sLlflID := gen.ID()
		sJetID := gen.JetID()
		sBuck := FilamentIndex{
			ObjID: objID,
			Lifeline: Lifeline{
				LatestState: &sLlflID,
				JetID:       sJetID,
				Delegates:   []LifelineDelegate{},
			},
		}

		err = index.SetIndex(ctx, pn, sBuck)
		require.NoError(t, err)

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

	index.buckets[fPn] = map[insolar.ID]*FilamentIndex{}
	index.buckets[sPn] = map[insolar.ID]*FilamentIndex{}
	index.buckets[tPn] = map[insolar.ID]*FilamentIndex{}

	index.DeleteForPN(ctx, sPn)

	_, ok := index.buckets[fPn]
	require.Equal(t, true, ok)
	_, ok = index.buckets[sPn]
	require.Equal(t, false, ok)
	_, ok = index.buckets[tPn]
	require.Equal(t, true, ok)
}
