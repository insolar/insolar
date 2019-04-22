//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package node

import (
	"sync"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/insolar/node.Accessor -o ./ -s _mock.go

// Accessor provides info about active nodes.
type Accessor interface {
	All(pulse insolar.PulseNumber) ([]insolar.Node, error)
	InRole(pulse insolar.PulseNumber, role insolar.StaticRole) ([]insolar.Node, error)
}

//go:generate minimock -i github.com/insolar/insolar/insolar/node.Modifier -o ./ -s _mock.go

// Modifier provides methods for setting active nodes.
type Modifier interface {
	Set(pulse insolar.PulseNumber, nodes []insolar.Node) error
	DeleteForPN(pulse insolar.PulseNumber)
}

// Storage is an in-memory active node storage for each pulse. It's required to calculate node roles
// for past pulses to locate data.
// It should only contain previous N pulses. It should be stored on disk.
type Storage struct {
	lock  sync.RWMutex
	nodes map[insolar.PulseNumber][]insolar.Node
}

// NewStorage create new instance of Storage
func NewStorage() *Storage {
	// return new(nodeStorage)
	return &Storage{nodes: map[insolar.PulseNumber][]insolar.Node{}}
}

// Set saves active nodes for pulse in memory.
func (a *Storage) Set(pulse insolar.PulseNumber, nodes []insolar.Node) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if _, ok := a.nodes[pulse]; ok {
		return ErrOverride
	}

	if len(nodes) != 0 {
		a.nodes[pulse] = append([]insolar.Node{}, nodes...)
	} else {
		a.nodes[pulse] = nil
	}

	return nil
}

// All return active nodes for specified pulse.
func (a *Storage) All(pulse insolar.PulseNumber) ([]insolar.Node, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	nodes, ok := a.nodes[pulse]
	if !ok {
		return nil, ErrNoNodes
	}
	res := append(nodes[:0:0], nodes...)

	return res, nil
}

// InRole return active nodes for specified pulse and role.
func (a *Storage) InRole(pulse insolar.PulseNumber, role insolar.StaticRole) ([]insolar.Node, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	nodes, ok := a.nodes[pulse]
	if !ok {
		return nil, ErrNoNodes
	}
	var inRole []insolar.Node
	for _, node := range nodes {
		if node.Role == role {
			inRole = append(inRole, node)
		}
	}

	return inRole, nil
}

// DeleteForPN erases nodes for specified pulse.
func (a *Storage) DeleteForPN(pulse insolar.PulseNumber) {
	a.lock.Lock()
	delete(a.nodes, pulse)
	a.lock.Unlock()
}
