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

// Package builtin is implementation of builtin contracts engine
package builtin

import (
	"context"
	"errors"

	"github.com/tylerb/gls"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	lrCommon "github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

type LogicRunnerRPCStub interface {
	GetCode(rpctypes.UpGetCodeReq, *rpctypes.UpGetCodeResp) error
	RouteCall(rpctypes.UpRouteReq, *rpctypes.UpRouteResp) error
	SaveAsChild(rpctypes.UpSaveAsChildReq, *rpctypes.UpSaveAsChildResp) error
	SaveAsDelegate(rpctypes.UpSaveAsDelegateReq, *rpctypes.UpSaveAsDelegateResp) error
	GetObjChildrenIterator(rpctypes.UpGetObjChildrenIteratorReq, *rpctypes.UpGetObjChildrenIteratorResp) error
	GetDelegate(rpctypes.UpGetDelegateReq, *rpctypes.UpGetDelegateResp) error
	DeactivateObject(rpctypes.UpDeactivateObjectReq, *rpctypes.UpDeactivateObjectResp) error
}

// BuiltIn is a contract runner engine
type BuiltIn struct {
	// Prototype -> Code + Versions
	// PrototypeRegistry    map[string]preprocessor.ContractWrapper
	// PrototypeRefRegistry map[insolar.Reference]string
	// Code ->
	CodeRegistry    map[string]insolar.ContractWrapper
	CodeRefRegistry map[insolar.Reference]string
}

// NewBuiltIn is an constructor
func NewBuiltIn(am artifacts.Client, stub LogicRunnerRPCStub) *BuiltIn {
	codeDescriptors := InitializeCodeDescriptors()
	for _, codeDescriptor := range codeDescriptors {
		am.InjectCodeDescriptor(*codeDescriptor.Ref(), codeDescriptor)
	}

	prototypeDescriptors := InitializePrototypeDescriptors()
	for _, prototypeDescriptor := range prototypeDescriptors {
		am.InjectObjectDescriptor(*prototypeDescriptor.HeadRef(), prototypeDescriptor)
	}

	lrCommon.CurrentProxyCtx = NewProxyHelper(stub)

	return &BuiltIn{
		CodeRefRegistry: InitializeCodeRefs(),
		CodeRegistry:    InitializeContractMethods(),
	}
}

func (b *BuiltIn) CallConstructor(ctx context.Context, callCtx *insolar.LogicCallContext, codeRef insolar.Reference,
	name string, args insolar.Arguments) ([]byte, error) {

	ctx, span := instracer.StartSpan(ctx, "builtin.CallConstructor")
	defer span.End()

	gls.Set("callCtx", callCtx)
	defer gls.Cleanup()

	contractName, ok := b.CodeRefRegistry[codeRef]
	if !ok {
		return nil, errors.New("failed to find contract with reference")
	}
	contract := b.CodeRegistry[contractName]

	constructorFunc, ok := contract.Constructors[name]
	if !ok {
		return nil, errors.New("failed to find contracts method")
	}

	return constructorFunc(args)
}

func (b *BuiltIn) CallMethod(ctx context.Context, callCtx *insolar.LogicCallContext, codeRef insolar.Reference,
	data []byte, method string, args insolar.Arguments) ([]byte, insolar.Arguments, error) {

	// TODO AALEKSEEV check SystemError here

	ctx, span := instracer.StartSpan(ctx, "builtin.CallMethod")
	defer span.End()

	gls.Set("callCtx", callCtx)
	defer gls.Cleanup()

	contractName, ok := b.CodeRefRegistry[codeRef]
	if !ok {
		return nil, nil, errors.New("failed to find contract with reference")
	}
	contract := b.CodeRegistry[contractName]

	methodFunc, ok := contract.Methods[method]
	if !ok {
		return nil, nil, errors.New("failed to find contracts method")
	}

	return methodFunc(data, args)
}
