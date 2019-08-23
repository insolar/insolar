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

	"github.com/insolar/rpc/v2"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/version"
)

// Get returns status info
func (s *NodeService) GetStatus(r *http.Request, args *interface{}, requestBody *rpc.RequestBody, reply *requester.StatusResponse) error {
	traceID := utils.RandTraceID()
	ctx, inslog := inslogger.WithTraceField(context.Background(), traceID)

	inslog.Infof("[ NodeService.GetStatus ] Incoming request: %s", r.RequestURI)

	statusReply := s.runner.NetworkStatus.GetNetworkStatus()

	reply.NetworkState = statusReply.NetworkState.String()
	reply.ActiveListSize = statusReply.ActiveListSize
	reply.WorkingListSize = statusReply.WorkingListSize

	nodes := make([]requester.Node, reply.ActiveListSize)
	for i, node := range statusReply.Nodes {
		nodes[i] = requester.Node{
			Reference: node.ID().String(),
			Role:      node.Role().String(),
			IsWorking: node.GetPower() > 0,
			ID:        uint32(node.ShortID()),
		}
	}
	reply.Nodes = nodes

	reply.Origin = requester.Node{
		Reference: statusReply.Origin.ID().String(),
		Role:      statusReply.Origin.Role().String(),
		IsWorking: statusReply.Origin.GetPower() > 0,
		ID:        uint32(statusReply.Origin.ShortID()),
	}

	reply.NetworkPulseNumber = uint32(statusReply.Pulse.PulseNumber)

	p, err := s.runner.PulseAccessor.Latest(ctx)
	if err != nil {
		p = *insolar.GenesisPulse
	}
	reply.PulseNumber = uint32(p.PulseNumber)

	reply.Entropy = statusReply.Pulse.Entropy[:]
	reply.Version = version.Version
	reply.StartTime = statusReply.StartTime
	reply.Timestamp = statusReply.Timestamp

	return nil
}
