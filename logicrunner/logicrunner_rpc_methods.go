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
	"sync"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

type ProxyImplementation interface {
	GetCode(context.Context, *Transcript, rpctypes.UpGetCodeReq, *rpctypes.UpGetCodeResp) error
	RouteCall(context.Context, *Transcript, rpctypes.UpRouteReq, *rpctypes.UpRouteResp) error
	SaveAsChild(context.Context, *Transcript, rpctypes.UpSaveAsChildReq, *rpctypes.UpSaveAsChildResp) error
	SaveAsDelegate(context.Context, *Transcript, rpctypes.UpSaveAsDelegateReq, *rpctypes.UpSaveAsDelegateResp) error
	GetObjChildrenIterator(context.Context, *Transcript, rpctypes.UpGetObjChildrenIteratorReq, *rpctypes.UpGetObjChildrenIteratorResp) error
	GetDelegate(context.Context, *Transcript, rpctypes.UpGetDelegateReq, *rpctypes.UpGetDelegateResp) error
	DeactivateObject(context.Context, *Transcript, rpctypes.UpDeactivateObjectReq, *rpctypes.UpDeactivateObjectResp) error
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
) *RPCMethods {
	return &RPCMethods{
		ss:         ss,
		execution:  NewExecutionProxyImplementation(dc, cr, am),
		validation: NewValidationProxyImplementation(),
	}
}

func (m *RPCMethods) getCurrent(
	obj insolar.Reference, mode insolar.CallMode, reqRef insolar.Reference,
) (
	ProxyImplementation, *Transcript, error,
) {
	os := m.ss.GetObjectState(obj)
	if os == nil {
		return nil, nil, errors.New("Failed to find requested object state. ref: " + obj.String())
	}
	switch mode {
	case insolar.ExecuteCallMode:
		es := os.ExecutionState
		if es == nil {
			return nil, nil, errors.New("No execution in the state")
		}

		cur := es.CurrentList.Get(reqRef)
		if cur == nil {
			return nil, nil, errors.New("No current execution in the state for request " + reqRef.String())
		}

		return m.execution, cur, nil
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

// SaveAsDelegate is an RPC saving data as memory of a contract as child a parent
func (m *RPCMethods) SaveAsDelegate(req rpctypes.UpSaveAsDelegateReq, rep *rpctypes.UpSaveAsDelegateResp) error {
	impl, current, err := m.getCurrent(req.Callee, req.Mode, req.Request)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch current execution")
	}
	return impl.SaveAsDelegate(current.Context, current, req, rep)
}

// GetObjChildrenIterator is an RPC returns an iterator over object children with specified prototype
func (m *RPCMethods) GetObjChildrenIterator(
	req rpctypes.UpGetObjChildrenIteratorReq,
	rep *rpctypes.UpGetObjChildrenIteratorResp,
) error {
	impl, current, err := m.getCurrent(req.Callee, req.Mode, req.Request)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch current execution")
	}

	return impl.GetObjChildrenIterator(current.Context, current, req, rep)
}

// GetDelegate is an RPC saving data as memory of a contract as child a parent
func (m *RPCMethods) GetDelegate(req rpctypes.UpGetDelegateReq, rep *rpctypes.UpGetDelegateResp) error {
	impl, current, err := m.getCurrent(req.Callee, req.Mode, req.Request)
	if err != nil {
		return errors.Wrap(err, "Failed to fetch current execution")
	}

	return impl.GetDelegate(current.Context, current, req, rep)
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
	dc artifacts.DescriptorsCache
	cr insolar.ContractRequester
	am artifacts.Client
}

func NewExecutionProxyImplementation(
	dc artifacts.DescriptorsCache,
	cr insolar.ContractRequester,
	am artifacts.Client,
) ProxyImplementation {
	return &executionProxyImplementation{
		dc: dc,
		cr: cr,
		am: am,
	}
}

func (m *executionProxyImplementation) GetCode(
	ctx context.Context, current *Transcript, req rpctypes.UpGetCodeReq, reply *rpctypes.UpGetCodeResp,
) error {
	ctx, span := instracer.StartSpan(ctx, "service.GetCode")
	defer span.End()

	codeDescriptor, err := m.dc.GetCode(ctx, req.Code)
	if err != nil {
		return err
	}
	reply.Code, err = codeDescriptor.Code()
	if err != nil {
		return err
	}
	return nil
}

// RouteCall routes call from a contract to a contract through event bus.
func (m *executionProxyImplementation) RouteCall(
	ctx context.Context, current *Transcript, req rpctypes.UpRouteReq, rep *rpctypes.UpRouteResp,
) error {
	inslogger.FromContext(ctx).Debug("RPC.RouteCall")

	if current.LogicContext.Immutable {
		return errors.New("Try to call route from immutable method")
	}

	// TODO: delegation token

	current.Nonce++

	reqRecord := record.IncomingRequest{
		Caller:          req.Callee,
		CallerPrototype: req.CalleePrototype,
		Nonce:           current.Nonce,

		Immutable: req.Immutable,

		Object:    &req.Object,
		Prototype: &req.Prototype,
		Method:    req.Method,
		Arguments: req.Arguments,

		APIRequestID: current.Request.APIRequestID,
		Reason:       *current.RequestRef,
	}

	if !req.Wait {
		reqRecord.ReturnMode = record.ReturnNoWait
	}

	msg := &message.CallMethod{IncomingRequest: reqRecord}
	res, err := m.cr.CallMethod(ctx, msg)
	current.AddOutgoingRequest(ctx, reqRecord, rep.Result, nil, err)
	if err != nil {
		return err
	}

	if req.Wait {
		rep.Result = res.(*reply.CallMethod).Result
	}

	return nil
}

// SaveAsChild is an RPC saving data as memory of a contract as child a parent
func (m *executionProxyImplementation) SaveAsChild(
	ctx context.Context, current *Transcript, req rpctypes.UpSaveAsChildReq, rep *rpctypes.UpSaveAsChildResp,
) error {
	inslogger.FromContext(ctx).Debug("RPC.SaveAsChild")
	ctx, span := instracer.StartSpan(ctx, "RPC.SaveAsChild")
	defer span.End()

	current.Nonce++

	reqRecord := record.IncomingRequest{
		Caller:          req.Callee,
		CallerPrototype: req.CalleePrototype,
		Nonce:           current.Nonce,

		CallType:  record.CTSaveAsChild,
		Base:      &req.Parent,
		Prototype: &req.Prototype,
		Method:    req.ConstructorName,
		Arguments: req.ArgsSerialized,

		APIRequestID: current.Request.APIRequestID,
		Reason:       *current.RequestRef,
	}

	msg := &message.CallMethod{IncomingRequest: reqRecord}

	ref, err := m.cr.CallConstructor(ctx, msg)
	current.AddOutgoingRequest(ctx, reqRecord, nil, ref, err)

	rep.Reference = ref

	return err
}

// SaveAsDelegate is an RPC saving data as memory of a contract as child a parent
func (m *executionProxyImplementation) SaveAsDelegate(
	ctx context.Context, current *Transcript, req rpctypes.UpSaveAsDelegateReq, rep *rpctypes.UpSaveAsDelegateResp,
) error {
	inslogger.FromContext(ctx).Debug("RPC.SaveAsDelegate")
	ctx, span := instracer.StartSpan(ctx, "RPC.SaveAsDelegate")
	defer span.End()

	current.Nonce++
	reqRecord := record.IncomingRequest{
		Caller:          req.Callee,
		CallerPrototype: req.CalleePrototype,
		Nonce:           current.Nonce,

		CallType:  record.CTSaveAsDelegate,
		Base:      &req.Into,
		Prototype: &req.Prototype,
		Method:    req.ConstructorName,
		Arguments: req.ArgsSerialized,

		APIRequestID: current.Request.APIRequestID,
		Reason:       *current.RequestRef,
	}
	msg := &message.CallMethod{IncomingRequest: reqRecord}

	ref, err := m.cr.CallConstructor(ctx, msg)
	current.AddOutgoingRequest(ctx, reqRecord, nil, ref, err)

	rep.Reference = ref
	return err
}

var iteratorBuffSize = 1000
var iteratorMap = make(map[string]artifacts.RefIterator)
var iteratorMapLock = sync.RWMutex{}

// GetObjChildrenIterator is an RPC returns an iterator over object children with specified prototype
func (m *executionProxyImplementation) GetObjChildrenIterator(
	ctx context.Context, current *Transcript,
	req rpctypes.UpGetObjChildrenIteratorReq,
	rep *rpctypes.UpGetObjChildrenIteratorResp,
) error {
	ctx, span := instracer.StartSpan(ctx, "RPC.GetObjChildrenIterator")
	defer span.End()

	iteratorID := req.IteratorID

	iteratorMapLock.RLock()
	iterator, ok := iteratorMap[iteratorID]
	iteratorMapLock.RUnlock()

	if !ok {
		newIterator, err := m.am.GetChildren(ctx, req.Object, nil)
		if err != nil {
			return errors.Wrap(err, "[ GetObjChildrenIterator ] Can't get children")
		}

		id, err := uuid.NewV4()
		if err != nil {
			return errors.Wrap(err, "[ GetObjChildrenIterator ] Can't generate UUID")
		}

		iteratorID = id.String()

		iteratorMapLock.Lock()
		iterator, ok = iteratorMap[iteratorID]
		if !ok {
			iteratorMap[iteratorID] = newIterator
			iterator = newIterator
		}
		iteratorMapLock.Unlock()
	}

	iter := iterator

	rep.Iterator.ID = iteratorID
	rep.Iterator.CanFetch = iter.HasNext()
	for len(rep.Iterator.Buff) < iteratorBuffSize && iter.HasNext() {
		r, err := iter.Next()
		if err != nil {
			return errors.Wrap(err, "[ GetObjChildrenIterator ] Can't get Next")
		}
		rep.Iterator.CanFetch = iter.HasNext()

		o, err := m.am.GetObject(ctx, *r)

		if err != nil {
			if err == insolar.ErrDeactivated {
				continue
			}
			return errors.Wrap(err, "[ GetObjChildrenIterator ] Can't call GetObject on Next")
		}
		protoRef, err := o.Prototype()
		if err != nil {
			return errors.Wrap(err, "[ GetObjChildrenIterator ] Can't get prototype reference")
		}

		if protoRef.Equal(req.Prototype) {
			rep.Iterator.Buff = append(rep.Iterator.Buff, *r)
		}
	}

	if !iter.HasNext() {
		iteratorMapLock.Lock()
		delete(iteratorMap, rep.Iterator.ID)
		iteratorMapLock.Unlock()
	}

	return nil
}

// GetDelegate is an RPC saving data as memory of a contract as child a parent
func (m *executionProxyImplementation) GetDelegate(
	ctx context.Context, current *Transcript, req rpctypes.UpGetDelegateReq, rep *rpctypes.UpGetDelegateResp,
) error {
	ref, err := m.am.GetDelegate(ctx, req.Object, req.OfType)
	if err != nil {
		return err
	}
	rep.Object = *ref
	return nil
}

// DeactivateObject is an RPC saving data as memory of a contract as child a parent
func (m *executionProxyImplementation) DeactivateObject(
	ctx context.Context, current *Transcript, req rpctypes.UpDeactivateObjectReq, rep *rpctypes.UpDeactivateObjectResp,
) error {

	current.Deactivate = true

	return nil
}

type validationProxyImplementation struct {
}

func NewValidationProxyImplementation() ProxyImplementation {
	return &validationProxyImplementation{}
}

func (m *validationProxyImplementation) GetCode(
	ctx context.Context, current *Transcript, req rpctypes.UpGetCodeReq, reply *rpctypes.UpGetCodeResp,
) error {
	panic("implement me")
}

// RouteCall routes call from a contract to a contract through event bus.
func (m *validationProxyImplementation) RouteCall(
	ctx context.Context, current *Transcript, req rpctypes.UpRouteReq, rep *rpctypes.UpRouteResp,
) error {
	panic("implement me")
}

// SaveAsChild is an RPC saving data as memory of a contract as child a parent
func (m *validationProxyImplementation) SaveAsChild(
	ctx context.Context, current *Transcript, req rpctypes.UpSaveAsChildReq, rep *rpctypes.UpSaveAsChildResp,
) error {
	panic("implement me")
}

// SaveAsDelegate is an RPC saving data as memory of a contract as child a parent
func (m *validationProxyImplementation) SaveAsDelegate(
	ctx context.Context, current *Transcript, req rpctypes.UpSaveAsDelegateReq, rep *rpctypes.UpSaveAsDelegateResp,
) error {
	panic("implement me")
}

// GetObjChildrenIterator is an RPC returns an iterator over object children with specified prototype
func (m *validationProxyImplementation) GetObjChildrenIterator(
	ctx context.Context, current *Transcript,
	req rpctypes.UpGetObjChildrenIteratorReq,
	rep *rpctypes.UpGetObjChildrenIteratorResp,
) error {
	panic("implement me")
}

// GetDelegate is an RPC saving data as memory of a contract as child a parent
func (m *validationProxyImplementation) GetDelegate(
	ctx context.Context, current *Transcript, req rpctypes.UpGetDelegateReq, rep *rpctypes.UpGetDelegateResp,
) error {
	panic("implement me")
}

// DeactivateObject is an RPC saving data as memory of a contract as child a parent
func (m *validationProxyImplementation) DeactivateObject(
	ctx context.Context, current *Transcript, req rpctypes.UpDeactivateObjectReq, rep *rpctypes.UpDeactivateObjectResp,
) error {
	panic("implement me")
}
