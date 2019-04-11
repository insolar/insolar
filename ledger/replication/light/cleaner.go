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

package light

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/node"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/ledger/storage/pulse"
)

type Cleaner interface {
	Clean(ctx context.Context, pn insolar.PulseNumber)
}

type cleaner struct {
	jetModifier    jet.Modifier
	jetAccessor    jet.Accessor
	nodeModifier   node.Modifier
	dropCleaner    drop.Cleaner
	blobCleaner    blob.Cleaner
	recCleaner     object.RecordCleaner
	indexCleaner   object.IndexCleaner
	recentProvider recentstorage.Provider
	pulseShifter   pulse.Shifter
}

func NewCleaner(
	jetModifier jet.Modifier,
	jetAccessor jet.Accessor,
	nodeModifier node.Modifier,
	dropCleaner drop.Cleaner,
	blobCleaner blob.Cleaner,
	recCleaner object.RecordCleaner,
	indexCleaner object.IndexCleaner,
	recentProvider recentstorage.Provider,
	pulseShifter pulse.Shifter,
) Cleaner {
	return &cleaner{
		jetModifier:    jetModifier,
		jetAccessor:    jetAccessor,
		nodeModifier:   nodeModifier,
		dropCleaner:    dropCleaner,
		blobCleaner:    blobCleaner,
		recCleaner:     recCleaner,
		indexCleaner:   indexCleaner,
		recentProvider: recentProvider,
		pulseShifter:   pulseShifter,
	}
}

func (c *cleaner) Clean(ctx context.Context, pn insolar.PulseNumber) {
	c.nodeModifier.Delete(pn)
	c.dropCleaner.Delete(pn)
	c.blobCleaner.Delete(ctx, pn)
	c.recCleaner.Remove(ctx, pn)

	c.jetModifier.Delete(ctx, pn)

	excIdx := c.getExcludedIndexes(ctx, pn)
	c.indexCleaner.RemoveUntil(ctx, pn, excIdx)

	err := c.pulseShifter.Shift(ctx, pn)
	if err != nil {
		inslogger.FromContext(ctx).Errorf("Can't clean pulse-tracker from pulse: %s", err)
	}
}

func (c *cleaner) getExcludedIndexes(ctx context.Context, pn insolar.PulseNumber) map[insolar.ID]struct{} {
	jets := c.jetAccessor.All(ctx, pn)
	res := make(map[insolar.ID]struct{})
	for _, j := range jets {
		storage := c.recentProvider.GetIndexStorage(ctx, insolar.ID(j))
		ids := storage.GetObjects()
		for id, ttl := range ids {
			if id.Pulse() > pn {
				continue
			}
			if ttl > 0 {
				res[id] = struct{}{}
			}
		}
	}
	return res
}
