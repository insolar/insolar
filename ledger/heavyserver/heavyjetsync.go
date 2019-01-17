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

// NewHeavyJetSync creates HeavyJetTreeSync
func NewHeavyJetSync(db *storage.DB) *HeavyJetSync {
	return &HeavyJetSync{db: db}
}

// SyncTree updates state of the heavy's jet tree
func (hs *HeavyJetSync) SyncTree(ctx context.Context, tree jet.Tree, pulse core.PulseNumber) error {
	return hs.db.AppendJetTree(ctx, pulse, &tree)
}
