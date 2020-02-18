// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/insolar/rpc/v2"

	"github.com/insolar/insolar/api/instrumenter"
	"github.com/insolar/insolar/applicationbase/genesisrefs"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// InfoArgs is arguments that Info service accepts.
type InfoArgs struct{}

// InfoService is a service that provides API for getting info about genesis objects.
type InfoService struct {
	runner *Runner
}

// NewInfoService creates new Info service instance.
func NewInfoService(runner *Runner) *InfoService {
	return &InfoService{runner: runner}
}

func (s *InfoService) getInfo(_ context.Context, _ *http.Request, _ *InfoArgs, _ *rpc.RequestBody, reply *map[string]interface{}) error {
	nodeDomain := genesisrefs.ContractNodeDomain
	if nodeDomain.IsEmpty() {
		return errors.New("nodeDomain ref is nil")
	}

	*reply = map[string]interface{}{
		"nodeDomain": nodeDomain.String(),
		"traceID":    utils.RandTraceID(),
	}

	return nil
}

// GetInfo returns info about genesis objects.
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
func (s *InfoService) GetInfo(r *http.Request, args *InfoArgs, requestBody *rpc.RequestBody, reply *map[string]interface{}) error {
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

	for name, val := range s.runner.Options.InfoResponse {
		(*reply)[name] = val
	}

	return err
}
