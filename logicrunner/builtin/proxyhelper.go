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

package builtin

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/tylerb/gls"

	"github.com/insolar/insolar/insolar"
	lrCommon "github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

const (
	glsCallContextKey = "callCtx"
	glsSystemErrorKey = "systemError"
)

type ProxyHelper struct {
	lrCommon.Serializer
	methods lrCommon.LogicRunnerRPCStub
}

func NewProxyHelper(runner lrCommon.LogicRunnerRPCStub) *ProxyHelper {
	return &ProxyHelper{
		Serializer: lrCommon.NewCBORSerializer(),
		methods:    runner,
	}
}

func (h *ProxyHelper) getSystemError() error {
	// SystemError means an error in the system (platform), not a particular contract.
	// For instance, timed out external call or failed deserialization means a SystemError.
	// In case of SystemError all following external calls during current method call return
	// an error and the result of the current method call is discarded (not registered).
	callContextInterface := gls.Get(glsSystemErrorKey)
	if callContextInterface == nil {
		return nil
	}
	return callContextInterface.(error)
}

func (h *ProxyHelper) setSystemError(err error) {
	gls.Set(glsSystemErrorKey, err)
}

func (h *ProxyHelper) getUpBaseReq() rpctypes.UpBaseReq {
	callContextInterface := gls.Get(glsCallContextKey)
	if callContextInterface == nil {
		panic("Failed to find call context")
	}
	callContext, ok := callContextInterface.(*insolar.LogicCallContext)
	if !ok {
		panic("Unknown value stored in '" + glsCallContextKey + "'")
	}

	return rpctypes.UpBaseReq{
		Mode:            callContext.Mode,
		Callee:          *callContext.Callee,
		CalleePrototype: *callContext.CallerPrototype,
		Request:         *callContext.Request,
	}
}

// CleanupSystemError should be called in a contract wrapper before actually
// executing a method.
func (h *ProxyHelper) CleanupSystemError() {
	h.setSystemError(nil)
}

// SystemError() should be checked in contract wrapper before returning a result.
// If system error occurred the result returned by the method should be discarded.
func (h *ProxyHelper) SystemError() error {
	return h.getSystemError()
}

func (h *ProxyHelper) RouteCall(ref insolar.Reference, wait bool, immutable bool, saga bool, method string, args []byte,
	proxyPrototype insolar.Reference) ([]byte, error) {

	sysErr := h.getSystemError()
	if sysErr != nil {
		return nil, sysErr
	}

	res := rpctypes.UpRouteResp{}
	req := rpctypes.UpRouteReq{
		UpBaseReq: h.getUpBaseReq(),

		Object:    ref,
		Wait:      wait,
		Immutable: immutable,
		Saga:      saga,
		Method:    method,
		Arguments: args,
		Prototype: proxyPrototype,
	}

	err := h.methods.RouteCall(req, &res)

	if err != nil {
		h.setSystemError(err)
		return nil, err
	}
	return res.Result, nil
}

func (h *ProxyHelper) SaveAsChild(parentRef, classRef insolar.Reference, constructorName string,
	argsSerialized []byte) (insolar.Reference, error) {

	sysErr := h.getSystemError()
	if sysErr != nil {
		return insolar.Reference{}, sysErr
	}

	res := rpctypes.UpSaveAsChildResp{}
	req := rpctypes.UpSaveAsChildReq{
		UpBaseReq: h.getUpBaseReq(),

		Parent:          parentRef,
		Prototype:       classRef,
		ConstructorName: constructorName,
		ArgsSerialized:  argsSerialized,
	}

	if err := h.methods.SaveAsChild(req, &res); err != nil {
		h.setSystemError(err)
		return insolar.Reference{}, err
	}
	if res.Reference == nil {
		err := errors.New("Unexpected result, empty reference")
		h.setSystemError(err)
		return insolar.Reference{}, err
	}
	return *res.Reference, nil
}

func (h *ProxyHelper) SaveAsDelegate(parentRef, classRef insolar.Reference, constructorName string,
	argsSerialized []byte) (insolar.Reference, error) {

	sysErr := h.getSystemError()
	if sysErr != nil {
		return insolar.Reference{}, sysErr
	}

	res := rpctypes.UpSaveAsDelegateResp{}
	req := rpctypes.UpSaveAsDelegateReq{
		UpBaseReq: h.getUpBaseReq(),

		Into:            parentRef,
		Prototype:       classRef,
		ConstructorName: constructorName,
		ArgsSerialized:  argsSerialized,
	}

	if err := h.methods.SaveAsDelegate(req, &res); err != nil {
		h.setSystemError(err)
		return insolar.Reference{}, err
	}
	if res.Reference == nil {
		err := errors.New("Unexpected result, empty reference")
		h.setSystemError(err)
		return insolar.Reference{}, err
	}
	return *res.Reference, nil

}

func (h *ProxyHelper) GetDelegate(object, ofType insolar.Reference) (insolar.Reference, error) {
	sysErr := h.getSystemError()
	if sysErr != nil {
		return insolar.Reference{}, sysErr
	}

	res := rpctypes.UpGetDelegateResp{}
	req := rpctypes.UpGetDelegateReq{
		UpBaseReq: h.getUpBaseReq(),

		Object: object,
		OfType: ofType,
	}

	if err := h.methods.GetDelegate(req, &res); err != nil {
		h.setSystemError(err)
		return insolar.Reference{}, err
	}
	return res.Object, nil
}

func (h *ProxyHelper) DeactivateObject(object insolar.Reference) error {
	sysErr := h.getSystemError()
	if sysErr != nil {
		return sysErr
	}

	res := rpctypes.UpDeactivateObjectResp{}
	req := rpctypes.UpDeactivateObjectReq{
		UpBaseReq: h.getUpBaseReq(),
	}

	if err := h.methods.DeactivateObject(req, &res); err != nil {
		h.setSystemError(err)
		return err
	}
	return nil
}

/*
func (h *ProxyHelper) Serialize(what interface{}, to *[]byte) error {
	panic("implement me")
}

func (h *ProxyHelper) Deserialize(from []byte, into interface{}) error {
	panic("implement me")
}
*/

func (h *ProxyHelper) MakeErrorSerializable(err error) error {
	if err == nil || err == (*foundation.Error)(nil) || reflect.ValueOf(err).IsNil() {
		return nil
	}
	return &foundation.Error{S: err.Error()}
}
