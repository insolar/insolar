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
	Acquire(ctx context.Context)
	Release(ctx context.Context)
}

// NetworkSwitcher is a network FSM using for bootstrapping
type NetworkSwitcher struct {
	NodeNetwork        core.NodeNetwork        `inject:""`
	SwitcherWorkAround core.SwitcherWorkAround `inject:""`
	MBLocker           messageBusLocker        `inject:""`

	state     core.NetworkState
	stateLock sync.RWMutex

	mbLocks     map[string]*lock
	mbLocksLock sync.RWMutex
}

type lock struct {
	sync.RWMutex
	lockCount int
}

// NewNetworkSwitcher creates new NetworkSwitcher
func NewNetworkSwitcher() (*NetworkSwitcher, error) {
	ns := &NetworkSwitcher{
		state:     core.NoNetworkState,
		stateLock: sync.RWMutex{},
		mbLocks:   make(map[string]*lock),
	}
	ns.mbLocksLock.Lock()
	ns.mbLocks["NetworkSwitcher"] = &lock{lockCount: 1}
	ns.mbLocksLock.Unlock()
	return ns, nil
}

// TODO: after INS-923 remove this func
func (ns *NetworkSwitcher) Start(ctx context.Context) error {
	ns.stateLock.Lock()
	defer ns.stateLock.Unlock()

	ns.ReleaseGlobalLock(ctx, "NetworkSwitcher")
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
		ns.ReleaseGlobalLock(ctx, "NetworkSwitcher")
		inslogger.FromContext(ctx).Info("Current NetworkSwitcher state switched to: %s", ns.state)
	}

	return nil
}

func (ns *NetworkSwitcher) AcquireGlobalLock(ctx context.Context, caller string) {
	ns.mbLocksLock.Lock()

	callerLock, ok := ns.mbLocks[caller]
	if !ok {
		ns.mbLocks[caller] = &lock{}
	}
	callerLock = ns.mbLocks[caller]
	ns.mbLocksLock.Unlock()

	callerLock.Lock()
	callerLock.lockCount = callerLock.lockCount + 1
	callerLock.Unlock()

	ns.mbLocksLock.Lock()
	defer ns.mbLocksLock.Unlock()
	for _, lock := range ns.mbLocks {
		if lock.lockCount != 0 {
			return
		}
	}
	ns.MBLocker.Acquire(ctx)
}

func (ns *NetworkSwitcher) ReleaseGlobalLock(ctx context.Context, caller string) {
	ns.mbLocksLock.Lock()

	callerLock, ok := ns.mbLocks[caller]
	ns.mbLocksLock.Unlock()
	if !ok {
		panic("You are trying to unlock GlobalLock without previously locking it!")
	}

	callerLock.Lock()
	if callerLock.lockCount == 0 {
		callerLock.Unlock()
		panic("You are trying to unlock GlobalLock without previously locking it!")
	}
	callerLock.lockCount = callerLock.lockCount - 1
	callerLock.Unlock()

	ns.mbLocksLock.Lock()
	defer ns.mbLocksLock.Unlock()
	for _, lock := range ns.mbLocks {
		if lock.lockCount != 0 {
			return
		}
	}
	ns.MBLocker.Release(ctx)
}
