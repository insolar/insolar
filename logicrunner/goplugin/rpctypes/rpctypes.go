// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
	Data []byte
	Ret  insolar.Arguments
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
	Result insolar.Arguments
}

// UpDeactivateObjectReq is a set of arguments for DeactivateObject RPC in goplugin
type UpDeactivateObjectReq struct {
	UpBaseReq
}

// UpDeactivateObjectResp is response from DeactivateObject RPC in goplugin
type UpDeactivateObjectResp struct {
}
