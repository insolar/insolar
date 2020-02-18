// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/insolar/insolar/api/instrumenter"
	"github.com/insolar/insolar/api/requester"
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

	setRootReferenceIfNeeded(args, runner.Options)

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
