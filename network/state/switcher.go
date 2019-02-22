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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network"
	"go.opencensus.io/trace"
)

//go:generate minimock -i github.com/insolar/insolar/network/state.messageBusLocker -o ./ -s _mock.go
type messageBusLocker interface {
	Lock(ctx context.Context)
	Unlock(ctx context.Context)
}

type transitionData struct {
	pulse        core.Pulse
	majorityRule bool
	minRoleRule  bool
}

type stateHandler func(context.Context, *transitionData)

// NetworkSwitcher is a network FSM using for bootstrapping
type NetworkSwitcher struct {
	NodeNetwork        core.NodeNetwork        `inject:""`
	Rules              network.Rules           `inject:""`
	SwitcherWorkAround core.SwitcherWorkAround `inject:""`
	MBLocker           messageBusLocker        `inject:""`

	counter uint64

	state     core.NetworkState
	stateMap  map[core.NetworkState]stateHandler
	stateLock sync.RWMutex

	span *trace.Span
}

// NewNetworkSwitcher creates new NetworkSwitcher
func NewNetworkSwitcher() (*NetworkSwitcher, error) {
	ns := &NetworkSwitcher{
		state:     core.NoNetworkState,
		stateLock: sync.RWMutex{},
		counter:   1,
	}

	ns.stateMap = map[core.NetworkState]stateHandler{
		core.NoNetworkState:       ns.handleNoNetworkState,
		core.VoidNetworkState:     ns.handleVoidNetworkState,
		core.CompleteNetworkState: ns.handleCompleteNetworkState,
	}

	return ns, nil
}

// GetState method returns current network state
func (ns *NetworkSwitcher) GetState() core.NetworkState {
	ns.stateLock.RLock()
	defer ns.stateLock.RUnlock()

	return ns.state
}

func (ns *NetworkSwitcher) handleNoNetworkState(ctx context.Context, transitionData *transitionData) {
	if ns.SwitcherWorkAround.IsBootstrapped() && transitionData.majorityRule {
		ns.state = core.VoidNetworkState
	}
}

func (ns *NetworkSwitcher) handleVoidNetworkState(ctx context.Context, transitionData *transitionData) {
	if !transitionData.majorityRule {
		ns.state = core.NoNetworkState
	}

	if transitionData.minRoleRule {
		defer ns.Release(ctx)

		ns.state = core.CompleteNetworkState
	}
}

func (ns *NetworkSwitcher) handleCompleteNetworkState(ctx context.Context, transitionData *transitionData) {
	if !transitionData.majorityRule || !transitionData.minRoleRule {
		defer ns.Acquire(ctx)

		if !transitionData.majorityRule {
			ns.state = core.NoNetworkState
		}

		if !transitionData.minRoleRule {
			ns.state = core.VoidNetworkState
		}
	}
}

func (ns *NetworkSwitcher) changeState(ctx context.Context, pulse core.Pulse) core.NetworkState {
	majorityOk, _ := ns.Rules.CheckMajorityRule()
	minRoleOk := ns.Rules.CheckMinRole()

	transitionData := &transitionData{
		pulse:        pulse,
		majorityRule: majorityOk,
		minRoleRule:  minRoleOk,
	}

	for {
		state := ns.state
		handler := ns.stateMap[state]

		handler(ctx, transitionData)

		stateChanged := state != ns.state
		if !stateChanged {
			break
		}
	}

	return ns.state
}

// OnPulse method checks current state and finds out reasons to update this state
func (ns *NetworkSwitcher) OnPulse(ctx context.Context, pulse core.Pulse) error {
	ns.stateLock.Lock()
	defer ns.stateLock.Unlock()

	oldState := ns.state

	ctx, span := instracer.StartSpan(ctx, "NetworkSwitcher.changeState")
	span.AddAttributes(
		trace.StringAttribute("NetworkSwitcher state: ", oldState.String()),
	)
	defer span.End()

	newState := ns.changeState(ctx, pulse)

	if oldState != newState {
		inslogger.FromContext(ctx).WithFields(map[string]interface{}{
			"oldState": oldState.String(),
			"newState": newState.String(),
		}).Infof("Current NetworkSwitcher state switched")
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
