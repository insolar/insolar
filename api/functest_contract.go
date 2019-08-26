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

// +build functest

package api

import (
	"context"
	"net/http"
	"reflect"

	"github.com/pkg/errors"

	"github.com/insolar/rpc/v2"

	"github.com/insolar/insolar/application/extractor"
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
	Code string
	Name string
}

// UploadReply is reply that Contract.Upload returns
type UploadReply struct {
	PrototypeRef string `json:"PrototypeRef"`
	TraceID      string `json:"TraceID"`
}

// Upload builds code and return prototype ref
func (s *FuncTestContractService) Upload(r *http.Request, args *UploadArgs, requestBody *rpc.RequestBody, reply *UploadReply) error {
	ctx, inslog := inslogger.WithTraceField(context.Background(), utils.RandTraceID())
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

	err := s.cb.Build(ctx, contractMap)
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
	ctx, inslog := inslogger.WithTraceField(context.Background(), utils.RandTraceID())
	reply.TraceID = utils.TraceID(ctx)

	inslog.Infof("[ FuncTestContractService.CallConstructor ] Incoming request: %s", r.RequestURI)

	if len(args.PrototypeRefString) == 0 {
		return errors.New("params.PrototypeRefString is missing")
	}

	protoRef, err := insolar.NewReferenceFromBase58(args.PrototypeRefString)
	if err != nil {
		return errors.Wrap(err, "can't get protoRef")
	}

	pulse, err := s.runner.PulseAccessor.Latest(ctx)
	if err != nil {
		return errors.Wrap(err, "can't get current pulse")
	}

	base := insolar.GenesisRecord.Ref()
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

// CallConstructor make an object from its prototype
func (s *FuncTestContractService) CallMethod(r *http.Request, args *CallMethodArgs, requestBody *rpc.RequestBody, re *CallMethodReply) error {
	ctx, inslog := inslogger.WithTraceField(context.Background(), utils.RandTraceID())
	re.TraceID = utils.TraceID(ctx)

	inslog.Infof("[ FuncTestContractService.CallMethod ] Incoming request: %s", r.RequestURI)

	if len(args.ObjectRefString) == 0 {
		return errors.New("params.ObjectRefString is missing")
	}

	objectRef, err := insolar.NewReferenceFromBase58(args.ObjectRefString)
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

	callReply, _, err := s.runner.ContractRequester.Call(ctx, msg)
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
