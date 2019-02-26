/*
 *    Copyright 2019 Insolar
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

package nodes

import (
	"sync"

	"github.com/insolar/insolar"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
)

// Accessor provides info about active nodes.
//go:generate minimock -i github.com/insolar/insolar/ledger/storage/nodes.Accessor -o ./ -s _mock.go
type Accessor interface {
	All(pulse core.PulseNumber) ([]insolar.Node, error)
	InRole(pulse core.PulseNumber, role core.StaticRole) ([]insolar.Node, error)
}

// Setter provides methods for setting active nodes.
//go:generate minimock -i github.com/insolar/insolar/ledger/storage/nodes.Setter -o ./ -s _mock.go
type Setter interface {
	Set(pulse core.PulseNumber, nodes []insolar.Node) error
	Delete(pulse core.PulseNumber)
}

// Storage is an in-memory active node storage for each pulse. It's required to calculate node roles
// for past pulses to locate data.
// It should only contain previous N pulses. It should be stored on disk.
type Storage struct {
	lock  sync.RWMutex
	nodes map[core.PulseNumber][]insolar.Node
}

// NewStorage create new instance of Storage
func NewStorage() *Storage {
	// return new(nodeStorage)
	return &Storage{nodes: map[core.PulseNumber][]insolar.Node{}}
}

// Set saves active nodes for pulse in memory.
func (a *Storage) Set(pulse core.PulseNumber, nodes []insolar.Node) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if _, ok := a.nodes[pulse]; ok {
		return storage.ErrOverride
	}

	if len(nodes) != 0 {
		a.nodes[pulse] = append([]insolar.Node{}, nodes...)
	} else {
		a.nodes[pulse] = nil
	}

	return nil
}

// All return active nodes for specified pulse.
func (a *Storage) All(pulse core.PulseNumber) ([]insolar.Node, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	nodes, ok := a.nodes[pulse]
	if !ok {
		return nil, core.ErrNoNodes
	}
	res := append(nodes[:0:0], nodes...)

	return res, nil
}

// InRole return active nodes for specified pulse and role.
func (a *Storage) InRole(pulse core.PulseNumber, role core.StaticRole) ([]insolar.Node, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	nodes, ok := a.nodes[pulse]
	if !ok {
		return nil, core.ErrNoNodes
	}
	var inRole []insolar.Node
	for _, node := range nodes {
		if node.Role == role {
			inRole = append(inRole, node)
		}
	}

	return inRole, nil
}

// Delete erases nodes for specified pulse.
func (a *Storage) Delete(pulse core.PulseNumber) {
	a.lock.Lock()
	delete(a.nodes, pulse)
	a.lock.Unlock()
}
