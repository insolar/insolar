//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package routing

import (
	"strconv"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/pkg/errors"
)

type Table struct {
	NodeKeeper network.NodeKeeper `inject:""`
}

func (t *Table) ResolveConsensus(id insolar.ShortNodeID) (*host.Host, error) {
	node := t.NodeKeeper.GetAccessor().GetActiveNodeByShortID(id)
	if node != nil {
		return host.NewHostNS(node.Address(), node.ID(), node.ShortID())
	}
	h := t.NodeKeeper.GetConsensusInfo().ResolveConsensus(id)
	if h == nil {
		return nil, errors.New("no such local node with ShortID: " + strconv.FormatUint(uint64(id), 10))
	}
	return h, nil
}

func (t *Table) ResolveConsensusRef(ref insolar.Reference) (*host.Host, error) {
	node := t.NodeKeeper.GetAccessor().GetActiveNode(ref)
	if node != nil {
		return host.NewHostNS(node.Address(), node.ID(), node.ShortID())
	}
	h := t.NodeKeeper.GetConsensusInfo().ResolveConsensusRef(ref)
	if h == nil {
		return nil, errors.New("no such local node with node ID: " + ref.String())
	}
	return h, nil
}

func (t *Table) isLocalNode(insolar.Reference) bool {
	return true
}

func (t *Table) resolveRemoteNode(_ insolar.Reference) (*host.Host, error) {
	return nil, errors.New("not implemented")
}

func (t *Table) addRemoteHost(_ *host.Host) {
	log.Warn("not implemented")
}

// Resolve NodeID -> ShortID, Address. Can initiate network requests.
func (t *Table) Resolve(ref insolar.Reference) (*host.Host, error) {
	if t.isLocalNode(ref) {
		node := t.NodeKeeper.GetAccessor().GetActiveNode(ref)
		if node == nil {
			return nil, errors.New("no such local node with NodeID: " + ref.String())
		}
		return host.NewHostNS(node.Address(), node.ID(), node.ShortID())
	}
	return t.resolveRemoteNode(ref)
}

// AddToKnownHosts add host to routing table.
func (t *Table) AddToKnownHosts(h *host.Host) {
	if t.isLocalNode(h.NodeID) {
		// we should already have this node in NodeNetwork active list, do nothing
		return
	}
	t.addRemoteHost(h)
}

// Rebalance recreate shards of routing table with known hosts according to new partition policy.
func (t *Table) Rebalance(network.PartitionPolicy) {
	log.Warn("not implemented")
}
