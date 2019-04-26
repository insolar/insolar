/*
 *    Copyright 2019 Insolar Technologies
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

package replication

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/replication.Cleaner -o ./ -s _mock.go

// Cleaner is an interface that represents a cleaner-component
// It's supposed, that all the process of cleaning data from LME will be doing by it
type Cleaner interface {
	// NotifyAboutPulse notifies a component about a pulse
	NotifyAboutPulse(ctx context.Context, pn insolar.PulseNumber)
}

// LightCleaner is an implementation of Cleaner interface
type LightCleaner struct {
	once          sync.Once
	pulseForClean chan insolar.PulseNumber

	jetStorage   jet.Storage
	nodeModifier node.Modifier
	dropCleaner  drop.Cleaner
	blobCleaner  blob.Cleaner
	recCleaner   object.RecordCleaner
	indexCleaner object.IndexCleaner
	pulseShifter pulse.Shifter

	pulseCalculator pulse.Calculator

	lightChainLimit int
}

// NewCleaner creates a new instance of LightCleaner
func NewCleaner(
	jetStorage jet.Storage,
	nodeModifier node.Modifier,
	dropCleaner drop.Cleaner,
	blobCleaner blob.Cleaner,
	recCleaner object.RecordCleaner,
	indexCleaner object.IndexCleaner,
	pulseShifter pulse.Shifter,
	pulseCalculator pulse.Calculator,
	lightChainLimit int,
) *LightCleaner {
	return &LightCleaner{
		jetStorage:      jetStorage,
		nodeModifier:    nodeModifier,
		dropCleaner:     dropCleaner,
		blobCleaner:     blobCleaner,
		recCleaner:      recCleaner,
		indexCleaner:    indexCleaner,
		pulseShifter:    pulseShifter,
		pulseCalculator: pulseCalculator,
		lightChainLimit: lightChainLimit,
		pulseForClean:   make(chan insolar.PulseNumber),
	}
}

// NotifyAboutPulse cleans a light's data. When it's called, it tries to fetch
// pulse, which is backwards by a size of lightChainLimit. If a pulse is fetched successfully,
// all the data for it will be cleaned
func (c *LightCleaner) NotifyAboutPulse(ctx context.Context, pn insolar.PulseNumber) {
	c.once.Do(func() {
		go c.clean(ctx)
	})
	inslogger.FromContext(ctx).Debugf("[Cleaner][NotifyAboutPulse] received pulse - %v", pn)
	c.pulseForClean <- pn
}

func (c *LightCleaner) clean(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	for pn := range c.pulseForClean {
		logger.Debugf("[Cleaner][NotifyAboutPulse] start cleaning pulse - %v", pn)

		expiredPn, err := c.pulseCalculator.Backwards(ctx, pn, c.lightChainLimit)
		if err == pulse.ErrNotFound {
			logger.Errorf("[Cleaner][NotifyAboutPulse] expiredPn for pn - %v doesn't exist. limit - %v", pn, c.lightChainLimit)
			continue
		}
		if err != nil {
			panic(err)
		}

		c.cleanPulse(ctx, expiredPn.PulseNumber)
	}
}

func (c *LightCleaner) cleanPulse(ctx context.Context, pn insolar.PulseNumber) {
	inslogger.FromContext(ctx).Debugf("[Cleaner][cleanPulse] start cleaning. pn - %v", pn)
	c.nodeModifier.DeleteForPN(pn)
	c.dropCleaner.DeleteForPN(ctx, pn)
	c.blobCleaner.DeleteForPN(ctx, pn)
	c.recCleaner.DeleteForPN(ctx, pn)

	c.jetStorage.DeleteForPN(ctx, pn)
	c.indexCleaner.DeleteForPN(ctx, pn)

	err := c.pulseShifter.Shift(ctx, pn)
	if err != nil {
		inslogger.FromContext(ctx).Errorf("[Cleaner][cleanPulse] Can't clean pulse-tracker from pulse: %s", err)
	}
	inslogger.FromContext(ctx).Debugf("[Cleaner][cleanPulse] end cleaning. pn - %v", pn)
}
