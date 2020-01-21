// Copyright 2020 Insolar Network Ltd.
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

package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/insolar/rpc/v2"

	"github.com/insolar/insolar/application/api/instrumenter"
	"github.com/insolar/insolar/application/genesisrefs"
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
	FeeMember              string   `json:"feeMember"`
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
//	Request structure:
//	{
//		"jsonrpc": "2.0",
//		"method": "network.getInfo",
//		"id": str|int|null
//		"params": { }
//	}
//
//	Response structure:
//	{
//		"jsonrpc": "2.0",
//		"result": {
//			"rootDomain": str, // reference to RootDomain instance
//			"rootMember": str, // reference to RootMember instance
//			"migrationAdminMember": str, // reference to migrationAdminMember
//			"migrationDaemonMembers": [ //array string
//				str, // reference to migrationDaemon
//				str, // reference to migrationDaemon
//				str, // reference to migrationDaemon
//			],
//			"nodeDomain": str, // reference to NodeDomain instance
//			"traceID": str // traceID for request
//		},
//		"id": str|int|null // same as in request
//	}
//
func (s *InfoService) getInfo(_ context.Context, _ *http.Request, _ *InfoArgs, _ *rpc.RequestBody, reply *InfoReply) error {
	rootDomain := genesisrefs.ContractRootDomain
	if rootDomain.IsEmpty() {
		return errors.New("rootDomain ref is nil")
	}

	rootMember := genesisrefs.ContractRootMember
	if rootMember.IsEmpty() {
		return errors.New("rootMember ref is nil")
	}

	migrationDaemonMembers := genesisrefs.ContractMigrationDaemonMembers
	migrationDaemonMembersStrs := make([]string, 0)
	for _, r := range migrationDaemonMembers {
		if r.IsEmpty() {
			return errors.New("migration daemon members refs are nil")
		}
		migrationDaemonMembersStrs = append(migrationDaemonMembersStrs, r.String())
	}

	migrationAdminMember := genesisrefs.ContractMigrationAdminMember
	if migrationAdminMember.IsEmpty() {
		return errors.New("migration admin member ref is nil")
	}
	feeMember := genesisrefs.ContractFeeMember
	if feeMember.IsEmpty() {
		return errors.New("feeMember ref is nil")
	}
	nodeDomain := genesisrefs.ContractNodeDomain
	if nodeDomain.IsEmpty() {
		return errors.New("nodeDomain ref is nil")
	}

	reply.RootDomain = rootDomain.String()
	reply.RootMember = rootMember.String()
	reply.MigrationAdminMember = migrationAdminMember.String()
	reply.FeeMember = feeMember.String()
	reply.MigrationDaemonMembers = migrationDaemonMembersStrs
	reply.NodeDomain = nodeDomain.String()
	reply.TraceID = utils.RandTraceID()

	return nil
}

func (s *InfoService) GetInfo(r *http.Request, args *InfoArgs, requestBody *rpc.RequestBody, reply *InfoReply) error {
	ctx, instr := instrumenter.NewMethodInstrument("InfoService.getInfo")
	defer instr.End()

	msg := fmt.Sprint("Incoming request: ", r.RequestURI)
	instr.Annotate(msg)

	logger := inslogger.FromContext(ctx)
	logger.Info("[ InfoService.getInfo ] ", msg)

	err := s.getInfo(ctx, r, args, requestBody, reply)
	if err != nil {
		logger.Error("[ InfoService.getInfo ] failed to execute: ", err.Error())
		err = errors.Wrap(err, "Failed to execute InfoService.getInfo")
		instr.SetError(err, InternalErrorShort)
	}
	return err
}
