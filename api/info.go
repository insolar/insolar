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

// InfoArgs is arguments that Info service accepts.
type InfoArgs struct{}

// InfoReply is reply for Info service requests.
type InfoReply struct {
	RootDomain string
	RootMember string
	NodeDomain string
	TraceID    string
}

// InfoService is a service that provides API for getting info about genesis objects.
type InfoService struct {
	runner *Runner
}

// NewInfoService creates new Info service instance.
func NewInfoService(runner *Runner) *InfoService {
	return &InfoService{runner: runner}
}

// Get returns info about genesis objects.
//
//   Request structure:
//   {
//     "jsonrpc": "2.0",
//     "method": "info.Get",
//     "id": str|int|null
//   }
//
//     Response structure:
// 	{
// 		"jsonrpc": "2.0",
// 		"result": {
// 			"RootDomain": str, // reference to RootDomain instance
// 			"RootMember": str, // reference to RootMember instance
// 			"NodeDomain": str, // reference to NodeDomain instance
// 			"TraceID": str // traceID for request
// 		},
// 		"id": str|int|null // same as in request
// 	}
//
func (s *InfoService) Get(r *http.Request, args *InfoArgs, reply *InfoReply) error {
	ctx, inslog := inslogger.WithTraceField(context.Background(), utils.RandTraceID())

	inslog.Infof("[ INFO ] Incoming request: %s", r.RequestURI)

	rootDomain := s.runner.GenesisDataProvider.GetRootDomain(ctx)
	if rootDomain == nil {
		inslog.Error("[ INFO ] rootDomain ref is nil")
		return errors.New("[ INFO ] rootDomain ref is nil")
	}
	rootMember, err := s.runner.GenesisDataProvider.GetRootMember(ctx)
	if err != nil {
		inslog.Error(errors.Wrap(err, "[ INFO ] Can't get rootMember ref"))
		return errors.Wrap(err, "[ INFO ] Can't get rootMember ref")
	}
	if rootMember == nil {
		inslog.Error("[ INFO ] rootMember ref is nil")
		return errors.New("[ INFO ] rootMember ref is nil")
	}
	nodeDomain, err := s.runner.GenesisDataProvider.GetNodeDomain(ctx)
	if err != nil {
		inslog.Error(errors.Wrap(err, "[ INFO ] Can't get nodeDomain ref"))
		return errors.Wrap(err, "[ INFO ] Can't get nodeDomain ref")
	}
	if nodeDomain == nil {
		inslog.Error("[ INFO ] nodeDomain ref is nil")
		return errors.New("[ INFO ] nodeDomain ref is nil")
	}

	reply.RootDomain = rootDomain.String()
	reply.RootMember = rootMember.String()
	reply.NodeDomain = nodeDomain.String()
	reply.TraceID = utils.RandTraceID()

	return nil
}
