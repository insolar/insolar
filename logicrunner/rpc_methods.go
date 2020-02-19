// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logicrunner

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
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

func getRequestReference(info *payload.RequestInfo) *insolar.Reference {
	return insolar.NewRecordReference(info.RequestID)
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
	ctx = instracer.WithParentSpan(ctx, instracer.TraceSpan{
		TraceID: []byte(inslogger.TraceID(ctx)),
		SpanID:  instracer.MakeBinarySpan(current.Request.Reason.Bytes()),
	})

	ctx, span := instracer.StartSpan(ctx, "service.GetCode")
	defer span.Finish()

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
	logger := inslogger.FromContext(ctx)
	logger.Debug(
		"call to others contract method ", req.Method,
		" on object ", req.Object,
	)

	logger = logger.WithFields(map[string]interface{}{
		"call_to":   req.Method,
		"on_object": req.Object,
	})

	outgoing := buildOutgoingRequest(ctx, current, req)

	// Step 1. Register outgoing request.

	logger.Debug("registering outgoing request")

	ctx = instracer.WithParentSpan(ctx, instracer.TraceSpan{
		TraceID: []byte(inslogger.TraceID(ctx)),
		SpanID:  instracer.MakeBinarySpan(current.Request.Reason.Bytes()),
	})

	// If pulse changes during registering of OutgoingRequest we don't care because
	// we _already_ are processing the request. We should continue to execute and
	// the next executor will wait for us in pending state. For this reason Flow is not
	// used for registering the outgoing request.
	outReqInfo, err := m.am.RegisterOutgoingRequest(ctx, outgoing)
	if err != nil {
		return err
	}

	logger.Debug("registered outgoing request")

	if req.Saga {
		// Saga methods are not executed right away. LME will send a method
		// to the VE when current object finishes the execution and validation.
		if outReqInfo.Result != nil {
			return errors.New("RegisterOutgoingRequest returns Result for Saga Call")
		}
		return nil
	}

	// if we replay abandoned request after node was down we can already have Result
	if outReqInfo.Result != nil {
		returns, err := unwrapResult(ctx, outReqInfo.Result)
		if err != nil {
			return errors.Wrap(err, "couldn't unwrap result from ledger")
		}
		rep.Result = returns
		return nil
	}

	logger.Debug("sending outgoing request")

	// Step 2. Send the request and register the result (both is done by outgoingSender)
	rep.Result, _, err = m.outgoingSender.SendOutgoingRequest(ctx, *getRequestReference(outReqInfo), outgoing)
	if err != nil {
		err = errors.Wrap(err, "failed to send outgoing request")
		logger.Error(err)
		return err
	}

	logger.Debug("got result of outgoing request")

	return nil
}

// SaveAsChild is an RPC saving data as memory of a contract as child a parent
func (m *executionProxyImplementation) SaveAsChild(
	ctx context.Context, current *common.Transcript, req rpctypes.UpSaveAsChildReq, rep *rpctypes.UpSaveAsChildResp,
) error {
	logger := inslogger.FromContext(ctx)
	logger.Debug(
		"call to others contract constructor ", req.ConstructorName,
		" on prototype ", req.Prototype.String(),
	)

	ctx = instracer.WithParentSpan(ctx, instracer.TraceSpan{
		TraceID: []byte(inslogger.TraceID(ctx)),
		SpanID:  instracer.MakeBinarySpan(current.Request.Reason.Bytes()),
	})

	ctx, span := instracer.StartSpan(ctx, "RPC.SaveAsChild")
	defer span.Finish()

	logger = logger.WithFields(map[string]interface{}{
		"call_to":   req.ConstructorName,
		"on_object": req.Prototype,
	})

	outgoing := buildOutgoingSaveAsChildRequest(ctx, current, req)

	logger.Debug("registering outgoing request")

	// Register outgoing request
	outReqInfo, err := m.am.RegisterOutgoingRequest(ctx, outgoing)
	if err != nil {
		return err
	}

	logger.Debug("registered outgoing request")

	// if we replay abandoned request after node was down we can already have Result
	if outReqInfo.Result != nil {
		returns, err := unwrapResult(ctx, outReqInfo.Result)
		if err != nil {
			return errors.Wrap(err, "couldn't unwrap result from ledger")
		}
		rep.Result = returns
		return nil
	}

	// Register result of the outgoing method
	outgoingReqRef := *getRequestReference(outReqInfo)

	logger.Debug("sending outgoing request")

	var incoming *record.IncomingRequest
	rep.Result, incoming, err = m.outgoingSender.SendOutgoingRequest(ctx, outgoingReqRef, outgoing)
	if incoming != nil {
		current.AddOutgoingRequest(ctx, *incoming, rep.Result, err)
	}

	logger.Debug("got result of outgoing request")

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

	rep.Result = reqRes.Response

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
	apiReqID := outgoing.APIRequestID

	if outgoing.ReturnMode == record.ReturnSaga {
		apiReqID += fmt.Sprintf("-saga-%d", outgoing.Nonce)
	}

	incoming := record.IncomingRequest{
		Caller:          outgoing.Caller,
		CallerPrototype: outgoing.CallerPrototype,
		Nonce:           outgoing.Nonce,

		Immutable:  outgoing.Immutable,
		ReturnMode: outgoing.ReturnMode,

		CallType:  outgoing.CallType, // used only for CTSaveAsChild
		Base:      outgoing.Base,     // used only for CTSaveAsChild
		Object:    outgoing.Object,
		Prototype: outgoing.Prototype,
		Method:    outgoing.Method,
		Arguments: outgoing.Arguments,

		APIRequestID: apiReqID,
		Reason:       outgoing.Reason,
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

func unwrapResult(
	_ context.Context, materialBlob []byte,
) ([]byte, error) {
	rec := record.Material{}
	err := rec.Unmarshal(materialBlob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal existing result")
	}
	virtual := record.Unwrap(&rec.Virtual)
	resultRecord, ok := virtual.(*record.Result)
	if !ok {
		return nil, fmt.Errorf("unexpected record %T", virtual)
	}
	return resultRecord.Payload, nil
}
