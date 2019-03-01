package db

import (
	"testing"

	"github.com/insolar/insolar/gen"
	"github.com/stretchr/testify/assert"
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
