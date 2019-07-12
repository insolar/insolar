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

// todo it may use foundation.Context
// DownCallMethodReq is a set of arguments for CallMethod RPC in the runner
type DownCallMethodReq struct {
	Context   *insolar.LogicCallContext
	Code      insolar.Reference
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
	Code      insolar.Reference
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
	Mode            insolar.CallMode
	Callee          insolar.Reference
	CalleePrototype insolar.Reference
	Request         insolar.Reference
}

// UpRespIface interface for UpBaseReq descendant responses
type UpRespIface interface{}

// UpGetCodeReq is a set of arguments for GetCode RPC in goplugin
type UpGetCodeReq struct {
	UpBaseReq
	MType insolar.MachineType
	Code  insolar.Reference
}

// UpGetCodeResp is response from GetCode RPC in goplugin
type UpGetCodeResp struct {
	Code []byte
}

// UpRouteReq is a set of arguments for Send RPC in goplugin
type UpRouteReq struct {
	UpBaseReq
	Wait      bool
	Immutable bool
	Saga      bool
	Object    insolar.Reference
	Method    string
	Arguments insolar.Arguments
	Prototype insolar.Reference
}

// UpRouteResp is response from Send RPC in goplugin
type UpRouteResp struct {
	Result insolar.Arguments
}

// UpSaveAsChildReq is a set of arguments for SaveAsChild RPC in goplugin
type UpSaveAsChildReq struct {
	UpBaseReq
	Parent          insolar.Reference
	Prototype       insolar.Reference
	ConstructorName string
	ArgsSerialized  []byte
}

// UpSaveAsChildResp is a set of arguments for SaveAsChild RPC in goplugin
type UpSaveAsChildResp struct {
	Reference *insolar.Reference
}

// UpGetObjChildrenIteratorReq is a set of arguments for GetObjChildrenIterator RPC in goplugin
type UpGetObjChildrenIteratorReq struct {
	UpBaseReq
	IteratorID string
	Object     insolar.Reference
	Prototype  insolar.Reference
}

// UpGetObjChildrenIteratorResp is response from GetObjChildren RPC in goplugin
type UpGetObjChildrenIteratorResp struct {
	Iterator ChildIterator
}

// ChildIterator hold an iterator data of GetObjChildrenIterator method
type ChildIterator struct {
	ID       string
	Buff     []insolar.Reference
	CanFetch bool
}

// UpSaveAsDelegateReq is a set of arguments for SaveAsDelegate RPC in goplugin
type UpSaveAsDelegateReq struct {
	UpBaseReq
	Into            insolar.Reference
	Prototype       insolar.Reference
	ConstructorName string
	ArgsSerialized  []byte
}

// UpSaveAsDelegateResp is response from SaveAsDelegate RPC in goplugin
type UpSaveAsDelegateResp struct {
	Reference *insolar.Reference
}

// UpGetDelegateReq is a set of arguments for GetDelegate RPC in goplugin
type UpGetDelegateReq struct {
	UpBaseReq
	Object insolar.Reference
	OfType insolar.Reference
}

// UpGetDelegateResp is response from GetDelegate RPC in goplugin
type UpGetDelegateResp struct {
	Object insolar.Reference
}

// UpDeactivateObjectReq is a set of arguments for DeactivateObject RPC in goplugin
type UpDeactivateObjectReq struct {
	UpBaseReq
}

// UpDeactivateObjectResp is response from DeactivateObject RPC in goplugin
type UpDeactivateObjectResp struct {
}
