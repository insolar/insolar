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
