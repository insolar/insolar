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

package rpctypes

import (
	"github.com/insolar/insolar/insolar"
)

// Types for RPC requests and responses between goplugin and goinsider.
// Calls from goplugin to goinsider go "downwards" and names are
// prefixed with "Down". Reverse calls go "upwards", so "Up" prefix

// DownCallMethodReq is a set of arguments for CallMethod RPC in the runner
type DownCallMethodReq struct { // todo it may use foundation.Context
	Context   *insolar.LogicCallContext
	Code      insolar.RecordRef
	Data      []byte
	Method    string
	Arguments insolar.Arguments
}

// DownCallMethodResp is response from CallMethod RPC in the runner
type DownCallMethodResp struct {
	Data []byte
	Ret  insolar.Arguments
}

// DownCallConstructorReq is a set of arguments for CallConstructor RPC
// in the runner
type DownCallConstructorReq struct {
	Code      insolar.RecordRef
	Name      string
	Arguments insolar.Arguments
	Context   *insolar.LogicCallContext
}

// DownCallConstructorResp is response from CallConstructor RPC in the runner
type DownCallConstructorResp struct {
	Ret insolar.Arguments
}

// UpBaseReq  is a base type for all insgorund -> logicrunner requests
type UpBaseReq struct {
	Mode      string
	Callee    insolar.RecordRef
	Prototype insolar.RecordRef
	Request   insolar.RecordRef
}

// UpRespIface interface for UpBaseReq descendant responses
type UpRespIface interface{}

// UpGetCodeReq is a set of arguments for GetCode RPC in goplugin
type UpGetCodeReq struct {
	UpBaseReq
	MType insolar.MachineType
	Code  insolar.RecordRef
}

// UpGetCodeResp is response from GetCode RPC in goplugin
type UpGetCodeResp struct {
	Code []byte
}

// UpRouteReq is a set of arguments for Send RPC in goplugin
type UpRouteReq struct {
	UpBaseReq
	Wait           bool
	Object         insolar.RecordRef
	Method         string
	Arguments      insolar.Arguments
	ProxyPrototype insolar.RecordRef
}

// UpRouteResp is response from Send RPC in goplugin
type UpRouteResp struct {
	Result insolar.Arguments
}

// UpSaveAsChildReq is a set of arguments for SaveAsChild RPC in goplugin
type UpSaveAsChildReq struct {
	UpBaseReq
	Parent          insolar.RecordRef
	Prototype       insolar.RecordRef
	ConstructorName string
	ArgsSerialized  []byte
}

// UpSaveAsChildResp is a set of arguments for SaveAsChild RPC in goplugin
type UpSaveAsChildResp struct {
	Reference *insolar.RecordRef
}

// UpGetObjChildrenIteratorReq is a set of arguments for GetObjChildrenIterator RPC in goplugin
type UpGetObjChildrenIteratorReq struct {
	UpBaseReq
	IteratorID string
	Obj        insolar.RecordRef
	Prototype  insolar.RecordRef
}

// UpGetObjChildrenIteratorResp is response from GetObjChildren RPC in goplugin
type UpGetObjChildrenIteratorResp struct {
	Iterator ChildIterator
}

// ChildIterator hold an iterator data of GetObjChildrenIterator method
type ChildIterator struct {
	ID       string
	Buff     []insolar.RecordRef
	CanFetch bool
}

// UpSaveAsDelegateReq is a set of arguments for SaveAsDelegate RPC in goplugin
type UpSaveAsDelegateReq struct {
	UpBaseReq
	Into            insolar.RecordRef
	Prototype       insolar.RecordRef
	ConstructorName string
	ArgsSerialized  []byte
}

// UpSaveAsDelegateResp is response from SaveAsDelegate RPC in goplugin
type UpSaveAsDelegateResp struct {
	Reference *insolar.RecordRef
}

// UpGetDelegateReq is a set of arguments for GetDelegate RPC in goplugin
type UpGetDelegateReq struct {
	UpBaseReq
	Object insolar.RecordRef
	OfType insolar.RecordRef
}

// UpGetDelegateResp is response from GetDelegate RPC in goplugin
type UpGetDelegateResp struct {
	Object insolar.RecordRef
}

// UpDeactivateObjectReq is a set of arguments for DeactivateObject RPC in goplugin
type UpDeactivateObjectReq struct {
	UpBaseReq
}

// UpDeactivateObjectResp is response from DeactivateObject RPC in goplugin
type UpDeactivateObjectResp struct {
}
