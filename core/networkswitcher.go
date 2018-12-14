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

package core

import (
	"context"
)

// NetworkState type for bootstrapping process
type NetworkState int

//go:generate stringer -type=NetworkState
const (
	// NoNetworkState state means that nodes doesn`t match majority_rule
	NoNetworkState NetworkState = iota
	// VoidNetworkState state means that nodes have not complete min_role_count rule for proper work
	VoidNetworkState
	// JetlessNetworkState state means that every Jet need proof completeness of stored data
	JetlessNetworkState
	// AuthorizationNetworkState state means that every node need to validate ActiveNodeList using NodeDomain
	AuthorizationNetworkState
	// CompleteNetworkState state means network is ok and ready for proper work
	CompleteNetworkState
)

// NetworkSwitcher is a network FSM using for bootstrapping
//go:generate minimock -i github.com/insolar/insolar/core.NetworkSwitcher -o ../testutils -s _mock.go
type NetworkSwitcher interface {
	// GetState method returns current network state
	GetState() NetworkState
	// OnPulse method checks current state and finds out reasons to update this state
	OnPulse(context.Context, Pulse) error
}

//go:generate minimock -i github.com/insolar/insolar/core.GlobalInsolarLock -o ../testutils -s _mock.go
// GlobalInsolarLock is lock of all incoming and outcoming network calls.
// It's not intended to be used in multiple threads. And main use of it is `Set` method of `PulseManager`.
type GlobalInsolarLock interface {
	Acquire(ctx context.Context)
	Release(ctx context.Context)
}
