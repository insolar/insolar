///
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
///

package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

// ContractService is a service that provides API for working with smart contracts.
type ContractService struct {
	runner *Runner
}

// NewContractService creates new Contract service instance.
func NewContractService(runner *Runner) *ContractService {
	return &ContractService{runner: runner}
}

//
// func (cs *ContractService) Call(response http.ResponseWriter, req *http.Request) error {
// 	cs.runner.callHandler()(response, req)
// 	return nil
// }

func (cs *ContractService) Call(req *http.Request, args *requester.Params, result *requester.Result) error {

	traceID := utils.RandTraceID()
	ctx, insLog := inslogger.WithTraceField(context.Background(), traceID)

	ctx, span := instracer.StartSpan(ctx, "callHandler")
	defer span.End()
	//
	// contractRequest := &requester.Request{}
	// contractAnswer := &requester.ContractAnswer{}
	// defer writeResponse(insLog, response, contractAnswer)

	// startTime := time.Now()
	// defer observeResultStatus(contractRequest.Method, contractAnswer, startTime)

	insLog.Infof("[ ContractService.Call ] Incoming request: %s", req.RequestURI)

	// ctx, rawBody, err := processRequest(ctx, req, contractRequest, result)
	// if err != nil {
	// 	processError(err, err.Error(), contractAnswer, insLog, traceID)
	// 	return err
	// }

	if args.Test != "" {
		insLog.Infof("Request related to %s", args.Test)
	}
	//
	// if contractRequest.Method != "contract.call" {
	// 	err := errors.New("rpc method does not exist")
	// 	processError(err, err.Error(), contractAnswer, insLog, traceID)
	// 	return err
	// }

	rawBody, err := json.Marshal(requester.Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "contract.call",
		Params:  *args,
	})

	signature, err := validateRequestHeaders(req.Header.Get(requester.Digest), req.Header.Get(requester.Signature), rawBody)
	if err != nil {
		return err
	}

	seedPulse, err := cs.runner.checkSeed(args.Seed)
	if err != nil {
		return err
	}

	setRootReferenceIfNeeded(args)

	callResult, err := cs.runner.makeCall(ctx, "contract.call", *args, rawBody, signature, 0, seedPulse)
	if err != nil {
		return err
	}

	result = &requester.Result{ContractResult: callResult, TraceID: traceID}

	return nil
}
