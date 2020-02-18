// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package api

import (
	"context"
	"net/http"
	"reflect"

	"github.com/insolar/insolar/applicationbase/genesis"
	"github.com/insolar/rpc/v2"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/api/instrumenter"
	"github.com/insolar/insolar/applicationbase/extractor"
	"github.com/insolar/insolar/insolar"
	insolarApi "github.com/insolar/insolar/insolar/api"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
)

// FuncTestContractService is a service that provides ability to add custom contracts
type FuncTestContractService struct {
	runner *Runner
	cb     *goplugintestutils.ContractsBuilder
}

// NewFuncTestContractService creates new Contract service instance.
func NewFuncTestContractService(runner *Runner) *FuncTestContractService {
	return &FuncTestContractService{runner: runner}
}

// UploadArgs is arguments that Contract.Upload accepts.
type UploadArgs struct {
	Code                string
	Name                string
	PanicIsLogicalError bool
}

// UploadReply is reply that Contract.Upload returns
type UploadReply struct {
	PrototypeRef string `json:"PrototypeRef"`
	TraceID      string `json:"TraceID"`
}

// Upload builds code and return prototype ref
func (s *FuncTestContractService) Upload(r *http.Request, args *UploadArgs, requestBody *rpc.RequestBody, reply *UploadReply) error {
	ctx, instr := instrumenter.NewMethodInstrument("FuncTestContractService.Upload")
	defer instr.End()

	inslog := inslogger.FromContext(ctx)

	reply.TraceID = utils.TraceID(ctx)

	inslog.Infof("[ FuncTestContractService.Upload ] Incoming request: %s", r.RequestURI)

	if len(args.Name) == 0 {
		return errors.New("params.name is missing")
	}

	if len(args.Code) == 0 {
		return errors.New("params.code is missing")
	}

	if s.cb == nil {
		insgocc, err := goplugintestutils.BuildPreprocessor()
		if err != nil {
			inslog.Infof("[ FuncTestContractService.Upload ] can't build preprocessor %#v", err)
			return errors.Wrap(err, "can't build preprocessor")
		}
		s.cb = goplugintestutils.NewContractBuilder(
			insgocc, s.runner.ArtifactManager, s.runner.PulseAccessor, s.runner.JetCoordinator,
		)
	}

	contractMap := make(map[string]string)
	contractMap[args.Name] = args.Code
	buildOptions := goplugintestutils.BuildOptions{PanicIsLogicalError: args.PanicIsLogicalError}

	err := s.cb.Build(ctx, contractMap, buildOptions)
	if err != nil {
		return errors.Wrap(err, "can't build contract")
	}
	reference := *s.cb.Prototypes[args.Name]
	reply.PrototypeRef = reference.String()
	return nil
}

// CallConstructorArgs is arguments that Contract.CallConstructor accepts.
type CallConstructorArgs struct {
	PrototypeRefString string
	Method             string
	MethodArgs         []byte
}

// CallConstructor make an object from its prototype
func (s *FuncTestContractService) CallConstructor(r *http.Request, args *CallConstructorArgs, requestBody *rpc.RequestBody, reply *CallMethodReply) error {
	ctx, instr := instrumenter.NewMethodInstrument("FuncTestContractService.CallConstructor")
	defer instr.End()

	inslog := inslogger.FromContext(ctx)

	reply.TraceID = utils.TraceID(ctx)

	inslog.Infof("[ FuncTestContractService.CallConstructor ] Incoming request: %s", r.RequestURI)

	if len(args.PrototypeRefString) == 0 {
		return errors.New("params.PrototypeRefString is missing")
	}

	protoRef, err := insolar.NewReferenceFromString(args.PrototypeRefString)
	if err != nil {
		return errors.Wrap(err, "can't get protoRef")
	}

	pulse, err := s.runner.PulseAccessor.Latest(ctx)
	if err != nil {
		return errors.Wrap(err, "can't get current pulse")
	}

	base := genesis.Record.Ref()
	msg := &payload.CallMethod{
		Request: &record.IncomingRequest{
			Method:          args.Method,
			Arguments:       args.MethodArgs,
			Base:            &base,
			CallerPrototype: gen.Reference(),
			Prototype:       protoRef,
			CallType:        record.CTSaveAsChild,
			APIRequestID:    utils.TraceID(ctx),
			Reason:          insolarApi.MakeReason(pulse.PulseNumber, args.MethodArgs),
			APINode:         s.runner.JetCoordinator.Me(),
		},
	}

	err = s.call(ctx, msg, reply)
	if err != nil {
		return err
	}

	return nil
}

// CallMethodArgs is arguments that Contract.CallMethod accepts.
type CallMethodArgs struct {
	ObjectRefString string
	Method          string
	MethodArgs      []byte
}

// CallMethodReply is reply that Contract.CallMethod returns
type CallMethodReply struct {
	Object         string            `json:"ObjectRef"`
	Result         []byte            `json:"Result"`
	ExtractedReply interface{}       `json:"ExtractedReply"`
	ExtractedError string            `json:"ExtractedError"`
	Error          *foundation.Error `json:"FoundationError"`
	TraceID        string            `json:"TraceID"`
}

// CallMethod make an object from its prototype
func (s *FuncTestContractService) CallMethod(r *http.Request, args *CallMethodArgs, requestBody *rpc.RequestBody, re *CallMethodReply) error {
	ctx, instr := instrumenter.NewMethodInstrument("FuncTestContractService.CallMethod")
	defer instr.End()

	inslog := inslogger.FromContext(ctx)

	re.TraceID = utils.TraceID(ctx)

	inslog.Infof("[ FuncTestContractService.CallMethod ] Incoming request: %s", r.RequestURI)

	if len(args.ObjectRefString) == 0 {
		return errors.New("params.ObjectRefString is missing")
	}

	objectRef, err := insolar.NewReferenceFromString(args.ObjectRefString)
	if err != nil {
		return errors.Wrap(err, "can't get objectRef")
	}

	pulse, err := s.runner.PulseAccessor.Latest(ctx)
	if err != nil {
		return errors.Wrap(err, "can't get current pulse")
	}

	msg := &payload.CallMethod{
		Request: &record.IncomingRequest{
			Object:       objectRef,
			Method:       args.Method,
			Arguments:    args.MethodArgs,
			APIRequestID: utils.TraceID(ctx),
			Reason:       insolarApi.MakeReason(pulse.PulseNumber, args.MethodArgs),
			APINode:      s.runner.JetCoordinator.Me(),
		},
	}

	err = s.call(ctx, msg, re)
	if err != nil {
		return err
	}

	return nil
}

// CallConstructor make an object from its prototype
func (s *FuncTestContractService) call(ctx context.Context, msg insolar.Payload, re *CallMethodReply) error {
	inslog := inslogger.FromContext(ctx)

	callReply, _, err := s.runner.ContractRequester.SendRequest(ctx, msg)
	if err != nil {
		inslog.Error("failed to call: ", err.Error())
		return errors.Wrap(err, "CallMethod failed with error")
	}

	typedReply := callReply.(*reply.CallMethod)
	if typedReply.Object != nil {
		re.Object = typedReply.Object.String()
	}
	re.Result = typedReply.Result

	extractedReply, foundationError, err := extractor.CallResponse(re.Result)
	if err != nil {
		return errors.Wrap(err, "Can't extract response")
	}

	// TODO need to understand why sometimes errors goes to reply
	// see tests TestConstructorReturnNil, TestContractCallingContract, TestPrototypeMismatch
	switch extractedReply.(type) {
	case map[string]interface{}:
		replyMap := extractedReply.(map[string]interface{})
		if len(replyMap) == 1 {
			for k, v := range replyMap {
				if reflect.ValueOf(k).String() == "S" && len(reflect.TypeOf(v).String()) > 0 {
					re.ExtractedError = reflect.ValueOf(v).String()
				}
			}
		}
	default:
		re.ExtractedReply = extractedReply
	}

	re.Error = foundationError

	return nil
}
