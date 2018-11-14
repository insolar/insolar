/*
 *    Copyright 2018 Insolar
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

package relay

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/transport/host"

	"errors"
)

// State is alias for relaying state
type State int

const (
	// Unknown unknown relay state.
	Unknown = State(iota + 1)
	// Started is relay type means relaying started.
	Started
	// Stopped is relay type means relaying stopped.
	Stopped
	// Error is relay type means error state change.
	Error
	// NoAuth - this error returns if host tries to send relay request but not authenticated.
	NoAuth
)

// Relay Interface for relaying
type Relay interface {
	// AddClient add client to relay list.
	AddClient(host *host.Host) error
	// RemoveClient removes client from relay list.
	RemoveClient(host *host.Host) error
	// ClientsCount - clients count.
	ClientsCount() int
	// NeedToRelay returns true if origin host is proxy for target host.
	NeedToRelay(targetAddress string) bool
}

type relay struct {
	clients []*host.Host
}

// NewRelay constructs relay list.
func NewRelay() Relay {
	return &relay{
		clients: make([]*host.Host, 0),
	}
}

// AddClient add client to relay list.
func (r *relay) AddClient(host *host.Host) error {
	if _, n := r.findClient(host.NodeID); n != nil {
		return errors.New("client exists already")
	}
	r.clients = append(r.clients, host)
	return nil
}

// RemoveClient removes client from relay list.
func (r *relay) RemoveClient(host *host.Host) error {
	idx, n := r.findClient(host.NodeID)
	if n == nil {
		return errors.New("client not found")
	}
	r.clients = append(r.clients[:idx], r.clients[idx+1:]...)
	return nil
}

// ClientsCount - returns clients count.
func (r *relay) ClientsCount() int {
	return len(r.clients)
}

// NeedToRelay returns true if origin host is proxy for target host.
func (r *relay) NeedToRelay(targetAddress string) bool {
	for i := 0; i < r.ClientsCount(); i++ {
		if r.clients[i].Address.String() == targetAddress {
			return true
		}
	}
	return false
}

func (r *relay) findClient(id core.RecordRef) (int, *host.Host) {
	for idx, hostIterator := range r.clients {
		if hostIterator.NodeID.Equal(id) {
			return idx, hostIterator
		}
	}
	return -1, nil
}
