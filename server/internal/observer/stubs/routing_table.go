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

package stubs

import (
	"errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
)

type RoutingTable struct {
}

// Resolve NodeID -> ShortID, Address. Can initiate network requests.
func (rt *RoutingTable) Resolve(insolar.Reference) (*host.Host, error) {
	return nil, errors.New("not implemented")
}

// ResolveConsensus ShortID -> NodeID, Address for node inside current globe for current consensus.
func (rt *RoutingTable) ResolveConsensus(insolar.ShortNodeID) (*host.Host, error) {
	return nil, errors.New("not implemented")
}

// ResolveConsensusRef NodeID -> ShortID, Address for node inside current globe for current consensus.
func (rt *RoutingTable) ResolveConsensusRef(insolar.Reference) (*host.Host, error) {
	return nil, errors.New("not implemented")
}

// AddToKnownHosts add host to routing table.
func (rt *RoutingTable) AddToKnownHosts(*host.Host) {
}

// Rebalance recreate shards of routing table with known hosts according to new partition policy.
func (rt *RoutingTable) Rebalance(pp network.PartitionPolicy) {
	panic(errors.New("not implemented"))
}
