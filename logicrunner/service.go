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
	"net"
	"net/rpc"
	"sync/atomic"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/pkg/errors"
)

// StartRPC starts RPC server for isolated executors to use
func StartRPC(lr *LogicRunner) *RPC {
	rpcService := &RPC{lr: lr}

	rpcServer := rpc.NewServer()
	err := rpcServer.Register(rpcService)
	if err != nil {
		panic("Fail to register LogicRunner RPC Service: " + err.Error())
	}

	l, e := net.Listen(lr.Cfg.RPCProtocol, lr.Cfg.RPCListen)
	if e != nil {
		log.Fatal("couldn't setup listener on '"+lr.Cfg.RPCListen+"' over "+lr.Cfg.RPCProtocol+": ", e)
	}
	lr.sock = l

	log.Infof("starting LogicRunner RPC service on %q over %s", lr.Cfg.RPCListen, lr.Cfg.RPCProtocol)
	go func() {
		rpcServer.Accept(l)
		log.Info("LogicRunner RPC service stopped")
	}()

	return rpcService
}

// RPC is a RPC interface for runner to use for various tasks, e.g. code fetching
type RPC struct {
	lr *LogicRunner
}

// GetCode is an RPC retrieving a code by its reference
func (gpr *RPC) GetCode(req rpctypes.UpGetCodeReq, reply *rpctypes.UpGetCodeResp) error {
	am := gpr.lr.ArtifactManager
	codeDescriptor, err := am.GetCode(req.Code)
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
		Caller:  req.Callee,
		Request: req.Request,
		Nonce:   atomicLoadAndIncrementUint64(&serial),
	}
}

// RouteCall routes call from a contract to a contract through event bus.
func (gpr *RPC) RouteCall(req rpctypes.UpRouteReq, rep *rpctypes.UpRouteResp) error {
	cr, step := gpr.lr.getNextValidationStep(req.Callee)
	if step >= 0 { // validate
		if core.CaseRecordTypeRouteCall != cr.Type {
			return errors.New("Wrong validation type on RouteCall")
		}
		sig := HashInterface(req)
		if !bytes.Equal(cr.ReqSig, sig) {
			return errors.New("Wrong validation sig on RouteCall")
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

	res, err := gpr.lr.MessageBus.Send(msg)
	if err != nil {
		return errors.Wrap(err, "couldn't dispatch event")
	}

	rep.Result = res.(*reply.CallMethod).Result
	gpr.lr.addObjectCaseRecord(req.Callee, core.CaseRecord{
		Type:   core.CaseRecordTypeRouteCall,
		ReqSig: HashInterface(req),
		Resp:   rep.Result,
	})

	return nil
}

// SaveAsChild is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) SaveAsChild(req rpctypes.UpSaveAsChildReq, rep *rpctypes.UpSaveAsChildResp) error {
	if gpr.lr.MessageBus == nil {
		return errors.New("event bus was not set during initialization")
	}

	cr, step := gpr.lr.getNextValidationStep(req.Callee)
	if step >= 0 { // validate
		if core.CaseRecordTypeSaveAsChild != cr.Type {
			return errors.New("Wrong validation type on SaveAsChild")
		}
		sig := HashInterface(req)
		if !bytes.Equal(cr.ReqSig, sig) {
			return errors.New("Wrong validation sig on SaveAsChild")
		}

		rep.Reference = cr.Resp.(*core.RecordRef)
		return nil
	}

	msg := &message.CallConstructor{
		BaseLogicMessage: MakeBaseMessage(req.UpBaseReq),
		ClassRef:         req.Class,
		ParentRef:        req.Parent,
		Name:             req.ConstructorName,
		Arguments:        req.ArgsSerialized,
		SaveAs:           message.Child,
	}

	res, err := gpr.lr.MessageBus.Send(msg)
	if err != nil {
		return errors.Wrap(err, "couldn't save new object as child")
	}

	rep.Reference = res.(*reply.CallConstructor).Object

	gpr.lr.addObjectCaseRecord(req.Callee, core.CaseRecord{
		Type:   core.CaseRecordTypeSaveAsChild,
		ReqSig: HashInterface(req),
		Resp:   rep.Reference,
	})

	return nil
}

// GetObjChildren is an RPC returns set of object children
func (gpr *RPC) GetObjChildren(req rpctypes.UpGetObjChildrenReq, rep *rpctypes.UpGetObjChildrenResp) error {
	// TODO: INS-408

	cr, step := gpr.lr.getNextValidationStep(req.Callee)
	if step >= 0 { // validate
		if core.CaseRecordTypeGetObjChildren != cr.Type {
			return errors.New("Wrong validation type on GetObjChildren")
		}
		sig := HashInterface(req)
		if !bytes.Equal(cr.ReqSig, sig) {
			return errors.New("Wrong validation sig on GetObjChildren")
		}

		rep.Children = cr.Resp.([]core.RecordRef)
		return nil
	}

	am := gpr.lr.ArtifactManager
	i, err := am.GetChildren(req.Obj, nil)
	if err != nil {
		return err
	}
	for i.HasNext() {
		r, err := i.Next()
		if err != nil {
			return err
		}
		o, err := am.GetObject(*r, nil)
		if err != nil {
			// TODO: we should detect deactivated objects
			continue
		}
		cd, err := o.ClassDescriptor(nil)
		if err != nil {
			return errors.Wrap(err, "Have ref, have no object")
		}
		ref := cd.HeadRef()
		if ref.Equal(req.Class) {
			rep.Children = append(rep.Children, *r)
		}
	}
	gpr.lr.addObjectCaseRecord(req.Callee, core.CaseRecord{ // bad idea, we can store gadzillion of children
		Type:   core.CaseRecordTypeGetObjChildren,
		ReqSig: HashInterface(req),
		Resp:   rep.Children,
	})
	return nil
}

// SaveAsDelegate is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) SaveAsDelegate(req rpctypes.UpSaveAsDelegateReq, rep *rpctypes.UpSaveAsDelegateResp) error {
	cr, step := gpr.lr.getNextValidationStep(req.Callee)
	if step >= 0 { // validate
		if core.CaseRecordTypeSaveAsDelegate != cr.Type {
			return errors.New("Wrong validation type on SaveAsDelegate")
		}
		sig := HashInterface(req)
		if !bytes.Equal(cr.ReqSig, sig) {
			return errors.New("Wrong validation sig on SaveAsDelegate")
		}

		rep.Reference = cr.Resp.(*core.RecordRef)
		return nil
	}

	msg := &message.CallConstructor{
		BaseLogicMessage: MakeBaseMessage(req.UpBaseReq),
		ClassRef:         req.Class,
		ParentRef:        req.Into,
		Name:             req.ConstructorName,
		Arguments:        req.ArgsSerialized,
		SaveAs:           message.Delegate,
	}

	res, err := gpr.lr.MessageBus.Send(msg)

	if err != nil {
		return errors.Wrap(err, "couldn't save new object as delegate")
	}

	rep.Reference = res.(*reply.CallConstructor).Object
	gpr.lr.addObjectCaseRecord(req.Callee, core.CaseRecord{
		Type:   core.CaseRecordTypeSaveAsDelegate,
		ReqSig: HashInterface(req),
		Resp:   rep.Reference,
	})

	return nil
}

// GetDelegate is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) GetDelegate(req rpctypes.UpGetDelegateReq, rep *rpctypes.UpGetDelegateResp) error {
	cr, step := gpr.lr.getNextValidationStep(req.Callee)
	if step >= 0 { // validate
		if core.CaseRecordTypeGetDelegate != cr.Type {
			return errors.New("Wrong validation type on RouteCall")
		}
		sig := HashInterface(req)
		if !bytes.Equal(cr.ReqSig, sig) {
			return errors.New("Wrong validation sig on RouteCall")
		}

		rep.Object = cr.Resp.(core.RecordRef)
		return nil
	}
	am := gpr.lr.ArtifactManager
	ref, err := am.GetDelegate(req.Object, req.OfType)
	if err != nil {
		return err
	}
	rep.Object = *ref
	gpr.lr.addObjectCaseRecord(req.Callee, core.CaseRecord{
		Type:   core.CaseRecordTypeGetDelegate,
		ReqSig: HashInterface(req),
		Resp:   rep.Object,
	})
	return nil
}

// DeactivateObject is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) DeactivateObject(req rpctypes.UpDeactivateObjectReq, rep *rpctypes.UpDeactivateObjectResp) error {
	cr, step := gpr.lr.getNextValidationStep(req.Callee)
	if step >= 0 { // validate
		if core.CaseRecordTypeDeactivateObject != cr.Type {
			return errors.New("Wrong validation type on RouteCall")
		}
		sig := HashInterface(req)
		if !bytes.Equal(cr.ReqSig, sig) {
			return errors.New("Wrong validation sig on RouteCall")
		}
		return nil
	}
	am := gpr.lr.ArtifactManager
	_, err := am.DeactivateObject(core.RecordRef{}, core.RecordRef{}, req.Object)
	if err != nil {
		return err
	}
	gpr.lr.addObjectCaseRecord(req.Callee, core.CaseRecord{
		Type:   core.CaseRecordTypeDeactivateObject,
		ReqSig: HashInterface(req),
	})
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
