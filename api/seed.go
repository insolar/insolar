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

package api

import (
	"context"
	"net/http"

	"github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

// SeedArgs is arguments that Seed service accepts.
type SeedArgs struct{}

// SeedReply is reply for Seed service requests.
type SeedReply struct {
	Seed    []byte
	TraceID string
}

// SeedService is a service that provides API for getting new seed.
type SeedService struct {
	runner *Runner
}

// NewSeedService creates new Seed service instance.
func NewSeedService(runner *Runner) *SeedService {
	return &SeedService{runner: runner}
}

// Get returns new active seed.
//
//   Request structure:
//   {
//     "jsonrpc": "2.0",
//     "method": "seed.Get",
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
func (s *SeedService) Get(r *http.Request, args *SeedArgs, reply *SeedReply) error {
	traceID := utils.RandTraceID()
	_, inslog := inslogger.WithTraceField(context.Background(), traceID)

	inslog.Infof("[ SeedService.Get ] Incoming request: %s", r.RequestURI)

	seed, err := s.runner.SeedGenerator.Next()
	if err != nil {
		return errors.Wrap(err, "[ GetSeed ]")
	}
	s.runner.SeedManager.Add(*seed)

	reply.Seed = seed[:]
	reply.TraceID = traceID

	return nil
}
