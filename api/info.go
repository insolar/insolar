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

	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// InfoArgs is arguments that Info service accepts.
type InfoArgs struct{}

// InfoReply is reply for Info service requests.
type InfoReply struct {
	RootDomain             string   `json:"rootDomain"`
	RootMember             string   `json:"rootMember"`
	MigrationAdminMember   string   `json:"migrationAdminMember"`
	MigrationDaemonMembers []string `json:"migrationDaemonMembers"`
	NodeDomain             string   `json:"nodeDomain"`
	TraceID                string   `json:"traceID"`
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
//     "method": "network.getInfo",
//     "id": str|int|null
//     "params": { }
//   }
//
//     Response structure:
// 	{
// 		"jsonrpc": "2.0",
// 		"result": {
// 			"rootDomain": str, // reference to RootDomain instance
// 			"rootMember": str, // reference to RootMember instance
//			"migrationAdminMember": str, // reference to migrationAdminMember
//			"migrationDaemonMembers": [ //array string
//			 str, // reference to migrationDaemon
//			 str, // reference to migrationDaemon
//			 str, // reference to migrationDaemon
//],
// 			"nodeDomain": str, // reference to NodeDomain instance
// 			"traceID": str // traceID for request
// 		},
// 		"id": str|int|null // same as in request
// 	}
//
func (s *InfoService) GetInfo(r *http.Request, args *InfoArgs, reply *InfoReply) error {
	ctx, inslog := inslogger.WithTraceField(context.Background(), utils.RandTraceID())

	inslog.Infof("[ INFO ] Incoming request: %s", r.RequestURI)

	_, span := instracer.StartSpan(ctx, "InfoService.GetInfo")
	defer span.End()

	span.Annotate(nil, fmt.Sprintf("Incoming request: %s", r.RequestURI))

	rootDomain := s.runner.GenesisDataProvider.GetRootDomain(ctx)
	if rootDomain == nil {
		msg := "[ INFO ] rootDomain ref is nil"
		inslog.Error(msg)
		err := errors.New(msg)
		instracer.AddError(span, err)
		return err
	}
	rootMember, err := s.runner.GenesisDataProvider.GetRootMember(ctx)
	if err != nil {
		msg := "[ INFO ] Can't get rootMember ref"
		inslog.Error(errors.Wrap(err, msg))
		err := errors.Wrap(err, msg)
		instracer.AddError(span, err)
		return err
	}
	if rootMember == nil {
		msg := "[ INFO ] rootMember ref is nil"
		inslog.Error(msg)
		err := errors.New(msg)
		instracer.AddError(span, err)
		return err
	}
	migrationDaemonMembers, err := s.runner.GenesisDataProvider.GetMigrationDaemonMembers(ctx)
	if err != nil {
		msg := "[ INFO ] Can't get migration daemon members refs"
		inslog.Error(errors.Wrap(err, msg))
		instracer.AddError(span, err)
		return errors.Wrap(err, msg)
	}
	migrationDaemonMembersStrs := []string{}
	for _, r := range migrationDaemonMembers {
		if r == nil {
			msg := "[ INFO ] migration daemon members refs are nil"
			inslog.Error(msg)
			instracer.AddError(span, err)
			return errors.New(msg)
		}
		migrationDaemonMembersStrs = append(migrationDaemonMembersStrs, r.String())
	}
	migrationAdminMember, err := s.runner.GenesisDataProvider.GetMigrationAdminMember(ctx)
	if err != nil {
		msg := "[ INFO ] Can't get migration admin member ref"
		inslog.Error(errors.Wrap(err, msg))
		instracer.AddError(span, err)
		return errors.Wrap(err, msg)
	}
	if migrationAdminMember == nil {
		msg := "[ INFO ] migration admin member ref is nil"
		inslog.Error(msg)
		err := errors.New(msg)
		instracer.AddError(span, err)
		return err
	}
	nodeDomain, err := s.runner.GenesisDataProvider.GetNodeDomain(ctx)
	if err != nil {
		msg := "[ INFO ] Can't get nodeDomain ref"
		inslog.Error(errors.Wrap(err, msg))
		instracer.AddError(span, err)
		return errors.Wrap(err, msg)
	}
	if nodeDomain == nil {
		msg := "[ INFO ] nodeDomain ref is nil"
		inslog.Error(msg)
		err := errors.New(msg)
		instracer.AddError(span, err)
		return err
	}

	reply.RootDomain = rootDomain.String()
	reply.RootMember = rootMember.String()
	reply.MigrationAdminMember = migrationAdminMember.String()
	reply.MigrationDaemonMembers = migrationDaemonMembersStrs
	reply.NodeDomain = nodeDomain.String()
	reply.TraceID = utils.RandTraceID()

	return nil
}
