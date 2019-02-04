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

package state

import (
	"context"
	"sync"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/metrics"
	"go.opencensus.io/trace"
)

//go:generate minimock -i github.com/insolar/insolar/network/state.messageBusLocker -o ./ -s _mock.go
type messageBusLocker interface {
	Lock(ctx context.Context)
	Unlock(ctx context.Context)
}

// NetworkSwitcher is a network FSM using for bootstrapping
type NetworkSwitcher struct {
	NodeNetwork        core.NodeNetwork        `inject:""`
	SwitcherWorkAround core.SwitcherWorkAround `inject:""`
	MBLocker           messageBusLocker        `inject:""`

	counter uint64

	state     core.NetworkState
	stateLock sync.RWMutex
	span      *trace.Span
}

// NewNetworkSwitcher creates new NetworkSwitcher
func NewNetworkSwitcher() (*NetworkSwitcher, error) {
	return &NetworkSwitcher{
		state:     core.NoNetworkState,
		stateLock: sync.RWMutex{},
		counter:   1,
	}, nil
}

// GetState method returns current network state
func (ns *NetworkSwitcher) GetState() core.NetworkState {
	ns.stateLock.RLock()
	defer ns.stateLock.RUnlock()

	return ns.state
}

// OnPulse method checks current state and finds out reasons to update this state
func (ns *NetworkSwitcher) OnPulse(ctx context.Context, pulse core.Pulse) error {
	ns.stateLock.Lock()
	defer ns.stateLock.Unlock()

	ctx, span := instracer.StartSpan(ctx, "NetworkSwitcher.OnPulse")
	span.AddAttributes(
		trace.StringAttribute("NetworkSwitcher state: ", ns.state.String()),
	)
	defer span.End()
	inslogger.FromContext(ctx).Infof("Current NetworkSwitcher state is: %s", ns.state)

	if ns.SwitcherWorkAround.IsBootstrapped() && ns.state != core.CompleteNetworkState {
		ns.state = core.CompleteNetworkState
		ns.Release(ctx)
		metrics.NetworkComplete.Set(float64(time.Now().Unix()))
		inslogger.FromContext(ctx).Infof("Current NetworkSwitcher state switched to: %s", ns.state)
	}

	return nil
}

// Acquire increases lock counter and locks message bus if it wasn't lock before
func (ns *NetworkSwitcher) Acquire(ctx context.Context) {
	ctx, span := instracer.StartSpan(ctx, "NetworkSwitcher.Acquire")
	defer span.End()
	inslogger.FromContext(ctx).Info("Call Acquire in NetworkSwitcher: ", ns.counter)
	ns.counter = ns.counter + 1
	if ns.counter-1 == 0 {
		inslogger.FromContext(ctx).Info("Lock MB")
		ctx, ns.span = instracer.StartSpan(context.Background(), "GIL Lock (Lock MB)")
		ns.MBLocker.Lock(ctx)
	}
}

// Release decreases lock counter and unlocks message bus if it wasn't lock by someone else
func (ns *NetworkSwitcher) Release(ctx context.Context) {
	ctx, span := instracer.StartSpan(ctx, "NetworkSwitcher.Release")
	defer span.End()
	inslogger.FromContext(ctx).Info("Call Release in NetworkSwitcher: ", ns.counter)
	if ns.counter == 0 {
		panic("Trying to unlock without locking")
	}
	ns.counter = ns.counter - 1
	if ns.counter == 0 {
		inslogger.FromContext(ctx).Info("Unlock MB")
		ns.MBLocker.Unlock(ctx)
		ns.span.End()
	}
}
