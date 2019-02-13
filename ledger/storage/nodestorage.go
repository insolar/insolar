/*
 *    Copyright 2019 Insolar Technologies
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

package storage

import (
	"fmt"
	"sync"

	"github.com/insolar/insolar/core"
)

// NodeStorage provides info about active nodes
//go:generate minimock -i github.com/insolar/insolar/ledger/storage.NodeStorage -o ./ -s _mock.go
type NodeStorage interface {
	SetActiveNodes(pulse core.PulseNumber, nodes []core.Node) error
	GetActiveNodes(pulse core.PulseNumber) ([]core.Node, error)
	GetActiveNodesByRole(pulse core.PulseNumber, role core.StaticRole) ([]core.Node, error)
	RemoveActiveNodesUntil(pulse core.PulseNumber)
}

type nodeStorage struct {
	DB DBContext `inject:""`

	// NodeHistory is an in-memory active node storage for each pulse. It's required to calculate node roles
	// for past pulses to locate data.
	// It should only contain previous N pulses. It should be stored on disk.
	nodeHistory     map[core.PulseNumber][]Node
	nodeHistoryLock sync.RWMutex
}

// NewNodeStorage create new instance of NodeStorage
func NewNodeStorage() NodeStorage {
	// return new(nodeStorage)
	return &nodeStorage{nodeHistory: map[core.PulseNumber][]Node{}}
}

// SetActiveNodes saves active nodes for pulse in memory.
func (a *nodeStorage) SetActiveNodes(pulse core.PulseNumber, nodes []core.Node) error {
	a.nodeHistoryLock.Lock()
	defer a.nodeHistoryLock.Unlock()

	if _, ok := a.nodeHistory[pulse]; ok {
		return ErrOverride
	}

	a.nodeHistory[pulse] = []Node{}
	for _, n := range nodes {
		a.nodeHistory[pulse] = append(a.nodeHistory[pulse], Node{
			FID:   n.ID(),
			FRole: n.Role(),
		})
	}

	return nil
}

// GetActiveNodes return active nodes for specified pulse.
func (a *nodeStorage) GetActiveNodes(pulse core.PulseNumber) ([]core.Node, error) {
	a.nodeHistoryLock.RLock()
	defer a.nodeHistoryLock.RUnlock()

	nodes, ok := a.nodeHistory[pulse]
	if !ok {
		return nil, fmt.Errorf("GetActiveNodes: no nodes for pulse %v", pulse)
	}
	res := make([]core.Node, len(nodes))
	for i, n := range nodes {
		res[i] = n
	}

	return res, nil
}

// GetActiveNodesByRole return active nodes for specified pulse and role.
func (a *nodeStorage) GetActiveNodesByRole(pulse core.PulseNumber, role core.StaticRole) ([]core.Node, error) {
	a.nodeHistoryLock.RLock()
	defer a.nodeHistoryLock.RUnlock()

	nodes, ok := a.nodeHistory[pulse]
	if !ok {
		return nil, fmt.Errorf("GetActiveNodesByRole: no nodes for pulse %v", pulse)
	}
	var inRole []core.Node
	for _, n := range nodes {
		if n.Role() == role {
			inRole = append(inRole, n)
		}
	}

	return inRole, nil
}

// RemoveActiveNodesUntil removes active nodes for all nodes less than provided pulse.
func (a *nodeStorage) RemoveActiveNodesUntil(pulse core.PulseNumber) {
	a.nodeHistoryLock.Lock()
	defer a.nodeHistoryLock.Unlock()

	for pn := range a.nodeHistory {
		if pn < pulse {
			delete(a.nodeHistory, pn)
		}
	}
}
