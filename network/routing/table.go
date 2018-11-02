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

package routing

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/pkg/errors"
)

type Table struct {
	keeper network.NodeKeeper
}

func (t *Table) isLocalNode(core.RecordRef) bool {
	return true
}

func (t *Table) resolveRemoteNode(ref core.RecordRef) (string, error) {
	return "", errors.New("not implemented")
}

// Resolve NodeID -> Address. Can initiate network requests.
func (t *Table) Resolve(ref core.RecordRef) (string, error) {
	if t.isLocalNode(ref) {
		node := t.keeper.GetActiveNode(ref)
		if node == nil {
			return "", errors.New("no such local node")
		}
		return node.PhysicalAddress(), nil
	}
	return t.resolveRemoteNode(ref)
}

// AddToKnownHosts add host to routing table.
func (t *Table) AddToKnownHosts(*host.Host) {
	log.Warn("not implemented")
}

// Rebalance recreate shards of routing table with known hosts according to new partition policy.
func (t *Table) Rebalance(PartitionPolicy) {
	log.Warn("not implemented")
}
