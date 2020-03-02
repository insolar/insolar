// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
	indexAccessor object.MemoryIndexAccessor

	pulseShifter    pulse.Shifter
	pulseCalculator pulse.Calculator

	filamentCleaner FilamentCleaner

	filamentLimit   int
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
	indexAccessor object.MemoryIndexAccessor,
	filamentCleaner FilamentCleaner,
	lightChainLimit int,
	cleanerDelay int,
	filamentLimit int,
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
		filamentLimit:   filamentLimit,
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
			logger.Panic(err)
		}
		c.cleanPulse(ctx, expiredPn.PulseNumber, pn)
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

func (c *LightCleaner) cleanPulse(ctx context.Context, cleanFrom, latest insolar.PulseNumber) {
	logger := inslogger.FromContext(ctx)

	logger.Infof("start cleaning pn:%v", cleanFrom)

	c.nodeModifier.DeleteForPN(cleanFrom)
	c.dropCleaner.DeleteForPN(ctx, cleanFrom)
	c.recCleaner.DeleteForPN(ctx, cleanFrom)
	c.jetCleaner.DeleteForPN(ctx, cleanFrom)
	c.indexCleaner.DeleteForPN(ctx, cleanFrom)

	prev, err := c.pulseCalculator.Backwards(ctx, latest, 1)
	if err == nil {
		indexes, err := c.indexAccessor.ForPulse(ctx, prev.PulseNumber)
		if err != nil && err != object.ErrIndexNotFound {
			logger.Errorf("Can't get indexes for pulse: %s", err)
		} else {
			ids := make([]insolar.ID, len(indexes))
			for i, index := range indexes {
				ids[i] = index.ObjID
			}
			c.filamentCleaner.ClearAllExcept(ids)
		}
	} else {
		logger.Error("Can't get prev pulse", err)
	}

	c.filamentCleaner.ClearIfLonger(c.filamentLimit)

	err = c.pulseShifter.Shift(ctx, cleanFrom)
	if err != nil {
		logger.Errorf("can't clean pulse-tracker from pulse: %s", err)
	}

	logger.Infof("finish cleaning pn:%v", cleanFrom)
}
