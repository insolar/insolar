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
	"errors"
	"net/http"
	"strings"

	"github.com/insolar/insolar/application/api/instrumenter"
	"github.com/insolar/insolar/application/api/requester"
	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/rpc/v2"
	"github.com/insolar/rpc/v2/json2"
)

func wrapCall(ctx context.Context, runner *Runner, allowedMethods map[string]bool, req *http.Request, args *requester.Params, requestBody *rpc.RequestBody, result *requester.ContractResult) error {
	instr := instrumenter.GetInstrumenter(ctx)
	traceID := instr.TraceID()
	logger := inslogger.FromContext(ctx)

	if !runner.AvailabilityChecker.IsAvailable(ctx) {
		logger.Error("API is not available")

		instr.SetError(errors.New(ServiceUnavailableErrorMessage), ServiceUnavailableErrorShort)
		return &json2.Error{
			Code:    ServiceUnavailableError,
			Message: ServiceUnavailableErrorMessage,
			Data: requester.Data{
				TraceID: traceID,
			},
		}
	}

	_, ok := allowedMethods[args.CallSite]
	if !ok {
		logger.Warnf("CallSite '%s' is not in list of allowed methods", args.CallSite)
		instr.SetError(errors.New(MethodNotFoundErrorMessage), MethodNotFoundErrorShort)
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
	instr.SetCallSite(args.CallSite)

	signature, err := validateRequestHeaders(req.Header.Get(requester.Digest), req.Header.Get(requester.Signature), requestBody.Raw)
	if err != nil {
		logger.Warn("validateRequestHeaders return error: ", err.Error())
		instr.SetError(err, InvalidParamsErrorShort)
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
		logger.Warn("checkSeed returned error: ", err.Error())
		instr.SetError(err, InvalidRequestErrorShort)
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

	callResult, requestRef, err := runner.makeCall(ctx, *args, requestBody.Raw, signature, seedPulse)

	var ref string
	if requestRef != nil {
		ref = requestRef.String()
	}

	if err != nil {
		// TODO: white list of errors that doesnt require log
		logger.Error("API return error: ", err.Error())
		if strings.Contains(err.Error(), "invalid signature") {
			instr.SetError(err, UnauthorizedErrorShort)
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
			instr.SetError(err, ParseErrorShort)
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

		instr.SetError(err, ExecutionErrorShort)
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
