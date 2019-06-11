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

	"github.com/insolar/insolar/application/extractor"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/insolar/insolar/testutils"
)

// ContractService is a service that provides ability to add custom contracts.
type ContractService interface {
	Upload(r *http.Request, args *UploadArgs, reply *UploadReply) error
	CallConstructor(r *http.Request, args *CallConstructorArgs, reply *CallConstructorReply) error
	CallMethod(r *http.Request, args *CallMethodArgs, reply *CallMethodReply) error
}

// ContractServiceReal is the ContractService implementation that provides ability to add custom contracts.
type ContractServiceReal struct {
	runner *Runner
}

// NewContractService creates new Contract service instance.
func NewContractService(runner *Runner, enable bool) ContractService {
	if enable {
		return &ContractServiceReal{runner: runner}
	}
	return &ContractServiceDummy{}
}

// UploadArgs is arguments that Contract.Upload accepts.
type UploadArgs struct {
	Code string
	Name string
}

// UploadReply is reply that Contract.Upload returns
type UploadReply struct {
	PrototypeRef string `json:"PrototypeRef"`
}

// Upload builds code and return prototype ref
func (s *ContractServiceReal) Upload(r *http.Request, args *UploadArgs, reply *UploadReply) error {
	_, inslog := inslogger.WithTraceField(context.Background(), utils.RandTraceID())

	inslog.Infof("[ ContractService.Upload ] Incoming request: %s", r.RequestURI)

	if len(args.Name) == 0 {
		return errors.New("params.name is missing")
	}

	if len(args.Code) == 0 {
		return errors.New("params.code is missing")
	}

	insgocc, err := goplugintestutils.BuildPreprocessor()
	if err != nil {
		inslog.Infof("[ ContractService.Upload ] can't build preprocessor %#v", err)
		return errors.Wrap(err, "can't build preprocessor")
	}
	cb := goplugintestutils.NewContractBuilder(s.runner.ArtifactManager, insgocc)

	contractMap := make(map[string]string)
	contractMap[args.Name] = args.Code

	err = cb.Build(contractMap)
	if err != nil {
		inslog.Infof("[ ContractService.Upload ] can't build contract %#v", err)
		return errors.Wrap(err, "can't build contract")
	}
	reference := *cb.Prototypes[args.Name]
	reply.PrototypeRef = reference.String()
	return nil
}

// CallConstructorArgs is arguments that Contract.CallConstructor accepts.
type CallConstructorArgs struct {
	PrototypeRefString string
}

// CallConstructorReply is reply that Contract.CallConstructor returns
type CallConstructorReply struct {
	ObjectRef string `json:"ObjectRef"`
}

// CallConstructor make an object from its prototype
func (s *ContractServiceReal) CallConstructor(r *http.Request, args *CallConstructorArgs, reply *CallConstructorReply) error {
	ctx, inslog := inslogger.WithTraceField(context.Background(), utils.RandTraceID())

	inslog.Infof("[ ContractService.CallConstructor ] Incoming request: %s", r.RequestURI)

	if len(args.PrototypeRefString) == 0 {
		return errors.New("params.PrototypeRefString is missing")
	}

	protoRef, err := insolar.NewReferenceFromBase58(args.PrototypeRefString)
	if err != nil {
		return errors.Wrap(err, "can't get protoRef")
	}

	domain, err := insolar.NewReferenceFromBase58("4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa")
	if err != nil {
		return errors.Wrap(err, "can't get domain reference")
	}

	base := testutils.RandomRef()

	contractID, err := s.runner.ArtifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType:     record.CTSaveAsChild,
			Prototype:    &base,
			APIRequestID: utils.TraceID(ctx),
		},
	)

	if err != nil {
		return errors.Wrap(err, "can't register request")
	}

	objectRef := insolar.Reference{}
	objectRef.SetRecord(*contractID)

	memory, _ := insolar.Serialize(nil)

	_, err = s.runner.ArtifactManager.ActivateObject(
		ctx,
		*domain,
		objectRef,
		insolar.GenesisRecord.Ref(),
		*protoRef,
		false,
		memory,
	)

	if err != nil {
		return errors.Wrap(err, "can't activate object")
	}

	reply.ObjectRef = objectRef.String()

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
	Reply          reply.CallMethod `json:"Reply"`
	ExtractedReply interface{}      `json:"ExtractedReply"`
}

// CallConstructor make an object from its prototype
func (s *ContractServiceReal) CallMethod(r *http.Request, args *CallMethodArgs, re *CallMethodReply) error {
	ctx, inslog := inslogger.WithTraceField(context.Background(), utils.RandTraceID())

	inslog.Infof("[ ContractService.CallMethod ] Incoming request: %s", r.RequestURI)

	if len(args.ObjectRefString) == 0 {
		return errors.New("params.ObjectRefString is missing")
	}

	objectRef, err := insolar.NewReferenceFromBase58(args.ObjectRefString)
	if err != nil {
		return errors.Wrap(err, "can't get objectRef")
	}

	msg := &message.CallMethod{
		Request: record.Request{
			Caller:       testutils.RandomRef(),
			Object:       objectRef,
			Method:       args.Method,
			Arguments:    args.MethodArgs,
			APIRequestID: utils.TraceID(ctx),
		},
	}

	callMethodReply, err := s.runner.ContractRequester.Call(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "CallMethod failed with error")
	}

	re.Reply = *callMethodReply.(*reply.CallMethod)

	var contractErr *foundation.Error
	re.ExtractedReply, contractErr, err = extractor.CallResponse(re.Reply.Result)

	if err != nil {
		return errors.Wrap(err, "Can't extract response")
	}

	if contractErr != nil {
		return errors.Wrap(errors.New(contractErr.S), "Error in called method")
	}

	return nil
}
