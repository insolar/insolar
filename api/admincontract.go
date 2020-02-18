// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package api

import (
	"net/http"

	"github.com/insolar/insolar/api/instrumenter"
	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/rpc/v2"
)

// AdminContractService is a service that provides API for working with smart contracts.
type AdminContractService struct {
	runner         *Runner
	allowedMethods map[string]bool
}

// NewAdminContractService creates new AdminContract service instance.
func NewAdminContractService(runner *Runner) *AdminContractService {
	allowedMethods := map[string]bool{}
	if runner.Options.AdminContractMethods != nil {
		allowedMethods = runner.Options.AdminContractMethods
	}
	allowedMethods["contract.registerNode"] = true
	allowedMethods["contract.getNodeRef"] = true
	return &AdminContractService{runner: runner, allowedMethods: allowedMethods}
}

func (cs *AdminContractService) Call(req *http.Request, args *requester.Params, requestBody *rpc.RequestBody, result *requester.ContractResult) error {
	ctx, instr := instrumenter.NewMethodInstrument("AdminContractService.call")
	defer instr.End()

	inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"callSite": args.CallSite,
		"uri":      req.RequestURI,
		"service":  "AdminContractService",
		"params":   args.CallParams,
		"seed":     args.Seed,
	}).Infof("Incoming request")

	return wrapCall(ctx, cs.runner, cs.allowedMethods, req, args, requestBody, result)
}
