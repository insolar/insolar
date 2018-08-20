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

package relay

import (
	"github.com/insolar/insolar/network/host/id"
	"github.com/insolar/insolar/network/host/node"

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
	// NoAuth - this error returns if node tries to send relay request but not authenticated.
	NoAuth
)

// Relay Interface for relaying
type Relay interface {
	// AddClient add client to relay list.
	AddClient(node *node.Node) error
	// RemoveClient removes client from relay list.
	RemoveClient(node *node.Node) error
	// ClientsCount - clients count.
	ClientsCount() int
	// NeedToRelay returns true if origin node is proxy for target node.
	NeedToRelay(targetAddress string) bool
}

type relay struct {
	clients []*node.Node
}

// NewRelay constructs relay list.
func NewRelay() Relay {
	return &relay{
		clients: make([]*node.Node, 0),
	}
}

// AddClient add client to relay list.
func (r *relay) AddClient(node *node.Node) error {
	if _, n := r.findClient(node.ID); n != nil {
		return errors.New("client exists already")
	}
	r.clients = append(r.clients, node)
	return nil
}

// RemoveClient removes client from relay list.
func (r *relay) RemoveClient(node *node.Node) error {
	idx, n := r.findClient(node.ID)
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

// NeedToRelay returns true if origin node is proxy for target node.
func (r *relay) NeedToRelay(targetAddress string) bool {
	for i := 0; i < r.ClientsCount(); i++ {
		if r.clients[i].Address.String() == targetAddress {
			return true
		}
	}
	return false
}

func (r *relay) findClient(id id.ID) (int, *node.Node) {
	for idx, nodeIterator := range r.clients {
		if nodeIterator.ID.HashEqual(id.Hash) {
			return idx, nodeIterator
		}
	}
	return -1, nil
}
