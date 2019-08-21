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

package executor

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.Cleaner -o ./ -s _mock.go -g

// Cleaner is an interface that represents a cleaner-component
// It's supposed, that all the process of cleaning data from LME will be doing by it
type Cleaner interface {
	// NotifyAboutPulse notifies a component about a pulse
	NotifyAboutPulse(ctx context.Context, pn insolar.PulseNumber)

	Stop()
}

// LightCleaner is an implementation of Cleaner interface
type LightCleaner struct {
	once          sync.Once
	pulseForClean chan insolar.PulseNumber
	done          chan struct{}

	jetCleaner   jet.Cleaner
	nodeModifier node.Modifier
	dropCleaner  drop.Cleaner
	recCleaner   object.RecordCleaner

	indexCleaner  object.IndexCleaner
	indexAccessor object.IndexAccessor

	pulseShifter    pulse.Shifter
	pulseCalculator pulse.Calculator

	filamentCleaner FilamentCleaner

	lightChainLimit int
	cleanerDelay    int
}

// NewCleaner creates a new instance of LightCleaner
func NewCleaner(
	jetCleaner jet.Cleaner,
	nodeModifier node.Modifier,
	dropCleaner drop.Cleaner,
	recCleaner object.RecordCleaner,
	indexCleaner object.IndexCleaner,
	pulseShifter pulse.Shifter,
	pulseCalculator pulse.Calculator,
	indexAccessor object.IndexAccessor,
	filamentCleaner FilamentCleaner,
	lightChainLimit int,
	cleanerDelay int,
) *LightCleaner {
	return &LightCleaner{
		jetCleaner:      jetCleaner,
		nodeModifier:    nodeModifier,
		dropCleaner:     dropCleaner,
		recCleaner:      recCleaner,
		indexCleaner:    indexCleaner,
		pulseShifter:    pulseShifter,
		pulseCalculator: pulseCalculator,
		lightChainLimit: lightChainLimit,
		cleanerDelay:    cleanerDelay,
		filamentCleaner: filamentCleaner,
		indexAccessor:   indexAccessor,
		pulseForClean:   make(chan insolar.PulseNumber),
		done:            make(chan struct{}),
	}
}

// NotifyAboutPulse cleans a light's data. When it's called, it tries to fetch
// pulse, which is backwards by a size of lightChainLimit. If a pulse is fetched successfully,
// all the data for it will be cleaned
func (c *LightCleaner) NotifyAboutPulse(ctx context.Context, pn insolar.PulseNumber) {
	logger := inslogger.FromContext(ctx)
	logger.Info("cleaner got notification about new pulse :%v", pn)
	c.once.Do(func() {
		go c.clean(context.Background())
	})

	logger.Info("cleaner before pulse :%v to waiting channel", pn)
	c.pulseForClean <- pn
}

func (c *LightCleaner) Stop() {
	close(c.done)
}

func (c *LightCleaner) clean(ctx context.Context) {
	work := func(pn insolar.PulseNumber) {
		ctx, logger := inslogger.WithTraceField(ctx, utils.RandTraceID())
		logger.Infof("cleaner reads pn:%v from queue", pn)

		// A few steps back to eliminate race conditions on pulse change.
		// Message handlers don't hold locks on data. A particular case is when we check if data is beyond limit
		// and then access nodes. Between message receive and data access cleaner can remove data for the
		// pulse on lightChainLimit. This will lead to data fetch failure. We need to give handlers time to
		// finish before removing data.
		cleanFrom := c.lightChainLimit + c.cleanerDelay
		expiredPn, err := c.pulseCalculator.Backwards(ctx, pn, cleanFrom)
		if err == pulse.ErrNotFound {
			logger.Warnf("[Cleaner][NotifyAboutPulse] expiredPn for pn - %v doesn't exist. limit - %v",
				pn, c.lightChainLimit)
			return
		}
		if err != nil {
			panic(err)
		}
		c.cleanPulse(ctx, expiredPn.PulseNumber)
	}

	for {
		select {
		case pn, ok := <-c.pulseForClean:
			if !ok {
				return
			}
			work(pn)
		case <-c.done:
			inslogger.FromContext(ctx).Info("light cleaner stopped")
			return
		}
	}
}

func (c *LightCleaner) cleanPulse(ctx context.Context, pn insolar.PulseNumber) {
	logger := inslogger.FromContext(ctx)
	logger.Infof("start cleaning pn:%v", pn)

	c.nodeModifier.DeleteForPN(pn)
	c.dropCleaner.DeleteForPN(ctx, pn)
	c.recCleaner.DeleteForPN(ctx, pn)

	c.jetCleaner.DeleteForPN(ctx, pn)

	idxs := c.indexAccessor.ForPulse(ctx, pn)
	for _, idx := range idxs {
		if idx.LifelineLastUsed < pn {
			c.filamentCleaner.Clear(idx.ObjID)
		}
	}

	c.indexCleaner.DeleteForPN(ctx, pn)

	err := c.pulseShifter.Shift(ctx, pn)
	if err != nil {
		logger.Errorf("can't clean pulse-tracker from pulse: %s", err)
	}

	logger.Infof("finish cleaning pn:%v", pn)
}
