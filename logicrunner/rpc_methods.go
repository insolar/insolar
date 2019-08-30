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

package logicrunner

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.ProxyImplementation -o ./ -s _mock.go -g

type ProxyImplementation interface {
	GetCode(context.Context, *common.Transcript, rpctypes.UpGetCodeReq, *rpctypes.UpGetCodeResp) error
	RouteCall(context.Context, *common.Transcript, rpctypes.UpRouteReq, *rpctypes.UpRouteResp) error
	SaveAsChild(context.Context, *common.Transcript, rpctypes.UpSaveAsChildReq, *rpctypes.UpSaveAsChildResp) error
	DeactivateObject(context.Context, *common.Transcript, rpctypes.UpDeactivateObjectReq, *rpctypes.UpDeactivateObjectResp) error
}

type RPCMethods struct {
	ss         StateStorage
	execution  ProxyImplementation
	validation ProxyImplementation
}

func NewRPCMethods(
	am artifacts.Client,
	dc artifacts.DescriptorsCache,
	cr insolar.ContractRequester,
	ss StateStorage,
	outgoingSender OutgoingRequestSender,
) *RPCMethods {
	return &RPCMethods{
		ss:         ss,
		execution:  NewExecutionProxyImplementation(dc, cr, am, outgoingSender),
		validation: NewValidationProxyImplementation(dc),
	}
}

func (m *RPCMethods) getCurrent(
	obj insolar.Reference, mode insolar.CallMode, reqRef insolar.Reference,
) (
	ProxyImplementation, *common.Transcript, error,
) {
	switch mode {
	case insolar.ExecuteCallMode:
		registry := m.ss.GetExecutionRegistry(obj)
		if registry == nil {
			return nil, nil, errors.New("No execution registry in the state")
		}

		transcript := registry.GetActiveTranscript(reqRef)
		if transcript == nil {
			return nil, nil, errors.New("No transcript in the execution registry")
		}

		return m.execution, transcript, nil
	default:
		panic("not implemented")
	}
}

// GetCode is an RPC retrieving a code by its reference
func (m *RPCMethods) GetCode(req rpctypes.UpGetCodeReq, rep *rpctypes.UpGetCodeResp) error {
	impl, current, err := m.getCurrent(req.Callee, req.Mode, req.Request)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch current execution")
	}

	return impl.GetCode(current.Context, current, req, rep)
}

// RouteCall routes call from a contract to a contract through event bus.
func (m *RPCMethods) RouteCall(req rpctypes.UpRouteReq, rep *rpctypes.UpRouteResp) error {
	impl, current, err := m.getCurrent(req.Callee, req.Mode, req.Request)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch current execution")
	}

	return impl.RouteCall(current.Context, current, req, rep)
}

// SaveAsChild is an RPC saving data as memory of a contract as child a parent
func (m *RPCMethods) SaveAsChild(req rpctypes.UpSaveAsChildReq, rep *rpctypes.UpSaveAsChildResp) error {
	impl, current, err := m.getCurrent(req.Callee, req.Mode, req.Request)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch current execution")
	}

	return impl.SaveAsChild(current.Context, current, req, rep)
}

// DeactivateObject is an RPC saving data as memory of a contract as child a parent
func (m *RPCMethods) DeactivateObject(req rpctypes.UpDeactivateObjectReq, rep *rpctypes.UpDeactivateObjectResp) error {
	impl, current, err := m.getCurrent(req.Callee, req.Mode, req.Request)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch current execution")
	}

	return impl.DeactivateObject(current.Context, current, req, rep)
}

type executionProxyImplementation struct {
	dc             artifacts.DescriptorsCache
	cr             insolar.ContractRequester
	am             artifacts.Client
	outgoingSender OutgoingRequestSender
}

func NewExecutionProxyImplementation(
	dc artifacts.DescriptorsCache,
	cr insolar.ContractRequester,
	am artifacts.Client,
	outgoingSender OutgoingRequestSender,
) ProxyImplementation {
	return &executionProxyImplementation{
		dc:             dc,
		cr:             cr,
		am:             am,
		outgoingSender: outgoingSender,
	}
}

func (m *executionProxyImplementation) GetCode(
	ctx context.Context, current *common.Transcript, req rpctypes.UpGetCodeReq, reply *rpctypes.UpGetCodeResp,
) error {
	ctx, span := instracer.StartSpan(ctx, "service.GetCode")
	defer span.End()

	codeDescriptor, err := m.dc.GetCode(ctx, req.Code)
	if err != nil {
		return errors.Wrap(err, "couldn't get code descriptor")
	}
	reply.Code, err = codeDescriptor.Code()
	if err != nil {
		return errors.Wrap(err, "couldn't get code content")
	}
	return nil
}

func (m *executionProxyImplementation) RouteCall(
	ctx context.Context, current *common.Transcript, req rpctypes.UpRouteReq, rep *rpctypes.UpRouteResp,
) error {
	inslogger.FromContext(ctx).Debug(
		"call to others contract method ", req.Method,
		" on object ", req.Object,
	)

	outgoing := buildOutgoingRequest(ctx, current, req)

	// Step 1. Register outgoing request.

	// If pulse changes during registering of OutgoingRequest we don't care because
	// we _already_ are processing the request. We should continue to execute and
	// the next executor will wait for us in pending state. For this reason Flow is not
	// used for registering the outgoing request.
	outReqInfo, err := m.am.RegisterOutgoingRequest(ctx, outgoing)
	if err != nil {
		return err
	}

	if req.Saga {
		// Saga methods are not executed right away. LME will send a method
		// to the VE when current object finishes the execution and validation.
		return nil
	}

	// Step 2. Send the request and register the result (both is done by outgoingSender)

	outgoingReqRef := *outReqInfo.RequestReference()

	var incoming *record.IncomingRequest
	_, rep.Result, incoming, err = m.outgoingSender.SendOutgoingRequest(ctx, outgoingReqRef, outgoing)
	if incoming != nil {
		current.AddOutgoingRequest(ctx, *incoming, rep.Result, nil, err)
	}
	return err
}

// SaveAsChild is an RPC saving data as memory of a contract as child a parent
func (m *executionProxyImplementation) SaveAsChild(
	ctx context.Context, current *common.Transcript, req rpctypes.UpSaveAsChildReq, rep *rpctypes.UpSaveAsChildResp,
) error {
	inslogger.FromContext(ctx).Debug(
		"call to others contract constructor ", req.ConstructorName,
		" on prototype ", req.Prototype.String(),
	)
	ctx, span := instracer.StartSpan(ctx, "RPC.SaveAsChild")
	defer span.End()

	outgoing := buildOutgoingSaveAsChildRequest(ctx, current, req)

	// Register outgoing request
	outReqInfo, err := m.am.RegisterOutgoingRequest(ctx, outgoing)
	if err != nil {
		return err
	}

	// Register result of the outgoing method
	outgoingReqRef := *outReqInfo.RequestReference()

	var incoming *record.IncomingRequest
	rep.Reference, rep.Result, incoming, err = m.outgoingSender.SendOutgoingRequest(ctx, outgoingReqRef, outgoing)
	if incoming != nil {
		current.AddOutgoingRequest(ctx, *incoming, rep.Result, nil, err)
	}
	return err
}

func (m *executionProxyImplementation) DeactivateObject(
	ctx context.Context, current *common.Transcript, req rpctypes.UpDeactivateObjectReq, rep *rpctypes.UpDeactivateObjectResp,
) error {
	inslogger.FromContext(ctx).Debug("contract deactivating itself")

	current.Deactivate = true

	return nil
}

type validationProxyImplementation struct {
	dc artifacts.DescriptorsCache
}

func NewValidationProxyImplementation(
	dc artifacts.DescriptorsCache,
) ProxyImplementation {
	return &validationProxyImplementation{
		dc: dc,
	}
}

func (m *validationProxyImplementation) GetCode(
	ctx context.Context, current *common.Transcript, req rpctypes.UpGetCodeReq, reply *rpctypes.UpGetCodeResp,
) error {
	codeDescriptor, err := m.dc.GetCode(ctx, req.Code)
	if err != nil {
		return errors.Wrap(err, "couldn't get code descriptor")
	}

	reply.Code, err = codeDescriptor.Code()
	if err != nil {
		return errors.Wrap(err, "couldn't get code content")
	}
	return nil
}

func (m *validationProxyImplementation) RouteCall(
	ctx context.Context, current *common.Transcript, req rpctypes.UpRouteReq, rep *rpctypes.UpRouteResp,
) error {
	if current.Request.Immutable {
		return errors.New("immutable method can't make calls")
	}

	outgoing := buildOutgoingRequest(ctx, current, req)
	incoming := buildIncomingRequestFromOutgoing(outgoing)

	reqRes := current.HasOutgoingRequest(ctx, *incoming)
	if reqRes == nil {
		return errors.New("unexpected outgoing call during validation")
	}
	if reqRes.Error != nil {
		return reqRes.Error
	}

	if req.Wait {
		rep.Result = reqRes.Response
	}

	return nil
}

func (m *validationProxyImplementation) SaveAsChild(
	ctx context.Context, current *common.Transcript, req rpctypes.UpSaveAsChildReq, rep *rpctypes.UpSaveAsChildResp,
) error {
	outgoing := buildOutgoingSaveAsChildRequest(ctx, current, req)
	incoming := buildIncomingRequestFromOutgoing(outgoing)

	reqRes := current.HasOutgoingRequest(ctx, *incoming)
	if reqRes == nil {
		return errors.New("unexpected outgoing call during validation")
	}
	if reqRes.Error != nil {
		return reqRes.Error
	}

	rep.Reference = reqRes.NewObject

	return nil
}

func (m *validationProxyImplementation) DeactivateObject(
	ctx context.Context, current *common.Transcript, req rpctypes.UpDeactivateObjectReq, rep *rpctypes.UpDeactivateObjectResp,
) error {

	current.Deactivate = true

	return nil
}

func buildIncomingRequestFromOutgoing(outgoing *record.OutgoingRequest) *record.IncomingRequest {
	// Currently IncomingRequest and OutgoingRequest are almost exact copies of each other
	// thus the following code is a bit ugly. However this will change when we'll
	// figure out which fields are actually needed in OutgoingRequest and which are
	// not. Thus please keep the code the way it is for now, dont't introduce any
	// CommonRequestData structures or something like this.
	// This being said the implementation of Request interface differs for Incoming and
	// OutgoingRequest. See corresponding implementation of the interface methods.
	incoming := record.IncomingRequest{
		Caller:          outgoing.Caller,
		CallerPrototype: outgoing.CallerPrototype,
		Nonce:           outgoing.Nonce,

		Immutable: outgoing.Immutable,

		CallType:  outgoing.CallType, // used only for CTSaveAsChild
		Base:      outgoing.Base,     // used only for CTSaveAsChild
		Object:    outgoing.Object,
		Prototype: outgoing.Prototype,
		Method:    outgoing.Method,
		Arguments: outgoing.Arguments,

		APIRequestID: outgoing.APIRequestID,
		Reason:       outgoing.Reason,
	}

	if outgoing.ReturnMode == record.ReturnSaga {
		// We never wait for a result of saga call
		incoming.ReturnMode = record.ReturnNoWait
	} else {
		// If this is not a saga call just copy the ReturnMode
		incoming.ReturnMode = outgoing.ReturnMode
	}

	return &incoming
}

func buildOutgoingRequest(
	_ context.Context, current *common.Transcript, req rpctypes.UpRouteReq,
) *record.OutgoingRequest {

	current.Nonce++

	outgoing := &record.OutgoingRequest{
		Caller:          req.Callee,
		CallerPrototype: req.CalleePrototype,
		Nonce:           current.Nonce,

		Immutable: req.Immutable,

		Object:    &req.Object,
		Prototype: &req.Prototype,
		Method:    req.Method,
		Arguments: req.Arguments,

		APIRequestID: current.Request.APIRequestID,
		Reason:       current.RequestRef,
	}

	if req.Saga {
		// OutgoingRequest with ReturnMode = ReturnSaga will be called by LME
		// when current object finishes the execution and validation.
		outgoing.ReturnMode = record.ReturnSaga
	} else if !req.Wait {
		outgoing.ReturnMode = record.ReturnNoWait
	}

	return outgoing
}

func buildOutgoingSaveAsChildRequest(
	_ context.Context, current *common.Transcript, req rpctypes.UpSaveAsChildReq,
) *record.OutgoingRequest {

	current.Nonce++

	outgoing := record.OutgoingRequest{
		Caller:          req.Callee,
		CallerPrototype: req.CalleePrototype,
		Nonce:           current.Nonce,

		CallType:  record.CTSaveAsChild,
		Base:      &req.Parent,
		Prototype: &req.Prototype,
		Method:    req.ConstructorName,
		Arguments: req.ArgsSerialized,

		APIRequestID: current.Request.APIRequestID,
		Reason:       current.RequestRef,
	}

	return &outgoing
}
