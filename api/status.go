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
	"errors"
	"net/http"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/version"
	"github.com/insolar/rpc/v2"
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
	Version            string
	Timestamp          time.Time
	StartTime          time.Time
}

// Get returns status info
func (s *NodeService) GetStatus(r *http.Request, args *interface{}, requestBody *rpc.RequestBody, reply *StatusReply) error {

	if s.runner.cfg.IsPublic {
		return errors.New("method not allowed")
	}

	traceID := utils.RandTraceID()
	ctx, inslog := inslogger.WithTraceField(context.Background(), traceID)

	inslog.Infof("[ NodeService.GetStatus ] Incoming request: %s", r.RequestURI)
	statusReply := s.runner.NetworkStatus.GetNetworkStatus()

	reply.NetworkState = statusReply.NetworkState.String()
	reply.ActiveListSize = statusReply.ActiveListSize
	reply.WorkingListSize = statusReply.WorkingListSize

	nodes := make([]Node, reply.ActiveListSize)
	for i, node := range statusReply.Nodes {
		nodes[i] = Node{
			Reference: node.ID().String(),
			Role:      node.Role().String(),
			IsWorking: node.GetPower() > 0,
			ID:        uint32(node.ShortID()),
		}
	}
	reply.Nodes = nodes

	reply.Origin = Node{
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
