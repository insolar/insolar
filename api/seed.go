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
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/rpc/v2"
)

// SeedArgs is arguments that Seed service accepts.
type SeedArgs struct{}

// NodeService is a service that provides API for getting new seed and status.
type NodeService struct {
	runner *Runner
}

// NewNodeService creates new Node service instance.
func NewNodeService(runner *Runner) *NodeService {
	return &NodeService{runner: runner}
}

// Get returns new active seed.
//
//   Request structure:
//   {
//     "jsonrpc": "2.0",
//     "method": "node.getSeed",
//     "id": str|int|null
//   }
//
//     Response structure:
// 	{
// 		"jsonrpc": "2.0",
// 		"result": {
// 			"Seed": str, // correct seed for new Call request
// 			"TraceID": str // traceID for request
// 		},
// 		"id": str|int|null // same as in request
// 	}
//
func (s *NodeService) GetSeed(r *http.Request, args *SeedArgs, requestBody *rpc.RequestBody, reply *requester.SeedReply) error {
	traceID := utils.RandTraceID()
	ctx, inslog := inslogger.WithTraceField(context.Background(), traceID)

	_, span := instracer.StartSpan(ctx, "NodeService.getSeed")
	defer span.End()

	info := fmt.Sprintf("[ NodeService.GetSeed ] Incoming request: %s", r.RequestURI)
	inslog.Infof(info)
	span.Annotate(nil, info)

	seed, err := s.runner.SeedGenerator.Next()
	if err != nil {
		instracer.AddError(span, err)
		return errors.Wrap(err, "failed to get next seed")
	}

	pulse, err := s.runner.PulseAccessor.Latest(context.Background())
	if err != nil {
		instracer.AddError(span, err)
		return errors.Wrap(err, "[ NodeService::GetSeed ] Couldn't receive pulse")
	}
	s.runner.SeedManager.Add(*seed, pulse.PulseNumber)

	reply.Seed = seed[:]
	reply.TraceID = traceID

	return nil
}
