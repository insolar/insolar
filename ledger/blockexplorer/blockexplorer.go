package blockexplorer

import (
	"context"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
)

const (
	getHistoryChunkSize = 10 * 1000
)

// BlockExplorerManager provides concrete API to block explorer module.
type BlockExplorerManager struct {
	db         *storage.DB
	messageBus core.MessageBus

	getHistoryChunkSize int
}

// NewArtifactManger creates new manager instance.
func NewBlockExplorer(db *storage.DB) (*BlockExplorerManager, error) {
	return &BlockExplorerManager{db: db, getHistoryChunkSize: getHistoryChunkSize}, nil
}

// Link links external components.
func (m *BlockExplorerManager) Link(components core.Components) error {
	m.messageBus = components.MessageBus

	return nil
}

// GetHistory returns history iterator.
//
// During iteration history will be fetched from remote source.
func (m *BlockExplorerManager) GetHistory(ctx context.Context, object core.RecordRef,
	pulse *core.PulseNumber) (core.RefIterator, error) {
	var err error
	defer instrument(ctx, "GetHistory").err(&err).end()
	iter, err := NewHistoryIterator(ctx, m.messageBus, object, pulse, m.getHistoryChunkSize)
	return iter, err
}
