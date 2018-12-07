/*
 *
 *  *    Copyright 2018 Insolar
 *  *
 *  *    Licensed under the Apache License, Version 2.0 (the "License");
 *  *    you may not use this file except in compliance with the License.
 *  *    You may obtain a copy of the License at
 *  *
 *  *        http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  *    Unless required by applicable law or agreed to in writing, software
 *  *    distributed under the License is distributed on an "AS IS" BASIS,
 *  *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  *    See the License for the specific language governing permissions and
 *  *    limitations under the License.
 *
 */

package logicrunner

import (
	"bytes"
	"context"
	"net"
	"net/rpc"
	"sync/atomic"

	"github.com/satori/go.uuid"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/pkg/errors"
)

// StartRPC starts RPC server for isolated executors to use
func StartRPC(ctx context.Context, lr *LogicRunner) *RPC {
	rpcService := &RPC{lr: lr}

	rpcServer := rpc.NewServer()
	err := rpcServer.Register(rpcService)
	if err != nil {
		panic("Fail to register LogicRunner RPC Service: " + err.Error())
	}

	l, e := net.Listen(lr.Cfg.RPCProtocol, lr.Cfg.RPCListen)
	if e != nil {
		inslogger.FromContext(ctx).Fatal("couldn't setup listener on '"+lr.Cfg.RPCListen+"' over "+lr.Cfg.RPCProtocol+": ", e)
	}
	lr.sock = l

	inslogger.FromContext(ctx).Infof("starting LogicRunner RPC service on %q over %s", lr.Cfg.RPCListen, lr.Cfg.RPCProtocol)
	go func() {
		rpcServer.Accept(l)
		inslogger.FromContext(ctx).Info("LogicRunner RPC service stopped")
	}()

	return rpcService
}

// RPC is a RPC interface for runner to use for various tasks, e.g. code fetching
type RPC struct {
	lr *LogicRunner
}

// GetCode is an RPC retrieving a code by its reference
func (gpr *RPC) GetCode(req rpctypes.UpGetCodeReq, reply *rpctypes.UpGetCodeResp) error {
	es := gpr.lr.UpsertExecution(req.Callee)
	ctx := es.insContext

	am := gpr.lr.ArtifactManager
	codeDescriptor, err := am.GetCode(ctx, req.Code)
	if err != nil {
		return err
	}
	reply.Code, err = codeDescriptor.Code()
	if err != nil {
		return err
	}
	return nil
}

var serial uint64 = 1

// MakeBaseMessage makes base of logicrunner event from base of up request
func MakeBaseMessage(req rpctypes.UpBaseReq) message.BaseLogicMessage {
	return message.BaseLogicMessage{
		Caller:          req.Callee,
		CallerPrototype: req.Prototype,
		Request:         req.Request,
		Nonce:           atomicLoadAndIncrementUint64(&serial),
	}
}

// RouteCall routes call from a contract to a contract through event bus.
func (gpr *RPC) RouteCall(req rpctypes.UpRouteReq, rep *rpctypes.UpRouteResp) error {
	es := gpr.lr.UpsertExecution(req.Callee)
	ctx := es.insContext

	cr, step := gpr.lr.nextValidationStep(req.Callee)
	if step >= 0 { // validate
		if core.CaseRecordTypeRouteCall != cr.Type {
			return errors.New("wrong validation type on RouteCall")
		}
		sig := HashInterface(gpr.lr.PlatformCryptographyScheme, req)
		if !bytes.Equal(cr.ReqSig, sig) {
			return errors.New("wrong validation sig on RouteCall")
		}

		rep.Result = cr.Resp.(core.Arguments)
		return nil
	}

	var mode message.MethodReturnMode
	if req.Wait {
		mode = message.ReturnResult
	} else {
		mode = message.ReturnNoWait
	}

	msg := &message.CallMethod{
		BaseLogicMessage: MakeBaseMessage(req.UpBaseReq),
		ReturnMode:       mode,
		ObjectRef:        req.Object,
		Method:           req.Method,
		Arguments:        req.Arguments,
	}

	res, err := gpr.lr.MessageBus.Send(ctx, msg, nil)
	if err != nil {
		return errors.Wrap(err, "couldn't dispatch event")
	}

	rep.Result = res.(*reply.CallMethod).Result
	gpr.lr.addObjectCaseRecord(req.Callee, core.CaseRecord{
		Type:   core.CaseRecordTypeRouteCall,
		ReqSig: HashInterface(gpr.lr.PlatformCryptographyScheme, req),
		Resp:   rep.Result,
	})

	return nil
}

// SaveAsChild is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) SaveAsChild(req rpctypes.UpSaveAsChildReq, rep *rpctypes.UpSaveAsChildResp) error {
	es := gpr.lr.UpsertExecution(req.Callee)
	ctx := es.insContext

	if gpr.lr.MessageBus == nil {
		return errors.New("event bus was not set during initialization")
	}

	cr, step := gpr.lr.nextValidationStep(req.Callee)
	if step >= 0 { // validate
		if core.CaseRecordTypeSaveAsChild != cr.Type {
			return errors.New("wrong validation type on SaveAsChild")
		}
		sig := HashInterface(gpr.lr.PlatformCryptographyScheme, req)
		if !bytes.Equal(cr.ReqSig, sig) {
			return errors.New("wrong validation sig on SaveAsChild")
		}

		rep.Reference = cr.Resp.(*core.RecordRef)
		return nil
	}

	msg := &message.CallConstructor{
		BaseLogicMessage: MakeBaseMessage(req.UpBaseReq),
		PrototypeRef:     req.Prototype,
		ParentRef:        req.Parent,
		Name:             req.ConstructorName,
		Arguments:        req.ArgsSerialized,
		SaveAs:           message.Child,
	}

	res, err := gpr.lr.MessageBus.Send(ctx, msg, nil)
	if err != nil {
		return errors.Wrap(err, "couldn't save new object as child")
	}

	rep.Reference = res.(*reply.CallConstructor).Object

	gpr.lr.addObjectCaseRecord(req.Callee, core.CaseRecord{
		Type:   core.CaseRecordTypeSaveAsChild,
		ReqSig: HashInterface(gpr.lr.PlatformCryptographyScheme, req),
		Resp:   rep.Reference,
	})

	return nil
}

var iteratorMap = make(map[string]*core.RefIterator)
var iteratorBuffSize = 1000

// GetObjChildrenIterator is an RPC returns set of object children
func (gpr *RPC) GetObjChildrenIterator(req rpctypes.UpGetObjChildrenIteratorReq, rep *rpctypes.UpGetObjChildrenIteratorResp) error {
	es := gpr.lr.UpsertExecution(req.Callee)
	ctx := es.insContext

	cr, step := gpr.lr.nextValidationStep(req.Callee)
	if step >= 0 { // validate
		if core.CaseRecordTypeGetObjChildrenIterator != cr.Type {
			return errors.New("wrong validation type on GetObjChildrenIterator")
		}
		sig := HashInterface(gpr.lr.PlatformCryptographyScheme, req)
		if !bytes.Equal(cr.ReqSig, sig) {
			return errors.New("wrong validation sig on GetObjChildrenIterator")
		}

		rep.Iterator = cr.Resp.(rpctypes.ChildIterator)
		return nil
	}
	am := gpr.lr.ArtifactManager
	iteratorID := req.IteratorID
	if _, ok := iteratorMap[iteratorID]; !ok {
		i, err := am.GetChildren(ctx, req.Obj, nil)
		if err != nil {
			return errors.Wrap(err, "[ GetObjChildrenIterator ] Can't get children")
		}

		id, err := uuid.NewV4()
		if err != nil {
			return errors.Wrap(err, "[ GetObjChildrenIterator ] Can't generate UUID")
		}

		iteratorID = id.String()
		iteratorMap[iteratorID] = &i
	}

	i := *iteratorMap[iteratorID]
	rep.Iterator.ID = iteratorID
	rep.Iterator.CanFetch = i.HasNext()
	for len(rep.Iterator.Buff) < iteratorBuffSize && i.HasNext() {
		r, err := i.Next()
		if err != nil {
			return errors.Wrap(err, "[ GetObjChildrenIterator ] Can't get Next")
		}
		rep.Iterator.CanFetch = i.HasNext()

		o, err := am.GetObject(ctx, *r, nil, false)

		if err != nil {
			if err == core.ErrDeactivated {
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

	if !i.HasNext() {
		delete(iteratorMap, rep.Iterator.ID)
	}

	gpr.lr.addObjectCaseRecord(req.Callee, core.CaseRecord{ // bad idea, we can store gadzillion of children
		Type:   core.CaseRecordTypeGetObjChildrenIterator,
		ReqSig: HashInterface(gpr.lr.PlatformCryptographyScheme, req),
		Resp:   rep.Iterator,
	})
	return nil
}

// SaveAsDelegate is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) SaveAsDelegate(req rpctypes.UpSaveAsDelegateReq, rep *rpctypes.UpSaveAsDelegateResp) error {
	es := gpr.lr.UpsertExecution(req.Callee)
	ctx := es.insContext

	cr, step := gpr.lr.nextValidationStep(req.Callee)
	if step >= 0 { // validate
		if core.CaseRecordTypeSaveAsDelegate != cr.Type {
			return errors.New("wrong validation type on SaveAsDelegate")
		}
		sig := HashInterface(gpr.lr.PlatformCryptographyScheme, req)
		if !bytes.Equal(cr.ReqSig, sig) {
			return errors.New("wrong validation sig on SaveAsDelegate")
		}

		rep.Reference = cr.Resp.(*core.RecordRef)
		return nil
	}

	msg := &message.CallConstructor{
		BaseLogicMessage: MakeBaseMessage(req.UpBaseReq),
		PrototypeRef:     req.Prototype,
		ParentRef:        req.Into,
		Name:             req.ConstructorName,
		Arguments:        req.ArgsSerialized,
		SaveAs:           message.Delegate,
	}

	res, err := gpr.lr.MessageBus.Send(ctx, msg, nil)

	if err != nil {
		return errors.Wrap(err, "couldn't save new object as delegate")
	}

	rep.Reference = res.(*reply.CallConstructor).Object
	gpr.lr.addObjectCaseRecord(req.Callee, core.CaseRecord{
		Type:   core.CaseRecordTypeSaveAsDelegate,
		ReqSig: HashInterface(gpr.lr.PlatformCryptographyScheme, req),
		Resp:   rep.Reference,
	})

	return nil
}

// GetDelegate is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) GetDelegate(req rpctypes.UpGetDelegateReq, rep *rpctypes.UpGetDelegateResp) error {
	es := gpr.lr.UpsertExecution(req.Callee)
	ctx := es.insContext

	cr, step := gpr.lr.nextValidationStep(req.Callee)
	if step >= 0 { // validate
		if core.CaseRecordTypeGetDelegate != cr.Type {
			return errors.New("wrong validation type on RouteCall")
		}
		sig := HashInterface(gpr.lr.PlatformCryptographyScheme, req)
		if !bytes.Equal(cr.ReqSig, sig) {
			return errors.New("wrong validation sig on RouteCall")
		}

		rep.Object = cr.Resp.(core.RecordRef)
		return nil
	}
	am := gpr.lr.ArtifactManager
	ref, err := am.GetDelegate(ctx, req.Object, req.OfType)
	if err != nil {
		return err
	}
	rep.Object = *ref
	gpr.lr.addObjectCaseRecord(req.Callee, core.CaseRecord{
		Type:   core.CaseRecordTypeGetDelegate,
		ReqSig: HashInterface(gpr.lr.PlatformCryptographyScheme, req),
		Resp:   rep.Object,
	})
	return nil
}

// DeactivateObject is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) DeactivateObject(req rpctypes.UpDeactivateObjectReq, rep *rpctypes.UpDeactivateObjectResp) error {
	state := gpr.lr.GetExecution(req.Callee)
	if state == nil {
		return errors.New("no execution state, impossible, shouldn't be")
	}

	// TODO: is it race? make sure it's not!
	state.deactivate = true

	return nil
}

// atomicLoadAndIncrementUint64 performs CAS loop, increments counter and returns old value.
func atomicLoadAndIncrementUint64(addr *uint64) uint64 {
	for {
		val := atomic.LoadUint64(addr)
		if atomic.CompareAndSwapUint64(addr, val, val+1) {
			return val
		}
	}
}
