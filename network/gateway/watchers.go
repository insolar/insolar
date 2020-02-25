// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
)

func pulseProcessingWatchdog(ctx context.Context, gateway *Base, pulse insolar.Pulse, done chan struct{}) {
	logger := inslogger.FromContext(ctx)

	go func() {
		select {
		case <-time.After(time.Second * time.Duration(pulse.NextPulseNumber-pulse.PulseNumber)):
			gateway.FailState(ctx, fmt.Sprintf("Node stopped due to long pulse processing, pulse:%v", pulse.PulseNumber))
		case <-done:
			logger.Debug("Resetting pulse processing watchdog")
		}
	}()
}

type pulseWatchdog struct {
	ctx       context.Context
	gateway   network.Gateway
	timer     *time.Timer
	timeout   time.Duration
	stopChan  chan struct{}
	resetChan chan struct{}
	started   bool
}

func newPulseWatchdog(ctx context.Context, gateway network.Gateway, timeout time.Duration) *pulseWatchdog {
	w := &pulseWatchdog{
		ctx:       ctx,
		gateway:   gateway,
		timeout:   timeout,
		stopChan:  make(chan struct{}, 1),
		resetChan: make(chan struct{}, 1),
		started:   false,
	}

	return w
}

func (w *pulseWatchdog) start() {
	go func(w *pulseWatchdog) {
		w.timer = time.NewTimer(w.timeout)
		for {
			select {
			case <-w.resetChan:
				w.timer.Reset(w.timeout)
			case <-w.stopChan:
				w.timer.Stop()
				return
			case <-w.timer.C:
				w.timer.Stop()
				w.gateway.FailState(w.ctx, "New valid pulse timeout exceeded")
			}
		}
	}(w)
}

func (w *pulseWatchdog) Stop() {
	w.stopChan <- struct{}{}
}

func (w *pulseWatchdog) Reset() {
	if !w.started {
		w.start()
		w.started = true
	} else {
		inslogger.FromContext(w.ctx).Debug("Resetting new pulse watchdog")
		w.resetChan <- struct{}{}
	}
}
