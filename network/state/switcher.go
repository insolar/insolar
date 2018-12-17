/*
 *    Copyright 2018 INS Ecosystem
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

package state

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
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
}

// NewNetworkSwitcher creates new NetworkSwitcher
func NewNetworkSwitcher() (*NetworkSwitcher, error) {
	return &NetworkSwitcher{
		state:     core.NoNetworkState,
		stateLock: sync.RWMutex{},
		counter:   1,
	}, nil
}

// TODO: after INS-923 remove this func
func (ns *NetworkSwitcher) Start(ctx context.Context) error {
	ns.stateLock.Lock()
	defer ns.stateLock.Unlock()

	ns.Release(ctx)
	ns.state = core.CompleteNetworkState
	return nil
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

	inslogger.FromContext(ctx).Info("Current NetworkSwitcher state is: %s", ns.state)

	if ns.SwitcherWorkAround.IsBootstrapped() && ns.state != core.CompleteNetworkState {
		ns.state = core.CompleteNetworkState
		ns.Release(ctx)
		inslogger.FromContext(ctx).Info("Current NetworkSwitcher state switched to: %s", ns.state)
	}

	return nil
}

// Acquire increases lock counter and locks message bus if it wasn't lock before
func (ns *NetworkSwitcher) Acquire(ctx context.Context) {
	inslogger.FromContext(ctx).Info("Call Acquire in NetworkSwitcher: ", ns.counter)
	ns.counter = ns.counter + 1
	if ns.counter-1 == 0 {
		inslogger.FromContext(ctx).Info("Lock MB")
		ns.MBLocker.Lock(ctx)
	}
}

// Release decreases lock counter and unlocks message bus if it wasn't lock by someone else
func (ns *NetworkSwitcher) Release(ctx context.Context) {
	inslogger.FromContext(ctx).Info("Call Release in NetworkSwitcher: ", ns.counter)
	if ns.counter == 0 {
		panic("Trying to unlock without locking")
	}
	ns.counter = ns.counter - 1
	if ns.counter == 0 {
		inslogger.FromContext(ctx).Info("Unlock MB")
		ns.MBLocker.Unlock(ctx)
	}
}
