/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

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
type ExplorerManager struct {
	db         *storage.DB
	DefaultBus core.MessageBus `inject:""`

	getHistoryChunkSize int
}

// NewArtifactManger creates new manager instance.
func NewExplorerManager(db *storage.DB) *ExplorerManager {
	return &ExplorerManager{db: db, getHistoryChunkSize: getHistoryChunkSize}
}

// GetHistory returns history iterator.
//
// During iteration history will be fetched from remote source.
func (m *ExplorerManager) GetHistory(ctx context.Context, object core.RecordRef,
	pulse *core.PulseNumber) (core.RefIterator, error) {
	var err error
	defer instrument(ctx, "GetHistory").err(&err).end()
	iter, err := NewHistoryIterator(ctx, m.bus(ctx), object, pulse, m.getHistoryChunkSize)
	return iter, err
}

func (m *ExplorerManager) bus(ctx context.Context) core.MessageBus {
	return core.MessageBusFromContext(ctx, m.DefaultBus)
}
