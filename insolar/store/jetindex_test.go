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

package store

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndex_Add(t *testing.T) {
	t.Parallel()

	idx := NewJetIndex()
	id := gen.ID()
	jetID := gen.JetID()
	idx.Add(id, jetID)
	assert.Equal(t, idx.storage[jetID], recordSet{id: struct{}{}})
}

func TestJetIndex_Delete(t *testing.T) {
	t.Parallel()

	idx := NewJetIndex()
	id := gen.ID()
	jetID := gen.JetID()
	idx.storage[jetID] = recordSet{}
	idx.storage[jetID][id] = struct{}{}
	idx.Delete(id, jetID)
	assert.Nil(t, idx.storage[jetID])
}

func TestJetIndex_For(t *testing.T) {
	t.Parallel()

	idx := NewJetIndex()
	id := insolar.NewID(insolar.PulseNumber(4), []byte{1})
	sID := insolar.NewID(insolar.PulseNumber(4), []byte{2})
	tID := insolar.NewID(insolar.PulseNumber(4), []byte{3})
	jetID := gen.JetID()
	idx.Add(*id, jetID)
	idx.Add(*sID, jetID)
	idx.Add(*tID, jetID)

	for i := 0; i < 100; i++ {
		id := gen.ID()
		rJetID := gen.JetID()
		if id.Pulse() != insolar.PulseNumber(4) && rJetID != jetID {
			idx.Add(id, rJetID)
		}
	}

	res := idx.For(jetID)

	require.Equal(t, 3, len(res))
	_, ok := res[*id]
	require.Equal(t, true, ok)
	_, ok = res[*sID]
	require.Equal(t, true, ok)
	_, ok = res[*tID]
	require.Equal(t, true, ok)
}
