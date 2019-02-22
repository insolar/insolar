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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage"
)

// Accessor provides info about active nodes.
//go:generate minimock -i github.com/insolar/insolar/ledger/storage/nodes.Accessor -o ./ -s _mock.go
type Accessor interface {
	All(pulse core.PulseNumber) ([]core.Node, error)
	InRole(pulse core.PulseNumber, role core.StaticRole) ([]core.Node, error)
}

// Setter provides methods for setting active nodes.
//go:generate minimock -i github.com/insolar/insolar/ledger/storage/nodes.Setter -o ./ -s _mock.go
type Setter interface {
	Set(pulse core.PulseNumber, nodes []core.Node) error
	Delete(pulse core.PulseNumber)
}

// Storage is an in-memory active node storage for each pulse. It's required to calculate node roles
// for past pulses to locate data.
// It should only contain previous N pulses. It should be stored on disk.
type Storage struct {
	lock  sync.RWMutex
	nodes map[core.PulseNumber][]Node
}

// NewStorage create new instance of Storage
func NewStorage() *Storage {
	// return new(nodeStorage)
	return &Storage{nodes: map[core.PulseNumber][]Node{}}
}

// Set saves active nodes for pulse in memory.
func (a *Storage) Set(pulse core.PulseNumber, nodes []core.Node) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if _, ok := a.nodes[pulse]; ok {
		return storage.ErrOverride
	}

	a.nodes[pulse] = []Node{}
	for _, n := range nodes {
		a.nodes[pulse] = append(a.nodes[pulse], Node{
			FID:   n.ID(),
			FRole: n.Role(),
		})
	}

	return nil
}

// All return active nodes for specified pulse.
func (a *Storage) All(pulse core.PulseNumber) ([]core.Node, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	nodes, ok := a.nodes[pulse]
	if !ok {
		return nil, core.ErrNoNodes
	}
	res := make([]core.Node, len(nodes))
	for i, n := range nodes {
		res[i] = n
	}

	return res, nil
}

// InRole return active nodes for specified pulse and role.
func (a *Storage) InRole(pulse core.PulseNumber, role core.StaticRole) ([]core.Node, error) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	nodes, ok := a.nodes[pulse]
	if !ok {
		return nil, core.ErrNoNodes
	}
	var inRole []core.Node
	for _, n := range nodes {
		if n.Role() == role {
			inRole = append(inRole, n)
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
