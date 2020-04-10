// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// Package builtin is implementation of builtin contracts engine
package builtin

import (
	"context"
	"errors"
	"time"

	"github.com/insolar/insolar/applicationbase/builtin"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	lrCommon "github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/insolar/insolar/logicrunner/metrics"
)

type LogicRunnerRPCStub interface {
	GetCode(rpctypes.UpGetCodeReq, *rpctypes.UpGetCodeResp) error
	RouteCall(rpctypes.UpRouteReq, *rpctypes.UpRouteResp) error
	SaveAsChild(rpctypes.UpSaveAsChildReq, *rpctypes.UpSaveAsChildResp) error
	DeactivateObject(rpctypes.UpDeactivateObjectReq, *rpctypes.UpDeactivateObjectResp) error
}

// BuiltIn is a contract runner engine
type BuiltIn struct {
	// Prototype -> Code + Versions
	// PrototypeRegistry    map[string]preprocessor.ContractWrapper
	// PrototypeRefRegistry map[insolar.Reference]string
	// Code ->
	CodeRegistry         map[string]insolar.ContractWrapper
	CodeRefRegistry      map[insolar.Reference]string
	PrototypeRefRegistry map[insolar.Reference]string
}

type BuiltinContracts struct {
	CodeRegistry         map[string]insolar.ContractWrapper
	CodeRefRegistry      map[insolar.Reference]string
	CodeDescriptors      []artifacts.CodeDescriptor
	PrototypeDescriptors []artifacts.PrototypeDescriptor
}

// NewBuiltIn is an constructor
func NewBuiltIn(
	am artifacts.Client, stub LogicRunnerRPCStub, builtinContracts BuiltinContracts,
) *BuiltIn {
	fullCodeDescriptors := append(builtin.InitializeCodeDescriptors(), builtinContracts.CodeDescriptors...)
	for _, codeDescriptor := range fullCodeDescriptors {
		am.InjectCodeDescriptor(*codeDescriptor.Ref(), codeDescriptor)
	}

	fullPrototypeDescriptors := append(builtin.InitializePrototypeDescriptors(), builtinContracts.PrototypeDescriptors...)
	for _, prototypeDescriptor := range fullPrototypeDescriptors {
		am.InjectPrototypeDescriptor(*prototypeDescriptor.HeadRef(), prototypeDescriptor)
	}

	lrCommon.CurrentProxyCtx = NewProxyHelper(stub)

	fullCodeRefRegistry := builtin.InitializeCodeRefs()
	for k, v := range builtinContracts.CodeRefRegistry {
		fullCodeRefRegistry[k] = v
	}
	fullCodeRegistry := builtin.InitializeContractMethods()
	for k, v := range builtinContracts.CodeRegistry {
		fullCodeRegistry[k] = v
	}
	return &BuiltIn{
		CodeRefRegistry: fullCodeRefRegistry,
		CodeRegistry:    fullCodeRegistry,
	}
}

func (b *BuiltIn) CallConstructor(
	ctx context.Context,
	callCtx *insolar.LogicCallContext,
	codeRef insolar.Reference,
	name string,
	args insolar.Arguments,
) (
	[]byte, insolar.Arguments, error,
) {
	executeStart := time.Now()
	ctx = insmetrics.InsertTag(ctx, metrics.TagContractPrototype, b.PrototypeRefRegistry[codeRef])
	ctx = insmetrics.InsertTag(ctx, metrics.TagContractMethodName, "Constructor")

	defer func(ctx context.Context) {
		executionTime := time.Since(executeStart).Nanoseconds()
		stats.Record(ctx, metrics.ContractExecutionTime.M(float64(executionTime)/1e6))
	}(ctx)

	ctx, span := instracer.StartSpan(ctx, "builtin.CallConstructor")
	defer span.Finish()

	foundation.SetLogicalContext(callCtx)
	defer foundation.ClearContext()

	contractName, ok := b.CodeRefRegistry[codeRef]
	if !ok {
		return nil, nil, errors.New("failed to find contract with reference")
	}
	contract := b.CodeRegistry[contractName]

	constructorFunc, ok := contract.Constructors[name]
	if !ok {
		return nil, nil, errors.New("failed to find contracts method")
	}

	objRef := insolar.NewReference(*callCtx.Request.GetLocal())
	return constructorFunc(*objRef, args)
}

func (b *BuiltIn) CallMethod(
	ctx context.Context,
	callCtx *insolar.LogicCallContext,
	codeRef insolar.Reference,
	data []byte,
	method string,
	args insolar.Arguments,
) (
	[]byte, insolar.Arguments, error,
) {
	executeStart := time.Now()
	ctx = insmetrics.InsertTag(ctx, metrics.TagContractPrototype, b.PrototypeRefRegistry[codeRef])
	ctx = insmetrics.InsertTag(ctx, metrics.TagContractMethodName, method)

	defer func(ctx context.Context) {
		executionTime := time.Since(executeStart).Nanoseconds()
		stats.Record(ctx, metrics.ContractExecutionTime.M(float64(executionTime)/1e6))
	}(ctx)

	ctx, span := instracer.StartSpan(ctx, "builtin.CallMethod")
	defer span.Finish()

	foundation.SetLogicalContext(callCtx)
	defer foundation.ClearContext()

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
