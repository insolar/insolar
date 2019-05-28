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

func (h *ProxyHelper) RouteCall(ref insolar.Reference, wait bool, immutable bool, method string, args []byte,
	proxyPrototype insolar.Reference) ([]byte, error) {

	res := rpctypes.UpRouteResp{}
	req := rpctypes.UpRouteReq{
		UpBaseReq: h.getUpBaseReq(),

		Object:    ref,
		Wait:      wait,
		Immutable: immutable,
		Method:    method,
		Arguments: args,
		Prototype: proxyPrototype,
	}

	err := h.methods.RouteCall(req, &res)

	if err != nil {
		return nil, err
	}
	return res.Result, nil
}

func (h *ProxyHelper) SaveAsChild(parentRef, classRef insolar.Reference, constructorName string,
	argsSerialized []byte) (insolar.Reference, error) {

	res := rpctypes.UpSaveAsChildResp{}
	req := rpctypes.UpSaveAsChildReq{
		UpBaseReq: h.getUpBaseReq(),

		Parent:          parentRef,
		Prototype:       classRef,
		ConstructorName: constructorName,
		ArgsSerialized:  argsSerialized,
	}

	if err := h.methods.SaveAsChild(req, &res); err != nil {
		return insolar.Reference{}, err
	}
	if res.Reference == nil {
		return insolar.Reference{}, errors.New("Unexpected result, empty reference")
	}
	return *res.Reference, nil
}

func (h *ProxyHelper) GetObjChildrenIterator(head insolar.Reference, prototype insolar.Reference,
	iteratorID string) (*lrCommon.ChildrenTypedIterator, error) {

	res := rpctypes.UpGetObjChildrenIteratorResp{}
	req := rpctypes.UpGetObjChildrenIteratorReq{
		UpBaseReq: h.getUpBaseReq(),

		IteratorID: iteratorID,
		Object:     head,
		Prototype:  prototype,
	}

	err := h.methods.GetObjChildrenIterator(req, &res)

	if err != nil {
		return &lrCommon.ChildrenTypedIterator{}, errors.Wrap(err, "on calling GetObjChildren")
	}
	return &lrCommon.ChildrenTypedIterator{
		Parent:         head,
		ChildPrototype: prototype,
		IteratorID:     res.Iterator.ID,
		Buff:           res.Iterator.Buff,
		CanFetch:       res.Iterator.CanFetch,
	}, nil
}

func (h *ProxyHelper) SaveAsDelegate(parentRef, classRef insolar.Reference, constructorName string,
	argsSerialized []byte) (insolar.Reference, error) {

	res := rpctypes.UpSaveAsDelegateResp{}
	req := rpctypes.UpSaveAsDelegateReq{
		UpBaseReq: h.getUpBaseReq(),

		Into:            parentRef,
		Prototype:       classRef,
		ConstructorName: constructorName,
		ArgsSerialized:  argsSerialized,
	}

	if err := h.methods.SaveAsDelegate(req, &res); err != nil {
		return insolar.Reference{}, err
	}
	if res.Reference == nil {
		return insolar.Reference{}, errors.New("Unexpected result, empty reference")
	}
	return *res.Reference, nil

}

func (h *ProxyHelper) GetDelegate(object, ofType insolar.Reference) (insolar.Reference, error) {
	res := rpctypes.UpGetDelegateResp{}
	req := rpctypes.UpGetDelegateReq{
		UpBaseReq: h.getUpBaseReq(),

		Object: object,
		OfType: ofType,
	}

	if err := h.methods.GetDelegate(req, &res); err != nil {
		return insolar.Reference{}, err
	}
	return res.Object, nil
}

func (h *ProxyHelper) DeactivateObject(object insolar.Reference) error {
	res := rpctypes.UpDeactivateObjectResp{}
	req := rpctypes.UpDeactivateObjectReq{
		UpBaseReq: h.getUpBaseReq(),
	}

	if err := h.methods.DeactivateObject(req, &res); err != nil {
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
