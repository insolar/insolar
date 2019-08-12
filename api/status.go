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

package api

import (
	"context"
	"net/http"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/version"
)

type Node struct {
	Reference string
	Role      string
	IsWorking bool
	ID        uint32
}

// StatusReply is reply for Status service requests.
type StatusReply struct {
	NetworkState       string
	Origin             Node
	ActiveListSize     int
	WorkingListSize    int
	Nodes              []Node
	PulseNumber        uint32
	NetworkPulseNumber uint32
	Entropy            []byte
	NodeState          string
	Version            string
}

// Get returns status info
func (s *NodeService) GetStatus(r *http.Request, args *interface{}, requestBody *RequestBody, reply *StatusReply) error {
	traceID := utils.RandTraceID()
	ctx, inslog := inslogger.WithTraceField(context.Background(), traceID)

	inslog.Infof("[ NodeService.GetStatus ] Incoming request: %s", r.RequestURI)

	reply.NetworkState = s.runner.ServiceNetwork.GetState().String()
	reply.NodeState = s.runner.NodeNetwork.GetOrigin().GetState().String()

	np, err := s.runner.ServiceNetwork.(*servicenetwork.ServiceNetwork).PulseAccessor.GetLatestPulse(ctx)
	if err != nil {
		np = *insolar.GenesisPulse
	}

	p, err := s.runner.PulseAccessor.Latest(ctx)
	if err != nil {
		p = *insolar.GenesisPulse
	}

	activeNodes := s.runner.NodeNetwork.GetAccessor(np.PulseNumber).GetActiveNodes()
	workingNodes := s.runner.NodeNetwork.GetAccessor(np.PulseNumber).GetWorkingNodes()

	reply.ActiveListSize = len(activeNodes)
	reply.WorkingListSize = len(workingNodes)

	nodes := make([]Node, reply.ActiveListSize)
	for i, node := range activeNodes {
		nodes[i] = Node{
			Reference: node.ID().String(),
			Role:      node.Role().String(),
			IsWorking: node.GetState() == insolar.NodeReady,
			ID:        uint32(node.ShortID()),
		}
	}

	reply.Nodes = nodes
	origin := s.runner.NodeNetwork.GetOrigin()
	reply.Origin = Node{
		Reference: origin.ID().String(),
		Role:      origin.Role().String(),
		IsWorking: origin.GetState() == insolar.NodeReady,
		ID:        uint32(origin.ShortID()),
	}

	reply.NetworkPulseNumber = uint32(np.PulseNumber)
	reply.PulseNumber = uint32(p.PulseNumber)
	reply.Entropy = np.Entropy[:]
	reply.Version = version.Version

	return nil
}
