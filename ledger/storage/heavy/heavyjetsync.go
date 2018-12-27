package heavy

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/jet"
)

type HeavyJetSync interface {
	SyncTree(ctx context.Context, tree jet.Tree, pulse core.PulseNumber) error
}