package db

import (
	"testing"

	"github.com/insolar/insolar/gen"
)

func TestIndex_Add(t *testing.T) {
	t.Parallel()

	idx := NewJetIndex()
	id := gen.ID()
	jetID := gen.JetID()
	idx.Add(id, jetID)
}
