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
	"fmt"
	"net"
	"net/rpc"
	"sync"

	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
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

func recoverRPC(err *error) {
	if r := recover(); r != nil {
		if err != nil {
			if *err == nil {
				*err = errors.New(fmt.Sprint(r))
			} else {
				*err = errors.New(fmt.Sprint(*err, r))
			}
		}
	}
}

// GetCode is an RPC retrieving a code by its reference
func (gpr *RPC) GetCode(req rpctypes.UpGetCodeReq, reply *rpctypes.UpGetCodeResp) (err error) {
	defer recoverRPC(&err)
	os := gpr.lr.MustObjectState(req.Callee)
	es := os.MustModeState(req.Mode)
	ctx := es.Current.Context
	inslogger.FromContext(ctx).Debug("In RPC.GetCode ....")

	am := gpr.lr.ArtifactManager

	ctx, span := instracer.StartSpan(ctx, "service.GetCode")
	defer span.End()

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

// MakeBaseMessage makes base of logicrunner event from base of up request
func MakeBaseMessage(req rpctypes.UpBaseReq, es *ExecutionState) message.BaseLogicMessage {
	es.nonce++
	return message.BaseLogicMessage{
		Caller:          req.Callee,
		CallerPrototype: req.Prototype,
		Request:         req.Request,
		Nonce:           es.nonce,
	}
}

// RouteCall routes call from a contract to a contract through event bus.
func (gpr *RPC) RouteCall(req rpctypes.UpRouteReq, rep *rpctypes.UpRouteResp) (err error) {
	defer recoverRPC(&err)

	os := gpr.lr.MustObjectState(req.Callee)

	if os.ExecutionState.Current.LogicContext.Immutable {
		return errors.New("Try to call route from immutable method")
	}

	es := os.MustModeState(req.Mode)
	ctx := es.Current.Context

	// TODO: delegation token

	bm := MakeBaseMessage(req.UpBaseReq, es)
	res, err := gpr.lr.ContractRequester.CallMethod(ctx,
		&bm,
		!req.Wait,
		req.Immutable,
		&req.Object,
		req.Method,
		req.Arguments,
		&req.ProxyPrototype,
	)
	if err != nil {
		return err
	}

	if req.Wait {
		rep.Result = res.(*reply.CallMethod).Result
	}

	return nil
}

// SaveAsChild is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) SaveAsChild(req rpctypes.UpSaveAsChildReq, rep *rpctypes.UpSaveAsChildResp) (err error) {
	defer recoverRPC(&err)

	os := gpr.lr.MustObjectState(req.Callee)
	es := os.MustModeState(req.Mode)
	ctx := es.Current.Context

	bm := MakeBaseMessage(req.UpBaseReq, es)
	ref, err := gpr.lr.ContractRequester.CallConstructor(ctx, &bm, false, &req.Prototype, &req.Parent, req.ConstructorName, req.ArgsSerialized, int(message.Child))

	rep.Reference = ref

	return err
}

// SaveAsDelegate is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) SaveAsDelegate(req rpctypes.UpSaveAsDelegateReq, rep *rpctypes.UpSaveAsDelegateResp) (err error) {
	defer recoverRPC(&err)

	os := gpr.lr.MustObjectState(req.Callee)
	es := os.MustModeState(req.Mode)
	ctx := es.Current.Context

	bm := MakeBaseMessage(req.UpBaseReq, es)
	ref, err := gpr.lr.ContractRequester.CallConstructor(ctx, &bm, false, &req.Prototype, &req.Into, req.ConstructorName, req.ArgsSerialized, int(message.Delegate))

	rep.Reference = ref
	return err
}

var iteratorMap = make(map[string]artifacts.RefIterator)
var iteratorMapLock = sync.RWMutex{}
var iteratorBuffSize = 1000

// GetObjChildrenIterator is an RPC returns an iterator over object children with specified prototype
func (gpr *RPC) GetObjChildrenIterator(
	req rpctypes.UpGetObjChildrenIteratorReq,
	rep *rpctypes.UpGetObjChildrenIteratorResp,
) (
	err error,
) {
	defer recoverRPC(&err)

	os := gpr.lr.MustObjectState(req.Callee)
	es := os.MustModeState(req.Mode)
	ctx := es.Current.Context

	am := gpr.lr.ArtifactManager
	iteratorID := req.IteratorID

	iteratorMapLock.RLock()
	iterator, ok := iteratorMap[iteratorID]
	iteratorMapLock.RUnlock()

	if !ok {
		newIterator, err := am.GetChildren(ctx, req.Obj, nil)
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

		o, err := am.GetObject(ctx, *r, nil, false)

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
func (gpr *RPC) GetDelegate(req rpctypes.UpGetDelegateReq, rep *rpctypes.UpGetDelegateResp) (err error) {
	defer recoverRPC(&err)

	os := gpr.lr.MustObjectState(req.Callee)
	es := os.MustModeState(req.Mode)
	ctx := es.Current.Context

	am := gpr.lr.ArtifactManager
	ref, err := am.GetDelegate(ctx, req.Object, req.OfType)
	if err != nil {
		return err
	}
	rep.Object = *ref
	return nil
}

// DeactivateObject is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) DeactivateObject(req rpctypes.UpDeactivateObjectReq, rep *rpctypes.UpDeactivateObjectResp) (err error) {
	defer recoverRPC(&err)

	os := gpr.lr.MustObjectState(req.Callee)
	es := os.MustModeState(req.Mode)
	es.deactivate = true
	return nil
}
