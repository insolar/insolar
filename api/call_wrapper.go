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
	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/rpc/v2"
	"github.com/pkg/errors"
	"net/http"
)

func wrapCall(runner *Runner, allowedMethods map[string]bool, req *http.Request, args *requester.Params, requestBody *rpc.RequestBody, result *requester.ContractResult) error {
	traceID := utils.RandTraceID()
	ctx, logger := inslogger.WithTraceField(context.Background(), traceID)

	ctx, span := instracer.StartSpan(ctx, "Call")
	defer span.End()

	logger.Infof("[ ContractService.Call ] Incoming request: %s", req.RequestURI)

	_, ok := allowedMethods[args.CallSite]
	if !ok {
		return errors.New("method not allowed")
	}

	if args.Test != "" {
		logger.Infof("ContractRequest related to %s", args.Test)
	}

	signature, err := validateRequestHeaders(req.Header.Get(requester.Digest), req.Header.Get(requester.Signature), requestBody.Raw)
	if err != nil {
		return err
	}

	seedPulse, err := runner.checkSeed(args.Seed)
	if err != nil {
		return err
	}

	setRootReferenceIfNeeded(args)

	callResult, requestRef, err := runner.makeCall(ctx, "contract.call", *args, requestBody.Raw, signature, 0, seedPulse)
	if err != nil {
		return err
	}

	if requestRef != nil {
		result.RequestReference = requestRef.String()
	}
	result.CallResult = callResult
	result.TraceID = traceID
	return nil
}
