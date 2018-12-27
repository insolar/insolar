package heavyserver

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/jet"
)

type HeavyJetSync struct {
	db *storage.DB
}

func NewHeavyJetSync(db *storage.DB) *HeavyJetSync {
	return &HeavyJetSync{db: db}
}

func (hs *HeavyJetSync) SyncTree(ctx context.Context, tree jet.Tree, pulse core.PulseNumber) error {
	// jetTree, err := hs.db.GetJetTree(ctx, pulse)
	// if err != nil {
	// 	return err
	// }
	// //x, a = a[0], a[1:]
	// jetQueue := []jet
	panic("")
}
