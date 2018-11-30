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
	"strconv"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/pkg/errors"
)

type Table struct {
	NodeKeeper network.NodeKeeper
}

func (t *Table) isLocalNode(core.RecordRef) bool {
	return true
}

func (t *Table) resolveRemoteNode(ref core.RecordRef) (*host.Host, error) {
	return nil, errors.New("not implemented")
}

func (t *Table) addRemoteHost(h *host.Host) {
	log.Warn("not implemented")
}

// Resolve NodeID -> ShortID, Address. Can initiate network requests.
func (t *Table) Resolve(ref core.RecordRef) (*host.Host, error) {
	if t.isLocalNode(ref) {
		node := t.NodeKeeper.GetActiveNode(ref)
		if node == nil {
			return nil, errors.New("no such local node with NodeID: " + ref.String())
		}
		return host.NewHostNS(node.PhysicalAddress(), node.ID(), node.ShortID())
	}
	return t.resolveRemoteNode(ref)
}

// ResolveS ShortID -> NodeID, Address for node inside current globe.
func (t *Table) ResolveS(id core.ShortNodeID) (*host.Host, error) {
	node := t.NodeKeeper.GetActiveNodeByShortID(id)
	if node == nil {
		return nil, errors.New("no such local node with ShortID: " + strconv.FormatUint(uint64(id), 10))
	}
	return host.NewHostNS(node.PhysicalAddress(), node.ID(), node.ShortID())
}

// AddToKnownHosts add host to routing table.
func (t *Table) AddToKnownHosts(h *host.Host) {
	if t.isLocalNode(h.NodeID) {
		// we should already have this node in NodeNetwork active list, do nothing
		return
	}
	t.addRemoteHost(h)
}

// GetRandomNodes get a specified number of random nodes. Returns less if there are not enough nodes in network.
func (t *Table) GetRandomNodes(count int) []host.Host {
	// not so random for now
	nodes := t.NodeKeeper.GetActiveNodes()
	resultCount := count
	if count > len(nodes) {
		resultCount = len(nodes)
	}
	result := make([]host.Host, 0)
	for i := 0; i < resultCount; i++ {
		address, err := host.NewAddress(nodes[i].PhysicalAddress())
		if err != nil {
			log.Error(err)
			continue
		}
		h := host.Host{NodeID: nodes[i].ID(), Address: address}
		result = append(result, h)
	}
	return result
}

// Rebalance recreate shards of routing table with known hosts according to new partition policy.
func (t *Table) Rebalance(network.PartitionPolicy) {
	log.Warn("not implemented")
}

func (t *Table) Inject(nodeKeeper network.NodeKeeper) {
	t.NodeKeeper = nodeKeeper
}
