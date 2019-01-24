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

package storage

import (
	"fmt"

	"github.com/insolar/insolar/core"
)

// SetActiveNodes saves active nodes for pulse in memory.
func (db *DB) SetActiveNodes(pulse core.PulseNumber, nodes []core.Node) error {
	db.nodeHistoryLock.Lock()
	defer db.nodeHistoryLock.Unlock()

	if _, ok := db.nodeHistory[pulse]; ok {
		return ErrOverride
	}

	db.nodeHistory[pulse] = []Node{}
	for _, n := range nodes {
		db.nodeHistory[pulse] = append(db.nodeHistory[pulse], Node{
			FID:   n.ID(),
			FRole: n.Role(),
		})
	}

	return nil
}

// GetActiveNodes return active nodes for specified pulse.
func (db *DB) GetActiveNodes(pulse core.PulseNumber) ([]core.Node, error) {
	db.nodeHistoryLock.RLock()
	defer db.nodeHistoryLock.RUnlock()

	nodes, ok := db.nodeHistory[pulse]
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
func (db *DB) GetActiveNodesByRole(pulse core.PulseNumber, role core.StaticRole) ([]core.Node, error) {
	db.nodeHistoryLock.RLock()
	defer db.nodeHistoryLock.RUnlock()

	nodes, ok := db.nodeHistory[pulse]
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
func (db *DB) RemoveActiveNodesUntil(pulse core.PulseNumber) {
	db.nodeHistoryLock.Lock()
	defer db.nodeHistoryLock.Unlock()
	fmt.Printf("cleanLightData: RemoveActiveNodesUntil: %v\n", pulse)

	for pn := range db.nodeHistory {
		if pn < pulse {
			delete(db.nodeHistory, pulse)
		}
	}
}
