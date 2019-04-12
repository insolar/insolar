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
	"sync"

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
	NotifyAboutPulse(ctx context.Context, pn insolar.PulseNumber)
}

type cleaner struct {
	pulsesMux sync.Mutex
	pulses    []insolar.PulseNumber

	lightChainLimint int

	jetStorage     jet.Storage
	nodeModifier   node.Modifier
	dropCleaner    drop.Cleaner
	blobCleaner    blob.Cleaner
	recCleaner     object.RecordCleaner
	indexCleaner   object.IndexCleaner
	recentProvider recentstorage.Provider
	pulseShifter   pulse.Shifter
}

func NewCleaner(
	jetStorage jet.Storage,
	nodeModifier node.Modifier,
	dropCleaner drop.Cleaner,
	blobCleaner blob.Cleaner,
	recCleaner object.RecordCleaner,
	indexCleaner object.IndexCleaner,
	recentProvider recentstorage.Provider,
	pulseShifter pulse.Shifter,
	lightChainLimint int,
) Cleaner {
	return &cleaner{
		jetStorage:       jetStorage,
		nodeModifier:     nodeModifier,
		dropCleaner:      dropCleaner,
		blobCleaner:      blobCleaner,
		recCleaner:       recCleaner,
		indexCleaner:     indexCleaner,
		recentProvider:   recentProvider,
		pulseShifter:     pulseShifter,
		lightChainLimint: lightChainLimint,
	}
}

func (c *cleaner) NotifyAboutPulse(ctx context.Context, pn insolar.PulseNumber) {
	c.pulsesMux.Lock()
	defer c.pulsesMux.Unlock()

	c.pulses = append(c.pulses, pn)

	for len(c.pulses) > c.lightChainLimint && len(c.pulses) > 0 {
		pnForClean := c.pulses[0]
		c.pulses = c.pulses[1:]
		c.cleanPulse(ctx, pnForClean)
	}
}

func (c cleaner) cleanPulse(ctx context.Context, pn insolar.PulseNumber) {
	c.nodeModifier.Delete(pn)
	c.dropCleaner.Delete(pn)
	c.blobCleaner.Delete(ctx, pn)
	c.recCleaner.Remove(ctx, pn)

	c.jetStorage.Delete(ctx, pn)

	excIdx := c.getExcludedIndexes(ctx, pn)
	c.indexCleaner.RemoveUntil(ctx, pn, excIdx)

	err := c.pulseShifter.Shift(ctx, pn)
	if err != nil {
		inslogger.FromContext(ctx).Errorf("Can't clean pulse-tracker from pulse: %s", err)
	}
}

func (c *cleaner) getExcludedIndexes(ctx context.Context, pn insolar.PulseNumber) map[insolar.ID]struct{} {
	jets := c.jetStorage.All(ctx, pn)
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
