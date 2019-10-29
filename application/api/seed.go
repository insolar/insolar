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
	"strings"

	"github.com/pkg/errors"

	"github.com/insolar/rpc/v2"
	"github.com/insolar/rpc/v2/json2"

	"github.com/insolar/insolar/application/api/instrumenter"
	"github.com/insolar/insolar/application/api/requester"
	"github.com/insolar/insolar/instrumentation/inslogger"
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
//	Request structure:
//	{
//		"jsonrpc": "2.0",
//		"method": "node.getSeed",
//		"id": str|int|null
//	}
//
//	Response structure:
//	{
//		"jsonrpc": "2.0",
//		"result": {
//			"seed": str, // correct seed for new Call request
//			"traceID": str // traceID for request
//		},
//		"id": str|int|null // same as in request
//	}
//
func (s *NodeService) getSeed(ctx context.Context, _ *http.Request, _ *SeedArgs, _ *rpc.RequestBody, reply *requester.SeedReply) error {
	traceID := instrumenter.GetTraceID(ctx)

	seed, err := s.runner.SeedGenerator.Next()
	if err != nil {
		return err
	}

	pulse, err := s.runner.PulseAccessor.Latest(context.Background())
	if err != nil {
		return errors.Wrap(err, "couldn't receive pulse")
	}
	s.runner.SeedManager.Add(*seed, pulse.PulseNumber)

	reply.Seed = seed[:]
	reply.TraceID = traceID

	return nil
}

func (s *NodeService) GetSeed(r *http.Request, args *SeedArgs, requestBody *rpc.RequestBody, reply *requester.SeedReply) error {
	ctx, instr := instrumenter.NewMethodInstrument("NodeService.getSeed")
	defer instr.End()

	msg := fmt.Sprint("Incoming request: ", r.RequestURI)
	instr.Annotate(msg)

	logger := inslogger.FromContext(ctx)
	logger.Info("[ NodeService.getSeed ] ", msg)

	if !s.runner.AvailabilityChecker.IsAvailable(ctx) {
		logger.Error("[ NodeService.getSeed ] API is not available")

		instr.SetError(errors.New(ServiceUnavailableErrorMessage), ServiceUnavailableErrorShort)
		return &json2.Error{
			Code:    ServiceUnavailableError,
			Message: ServiceUnavailableErrorMessage,
			Data: requester.Data{
				TraceID: instr.TraceID(),
			},
		}
	}

	err := s.getSeed(ctx, r, args, requestBody, reply)
	if err != nil {
		logger.Error("[ NodeService.getSeed ] failed to execute: ", err.Error())

		instr.SetError(err, InternalErrorShort)
		return &json2.Error{
			Code:    InternalError,
			Message: InternalErrorMessage,
			Data: requester.Data{
				Trace:   strings.Split(err.Error(), ": "),
				TraceID: instr.TraceID(),
			},
		}
	}
	return nil
}
