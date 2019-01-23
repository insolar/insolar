/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package fakepulsar

import (
	"context"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/pulsar/entropygenerator"
)

// Fakepulsar needed when the network starts and can't receive a real pulse.

// onPulse is a callbaback for pulse recv.
type callbackOnPulse func(ctx context.Context, pulse core.Pulse)

// FakePulsar is a struct which uses at void network state.
type FakePulsar struct {
	onPulse   callbackOnPulse
	stop      chan bool
	timeoutMs int32 // ms
	running   bool
}

// NewFakePulsar creates and returns a new FakePulsar.
func NewFakePulsar(callback callbackOnPulse, timeoutMs int32) *FakePulsar {
	return &FakePulsar{
		onPulse:   callback,
		timeoutMs: timeoutMs,
		stop:      make(chan bool),
		running:   false,
	}
}

// GetFakePulse creates and returns a fake pulse.
func (fp *FakePulsar) GetFakePulse() *core.Pulse {
	return fp.newPulse()
}

// Start starts sending a fake pulse.
func (fp *FakePulsar) Start(ctx context.Context) {
	fp.running = true
	go func(fp *FakePulsar) {
		for {
			select {
			case <-time.After(time.Millisecond * time.Duration(fp.timeoutMs)):
				{
					fp.onPulse(ctx, *fp.GetFakePulse())
				}
			case <-fp.stop:
				return
			}
		}

	}(fp)
}

// Stop sending a fake pulse.
func (fp *FakePulsar) Stop(ctx context.Context) {
	if fp.running {
		fp.stop <- true
		close(fp.stop)
		fp.running = false
	}
}

func (fp *FakePulsar) newPulse() *core.Pulse {
	generator := entropygenerator.StandardEntropyGenerator{}
	return &core.Pulse{
		PulseNumber:     0,
		NextPulseNumber: 0,
		Entropy:         generator.GenerateEntropy(),
	}
}
