package heavyserver

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/jet"
)

// HeavyJetSync is used for sync jets on heavy
type HeavyJetSync struct {
	db *storage.DB
}

// NewHeavyJetSync creates HeavyJetSync
func NewHeavyJetSync(db *storage.DB) *HeavyJetSync {
	return &HeavyJetSync{db: db}
}

// SyncTree updates state of the heavy's jet tree
func (hs *HeavyJetSync) SyncTree(ctx context.Context, tree jet.Tree, pulse core.PulseNumber) error {
	savedTree, err := hs.db.GetJetTree(ctx, pulse)
	if err != nil {
		return err
	}

	mergedTree := savedTree.Merge(&tree)
	return hs.db.SetJetTree(ctx, pulse, mergedTree)
}
