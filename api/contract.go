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

	"github.com/insolar/insolar/application/extractor"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/insolar/insolar/testutils"
)

// ContractService is a service that provides ability to add custom contracts
type ContractService struct {
	runner *Runner
	cb     *goplugintestutils.ContractsBuilder
}

// NewContractService creates new Contract service instance.
func NewContractService(runner *Runner) *ContractService {
	return &ContractService{runner: runner}
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
func (s *ContractService) Upload(r *http.Request, args *UploadArgs, reply *UploadReply) error {
	ctx, inslog := inslogger.WithTraceField(context.Background(), utils.RandTraceID())
	reply.TraceID = utils.TraceID(ctx)

	inslog.Infof("[ ContractService.Upload ] Incoming request: %s", r.RequestURI)

	if len(args.Name) == 0 {
		return errors.New("params.name is missing")
	}

	if len(args.Code) == 0 {
		return errors.New("params.code is missing")
	}

	if s.cb == nil {
		insgocc, err := goplugintestutils.BuildPreprocessor()
		if err != nil {
			inslog.Infof("[ ContractService.Upload ] can't build preprocessor %#v", err)
			return errors.Wrap(err, "can't build preprocessor")
		}
		s.cb = goplugintestutils.NewContractBuilder(s.runner.ArtifactManager, insgocc, s.runner.PulseAccessor)
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

// CallConstructorReply is reply that Contract.CallConstructor returns
type CallConstructorReply struct {
	ObjectRef string `json:"ObjectRef"`
	TraceID   string `json:"TraceID"`
}

// CallConstructor make an object from its prototype
func (s *ContractService) CallConstructor(r *http.Request, args *CallConstructorArgs, reply *CallConstructorReply) error {
	ctx, inslog := inslogger.WithTraceField(context.Background(), utils.RandTraceID())
	reply.TraceID = utils.TraceID(ctx)

	inslog.Infof("[ ContractService.CallConstructor ] Incoming request: %s", r.RequestURI)

	if len(args.PrototypeRefString) == 0 {
		return errors.New("params.PrototypeRefString is missing")
	}

	protoRef, err := insolar.NewReferenceFromBase58(args.PrototypeRefString)
	if err != nil {
		return errors.Wrap(err, "can't get protoRef")
	}

	base := insolar.GenesisRecord.Ref()
	msg := &message.CallMethod{
		IncomingRequest: record.IncomingRequest{
			Method:          args.Method,
			Arguments:       args.MethodArgs,
			Base:            &base,
			Caller:          testutils.RandomRef(),
			CallerPrototype: testutils.RandomRef(),
			Prototype:       protoRef,
			CallType:        record.CTSaveAsChild,
			APIRequestID:    utils.TraceID(ctx),
		},
	}

	callConstructorReply, err := s.runner.ContractRequester.CallConstructor(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "CallConstructor error")
	}

	reply.ObjectRef = callConstructorReply.String()

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
	Reply          reply.CallMethod  `json:"Reply"`
	ExtractedReply interface{}       `json:"ExtractedReply"`
	Error          *foundation.Error `json:"Error"`
	ExtractedError string            `json:"ExtractedError"`
	TraceID        string            `json:"TraceID"`
}

// CallConstructor make an object from its prototype
func (s *ContractService) CallMethod(r *http.Request, args *CallMethodArgs, re *CallMethodReply) error {
	ctx, inslog := inslogger.WithTraceField(context.Background(), utils.RandTraceID())
	re.TraceID = utils.TraceID(ctx)

	inslog.Infof("[ ContractService.CallMethod ] Incoming request: %s", r.RequestURI)

	if len(args.ObjectRefString) == 0 {
		return errors.New("params.ObjectRefString is missing")
	}

	objectRef, err := insolar.NewReferenceFromBase58(args.ObjectRefString)
	if err != nil {
		return errors.Wrap(err, "can't get objectRef")
	}

	msg := &message.CallMethod{
		IncomingRequest: record.IncomingRequest{
			Caller:       testutils.RandomRef(),
			Object:       objectRef,
			Method:       args.Method,
			Arguments:    args.MethodArgs,
			APIRequestID: utils.TraceID(ctx),
		},
	}

	callMethodReply, err := s.runner.ContractRequester.Call(ctx, msg)
	if err != nil {
		inslogger.FromContext(ctx).Error("failed to call: ", err.Error())
		return errors.Wrap(err, "CallMethod failed with error")
	}

	re.Reply = *callMethodReply.(*reply.CallMethod)

	var extractedReply interface{}
	extractedReply, _, err = extractor.CallResponse(re.Reply.Result)
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

	return nil
}
