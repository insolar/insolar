//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

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
