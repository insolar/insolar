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
	"net"
	"net/http"
	"net/rpc"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/pkg/errors"
)

var rpcService *RPC

func StartRPC(lr *LogicRunner) *RPC {
	if rpcService == nil {
		rpcService = &RPC{lr: lr}
		err := rpc.Register(rpcService)
		if err != nil {
			panic("Fail to register LogicRunner RPC Service: " + err.Error())
		}
		rpc.HandleHTTP()
	}
	rpcService.lr = lr

	l, e := net.Listen("tcp", lr.Cfg.RPCListen)
	if e != nil {
		log.Fatal("couldn't setup listener on '"+lr.Cfg.RPCListen+"': ", e)
	}
	lr.sock = l
	log.Infof("starting LogicRunner RPC service on %q", lr.Cfg.RPCListen)
	go func() {
		if err := http.Serve(l, nil); err != nil {
			log.Warn("Can't Listen LogicRunner RPC Socket: ", err)
		}
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
	codeDescriptor, err := am.GetCode(req.Code, []core.MachineType{req.MType})
	if err != nil {
		return err
	}
	reply.Code = codeDescriptor.Code()
	return nil
}

// MakeBaseEvent makes base of logicrunner event from base of up request
func MakeBaseEvent(req rpctypes.UpBaseReq) message.BaseLogicEvent {
	return message.BaseLogicEvent{
		Caller: req.Me,
	}
}

// RouteCall routes call from a contract to a contract through event bus.
func (gpr *RPC) RouteCall(req rpctypes.UpRouteReq, rep *rpctypes.UpRouteResp) error {
	if gpr.lr.MessageBus == nil {
		return errors.New("event bus was not set during initialization")
	}

	var mode message.MethodReturnMode
	if req.Wait {
		mode = message.ReturnResult
	} else {
		mode = message.ReturnNoWait
	}

	msg := &message.CallMethod{
		BaseLogicEvent: MakeBaseEvent(req.UpBaseReq),
		ReturnMode:     mode,
		ObjectRef:      req.Object,
		Method:         req.Method,
		Arguments:      req.Arguments,
	}

	res, err := gpr.lr.MessageBus.Send(msg)
	if err != nil {
		return errors.Wrap(err, "couldn't dispatch event")
	}

	gpr.lr.addObjectCaseRecord(req.Me, CaseRecord{
		Type:   CaseRecordTypeRouteCall,
		ReqSig: HashInterface(req),
		Resp:   rep,
	})
	rep.Result = res.(*reply.Common).Result

	return nil
}

// RouteConstructorCall routes call from a contract to a constructor of another contract
func (gpr *RPC) RouteConstructorCall(req rpctypes.UpRouteConstructorReq, rep *rpctypes.UpRouteConstructorResp) error {
	if gpr.lr.MessageBus == nil {
		return errors.New("event bus was not set during initialization")
	}

	msg := &message.CallConstructor{
		BaseLogicEvent: MakeBaseEvent(req.UpBaseReq),
		ClassRef:       req.Reference,
		Name:           req.Constructor,
		Arguments:      req.Arguments,
	}

	res, err := gpr.lr.MessageBus.Send(msg)
	if err != nil {
		return errors.Wrap(err, "couldn't dispatch event")
	}
	gpr.lr.addObjectCaseRecord(req.Me, CaseRecord{
		Type:   CaseRecordTypeRouteCall,
		ReqSig: HashInterface(req),
		Resp:   rep,
	})
	rep.Data = res.(*reply.Common).Data
	return nil
}

// SaveAsChild is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) SaveAsChild(req rpctypes.UpSaveAsChildReq, reply *rpctypes.UpSaveAsChildResp) error {
	ref, err := gpr.lr.ArtifactManager.ActivateObject(
		core.RecordRef{}, core.RandomRef(), req.Class, req.Parent, req.Data,
	)
	if err != nil {
		return errors.Wrap(err, "couldn't save new object")
	}
	gpr.lr.addObjectCaseRecord(req.Me, CaseRecord{
		Type:   CaseRecordTypeSaveAsChild,
		ReqSig: HashInterface(req),
		Resp:   reply,
	})
	reply.Reference = *ref
	return nil
}

// GetObjChildren is an RPC returns set of object children
func (gpr *RPC) GetObjChildren(req rpctypes.UpGetObjChildrenReq, reply *rpctypes.UpGetObjChildrenResp) error {
	// TODO: INS-408
	am := gpr.lr.ArtifactManager
	obj, err := am.GetObject(req.Obj, nil)
	if err != nil {
		return errors.Wrap(err, "am.GetObject failed")
	}
	i := obj.Children()
	for i.HasNext() {
		r, err := i.Next()
		if err != nil {
			return err
		}
		o, err := am.GetObject(r, nil)
		if err != nil {
			return errors.Wrap(err, "Have ref, have no object")
		}
		cd, err := o.ClassDescriptor(nil)
		if err != nil {
			return errors.Wrap(err, "Have ref, have no object")
		}
		ref := cd.HeadRef()
		if ref.Equal(req.Class) {
			reply.Children = append(reply.Children, r)
		}
	}
	gpr.lr.addObjectCaseRecord(req.Me, CaseRecord{ // bad idea, we can store gadzillion of children
		Type:   CaseRecordTypeGetObjChildren,
		ReqSig: HashInterface(req),
		Resp:   reply,
	})
	return nil
}

// SaveAsDelegate is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) SaveAsDelegate(req rpctypes.UpSaveAsDelegateReq, reply *rpctypes.UpSaveAsDelegateResp) error {
	ref, err := gpr.lr.ArtifactManager.ActivateObjectDelegate(
		core.RecordRef{}, core.RandomRef(), req.Class, req.Into, req.Data,
	)
	if err != nil {
		return errors.Wrap(err, "couldn't save delegate")
	}
	gpr.lr.addObjectCaseRecord(req.Me, CaseRecord{
		Type:   CaseRecordTypeSaveAsDelegate,
		ReqSig: HashInterface(req),
		Resp:   reply,
	})
	reply.Reference = *ref
	return nil
}

// GetDelegate is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) GetDelegate(req rpctypes.UpGetDelegateReq, reply *rpctypes.UpGetDelegateResp) error {
	am := gpr.lr.ArtifactManager
	ref, err := am.GetDelegate(req.Object, req.OfType)
	if err != nil {
		return err
	}
	gpr.lr.addObjectCaseRecord(req.Me, CaseRecord{
		Type:   CaseRecordTypeGetDelegate,
		ReqSig: HashInterface(req),
		Resp:   reply,
	})
	reply.Object = *ref
	return nil
}
