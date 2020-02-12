// Copyright 2020 Insolar Network Ltd.
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

package builtin

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	lrCommon "github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

type ProxyHelper struct {
	lrCommon.Serializer
	lrCommon.SystemError
	methods lrCommon.LogicRunnerRPCStub
}

func NewProxyHelper(runner lrCommon.LogicRunnerRPCStub) *ProxyHelper {
	return &ProxyHelper{
		Serializer:  lrCommon.NewCBORSerializer(),
		SystemError: lrCommon.NewSystemError(),
		methods:     runner,
	}
}

func (h *ProxyHelper) getUpBaseReq() rpctypes.UpBaseReq {
	callContext := foundation.GetLogicalContext()

	return rpctypes.UpBaseReq{
		Mode:            callContext.Mode,
		Callee:          *callContext.Callee,
		CalleePrototype: *callContext.CallerPrototype,
		Request:         *callContext.Request,
	}
}

func (h *ProxyHelper) RouteCall(ref insolar.Reference, immutable bool, saga bool, method string, args []byte,
	proxyPrototype insolar.Reference) ([]byte, error) {

	if h.GetSystemError() != nil {
		return nil, h.GetSystemError()
	}

	res := rpctypes.UpRouteResp{}
	req := rpctypes.UpRouteReq{
		UpBaseReq: h.getUpBaseReq(),

		Object:    ref,
		Immutable: immutable,
		Saga:      saga,
		Method:    method,
		Arguments: args,
		Prototype: proxyPrototype,
	}

	err := h.methods.RouteCall(req, &res)
	if err != nil {
		h.SetSystemError(err)
		return nil, err
	}

	return res.Result, nil
}

func (h *ProxyHelper) SaveAsChild(
	parentRef, classRef insolar.Reference,
	constructorName string, argsSerialized []byte,
) (
	[]byte, error,
) {
	if !parentRef.IsObjectReference() {
		return nil, errors.Errorf("Failed to save AsChild: objRef should be ObjectReference; ref=%s", parentRef.String())
	}

	if h.GetSystemError() != nil {
		// There was a system error during execution of the contract.
		// Immediately return this error to the calling contract - any
		// results will not be registered on LME anyway.
		return nil, h.GetSystemError()
	}

	res := rpctypes.UpSaveAsChildResp{}
	req := rpctypes.UpSaveAsChildReq{
		UpBaseReq: h.getUpBaseReq(),

		Parent:          parentRef,
		Prototype:       classRef,
		ConstructorName: constructorName,
		ArgsSerialized:  argsSerialized,
	}

	err := h.methods.SaveAsChild(req, &res)
	if err != nil {
		h.SetSystemError(err)
		return nil, err
	}

	return res.Result, nil
}

func (h *ProxyHelper) DeactivateObject(object insolar.Reference) error {
	if h.GetSystemError() != nil {
		return h.GetSystemError()
	}

	res := rpctypes.UpDeactivateObjectResp{}
	req := rpctypes.UpDeactivateObjectReq{
		UpBaseReq: h.getUpBaseReq(),
	}

	if err := h.methods.DeactivateObject(req, &res); err != nil {
		h.SetSystemError(err)
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
