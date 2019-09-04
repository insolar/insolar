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
	"net/http"
	"strings"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/rpc/v2"
	"github.com/insolar/rpc/v2/json2"
)

const (
	ParseError                 = -31700
	ParseErrorMessage          = "Parsing error on the server side: received an invalid JSON."
	InvalidRequestError        = -31600
	InvalidRequestErrorMessage = "The JSON received is not a valid request payload."
	MethodNotFoundError        = -31601
	MethodNotFoundErrorMessage = "Method does not exist / is not available."
	InvalidParamsError         = -31602
	InvalidParamsErrorMessage  = "Invalid method parameter(s)."
	InternalError              = -31603
	InternalErrorMessage       = "Internal JSON RPC error."
	TimeoutError               = -31106
	TimeoutErrorMessage        = "Request's timeout has expired."
	UnauthorizedError          = -31401
	UnauthorizedErrorMessage   = "Action is not authorized."
	ExecutionError             = -31103
	ExecutionErrorMessage      = "Execution error."
)

func wrapCall(runner *Runner, allowedMethods map[string]bool, req *http.Request, args *requester.Params, requestBody *rpc.RequestBody, result *requester.ContractResult) error {
	traceID := utils.RandTraceID()
	ctx, logger := inslogger.WithTraceField(context.Background(), traceID)

	ctx, span := instracer.StartSpan(ctx, "Call")
	defer span.End()

	logger.Infof("[ ContractService.Call ] Incoming request: %s", req.RequestURI)

	_, ok := allowedMethods[args.CallSite]
	if !ok {
		return &json2.Error{
			Code:    MethodNotFoundError,
			Message: MethodNotFoundErrorMessage,
			Data: requester.Data{
				TraceID: traceID,
			},
		}
	}

	if args.Test != "" {
		logger.Infof("ContractRequest related to %s", args.Test)
	}

	signature, err := validateRequestHeaders(req.Header.Get(requester.Digest), req.Header.Get(requester.Signature), requestBody.Raw)
	if err != nil {
		return &json2.Error{
			Code:    InvalidParamsError,
			Message: InvalidParamsErrorMessage,
			Data: requester.Data{
				Trace:   strings.Split(err.Error(), ": "),
				TraceID: traceID,
			},
		}
	}

	seedPulse, err := runner.checkSeed(args.Seed)
	if err != nil {
		return &json2.Error{
			Code:    InvalidRequestError,
			Message: InvalidRequestErrorMessage,
			Data: requester.Data{
				Trace:   []string{err.Error()},
				TraceID: traceID,
			},
		}
	}

	setRootReferenceIfNeeded(args)

	callResult, requestRef, err := runner.makeCall(ctx, "contract.call", *args, requestBody.Raw, signature, seedPulse)

	var ref string
	if requestRef != nil {
		ref = requestRef.String()
	}

	if err != nil {
		// TODO: white list of errors that doesnt require log
		logger.Error(err.Error())
		if strings.Contains(err.Error(), "invalid signature") {
			return &json2.Error{
				Code:    UnauthorizedError,
				Message: UnauthorizedErrorMessage,
				Data: requester.Data{
					Trace:            strings.Split(err.Error(), ": "),
					TraceID:          traceID,
					RequestReference: ref,
				},
			}
		}

		if strings.Contains(err.Error(), "failed to parse") {
			return &json2.Error{
				Code:    ParseError,
				Message: ParseErrorMessage,
				Data: requester.Data{
					Trace:            strings.Split(err.Error(), ": "),
					TraceID:          traceID,
					RequestReference: ref,
				},
			}
		}

		return &json2.Error{
			Code:    ExecutionError,
			Message: ExecutionErrorMessage,
			Data: requester.Data{
				Trace:            strings.Split(err.Error(), ": "),
				TraceID:          traceID,
				RequestReference: ref,
			},
		}
	}

	result.RequestReference = ref
	result.CallResult = callResult
	result.TraceID = traceID
	return nil
}
